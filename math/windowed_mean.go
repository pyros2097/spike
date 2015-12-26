// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

import (
	"math"
)

// A simple class keeping track of the mean of a stream of values within a certain window. the WindowedMean will only return a
// value in case enough data has been sampled. After enough data has been sampled the oldest sample will be replaced by the newest
// in case a new sample is added.
type WindowedMean struct {
	values       []float32
	added_values int
	last_value   int
	mean         float32
	dirty        bool
}

// constructor, window_size specifies the number of samples we will continuously get the mean and variance from. the class will
// only return meaning full values if at least window_size values have been added.
// window_size size of the sample window
func NewWindowedMean(window_size int) *WindowedMean {
	return &WindowedMean{dirty: true, values: make([]float32, window_size)}
}

// return whether the value returned will be meaningful
func (self *WindowedMean) HasEnoughData() bool {
	return self.added_values >= len(self.values)
}

// clears this WindowedMean. The class will only return meaningful values after enough data has been added again.
func (self *WindowedMean) Clear() {
	self.added_values = 0
	self.last_value = 0
	for i := 0; i < len(self.values); i++ {
		self.values[i] = 0
	}
	self.dirty = true
}

// adds a new sample to this mean. In case the window is full the oldest value will be replaced by this new value.
func (self *WindowedMean) AddValue(value float32) {
	if self.added_values < len(self.values) {
		self.added_values++
	}
	self.last_value++
	self.values[self.last_value] = value
	if self.last_value > len(self.values)-1 {
		self.last_value = 0
	}
	self.dirty = true
}

// returns the mean of the samples added to this instance. Only returns meaningful results when at least window_size samples
// as specified in the constructor have been added.
func (self *WindowedMean) GetMean() float32 {
	if self.HasEnoughData() {
		if self.dirty == true {
			var mean float32
			for i := 0; i < len(self.values); i++ {
				mean += self.values[i]
			}
			self.mean = mean / float32(len(self.values))
			self.dirty = false
		}
		return self.mean
	} else {
		return 0
	}
}

// return the oldest value in the window
func (self *WindowedMean) GetOldest() float32 {
	if self.last_value == len(self.values)-1 {
		return self.values[0]
	}
	return self.values[self.last_value+1]
}

// return the value last added
func (self *WindowedMean) GetLatest() float32 {
	if self.last_value-1 == -1 {
		return self.values[len(self.values)-1]
	}
	return self.values[self.last_value-1]
}

// return The standard deviation
func (self *WindowedMean) StandardDeviation() float32 {
	if !self.HasEnoughData() {
		return 0
	}

	mean := self.GetMean()
	var sum float32 = 0
	for i := 0; i < len(self.values); i++ {
		sum += (self.values[i] - mean) * (self.values[i] - mean)
	}

	return float32(math.Sqrt(float64(sum / float32(len(self.values)))))
}

func (self *WindowedMean) GetWindowSize() int {
	return len(self.values)
}
