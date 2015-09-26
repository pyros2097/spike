// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

import (
	"github.com/pyros2097/spike/math/geom"
	. "github.com/pyros2097/spike/math/shape"
	"github.com/pyros2097/spike/math/utils"
)

// Encapsulates a 2D polygon defined by it's vertices relative to an origin point (default of 0, 0).
type Polygon struct {
	localVertices    []float32
	worldVertices    []float32
	x, y             float32
	originX, originY float32
	rotation         float32
	scaleX, scaleY   float32
	dirty            bool
	bounds           *Rectangle
}

// Constructs a new polygon with no vertices
func NewPolygonEmpty() *Polygon {
	return &Polygon{scaleX: 1, scaleY: 1, dirty: true, localVertices: make([]float32, 0)}
}

// Constructs a new polygon from a float array of parts of vertex points.
// @param vertices an array where every even element represents the horizontal part of a point, and the following element
//          representing the vertical part
// panics if less than 6 elements, representing 3 points, are provided
func NewPolygon(vertices []float32) *Polygon {
	if len(vertices) < 6 {
		panic("polygons must contain at least 3 points.")
	}
	return &Polygon{scaleX: 1, scaleY: 1, dirty: true, localVertices: vertices}
}

// Returns the polygon's local vertices without scaling or rotation and without being offset by the polygon position.
func (self *Polygon) GetVertices() []float32 {
	return self.localVertices
}

// Calculates and returns the vertices of the polygon after scaling, rotation, and positional translations have been applied,
// as they are position within the world.
//
// @return vertices scaled, rotated, and offset by the polygon position.
func (self *Polygon) GetTransformedVertices() []float32 {
	if !self.dirty {
		return self.worldVertices
	}
	self.dirty = false
	var localVertices []float32
	copy(localVertices, self.localVertices)
	if self.worldVertices == nil || len(self.worldVertices) != len(localVertices) {
		self.worldVertices = make([]float32, len(localVertices))
	}

	var worldVertices []float32
	copy(worldVertices, self.worldVertices)
	positionX := self.x
	positionY := self.y
	originX := self.originX
	originY := self.originY
	scaleX := self.scaleX
	scaleY := self.scaleY
	scale := scaleX != 1 || scaleY != 1
	rotation := self.rotation
	cos := utils.CosDeg(rotation)
	sin := utils.SinDeg(rotation)
	n := len(localVertices)
	for i := 0; i < n; i += 2 {
		x := localVertices[i] - originX
		y := localVertices[i+1] - originY

		// scale if needed
		if scale {
			x *= scaleX
			y *= scaleY
		}

		// rotate if needed
		if rotation != 0 {
			oldX := x
			x = cos*x - sin*y
			y = sin*oldX + cos*y
		}

		worldVertices[i] = positionX + x + originX
		worldVertices[i+1] = positionY + y + originY
	}
	return worldVertices
}

// Sets the origin point to which all of the polygon's local vertices are relative to.
func (self *Polygon) setOrigin(originX, originY float32) {
	self.originX = originX
	self.originY = originY
	self.dirty = true
}

// Sets the polygon's position within the world.
func (self *Polygon) SetPosition(x, y float32) {
	self.x = x
	self.y = y
	self.dirty = true
}

// Sets the polygon's local vertices relative to the origin point, without any scaling, rotating or translations being applied.
//
// @param vertices float array where every even element represents the x-coordinate of a vertex, and the proceeding element
// representing the y-coordinate.
// panics if less than 6 elements, representing 3 points, are provided
func (self *Polygon) SetVertices(vertices []float32) {
	if len(vertices) < 6 {
		panic("polygons must contain at least 3 points.")
	}
	self.localVertices = vertices
	self.dirty = true
}

// Translates the polygon's position by the specified horizontal and vertical amounts.
func (self *Polygon) Translate(x, y float32) {
	self.x += x
	self.y += y
	self.dirty = true
}

// Sets the polygon to be rotated by the supplied degrees.
func (self *Polygon) SetRotation(degrees float32) {
	self.rotation = degrees
	self.dirty = true
}

// Applies additional rotation to the polygon by the supplied degrees.
func (self *Polygon) Rotate(degrees float32) {
	self.rotation += degrees
	self.dirty = true
}

// Sets the amount of scaling to be applied to the polygon.
func (self *Polygon) SetScale(scaleX, scaleY float32) {
	self.scaleX = scaleX
	self.scaleY = scaleY
	self.dirty = true
}

// Applies additional scaling to the polygon by the supplied amount.
func (self *Polygon) Scale(amount float32) {
	self.scaleX += amount
	self.scaleY += amount
	self.dirty = true
}

// Sets the polygon's world vertices to be recalculated when calling {@link #getTransformedVertices() getTransformedVertices}.
func (self *Polygon) Dirty() {
	self.dirty = true
}

// Returns the area contained within the polygon.
func (self *Polygon) Area() float32 {
	vertices := self.GetTransformedVertices()
	return geom.PolygonArea(vertices, 0, len(vertices))
}

// Returns an axis-aligned bounding box of this polygon.
// Note the returned Rectangle is cached in this polygon, and will be reused if this Polygon is changed.
// @return this polygon's bounding box {@link Rectangle}
func (self *Polygon) GetBoundingRectangle() *Rectangle {
	vertices := self.GetTransformedVertices()

	minX := vertices[0]
	minY := vertices[1]
	maxX := vertices[0]
	maxY := vertices[1]

	numFloats := len(vertices)
	for i := 2; i < numFloats; i += 2 {
		switch {
		case minX > vertices[i]:
			minX = vertices[i]
		case minY > vertices[i+1]:
			minY = vertices[i+1]
		case maxX < vertices[i]:
			maxX = vertices[i]
		case maxY < vertices[i+1]:
			maxY = vertices[i+1]
		}
	}

	if self.bounds == nil {
		self.bounds = NewRectangleEmpty()
	}
	self.bounds.X = minX
	self.bounds.Y = minY
	self.bounds.W = maxX - minX
	self.bounds.H = maxY - minY

	return self.bounds
}

// Returns whether an x, y pair is contained within the polygon.
func (self *Polygon) Contains(x, y float32) bool {
	vertices := self.GetTransformedVertices()
	numFloats := len(vertices)
	intersects := 0

	for i := 0; i < numFloats; i += 2 {
		x1 := vertices[i]
		y1 := vertices[i+1]
		x2 := vertices[(i+2)%numFloats]
		y2 := vertices[(i+3)%numFloats]
		if ((y1 <= y && y < y2) || (y2 <= y && y < y1)) && x < ((x2-x1)/(y2-y1)*(y-y1)+x1) {
			intersects++
		}
	}
	return (intersects & 1) == 1
}

// Returns the x-coordinate of the polygon's position within the world.
func (self *Polygon) getX() float32 {
	return self.x
}

// Returns the y-coordinate of the polygon's position within the world.
func (self *Polygon) getY() float32 {
	return self.y
}

// Returns the x-coordinate of the polygon's origin point.
func (self *Polygon) getOriginX() float32 {
	return self.originX
}

// Returns the y-coordinate of the polygon's origin point.
func (self *Polygon) getOriginY() float32 {
	return self.originY
}

// Returns the total rotation applied to the polygon.
func (self *Polygon) getRotation() float32 {
	return self.rotation
}

// Returns the total horizontal scaling applied to the polygon.
func (self *Polygon) getScaleX() float32 {
	return self.scaleX
}

// Returns the total vertical scaling applied to the polygon.
func (self *Polygon) getScaleY() float32 {
	return self.scaleY
}
