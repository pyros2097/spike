// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package shape

import (
	"math"

	"github.com/pyros2097/gdx/math/utils"
	. "github.com/pyros2097/gdx/math/vector"
)

// A convenient 2D ellipse class, based on the circle class
// Implements Shape2D
type Ellipse struct {
	x, y, w, h float32
}

// Construct a new ellipse with all values set to zero
func NewEllipseEmpty() *Ellipse {
	return &Ellipse{}
}

// Copy constructor
// param ellipse Ellipse to construct a copy of.
func NewEllipseCopy(ellipse *Ellipse) *Ellipse {
	return &Ellipse{ellipse.x, ellipse.y, ellipse.w, ellipse.h}
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
	return &Ellipse{position.x, position.y, w, h}
}

func NewEllipseV2(position, size *Vector2) *Ellipse {
	return &Ellipse{position.x, position.y, size.x, size.y}
}

// Constructs a new {@link Ellipse} from the position and radius of a {@link Circle} (since circles are special cases of
// ellipses).
// param circle The circle to take the values of
func NewEllipseCircle(circle *Circle) *Ellipse {
	return &Ellipse{circle.x, circle.y, circle.radius, circle.radius}
}

// Checks whether or not this ellipse contains the given point.
// param x X coordinate
// param y Y coordinate
// return true if this ellipse contains the given point; false otherwise.
func (self *Ellipse) Contains(x, y float32) bool {
	x = x - self.x
	y = y - self.y

	return (x*x)/(self.w*0.5*self.w*0.5)+(y*y)/(self.h*0.5*self.h*0.5) <= 1.0
}

// Checks whether or not this ellipse contains the given point.
// param point Position vector
// return true if this ellipse contains the given point; false otherwise.
func (self *Ellipse) ContainsV(point *Vector2) bool {
	return self.Contains(point.x, point.y)
}

// Sets a new position and size for this ellipse.
// param x X coordinate
// param y Y coordinate
// param width Width in pixels
// param height Height in pixels
func (self *Ellipse) Set(x, y, w, h float32) {
	self.x = x
	self.y = y
	self.w = w
	self.h = h
}

// Sets a new position and size for this ellipse based upon another ellipse.
// param ellipse The ellipse to copy the position and size of.
func (self *Ellipse) SetCopy(ellipse *Ellipse) {
	self.x = ellipse.x
	self.y = ellipse.y
	self.w = ellipse.w
	self.h = ellipse.h
}

func (self *Ellipse) SetCircle(circle *Circle) {
	self.x = circle.x
	self.y = circle.y
	self.w = circle.radius
	self.h = circle.radius
}

func (self *Ellipse) SetV2(position, size *Vector2) {
	self.x = position.x
	self.y = position.y
	self.w = size.x
	self.h = size.y
}

// Sets the x and y-coordinates of ellipse center from a {@link Vector2}.
// param position The position vector
// return this ellipse for chaining
func (self *Ellipse) SetPositionV(position *Vector2) *Ellipse {
	self.x = position.x
	self.y = position.y
	return self
}

// Sets the x and y-coordinates of ellipse center
// param x The x-coordinate
// param y The y-coordinate
// return this ellipse for chaining
func (self *Ellipse) SetPosition(x, y float32) *Ellipse {
	self.x = x
	self.y = y
	return self
}

// Sets the width and height of this ellipse
// param width The width
// param height The height
// return this ellipse for chaining
func (self *Ellipse) SetSize(w, h float32) *Ellipse {
	self.w = w
	self.h = h
	return self
}

// return The area of this {@link Ellipse} as {@link MathUtils#PI} * {@link Ellipse#width} * {@link Ellipse#height}
func (self *Ellipse) Area() float32 {
	return utils.PI * (self.w * self.h) / 4
}

// Approximates the circumference of this {@link Ellipse}. Oddly enough, the circumference of an ellipse is actually difficult
// to compute exactly.
// return The Ramanujan approximation to the circumference of an ellipse if one dimension is at least three times longer than
// the other, else the simpler approximation
func (self *Ellipse) Circumference() float32 {
	a := self.w / 2
	b := self.h / 2
	if a*3 > b || b*3 > a {
		// If one dimension is three times as long as the other...
		return float32(float64(utils.PI*(3*(a+b)) - math.Sqrt(float64((3*a+b)*(a+3*b)))))
	} else {
		// We can use the simpler approximation, then
		return float32(float64(utils.PI2) * math.Sqrt(float64((a*a+b*b)/2)))
	}
}

// func (self *Ellipse) Equals(Object o) bool {
// 	if (o == this) return true;
// 	if (o == null || o.getClass() != self.getClass()) return false;
// 	Ellipse e = (Ellipse)o;
// 	return self.x == e.x && self.y == e.y && self.w == e.width && self.h == e.h
// }

// public int hashCode () {
// 	final int prime = 53;
// 	int result = 1;
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.h);
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.w);
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.x);
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.y);
// 	return result;
// }
