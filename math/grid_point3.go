// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

// A point in a 3D grid, with integer x and y coordinates
type GridPoint3 struct {
	x, y, z int
}

// Constructs a 3D grid point with all coordinates pointing to the origin (0, 0, 0).
func NewGridPoint3Empty() *GridPoint3 {
	return &GridPoint3{}
}

// Constructs a 3D grid point.
func NewGridPoint3(x, y, z int) *GridPoint3 {
	return &GridPoint3{x, y, z}
}

// Copy constructor
func NewGridPoint3Copy(point *GridPoint3) *GridPoint3 {
	return &GridPoint3{point.x, point.y, point.z}
}

// Sets the coordinates of this 3D grid point to that of another.
func (self *GridPoint3) SetG(point *GridPoint3) {
	self.x = point.x
	self.y = point.y
	self.z = point.z
	return self
}

// Sets the coordinates of this GridPoint3D.
func (self *GridPoint3) Set(x, y, z int) {
	self.x = x
	self.y = y
	self.z = z
	return self
}

// func (self *GridPoint3) Equals() bool {
// 	if (this == o) return true;
// 	if (o == null || o.getClass() != this.getClass()) return false;
// 	GridPoint3 g = (GridPoint3)o;
// 	return this.x == g.x && this.y == g.y && this.z == g.z;
// }

func (self *GridPoint3) HashCode() int {
	prime := 17
	result := 1
	result = prime*result + self.x
	result = prime*result + self.y
	result = prime*result + self.z
	return result
}

func (self *GridPoint3) String() string {
	return "(" + self.x + ", " + self.y + ", " + self.z + ")"
}
