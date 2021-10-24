// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ratecount

import (
	"fmt"
	"testing"
	"time"
)

func ExampleCounter() {
	c := NewAvgCounter(100*time.Millisecond, 3)

	for i := 0; i < 51; i++ {
		c.Incr(1)
		time.Sleep(10 * time.Millisecond)
	}

	// should be 9 or 10
	fmt.Println(c.Rate() == 9 || c.Rate() == 10)
	// output: true
}

func BenchmarkIncrSmall(b *testing.B) {
	c := NewAvgCounter(time.Millisecond, 5)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Incr(1)
	}
}

func BenchmarkIncrMid(b *testing.B) {
	c := NewAvgCounter(500*time.Millisecond, 5)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Incr(1)
	}
}

func BenchmarkIncrLarge(b *testing.B) {
	c := NewAvgCounter(time.Minute, 5)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Incr(1)
	}
}
