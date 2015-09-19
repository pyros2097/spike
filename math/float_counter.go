// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

import (
	"math"
)

// Track properties of a stream of float values. The properties (total value, minimum, etc) are updated as values are
// put(float) into the stream.
type FloatCounter struct {
	// The amount of values adde
	count int
	// The sum of all value
	total float32
	// The smallest valu
	min float32
	// The largest valu
	max float32
	// The average value (total / count
	average float32
	// The latest raw valu
	latest float32
	// The current windowed mean valu
	value float32
	// Provides access to the WindowedMean if any (can be null
	mean *WindowedMean
}

// Construct a new FloatCounter
// windowSize The size of the mean window or 1 or below to not use a windowed mean
func NewFloatCounter(windowSize int) {
	if windowSize > 1 {
		self.mean = NewWindowedMean(windowSize)
	}
	self.Reset()
}

// Add a value and update all fields.
// value The value to add
func (self *FloatCounter) Put(value float32) {
	self.latest = value
	self.total += value
	self.count++
	self.average = self.total / self.count

	if self.mean != nil {
		self.mean.AddValue(value)
		self.value = self.mean.GetMean()
	} else {
		self.value = self.latest
	}

	if self.mean == nil || self.mean.HasEnoughData() {
		if self.value < min {
			min = self.value
		}
		if self.value > max {
			max = self.value
		}
	}
}

// Reset all values to their default value
func (self *FloatCounter) Reset() {
	self.count = 0
	self.total = 0
	self.min = math.MaxFloat32
	self.max = math.SmallestNonzeroFloat32
	self.average = 0
	self.latest = 0
	self.value = 0
	if self.mean != nil {
		self.mean.Clear()
	}
}
