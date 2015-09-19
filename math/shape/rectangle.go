// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package shape

import (
	"math"

	. "github.com/pyros2097/gdx/math/vector"
)

var (
	// Static temporary rectangle. Use with care! Use only when sure other code will not also use self.
	Tmp = NewRectangleEmpty()

	// Static temporary rectangle. Use with care! Use only when sure other code will not also use self.
	Tmp2 = NewRectangleEmpty()
)

// Encapsulates a 2D rectangle defined by its corner point in the bottom left and its extents in x (width) and y (height).
// Implements Shape2D
type Rectangle struct {
	x, y, w, h float32
}

// Constructs a new rectangle with all values set to zero
func NewRectangleEmpty() *Rectangle {
	return &Rectangle{}
}

// Constructs a new rectangle with the given corner point in the bottom left and dimensions.
func NewRectangle(x, y, w, h float32) {
	return &Rectangle{x, y, w, h}
}

// Constructs a rectangle based on the given rectangle
func NewRectangleCopy(rect *Rectangle) {
	return &Rectangle{rect.x, rect.y, rect.w, rect.h}
}

func (self *Rectangle) Set(x, y, w, h float32) *Rectangle {
	self.x = x
	self.y = y
	self.w = w
	self.h = h
	return self
}

// Sets the values of the given rectangle to this rectangle.
func (self *Rectangle) SetR(rect *Rectangle) *Rectangle {
	self.x = rect.x
	self.y = rect.y
	self.w = rect.w
	self.h = rect.h
	return self
}

// returns the x-coordinate of the bottom left corner
func (self *Rectangle) GetX() float32 {
	return self.x
}

// Sets the x-coordinate of the bottom left corner
// param x The x-coordinate
// returns this rectangle for chaining
func (self *Rectangle) SetX(x float32) *Rectangle {
	self.x = x
	return self
}

// returns the y-coordinate of the bottom left corner
func (self *Rectangle) GetY() float32 {
	return self.y
}

// Sets the y-coordinate of the bottom left corner
// param y The y-coordinate
// returns this rectangle for chaining
func (self *Rectangle) SetY(y float32) *Rectangle {
	self.y = y
	return self
}

// returns the width
func (self *Rectangle) GetW() float32 {
	return self.w
}

// Sets the width of this rectangle
func (self *Rectangle) SetW(w float32) *Rectangle {
	self.w = w
	return self
}

// returns the height
func (self *Rectangle) GetH() float32 {
	return self.w
}

// Sets the height of this rectangle
func (self *Rectangle) SetH(h float32) *Rectangle {
	self.h = h
	return self
}

// return the Vector2 with coordinates of this rectangle
func (self *Rectangle) GetPosition(position *Vector2) *Vector2 {
	return position.Set(self.x, self.y)
}

// Sets the x and y-coordinates of the bottom left corner from vector
func (self *Rectangle) SetPositionV(position *Vector2) *Rectangle {
	self.x = position.x
	self.y = position.y
	return self
}

// Sets the x and y-coordinates of the bottom left corner
func (self *Rectangle) SetPosition(x, y float32) *Rectangle {
	self.x = x
	self.y = y
	return self
}

// Sets the width and height of this rectangle
func (self *Rectangle) SetSize(width, height float32) *Rectangle {
	self.width = width
	self.height = height
	return self
}

// Sets the squared size of this rectangle
func (self *Rectangle) SetSizeXY(sizeXY float32) *Rectangle {
	self.width = sizeXY
	self.height = sizeXY
	return self
}

// returns the Vector2 with size of this rectangle
func (self *Rectangle) GetSize(size *Vector2) *Vector2 {
	return size.Set(self.w, self.h)
}

// returns whether the point is contained in the rectangle
func (self *Rectangle) Contains(x, y float32) bool {
	return self.x <= x && self.x+self.width >= x && self.y <= y && self.y+self.height >= y
}

// returns whether the point is contained in the rectangle
func (self *Rectangle) ContainsV(point *Vector2) bool {
	return self.Contains(point.x, point.y)
}

// returns whether the other rectangle is contained in this rectangle.
func (self *Rectangle) ContainsRect(rectangle *Rectangle) bool {
	xmin := rectangle.x
	xmax := xmin + rectangle.width

	ymin := rectangle.y
	ymax := ymin + rectangle.height

	return ((xmin > self.x && xmin < self.x+self.w) && (xmax > self.x && xmax < self.x+w)) &&
		((ymin > self.y && ymin < self.y+h) && (ymax > self.y && ymax < self.y+h))
}

// returns whether this rectangle overlaps the other rectangle.
func (self *Rectangle) Overlaps(r *Rectangle) bool {
	return self.x < r.x+r.width && self.x+self.w > r.x && self.y < r.y+r.height && self.y+self.h > r.y
}

// Merges this rectangle with the other rectangle. The rectangle should not have negative width or negative height.
func (self *Rectangle) MergeRect(rect *Rectangle) *Rectangle {
	minX := float32(math.Min(self.x, rect.x))
	maxX := float32(math.Max(self.x+self.w, rect.x+rect.width))
	self.x = minX
	self.w = maxX - minX

	minY := float32(math.Min(self.y, rect.y))
	maxY := float32(math.Max(self.y+self.h, rect.y+rect.height))
	self.y = minY
	self.h = maxY - minY

	return self
}

// Merges this rectangle with a point. The rectangle should not have negative width or negative height.
func (self *Rectangle) Merge(x, y float32) *Rectangle {
	minX := float32(math.Min(self.x, x))
	maxX := float32(math.Max(self.x+self.w, x))
	self.x = minX
	self.w = maxX - minX

	minY := float32(math.Min(self.y, y))
	maxY := float32(math.Max(self.y+self.h, y))
	self.y = minY
	self.h = maxY - minY

	return self
}

// Merges this rectangle with a point. The rectangle should not have negative width or negative height.
func (self *Rectangle) MergeVec(vec *Vector2) *Rectangle {
	return self.Merge(vec.x, vec.y)
}

// Merges this rectangle with a list of points. The rectangle should not have negative width or negative height.
func (self *Rectangle) MergeVecs(vecs []*Vector2) *Rectangle {
	minX := x
	maxX := x + width
	minY := y
	maxY := y + height // TODO check ++i
	for i := 0; i < len(vecs); i++ {
		v := vecs[i]
		minX = float32(math.Min(minX, v.x))
		maxX = float32(math.Max(maxX, v.x))
		minY = float32(math.Min(minY, v.y))
		maxY = float32(math.Max(maxY, v.y))
	}
	self.x = minX
	self.w = maxX - minX
	self.y = minY
	self.h = maxY - minY
	return self
}

// Calculates the aspect ratio ( width / height ) of this rectangle
// returns the aspect ratio of this rectangle. Returns Float.NaN if height is 0 to avoid ArithmeticException
func (self *Rectangle) GetAspectRatio() float32 {
	if height == 0 {
		return float32(math.NaN)
	}
	return self.w / self.h
}

// Calculates the center of the rectangle. Results are located in the given Vector2
// returns the given vector with results stored inside
func (self *Rectangle) GetCenter(vector *Vector2) *Vector2 {
	vector.x = self.x + self.w/2
	vector.y = self.y + self.h/2
	return vector
}

// Moves this rectangle so that its center point is located at a given position
// returns this for chaining
func (self *Rectangle) SetCenter(x, y float32) *Rectangle {
	self.SetPosition(self.x-self.w/2, self.y-self.h/2)
	return self
}

// Moves this rectangle so that its center point is located at a given position
// returns this for chaining
func (self *Rectangle) SetCenterV(position *Vector2) *Rectangle {
	self.SetPosition(position.x-self.w/2, position.y-self.h/2)
	return self
}

// Fits this rectangle around another rectangle while maintaining aspect ratio. This scales and centers the rectangle to the
// other rectangle (e.g. Having a camera translate and scale to show a given area)
// param rect the other rectangle to fit this rectangle around
// returns this rectangle for chaining
func (self *Rectangle) FitOutside(rect *Rectangle) *Rectangle {
	ratio := self.GetAspectRatio()

	if ratio > rect.getAspectRatio() {
		// Wider than tall
		self.SetSize(rect.height*ratio, rect.height)
	} else {
		// Taller than wide
		self.SetSize(rect.width, rect.width/ratio)
	}

	self.SetPosition((rect.x+rect.width/2)-self.w/2, (rect.y+rect.height/2)-self.h/2)
	return self
}

// Fits this rectangle into another rectangle while maintaining aspect ratio. This scales and centers the rectangle to the
// other rectangle (e.g. Scaling a texture within a arbitrary cell without squeezing)
// param rect the other rectangle to fit this rectangle inside
// returns this rectangle for chaining
func (self *Rectangle) FitInside(rect *Rectangle) *Rectangle {
	ratio := self.GetAspectRatio()

	if ratio < rect.GetAspectRatio() {
		// Taller than wide
		self.SetSize(rect.height*ratio, rect.height)
	} else {
		// Wider than tall
		self.SetSize(rect.width, rect.width/ratio)
	}

	self.SetPosition((rect.x+rect.width/2)-self.w/2, (rect.y+rect.height/2)-self.h/2)
	return self
}

func (self *Rectangle) Area() float32 {
	return self.width * self.height
}

func (self *Rectangle) Perimeter() float32 {
	return 2 * (self.width + self.height)
}

func (self *Rectangle) HashCode() int {
	//     final int prime = 31;
	//     int result = 1;
	//     result = prime * result + NumberUtils.floatToRawIntBits(height);
	//     result = prime * result + NumberUtils.floatToRawIntBits(width);
	//     result = prime * result + NumberUtils.floatToRawIntBits(x);
	//     result = prime * result + NumberUtils.floatToRawIntBits(y);
	//     return result;
	return 0
}

func (self *Rectangle) Equals(rect *Rectangle) bool {
	if self == rect {
		return true
	}
	if rect == nil {
		return false
	}
	// if (NumberUtils.floatToRawIntBits(height) != NumberUtils.floatToRawIntBits(other.height)) return false;
	// if (NumberUtils.floatToRawIntBits(width) != NumberUtils.floatToRawIntBits(other.width)) return false;
	// if (NumberUtils.floatToRawIntBits(x) != NumberUtils.floatToRawIntBits(other.x)) return false;
	// if (NumberUtils.floatToRawIntBits(y) != NumberUtils.floatToRawIntBits(other.y)) return false;
	return true
}

func (self *Rectangle) String() string {
	return ""
	// return x + "," + y + "," + width + "," + height
}
