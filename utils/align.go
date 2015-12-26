// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package utils

// Provides bit flag constants for alignment
type Alignment int

const (
	AlignmentCenter = 1 << 0
	AlignmentTop    = 1 << 1
	AlignmentBottom = 1 << 2
	AlignmentLeft   = 1 << 3
	AlignmentRight  = 1 << 4

	AlignmentTopLeft     = AlignmentTop | AlignmentLeft
	AlignmentTopRight    = AlignmentTop | AlignmentRight
	AlignmentBottomLeft  = AlignmentBottom | AlignmentLeft
	AlignmentBottomRight = AlignmentBottom | AlignmentRight
)
