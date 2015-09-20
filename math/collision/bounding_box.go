// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package collision

import (
	"math"

	. "github.com/pyros2097/gdx/math/vector"
)

var (
	tmpVector = NewVector3Empty()
)

// Encapsulates an axis aligned bounding box represented by a minimum and a maximum Vector. Additionally you can query for the
// bounding box's center, dimensions and corner points.
type BoundingBox struct {
	Min, Max, Cnt, Dim *Vector3
}

// Constructs a new bounding box with the minimum and maximum vector set to zeros.
func NewBoundingBox() *BoundingBox {
	return (&BoundingBox{&Vector3{}, &Vector3{}, &Vector3{}, &Vector3{}}).Clr()
}

// Constructs a new bounding box from the given bounding box.
// param bounds The bounding box to copy
func NewBoundingBoxCopy(bounds *BoundingBox) *BoundingBox {
	return (&BoundingBox{&Vector3{}, &Vector3{}, &Vector3{}, &Vector3{}}).SetB(bounds)
}

// Constructs the new bounding box using the given minimum and maximum vector.
// param minimum The minimum vector
// param maximum The maximum vector
func NewBoundingBoxVector3(minimum, maximum *Vector3) *BoundingBox {
	return (&BoundingBox{&Vector3{}, &Vector3{}, &Vector3{}, &Vector3{}}).Set(minimum, maximum)
}

// @param out The {@link Vector3} to receive the center of the bounding box.
// return The vector specified with the out argument.
func (self *BoundingBox) GetCenter(out *Vector3) *Vector3 {
	return out.SetV(self.Cnt)
}

func (self *BoundingBox) GetCenterX() float32 {
	return self.Cnt.X
}

func (self *BoundingBox) GetCenterY() float32 {
	return self.Cnt.Y
}

func (self *BoundingBox) GetCenterZ() float32 {
	return self.Cnt.Z
}

func (self *BoundingBox) GetCorner000(out *Vector3) *Vector3 {
	return out.Set(self.Min.X, self.Min.Y, self.Min.Z)
}

func (self *BoundingBox) GetCorner001(out *Vector3) *Vector3 {
	return out.Set(self.Min.X, self.Min.Y, self.Max.Z)
}

func (self *BoundingBox) GetCorner010(out *Vector3) *Vector3 {
	return out.Set(self.Min.X, self.Max.Y, self.Min.Z)
}

func (self *BoundingBox) GetCorner011(out *Vector3) *Vector3 {
	return out.Set(self.Min.X, self.Max.Y, self.Max.Z)
}

func (self *BoundingBox) GetCorner100(out *Vector3) *Vector3 {
	return out.Set(self.Max.X, self.Min.Y, self.Min.Z)
}

func (self *BoundingBox) GetCorner101(out *Vector3) *Vector3 {
	return out.Set(self.Max.X, self.Min.Y, self.Max.Z)
}

func (self *BoundingBox) GetCorner110(out *Vector3) *Vector3 {
	return out.Set(self.Max.X, self.Max.Y, self.Min.Z)
}

func (self *BoundingBox) GetCorner111(out *Vector3) *Vector3 {
	return out.Set(self.Max.X, self.Max.Y, self.Max.Z)
}

// @param out The {@link Vector3} to receive the dimensions of this bounding box on all three axis.
// return The vector specified with the out argument
func (self *BoundingBox) GetDimensions(out *Vector3) *Vector3 {
	return out.SetV(self.Dim)
}

func (self *BoundingBox) GetWidth() float32 {
	return self.Dim.X
}

func (self *BoundingBox) GetHeight() float32 {
	return self.Dim.Y
}

func (self *BoundingBox) GetDepth() float32 {
	return self.Dim.Z
}

// @param out The {@link Vector3} to receive the minimum values.
// return The vector specified with the out argument
func (self *BoundingBox) GetMin(out *Vector3) *Vector3 {
	return out.SetV(self.Min)
}

// @param out The {@link Vector3} to receive the maximum values.
// return The vector specified with the out argument
func (self *BoundingBox) GetMax(out *Vector3) *Vector3 {
	return out.SetV(self.Max)
}

// Sets the given bounding box.
// param bounds The bounds.
// return This bounding box for chaining.
func (self *BoundingBox) SetB(bounds *BoundingBox) *BoundingBox {
	return self.Set(bounds.Min, bounds.Max)
}

// Sets the given minimum and maximum vector.
// param minimum The minimum vector
// param maximum The maximum vector
// return This bounding box for chaining.
func (self *BoundingBox) Set(min, max *Vector3) *BoundingBox {
	var mx, my, mz, mmx, mmy, mmz float32
	switch {
	case min.X < max.X:
		mx = min.X
	case min.X > max.X:
		mx = max.X
	case min.Y < max.Y:
		my = min.Y
	case min.Y > max.Y:
		my = max.Y
	case min.Z < max.Z:
		mz = min.Z
	case min.Z > max.Z:
		mz = max.Z
	}
	self.Min.Set(mx, my, mz)
	switch {
	case min.X < max.X:
		mmx = max.X
	case min.X > max.X:
		mmx = min.X
	case min.Y < max.Y:
		mmy = max.Y
	case min.Y > max.Y:
		mmy = min.Y
	case min.Z < max.Z:
		mmz = max.Z
	case min.Z > max.Z:
		mmz = min.Z
	}
	self.Max.Set(mmx, mmy, mmz)
	self.Cnt.SetV(min).AddV(max).SclScalar(0.5)
	self.Dim.SetV(max).SubV(min)
	return self
}

// Sets the bounding box minimum and maximum vector from the given points.
// param points The points.
// return This bounding box for chaining.
func (self *BoundingBox) SetV(points []*Vector3) *BoundingBox {
	self.Inf()
	for _, l_point := range points {
		self.ExtV3(l_point)
	}
	return self
}

// Todo check Infinity
// Sets the minimum and maximum vector to positive and negative infinity.
// return This bounding box for chaining.
func (self *BoundingBox) Inf() *BoundingBox {
	self.Min.Set(math.MaxFloat32, math.MaxFloat32, math.MaxFloat32)
	self.Max.Set(math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32)
	self.Cnt.Set(0, 0, 0)
	self.Dim.Set(0, 0, 0)
	return self
}

// Sets the minimum and maximum vector to zeros.
// return This bounding box for chaining.
func (self *BoundingBox) Clr() *BoundingBox {
	return self.Set(self.Min.Set(0, 0, 0), self.Max.Set(0, 0, 0))
}

// Returns whether this bounding box is valid. This means that {@link #max} is greater than {@link #min}.
// return True in case the bounding box is valid, false otherwise
func (self *BoundingBox) IsValid() bool {
	return self.Min.X < self.Max.X && self.Min.Y < self.Max.Y && self.Min.Z < self.Max.Z
}

// Extends the bounding box to incorporate the given {@link Vector3}.
// param point The vector
// return This bounding box for chaining.
func (self *BoundingBox) ExtV3(point *Vector3) *BoundingBox {
	return self.Set(self.Min.Set(min(self.Min.X, point.X), min(self.Min.Y, point.Y), min(self.Min.Z, point.Z)),
		self.Max.Set(max(self.Max.X, point.X), max(self.Max.Y, point.Y), max(self.Max.Z, point.Z)))
}

// Extends the bounding box by the given vector.
// param x The x-coordinate
// param y The y-coordinate
// param z The z-coordinate
// return This bounding box for chaining.
func (self *BoundingBox) Ext(x, y, z float32) *BoundingBox {
	return self.Set(self.Min.Set(min(self.Min.X, x), min(self.Min.Y, y), min(self.Min.Z, z)),
		self.Max.Set(max(self.Max.X, x), max(self.Max.Y, y), max(self.Max.Z, z)))
}

// Extends this bounding box by the given bounding box.
// param bounds The bounding box
// return This bounding box for chaining.
func (self *BoundingBox) ExtB(bounds *BoundingBox) *BoundingBox {
	return self.Set(self.Min.Set(min(self.Min.X, bounds.Min.X), min(self.Min.Y, bounds.Min.Y), min(self.Min.Z, bounds.Min.Z)),
		self.Max.Set(max(self.Max.X, bounds.Max.X), max(self.Max.Y, bounds.Max.Y), max(self.Max.Z, bounds.Max.Z)))
}

// Extends this bounding box by the given transformed bounding box.
// param bounds The bounding box
// param transform The transformation matrix to apply to bounds, before using it to extend this bounding box.
// return This bounding box for chaining.
func (self *BoundingBox) ExtBT(bounds *BoundingBox, transform *Matrix4) *BoundingBox {
	self.ExtV3(tmpVector.Set(bounds.Min.X, bounds.Min.Y, bounds.Min.Z).Mul(transform))
	self.ExtV3(tmpVector.Set(bounds.Min.X, bounds.Min.Y, bounds.Max.Z).Mul(transform))
	self.ExtV3(tmpVector.Set(bounds.Min.X, bounds.Max.Y, bounds.Min.Z).Mul(transform))
	self.ExtV3(tmpVector.Set(bounds.Min.X, bounds.Max.Y, bounds.Max.Z).Mul(transform))
	self.ExtV3(tmpVector.Set(bounds.Max.X, bounds.Min.Y, bounds.Min.Z).Mul(transform))
	self.ExtV3(tmpVector.Set(bounds.Max.X, bounds.Min.Y, bounds.Max.Z).Mul(transform))
	self.ExtV3(tmpVector.Set(bounds.Max.X, bounds.Max.Y, bounds.Min.Z).Mul(transform))
	self.ExtV3(tmpVector.Set(bounds.Max.X, bounds.Max.Y, bounds.Max.Z).Mul(transform))
	return self
}

// Multiplies the bounding box by the given matrix. This is achieved by multiplying the 8 corner points and then calculating
// the minimum and maximum vectors from the transformed points.
// param transform The matrix
// return This bounding box for chaining.
func (self *BoundingBox) Mul(transform *Matrix4) *BoundingBox {
	x0 := self.Min.X
	y0 := self.Min.Y
	z0 := self.Min.Z
	x1 := self.Max.X
	y1 := self.Max.Y
	z1 := self.Max.Z
	self.Inf()
	self.ExtV3(tmpVector.Set(x0, y0, z0).Mul(transform))
	self.ExtV3(tmpVector.Set(x0, y0, z1).Mul(transform))
	self.ExtV3(tmpVector.Set(x0, y1, z0).Mul(transform))
	self.ExtV3(tmpVector.Set(x0, y1, z1).Mul(transform))
	self.ExtV3(tmpVector.Set(x1, y0, z0).Mul(transform))
	self.ExtV3(tmpVector.Set(x1, y0, z1).Mul(transform))
	self.ExtV3(tmpVector.Set(x1, y1, z0).Mul(transform))
	self.ExtV3(tmpVector.Set(x1, y1, z1).Mul(transform))
	return self
}

// Returns whether the given bounding box is contained in this bounding box.
// param b The bounding box
// return Whether the given bounding box is contained
func (self *BoundingBox) Contains(b *BoundingBox) bool {
	return !self.IsValid() || (self.Min.X <= b.Min.X && self.Min.Y <= b.Min.Y &&
		self.Min.Z <= b.Min.Z && self.Max.X >= b.Max.X && self.Max.Y >= b.Max.Y && self.Max.Z >= b.Max.Z)
}

// Returns whether the given bounding box is intersecting this bounding box (at least one point in).
// param b The bounding box
// return Whether the given bounding box is intersected
func (self *BoundingBox) Intersects(b *BoundingBox) bool {
	if !self.IsValid() {
		return false
	}
	// test using SAT (separating axis theorem)

	lx := float32(math.Abs(float64(self.Cnt.X - b.Cnt.X)))
	sumx := (self.Dim.X / 2.0) + (b.Dim.X / 2.0)

	ly := float32(math.Abs(float64(self.Cnt.Y - b.Cnt.Y)))
	sumy := (self.Dim.Y / 2.0) + (b.Dim.Y / 2.0)

	lz := float32(math.Abs(float64(self.Cnt.Z - b.Cnt.Z)))
	sumz := (self.Dim.Z / 2.0) + (b.Dim.Z / 2.0)

	return (lx <= sumx && ly <= sumy && lz <= sumz)
}

// Returns whether the given vector is contained in this bounding box.
// param v The vector
// return Whether the vector is contained or not.
func (self *BoundingBox) ContainsV3(v *Vector3) bool {
	return self.Min.X <= v.X && self.Max.X >= v.X && self.Min.Y <= v.Y && self.Max.Y >= v.Y && self.Min.Z <= v.Z && self.Max.Z >= v.Z
}

func (self *BoundingBox) String() string {
	return ""
	// return "[" + self.Min + "|" + self.Max + "]"
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
