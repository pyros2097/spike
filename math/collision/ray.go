// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package collision

import (
	. "github.com/pyros2097/spike/math/vector"
)

var (
	tmp = NewVector3Empty()
)

// Encapsulates a ray having a starting position and a unit length direction.
type Ray struct {
	Origin    *Vector3
	Direction *Vector3
}

func NewRayEmpty() *Ray {
	return &Ray{}
}

// Constructor, sets the starting position of the ray and the direction.
// param origin The starting position
// param direction The direction
func NewRay(origin, direction *Vector3) *Ray {
	ray := &Ray{}
	ray.Origin.SetV(origin)
	ray.Direction.SetV(direction).Nor()
	return ray
}

// return a copy of this ray.
func (self *Ray) Copy() *Ray {
	return NewRay(self.Origin, self.Direction)
}

// Returns the endpoint given the distance. This is calculated as startpoint + distance * direction.
// param out The vector to set to the result
// param distance The distance from the end point to the start point.
func (self *Ray) GetEndPoint(out *Vector3, distance float32) *Vector3 {
	return out.SetV(self.Direction).SclScalar(distance).AddV(self.Origin)
}

// Multiplies the ray by the given matrix. Use this to transform a ray into another coordinate system.
// param matrix The matrix
func (self *Ray) Mul(matrix *Matrix4) *Ray {
	tmp.SetV(self.Origin).AddV(self.Direction)
	tmp.Mul(matrix)
	self.Origin.Mul(matrix)
	self.Direction.SetV(tmp.SubV(self.Origin))
	return self
}

// Sets the starting position and the direction of this ray.
// param origin The starting position
// param direction The direction
func (self *Ray) SetV3(origin, direction *Vector3) *Ray {
	self.Origin.SetV(origin)
	self.Direction.SetV(direction)
	return self
}

// Sets this ray from the given starting position and direction.
// param x The x-component of the starting position
// param y The y-component of the starting position
// param z The z-component of the starting position
// param dx The x-component of the direction
// param dy The y-component of the direction
// param dz The z-component of the direction
func (self *Ray) Set(x, y, z, dx, dy, dz float32) *Ray {
	self.Origin.Set(x, y, z)
	self.Direction.Set(dx, dy, dz)
	return self
}

// Sets the starting position and direction from the given ray
// param ray The ray
func (self *Ray) SetRay(ray *Ray) *Ray {
	self.Origin.SetV(ray.Origin)
	self.Direction.SetV(ray.Direction)
	return self
}

// func (self *Ray) Equals(other *Ray) bool {
// 	if (o == this) return true;
// 	if (o == null || o.getClass() != self.getClass()) return false;
// 	Ray r = (Ray)o;
// 	return self.Direction.equals(r.direction) && self.Origin.equals(r.origin);
// }

// func (self *Ray) HashCode() int {
// 	final int prime = 73;
// 	int result = 1;
// 	result = prime * result + self.Direction.hashCode();
// 	result = prime * result + self.Origin.hashCode();
// 	return result;
// }

func (self *Ray) String() string {
	return ""
	// return "ray [" + self.Origin + ":" + self.Direction + "]"
}
