// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// The type of Gestures that can occurr
package gesture

type Type uint8

const (
	SWIPE_LEFT Type = iota
	SWIPE_RIGHT
	SWIPE_UP
	SWIPE_DOWN
	SWIPE_NONE
)
