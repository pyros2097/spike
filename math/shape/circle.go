// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package shape

import (
	"github.com/pyros2097/gdx/math/vector"
)

// A convenient 2D circle class.
// Implements Shape2D
type Circle struct {
	x, y, radius float32
}

// Constructs a new circle with all values set to zero
func NewCircleEmpty() *Circle {
	return &Circle{}
}

// Constructs a new circle with the given X and Y coordinates and the given radius.
func NewCircle(x, y, radius float32) *Circle {
	return &Circle{x, y, radius}
}

// Constructs a new circle using a given {@link Vector2} that contains the desired X and Y coordinates, and a given radius.
// position The position {@link Vector2}.
// radius The radius
func NewCircleV(position *Vector2, radius float32) *Circle {
	return &Circle{position.x, position.y, radius}
}

// Copy constructor
// circle The circle to construct a copy of.
func NewCircleC(c *Circle) *Circle {
	return &Circle{c.x, c.y, radius}
}

// Creates a new {@link Circle} in terms of its center and a point on its edge.
// center The center of the new circle
// edge Any point on the edge of the given circle
func NewCircleV2(center, edge *Vector2) *Circle {
	return &Circle{center.x, center.y, LenV2(center.x-edge.x, center.y-edge.y)}
}

// Sets a new location and radius for this circle.
func (self *Circle) Set(x, y, radius float32) {
	self.x = x
	self.y = y
	self.radius = radius
}

// Sets a new location and radius for this circle.
func (self *Circle) SetV(position *Vector2, radius float32) {
	self.x = position.x
	self.y = position.y
	self.radius = radius
}

// Sets a new location and radius for this circle, based upon another circle.
// circle The circle to copy the position and radius of.
func (self *Circle) SetC(circle *Circle) {
	self.x = circle.x
	self.y = circle.y
	self.radius = circle.radius
}

// Sets this {@link Circle}'s values in terms of its center and a point on its edge.
// center The new center of the circle
// edge Any point on the edge of the given circle
func (self *Circle) SetV2(center, edge *Vector2) {
	self.x = center.x
	self.y = center.y
	self.radius = LenV2(center.x-edge.x, center.y-edge.y)
}

// Sets the x and y-coordinates of circle center from vector
// position The position vector
func (self *Circle) SetPosition(position *Vector2) {
	self.x = position.x
	self.y = position.y
}

// Sets the x and y-coordinates of circle center
// x The x-coordinate
// y The y-coordinate
func (self *Circle) setPosition(x, y float32) {
	self.x = x
	self.y = y
}

// Sets the x-coordinate of circle center
// x The x-coordinate
func (self *Circle) setX(float x) {
	self.x = x
}

// Sets the y-coordinate of circle center
// y The y-coordinate
func (self *Circle) setY(float y) {
	self.y = y
}

// Sets the radius of circle
// radius The radius
func (self *Circle) SetRadius(radius float32) {
	self.radius = radius
}

// Checks whether or not this circle contains a given point.
// x X coordinate
// y Y coordinate
// true if this circle contains the given point.
func (self *Circle) Contains(x, y float32) bool {
	x = self.x - x
	y = self.y - y
	return x*x+y*y <= radius*radius
}

// Checks whether or not this circle contains a given point.
// point The {@link Vector2} that contains the point coordinates.
// true if this circle contains this point; false otherwise.
func (self *Circle) ContainsV(point *Vector2) {
	dx := x - point.x
	dy := y - point.y
	return dx*dx+dy*dy <= radius*radius
}

// c the other {@link Circle}
// whether this circle contains the other circle.
func (self *Circle) ContainsC(c *Circle) bool {
	radiusDiff := radius - c.radius
	if radiusDiff < 0 {
		return false // Can't contain bigger circle
	}
	dx := x - c.x
	dy := y - c.y
	dst := dx*dx + dy*dy
	radiusSum := radius + c.radius
	return (!(radiusDiff*radiusDiff < dst) && (dst < radiusSum*radiusSum))
}

// // c the other {@link Circle}
// // whether this circle overlaps the other circle.
func (self *Circle) Overlaps(c *Circle) bool {
	dx := x - c.x
	dy := y - c.y
	distance := dx*dx + dy*dy
	radiusSum := radius + c.radius
	return distance < radiusSum*radiusSum
}

// Returns a {@link String} representation of this {@link Circle} of the form {@code x,y,radius}.
func (self *Circle) String() {
	return x + "," + y + "," + radius
}

// The circumference of this circle (as 2 * {@link MathUtils#PI2}) * {@code radius}
func (self *Circle) Circumference() float32 {
	return self.radius * PI2
}

// The area of this circle (as {@link MathUtils#PI} * radius * radius).
func (self *Circle) Area() float32 {
	return self.radius * self.radius * PI
}

// @Override
// func (self *Circle) equals (Object o) {
// 	if (o == this) return true;
// 	if (o == null || o.getClass() != self.getClass()) return false;
// 	Circle c = (Circle)o;
// 	return self.x == c.x && self.y == c.y && self.radius == c.radius;
// }

// @Override
// public int hashCode () {
// 	final int prime = 41;
// 	int result = 1;
// 	result = prime * result + NumberUtils.floatToRawIntBits(radius);
// 	result = prime * result + NumberUtils.floatToRawIntBits(x);
// 	result = prime * result + NumberUtils.floatToRawIntBits(y);
// 	return result;
// }
