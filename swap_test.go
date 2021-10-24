// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ratecount

import (
	"reflect"
	"testing"
	"time"
)

func TestSwapLarge(t *testing.T) {
	c := NewAvgCounter(time.Second, 3)
	c.values[c.length] = 100
	c.swap(100)
	expect := []int64{0, 0, 0, 0}
	if !reflect.DeepEqual(c.values, expect) {
		t.Fatal("unexpected result: ", c.values)
	}
}

func TestSwapNormal(t *testing.T) {
	c := NewAvgCounter(time.Second, 3)
	c.values[c.length] = 100
	c.position = 99
	c.swap(100)
	expect := []int64{0, 0, 100, 0}
	if !reflect.DeepEqual(c.values, expect) {
		t.Fatal("unexpected result: ", c.values)
	}
}

func TestSwapSkipSome(t *testing.T) {
	c := NewAvgCounter(time.Second, 3)
	c.values[c.length] = 100
	c.position = 98
	c.swap(100)
	expect := []int64{0, 100, 0, 0}
	if !reflect.DeepEqual(c.values, expect) {
		t.Fatal("unexpected result: ", c.values)
	}
}
