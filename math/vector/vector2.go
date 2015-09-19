// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package vector

import (
	"math"

	. "github.com/pyros2097/gdx/math/interpolation"
)

var (
	XV2    = NewVector2(1, 0)
	YV2    = NewVector2(0, 1)
	ZeroV2 = NewVector2(0, 0)
)

// Encapsulates a 2D vector. Allows chaining methods by returning a reference to itself
// Implements IVector
type Vector2 struct {
	x, y float32
}

// Constructs a vector with the given components
func NewVector2(x, y float32) *Vector2 {
	return &Vector2{x, y}
}

// Constructs a new vector at (0,0)
func NewVector2Empty() *Vector2 {
	return &Vector2{}
}

// Constructs a vector from the given vector
func NewVector2Copy(v *Vector2) *Vector2 {
	return &Vector2{v.x, v.y}
}

func (self *Vector2) Copy() *Vector2 {
	return NewVector2Copy(self)
}

func (self *Vector2) Len() float32 {
	return float32(math.Sqrt(float64(self.x*self.x + self.y*self.y)))
}

func (self *Vector2) Len2() float32 {
	return self.x*self.x + self.y*self.y
}

func LenD(x, y float32) float32 {
	return float32(math.Sqrt(float64(x*x + y*y)))
}

func LenD2(x, y float32) float32 {
	return x*x + y*y
}

func (self *Vector2) SetV(v *Vector2) *Vector2 {
	self.x = v.x
	self.y = v.y
	return self
}

// Sets the components of this vector
func (self *Vector2) Set(x, y float32) *Vector2 {
	self.x = x
	self.y = y
	return self
}

func (self *Vector2) SubV(v *Vector2) *Vector2 {
	x -= v.x
	y -= v.y
	return self
}

// Substracts the other vector from this vector.
func (self *Vector2) Sub(x, y float32) *Vector2 {
	self.x -= x
	self.y -= y
	return self
}

func (self *Vector2) Nor() *Vector2 {
	length := self.Len()
	if length != 0 {
		self.x /= length
		self.y /= length
	}
	return self
}

func (self *Vector2) AddV(v *Vector2) *Vector2 {
	x += v.x
	y += v.y
	return self
}

// Adds the given components to this vector
func (self *Vector2) Add(x, y float32) {
	self.x += x
	self.y += y
	return self
}

func DotV2(x1, y1, x2, y2 float32) float32 {
	return x1*x2 + y1*y2
}

func (self *Vector2) DotV(v *Vector2) float32 {
	return self.x*v.x + self.y*v.y
}

func (self *Vector2) Dot(ox, oy float32) float32 {
	return self.x*ox + self.y*oy
}

func (self *Vector2) SclScalar(scalar float32) *Vector2 {
	x *= scalar
	y *= scalar
	return self
}

// Multiplies this vector by a scalar
func (self *Vector2) Scl(x, y float32) *Vector2 {
	self.x *= x
	self.y *= y
	return self
}

func (self *Vector2) SclV(v *Vector2) *Vector2 {
	self.x *= v.x
	self.y *= v.y
	return self
}

func (self *Vector2) MulAddScalar(v *Vector2, scalar float32) *Vector2 {
	self.x += v.x * scalar
	self.y += v.y * scalar
	return self
}

func (self *Vector2) MulAdd(v *Vector2, mulVec *Vector2) *Vector2 {
	self.x += v.x * mulVec.x
	self.y += v.y * mulVec.y
	return self
}

func DstV2(x1, y1, x2, y2 float32) {
	x_d := x2 - x1
	y_d := y2 - y1
	return float32(math.Sqrt(float64(x_d*x_d + y_d*y_d)))
}

func (self *Vector2) DstV(v *Vector2) float32 {
	x_d := v.x - self.x
	y_d := v.y - self.y
	return float32(math.Sqrt(float64(x_d*x_d + y_d*y_d)))
}

func (self *Vector2) Dst(x, y float32) float32 {
	x_d := x - self.x
	y_d := y - self.y
	return float32(math.Sqrt(float64(x_d*x_d + y_d*y_d)))
}

func Dst2V2(x1, y1, x2, y2 float32) {
	x_d := x2 - x1
	y_d := y2 - y1
	return x_d*x_d + y_d*y_d
}

func (self *Vector2) Dst2V(v *Vector2) float32 {
	x_d := v.x - self.x
	y_d := v.y - self.y
	return x_d*x_d + y_d*y_d
}

func (self *Vector2) Dst2(x, y float32) float32 {
	x_d := x - self.x
	y_d := y - self.y
	return x_d*x_d + y_d*y_d
}

func (self *Vector2) Limit(limit float32) *Vector2 {
	return self.Limit2(limit * limit)
}

func (self *Vector2) Limit2(limit2 float32) *Vector2 {
	len2 := self.Len2()
	if len2 > limit2 {
		return self.Scl(float32(Math.sqrt(float64(limit2 / len2))))
	}
	return self
}

func (self *Vector2) Clamp(min, max float32) *Vector2 {
	len2 := self.Len2()
	if len2 == 0 {
		return self
	}
	max2 := max * max
	if len2 > max2 {
		return self.Scl(float32(Math.sqrt(float64(max2 / len2))))
	}
	min2 := min * min
	if len2 < min2 {
		return self.Scl(float32(Math.sqrt(float64(min2 / len2))))
	}
	return self
}

func (self *Vector2) SetLength(length float32) *Vector2 {
	return self.SetLength2(length * length)
}

func (self *Vector2) SetLength2(len2 float32) {
	oldLen2 := self.Len2()
	if oldLen2 == 0 || oldLen2 == len2 {
		return self
	}
	return self.Scl(float32(Math.sqrt(float64(len2 / oldLen2))))
}

func (self *Vector2) String() string {
	return "[" + x + ":" + y + "]"
}

// Left-multiplies this vector by the given matrix
func (self *Vector2) Mul(mat Matrix3) *Vector2 {
	self.x = self.x*mat.val[0] + self.y*mat.val[3] + mat.val[6]
	self.y = self.x*mat.val[1] + self.y*mat.val[4] + mat.val[7]
	return self
}

// Calculates the 2D cross product between this and the given vector
func (self *Vector2) CrsVector2(v *Vector2) float32 {
	return self.x*v.y - self.y*v.x
}

// Calculates the 2D cross product between this and the given vector.
func (self *Vector2) Crs(x, y float32) float32 {
	return self.x*y - self.y*x
}

// @return the angle in degrees of this vector (point) relative to the x-axis.
// Angles are towards the positive y-axis (typically counter-clockwise) and between 0 and 360.
func (self *Vector2) Angle() float32 {
	angle := float32(math.Atan2(float64(self.y), float64(self.x)) * RadiansToDegrees)
	if angle < 0 {
		angle += 360
	}
	return angle
}

// @return the angle in degrees of this vector (point) relative to the given vector.
// Angles are towards the positive y-axis (typically counter-clockwise.) between -180 and +180
func (self *Vector2) AngleVector2(reference *Vector2) float32 {
	return float32(Math.atan2(float64(self.Crs(reference), self.Dot(reference))) * RadiansToDegrees)
}

// @return the angle in radians of this vector (point) relative to the x-axis.
// Angles are towards the positive y-axis. (typically counter-clockwise)
func (self *Vector2) AngleRad() float32 {
	return float32(Math.atan2(float64(self.y), float64(self.x)))
}

// @return the angle in radians of this vector (point) relative to the given vector.
// Angles are towards the positive y-axis. (typically counter-clockwise.)
func (self *Vector2) AngleRadVector2(reference *Vector2) float32 {
	return float32(Math.atan2(float64(self.Crs(reference)), float64(self.Dot(reference))))
}

// Sets the angle of the vector in degrees relative to the x-axis, towards the positive y-axis (typically counter-clockwise).
// @param degrees The angle in degrees to set.
func (self *Vector2) SetAngle(degrees float32) *Vector2 {
	return self.SetAngleRad(degrees * DegreesToRadians)
}

// Sets the angle of the vector in radians relative to the x-axis, towards the positive y-axis (typically counter-clockwise).
// @param radians The angle in radians to set.
func (self *Vector2) SetAngleRad(radians float32) *Vector2 {
	self.Set(self.Len(), 0)
	self.RotateRad(radians)
	return self
}

// Rotates the *Vector2 by the given angle, counter-clockwise assuming the y-axis points up.
// @param degrees the angle in degrees */
func (self *Vector2) Rotate(degrees float32) *Vector2 {
	return self.RotateRad(degrees * DegreesToRadians)
}

// Rotates the *Vector2 by the given angle, counter-clockwise assuming the y-axis points up.
// @param radians the angle in radians
func (self *Vector2) RotateRad(radians float32) *Vector2 {
	cos := float32(math.Cos(float64(radians)))
	sin := float32(math.Sin(float64(radians)))
	self.x = self.x*cos - self.y*sin
	self.y = self.x*sin + self.y*cos
	return self
}

// Rotates the *Vector2 by 90 degrees in the specified direction,
// where >= 0 is counter-clockwise and < 0 is clockwise.
func (self *Vector2) Rotate90(dir int) *Vector2 {
	x := self.x
	if dir >= 0 {
		self.x = -self.y
		self.y = self.x
	} else {
		self.x = self.y
		self.y = -x
	}
	return self
}

func (self *Vector2) Lerp(target *Vector2, alpha float32) *Vector2 {
	invAlpha := 1.0 - alpha
	self.x = (self.x * invAlpha) + (target.x * alpha)
	self.y = (self.y * invAlpha) + (target.y * alpha)
	return self
}

func (self *Vector2) Interpolate(target *Vector2, alpha float32, interpolation Interpolation) *Vector2 {
	return self.Lerp(target, interpolation(alpha))
}

func (self *Vector2) HashCode() int {
	// final int prime = 31;
	// int result = 1;
	// result = prime * result + NumberUtils.floatToIntBits(x);
	// result = prime * result + NumberUtils.floatToIntBits(y);
	// return result;
	return 0
}

func (self *Vector2) Equals(v *Vector2) bool {
	if self == v {
		return true
	}
	if v == nil {
		return false
	}
	// if (NumberUtils.floatToIntBits(x) != NumberUtils.floatToIntBits(other.x)) return false;
	// if (NumberUtils.floatToIntBits(y) != NumberUtils.floatToIntBits(other.y)) return false;
	return true
}

func (self *Vector2) EpsilonEqualsV(other *Vector2, epsilon float32) bool {
	if other == nil {
		return false
	}
	if math.Abs(other.x-self.x) > epsilon {
		return false
	}
	if math.Abs(other.y-self.y) > epsilon {
		return false
	}
	return true
}

// Compares this vector with the other vector, using the supplied epsilon for fuzzy equality testing.
// @return whether the vectors are the same
func (self *Vector2) EpsilonEquals(x, y, epsilon float32) bool {
	if math.Abs(x-self.x) > epsilon {
		return false
	}
	if math.Abs(y-self.y) > epsilon {
		return false
	}
	return true
}

func (self *Vector2) IsUnit() bool {
	return isUnit(0.000000001)
}

func (self *Vector2) IsUnitMargin(margin float32) bool {
	return math.Abs(self.Len2()-1) < margin
}

func (self *Vector2) IsZero() bool {
	return self.x == 0 && self.y == 0
}

func (self *Vector2) IsZeroMargin(margin float32) bool {
	return self.Len2() < margin
}

func (self *Vector2) IsOnLine(other *Vector2) bool {
	return IsZero(x*other.y - y*other.x)
}

func (self *Vector2) IsOnLineEpsilon(other *Vector2, epsilon float32) bool {
	return IsZero(x*other.y-y*other.x, epsilon)
}

func (self *Vector2) IsCollinearEpsilon(other *Vector2, epsilon float32) bool {
	return self.IsOnLine(other, epsilon) && self.Dot(other) > 0
}

func (self *Vector2) IsCollinear(other *Vector2) bool {
	return self.IsOnLine(other) && self.Dot(other) > 0
}

func (self *Vector2) IsCollinearOppositeEpsilon(other *Vector2, epsilon float32) bool {
	return self.IsOnLine(other, epsilon) && self.Dot(other) < 0
}

func (self *Vector2) IsCollinearOpposite(other *Vector2) bool {
	return self.IsOnLine(other) && self.Dot(other) < 0
}

func (self *Vector2) IsPerpendicular(v *Vector2) bool {
	return IsZero(self.Dot(v))
}

func (self *Vector2) IsPerpendicularEpsilon(vector *Vector2, epsilon float32) bool {
	return IsZero(self.Dot(vector), epsilon)
}

func (self *Vector2) HasSameDirection(v *Vector2) bool {
	return self.Dot(vector) > 0
}

func (self *Vector2) HasOppositeDirection(v *Vector2) bool {
	return self.Dot(vector) < 0
}

func (self *Vector2) SetZero() *Vector2 {
	self.x = 0
	self.y = 0
	return self
}
