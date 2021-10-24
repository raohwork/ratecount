// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package ratecount provides thread-safe rate calculator and few helpers
package ratecount

import (
	"sync"
	"time"
)

// common constants used with *Counter.In()
const (
	KB  = 1e3
	MB  = 1e6
	GB  = 1e9
	TB  = 1e12
	KiB = 1024
	MiB = 1024 * 1024
	GiB = 1024 * 1024 * 1024
	TiB = 1024 * 1024 * 1024 * 1024
)

// NewAvgCounter creates a rate counter computes average of last n time windows.
//
// IT RETURNS nil IF res == 0 || n == 0
//
// It calculates how many "tokens" (in most cases, bytes) is transfered in given
// time window res, caches last n windows, and calculates average transfer rate. As
// it calculates AVERAGE OF LAST N windows, you'll get lower rate within first n
// windows.
//
// For example, NewAvgCounter(time.Second, 5) calculates transfer speed (token per
// second) for last 5 seconds. NewAvgCounter(10*time.Second, 5) calculates transfer
// speed (token per 10 seconds) for last 50 seconds.
//
// To be more specific, say you have a precise ticker, which calls Incr(1) at
// t = 0.5 / 1.5 / 2.5 / ... (1 tick per second), so calling Rate() at t=6 returns
// average speed of t=1~5, and calling Rate() at t=3 returns result from t=-2~2, or
// (t1+t2) / 5.
func NewAvgCounter(res time.Duration, n uint) *Counter {
	if res == 0 || n == 0 {
		return nil
	}
	ret := &Counter{
		values:     make([]int64, n+1),
		length:     n,
		resolution: res,
	}

	return ret
}

// NewCounter is shortcut of NewAvgCounter(res, 1)
func NewCounter(res time.Duration) *Counter {
	return NewAvgCounter(res, 1)
}

// Counter implements a thread-safe rate counter, See NewAvgCounter for detail.
type Counter struct {
	lock       sync.Mutex
	values     []int64
	position   int64
	length     uint
	resolution time.Duration
}

func (c *Counter) swapAndLock() {
	c.lock.Lock()

	pos := time.Now().UnixNano() / c.resolution.Nanoseconds()
	if c.position == pos {
		return
	}
	c.swap(pos)
}

func (c *Counter) swap(pos int64) {
	delta := uint(pos - c.position)
	if delta > c.length {
		for i := uint(0); i <= c.length; i++ {
			c.values[i] = 0
		}
		c.position = pos
		return
	}

	for i := delta; i <= c.length; i++ {
		c.values[i-delta] = c.values[i]
	}
	for i := uint(0); i < delta; i++ {
		c.values[c.length-i] = 0
	}
	c.values[c.length] = 0
	c.position = pos
}

// Rate retrieves average rate
func (c *Counter) Rate() int64 {
	c.swapAndLock()
	defer c.lock.Unlock()

	sum := int64(0)
	for i := uint(0); i < c.length; i++ {
		sum += c.values[i]
	}
	return sum / int64(c.length)
}

// In returns Rate()/i. Use it like a pro: fmt.Printf("%d kb", c.In(KB))
func (c *Counter) In(i int64) int64 {
	return c.Rate() / i
}

// Incr increases counter in current window
func (c *Counter) Incr(n int64) {
	c.swapAndLock()
	defer c.lock.Unlock()

	c.values[c.length] += n
}
