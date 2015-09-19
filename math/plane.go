// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

// Enum specifying on which side a point lies respective to the plane and it's normal. {@link PlaneSide#Front} is the side to
// which the normal points.
type PlaneSide int

const (
	PlaneOn PlaneSide = iota
	PlaneBack
	PlaneFront
)

// A plane defined via a unit length normal and the distance from the origin, as you learned in your math class.
type Plane struct {
	normal *Vector3
	d      float32
}

// Constructs a new plane with all values set to 0
func NewPlaneEmpty() *Plane {
	return &Plane{normal: NewVector3Empty(), d: 0}
}

// Constructs a new plane based on the normal and distance to the origin.
// param normal The plane normal
// param d The distance to the origin
func NewPlane(normal *Vector3, d float32) *Plane {
	return &Plane{normal: NewVector3Empty().Set(normal).Nor(), d: d}
}

// Constructs a new plane based on the normal and a point on the plane.
// param normal The normal
// param point The point on the plane
func NewPlanePoint(normal, point *Vector3) *Plane {
	normal = NewVector3Empty().Set(normal).Nor()
	return &Plane{normal: normal, d: -normal.DotV(point)}
}

// Constructs a new plane out of the three given points that are considered to be on the plane. The normal is calculated via a
// cross product between (point1-point2)x(point2-point3)
// param point1 The first point
// param point2 The second point
// param point3 The third point
func NewPlanePoint3(point1, point2, point3 *Vector3) *Plane {
	plane := &Plane{normal: NewVector3Empty(), d: 0}
	plane.Set(point1, point2, point3)
	return plane
}

// Sets the plane normal and distance to the origin based on the three given points which are considered to be on the plane.
// The normal is calculated via a cross product between (point1-point2)x(point2-point3)
// param point1
// param point2
// param point3
func (self *Plane) SetP3(point1, point2, point3 *Vector3) {
	self.normal.SetV(point1).SubV(point2).Crs(point2.x-point3.x, point2.y-point3.y, point2.z-point3.z).Nor()
	self.d = -point1.DotV(self.normal)
}

// Sets the plane normal and distance
// param nx normal x-component
// param ny normal y-component
// param nz normal z-component
// param d distance to origin
func (self *Plane) Set(nx, ny, nz, d float32) {
	self.normal.Set(nx, ny, nz)
	self.d = d
}

// Sets the plane to the given point and normal.
// param normal the normal of the plane
func (self *Plane) SetV(point, normal *Vector3) {
	self.normal.SetV(normal)
	self.d = -point.DotV(normal)
}

func (self *Plane) SetP6(pointX, pointY, pointZ, norX, norY, norZ float32) {
	self.normal.Set(norX, norY, norZ)
	self.d = -(pointX*norX + pointY*norY + pointZ*norZ)
}

// Sets this plane from the given plane
// param plane the plane
func (self *Plane) setPlane(plane *Plane) {
	self.normal.SetV(plane.normal)
	self.d = plane.d
}

// Calculates the shortest signed distance between the plane and the given point.
// param point The point
// return the shortest signed distance between the plane and the point
func (self *Plane) Distance(point *Vector3) float32 {
	return self.normal.DotV(point) + self.d
}

// Returns on which side the given point lies relative to the plane and its normal. PlaneSide.Front refers to the side the
// plane normal points to.
// param point The point
// return The side the point lies relative to the plane
func (self *Plane) TestPointV(point *Vector3) PlaneSide {
	dist := self.normal.DotV(point) + self.d
	if dist == 0 {
		return PlaneOn
	} else if dist < 0 {
		return PlaneBack
	} else {
		return PlaneFront
	}
}

// Returns on which side the given point lies relative to the plane and its normal. PlaneSide.Front refers to the side the
// plane normal points to.
// return The side the point lies relative to the plane
func (self *Plane) TestPoint(x, y, z float32) PlaneSide {
	dist := self.normal.Dot(x, y, z) + self.d
	if dist == 0 {
		return PlaneOn
	} else if dist < 0 {
		return PlaneBack
	} else {
		return PlaneFront
	}
}

// Returns whether the plane is facing the direction vector. Think of the direction vector as the direction a camera looks in.
// This method will return true if the front side of the plane determined by its normal faces the camera.
// param direction the direction
// return whether the plane is front facing
func (self *Plane) IsFrontFacing(direction *Vector3) bool {
	return dot <= self.normal.DotV(direction)
}

func (self *Plane) GetNormal() *Vector3 {
	return self.normal
}

func (self *Plane) GetD() float32 {
	return self.d
}

func (self *Plane) String() string {
	return self.normal.String() + ", " + self.d
}
