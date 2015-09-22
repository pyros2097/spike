// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package shape

import (
	"math"

	. "github.com/pyros2097/spike/math/vector"
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
	X, Y, W, H float32
}

// Constructs a new rectangle with all values set to zero
func NewRectangleEmpty() *Rectangle {
	return &Rectangle{}
}

// Constructs a new rectangle with the given corner point in the bottom left and dimensions.
func NewRectangle(x, y, w, h float32) *Rectangle {
	return &Rectangle{x, y, w, h}
}

// Constructs a rectangle based on the given rectangle
func NewRectangleCopy(rect *Rectangle) *Rectangle {
	return &Rectangle{rect.X, rect.Y, rect.W, rect.H}
}

func (self *Rectangle) Set(x, y, w, h float32) *Rectangle {
	self.X = x
	self.Y = y
	self.W = w
	self.H = h
	return self
}

// Sets the values of the given rectangle to this rectangle.
func (self *Rectangle) SetR(rect *Rectangle) *Rectangle {
	self.X = rect.X
	self.Y = rect.Y
	self.W = rect.W
	self.H = rect.H
	return self
}

// returns the x-coordinate of the bottom left corner
func (self *Rectangle) GetX() float32 {
	return self.X
}

// Sets the x-coordinate of the bottom left corner
// param x The x-coordinate
// returns this rectangle for chaining
func (self *Rectangle) SetX(x float32) *Rectangle {
	self.X = x
	return self
}

// returns the y-coordinate of the bottom left corner
func (self *Rectangle) GetY() float32 {
	return self.Y
}

// Sets the y-coordinate of the bottom left corner
// param y The y-coordinate
// returns this rectangle for chaining
func (self *Rectangle) SetY(y float32) *Rectangle {
	self.Y = y
	return self
}

// returns the width
func (self *Rectangle) GetW() float32 {
	return self.W
}

// Sets the width of this rectangle
func (self *Rectangle) SetW(w float32) *Rectangle {
	self.W = w
	return self
}

// returns the height
func (self *Rectangle) GetH() float32 {
	return self.W
}

// Sets the height of this rectangle
func (self *Rectangle) SetH(h float32) *Rectangle {
	self.H = h
	return self
}

// return the Vector2 with coordinates of this rectangle
func (self *Rectangle) GetPosition(position *Vector2) *Vector2 {
	return position.Set(self.X, self.Y)
}

// Sets the x and y-coordinates of the bottom left corner from vector
func (self *Rectangle) SetPositionV(position *Vector2) *Rectangle {
	self.X = position.X
	self.Y = position.Y
	return self
}

// Sets the x and y-coordinates of the bottom left corner
func (self *Rectangle) SetPosition(x, y float32) *Rectangle {
	self.X = x
	self.Y = y
	return self
}

// Sets the width and height of this rectangle
func (self *Rectangle) SetSize(width, height float32) *Rectangle {
	self.W = width
	self.H = height
	return self
}

// Sets the squared size of this rectangle
func (self *Rectangle) SetSizeXY(sizeXY float32) *Rectangle {
	self.W = sizeXY
	self.H = sizeXY
	return self
}

// returns the Vector2 with size of this rectangle
func (self *Rectangle) GetSize(size *Vector2) *Vector2 {
	return size.Set(self.W, self.H)
}

// returns whether the point is contained in the rectangle
func (self *Rectangle) Contains(x, y float32) bool {
	return self.X <= x && self.X+self.W >= x && self.Y <= y && self.Y+self.H >= y
}

// returns whether the point is contained in the rectangle
func (self *Rectangle) ContainsV(point *Vector2) bool {
	return self.Contains(point.X, point.Y)
}

// returns whether the other rectangle is contained in this rectangle.
func (self *Rectangle) ContainsRect(r *Rectangle) bool {
	xmin := r.X
	xmax := xmin + r.W

	ymin := r.Y
	ymax := ymin + r.H

	return ((xmin > self.X && xmin < self.X+self.W) && (xmax > self.X && xmax < self.X+self.W)) &&
		((ymin > self.Y && ymin < self.Y+self.H) && (ymax > self.Y && ymax < self.Y+self.H))
}

// returns whether this rectangle overlaps the other rectangle.
func (self *Rectangle) Overlaps(r *Rectangle) bool {
	return self.X < r.X+r.W && self.X+self.W > r.X && self.Y < r.Y+r.H && self.Y+self.H > r.Y
}

// Merges this rectangle with the other rectangle. The rectangle should not have negative width or negative height.
func (self *Rectangle) MergeRect(rect *Rectangle) *Rectangle {
	minX := float32(math.Min(float64(self.X), float64(rect.X)))
	maxX := float32(math.Max(float64(self.X+self.W), float64(rect.X+rect.W)))
	self.X = minX
	self.W = maxX - minX

	minY := float32(math.Min(float64(self.Y), float64(rect.Y)))
	maxY := float32(math.Max(float64(self.Y+self.H), float64(rect.Y+rect.H)))
	self.Y = minY
	self.H = maxY - minY

	return self
}

// Merges this rectangle with a point. The rectangle should not have negative width or negative height.
func (self *Rectangle) Merge(x, y float32) *Rectangle {
	minX := float32(math.Min(float64(self.X), float64(x)))
	maxX := float32(math.Max(float64(self.X+self.W), float64(x)))
	self.X = minX
	self.W = maxX - minX

	minY := float32(math.Min(float64(self.Y), float64(y)))
	maxY := float32(math.Max(float64(self.Y+self.H), float64(y)))
	self.Y = minY
	self.H = maxY - minY

	return self
}

// Merges this rectangle with a point. The rectangle should not have negative width or negative height.
func (self *Rectangle) MergeVec(vec *Vector2) *Rectangle {
	return self.Merge(vec.X, vec.Y)
}

// Merges this rectangle with a list of points. The rectangle should not have negative width or negative height.
func (self *Rectangle) MergeVecs(vecs []*Vector2) *Rectangle {
	minX := self.X
	maxX := self.X + self.W
	minY := self.Y
	maxY := self.Y + self.H // TODO check ++i
	for i := 0; i < len(vecs); i++ {
		v := vecs[i]
		minX = float32(math.Min(float64(minX), float64(v.X)))
		maxX = float32(math.Max(float64(maxX), float64(v.X)))
		minY = float32(math.Min(float64(minY), float64(v.Y)))
		maxY = float32(math.Max(float64(maxY), float64(v.Y)))
	}
	self.X = minX
	self.W = maxX - minX
	self.Y = minY
	self.H = maxY - minY
	return self
}

// Calculates the aspect ratio ( width / height ) of this rectangle
// returns the aspect ratio of this rectangle. Returns Float.NaN if height is 0 to avoid ArithmeticException
func (self *Rectangle) GetAspectRatio() float32 {
	if self.H == 0 {
		return float32(math.NaN())
	}
	return self.W / self.H
}

// Calculates the center of the rectangle. Results are located in the given Vector2
// returns the given vector with results stored inside
func (self *Rectangle) GetCenter(vector *Vector2) *Vector2 {
	vector.X = self.X + self.W/2
	vector.Y = self.Y + self.H/2
	return vector
}

// Moves this rectangle so that its center point is located at a given position
// returns this for chaining
func (self *Rectangle) SetCenter(x, y float32) *Rectangle {
	self.SetPosition(self.X-self.W/2, self.Y-self.H/2)
	return self
}

// Moves this rectangle so that its center point is located at a given position
// returns this for chaining
func (self *Rectangle) SetCenterV(position *Vector2) *Rectangle {
	self.SetPosition(position.X-self.W/2, position.Y-self.H/2)
	return self
}

// Fits this rectangle around another rectangle while maintaining aspect ratio. This scales and centers the rectangle to the
// other rectangle (e.g. Having a camera translate and scale to show a given area)
// param rect the other rectangle to fit this rectangle around
// returns this rectangle for chaining
func (self *Rectangle) FitOutside(rect *Rectangle) *Rectangle {
	ratio := self.GetAspectRatio()

	if ratio > rect.GetAspectRatio() {
		// Wider than tall
		self.SetSize(rect.H*ratio, rect.H)
	} else {
		// Taller than wide
		self.SetSize(rect.W, rect.W/ratio)
	}

	self.SetPosition((rect.X+rect.W/2)-self.W/2, (rect.Y+rect.H/2)-self.H/2)
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
		self.SetSize(rect.H*ratio, rect.H)
	} else {
		// Wider than tall
		self.SetSize(rect.W, rect.W/ratio)
	}

	self.SetPosition((rect.X+rect.W/2)-self.W/2, (rect.Y+rect.H/2)-self.H/2)
	return self
}

func (self *Rectangle) Area() float32 {
	return self.W * self.H
}

func (self *Rectangle) Perimeter() float32 {
	return 2 * (self.W + self.H)
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
	// if (NumberUtils.floatToRawIntBits(height) != NumberUtils.floatToRawIntBits(other.H)) return false;
	// if (NumberUtils.floatToRawIntBits(width) != NumberUtils.floatToRawIntBits(other.W)) return false;
	// if (NumberUtils.floatToRawIntBits(x) != NumberUtils.floatToRawIntBits(other.x)) return false;
	// if (NumberUtils.floatToRawIntBits(y) != NumberUtils.floatToRawIntBits(other.y)) return false;
	return true
}

func (self *Rectangle) String() string {
	return ""
	// return x + "," + y + "," + width + "," + height
}
