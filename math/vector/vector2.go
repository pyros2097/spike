// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package vector

import (
	"math"

	. "github.com/pyros2097/spike/interpolation"
	"github.com/pyros2097/spike/math/utils"
)

var (
	XV2    = NewVector2(1, 0)
	YV2    = NewVector2(0, 1)
	ZeroV2 = NewVector2(0, 0)
)

// Encapsulates a 2D vector. Allows chaining methods by returning a reference to itself
// Implements IVector
type Vector2 struct {
	X, Y float32
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
	return &Vector2{v.X, v.Y}
}

func (self *Vector2) Copy() *Vector2 {
	return NewVector2Copy(self)
}

func (self *Vector2) Len() float32 {
	return float32(math.Sqrt(float64(self.X*self.X + self.Y*self.Y)))
}

func (self *Vector2) Len2() float32 {
	return self.X*self.X + self.Y*self.Y
}

func LenV2(x, y float32) float32 {
	return float32(math.Sqrt(float64(x*x + y*y)))
}

func Len2V2(x, y float32) float32 {
	return x*x + y*y
}

func (self *Vector2) SetV(v *Vector2) *Vector2 {
	self.X = v.X
	self.Y = v.Y
	return self
}

// Sets the components of this vector
func (self *Vector2) Set(x, y float32) *Vector2 {
	self.X = x
	self.Y = y
	return self
}

func (self *Vector2) SubV(v *Vector2) *Vector2 {
	self.X -= v.X
	self.Y -= v.Y
	return self
}

// Substracts the other vector from this vector.
func (self *Vector2) Sub(x, y float32) *Vector2 {
	self.X -= x
	self.Y -= y
	return self
}

func (self *Vector2) Nor() *Vector2 {
	length := self.Len()
	if length != 0 {
		self.X /= length
		self.Y /= length
	}
	return self
}

func (self *Vector2) AddV(v *Vector2) *Vector2 {
	self.X += v.X
	self.Y += v.Y
	return self
}

// Adds the given components to this vector
func (self *Vector2) Add(x, y float32) *Vector2 {
	self.X += x
	self.Y += y
	return self
}

func DotV2(x1, y1, x2, y2 float32) float32 {
	return x1*x2 + y1*y2
}

func (self *Vector2) DotV(v *Vector2) float32 {
	return self.X*v.X + self.Y*v.Y
}

func (self *Vector2) Dot(ox, oy float32) float32 {
	return self.X*ox + self.Y*oy
}

func (self *Vector2) SclScalar(scalar float32) *Vector2 {
	self.X *= scalar
	self.Y *= scalar
	return self
}

// Multiplies this vector by a scalar
func (self *Vector2) Scl(x, y float32) *Vector2 {
	self.X *= x
	self.Y *= y
	return self
}

func (self *Vector2) SclV(v *Vector2) *Vector2 {
	self.X *= v.X
	self.Y *= v.Y
	return self
}

func (self *Vector2) MulAddScalar(v *Vector2, scalar float32) *Vector2 {
	self.X += v.X * scalar
	self.Y += v.Y * scalar
	return self
}

func (self *Vector2) MulAdd(v *Vector2, mulVec *Vector2) *Vector2 {
	self.X += v.X * mulVec.X
	self.Y += v.Y * mulVec.Y
	return self
}

func DstV2(x1, y1, x2, y2 float32) float32 {
	x_d := x2 - x1
	y_d := y2 - y1
	return float32(math.Sqrt(float64(x_d*x_d + y_d*y_d)))
}

func (self *Vector2) DstV(v *Vector2) float32 {
	x_d := v.X - self.X
	y_d := v.Y - self.Y
	return float32(math.Sqrt(float64(x_d*x_d + y_d*y_d)))
}

func (self *Vector2) Dst(x, y float32) float32 {
	x_d := x - self.X
	y_d := y - self.Y
	return float32(math.Sqrt(float64(x_d*x_d + y_d*y_d)))
}

func Dst2V2(x1, y1, x2, y2 float32) float32 {
	x_d := x2 - x1
	y_d := y2 - y1
	return x_d*x_d + y_d*y_d
}

func (self *Vector2) Dst2V(v *Vector2) float32 {
	x_d := v.X - self.X
	y_d := v.Y - self.Y
	return x_d*x_d + y_d*y_d
}

func (self *Vector2) Dst2(x, y float32) float32 {
	x_d := x - self.X
	y_d := y - self.Y
	return x_d*x_d + y_d*y_d
}

func (self *Vector2) Limit(limit float32) *Vector2 {
	return self.Limit2(limit * limit)
}

func (self *Vector2) Limit2(limit2 float32) *Vector2 {
	len2 := self.Len2()
	if len2 > limit2 {
		return self.SclScalar(float32(math.Sqrt(float64(limit2 / len2))))
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
		return self.SclScalar(float32(math.Sqrt(float64(max2 / len2))))
	}
	min2 := min * min
	if len2 < min2 {
		return self.SclScalar(float32(math.Sqrt(float64(min2 / len2))))
	}
	return self
}

func (self *Vector2) SetLength(length float32) *Vector2 {
	return self.SetLength2(length * length)
}

func (self *Vector2) SetLength2(len2 float32) *Vector2 {
	oldLen2 := self.Len2()
	if oldLen2 == 0 || oldLen2 == len2 {
		return self
	}
	return self.SclScalar(float32(math.Sqrt(float64(len2 / oldLen2))))
}

// Left-multiplies this vector by the given matrix
func (self *Vector2) Mul(mat Matrix3) *Vector2 {
	self.X = self.X*mat.val[0] + self.Y*mat.val[3] + mat.val[6]
	self.Y = self.X*mat.val[1] + self.Y*mat.val[4] + mat.val[7]
	return self
}

// Calculates the 2D cross product between this and the given vector
func (self *Vector2) CrsV(v *Vector2) float32 {
	return self.X*v.Y - self.Y*v.X
}

// Calculates the 2D cross product between this and the given vector.
func (self *Vector2) Crs(x, y float32) float32 {
	return self.X*y - self.Y*x
}

// @return the angle in degrees of this vector (point) relative to the x-axis.
// Angles are towards the positive y-axis (typically counter-clockwise) and between 0 and 360.
func (self *Vector2) Angle() float32 {
	angle := float32(math.Atan2(float64(self.Y), float64(self.X)) * float64(utils.RadiansToDegrees))
	if angle < 0 {
		angle += 360
	}
	return angle
}

// @return the angle in degrees of this vector (point) relative to the given vector.
// Angles are towards the positive y-axis (typically counter-clockwise.) between -180 and +180
func (self *Vector2) AngleVector2(reference *Vector2) float32 {
	return float32(math.Atan2(float64(self.CrsV(reference)), float64(self.DotV(reference))*float64(utils.RadiansToDegrees)))
}

// @return the angle in radians of this vector (point) relative to the x-axis.
// Angles are towards the positive y-axis. (typically counter-clockwise)
func (self *Vector2) AngleRad() float32 {
	return float32(math.Atan2(float64(self.Y), float64(self.X)))
}

// @return the angle in radians of this vector (point) relative to the given vector.
// Angles are towards the positive y-axis. (typically counter-clockwise.)
func (self *Vector2) AngleRadVector2(reference *Vector2) float32 {
	return float32(math.Atan2(float64(self.CrsV(reference)), float64(self.DotV(reference))))
}

// Sets the angle of the vector in degrees relative to the x-axis, towards the positive y-axis (typically counter-clockwise).
// @param degrees The angle in degrees to set.
func (self *Vector2) SetAngle(degrees float32) *Vector2 {
	return self.SetAngleRad(degrees * utils.DegreesToRadians)
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
	return self.RotateRad(degrees * utils.DegreesToRadians)
}

// Rotates the *Vector2 by the given angle, counter-clockwise assuming the y-axis points up.
// @param radians the angle in radians
func (self *Vector2) RotateRad(radians float32) *Vector2 {
	cos := float32(math.Cos(float64(radians)))
	sin := float32(math.Sin(float64(radians)))
	self.X = self.X*cos - self.Y*sin
	self.Y = self.X*sin + self.Y*cos
	return self
}

// Rotates the *Vector2 by 90 degrees in the specified direction,
// where >= 0 is counter-clockwise and < 0 is clockwise.
func (self *Vector2) Rotate90(dir int) *Vector2 {
	x := self.X
	if dir >= 0 {
		self.X = -self.Y
		self.Y = self.X
	} else {
		self.X = self.Y
		self.Y = -x
	}
	return self
}

func (self *Vector2) Lerp(target *Vector2, alpha float32) *Vector2 {
	invAlpha := 1.0 - alpha
	self.X = (self.X * invAlpha) + (target.X * alpha)
	self.Y = (self.Y * invAlpha) + (target.Y * alpha)
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
	// if (NumberUtils.floatToIntBits(x) != NumberUtils.floatToIntBits(other.X)) return false;
	// if (NumberUtils.floatToIntBits(y) != NumberUtils.floatToIntBits(other.Y)) return false;
	return true
}

func (self *Vector2) EpsilonEqualsV(other *Vector2, epsilon float32) bool {
	if other == nil {
		return false
	}
	if math.Abs(float64(other.X-self.X)) > float64(epsilon) {
		return false
	}
	if math.Abs(float64(other.Y-self.Y)) > float64(epsilon) {
		return false
	}
	return true
}

// Compares this vector with the other vector, using the supplied epsilon for fuzzy equality testing.
// @return whether the vectors are the same
func (self *Vector2) EpsilonEquals(x, y, epsilon float32) bool {
	if math.Abs(float64(x-self.X)) > float64(epsilon) {
		return false
	}
	if math.Abs(float64(y-self.Y)) > float64(epsilon) {
		return false
	}
	return true
}

func (self *Vector2) IsUnit() bool {
	return self.IsUnitMargin(0.000000001)
}

func (self *Vector2) IsUnitMargin(margin float32) bool {
	return math.Abs(float64(self.Len2()-1)) < float64(margin)
}

func (self *Vector2) IsZero() bool {
	return self.X == 0 && self.Y == 0
}

func (self *Vector2) IsZeroMargin(margin float32) bool {
	return self.Len2() < margin
}

func (self *Vector2) IsOnLine(other *Vector2) bool {
	return utils.IsZero(self.X*other.Y - self.Y*other.X)
}

func (self *Vector2) IsOnLineEpsilon(other *Vector2, epsilon float32) bool {
	return utils.IsZeroTolerance(self.X*other.Y-self.Y*other.X, epsilon)
}

func (self *Vector2) IsCollinearEpsilon(other *Vector2, epsilon float32) bool {
	return self.IsOnLineEpsilon(other, epsilon) && self.DotV(other) > 0
}

func (self *Vector2) IsCollinear(other *Vector2) bool {
	return self.IsOnLine(other) && self.DotV(other) > 0
}

func (self *Vector2) IsCollinearOppositeEpsilon(other *Vector2, epsilon float32) bool {
	return self.IsOnLineEpsilon(other, epsilon) && self.DotV(other) < 0
}

func (self *Vector2) IsCollinearOpposite(other *Vector2) bool {
	return self.IsOnLine(other) && self.DotV(other) < 0
}

func (self *Vector2) IsPerpendicular(v *Vector2) bool {
	return utils.IsZero(self.DotV(v))
}

func (self *Vector2) IsPerpendicularEpsilon(vector *Vector2, epsilon float32) bool {
	return utils.IsZeroTolerance(self.DotV(vector), epsilon)
}

func (self *Vector2) HasSameDirection(v *Vector2) bool {
	return self.DotV(v) > 0
}

func (self *Vector2) HasOppositeDirection(v *Vector2) bool {
	return self.DotV(v) < 0
}

func (self *Vector2) SetZero() *Vector2 {
	self.X = 0
	self.Y = 0
	return self
}

func (self *Vector2) String() string {
	return ""
	// return "[" + x + ":" + y + "]"
}
