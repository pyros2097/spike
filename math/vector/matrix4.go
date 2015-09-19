// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package vector

import (
	"math"

	"github.com/pyros2097/gdx/math/utils"
)

const (
	// XX: Typically the unrotated X component for scaling, also the cosine of the angle when rotated on the Y and/or Z axis. On
	// Vector3 multiplication this value is multiplied with the source X component and added to the target X component.
	M4_00 = 0
	// XY: Typically the negative sine of the angle when rotated on the Z axis. On Vector3 multiplication this value is multiplied
	// with the source Y component and added to the target X component.
	M4_01 = 4
	// XZ: Typically the sine of the angle when rotated on the Y axis. On Vector3 multiplication this value is multiplied with the
	// source Z component and added to the target X component.
	M4_02 = 8
	// XW: Typically the translation of the X component. On Vector3 multiplication this value is added to the target X component.
	M4_03 = 12
	// YX: Typically the sine of the angle when rotated on the Z axis. On Vector3 multiplication this value is multiplied with the
	// source X component and added to the target Y component.
	M4_10 = 1
	// YY: Typically the unrotated Y component for scaling, also the cosine of the angle when rotated on the X and/or Z axis. On
	// Vector3 multiplication this value is multiplied with the source Y component and added to the target Y component.
	M4_11 = 5
	// YZ: Typically the negative sine of the angle when rotated on the X axis. On Vector3 multiplication this value is multiplied
	// with the source Z component and added to the target Y component.
	M4_12 = 9
	// YW: Typically the translation of the Y component. On Vector3 multiplication this value is added to the target Y component.
	M4_13 = 13
	// ZX: Typically the negative sine of the angle when rotated on the Y axis. On Vector3 multiplication this value is multiplied
	// with the source X component and added to the target Z component.
	M4_20 = 2
	// ZY: Typical the sine of the angle when rotated on the X axis. On Vector3 multiplication this value is multiplied with the
	// source Y component and added to the target Z component.
	M4_21 = 6
	// ZZ: Typically the unrotated Z component for scaling, also the cosine of the angle when rotated on the X and/or Y axis. On
	// Vector3 multiplication this value is multiplied with the source Z component and added to the target Z component.
	M4_22 = 10
	// ZW: Typically the translation of the Z component. On Vector3 multiplication this value is added to the target Z component.
	M4_23 = 14
	// WX: Typically the value zero. On Vector3 multiplication this value is ignored.
	M4_30 = 3
	// WY: Typically the value zero. On Vector3 multiplication this value is ignored.
	M4_31 = 7
	// WZ: Typically the value zero. On Vector3 multiplication this value is ignored.
	M4_32 = 11
	// WW: Typically the value one. On Vector3 multiplication this value is ignored.
	M4_33 = 15
)

var (
	tmp [16]float32

	quat  = NewQuaternionEmpty()
	quat2 = NewQuaternionEmpty()

	l_vez = NewVector3Empty()
	l_vex = NewVector3Empty()
	l_vey = NewVector3Empty()

	tmpVec = NewVector3Empty()
	tmpMat = NewMatrix4Empty()

	right      = NewVector3Empty()
	tmpForward = NewVector3Empty()
	tmpUp      = NewVector3Empty()
)

// Encapsulates a <a href="http://en.wikipedia.org/wiki/Row-major_order#Column-major_order">column major</a> 4 by 4 matrix. Like
// the {@link Vector3} class it allows the chaining of methods by returning a reference to itself. For example:
// Matrix4 mat = new Matrix4().trn(position).mul(camera.combined);
type Matrix4 struct {
	val [16]float32
}

// Constructs an identity matrix
func NewMatrix4Empty() *Matrix4 {
	matrix := &Matrix4{}
	matrix.val[M4_00] = 1
	matrix.val[M4_11] = 1
	matrix.val[M4_22] = 1
	matrix.val[M4_33] = 1
	return matrix
}

// Constructs a matrix from the given matrix.
// param matrix The matrix to copy. (This matrix is not modified)
func NewMatrix4Copy(matrix *Matrix4) *Matrix4 {
	m := &Matrix4{}
	return m.SetM4(matrix)
}

// Constructs a matrix from the given float array. The array must have at least 16 elements; the first 16 will be copied.
// param values The float array to copy. Remember that this matrix is in
// http://en.wikipedia.org/wiki/Row-major_order column major order. (The float array is not modified)
func NewMatrix4(values [16]float32) *Matrix4 {
	m := &Matrix4{}
	return m.SetValues(values)
}

// Constructs a rotation matrix from the given {@link Quaternion}.
// param quaternion The quaternion to be copied. (The quaternion is not modified)
func NewMatrix4Q(quaternion *Quaternion) *Matrix4 {
	m := &Matrix4{}
	return m.SetQ(quaternion)
}

// Construct a matrix from the given translation, rotation and scale.
// param position The translation
// param rotation The rotation, must be normalized
// param scale The scale
func NewMatrix4V(position *Vector3, rotation *Quaternion, scale *Vector3) *Matrix4 {
	m := &Matrix4{}
	return m.SetVRS(position, rotation, scale)
}

// Sets the matrix to the given matrix.
// param matrix The matrix that is to be copied. (The given matrix is not modified)
func (self *Matrix4) SetM4(matrix *Matrix4) *Matrix4 {
	return self.SetValues(matrix.val)
}

// Sets the matrix to the given matrix as a float array. The float array must have at least 16 elements; the first 16 will be
// copied.
// param values The matrix, in float form, that is to be copied. Remember that this matrix is in
// http://en.wikipedia.org/wiki/Row-major_order column major order.
func (self *Matrix4) SetValues(values [16]float32) *Matrix4 {
	copy(values, 0, val, 0, val.length)
	return self
}

// Sets the matrix to a rotation matrix representing the quaternion.
// param quaternion The quaternion that is to be used to set this matrix.
func (self *Matrix4) SetQ(quaternion *Quaternion) *Matrix4 {
	return self.SetQR(quaternion.x, quaternion.y, quaternion.z, quaternion.w)
}

// Sets the matrix to a rotation matrix representing the quaternion.
// param qX The X component of the quaternion that is to be used to set this matrix.
// param qY The Y component of the quaternion that is to be used to set this matrix.
// param qZ The Z component of the quaternion that is to be used to set this matrix.
// param qW The W component of the quaternion that is to be used to set this matrix.
func (self *Matrix4) SetQR(qX, qY, qZ, qW float32) *Matrix4 {
	return self.SetTQ(0, 0, 0, qX, qY, qZ, qW)
}

// Set this matrix to the specified translation and rotation.
// param position The translation
// param orientation The rotation, must be normalized
func (self *Matrix4) SetV(position *Vector3, orientation *Quaternion) *Matrix4 {
	return self.SetTQ(position.x, position.y, position.z, orientation.x, orientation.y, orientation.z, orientation.w)
}

// Sets the matrix to a rotation matrix representing the translation and quaternion.
// param tX The X component of the translation that is to be used to set this matrix.
// param tY The Y component of the translation that is to be used to set this matrix.
// param tZ The Z component of the translation that is to be used to set this matrix.
// param qX The X component of the quaternion that is to be used to set this matrix.
// param qY The Y component of the quaternion that is to be used to set this matrix.
// param qZ The Z component of the quaternion that is to be used to set this matrix.
// param qW The W component of the quaternion that is to be used to set this matrix.
func (self *Matrix4) SetTQ(tX, tY, tZ, qX, qY, qZ, qW float32) *Matrix4 {
	xs := qX * 2
	ys := qY * 2
	zs := qZ * 2
	wx := qW * xs
	wy := qW * ys
	wz := qW * zs
	xx := qX * xs
	xy := qX * ys
	xz := qX * zs
	yy := qY * ys
	yz := qY * zs
	zz := qZ * zs

	self.val[M4_00] = (1.0 - (yy + zz))
	self.val[M4_01] = (xy - wz)
	self.val[M4_02] = (xz + wy)
	self.val[M4_03] = tX

	self.val[M4_10] = (xy + wz)
	self.val[M4_11] = (1.0 - (xx + zz))
	self.val[M4_12] = (yz - wx)
	self.val[M4_13] = tY

	self.val[M4_20] = (xz - wy)
	self.val[M4_21] = (yz + wx)
	self.val[M4_22] = (1.0 - (xx + yy))
	self.val[M4_23] = tZ

	self.val[M4_30] = 0
	self.val[M4_31] = 0
	self.val[M4_32] = 0
	self.val[M4_33] = 1.0
	return self
}

// Set this matrix to the specified translation, rotation and scale.
// param p The translation
// param orient The rotation, must be normalized
// param scale The scale
func (self *Matrix4) SetVRS(p *Vector3, orient *Quaternion, scale *Vector3) *Matrix4 {
	return self.Set(p.x, p.y, p.z, orient.x, orient.y, orient.z, orient.w, scale.x, scale.y, scale.z)
}

// Sets the matrix to a rotation matrix representing the translation and quaternion.
// param tX The X component of the translation that is to be used to set this matrix.
// param tY The Y component of the translation that is to be used to set this matrix.
// param tZ The Z component of the translation that is to be used to set this matrix.
// param qX The X component of the quaternion that is to be used to set this matrix.
// param qY The Y component of the quaternion that is to be used to set this matrix.
// param qZ The Z component of the quaternion that is to be used to set this matrix.
// param qW The W component of the quaternion that is to be used to set this matrix.
// param scaleX The X component of the scaling that is to be used to set this matrix.
// param scaleY The Y component of the scaling that is to be used to set this matrix.
// param scaleZ The Z component of the scaling that is to be used to set this matrix.
func (self *Matrix4) Set(tX, tY, tZ, qX, qY, qZ, qW, scaleX, scaleY, scaleZ float32) *Matrix4 {
	xs := qX * 2
	ys := qY * 2
	zs := qZ * 2
	wx := qW * xs
	wy := qW * ys
	wz := qW * zs
	xx := qX * xs
	xy := qX * ys
	xz := qX * zs
	yy := qY * ys
	yz := qY * zs
	zz := qZ * zs

	self.val[M4_00] = scaleX * (1.0 - (yy + zz))
	self.val[M4_01] = scaleY * (xy - wz)
	self.val[M4_02] = scaleZ * (xz + wy)
	self.val[M4_03] = tX

	self.val[M4_10] = scaleX * (xy + wz)
	self.val[M4_11] = scaleY * (1.0 - (xx + zz))
	self.val[M4_12] = scaleZ * (yz - wx)
	self.val[M4_13] = tY

	self.val[M4_20] = scaleX * (xz - wy)
	self.val[M4_21] = scaleY * (yz + wx)
	self.val[M4_22] = scaleZ * (1.0 - (xx + yy))
	self.val[M4_23] = tZ

	self.val[M4_30] = 0
	self.val[M4_31] = 0
	self.val[M4_32] = 0
	self.val[M4_33] = 1.0
	return self
}

// Sets the four columns of the matrix which correspond to the x-, y- and z-axis of the vector space this matrix creates as
// well as the 4th column representing the translation of any point that is multiplied by this matrix.
// param xAxis The x-axis.
// param yAxis The y-axis.
// param zAxis The z-axis.
// param pos The translation vector.
func (self *Matrix4) SetVAxis(xAxis, yAxis, zAxis, pos *Vector3) *Matrix4 {
	self.val[M4_00] = xAxis.x
	self.val[M4_01] = xAxis.y
	self.val[M4_02] = xAxis.z
	self.val[M4_10] = yAxis.x
	self.val[M4_11] = yAxis.y
	self.val[M4_12] = yAxis.z
	self.val[M4_20] = zAxis.x
	self.val[M4_21] = zAxis.y
	self.val[M4_22] = zAxis.z
	self.val[M4_03] = pos.x
	self.val[M4_13] = pos.y
	self.val[M4_23] = pos.z
	self.val[M4_30] = 0
	self.val[M4_31] = 0
	self.val[M4_32] = 0
	self.val[M4_33] = 1
	return self
}

// @return a copy of this matrix
func (self *Matrix4) Copy() *Matrix4 {
	return NewMatrix4Copy(self)
}

// Adds a translational component to the matrix in the 4th column. The other columns are untouched.
// param vector The translation vector to add to the current matrix. (This vector is not modified)
func (self *Matrix4) TrnV(vector *Vector3) *Matrix4 {
	self.val[M4_03] += vector.x
	self.val[M4_13] += vector.y
	self.val[M4_23] += vector.z
	return self
}

// Adds a translational component to the matrix in the 4th column. The other columns are untouched.
// param x The x-component of the translation vector.
// param y The y-component of the translation vector.
// param z The z-component of the translation vector.
func (self *Matrix4) Trn(x, y, z float32) *Matrix4 {
	self.val[M4_03] += x
	self.val[M4_13] += y
	self.val[M4_23] += z
	return self
}

// @return the backing float array
func (self *Matrix4) GetValues() [16]float32 {
	return self.val
}

// Postmultiplies this matrix with the given matrix, storing the result in this matrix. For example:
//
// A.mul(B) results in A := AB.
//
// param matrix The other matrix to multiply by.
func (self *Matrix4) MulM4(matrix *Matrix4) *Matrix4 {
	MulM4(self.val, matrix.val)
	return self
}

// Premultiplies this matrix with the given matrix, storing the result in this matrix. For example:
//
// A.mulLeft(B) results in A := BA.
//
// param matrix The other matrix to multiply by.
func (self *Matrix4) MulLeft(matrix *Matrix4) *Matrix4 {
	tmpMat.SetM4(matrix)
	MulM4(tmpMat.val, self.val)
	return self.SetM4(tmpMat)
}

// Transposes the matrix.
func (self *Matrix4) Tra() *Matrix4 {
	tmp[M4_00] = self.val[M4_00]
	tmp[M4_01] = self.val[M4_10]
	tmp[M4_02] = self.val[M4_20]
	tmp[M4_03] = self.val[M4_30]
	tmp[M4_10] = self.val[M4_01]
	tmp[M4_11] = self.val[M4_11]
	tmp[M4_12] = self.val[M4_21]
	tmp[M4_13] = self.val[M4_31]
	tmp[M4_20] = self.val[M4_02]
	tmp[M4_21] = self.val[M4_12]
	tmp[M4_22] = self.val[M4_22]
	tmp[M4_23] = self.val[M4_32]
	tmp[M4_30] = self.val[M4_03]
	tmp[M4_31] = self.val[M4_13]
	tmp[M4_32] = self.val[M4_23]
	tmp[M4_33] = self.val[M4_33]
	return self.SetValues(tmp)
}

// Sets the matrix to an identity matrix.
func (self *Matrix4) Idt() *Matrix4 {
	self.val[M4_00] = 1
	self.val[M4_01] = 0
	self.val[M4_02] = 0
	self.val[M4_03] = 0
	self.val[M4_10] = 0
	self.val[M4_11] = 1
	self.val[M4_12] = 0
	self.val[M4_13] = 0
	self.val[M4_20] = 0
	self.val[M4_21] = 0
	self.val[M4_22] = 1
	self.val[M4_23] = 0
	self.val[M4_30] = 0
	self.val[M4_31] = 0
	self.val[M4_32] = 0
	self.val[M4_33] = 1
	return self
}

// Inverts the matrix. Stores the result in this matrix.
// @throws RuntimeException if the matrix is singular (not invertible)
func (self *Matrix4) Inv() *Matrix4 {
	val := self.val
	l_det := val[M4_30]*val[M4_21]*val[M4_12]*val[M4_03] - val[M4_20]*val[M4_31]*val[M4_12]*val[M4_03] - val[M4_30]*val[M4_11]*
		self.val[M4_22]*val[M4_03] + val[M4_10]*val[M4_31]*val[M4_22]*val[M4_03] + val[M4_20]*val[M4_11]*val[M4_32]*val[M4_03] - val[M4_10]*
		self.val[M4_21]*val[M4_32]*val[M4_03] - val[M4_30]*val[M4_21]*val[M4_02]*val[M4_13] + val[M4_20]*val[M4_31]*val[M4_02]*val[M4_13] +
		self.val[M4_30]*val[M4_01]*val[M4_22]*val[M4_13] - val[M4_00]*val[M4_31]*val[M4_22]*val[M4_13] - val[M4_20]*val[M4_01]*val[M4_32]*
		self.val[M4_13] + val[M4_00]*val[M4_21]*val[M4_32]*val[M4_13] + val[M4_30]*val[M4_11]*val[M4_02]*val[M4_23] - val[M4_10]*val[M4_31]*
		self.val[M4_02]*val[M4_23] - val[M4_30]*val[M4_01]*val[M4_12]*val[M4_23] + val[M4_00]*val[M4_31]*val[M4_12]*val[M4_23] + val[M4_10]*
		self.val[M4_01]*val[M4_32]*val[M4_23] - val[M4_00]*val[M4_11]*val[M4_32]*val[M4_23] - val[M4_20]*val[M4_11]*val[M4_02]*val[M4_33] +
		self.val[M4_10]*val[M4_21]*val[M4_02]*val[M4_33] + val[M4_20]*val[M4_01]*val[M4_12]*val[M4_33] - val[M4_00]*val[M4_21]*val[M4_12]*
		self.val[M4_33] - val[M4_10]*val[M4_01]*val[M4_22]*val[M4_33] + val[M4_00]*val[M4_11]*val[M4_22]*val[M4_33]
	if l_det == 0 {
		panic("non-invertible matrix")
	}
	inv_det := 1.0 / l_det
	tmp[M4_00] = val[M4_12]*val[M4_23]*val[M4_31] - val[M4_13]*val[M4_22]*val[M4_31] + val[M4_13]*val[M4_21]*val[M4_32] - val[M4_11]*val[M4_23]*val[M4_32] - val[M4_12]*val[M4_21]*val[M4_33] + val[M4_11]*val[M4_22]*val[M4_33]
	tmp[M4_01] = val[M4_03]*val[M4_22]*val[M4_31] - val[M4_02]*val[M4_23]*val[M4_31] - val[M4_03]*val[M4_21]*val[M4_32] + val[M4_01]*val[M4_23]*val[M4_32] + val[M4_02]*val[M4_21]*val[M4_33] - val[M4_01]*val[M4_22]*val[M4_33]
	tmp[M4_02] = val[M4_02]*val[M4_13]*val[M4_31] - val[M4_03]*val[M4_12]*val[M4_31] + val[M4_03]*val[M4_11]*val[M4_32] - val[M4_01]*val[M4_13]*val[M4_32] - val[M4_02]*val[M4_11]*val[M4_33] + val[M4_01]*val[M4_12]*val[M4_33]
	tmp[M4_03] = val[M4_03]*val[M4_12]*val[M4_21] - val[M4_02]*val[M4_13]*val[M4_21] - val[M4_03]*val[M4_11]*val[M4_22] + val[M4_01]*val[M4_13]*val[M4_22] + val[M4_02]*val[M4_11]*val[M4_23] - val[M4_01]*val[M4_12]*val[M4_23]
	tmp[M4_10] = val[M4_13]*val[M4_22]*val[M4_30] - val[M4_12]*val[M4_23]*val[M4_30] - val[M4_13]*val[M4_20]*val[M4_32] + val[M4_10]*val[M4_23]*val[M4_32] + val[M4_12]*val[M4_20]*val[M4_33] - val[M4_10]*val[M4_22]*val[M4_33]
	tmp[M4_11] = val[M4_02]*val[M4_23]*val[M4_30] - val[M4_03]*val[M4_22]*val[M4_30] + val[M4_03]*val[M4_20]*val[M4_32] - val[M4_00]*val[M4_23]*val[M4_32] - val[M4_02]*val[M4_20]*val[M4_33] + val[M4_00]*val[M4_22]*val[M4_33]
	tmp[M4_12] = val[M4_03]*val[M4_12]*val[M4_30] - val[M4_02]*val[M4_13]*val[M4_30] - val[M4_03]*val[M4_10]*val[M4_32] + val[M4_00]*val[M4_13]*val[M4_32] + val[M4_02]*val[M4_10]*val[M4_33] - val[M4_00]*val[M4_12]*val[M4_33]
	tmp[M4_13] = val[M4_02]*val[M4_13]*val[M4_20] - val[M4_03]*val[M4_12]*val[M4_20] + val[M4_03]*val[M4_10]*val[M4_22] - val[M4_00]*val[M4_13]*val[M4_22] - val[M4_02]*val[M4_10]*val[M4_23] + val[M4_00]*val[M4_12]*val[M4_23]
	tmp[M4_20] = val[M4_11]*val[M4_23]*val[M4_30] - val[M4_13]*val[M4_21]*val[M4_30] + val[M4_13]*val[M4_20]*val[M4_31] - val[M4_10]*val[M4_23]*val[M4_31] - val[M4_11]*val[M4_20]*val[M4_33] + val[M4_10]*val[M4_21]*val[M4_33]
	tmp[M4_21] = val[M4_03]*val[M4_21]*val[M4_30] - val[M4_01]*val[M4_23]*val[M4_30] - val[M4_03]*val[M4_20]*val[M4_31] + val[M4_00]*val[M4_23]*val[M4_31] + val[M4_01]*val[M4_20]*val[M4_33] - val[M4_00]*val[M4_21]*val[M4_33]
	tmp[M4_22] = val[M4_01]*val[M4_13]*val[M4_30] - val[M4_03]*val[M4_11]*val[M4_30] + val[M4_03]*val[M4_10]*val[M4_31] - val[M4_00]*val[M4_13]*val[M4_31] - val[M4_01]*val[M4_10]*val[M4_33] + val[M4_00]*val[M4_11]*val[M4_33]
	tmp[M4_23] = val[M4_03]*val[M4_11]*val[M4_20] - val[M4_01]*val[M4_13]*val[M4_20] - val[M4_03]*val[M4_10]*val[M4_21] + val[M4_00]*val[M4_13]*val[M4_21] + val[M4_01]*val[M4_10]*val[M4_23] - val[M4_00]*val[M4_11]*val[M4_23]
	tmp[M4_30] = val[M4_12]*val[M4_21]*val[M4_30] - val[M4_11]*val[M4_22]*val[M4_30] - val[M4_12]*val[M4_20]*val[M4_31] + val[M4_10]*val[M4_22]*val[M4_31] + val[M4_11]*val[M4_20]*val[M4_32] - val[M4_10]*val[M4_21]*val[M4_32]
	tmp[M4_31] = val[M4_01]*val[M4_22]*val[M4_30] - val[M4_02]*val[M4_21]*val[M4_30] + val[M4_02]*val[M4_20]*val[M4_31] - val[M4_00]*val[M4_22]*val[M4_31] - val[M4_01]*val[M4_20]*val[M4_32] + val[M4_00]*val[M4_21]*val[M4_32]
	tmp[M4_32] = val[M4_02]*val[M4_11]*val[M4_30] - val[M4_01]*val[M4_12]*val[M4_30] - val[M4_02]*val[M4_10]*val[M4_31] + val[M4_00]*val[M4_12]*val[M4_31] + val[M4_01]*val[M4_10]*val[M4_32] - val[M4_00]*val[M4_11]*val[M4_32]
	tmp[M4_33] = val[M4_01]*val[M4_12]*val[M4_20] - val[M4_02]*val[M4_11]*val[M4_20] + val[M4_02]*val[M4_10]*val[M4_21] - val[M4_00]*val[M4_12]*val[M4_21] - val[M4_01]*val[M4_10]*val[M4_22] + val[M4_00]*val[M4_11]*val[M4_22]
	self.val[M4_00] = tmp[M4_00] * inv_det
	self.val[M4_01] = tmp[M4_01] * inv_det
	self.val[M4_02] = tmp[M4_02] * inv_det
	self.val[M4_03] = tmp[M4_03] * inv_det
	self.val[M4_10] = tmp[M4_10] * inv_det
	self.val[M4_11] = tmp[M4_11] * inv_det
	self.val[M4_12] = tmp[M4_12] * inv_det
	self.val[M4_13] = tmp[M4_13] * inv_det
	self.val[M4_20] = tmp[M4_20] * inv_det
	self.val[M4_21] = tmp[M4_21] * inv_det
	self.val[M4_22] = tmp[M4_22] * inv_det
	self.val[M4_23] = tmp[M4_23] * inv_det
	self.val[M4_30] = tmp[M4_30] * inv_det
	self.val[M4_31] = tmp[M4_31] * inv_det
	self.val[M4_32] = tmp[M4_32] * inv_det
	self.val[M4_33] = tmp[M4_33] * inv_det
	return self
}

// @return The determinant of this matrix
func (self *Matrix4) Det() float32 {
	val := self.val
	return val[M4_30]*val[M4_21]*val[M4_12]*val[M4_03] - val[M4_20]*val[M4_31]*val[M4_12]*val[M4_03] - val[M4_30]*val[M4_11]*
		val[M4_22]*val[M4_03] + val[M4_10]*val[M4_31]*val[M4_22]*val[M4_03] + val[M4_20]*val[M4_11]*val[M4_32]*val[M4_03] - val[M4_10]*
		val[M4_21]*val[M4_32]*val[M4_03] - val[M4_30]*val[M4_21]*val[M4_02]*val[M4_13] + val[M4_20]*val[M4_31]*val[M4_02]*val[M4_13] +
		val[M4_30]*val[M4_01]*val[M4_22]*val[M4_13] - val[M4_00]*val[M4_31]*val[M4_22]*val[M4_13] - val[M4_20]*val[M4_01]*val[M4_32]*
		val[M4_13] + val[M4_00]*val[M4_21]*val[M4_32]*val[M4_13] + val[M4_30]*val[M4_11]*val[M4_02]*val[M4_23] - val[M4_10]*val[M4_31]*
		val[M4_02]*val[M4_23] - val[M4_30]*val[M4_01]*val[M4_12]*val[M4_23] + val[M4_00]*val[M4_31]*val[M4_12]*val[M4_23] + val[M4_10]*
		val[M4_01]*val[M4_32]*val[M4_23] - val[M4_00]*val[M4_11]*val[M4_32]*val[M4_23] - val[M4_20]*val[M4_11]*val[M4_02]*val[M4_33] +
		val[M4_10]*val[M4_21]*val[M4_02]*val[M4_33] + val[M4_20]*val[M4_01]*val[M4_12]*val[M4_33] - val[M4_00]*val[M4_21]*val[M4_12]*
		val[M4_33] - val[M4_10]*val[M4_01]*val[M4_22]*val[M4_33] + val[M4_00]*val[M4_11]*val[M4_22]*val[M4_33]
}

// @return The determinant of the 3x3 upper left matrix
func (self *Matrix4) Det3x3() float32 {
	val := self.val
	return val[M4_00]*val[M4_11]*val[M4_22] + val[M4_01]*val[M4_12]*val[M4_20] + val[M4_02]*val[M4_10]*val[M4_21] - val[M4_00]*
		val[M4_12]*val[M4_21] - val[M4_01]*val[M4_10]*val[M4_22] - val[M4_02]*val[M4_11]*val[M4_20]
}

// Sets the matrix to a projection matrix with a near- and far plane, a field of view in degrees and an aspect ratio. Note that
// the field of view specified is the angle in degrees for the height, the field of view for the width will be calculated
// according to the aspect ratio.
//
// param near The near plane
// param far The far plane
// param fovy The field of view of the height in degrees
// param aspectRatio The "width over height" aspect ratio
func (self *Matrix4) SetToProjectionNear(near, far, fovy, aspectRatio float32) *Matrix4 {
	self.Idt()
	l_fd := float32(1.0 / math.Tan(float64((fovy*(utils.PI/180))/2.0)))
	l_a1 := (far + near) / (near - far)
	l_a2 := (2 * far * near) / (near - far)
	self.val[M4_00] = l_fd / aspectRatio
	self.val[M4_10] = 0
	self.val[M4_20] = 0
	self.val[M4_30] = 0
	self.val[M4_01] = 0
	self.val[M4_11] = l_fd
	self.val[M4_21] = 0
	self.val[M4_31] = 0
	self.val[M4_02] = 0
	self.val[M4_12] = 0
	self.val[M4_22] = l_a1
	self.val[M4_32] = -1
	self.val[M4_03] = 0
	self.val[M4_13] = 0
	self.val[M4_23] = l_a2
	self.val[M4_33] = 0

	return self
}

// Sets the matrix to a projection matrix with a near/far plane, and left, bottom, right and top specifying the points on the
// near plane that are mapped to the lower left and upper right corners of the viewport. This allows to create projection
// matrix with off-center vanishing point.
//
// param left
// param right
// param bottom
// param top
// param near The near plane
// param far The far plane
func (self *Matrix4) SetToProjectionLeft(left, right, bottom, top, near, far float32) *Matrix4 {
	x := 2.0 * near / (right - left)
	y := 2.0 * near / (top - bottom)
	a := (right + left) / (right - left)
	b := (top + bottom) / (top - bottom)
	l_a1 := (far + near) / (near - far)
	l_a2 := (2 * far * near) / (near - far)
	self.val[M4_00] = x
	self.val[M4_10] = 0
	self.val[M4_20] = 0
	self.val[M4_30] = 0
	self.val[M4_01] = 0
	self.val[M4_11] = y
	self.val[M4_21] = 0
	self.val[M4_31] = 0
	self.val[M4_02] = a
	self.val[M4_12] = b
	self.val[M4_22] = l_a1
	self.val[M4_32] = -1
	self.val[M4_03] = 0
	self.val[M4_13] = 0
	self.val[M4_23] = l_a2
	self.val[M4_33] = 0

	return self
}

// Sets this matrix to an orthographic projection matrix with the origin at (x,y) extending by width and height. The near plane
// is set to 0, the far plane is set to 1.
//
// param x The x-coordinate of the origin
// param y The y-coordinate of the origin
// param width The width
// param height The height
func (self *Matrix4) SetToOrtho2D(x, y, width, height float32) *Matrix4 {
	self.SetToOrtho(x, x+width, y, y+height, 0, 1)
	return self
}

// Sets this matrix to an orthographic projection matrix with the origin at (x,y) extending by width and height, having a near
// and far plane.
//
// param x The x-coordinate of the origin
// param y The y-coordinate of the origin
// param width The width
// param height The height
// param near The near plane
// param far The far plane
func (self *Matrix4) SetToOrtho2DNear(x, y, width, height, near, far float32) *Matrix4 {
	self.SetToOrtho(x, x+width, y, y+height, near, far)
	return self
}

// Sets the matrix to an orthographic projection like glOrtho (http://www.opengl.org/sdk/docs/man/xhtml/glOrtho.xml) following
// the OpenGL equivalent
//
// param left The left clipping plane
// param right The right clipping plane
// param bottom The bottom clipping plane
// param top The top clipping plane
// param near The near clipping plane
// param far The far clipping plane
func (self *Matrix4) SetToOrtho(left, right, bottom, top, near, far float32) *Matrix4 {
	self.Idt()
	x_orth := 2 / (right - left)
	y_orth := 2 / (top - bottom)
	z_orth := -2 / (far - near)

	tx := -(right + left) / (right - left)
	ty := -(top + bottom) / (top - bottom)
	tz := -(far + near) / (far - near)

	self.val[M4_00] = x_orth
	self.val[M4_10] = 0
	self.val[M4_20] = 0
	self.val[M4_30] = 0
	self.val[M4_01] = 0
	self.val[M4_11] = y_orth
	self.val[M4_21] = 0
	self.val[M4_31] = 0
	self.val[M4_02] = 0
	self.val[M4_12] = 0
	self.val[M4_22] = z_orth
	self.val[M4_32] = 0
	self.val[M4_03] = tx
	self.val[M4_13] = ty
	self.val[M4_23] = tz
	self.val[M4_33] = 1

	return self
}

// Sets the 4th column to the translation vector.
//
// param vector The translation vector
func (self *Matrix4) SetTranslationV3(vector *Vector3) *Matrix4 {
	self.val[M4_03] = vector.x
	self.val[M4_13] = vector.y
	self.val[M4_23] = vector.z
	return self
}

// Sets the 4th column to the translation vector.
//
// param x The X coordinate of the translation vector
// param y The Y coordinate of the translation vector
// param z The Z coordinate of the translation vector
func (self *Matrix4) SetTranslation(x, y, z float32) *Matrix4 {
	self.val[M4_03] = x
	self.val[M4_13] = y
	self.val[M4_23] = z
	return self
}

// Sets this matrix to a translation matrix, overwriting it first by an identity matrix and then setting the 4th column to the
// translation vector.
//
// param vector The translation vector
func (self *Matrix4) SetToTranslationV3(vector *Vector3) *Matrix4 {
	self.Idt()
	self.val[M4_03] = vector.x
	self.val[M4_13] = vector.y
	self.val[M4_23] = vector.z
	return self
}

// Sets this matrix to a translation matrix, overwriting it first by an identity matrix and then setting the 4th column to the
// translation vector.
//
// param x The x-component of the translation vector.
// param y The y-component of the translation vector.
// param z The z-component of the translation vector.
func (self *Matrix4) SetToTranslation(x, y, z float32) *Matrix4 {
	self.Idt()
	self.val[M4_03] = x
	self.val[M4_13] = y
	self.val[M4_23] = z
	return self
}

// Sets this matrix to a translation and scaling matrix by first overwriting it with an identity and then setting the
// translation vector in the 4th column and the scaling vector in the diagonal.
//
// param translation The translation vector
// param scaling The scaling vector
func (self *Matrix4) SetToTranslationAndScalingV3(translation, scaling *Vector3) *Matrix4 {
	self.Idt()
	self.val[M4_03] = translation.x
	self.val[M4_13] = translation.y
	self.val[M4_23] = translation.z
	self.val[M4_00] = scaling.x
	self.val[M4_11] = scaling.y
	self.val[M4_22] = scaling.z
	return self
}

// Sets this matrix to a translation and scaling matrix by first overwriting it with an identity and then setting the
// translation vector in the 4th column and the scaling vector in the diagonal.
//
// param tX The x-component of the translation vector
// param tY The y-component of the translation vector
// param tZ The z-component of the translation vector
// param scalingX The x-component of the scaling vector
// param scalingY The x-component of the scaling vector
// param scalingZ The x-component of the scaling vector
func (self *Matrix4) SetToTranslationAndScaling(translationX, tY, tZ, scalingX, scalingY, scalingZ float32) *Matrix4 {
	self.Idt()
	self.val[M4_03] = translationX
	self.val[M4_13] = tY
	self.val[M4_23] = tZ
	self.val[M4_00] = scalingX
	self.val[M4_11] = scalingY
	self.val[M4_22] = scalingZ
	return self
}

// Sets the matrix to a rotation matrix around the given axis.
//
// param axis The axis
// param degrees The angle in degrees
func (self *Matrix4) SetToRotationV3(axis *Vector3, degrees float32) *Matrix4 {
	if degrees == 0 {
		self.Idt()
		return self
	}
	return self.SetQ(quat.SetV3(axis, degrees))
}

// Sets the matrix to a rotation matrix around the given axis.
//
// param axis The axis
// param radians The angle in radians
func (self *Matrix4) SetToRotationRadV3(axis *Vector3, radians float32) *Matrix4 {
	if radians == 0 {
		self.Idt()
		return self
	}
	return self.SetQ(quat.SetFromAxisRadV3(axis, radians))
}

// Sets the matrix to a rotation matrix around the given axis.
//
// param axisX The x-component of the axis
// param axisY The y-component of the axis
// param axisZ The z-component of the axis
// param degrees The angle in degrees
func (self *Matrix4) SetToRotationAxis(axisX, axisY, axisZ, degrees float32) *Matrix4 {
	if degrees == 0 {
		self.Idt()
		return self
	}
	return self.SetQ(quat.SetFromAxis(axisX, axisY, axisZ, degrees))
}

// Sets the matrix to a rotation matrix around the given axis.
//
// param axisX The x-component of the axis
// param axisY The y-component of the axis
// param axisZ The z-component of the axis
// param radians The angle in radians
func (self *Matrix4) SetToRotationRad(axisX, axisY, axisZ, radians float32) *Matrix4 {
	if radians == 0 {
		self.Idt()
		return self
	}
	return self.SetQ(quat.SetFromAxisRad(axisX, axisY, axisZ, radians))
}

// Set the matrix to a rotation matrix between two vectors.
// param v1 The base vector
// param v2 The target vector
func (self *Matrix4) SetToRotationV(v1, v2 *Vector3) *Matrix4 {
	return self.SetQ(quat.SetFromCrossV3(v1, v2))
}

// Set the matrix to a rotation matrix between two vectors.
// param x1 The base vectors x value
// param y1 The base vectors y value
// param z1 The base vectors z value
// param x2 The target vector x value
// param y2 The target vector y value
// param z2 The target vector z value
func (self *Matrix4) SetToRotation(x1, y1, z1, x2, y2, z2 float32) *Matrix4 {
	return self.SetQ(quat.SetFromCross(x1, y1, z1, x2, y2, z2))
}

// Sets this matrix to a rotation matrix from the given euler angles.
// param yaw the yaw in degrees
// param pitch the pitch in degrees
// param roll the roll in degrees
func (self *Matrix4) SetFromEulerAngles(yaw, pitch, roll float32) *Matrix4 {
	quat.SetEulerAngles(yaw, pitch, roll)
	return self.SetQ(quat)
}

// Sets this matrix to a rotation matrix from the given euler angles.
// param yaw the yaw in radians
// param pitch the pitch in radians
// param roll the roll in radians
func (self *Matrix4) SetFromEulerAnglesRad(yaw, pitch, roll float32) *Matrix4 {
	quat.SetEulerAnglesRad(yaw, pitch, roll)
	return self.SetQ(quat)
}

// Sets this matrix to a scaling matrix
// param vector The scaling vector.
func (self *Matrix4) SetToScalingV3(vector *Vector3) *Matrix4 {
	self.Idt()
	self.val[M4_00] = vector.x
	self.val[M4_11] = vector.y
	self.val[M4_22] = vector.z
	return self
}

// Sets this matrix to a scaling matrix
//
// param x The x-component of the scaling vector
// param y The y-component of the scaling vector
// param z The z-component of the scaling vector.
func (self *Matrix4) SetToScaling(x, y, z float32) *Matrix4 {
	self.Idt()
	self.val[M4_00] = x
	self.val[M4_11] = y
	self.val[M4_22] = z
	return self
}

// Sets the matrix to a look at matrix with a direction and an up vector. Multiply with a translation matrix to get a camera
// model view matrix.
//
// param direction The direction vector
// param up The up vector
func (self *Matrix4) SetToLookAt(direction, up *Vector3) *Matrix4 {
	l_vez.SetV(direction).Nor()
	l_vex.SetV(direction).Nor()
	l_vex.CrsV(up).Nor()
	l_vey.SetV(l_vex).CrsV(l_vez).Nor()
	self.Idt()
	self.val[M4_00] = l_vex.x
	self.val[M4_01] = l_vex.y
	self.val[M4_02] = l_vex.z
	self.val[M4_10] = l_vey.x
	self.val[M4_11] = l_vey.y
	self.val[M4_12] = l_vey.z
	self.val[M4_20] = -l_vez.x
	self.val[M4_21] = -l_vez.y
	self.val[M4_22] = -l_vez.z

	return self
}

// Sets this matrix to a look at matrix with the given position, target and up vector.
//
// param position the position
// param target the target
// param up the up vector
func (self *Matrix4) SetToLookAtPos(position, target, up *Vector3) *Matrix4 {
	tmpVec.SetV(target).SubV(position)
	self.SetToLookAt(tmpVec, up)
	self.MulM4(tmpMat.SetToTranslation(-position.x, -position.y, -position.z))

	return self
}

func (self *Matrix4) SetToWorld(position, forward, up *Vector3) *Matrix4 {
	tmpForward.SetV(forward).Nor()
	right.SetV(tmpForward).CrsV(up).Nor()
	tmpUp.SetV(right).CrsV(tmpForward).Nor()

	self.SetVAxis(right, tmpUp, tmpForward.SclValue(-1), position)
	return self
}

// Linearly interpolates between this matrix and the given matrix mixing by alpha
// param matrix the matrix
// param alpha the alpha value in the range [0,1]
func (self *Matrix4) Lerp(matrix *Matrix4, alpha float32) *Matrix4 {
	for i := 0; i < 16; i++ {
		self.val[i] = self.val[i]*(1-alpha) + matrix.val[i]*alpha
	}
	return self
}

// Averages the given transform with this one and stores the result in this matrix. Translations and scales are lerped while
// rotations are slerped.
// param other The other transform
// param w Weight of this transform; weight of the other transform is (1 - w)
func (self *Matrix4) Avg(other *Matrix4, w float32) *Matrix4 {
	self.GetScale(tmpVec)
	other.GetScale(tmpForward)

	self.GetRotationQ(quat)
	other.GetRotationQ(quat2)

	self.GetTranslation(tmpUp)
	other.GetTranslation(right)

	self.SetToScalingV3(tmpVec.SclValue(w).AddV(tmpForward.SclValue(1 - w)))
	self.RotateQ(quat.Slerp(quat2, 1-w))
	self.SetTranslationV3(tmpUp.SclValue(w).AddV(right.SclValue(1 - w)))

	return self
}

// Averages the given transforms and stores the result in this matrix. Translations and scales are lerped while rotations are
// slerped. Does not destroy the data contained in t.
// param t List of transforms
func (self *Matrix4) AvgM4(t []*Matrix4) *Matrix4 {
	w := float32(1 / len(t))

	tmpVec.SetV(t[0].GetScale(tmpUp).SclValue(w))
	quat.SetQ(t[0].GetRotationQ(quat2).Exp(w))
	tmpForward.SetV(t[0].GetTranslation(tmpUp).SclValue(w))

	for i := 1; i < len(t); i++ {
		tmpVec.AddV(t[i].GetScale(tmpUp).SclValue(w))
		quat.MulQ(t[i].GetRotationQ(quat2).Exp(w))
		tmpForward.AddV(t[i].GetTranslation(tmpUp).SclValue(w))
	}
	quat.Nor()

	self.SetToScalingV3(tmpVec)
	self.RotateQ(quat)
	self.SetTranslationV3(tmpForward)

	return self
}

// Averages the given transforms with the given weights and stores the result in this matrix. Translations and scales are
// lerped while rotations are slerped. Does not destroy the data contained in t or w; Sum of w_i must be equal to 1, or
// unexpected results will occur.
// param t List of transforms
// param w List of weights
func (self *Matrix4) AvgM4W(t []*Matrix4, w []float32) *Matrix4 {
	tmpVec.SetV(t[0].GetScale(tmpUp).SclValue(w[0]))
	quat.SetQ(t[0].GetRotationQ(quat2).Exp(w[0]))
	tmpForward.SetV(t[0].GetTranslation(tmpUp).SclValue(w[0]))

	for i := 1; i < len(t); i++ {
		tmpVec.AddV(t[i].GetScale(tmpUp).SclValue(w[i]))
		quat.MulQ(t[i].GetRotationQ(quat2).Exp(w[i]))
		tmpForward.AddV(t[i].GetTranslation(tmpUp).SclValue(w[i]))
	}
	quat.Nor()

	self.SetToScalingV3(tmpVec)
	self.RotateQ(quat)
	self.SetTranslationV3(tmpForward)

	return self
}

// Sets this matrix to the given 3x3 matrix. The third column of this matrix is set to (0,0,1,0).
// param mat the matrix
func (self *Matrix4) SetM3(mat *Matrix3) *Matrix4 {
	self.val[0] = mat.val[0]
	self.val[1] = mat.val[1]
	self.val[2] = mat.val[2]
	self.val[3] = 0
	self.val[4] = mat.val[3]
	self.val[5] = mat.val[4]
	self.val[6] = mat.val[5]
	self.val[7] = 0
	self.val[8] = 0
	self.val[9] = 0
	self.val[10] = 1
	self.val[11] = 0
	self.val[12] = mat.val[6]
	self.val[13] = mat.val[7]
	self.val[14] = 0
	self.val[15] = mat.val[8]
	return self
}

// Sets this matrix to the given affine matrix. The values are mapped as follows:
//      [  M4_00  M4_01   0   M4_02  ]
//      [  M4_10  M4_11   0   M4_12  ]
//      [   0    0    1    0   ]
//      [   0    0    0    1   ]
// param affine the affine matrix
func (self *Matrix4) SetA2(affine *Affine2) *Matrix4 {
	self.val[M4_00] = affine.m00
	self.val[M4_10] = affine.m10
	self.val[M4_20] = 0
	self.val[M4_30] = 0
	self.val[M4_01] = affine.m01
	self.val[M4_11] = affine.m11
	self.val[M4_21] = 0
	self.val[M4_31] = 0
	self.val[M4_02] = 0
	self.val[M4_12] = 0
	self.val[M4_22] = 1
	self.val[M4_32] = 0
	self.val[M4_03] = affine.m02
	self.val[M4_13] = affine.m12
	self.val[M4_23] = 0
	self.val[M4_33] = 1
	return self
}

// Assumes that this matrix is a 2D affine transformation, copying only the relevant components. The values are mapped as
// follows:
//     [  M4_00  M4_01   _   M4_02  ]
//     [  M4_10  M4_11   _   M4_12  ]
//     [   _    _    _    _   ]
//     [   _    _    _    _   ]
// param affine the source matrix
func (self *Matrix4) SetAsAffine(affine *Affine2) *Matrix4 {
	self.val[M4_00] = affine.m00
	self.val[M4_10] = affine.m10
	self.val[M4_01] = affine.m01
	self.val[M4_11] = affine.m11
	self.val[M4_03] = affine.m02
	self.val[M4_13] = affine.m12
	return self
}

// Assumes that both matrices are 2D affine transformations, copying only the relevant components. The copied values are:
//      [  M4_00  M4_01   _   M4_03  ]
//      [  M4_10  M4_11   _   M4_13  ]
//      [   _    _    _    _   ]
//      [   _    _    _    _   ]
// param mat the source matrix
func (self *Matrix4) SetAsAffineM4(mat *Matrix4) *Matrix4 {
	self.val[M4_00] = mat.val[M4_00]
	self.val[M4_10] = mat.val[M4_10]
	self.val[M4_01] = mat.val[M4_01]
	self.val[M4_11] = mat.val[M4_11]
	self.val[M4_03] = mat.val[M4_03]
	self.val[M4_13] = mat.val[M4_13]
	return self
}

func (self *Matrix4) SclV(scale *Vector3) *Matrix4 {
	self.val[M4_00] *= scale.x
	self.val[M4_11] *= scale.y
	self.val[M4_22] *= scale.z
	return self
}

func (self *Matrix4) Scl(x, y, z float32) *Matrix4 {
	self.val[M4_00] *= x
	self.val[M4_11] *= y
	self.val[M4_22] *= z
	return self
}

func (self *Matrix4) SclValue(scale float32) *Matrix4 {
	self.val[M4_00] *= scale
	self.val[M4_11] *= scale
	self.val[M4_22] *= scale
	return self
}

func (self *Matrix4) GetTranslation(position *Vector3) *Vector3 {
	position.x = self.val[M4_03]
	position.y = self.val[M4_13]
	position.z = self.val[M4_23]
	return position
}

// Gets the rotation of this matrix.
// param rotation The {@link Quaternion} to receive the rotation
// param normalizeAxes True to normalize the axes, necessary when the matrix might also include scaling.
// return The provided {@link Quaternion} for chaining.
func (self *Matrix4) GetRotation(rotation *Quaternion, normalizeAxes bool) *Quaternion {
	return rotation.SetFromMatrixM4Normalize(normalizeAxes, self)
}

// Gets the rotation of this matrix.
// param rotation The {@link Quaternion} to receive the rotation
// @return The provided {@link Quaternion} for chaining.
func (self *Matrix4) GetRotationQ(rotation *Quaternion) *Quaternion {
	return rotation.SetFromMatrixM4(self)
}

// @return the squared scale factor on the X axis
func (self *Matrix4) GetScaleXSquared() float32 {
	return self.val[M4_00]*self.val[M4_00] + self.val[M4_01]*self.val[M4_01] + self.val[M4_02]*self.val[M4_02]
}

// @return the squared scale factor on the Y axis
func (self *Matrix4) GetScaleYSquared() float32 {
	return self.val[M4_10]*self.val[M4_10] + self.val[M4_11]*self.val[M4_11] + self.val[M4_12]*self.val[M4_12]
}

// @return the squared scale factor on the Z axis
func (self *Matrix4) GetScaleZSquared() float32 {
	return self.val[M4_20]*self.val[M4_20] + self.val[M4_21]*self.val[M4_21] + self.val[M4_22]*self.val[M4_22]
}

// @return the scale factor on the X axis (non-negative)
func (self *Matrix4) GetScaleX() float32 {
	if utils.IsZero(self.val[M4_01]) && utils.IsZero(self.val[M4_02]) {
		return float32(math.Abs(float64((self.val[M4_00]))))
	}
	return float32(math.Sqrt(float64((self.GetScaleXSquared()))))
}

// @return the scale factor on the Y axis (non-negative)
func (self *Matrix4) GetScaleY() float32 {
	if utils.IsZero(self.val[M4_10]) && utils.IsZero(self.val[M4_12]) {
		return float32(math.Abs(float64((self.val[M4_11]))))
	}
	return float32(math.Sqrt(float64((self.GetScaleYSquared()))))
}

// @return the scale factor on the X axis (non-negative)
func (self *Matrix4) GetScaleZ() float32 {
	if utils.IsZero(self.val[M4_20]) && utils.IsZero(self.val[M4_21]) {
		return float32(math.Abs(float64((self.val[M4_22]))))
	}
	return float32(math.Sqrt(float64((self.GetScaleZSquared()))))
}

// param scale The vector which will receive the (non-negative) scale components on each axis.
// return The provided vector for chaining.
func (self *Matrix4) GetScale(scale *Vector3) *Vector3 {
	return scale.Set(self.GetScaleX(), self.GetScaleY(), self.GetScaleZ())
}

// removes the translational part and transposes the matrix.
func (self *Matrix4) ToNormalMatrix() *Matrix4 {
	self.val[M4_03] = 0
	self.val[M4_13] = 0
	self.val[M4_23] = 0
	return self.Inv().Tra()
}

// JNI TO GO
// Multiplies the matrix mata with matrix matb, storing the result in mata. The arrays are assumed to hold 4x4 column major
// matrices as you can get from {@link Matrix4#val}. This is the same as {@link Matrix4#mul(Matrix4)}.
//
// param mata the first matrix.
// param matb the second matrix.
// public static native void mul (float[] mata, float[] matb) /*-{ }; /*
// matrix4_mul(mata, matb);
func MulM4(mata [16]float32, matb [16]float32) {
	var tmp [16]float32
	tmp[M4_00] = mata[M4_00]*matb[M4_00] + mata[M4_01]*matb[M4_10] + mata[M4_02]*matb[M4_20] + mata[M4_03]*matb[M4_30]
	tmp[M4_01] = mata[M4_00]*matb[M4_01] + mata[M4_01]*matb[M4_11] + mata[M4_02]*matb[M4_21] + mata[M4_03]*matb[M4_31]
	tmp[M4_02] = mata[M4_00]*matb[M4_02] + mata[M4_01]*matb[M4_12] + mata[M4_02]*matb[M4_22] + mata[M4_03]*matb[M4_32]
	tmp[M4_03] = mata[M4_00]*matb[M4_03] + mata[M4_01]*matb[M4_13] + mata[M4_02]*matb[M4_23] + mata[M4_03]*matb[M4_33]
	tmp[M4_10] = mata[M4_10]*matb[M4_00] + mata[M4_11]*matb[M4_10] + mata[M4_12]*matb[M4_20] + mata[M4_13]*matb[M4_30]
	tmp[M4_11] = mata[M4_10]*matb[M4_01] + mata[M4_11]*matb[M4_11] + mata[M4_12]*matb[M4_21] + mata[M4_13]*matb[M4_31]
	tmp[M4_12] = mata[M4_10]*matb[M4_02] + mata[M4_11]*matb[M4_12] + mata[M4_12]*matb[M4_22] + mata[M4_13]*matb[M4_32]
	tmp[M4_13] = mata[M4_10]*matb[M4_03] + mata[M4_11]*matb[M4_13] + mata[M4_12]*matb[M4_23] + mata[M4_13]*matb[M4_33]
	tmp[M4_20] = mata[M4_20]*matb[M4_00] + mata[M4_21]*matb[M4_10] + mata[M4_22]*matb[M4_20] + mata[M4_23]*matb[M4_30]
	tmp[M4_21] = mata[M4_20]*matb[M4_01] + mata[M4_21]*matb[M4_11] + mata[M4_22]*matb[M4_21] + mata[M4_23]*matb[M4_31]
	tmp[M4_22] = mata[M4_20]*matb[M4_02] + mata[M4_21]*matb[M4_12] + mata[M4_22]*matb[M4_22] + mata[M4_23]*matb[M4_32]
	tmp[M4_23] = mata[M4_20]*matb[M4_03] + mata[M4_21]*matb[M4_13] + mata[M4_22]*matb[M4_23] + mata[M4_23]*matb[M4_33]
	tmp[M4_30] = mata[M4_30]*matb[M4_00] + mata[M4_31]*matb[M4_10] + mata[M4_32]*matb[M4_20] + mata[M4_33]*matb[M4_30]
	tmp[M4_31] = mata[M4_30]*matb[M4_01] + mata[M4_31]*matb[M4_11] + mata[M4_32]*matb[M4_21] + mata[M4_33]*matb[M4_31]
	tmp[M4_32] = mata[M4_30]*matb[M4_02] + mata[M4_31]*matb[M4_12] + mata[M4_32]*matb[M4_22] + mata[M4_33]*matb[M4_32]
	tmp[M4_33] = mata[M4_30]*matb[M4_03] + mata[M4_31]*matb[M4_13] + mata[M4_32]*matb[M4_23] + mata[M4_33]*matb[M4_33]
	copy(mata, tmp)
}

// @off
/*JNI
#include <memory.h>
#include <stdio.h>
#include <string.h>

#define M4_00 0
#define M4_01 4
#define M4_02 8
#define M4_03 12
#define M4_10 1
#define M4_11 5
#define M4_12 9
#define M4_13 13
#define M4_20 2
#define M4_21 6
#define M4_22 10
#define M4_23 14
#define M4_30 3
#define M4_31 7
#define M4_32 11
#define M4_33 15

static inline void matrix4_mul(float* mata, float* matb) *Matrix4 {
	float tmp[16];
	tmp[M4_00] = mata[M4_00] * matb[M4_00] + mata[M4_01] * matb[M4_10] + mata[M4_02] * matb[M4_20] + mata[M4_03] * matb[M4_30];
	tmp[M4_01] = mata[M4_00] * matb[M4_01] + mata[M4_01] * matb[M4_11] + mata[M4_02] * matb[M4_21] + mata[M4_03] * matb[M4_31];
	tmp[M4_02] = mata[M4_00] * matb[M4_02] + mata[M4_01] * matb[M4_12] + mata[M4_02] * matb[M4_22] + mata[M4_03] * matb[M4_32];
	tmp[M4_03] = mata[M4_00] * matb[M4_03] + mata[M4_01] * matb[M4_13] + mata[M4_02] * matb[M4_23] + mata[M4_03] * matb[M4_33];
	tmp[M4_10] = mata[M4_10] * matb[M4_00] + mata[M4_11] * matb[M4_10] + mata[M4_12] * matb[M4_20] + mata[M4_13] * matb[M4_30];
	tmp[M4_11] = mata[M4_10] * matb[M4_01] + mata[M4_11] * matb[M4_11] + mata[M4_12] * matb[M4_21] + mata[M4_13] * matb[M4_31];
	tmp[M4_12] = mata[M4_10] * matb[M4_02] + mata[M4_11] * matb[M4_12] + mata[M4_12] * matb[M4_22] + mata[M4_13] * matb[M4_32];
	tmp[M4_13] = mata[M4_10] * matb[M4_03] + mata[M4_11] * matb[M4_13] + mata[M4_12] * matb[M4_23] + mata[M4_13] * matb[M4_33];
	tmp[M4_20] = mata[M4_20] * matb[M4_00] + mata[M4_21] * matb[M4_10] + mata[M4_22] * matb[M4_20] + mata[M4_23] * matb[M4_30];
	tmp[M4_21] = mata[M4_20] * matb[M4_01] + mata[M4_21] * matb[M4_11] + mata[M4_22] * matb[M4_21] + mata[M4_23] * matb[M4_31];
	tmp[M4_22] = mata[M4_20] * matb[M4_02] + mata[M4_21] * matb[M4_12] + mata[M4_22] * matb[M4_22] + mata[M4_23] * matb[M4_32];
	tmp[M4_23] = mata[M4_20] * matb[M4_03] + mata[M4_21] * matb[M4_13] + mata[M4_22] * matb[M4_23] + mata[M4_23] * matb[M4_33];
	tmp[M4_30] = mata[M4_30] * matb[M4_00] + mata[M4_31] * matb[M4_10] + mata[M4_32] * matb[M4_20] + mata[M4_33] * matb[M4_30];
	tmp[M4_31] = mata[M4_30] * matb[M4_01] + mata[M4_31] * matb[M4_11] + mata[M4_32] * matb[M4_21] + mata[M4_33] * matb[M4_31];
	tmp[M4_32] = mata[M4_30] * matb[M4_02] + mata[M4_31] * matb[M4_12] + mata[M4_32] * matb[M4_22] + mata[M4_33] * matb[M4_32];
	tmp[M4_33] = mata[M4_30] * matb[M4_03] + mata[M4_31] * matb[M4_13] + mata[M4_32] * matb[M4_23] + mata[M4_33] * matb[M4_33];
	memcpy(mata, tmp, sizeof(float) *  16);
}

static inline float matrix4_det(float* val) *Matrix4 {
	return val[M4_30] * val[M4_21] * val[M4_12] * val[M4_03] - val[M4_20] * val[M4_31] * val[M4_12] * val[M4_03] - val[M4_30] * val[M4_11]
			* val[M4_22] * val[M4_03] + val[M4_10] * val[M4_31] * val[M4_22] * val[M4_03] + val[M4_20] * val[M4_11] * val[M4_32] * val[M4_03] - val[M4_10]
			* val[M4_21] * val[M4_32] * val[M4_03] - val[M4_30] * val[M4_21] * val[M4_02] * val[M4_13] + val[M4_20] * val[M4_31] * val[M4_02] * val[M4_13]
			+ val[M4_30] * val[M4_01] * val[M4_22] * val[M4_13] - val[M4_00] * val[M4_31] * val[M4_22] * val[M4_13] - val[M4_20] * val[M4_01] * val[M4_32]
			* val[M4_13] + val[M4_00] * val[M4_21] * val[M4_32] * val[M4_13] + val[M4_30] * val[M4_11] * val[M4_02] * val[M4_23] - val[M4_10] * val[M4_31]
			* val[M4_02] * val[M4_23] - val[M4_30] * val[M4_01] * val[M4_12] * val[M4_23] + val[M4_00] * val[M4_31] * val[M4_12] * val[M4_23] + val[M4_10]
			* val[M4_01] * val[M4_32] * val[M4_23] - val[M4_00] * val[M4_11] * val[M4_32] * val[M4_23] - val[M4_20] * val[M4_11] * val[M4_02] * val[M4_33]
			+ val[M4_10] * val[M4_21] * val[M4_02] * val[M4_33] + val[M4_20] * val[M4_01] * val[M4_12] * val[M4_33] - val[M4_00] * val[M4_21] * val[M4_12]
			* val[M4_33] - val[M4_10] * val[M4_01] * val[M4_22] * val[M4_33] + val[M4_00] * val[M4_11] * val[M4_22] * val[M4_33];
}

static inline bool matrix4_inv(float* val) *Matrix4 {
	float tmp[16];
	float l_det = matrix4_det(val);
	if (l_det == 0) return false;
	tmp[M4_00] = val[M4_12] * val[M4_23] * val[M4_31] - val[M4_13] * val[M4_22] * val[M4_31] + val[M4_13] * val[M4_21] * val[M4_32] - val[M4_11]
		* val[M4_23] * val[M4_32] - val[M4_12] * val[M4_21] * val[M4_33] + val[M4_11] * val[M4_22] * val[M4_33];
	tmp[M4_01] = val[M4_03] * val[M4_22] * val[M4_31] - val[M4_02] * val[M4_23] * val[M4_31] - val[M4_03] * val[M4_21] * val[M4_32] + val[M4_01]
		* val[M4_23] * val[M4_32] + val[M4_02] * val[M4_21] * val[M4_33] - val[M4_01] * val[M4_22] * val[M4_33];
	tmp[M4_02] = val[M4_02] * val[M4_13] * val[M4_31] - val[M4_03] * val[M4_12] * val[M4_31] + val[M4_03] * val[M4_11] * val[M4_32] - val[M4_01]
		* val[M4_13] * val[M4_32] - val[M4_02] * val[M4_11] * val[M4_33] + val[M4_01] * val[M4_12] * val[M4_33];
	tmp[M4_03] = val[M4_03] * val[M4_12] * val[M4_21] - val[M4_02] * val[M4_13] * val[M4_21] - val[M4_03] * val[M4_11] * val[M4_22] + val[M4_01]
		* val[M4_13] * val[M4_22] + val[M4_02] * val[M4_11] * val[M4_23] - val[M4_01] * val[M4_12] * val[M4_23];
	tmp[M4_10] = val[M4_13] * val[M4_22] * val[M4_30] - val[M4_12] * val[M4_23] * val[M4_30] - val[M4_13] * val[M4_20] * val[M4_32] + val[M4_10]
		* val[M4_23] * val[M4_32] + val[M4_12] * val[M4_20] * val[M4_33] - val[M4_10] * val[M4_22] * val[M4_33];
	tmp[M4_11] = val[M4_02] * val[M4_23] * val[M4_30] - val[M4_03] * val[M4_22] * val[M4_30] + val[M4_03] * val[M4_20] * val[M4_32] - val[M4_00]
		* val[M4_23] * val[M4_32] - val[M4_02] * val[M4_20] * val[M4_33] + val[M4_00] * val[M4_22] * val[M4_33];
	tmp[M4_12] = val[M4_03] * val[M4_12] * val[M4_30] - val[M4_02] * val[M4_13] * val[M4_30] - val[M4_03] * val[M4_10] * val[M4_32] + val[M4_00]
		* val[M4_13] * val[M4_32] + val[M4_02] * val[M4_10] * val[M4_33] - val[M4_00] * val[M4_12] * val[M4_33];
	tmp[M4_13] = val[M4_02] * val[M4_13] * val[M4_20] - val[M4_03] * val[M4_12] * val[M4_20] + val[M4_03] * val[M4_10] * val[M4_22] - val[M4_00]
		* val[M4_13] * val[M4_22] - val[M4_02] * val[M4_10] * val[M4_23] + val[M4_00] * val[M4_12] * val[M4_23];
	tmp[M4_20] = val[M4_11] * val[M4_23] * val[M4_30] - val[M4_13] * val[M4_21] * val[M4_30] + val[M4_13] * val[M4_20] * val[M4_31] - val[M4_10]
		* val[M4_23] * val[M4_31] - val[M4_11] * val[M4_20] * val[M4_33] + val[M4_10] * val[M4_21] * val[M4_33];
	tmp[M4_21] = val[M4_03] * val[M4_21] * val[M4_30] - val[M4_01] * val[M4_23] * val[M4_30] - val[M4_03] * val[M4_20] * val[M4_31] + val[M4_00]
		* val[M4_23] * val[M4_31] + val[M4_01] * val[M4_20] * val[M4_33] - val[M4_00] * val[M4_21] * val[M4_33];
	tmp[M4_22] = val[M4_01] * val[M4_13] * val[M4_30] - val[M4_03] * val[M4_11] * val[M4_30] + val[M4_03] * val[M4_10] * val[M4_31] - val[M4_00]
		* val[M4_13] * val[M4_31] - val[M4_01] * val[M4_10] * val[M4_33] + val[M4_00] * val[M4_11] * val[M4_33];
	tmp[M4_23] = val[M4_03] * val[M4_11] * val[M4_20] - val[M4_01] * val[M4_13] * val[M4_20] - val[M4_03] * val[M4_10] * val[M4_21] + val[M4_00]
		* val[M4_13] * val[M4_21] + val[M4_01] * val[M4_10] * val[M4_23] - val[M4_00] * val[M4_11] * val[M4_23];
	tmp[M4_30] = val[M4_12] * val[M4_21] * val[M4_30] - val[M4_11] * val[M4_22] * val[M4_30] - val[M4_12] * val[M4_20] * val[M4_31] + val[M4_10]
		* val[M4_22] * val[M4_31] + val[M4_11] * val[M4_20] * val[M4_32] - val[M4_10] * val[M4_21] * val[M4_32];
	tmp[M4_31] = val[M4_01] * val[M4_22] * val[M4_30] - val[M4_02] * val[M4_21] * val[M4_30] + val[M4_02] * val[M4_20] * val[M4_31] - val[M4_00]
		* val[M4_22] * val[M4_31] - val[M4_01] * val[M4_20] * val[M4_32] + val[M4_00] * val[M4_21] * val[M4_32];
	tmp[M4_32] = val[M4_02] * val[M4_11] * val[M4_30] - val[M4_01] * val[M4_12] * val[M4_30] - val[M4_02] * val[M4_10] * val[M4_31] + val[M4_00]
		* val[M4_12] * val[M4_31] + val[M4_01] * val[M4_10] * val[M4_32] - val[M4_00] * val[M4_11] * val[M4_32];
	tmp[M4_33] = val[M4_01] * val[M4_12] * val[M4_20] - val[M4_02] * val[M4_11] * val[M4_20] + val[M4_02] * val[M4_10] * val[M4_21] - val[M4_00]
		* val[M4_12] * val[M4_21] - val[M4_01] * val[M4_10] * val[M4_22] + val[M4_00] * val[M4_11] * val[M4_22];

	float inv_det = 1.0f / l_det;
	self.val[M4_00] = tmp[M4_00] * inv_det;
	self.val[M4_01] = tmp[M4_01] * inv_det;
	self.val[M4_02] = tmp[M4_02] * inv_det;
	self.val[M4_03] = tmp[M4_03] * inv_det;
	self.val[M4_10] = tmp[M4_10] * inv_det;
	self.val[M4_11] = tmp[M4_11] * inv_det;
	self.val[M4_12] = tmp[M4_12] * inv_det;
	self.val[M4_13] = tmp[M4_13] * inv_det;
	self.val[M4_20] = tmp[M4_20] * inv_det;
	self.val[M4_21] = tmp[M4_21] * inv_det;
	self.val[M4_22] = tmp[M4_22] * inv_det;
	self.val[M4_23] = tmp[M4_23] * inv_det;
	self.val[M4_30] = tmp[M4_30] * inv_det;
	self.val[M4_31] = tmp[M4_31] * inv_det;
	self.val[M4_32] = tmp[M4_32] * inv_det;
	self.val[M4_33] = tmp[M4_33] * inv_det;
	return true;
}

static inline void matrix4_mulVec(float* mat, float* vec) *Matrix4 {
	float x = vec[0] * mat[M4_00] + vec[1] * mat[M4_01] + vec[2] * mat[M4_02] + mat[M4_03];
	float y = vec[0] * mat[M4_10] + vec[1] * mat[M4_11] + vec[2] * mat[M4_12] + mat[M4_13];
	float z = vec[0] * mat[M4_20] + vec[1] * mat[M4_21] + vec[2] * mat[M4_22] + mat[M4_23];
	vec[0] = x;
	vec[1] = y;
	vec[2] = z;
}

static inline void matrix4_proj(float* mat, float* vec) *Matrix4 {
	float inv_w = 1.0f / (vec[0] * mat[M4_30] + vec[1] * mat[M4_31] + vec[2] * mat[M4_32] + mat[M4_33]);
	float x = (vec[0] * mat[M4_00] + vec[1] * mat[M4_01] + vec[2] * mat[M4_02] + mat[M4_03]) * inv_w;
	float y = (vec[0] * mat[M4_10] + vec[1] * mat[M4_11] + vec[2] * mat[M4_12] + mat[M4_13]) * inv_w;
	float z = (vec[0] * mat[M4_20] + vec[1] * mat[M4_21] + vec[2] * mat[M4_22] + mat[M4_23]) * inv_w;
	vec[0] = x;
	vec[1] = y;
	vec[2] = z;
}

static inline void matrix4_rot(float* mat, float* vec) *Matrix4 {
	float x = vec[0] * mat[M4_00] + vec[1] * mat[M4_01] + vec[2] * mat[M4_02];
	float y = vec[0] * mat[M4_10] + vec[1] * mat[M4_11] + vec[2] * mat[M4_12];
	float z = vec[0] * mat[M4_20] + vec[1] * mat[M4_21] + vec[2] * mat[M4_22];
	vec[0] = x;
	vec[1] = y;
	vec[2] = z;
}

// Multiplies the vector with the given matrix. The matrix array is assumed to hold a 4x4 column major matrix as you can get
 * from {@link Matrix4#val}. The vector array is assumed to hold a 3-component vector, with x being the first element, y being
 * the second and z being the last component. The result is stored in the vector array. This is the same as
 * {@link Vector3#mul(Matrix4)}.
// param mat the matrix
// param vec the vector.
public static native void mulVec (float[] mat, float[] vec) /*-{ }; /*
	matrix4_mulVec(mat, vec);


// Multiplies the vectors with the given matrix. The matrix array is assumed to hold a 4x4 column major matrix as you can get
 * from {@link Matrix4#val}. The vectors array is assumed to hold 3-component vectors. Offset specifies the offset into the
 * array where the x-component of the first vector is located. The numVecs parameter specifies the number of vectors stored in
 * the vectors array. The stride parameter specifies the number of floats between subsequent vectors and must be >= 3. This is
 * the same as {@link Vector3#mul(Matrix4)} applied to multiple vectors.
 *
// param mat the matrix
// param vecs the vectors
// param offset the offset into the vectors array
// param numVecs the number of vectors
// param stride the stride between vectors in floats
public static native void mulVec (float[] mat, float[] vecs, int offset, int numVecs, int stride) /*-{ }; /*
	float* vecPtr = vecs + offset;
	for(int i = 0; i < numVecs; i++) *Matrix4 {
		matrix4_mulVec(mat, vecPtr);
		vecPtr += stride;
	}


// Multiplies the vector with the given matrix, performing a division by w. The matrix array is assumed to hold a 4x4 column
 * major matrix as you can get from {@link Matrix4#val}. The vector array is assumed to hold a 3-component vector, with x being
 * the first element, y being the second and z being the last component. The result is stored in the vector array. This is the
 * same as {@link Vector3#prj(Matrix4)}.
// param mat the matrix
// param vec the vector.
public static native void prj (float[] mat, float[] vec) /*-{ }; /*
	matrix4_proj(mat, vec);


// Multiplies the vectors with the given matrix, , performing a division by w. The matrix array is assumed to hold a 4x4 column
 * major matrix as you can get from {@link Matrix4#val}. The vectors array is assumed to hold 3-component vectors. Offset
 * specifies the offset into the array where the x-component of the first vector is located. The numVecs parameter specifies
 * the number of vectors stored in the vectors array. The stride parameter specifies the number of floats between subsequent
 * vectors and must be >= 3. This is the same as {@link Vector3#prj(Matrix4)} applied to multiple vectors.
 *
// param mat the matrix
// param vecs the vectors
// param offset the offset into the vectors array
// param numVecs the number of vectors
// param stride the stride between vectors in floats
public static native void prj (float[] mat, float[] vecs, int offset, int numVecs, int stride) /*-{ }; /*
	float* vecPtr = vecs + offset;
	for(int i = 0; i < numVecs; i++) *Matrix4 {
		matrix4_proj(mat, vecPtr);
		vecPtr += stride;
	}


// Multiplies the vector with the top most 3x3 sub-matrix of the given matrix. The matrix array is assumed to hold a 4x4 column
 * major matrix as you can get from {@link Matrix4#val}. The vector array is assumed to hold a 3-component vector, with x being
 * the first element, y being the second and z being the last component. The result is stored in the vector array. This is the
 * same as {@link Vector3#rot(Matrix4)}.
// param mat the matrix
// param vec the vector.
public static native void rot (float[] mat, float[] vec) /*-{ }; /*
	matrix4_rot(mat, vec);


// Multiplies the vectors with the top most 3x3 sub-matrix of the given matrix. The matrix array is assumed to hold a 4x4
 * column major matrix as you can get from {@link Matrix4#val}. The vectors array is assumed to hold 3-component vectors.
 * Offset specifies the offset into the array where the x-component of the first vector is located. The numVecs parameter
 * specifies the number of vectors stored in the vectors array. The stride parameter specifies the number of floats between
 * subsequent vectors and must be >= 3. This is the same as {@link Vector3#rot(Matrix4)} applied to multiple vectors.
 *
// param mat the matrix
// param vecs the vectors
// param offset the offset into the vectors array
// param numVecs the number of vectors
// param stride the stride between vectors in floats
public static native void rot (float[] mat, float[] vecs, int offset, int numVecs, int stride) /*-{ }; /*
	float* vecPtr = vecs + offset;
	for(int i = 0; i < numVecs; i++) *Matrix4 {
		matrix4_rot(mat, vecPtr);
		vecPtr += stride;
	}


// Computes the inverse of the given matrix. The matrix array is assumed to hold a 4x4 column major matrix as you can get from
 * {@link Matrix4#val}.
// param values the matrix values.
 * @return false in case the inverse could not be calculated, true otherwise.
public static native boolean inv (float[] values) /*-{ }; /*
	return matrix4_inv(values);


// Computes the determinante of the given matrix. The matrix array is assumed to hold a 4x4 column major matrix as you can get
 * from {@link Matrix4#val}.
// param values the matrix values.
 * @return the determinante.
public static native float det (float[] values) /*-{ }; /*
	return matrix4_det(values);

// @on
*/
// Postmultiplies this matrix by a translation matrix. Postmultiplication is also used by OpenGL ES'
// glTranslate/glRotate/glScale
// param translation
func (self *Matrix4) TranslateV(translation *Vector3) *Matrix4 {
	return self.Translate(translation.x, translation.y, translation.z)
}

// Postmultiplies this matrix by a translation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
// param x Translation in the x-axis.
// param y Translation in the y-axis.
// param z Translation in the z-axis.
func (self *Matrix4) Translate(x, y, z float32) *Matrix4 {
	tmp[M4_00] = 1
	tmp[M4_01] = 0
	tmp[M4_02] = 0
	tmp[M4_03] = x
	tmp[M4_10] = 0
	tmp[M4_11] = 1
	tmp[M4_12] = 0
	tmp[M4_13] = y
	tmp[M4_20] = 0
	tmp[M4_21] = 0
	tmp[M4_22] = 1
	tmp[M4_23] = z
	tmp[M4_30] = 0
	tmp[M4_31] = 0
	tmp[M4_32] = 0
	tmp[M4_33] = 1

	MulM4(self.val, tmp)
	return self
}

// Postmultiplies this matrix with a (counter-clockwise) rotation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
//
// param axis The vector axis to rotate around.
// param degrees The angle in degrees.
func (self *Matrix4) RotateV3Axis(axis *Vector3, degrees float32) *Matrix4 {
	if degrees == 0 {
		return self
	}
	quat.SetV3(axis, degrees)
	return self.RotateQ(quat)
}

// Postmultiplies this matrix with a (counter-clockwise) rotation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
//
// param axis The vector axis to rotate around.
// param radians The angle in radians.
func (self *Matrix4) RotateRad(axis *Vector3, radians float32) *Matrix4 {
	if radians == 0 {
		return self
	}
	quat.SetFromAxisRadV3(axis, radians)
	return self.RotateQ(quat)
}

// Postmultiplies this matrix with a (counter-clockwise) rotation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale
// param axisX The x-axis component of the vector to rotate around.
// param axisY The y-axis component of the vector to rotate around.
// param axisZ The z-axis component of the vector to rotate around.
// param degrees The angle in degrees
func (self *Matrix4) RotateAxis(axisX, axisY, axisZ, degrees float32) *Matrix4 {
	if degrees == 0 {
		return self
	}
	quat.SetFromAxis(axisX, axisY, axisZ, degrees)
	return self.RotateQ(quat)
}

// Postmultiplies this matrix with a (counter-clockwise) rotation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale
// param axisX The x-axis component of the vector to rotate around.
// param axisY The y-axis component of the vector to rotate around.
// param axisZ The z-axis component of the vector to rotate around.
// param radians The angle in radians
func (self *Matrix4) RotateRadAxis(axisX, axisY, axisZ, radians float32) *Matrix4 {
	if radians == 0 {
		return self
	}
	quat.SetFromAxisRad(axisX, axisY, axisZ, radians)
	return self.RotateQ(quat)
}

// Postmultiplies this matrix with a (counter-clockwise) rotation matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
//
// param rotation
func (self *Matrix4) RotateQ(rotation *Quaternion) *Matrix4 {
	rotation.ToMatrix(tmp)
	MulM4(self.val, tmp)
	return self
}

// Postmultiplies this matrix by the rotation between two vectors.
// param v1 The base vector
// param v2 The target vector
func (self *Matrix4) RotateV3(v1, v2 *Vector3) *Matrix4 {
	return self.RotateQ(quat.SetFromCrossV3(v1, v2))
}

// Postmultiplies this matrix with a scale matrix. Postmultiplication is also used by OpenGL ES' 1.x
// glTranslate/glRotate/glScale.
// param scaleX The scale in the x-axis.
// param scaleY The scale in the y-axis.
// param scaleZ The scale in the z-axis.
func (self *Matrix4) Scale(scaleX, scaleY, scaleZ float32) *Matrix4 {
	tmp[M4_00] = scaleX
	tmp[M4_01] = 0
	tmp[M4_02] = 0
	tmp[M4_03] = 0
	tmp[M4_10] = 0
	tmp[M4_11] = scaleY
	tmp[M4_12] = 0
	tmp[M4_13] = 0
	tmp[M4_20] = 0
	tmp[M4_21] = 0
	tmp[M4_22] = scaleZ
	tmp[M4_23] = 0
	tmp[M4_30] = 0
	tmp[M4_31] = 0
	tmp[M4_32] = 0
	tmp[M4_33] = 1

	MulM4(self.val, tmp)
	return self
}

// Copies the 4x3 upper-left sub-matrix into float array. The destination array is supposed to be a column major matrix.
// param dst the destination matrix
func (self *Matrix4) Extract4x3Matrix(dst []float32) {
	dst[0] = self.val[M4_00]
	dst[1] = self.val[M4_10]
	dst[2] = self.val[M4_20]
	dst[3] = self.val[M4_01]
	dst[4] = self.val[M4_11]
	dst[5] = self.val[M4_21]
	dst[6] = self.val[M4_02]
	dst[7] = self.val[M4_12]
	dst[8] = self.val[M4_22]
	dst[9] = self.val[M4_03]
	dst[10] = self.val[M4_13]
	dst[11] = self.val[M4_23]
}

func (self *Matrix4) String() string {
	return ""
	// return "[" + val[M4_00] + "|" + val[M4_01] + "|" + val[M4_02] + "|" + val[M4_03] + "]\n" + "[" + val[M4_10] + "|" + val[M4_11] + "|"
	// +val[M4_12] + "|" + val[M4_13] + "]\n" + "[" + val[M4_20] + "|" + val[M4_21] + "|" + val[M4_22] + "|" + val[M4_23] + "]\n" + "["
	// +val[M4_30] + "|" + val[M4_31] + "|" + val[M4_32] + "|" + val[M4_33] + "]\n"
}
