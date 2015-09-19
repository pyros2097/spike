// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package vector

import (
	"math"

	. "github.com/pyros2097/gdx/math/interpolation"
	"github.com/pyros2097/gdx/math/utils"
)

var (
	XV3      = NewVector3(1, 0, 0)
	YV3      = NewVector3(0, 1, 0)
	ZV3      = NewVector3(0, 0, 1)
	ZeroV3   = NewVector3(0, 0, 0)
	tmpMatV3 = NewMatrix4Empty()
)

// Encapsulates a 3D vector. Allows chaining operations by returning a reference to itself in all modification methods.
// Implements IVector
type Vector3 struct {
	x, y, z float32
}

// Constructs a vector at (0,0,0)
func NewVector3Empty() *Vector3 {
	return &Vector3{}
}

// Creates a vector with the given components
func NewVector3(x, y, z float32) *Vector3 {
	return &Vector3{x, y, z}
}

// Creates a vector from the given vector
// vector The vector
func NewVector3Copy(v *Vector3) *Vector3 {
	return &Vector3{v.x, v.y, v.z}
}

// Creates a vector from the given array. The array must have at least 3 elements.
// values The array
func NewVector3Values(values []float32) *Vector3 {
	return &Vector3{values[0], values[1], values[2]}
}

// Creates a vector from the given vector and z-component
// vector The vector
// z The z-component
func NewVector3VZ(vector *Vector2, z float32) *Vector3 {
	return &Vector3{vector.x, vector.y, z}
}

// Sets the vector to the given components
// x The x-component
// y The y-component
// z The z-component
// return this vector for chaining
func (self *Vector3) Set(x, y, z float32) *Vector3 {
	self.x = x
	self.y = y
	self.z = z
	return self
}

func (self *Vector3) SetV(vector *Vector3) *Vector3 {
	return self.Set(vector.x, vector.y, vector.z)
}

// Sets the components from the array. The array must have at least 3 elements
// values The array
// return this vector for chaining
func (self *Vector3) SetValues(values []float32) *Vector3 {
	return self.Set(values[0], values[1], values[2])
}

// Sets the components of the given vector and z-component
// vector The vector
// z The z-component
func (self *Vector3) SetVZ(vector *Vector2, z float32) *Vector3 {
	return self.Set(vector.x, vector.y, z)
}

func (self *Vector3) Copy() *Vector3 {
	return NewVector3Copy(self)
}

func (self *Vector3) AddV(vector *Vector3) *Vector3 {
	return self.Add(vector.x, vector.y, vector.z)
}

// Adds the given vector to this component
// x The x-component of the other vector
// y The y-component of the other vector
// z The z-component of the other vector.
func (self *Vector3) Add(x, y, z float32) *Vector3 {
	return self.Set(self.x+x, self.y+y, self.z+z)
}

// Adds the given value to all three components of the vector.
// value The value
func (self *Vector3) AddValue(value float32) *Vector3 {
	return self.Set(self.x+value, self.y+value, self.z+value)
}

func (self *Vector3) SubV(a_vec *Vector3) *Vector3 {
	return self.Sub(a_vec.x, a_vec.y, a_vec.z)
}

// Subtracts the other vector from this vector.
func (self *Vector3) Sub(x, y, z float32) *Vector3 {
	return self.Set(self.x-x, self.y-y, self.z-z)
}

// Subtracts the given value from all components of this vector
func (self *Vector3) SubValue(value float32) *Vector3 {
	return self.Set(self.x-value, self.y-value, self.z-value)
}

func (self *Vector3) SclScalar(scalar float32) *Vector3 {
	return self.Set(self.x*scalar, self.y*scalar, self.z*scalar)
}

func (self *Vector3) SclV(other *Vector3) *Vector3 {
	return self.Set(self.x*other.x, self.y*other.y, self.z*other.z)
}

// Scales this vector by the given values
func (self *Vector3) Scl(vx, vy, vz float32) *Vector3 {
	return self.Set(self.x*vx, self.y*vy, self.z*vz)
}

func (self *Vector3) MulAdd(vec *Vector3, scalar float32) *Vector3 {
	self.x += vec.x * scalar
	self.y += vec.y * scalar
	self.z += vec.z * scalar
	return self
}

func (self *Vector3) MulAddV(vec, mulVec *Vector3) *Vector3 {
	self.x += vec.x * mulVec.x
	self.y += vec.y * mulVec.y
	self.z += vec.z * mulVec.z
	return self
}

// return The euclidean length
func LenV3(x, y, z float32) float32 {
	return float32(math.Sqrt(float64(x*x + y*y + z*z)))
}

func (self *Vector3) Len() float32 {
	return float32(math.Sqrt(float64(self.x*self.x + self.y*self.y + self.z*self.z)))
}

// return The squared euclidean length
func Len2V3(x, y, z float32) float32 {
	return x*x + y*y + z*z
}

func (self *Vector3) Len2() float32 {
	return self.x*self.x + self.y*self.y + self.z*self.z
}

// @param vector The other vector
// return Whether this and the other vector are equal
func (self *Vector3) Idt(vector *Vector3) bool {
	return self.x == vector.x && self.y == vector.y && self.z == vector.z
}

// return The euclidean distance between the two specified vectors
func DstV3(x1, y1, z1, x2, y2, z2 float32) float32 {
	a := x2 - x1
	b := y2 - y1
	c := z2 - z1
	return float32(math.Sqrt(float64(a*a + b*b + c*c)))
}

func (self *Vector3) DstV(vector *Vector3) float32 {
	a := vector.x - self.x
	b := vector.y - self.y
	c := vector.z - self.z
	return float32(math.Sqrt(float64(a*a + b*b + c*c)))
}

// return the distance between this point and the given point
func (self *Vector3) Dst(x, y, z float32) float32 {
	a := x - self.x
	b := y - self.y
	c := z - self.z
	return float32(math.Sqrt(float64(a*a + b*b + c*c)))
}

// return the squared distance between the given points
func Dst2V3(x1, y1, z1, x2, y2, z2 float32) float32 {
	a := x2 - x1
	b := y2 - y1
	c := z2 - z1
	return a*a + b*b + c*c
}

func (self *Vector3) Dst2V(point *Vector3) float32 {
	a := point.x - self.x
	b := point.y - self.y
	c := point.z - self.z
	return a*a + b*b + c*c
}

// Returns the squared distance between this point and the given point
// x The x-component of the other point
// y The y-component of the other point
// z The z-component of the other point
// return The squared distance
func (self *Vector3) Dst2V3(x, y, z float32) float32 {
	a := x - self.x
	b := y - self.y
	c := z - self.z
	return a*a + b*b + c*c
}

func (self *Vector3) Nor() *Vector3 {
	len2 := self.Len2()
	if len2 == 0 || len2 == 1 {
		return self
	}
	return self.SclScalar(1 / float32(math.Sqrt(float64(len2))))
}

// return The dot product between the two vectors
func DotV3(x1, y1, z1, x2, y2, z2 float32) float32 {
	return x1*x2 + y1*y2 + z1*z2
}

func (self *Vector3) DotV(vector *Vector3) float32 {
	return self.x*vector.x + self.y*vector.y + self.z*vector.z
}

// Returns the dot product between this and the given vector.
// x The x-component of the other vector
// y The y-component of the other vector
// z The z-component of the other vector
// return The dot product
func (self *Vector3) Dot(x, y, z float32) float32 {
	return self.x*x + self.y*y + self.z*z
}

// Sets this vector to the cross product between it and the other vector.
// vector The other vector
func (self *Vector3) CrsV(vector *Vector3) *Vector3 {
	return self.Set(self.y*vector.z-self.z*vector.y, self.z*vector.x-self.x*vector.z, self.x*vector.y-self.y*vector.x)
}

// Sets this vector to the cross product between it and the other vector.
// x The x-component of the other vector
// y The y-component of the other vector
// z The z-component of the other vector
func (self *Vector3) Crs(x, y, z float32) *Vector3 {
	return self.Set(self.y*z-self.z*y, self.z*x-self.x*z, self.x*y-self.y*x)
}

// Left-multiplies the vector by the given 4x3 column major matrix. The matrix should be composed by a 3x3 matrix representing
// rotation and scale plus a 1x3 matrix representing the translation.
func (self *Vector3) Mul4x3(matrix []float32) *Vector3 {
	return self.Set(self.x*matrix[0]+self.y*matrix[3]+self.z*matrix[6]+matrix[9], self.x*matrix[1]+
		self.y*matrix[4]+self.z*matrix[7]+matrix[10], self.x*matrix[2]+self.y*matrix[5]+self.z*
		matrix[8]+matrix[11])
}

// Left-multiplies the vector by the given matrix, assuming the fourth (w) component of the vector is 1.
func (self *Vector3) Mul(matrix *Matrix4) *Vector3 {
	l_mat := matrix.val
	return self.Set(self.x*l_mat[M4_00]+self.y*l_mat[M4_01]+self.z*l_mat[M4_02]+
		l_mat[M4_03], self.x*l_mat[M4_10]+self.y*l_mat[M4_11]+self.z*l_mat[M4_12]+
		l_mat[M4_13], self.x*l_mat[M4_20]+self.y*l_mat[M4_21]+self.z*l_mat[M4_22]+
		l_mat[M4_23])
}

// Multiplies the vector by the transpose of the given matrix, assuming the fourth (w) component of the vector is 1.
func (self *Vector3) TraMul(matrix *Matrix4) *Vector3 {
	l_mat := matrix.val
	return self.Set(self.x*l_mat[M4_00]+self.y*l_mat[M4_10]+self.z*l_mat[M4_20]+
		l_mat[M4_30], self.x*l_mat[M4_01]+self.y*l_mat[M4_11]+self.z*l_mat[M4_21]+
		l_mat[M4_31], self.x*l_mat[M4_02]+self.y*l_mat[M4_12]+self.z*l_mat[M4_22]+
		l_mat[M4_32])
}

// Left-multiplies the vector by the given matrix.
func (self *Vector3) MulM3(matrix *Matrix3) *Vector3 {
	l_mat := matrix.val
	return self.Set(self.x*l_mat[M3_00]+self.y*l_mat[M3_01]+self.z*l_mat[M3_02], self.x*
		l_mat[M3_10]+self.y*l_mat[M3_11]+self.z*l_mat[M3_12], self.x*l_mat[M3_20]+
		self.y*l_mat[M3_21]+self.z*l_mat[M3_22])
}

// Multiplies the vector by the transpose of the given matrix.
func (self *Vector3) TraMulM3(matrix *Matrix3) *Vector3 {
	l_mat := matrix.val
	return self.Set(self.x*l_mat[M3_00]+self.y*l_mat[M3_10]+self.z*l_mat[M3_20], self.x*
		l_mat[M3_01]+self.y*l_mat[M3_11]+self.z*l_mat[M3_21], self.x*l_mat[M3_02]+
		self.y*l_mat[M3_12]+self.z*l_mat[M3_22])
}

// Multiplies the vector by the given {@link Quaternion}.
func (self *Vector3) MulQ(quat *Quaternion) *Vector3 {
	return quat.Transform(self)
}

// Multiplies this vector by the given matrix dividing by w, assuming the fourth (w) component of the vector is 1. This is
// mostly used to project/unproject vectors via a perspective projection matrix.
func (self *Vector3) Prj(matrix *Matrix4) *Vector3 {
	l_mat := matrix.val
	l_w := 1 / (self.x*l_mat[M4_30] + self.y*l_mat[M4_31] + self.z*l_mat[M4_32] + l_mat[M4_33])
	return self.Set((self.x*l_mat[M4_00]+self.y*l_mat[M4_01]+self.z*l_mat[M4_02]+
		l_mat[M4_03])*l_w, (self.x*l_mat[M4_10]+self.y*l_mat[M4_11]+self.z*l_mat[M4_12]+
		l_mat[M4_13])*l_w, (self.x*l_mat[M4_20]+self.y*l_mat[M4_21]+self.z*l_mat[M4_22]+
		l_mat[M4_23])*l_w)
}

// Multiplies this vector by the first three columns of the matrix, essentially only applying rotation and scaling.
func (self *Vector3) Rot(matrix *Matrix4) *Vector3 {
	l_mat := matrix.val
	return self.Set(self.x*l_mat[M4_00]+self.y*l_mat[M4_01]+self.z*l_mat[M4_02], self.x*
		l_mat[M4_10]+self.y*l_mat[M4_11]+self.z*l_mat[M4_12], self.x*l_mat[M4_20]+
		self.y*l_mat[M4_21]+self.z*l_mat[M4_22])
}

// Multiplies this vector by the transpose of the first three columns of the matrix. Note: only works for translation and
// rotation, does not work for scaling. For those, use {@link #rot(Matrix4)} with {@link Matrix4#inv()}.
func (self *Vector3) Unrotate(matrix *Matrix4) *Vector3 {
	l_mat := matrix.val
	return self.Set(self.x*l_mat[M4_00]+self.y*l_mat[M4_10]+self.z*l_mat[M4_20], self.x*
		l_mat[M4_01]+self.y*l_mat[M4_11]+self.z*l_mat[M4_21], self.x*l_mat[M4_02]+
		self.y*l_mat[M4_12]+self.z*l_mat[M4_22])
}

// Translates this vector in the direction opposite to the translation of the matrix and the multiplies this vector by the
// transpose of the first three columns of the matrix. Note: only works for translation and rotation, does not work for
// scaling. For those, use {@link #mul(Matrix4)} with {@link Matrix4#inv()}.
func (self *Vector3) Untransform(matrix *Matrix4) *Vector3 {
	l_mat := matrix.val
	self.x -= l_mat[M4_03]
	self.y -= l_mat[M4_03]
	self.z -= l_mat[M4_03]
	return self.Set(self.x*l_mat[M4_00]+self.y*l_mat[M4_10]+self.z*l_mat[M4_20], self.x*
		l_mat[M4_01]+self.y*l_mat[M4_11]+self.z*l_mat[M4_21], self.x*l_mat[M4_02]+
		self.y*l_mat[M4_12]+self.z*l_mat[M4_22])
}

// Rotates this vector by the given angle in degrees around the given axis.
// degrees the angle in degrees
// axisX the x-component of the axis
// axisY the y-component of the axis
// axisZ the z-component of the axis
func (self *Vector3) Rotate(degrees, axisX, axisY, axisZ float32) *Vector3 {
	return self.Mul(tmpMatV3.SetToRotationAxis(axisX, axisY, axisZ, degrees))
}

// Rotates this vector by the given angle in radians around the given axis.
// radians the angle in radians
// axisX the x-component of the axis
// axisY the y-component of the axis
// axisZ the z-component of the axis
func (self *Vector3) RotateRad(radians, axisX, axisY, axisZ float32) *Vector3 {
	return self.Mul(tmpMatV3.SetToRotationRad(axisX, axisY, axisZ, radians))
}

// Rotates this vector by the given angle in degrees around the given axis.
// axis the axis
// degrees the angle in degrees
func (self *Vector3) RotateV(axis *Vector3, degrees float32) *Vector3 {
	tmpMatV3.SetToRotationV3(axis, degrees)
	return self.Mul(tmpMatV3)
}

// Rotates this vector by the given angle in radians around the given axis.
// axis the axis
// radians the angle in radians
func (self *Vector3) RotateRadV(axis *Vector3, radians float32) *Vector3 {
	tmpMatV3.SetToRotationRadV3(axis, radians)
	return self.Mul(tmpMatV3)
}

func (self *Vector3) IsUnit() bool {
	return self.IsUnitMargin(0.000000001)
}

func (self *Vector3) IsUnitMargin(margin float32) bool {
	return math.Abs(float64(self.Len2()-1)) < float64(margin)
}

func (self *Vector3) IsZero() bool {
	return self.x == 0 && self.y == 0 && self.z == 0
}

func (self *Vector3) IsZeroMargin(margin float32) bool {
	return self.Len2() < margin
}

func (self *Vector3) IsOnLineEpsilon(other *Vector3, epsilon float32) bool {
	return Len2V3(self.y*other.z-self.z*other.y, self.z*other.x-self.x*other.z, self.x*other.y-self.y*other.x) <= epsilon
}

func (self *Vector3) IsOnLine(other *Vector3) bool {
	return Len2V3(self.y*other.z-self.z*other.y, self.z*other.x-self.x*other.z, self.x*other.y-self.y*other.x) <= utils.FLOAT_ROUNDING_ERROR
}

func (self *Vector3) IsCollinearEpsilon(other *Vector3, epsilon float32) bool {
	return self.IsOnLineEpsilon(other, epsilon) && self.HasSameDirection(other)
}

func (self *Vector3) IsCollinear(other *Vector3) bool {
	return self.IsOnLine(other) && self.HasSameDirection(other)
}

func (self *Vector3) IsCollinearOppositeEpsilon(other *Vector3, epsilon float32) bool {
	return self.IsOnLineEpsilon(other, epsilon) && self.HasOppositeDirection(other)
}

func (self *Vector3) IsCollinearOpposite(other *Vector3) bool {
	return self.IsOnLine(other) && self.HasOppositeDirection(other)
}

func (self *Vector3) IsPerpendicular(vector *Vector3) bool {
	return utils.IsZero(self.DotV(vector))
}

func (self *Vector3) IsPerpendicularEpsilon(vector *Vector3, epsilon float32) bool {
	return utils.IsZeroTolerance(self.DotV(vector), epsilon)
}

func (self *Vector3) HasSameDirection(vector *Vector3) bool {
	return self.DotV(vector) > 0
}

func (self *Vector3) HasOppositeDirection(vector *Vector3) bool {
	return self.DotV(vector) < 0
}

func (self *Vector3) Lerp(target *Vector3, alpha float32) *Vector3 {
	self.x += alpha * (target.x - self.x)
	self.y += alpha * (target.y - self.y)
	self.z += alpha * (target.z - self.z)
	return self
}

func (self *Vector3) Interpolate(target *Vector3, alpha float32, interpolator Interpolation) *Vector3 {
	return self.Lerp(target, InterpolationStartEnd(0, 1, alpha, interpolator))
}

// Spherically interpolates between this vector and the target vector by alpha which is in the range [0,1]. The result is
// stored in this vector.
// target The target vector
// alpha The interpolation coefficient.
func (self *Vector3) Slerp(target *Vector3, alpha float32) *Vector3 {
	dot := self.DotV(target)
	// If the inputs are too close for comfort, simply linearly interpolate.
	if dot > 0.9995 || dot < -0.9995 {
		return self.Lerp(target, alpha)
	}

	// theta0 = angle between input vectors
	theta0 := float32(math.Acos(float64(dot)))
	// theta = angle between this vector and result
	theta := theta0 * alpha

	st := float32(math.Sin(float64(theta)))
	tx := target.x - self.x*dot
	ty := target.y - self.y*dot
	tz := target.z - self.z*dot
	l2 := tx*tx + ty*ty + tz*tz
	var dl float32
	if l2 < 0.0001 {
		dl = st * 1
	} else {
		dl = st * 1 / float32(math.Sqrt(float64(l2)))
	}
	return self.SclScalar(float32(math.Cos(float64(theta)))).Add(tx*dl, ty*dl, tz*dl).Nor()
}

func (self *Vector3) Limit(limit float32) *Vector3 {
	return self.Limit2(limit * limit)
}

func (self *Vector3) Limit2(limit2 float32) *Vector3 {
	len2 := self.Len2()
	if len2 > limit2 {
		self.SclScalar(float32(math.Sqrt(float64(limit2 / len2))))
	}
	return self
}

func (self *Vector3) SetLength(length float32) *Vector3 {
	return self.SetLength2(length * length)
}

func (self *Vector3) SetLength2(len2 float32) *Vector3 {
	oldLen2 := self.Len2()
	if oldLen2 == 0 || oldLen2 == len2 {
		return self
	} else {
		return self.SclScalar(float32(math.Sqrt(float64(len2 / oldLen2))))
	}
}

func (self *Vector3) Clamp(min, max float32) *Vector3 {
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

// func (self *Vector3) hashCode () *Vector3 {
// 	final int prime = 31;
// 	int result = 1;
// 	result = prime * result + NumberUtils.floatToIntBits(x);
// 	result = prime * result + NumberUtils.floatToIntBits(y);
// 	result = prime * result + NumberUtils.floatToIntBits(z);
// 	return result;
// }

// func (self *Vector3) equals (Object obj) bool {
// 	if (this == obj) return true;
// 	if (obj == null) return false;
// 	if (getClass() != obj.getClass()) return false;
// 	other *Vector3 = (Vector3)obj;
// 	if (NumberUtils.floatToIntBits(x) != NumberUtils.floatToIntBits(other.x)) return false;
// 	if (NumberUtils.floatToIntBits(y) != NumberUtils.floatToIntBits(other.y)) return false;
// 	if (NumberUtils.floatToIntBits(z) != NumberUtils.floatToIntBits(other.z)) return false;
// 	return true;
// }

func (self *Vector3) EpsilonEqualsV(other *Vector3, epsilon float32) bool {
	if other == nil {
		return false
	}
	if math.Abs(float64(other.x-self.x)) > float64(epsilon) {
		return false
	}
	if math.Abs(float64(other.y-self.y)) > float64(epsilon) {
		return false
	}
	if math.Abs(float64(other.z-self.z)) > float64(epsilon) {
		return false
	}
	return true
}

// Compares this vector with the other vector, using the supplied epsilon for fuzzy equality testing.
// return whether the vectors are the same.
func (self *Vector3) EpsilonEquals(x, y, z, epsilon float32) bool {
	if math.Abs(float64(x-self.x)) > float64(epsilon) {
		return false
	}
	if math.Abs(float64(y-self.y)) > float64(epsilon) {
		return false
	}
	if math.Abs(float64(z-self.z)) > float64(epsilon) {
		return false
	}
	return true
}

func (self *Vector3) SetZero() *Vector3 {
	self.x = 0
	self.y = 0
	self.z = 0
	return self
}

func (self *Vector3) String() string {
	return ""
	// return "[" + self.x + ", " + self.y + ", " + self.z + "]"
}
