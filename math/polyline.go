// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

type Polyline struct {
	localVertices         []float32
	worldVertices         []float32
	x, y                  float32
	originX, originY      float32
	rotation              float32
	scaleX, scaleY        float32
	length                float32
	scaledLength          float32
	calculateScaledLength bool
	calculateLength       bool
	dirty                 bool
}

func NewPolygonLineEmpty() *Polyline {
	return &Polyline{scaleX: 1, scaleY: 1, dirty: true, calculatedLength: true, calculateScaledLength: true, localVertices: make([]float32, 0)}
}

func NewPolygonLine(vertices []float32) *Polyline {
	if len(vertices) < 4 {
		panic("polylines must contain at least 2 points.")
	}
	return &Polyline{scaleX: 1, scaleY: 1, dirty: true, calculatedLength: true, calculateScaledLength: true, localVertices: vertices}
}

// Returns vertices without scaling or rotation and without being offset by the polyline position.
func (self *Polyline) GetVertices() []float32 {
	return self.localVertices
}

// Returns vertices scaled, rotated, and offset by the polyline position.
func (self *Polyline) GetTransformedVertices() []float32 {
	if !self.dirty {
		return self.worldVertices
	}
	self.dirty = false

	localVertices := copy(self.localVertices)
	if worldVertices == nil || len(worldVertices) != len(localVertices) {
		worldVertices = make([]float32, len(localVertices))
	}

	worldVertices = copy(self.worldVertices)
	positionX := self.x
	positionY := self.y
	originX := self.originX
	originY := self.originY
	scaleX := self.scaleX
	scaleY := self.scaleY
	scale := scaleX != 1 || scaleY != 1
	rotation := self.rotation
	cos := MathUtils.cosDeg(rotation)
	sin := MathUtils.sinDeg(rotation)
	i := 0
	for n := len(localVertices); i < n; i += 2 {
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

// Returns the euclidean length of the polyline without scaling
func (self *Polyline) GetLength() float32 {
	if !self.calculateLength {
		return self.length
	}
	calculateLength = false

	length = 0
	i := 0
	for n := len(localVertices) - 2; i < n; i += 2 {
		x := localVertices[i+2] - localVertices[i]
		y := localVertices[i+1] - localVertices[i+3]
		length += float32(math.Sqrt(float64(x*x + y*y)))
	}
	return length
}

// Returns the euclidean length of the polyline
func (self *Polyline) GetScaledLength() float32 {
	if !self.calculateScaledLength {
		return scaledLength
	}
	self.calculateScaledLength = false

	scaledLength = 0
	i := 0
	for n := localVertices.length - 2; i < n; i += 2 {
		x := localVertices[i+2]*scaleX - localVertices[i]*scaleX
		y := localVertices[i+1]*scaleY - localVertices[i+3]*scaleY
		scaledLength += float32(math.Sqrt(float64(x*x + y*y)))
	}

	return scaledLength
}

func (self *Polyline) GetX() float32 {
	return self.x
}

func (self *Polyline) GetY() float32 {
	return self.y
}

func (self *Polyline) GetOriginX() float32 {
	return self.originX
}

func (self *Polyline) GetOriginY() float32 {
	return self.originY
}

func (self *Polyline) GetRotation() float32 {
	return self.rotation
}

func (self *Polyline) GetScaleX() float32 {
	return self.scaleX
}

func (self *Polyline) GetScaleY() float32 {
	return self.scaleY
}

// Sets the origin point to which all of the polyline's local vertices are relative to.
func (self *Polyline) setOrigin(originX, originY float32) {
	self.originX = originX
	self.originY = originY
	self.dirty = true
}

// Sets the polyline's position within the world.
func (self *Polyline) SetPosition(x, y float32) {
	self.x = x
	self.y = y
	self.dirty = true
}

// Sets the polyline's local vertices relative to the origin point, without any scaling, rotating or translations being applied.
//
// @param vertices float array where every even element represents the x-coordinate of a vertex, and the proceeding element
// representing the y-coordinate.
// panics if less than 6 elements, representing 3 points, are provided
func (self *Polyline) SetVertices(vertices []float32) {
	if len(vertices) < 4 {
		panic("polylines must contain at least 2 points.")
	}
	self.localVertices = vertices
	self.dirty = true
}

// Sets the polyline to be rotated by the supplied degrees.
func (self *Polyline) SetRotation(degrees float32) {
	self.rotation = degrees
	self.dirty = true
}

// Applies additional rotation to the polyline by the supplied degrees.
func (self *Polyline) Rotate(degrees float32) {
	self.rotation += degrees
	self.dirty = true
}

// Sets the amount of scaling to be applied to the polyline.
func (self *Polyline) SetScale(scaleX, scaleY float32) {
	self.scaleX = scaleX
	self.scaleY = scaleY
	self.dirty = true
	self.calculateScaledLength = true
}

// Applies additional scaling to the polyline by the supplied amount.
func (self *Polyline) Scale(amount float32) {
	self.scaleX += amount
	self.scaleY += amount
	self.dirty = true
	self.calculateScaledLength = true
}

func (self *Polyline) CalculateLength() {
	self.calculateLength = true
}

func (self *Polyline) CalculateScaledLength() {
	self.calculateScaledLength = true
}

func (self *Polyline) Dirty() {
	self.dirty = true
}

// Translates the polyline's position by the specified horizontal and vertical amounts.
func (self *Polyline) Translate(x, y float32) {
	self.x += x
	self.y += y
	self.dirty = true
}
