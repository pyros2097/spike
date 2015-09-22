// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package shape

import (
	"math"

	"github.com/pyros2097/spike/math/utils"
	. "github.com/pyros2097/spike/math/vector"
)

// A convenient 2D ellipse class, based on the circle class
// Implements Shape2D
type Ellipse struct {
	X, Y, W, H float32
}

// Construct a new ellipse with all values set to zero
func NewEllipseEmpty() *Ellipse {
	return &Ellipse{}
}

// Copy constructor
// param ellipse Ellipse to construct a copy of.
func NewEllipseCopy(ellipse *Ellipse) *Ellipse {
	return &Ellipse{ellipse.X, ellipse.Y, ellipse.W, ellipse.H}
}

// Constructs a new ellipse
// param x X coordinate
// param y Y coordinate
// param width Width in pixels
// param height Height in pixels
func NewEllipse(x, y, w, h float32) *Ellipse {
	return &Ellipse{x, y, w, h}
}

// Costructs a new ellipse
// param position Position vector
// param width Width in pixels
// param height Height in pixels
func NewEllipseV(position *Vector2, w, h float32) *Ellipse {
	return &Ellipse{position.X, position.Y, w, h}
}

func NewEllipseV2(position, size *Vector2) *Ellipse {
	return &Ellipse{position.X, position.Y, size.X, size.Y}
}

// Constructs a new {@link Ellipse} from the position and radius of a {@link Circle} (since circles are special cases of
// ellipses).
// param circle The circle to take the values of
func NewEllipseCircle(circle *Circle) *Ellipse {
	return &Ellipse{circle.X, circle.Y, circle.Radius, circle.Radius}
}

// Checks whether or not this ellipse contains the given point.
// param x X coordinate
// param y Y coordinate
// return true if this ellipse contains the given point; false otherwise.
func (self *Ellipse) Contains(x, y float32) bool {
	x = x - self.X
	y = y - self.Y

	return (x*x)/(self.W*0.5*self.W*0.5)+(y*y)/(self.H*0.5*self.H*0.5) <= 1.0
}

// Checks whether or not this ellipse contains the given point.
// param point Position vector
// return true if this ellipse contains the given point; false otherwise.
func (self *Ellipse) ContainsV(point *Vector2) bool {
	return self.Contains(point.X, point.Y)
}

// Sets a new position and size for this ellipse.
// param x X coordinate
// param y Y coordinate
// param width Width in pixels
// param height Height in pixels
func (self *Ellipse) Set(x, y, w, h float32) {
	self.X = x
	self.Y = y
	self.W = w
	self.H = h
}

// Sets a new position and size for this ellipse based upon another ellipse.
// param ellipse The ellipse to copy the position and size of.
func (self *Ellipse) SetCopy(ellipse *Ellipse) {
	self.X = ellipse.X
	self.Y = ellipse.Y
	self.W = ellipse.W
	self.H = ellipse.H
}

func (self *Ellipse) SetCircle(circle *Circle) {
	self.X = circle.X
	self.Y = circle.Y
	self.W = circle.Radius
	self.H = circle.Radius
}

func (self *Ellipse) SetV2(position, size *Vector2) {
	self.X = position.X
	self.Y = position.Y
	self.W = size.X
	self.H = size.Y
}

// Sets the x and y-coordinates of ellipse center from a {@link Vector2}.
// param position The position vector
// return this ellipse for chaining
func (self *Ellipse) SetPositionV(position *Vector2) *Ellipse {
	self.X = position.X
	self.Y = position.Y
	return self
}

// Sets the x and y-coordinates of ellipse center
// param x The x-coordinate
// param y The y-coordinate
// return this ellipse for chaining
func (self *Ellipse) SetPosition(x, y float32) *Ellipse {
	self.X = x
	self.Y = y
	return self
}

// Sets the width and height of this ellipse
// param width The width
// param height The height
// return this ellipse for chaining
func (self *Ellipse) SetSize(w, h float32) *Ellipse {
	self.W = w
	self.H = h
	return self
}

// return The area of this {@link Ellipse} as {@link MathUtils#PI} * {@link Ellipse#width} * {@link Ellipse#height}
func (self *Ellipse) Area() float32 {
	return utils.PI * (self.W * self.H) / 4
}

// Approximates the circumference of this {@link Ellipse}. Oddly enough, the circumference of an ellipse is actually difficult
// to compute exactly.
// return The Ramanujan approximation to the circumference of an ellipse if one dimension is at least three times longer than
// the other, else the simpler approximation
func (self *Ellipse) Circumference() float32 {
	a := self.W / 2
	b := self.H / 2
	if a*3 > b || b*3 > a {
		// If one dimension is three times as long as the other...
		return float32(float64(utils.PI*(3*(a+b))) - math.Sqrt(float64((3*a+b)*(a+3*b))))
	} else {
		// We can use the simpler approximation, then
		return float32(float64(utils.PI2) * math.Sqrt(float64((a*a+b*b)/2)))
	}
}

// func (self *Ellipse) Equals(Object o) bool {
// 	if (o == this) return true;
// 	if (o == null || o.getClass() != self.getClass()) return false;
// 	Ellipse e = (Ellipse)o;
// 	return self.X == e.x && self.Y == e.y && self.W == e.width && self.H == e.h
// }

// public int hashCode () {
// 	final int prime = 53;
// 	int result = 1;
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.H);
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.W);
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.X);
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.Y);
// 	return result;
// }
