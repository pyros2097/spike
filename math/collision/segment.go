// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package collision

import (
	. "github.com/pyros2097/gdx/math/vector"
)

// A Segment is a line in 3-space having a staring and an ending position.
type Segment struct {
	// the starting position
	A *Vector3

	// the ending position
	B *Vector3
}

// Constructs a new Segment from the two points given.
// param b the second point
func NewSegmentV(a, b *Vector3) *Segment {
	segment := &Segment{}
	segment.A.SetV(a)
	segment.B.SetV(b)
	return segment
}

// Constructs a new Segment from the two points given.
// param aX the x-coordinate of the first point
// param aY the y-coordinate of the first point
// param aZ the z-coordinate of the first point
// param bX the x-coordinate of the second point
// param bY the y-coordinate of the second point
// param bZ the z-coordinate of the second point
func NewSegment(aX, aY, aZ, bX, bY, bZ float32) *Segment {
	segment := &Segment{}
	segment.A.Set(aX, aY, aZ)
	segment.B.Set(bX, bY, bZ)
	return segment
}

func (self *Segment) Len() float32 {
	return self.A.DstV(self.B)
}

func (self *Segment) Len2() float32 {
	return self.A.Dst2V(self.B)
}

// public boolean equals (Object o) {
// 	if (o == this) return true;
// 	if (o == null || o.getClass() != self.getClass()) return false;
// 	Segment s = (Segment)o;
// 	return self.A.equals(s.a) && self.B.equals(s.b);
// }

// public int hashCode () {
// 	final int prime = 71;
// 	int result = 1;
// 	result = prime * result + self.A.hashCode();
// 	result = prime * result + self.B.hashCode();
// 	return result;
// }
