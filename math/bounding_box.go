// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package collision

import (
	"math"
)

var (
	tmpVector = NewVector3()
)

// Encapsulates an axis aligned bounding box represented by a minimum and a maximum Vector. Additionally you can query for the
// bounding box's center, dimensions and corner points.
type BoundingBox struct {
	min, max, cnt, dim *Vector3
}

// Constructs a new bounding box with the minimum and maximum vector set to zeros.
func NewBoundingBox() *BoundingBox {
	return &BoundingBox{&Vector3{}, &Vector3{}, &Vector3{}, &Vector3{}}.Clr()
}

// Constructs a new bounding box from the given bounding box.
// param bounds The bounding box to copy
func NewBoundingBoxCopy(bounds *BoundingBox) *BoundingBox {
	return &BoundingBox{&Vector3{}, &Vector3{}, &Vector3{}, &Vector3{}}.SetB(bounds)
}

// Constructs the new bounding box using the given minimum and maximum vector.
// param minimum The minimum vector
// param maximum The maximum vector
func NewBoundingBoxVector3(minimum, maximum *Vector3) *BoundingBox {
	return &BoundingBox{&Vector3{}, &Vector3{}, &Vector3{}, &Vector3{}}.Set(minimum, maximum)
}

// @param out The {@link Vector3} to receive the center of the bounding box.
// return The vector specified with the out argument.
func (self *BoundingBox) GetCenter(out *Vector3) *Vector3 {
	return out.SetV(self.cnt)
}

func (self *BoundingBox) GetCenterX() float32 {
	return self.cnt.x
}

func (self *BoundingBox) GetCenterY() float32 {
	return self.cnt.y
}

func (self *BoundingBox) GetCenterZ() float32 {
	return self.cnt.z
}

func (self *BoundingBox) GetCorner000(out *Vector3) *Vector3 {
	return out.Set(self.min.x, self.min.y, self.min.z)
}

func (self *BoundingBox) GetCorner001(out *Vector3) *Vector3 {
	return out.Set(self.min.x, self.min.y, self.max.z)
}

func (self *BoundingBox) GetCorner010(out *Vector3) *Vector3 {
	return out.Set(self.min.x, self.max.y, self.min.z)
}

func (self *BoundingBox) GetCorner011(out *Vector3) *Vector3 {
	return out.Set(self.min.x, self.max.y, self.max.z)
}

func (self *BoundingBox) GetCorner100(out *Vector3) *Vector3 {
	return out.Set(self.max.x, self.min.y, self.min.z)
}

func (self *BoundingBox) GetCorner101(out *Vector3) *Vector3 {
	return out.Set(self.max.x, self.min.y, self.max.z)
}

func (self *BoundingBox) GetCorner110(out *Vector3) *Vector3 {
	return out.Set(self.max.x, self.max.y, self.min.z)
}

func (self *BoundingBox) GetCorner111(out *Vector3) *Vector3 {
	return out.Set(self.max.x, self.max.y, self.max.z)
}

// @param out The {@link Vector3} to receive the dimensions of this bounding box on all three axis.
// return The vector specified with the out argument
func (self *BoundingBox) GetDimensions(out *Vector3) *Vector3 {
	return out.SetV(self.dim)
}

func (self *BoundingBox) GetWidth() float32 {
	return self.dim.x
}

func (self *BoundingBox) GetHeight() float32 {
	return self.dim.y
}

func (self *BoundingBox) GetDepth() float32 {
	return self.dim.z
}

// @param out The {@link Vector3} to receive the minimum values.
// return The vector specified with the out argument
func (self *BoundingBox) GetMin(out *Vector3) *Vector3 {
	return out.Set(self.min)
}

// @param out The {@link Vector3} to receive the maximum values.
// return The vector specified with the out argument
func (self *BoundingBox) GetMax(out *Vector3) *Vector3 {
	return out.Set(self.max)
}

// Sets the given bounding box.
// param bounds The bounds.
// return This bounding box for chaining.
func (self *BoundingBox) SetB(bounds *BoundingBox) *BoundingBox {
	return self.Set(bounds.min, bounds.max)
}

// Sets the given minimum and maximum vector.
// param minimum The minimum vector
// param maximum The maximum vector
// return This bounding box for chaining.
func (self *BoundingBox) Set(min, max *Vector3) *BoundingBox {
	switch {
	case min.x < max.x:
		mx := min.x
	case min.x > max.x:
		mx := max.x
	case min.y < max.y:
		my := min.y
	case min.y > max.y:
		my := max.y
	case min.z < max.z:
		mz := min.z
	case min.z > max.z:
		mz := max.z
	}
	self.min.Set(mx, my, mz)
	switch {
	case min.x < max.x:
		mmx := max.x
	case min.x > max.x:
		mmx := min.x
	case min.y < max.y:
		mmy := max.y
	case min.y > max.y:
		mmy := min.y
	case min.z < max.z:
		mmz := max.z
	case min.z > max.z:
		mmz := min.z
	}
	self.max.Set(mmx, mmy, mmz)
	self.cnt.Set(min).Add(max).Scl(0.5)
	self.dim.Set(max).Sub(min)
	return self
}

// Sets the bounding box minimum and maximum vector from the given points.
// param points The points.
// return This bounding box for chaining.
func (self *BoundingBox) SetV(points *[]Vector3) {
	self.Inf()
	for _, l_point := range points {
		self.Ext(l_point)
	}
	return self
}

// Todo check Infinity
// Sets the minimum and maximum vector to positive and negative infinity.
// return This bounding box for chaining.
func (self *BoundingBox) Inf() *BoundingBox {
	self.min.Set(math.MaxFloat32, math.MaxFloat32, math.MaxFloat32)
	self.max.Set(math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32)
	self.cnt.Set(0, 0, 0)
	self.dim.Set(0, 0, 0)
	return self
}

// Extends the bounding box to incorporate the given {@link Vector3}.
// param point The vector
// return This bounding box for chaining.
func (self *BoundingBox) Ext(point *Vector) *BoundingBox {
	return self.Set(min.Set(min(self.min.x, point.x), min(min.y, point.y), min(min.z, point.z)),
		self.max.set(math.Max(self.max.x, point.x), math.Max(max.y, point.y), math.Max(max.z, point.z)))
}

// Sets the minimum and maximum vector to zeros.
// return This bounding box for chaining.
func (self *BoundingBox) Clr() *BoundingBox {
	return self.Set(self.min.Set(0, 0, 0), self.max.Set(0, 0, 0))
}

// Returns whether this bounding box is valid. This means that {@link #max} is greater than {@link #min}.
// return True in case the bounding box is valid, false otherwise
func (self *BoundingBox) IsValid() bool {
	return self.min.x < self.max.x && self.min.y < self.max.y && self.min.z < self.max.z
}

// Extends the bounding box by the given vector.
// param x The x-coordinate
// param y The y-coordinate
// param z The z-coordinate
// return This bounding box for chaining.
func (self *BoundingBox) Ext(x, y, z float32) {
	return self.Set(self.min.Set(min(self.min.x, x), min(min.y, y), min(min.z, z)), self.max.Set(max(self.max.x, x), max(max.y, y), max(max.z, z)))
}

// Extends this bounding box by the given bounding box.
// param a_bounds The bounding box
// return This bounding box for chaining.
func (self *BoundingBox) ExtB(a_bounds *BoundingBox) *BoundingBox {
	return self.Set(self.min.Set(min(self.min.x, a_bounds.min.x), min(min.y, a_bounds.min.y), min(min.z, a_bounds.min.z)),
		self.max.set(max(self.max.x, a_bounds.max.x), max(max.y, a_bounds.max.y), max(max.z, a_bounds.max.z)))
}

// Extends this bounding box by the given transformed bounding box.
// param bounds The bounding box
// param transform The transformation matrix to apply to bounds, before using it to extend this bounding box.
// return This bounding box for chaining.
func (self *BoundingBox) ExtBT(bounds *BoundingBox, transform *Matrix4) {
	self.Ext(tmpVector.Set(bounds.min.x, bounds.min.y, bounds.min.z).Mul(transform))
	self.Ext(tmpVector.Set(bounds.min.x, bounds.min.y, bounds.max.z).Mul(transform))
	self.Ext(tmpVector.Set(bounds.min.x, bounds.max.y, bounds.min.z).Mul(transform))
	self.Ext(tmpVector.Set(bounds.min.x, bounds.max.y, bounds.max.z).Mul(transform))
	self.Ext(tmpVector.Set(bounds.max.x, bounds.min.y, bounds.min.z).Mul(transform))
	self.Ext(tmpVector.Set(bounds.max.x, bounds.min.y, bounds.max.z).Mul(transform))
	self.Ext(tmpVector.Set(bounds.max.x, bounds.max.y, bounds.min.z).Mul(transform))
	self.Ext(tmpVector.Set(bounds.max.x, bounds.max.y, bounds.max.z).Mul(transform))
	return self
}

// Multiplies the bounding box by the given matrix. This is achieved by multiplying the 8 corner points and then calculating
// the minimum and maximum vectors from the transformed points.
// param transform The matrix
// return This bounding box for chaining.
func (self *BoundingBox) Mul(transform *Matrix4) *BoundingBox {
	x0 := self.min.x
	y0 := self.min.y
	z0 := self.min.z
	x1 := self.max.x
	y1 := self.max.y
	z1 := self.max.z
	self.Inf()
	self.Ext(tmpVector.Set(x0, y0, z0).Mul(transform))
	self.Ext(tmpVector.Set(x0, y0, z1).Mul(transform))
	self.Ext(tmpVector.Set(x0, y1, z0).Mul(transform))
	self.Ext(tmpVector.Set(x0, y1, z1).Mul(transform))
	self.Ext(tmpVector.Set(x1, y0, z0).Mul(transform))
	self.Ext(tmpVector.Set(x1, y0, z1).Mul(transform))
	self.Ext(tmpVector.Set(x1, y1, z0).Mul(transform))
	self.Ext(tmpVector.Set(x1, y1, z1).Mul(transform))
	return self
}

// Returns whether the given bounding box is contained in this bounding box.
// param b The bounding box
// return Whether the given bounding box is contained
func (self *BoundingBox) Contains(b *BoundingBox) bool {
	return !self.IsValid() || (self.min.x <= b.min.x && self.min.y <= b.min.y && self.min.z <= b.min.z && max.x >= b.max.x && self.max.y >= b.max.y && self.max.z >= b.max.z)
}

// Returns whether the given bounding box is intersecting this bounding box (at least one point in).
// param b The bounding box
// return Whether the given bounding box is intersected
func (self *BoundingBox) Intersects(b *BoundingBox) bool {
	if !self.IsValid() {
		return false
	}
	// test using SAT (separating axis theorem)

	lx := math.Abs(self.cnt.x - b.cnt.x)
	sumx := (self.dim.x / 2.0) + (b.dim.x / 2.0)

	ly := math.Abs(self.cnt.y - b.cnt.y)
	sumy := (self.dim.y / 2.0) + (b.dim.y / 2.0)

	lz := math.Abs(self.cnt.z - b.cnt.z)
	sumz := (self.dim.z / 2.0) + (b.dim.z / 2.0)

	return (lx <= sumx && ly <= sumy && lz <= sumz)
}

// Returns whether the given vector is contained in this bounding box.
// param v The vector
// return Whether the vector is contained or not.
func (self *BoundingBox) ContainsV3(v *Vector3) bool {
	return self.min.x <= v.x && self.max.x >= v.x && self.min.y <= v.y && self.max.y >= v.y && self.min.z <= v.z && self.max.z >= v.z
}

func (self *BoundingBox) ToString() string {
	return "[" + self.min + "|" + self.max + "]"
}

func min(a, b float32) float32 {
	if a > b {
		return b
	}
	return a
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
