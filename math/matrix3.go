// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

import (
	"math"
)

const (
	M00 = 0
	M01 = 3
	M02 = 6
	M10 = 1
	M11 = 4
	M12 = 7
	M20 = 2
	M21 = 5
	M22 = 8
)

// A 3x3 http://en.wikipedia.org/wiki/Row-major_order#Column-major_order column major matrix; useful for 2D transforms.
type Matrix3 struct {
	val [9]float32
	tmp [9]float32
}

func NewMatrixEmpty() *Matrix3 {
	matrix := &Matrix3{}
	return matrix.Idt()
}

// Constructs a matrix from the given float array. The array must have at least 9 elements; the first 9 will be copied.
// values The float array to copy. Remember that this matrix is in
// href="http://en.wikipedia.org/wiki/Row-major_order#Column-major_order olumn major order. (The float array is
// not modified.)
func NewMatrix(values [9]float32) *Matrix3 {
	matrix := &Matrix3{}
	return matrix.Set(values)
}

func NewMatrixCopy(other *Matrix3) *Matrix3 {
	matrix := &Matrix3{}
	return matrix.Set(other)
}

// Sets this matrix to the identity matrix
// return This matrix for the purpose of chaining operations.
func (self *Matrix3) Idt() *Matrix3 {
	val := self.val
	val[M00] = 1
	val[M10] = 0
	val[M20] = 0
	val[M01] = 0
	val[M11] = 1
	val[M21] = 0
	val[M02] = 0
	val[M12] = 0
	val[M22] = 1
	return self
}

// Postmultiplies this matrix with the provided matrix and stores the result in this matrix. For example:
// A.mul(B) results in A := AB
// m Matrix to multiply by.
// return This matrix for the purpose of chaining operations together.
func (self *Matrix3) Mul(m *Matrix3) *Matrix3 {
	val := self.val

	v00 := val[M00]*m.val[M00] + val[M01]*m.val[M10] + val[M02]*m.val[M20]
	v01 := val[M00]*m.val[M01] + val[M01]*m.val[M11] + val[M02]*m.val[M21]
	v02 := val[M00]*m.val[M02] + val[M01]*m.val[M12] + val[M02]*m.val[M22]

	v10 := val[M10]*m.val[M00] + val[M11]*m.val[M10] + val[M12]*m.val[M20]
	v11 := val[M10]*m.val[M01] + val[M11]*m.val[M11] + val[M12]*m.val[M21]
	v12 := val[M10]*m.val[M02] + val[M11]*m.val[M12] + val[M12]*m.val[M22]

	v20 := val[M20]*m.val[M00] + val[M21]*m.val[M10] + val[M22]*m.val[M20]
	v21 := val[M20]*m.val[M01] + val[M21]*m.val[M11] + val[M22]*m.val[M21]
	v22 := val[M20]*m.val[M02] + val[M21]*m.val[M12] + val[M22]*m.val[M22]

	val[M00] = v00
	val[M10] = v10
	val[M20] = v20
	val[M01] = v01
	val[M11] = v11
	val[M21] = v21
	val[M02] = v02
	val[M12] = v12
	val[M22] = v22

	return self
}

// Premultiplies this matrix with the provided matrix and stores the result in this matrix. For example:
// A.mulLeft(B) results in A := BA
// m The other Matrix to multiply by
// return This matrix for the purpose of chaining operations.
func (self *Matrix3) MulLeft(m *Matrix3) *Matrix3 {
	val := self.val

	v00 := m.val[M00]*val[M00] + m.val[M01]*val[M10] + m.val[M02]*val[M20]
	v01 := m.val[M00]*val[M01] + m.val[M01]*val[M11] + m.val[M02]*val[M21]
	v02 := m.val[M00]*val[M02] + m.val[M01]*val[M12] + m.val[M02]*val[M22]

	v10 := m.val[M10]*val[M00] + m.val[M11]*val[M10] + m.val[M12]*val[M20]
	v11 := m.val[M10]*val[M01] + m.val[M11]*val[M11] + m.val[M12]*val[M21]
	v12 := m.val[M10]*val[M02] + m.val[M11]*val[M12] + m.val[M12]*val[M22]

	v20 := m.val[M20]*val[M00] + m.val[M21]*val[M10] + m.val[M22]*val[M20]
	v21 := m.val[M20]*val[M01] + m.val[M21]*val[M11] + m.val[M22]*val[M21]
	v22 := m.val[M20]*val[M02] + m.val[M21]*val[M12] + m.val[M22]*val[M22]

	val[M00] = v00
	val[M10] = v10
	val[M20] = v20
	val[M01] = v01
	val[M11] = v11
	val[M21] = v21
	val[M02] = v02
	val[M12] = v12
	val[M22] = v22

	return self
}

// Sets this matrix to a rotation matrix that will rotate any vector in counter-clockwise direction around the z-axis.
// degrees the angle in degrees.
// return This matrix for the purpose of chaining operations.
func (self *Matrix3) SetToRotation(degrees float32) *Matrix3 {
	return self.SetToRotationRad(DegreesToRadians * degrees)
}

// Sets this matrix to a rotation matrix that will rotate any vector in counter-clockwise direction around the z-axis.
// radians the angle in radians.
// return This matrix for the purpose of chaining operations.
func (self *Matrix3) SetToRotationRad(radians float32) *Matrix3 {
	cos := float32(math.Cos(float64(radians)))
	sin := float32(math.Sin(float64(radians)))
	val := self.val

	val[M00] = cos
	val[M10] = sin
	val[M20] = 0

	val[M01] = -sin
	val[M11] = cos
	val[M21] = 0

	val[M02] = 0
	val[M12] = 0
	val[M22] = 1

	return self
}

func (self *Matrix3) SetToRotationAxisDeg(axis *Vector3, degrees float32) *Matrix3 {
	return self.SetToRotation(axis, CosDeg(degrees), SinDeg(degrees))
}

func (self *Matrix3) SetToRotationAxis(axis *Vector3, cos, sin float32) *Matrix3 {
	val := self.val
	oc := 1 - cos
	val[M00] = oc*axis.x*axis.x + cos
	val[M10] = oc*axis.x*axis.y - axis.z*sin
	val[M20] = oc*axis.z*axis.x + axis.y*sin
	val[M01] = oc*axis.x*axis.y + axis.z*sin
	val[M11] = oc*axis.y*axis.y + cos
	val[M21] = oc*axis.y*axis.z - axis.x*sin
	val[M02] = oc*axis.z*axis.x - axis.y*sin
	val[M12] = oc*axis.y*axis.z + axis.x*sin
	val[M22] = oc*axis.z*axis.z + cos
	return self
}

// Sets this matrix to a translation matrix.
// x the translation in x
// y the translation in y
// return This matrix for the purpose of chaining operations.
func (self *Matrix3) SetToTranslation(x, y float32) *Matrix3 {
	val := self.val

	val[M00] = 1
	val[M10] = 0
	val[M20] = 0

	val[M01] = 0
	val[M11] = 1
	val[M21] = 0

	val[M02] = x
	val[M12] = y
	val[M22] = 1

	return self
}

// Sets this matrix to a translation matrix.
// translation The translation vector.
// return This matrix for the purpose of chaining operations.
func (self *Matrix3) SetToTranslationV(translation *Vector2) *Matrix3 {
	val := self.val

	val[M00] = 1
	val[M10] = 0
	val[M20] = 0

	val[M01] = 0
	val[M11] = 1
	val[M21] = 0

	val[M02] = translation.x
	val[M12] = translation.y
	val[M22] = 1

	return self
}

// Sets this matrix to a scaling matrix.
// scaleX the scale in x
// scaleY the scale in y
// return This matrix for the purpose of chaining operations.
func (self *Matrix3) SetToScaling(scaleX, scaleY float32) *Matrix3 {
	val := self.val
	val[M00] = scaleX
	val[M10] = 0
	val[M20] = 0
	val[M01] = 0
	val[M11] = scaleY
	val[M21] = 0
	val[M02] = 0
	val[M12] = 0
	val[M22] = 1
	return self
}

// Sets this matrix to a scaling matrix.
// scale The scale vector.
// return This matrix for the purpose of chaining operations.
func (self *Matrix3) SetToScalingV(scale *Vector2) *Matrix3 {
	val := self.val
	val[M00] = scale.x
	val[M10] = 0
	val[M20] = 0
	val[M01] = 0
	val[M11] = scale.y
	val[M21] = 0
	val[M02] = 0
	val[M12] = 0
	val[M22] = 1
	return self
}

func (self *Matrix3) String() string {
	val := self.val
	return "[" + val[M00] + "|" + val[M01] + "|" + val[M02] + "]\n" //
	+"[" + val[M10] + "|" + val[M11] + "|" + val[M12] + "]\n"       //
	+"[" + val[M20] + "|" + val[M21] + "|" + val[M22] + "]"
}

// return The determinant of this matrix
func (self *Matrix3) Det() float32 {
	val := self.val
	return val[M00]*val[M11]*val[M22] + val[M01]*val[M12]*val[M20] + val[M02]*val[M10]*val[M21] - val[M00]*val[M12]*val[M21] - val[M01]*val[M10]*val[M22] - val[M02]*val[M11]*val[M20]
}

// Inverts this matrix given that the determinant is != 0.
// return This matrix for the purpose of chaining operations.
// panics if the matrix is singular (not invertible)
func (self *Matrix3) Inv() *Matrix3 {
	det := self.Det()
	if det == 0 {
		panic("Can't invert a singular matrix")
	}

	inv_det := 1 / det
	tmp := self.tmp
	val := self.val

	tmp[M00] = val[M11]*val[M22] - val[M21]*val[M12]
	tmp[M10] = val[M20]*val[M12] - val[M10]*val[M22]
	tmp[M20] = val[M10]*val[M21] - val[M20]*val[M11]
	tmp[M01] = val[M21]*val[M02] - val[M01]*val[M22]
	tmp[M11] = val[M00]*val[M22] - val[M20]*val[M02]
	tmp[M21] = val[M20]*val[M01] - val[M00]*val[M21]
	tmp[M02] = val[M01]*val[M12] - val[M11]*val[M02]
	tmp[M12] = val[M10]*val[M02] - val[M00]*val[M12]
	tmp[M22] = val[M00]*val[M11] - val[M10]*val[M01]

	val[M00] = inv_det * tmp[M00]
	val[M10] = inv_det * tmp[M10]
	val[M20] = inv_det * tmp[M20]
	val[M01] = inv_det * tmp[M01]
	val[M11] = inv_det * tmp[M11]
	val[M21] = inv_det * tmp[M21]
	val[M02] = inv_det * tmp[M02]
	val[M12] = inv_det * tmp[M12]
	val[M22] = inv_det * tmp[M22]

	return self
}

// Copies the values from the provided matrix to this matrix.
// mat The matrix to copy.
// return This matrix for the purposes of chaining.
func (self *Matrix3) SetM(mat *Matrix3) *Matrix3 {
	copy(mat.val, 0, self.val, 0, len(self.val))
	return self
}

// Copies the values from the provided affine matrix to this matrix. The last row is set to (0, 0, 1).
// affine The affine matrix to copy.
// return This matrix for the purposes of chaining.
func (self *Matrix3) SetA(affine *Affine2) *Matrix3 {
	val := self.val

	val[M00] = affine.m00
	val[M10] = affine.m10
	val[M20] = 0
	val[M01] = affine.m01
	val[M11] = affine.m11
	val[M21] = 0
	val[M02] = affine.m02
	val[M12] = affine.m12
	val[M22] = 1

	return self
}

// Sets this 3x3 matrix to the top left 3x3 corner of the provided 4x4 matrix.
// mat The matrix whose top left corner will be copied. This matrix will not be modified.
// return This matrix for the purpose of chaining operations.
func (self *Matrix3) SetM4(mat *Matrix4) *Matrix3 {
	val := self.val
	val[M00] = mat.val[Matrix4.M00]
	val[M10] = mat.val[Matrix4.M10]
	val[M20] = mat.val[Matrix4.M20]
	val[M01] = mat.val[Matrix4.M01]
	val[M11] = mat.val[Matrix4.M11]
	val[M21] = mat.val[Matrix4.M21]
	val[M02] = mat.val[Matrix4.M02]
	val[M12] = mat.val[Matrix4.M12]
	val[M22] = mat.val[Matrix4.M22]
	return self
}

// Sets the matrix to the given matrix as a float array. The float array must have at least 9 elements; the first 9 will be
// copied.
// values The matrix, in float form, that is to be copied. Remember that this matrix is in
// http://en.wikipedia.org/wiki/Row-major_order#Column-major_order column major order.
// return This matrix for the purpose of chaining methods together.
func (self *Matrix3) Set(values []float32) *Matrix3 {
	copy(values, 0, self.val, 0, len(self.val))
	return self
}

// Adds a translational component to the matrix in the 3rd column. The other columns are untouched.
// vector The translation vector.
// return This matrix for the purpose of chaining.
func (self *Matrix3) TrnV(vector *Vector2) *Matrix3 {
	self.val[M02] += vector.x
	self.val[M12] += vector.y
	return self
}

// Adds a translational component to the matrix in the 3rd column. The other columns are untouched.
// x The x-component of the translation vector.
// y The y-component of the translation vector.
// return This matrix for the purpose of chaining.
func (self *Matrix3) Trn(x, y float32) *Matrix3 {
	self.val[M02] += x
	self.val[M12] += y
	return self
}

// Adds a translational component to the matrix in the 3rd column. The other columns are untouched.
// vector The translation vector. (The z-component of the vector is ignored because this is a 3x3 matrix)
// return This matrix for the purpose of chaining.
func (self *Matrix3) TrnV3(vector *Vector3) *Matrix3 {
	self.val[M02] += vector.x
	self.val[M12] += vector.y
	return self
}

// Postmultiplies this matrix by a translation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
// x The x-component of the translation vector.
// y The y-component of the translation vector.
// return This matrix for the purpose of chaining.
func (self *Matrix3) Translate(x, y float32) *Matrix3 {
	val := self.val
	tmp[M00] = 1
	tmp[M10] = 0
	tmp[M20] = 0

	tmp[M01] = 0
	tmp[M11] = 1
	tmp[M21] = 0

	tmp[M02] = x
	tmp[M12] = y
	tmp[M22] = 1
	self.Mul(self.val, self.tmp)
	return self
}

// Postmultiplies this matrix by a translation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
// translation The translation vector.
// return This matrix for the purpose of chaining.
func (self *Matrix3) TranslateV(translation *Vector2) *Matrix3 {
	val := self.val
	tmp[M00] = 1
	tmp[M10] = 0
	tmp[M20] = 0

	tmp[M01] = 0
	tmp[M11] = 1
	tmp[M21] = 0

	tmp[M02] = translation.x
	tmp[M12] = translation.y
	tmp[M22] = 1
	self.Mul(self.val, self.tmp)
	return self
}

// Postmultiplies this matrix with a (counter-clockwise) rotation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
// degrees The angle in degrees
// return This matrix for the purpose of chaining.
func (self *Matrix3) Rotate(degrees float32) *Matrix3 {
	return self.RotateRad(DegreesToRadians * degrees)
}

// Postmultiplies this matrix with a (counter-clockwise) rotation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
// radians The angle in radians
// return This matrix for the purpose of chaining.
func (self *Matrix3) RotateRad(radians float32) *Matrix3 {
	if radians == 0 {
		return self
	}
	cos = float32(math.Cos(float64(radians)))
	sin = float32(math.Sin(float64(radians)))
	tmp := self.tmp

	tmp[M00] = cos
	tmp[M10] = sin
	tmp[M20] = 0

	tmp[M01] = -sin
	tmp[M11] = cos
	tmp[M21] = 0

	tmp[M02] = 0
	tmp[M12] = 0
	tmp[M22] = 1
	self.Mul(self.val, self.tmp)
	return self
}

// Postmultiplies this matrix with a scale matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
// scaleX The scale in the x-axis.
// scaleY The scale in the y-axis.
// return This matrix for the purpose of chaining.
func (self *Matrix3) Scale(scaleX, scaleY float32) *Matrix3 {
	tmp := self.tmp
	tmp[M00] = scaleX
	tmp[M10] = 0
	tmp[M20] = 0
	tmp[M01] = 0
	tmp[M11] = scaleY
	tmp[M21] = 0
	tmp[M02] = 0
	tmp[M12] = 0
	tmp[M22] = 1
	self.Mul(self.val, self.tmp)
	return self
}

// Postmultiplies this matrix with a scale matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
// scale The vector to scale the matrix by.
// return This matrix for the purpose of chaining.
func (self *Matrix3) ScaleV(scale *Vector2) *Matrix3 {
	tmp := self.tmp
	tmp[M00] = scale.x
	tmp[M10] = 0
	tmp[M20] = 0
	tmp[M01] = 0
	tmp[M11] = scale.y
	tmp[M21] = 0
	tmp[M02] = 0
	tmp[M12] = 0
	tmp[M22] = 1
	self.Mul(self.val, self.tmp)
	return self
}

// Get the values in this matrix.
// return The float values that make up this matrix in column-major order.
func (self *Matrix3) GetValues() [9]float32 {
	return self.val
}

func (self *Matrix3) GetTranslation(position *Vector2) *Vector2 {
	position.x = self.val[M02]
	position.y = self.val[M12]
	return position
}

func (self *Matrix3) GetScale(scale *Vector2) *Vector2 {
	val := self.val
	scale.x = float32(Math.sqrt(float64(val[M00]*val[M00] + val[M01]*val[M01])))
	scale.y = float32(Math.sqrt(float64(val[M10]*val[M10] + val[M11]*val[M11])))
	return scale
}

func (self *Matrix3) GetRotation() float32 {
	return RadiansToDegrees * float32(math.Atan2(float64(self.val[M10]), float64(self.val[M00])))
}

func (self *Matrix3) GetRotationRad() float32 {
	return float32(math.Atan2(float64(self.val[M10]), float64(self.val[M00])))
}

// Scale the matrix in the both the x and y components by the scalar value.
// scale The single value that will be used to scale both the x and y components.
// return This matrix for the purpose of chaining methods together.
func (self *Matrix3) Scl(scale float32) *Matrix3 {
	self.val[M00] *= scale
	self.val[M11] *= scale
	return self
}

// Scale this matrix using the x and y components of the vector but leave the rest of the matrix alone.
// scale The {@link Vector3} to use to scale this matrix.
// return This matrix for the purpose of chaining methods together.
func (self *Matrix3) SclV(scale *Vector2) *Matrix3 {
	self.val[M00] *= scale.x
	self.val[M11] *= scale.y
	return self
}

// Scale this matrix using the x and y components of the vector but leave the rest of the matrix alone.
// scale The {@link Vector3} to use to scale this matrix. The z component will be ignored.
// return This matrix for the purpose of chaining methods together.
func (self *Matrix3) SclV3(scale *Vector3) *Matrix3 {
	self.val[M00] *= scale.x
	self.val[M11] *= scale.y
	return self
}

// Transposes the current matrix.
// return This matrix for the purpose of chaining methods together.
func (self *Matrix3) Transpose() *Matrix3 {
	// Where MXY you do not have to change MXX
	val := self.val
	v01 := val[M10]
	v02 := val[M20]
	v10 := val[M01]
	v12 := val[M21]
	v20 := val[M02]
	v21 := val[M12]
	val[M01] = v01
	val[M02] = v02
	val[M10] = v10
	val[M12] = v12
	val[M20] = v20
	val[M21] = v21
	return self
}

// Multiplies matrix a with matrix b in the following manner:
// mul(A, B) => A := AB
// mata The float array representing the first matrix. Must have at least 9 elements.
// matb The float array representing the second matrix. Must have at least 9 elements.
func Mul(mata []float32, matb []float32) {
	v00 := mata[M00]*matb[M00] + mata[M01]*matb[M10] + mata[M02]*matb[M20]
	v01 := mata[M00]*matb[M01] + mata[M01]*matb[M11] + mata[M02]*matb[M21]
	v02 := mata[M00]*matb[M02] + mata[M01]*matb[M12] + mata[M02]*matb[M22]

	v10 := mata[M10]*matb[M00] + mata[M11]*matb[M10] + mata[M12]*matb[M20]
	v11 := mata[M10]*matb[M01] + mata[M11]*matb[M11] + mata[M12]*matb[M21]
	v12 := mata[M10]*matb[M02] + mata[M11]*matb[M12] + mata[M12]*matb[M22]

	v20 := mata[M20]*matb[M00] + mata[M21]*matb[M10] + mata[M22]*matb[M20]
	v21 := mata[M20]*matb[M01] + mata[M21]*matb[M11] + mata[M22]*matb[M21]
	v22 := mata[M20]*matb[M02] + mata[M21]*matb[M12] + mata[M22]*matb[M22]

	mata[M00] = v00
	mata[M10] = v10
	mata[M20] = v20
	mata[M01] = v01
	mata[M11] = v11
	mata[M21] = v21
	mata[M02] = v02
	mata[M12] = v12
	mata[M22] = v22
}
