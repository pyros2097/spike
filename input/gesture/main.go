// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gesture

// The type of Gesture that can occurr
type Type uint8

const (
	SWIPE_LEFT Type = iota
	SWIPE_RIGHT
	SWIPE_UP
	SWIPE_DOWN
	SWIPE_NONE
)
