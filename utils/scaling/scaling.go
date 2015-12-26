// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package scaling

import (
	. "github.com/pyros2097/spike/math/vector"
)

// Various scaling types for fitting one rectangle into another.
type Scaling uint8

const (
	// Scales the source to fit the target while keeping the same aspect ratio. This may cause the source to be smaller than the
	// target in one direction.
	Fit = iota

	// Scales the source to fill the target while keeping the same aspect ratio. This may cause the source to be larger than the
	// target in one direction.
	Fill

	// Scales the source to fill the target in the x direction while keeping the same aspect ratio. This may cause the source to be
	// smaller or larger than the target in the y direction.
	FillX

	// Scales the source to fill the target in the y direction while keeping the same aspect ratio. This may cause the source to be
	// smaller or larger than the target in the x direction.
	FillY

	// Scales the source to fill the target. This may cause the source to not keep the same aspect ratio.
	Stretch

	// Scales the source to fill the target in the x direction, without changing the y direction. This may cause the source to not
	// keep the same aspect ratio.
	StretchX

	// Scales the source to fill the target in the y direction, without changing the x direction. This may cause the source to not
	// keep the same aspect ratio.
	StretchY

	// The source is not scaled.
	None
)

var (
	temp = NewVector2Empty()
)

// Returns the size of the source scaled to the target. Note the same Vector2 instance is always returned and should never be
//cached.
func (self Scaling) Apply(sourceWidth, sourceHeight, targetWidth, targetHeight float32) *Vector2 {
	switch self {
	case Fit:
		targetRatio := targetHeight / targetWidth
		sourceRatio := sourceHeight / sourceWidth
		var scale float32
		if targetRatio > sourceRatio {
			scale = targetWidth / sourceWidth
		} else {
			scale = targetHeight / sourceHeight
		}
		temp.X = sourceWidth * scale
		temp.Y = sourceHeight * scale
	case Fill:
		targetRatio := targetHeight / targetWidth
		sourceRatio := sourceHeight / sourceWidth
		var scale float32
		if targetRatio < sourceRatio {
			scale = targetWidth / sourceWidth
		} else {
			scale = targetHeight / sourceHeight
		}
		temp.X = sourceWidth * scale
		temp.Y = sourceHeight * scale
	case FillX:
		scale := targetWidth / sourceWidth
		temp.X = sourceWidth * scale
		temp.Y = sourceHeight * scale
	case FillY:
		scale := targetHeight / sourceHeight
		temp.X = sourceWidth * scale
		temp.Y = sourceHeight * scale
	case Stretch:
		temp.X = targetWidth
		temp.Y = targetHeight
	case StretchX:
		temp.X = targetWidth
		temp.Y = sourceHeight
	case StretchY:
		temp.X = sourceWidth
		temp.Y = targetHeight
	case None:
		temp.X = sourceWidth
		temp.Y = sourceHeight
	}
	return temp
}
