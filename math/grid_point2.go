// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

// A point in a 2D grid, with integer x and y coordinates
type GridPoint2 struct {
	x, y int
}

// Constructs a new 2D grid point.
func NewGridPoint2Empty() *GridPoint2 {
	return &GridPoint2{}
}

// Constructs a new 2D grid point.
func NewGridPoint2(x, y int) *GridPoint2 {
	return &GridPoint2{x, y}
}

// Copy constructor
func NewGridPoint2Copy(point *GridPoint2) *GridPoint2 {
	return &GridPoint2{point.x, point.y}
}

// Sets the coordinates of this 2D grid point to that of another.
// point The 2D grid point to copy the coordinates of.
// return this 2D grid point for chaining.
func (self *GridPoint2) SetG(point *GridPoint2) *GridPoint2 {
	self.x = point.x
	self.y = point.y
	return self
}

// Sets the coordinates of this 2D grid point.
func (self *GridPoint2) Set(x, y int) *GridPoint2 {
	self.x = x
	self.y = y
	return self
}

// func (self *GridPoint2) Equals() bool {
// 	if (this == o) return true;
// 	if (o == null || o.getClass() != self.getClass()) return false;
// 	GridPoint2 g = (GridPoint2)o;
// 	return self.x == g.x && self.y == g.y;
// }

func (self *GridPoint2) HashCode() int {
	prime := 53
	result := 1
	result = prime*result + self.x
	result = prime*result + self.y
	return result
}

func (self *GridPoint2) String() string {
	return ""
	// return "(" + self.x + ", " + self.y + ")"
}
