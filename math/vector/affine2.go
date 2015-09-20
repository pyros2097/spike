// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package vector

import (
	"github.com/pyros2097/gdx/math/utils"
)

// A specialized 3x3 matrix that can represent sequences of 2D translations, scales, flips, rotations, and shears. <a
// href="http://en.wikipedia.org/wiki/Affine_transformation">Affine transformations</a> preserve straight lines, and
// parallel lines remain parallel after the transformation. Operations on affine matrices are faster because the last row can
// always be assumed (0, 0, 1).
type Affine2 struct {
	m00, m01, m02, m10, m11, m12 float32
}

// constant: m21 = 0, m21 = 1, m22 = 1

// Constructs an identity matrix.
func NewAffine2() *Affine2 {
	return &Affine2{1, 0, 0, 0, 1, 0}
}

// Constructs a matrix from the given affine matrix.
// @param other The affine matrix to copy. This matrix will not be modified.
func NewAffine2Copy(other *Affine2) *Affine2 {
	return NewAffine2().Set(other)
}

// Sets this matrix to the identity matrix
// return This matrix for the purpose of chaining operations.
func (self *Affine2) Idt() *Affine2 {
	self.m00 = 1
	self.m01 = 0
	self.m02 = 0
	self.m10 = 0
	self.m11 = 1
	self.m12 = 0
	return self
}

// Copies the values from the provided affine matrix to this matrix.
// param other The affine matrix to copy.
// return This matrix for the purposes of chaining.
func (self *Affine2) Set(other *Affine2) *Affine2 {
	self.m00 = other.m00
	self.m01 = other.m01
	self.m02 = other.m02
	self.m10 = other.m10
	self.m11 = other.m11
	self.m12 = other.m12
	return self
}

// Copies the values from the provided matrix to this matrix.
// param matrix The matrix to copy, assumed to be an affine transformation.
// return This matrix for the purposes of chaining.
func (self *Affine2) SetM3(matrix *Matrix3) *Affine2 {
	other := matrix.val

	self.m00 = other[M3_00]
	self.m01 = other[M3_01]
	self.m02 = other[M3_02]
	self.m10 = other[M3_10]
	self.m11 = other[M3_11]
	self.m12 = other[M3_12]
	return self
}

// Copies the 2D transformation components from the provided 4x4 matrix. The values are mapped as follows:
//      [  M00  M01  M03  ]
//      [  M10  M11  M13  ]
//      [   0    0    1   ]
// param matrix The source matrix, assumed to be an affine transformation within XY plane. This matrix will not be modified.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetM4(matrix *Matrix4) *Affine2 {
	other := matrix.val

	self.m00 = other[M4_00]
	self.m01 = other[M4_01]
	self.m02 = other[M4_03]
	self.m10 = other[M4_10]
	self.m11 = other[M4_11]
	self.m12 = other[M4_13]
	return self
}

// Sets this matrix to a translation matrix.
// param x The translation in x
// param y The translation in y
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToTranslation(x, y float32) *Affine2 {
	self.m00 = 1
	self.m01 = 0
	self.m02 = x
	self.m10 = 0
	self.m11 = 1
	self.m12 = y
	return self
}

// Sets this matrix to a translation matrix.
// param trn The translation vector.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToTranslationV(trn *Vector2) *Affine2 {
	return self.SetToTranslation(trn.X, trn.Y)
}

// Sets this matrix to a scaling matrix.
// param scaleX The scale in x.
// param scaleY The scale in y.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToScaling(scaleX, scaleY float32) *Affine2 {
	self.m00 = scaleX
	self.m01 = 0
	self.m02 = 0
	self.m10 = 0
	self.m11 = scaleY
	self.m12 = 0
	return self
}

// Sets this matrix to a scaling matrix.
// param scale The scale vector.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToScalingV(scale *Vector2) *Affine2 {
	return self.SetToScaling(scale.X, scale.Y)
}

// Sets this matrix to a rotation matrix that will rotate any vector in counter-clockwise direction around the z-axis.
// param degrees The angle in degrees.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) setToRotation(degrees float32) *Affine2 {
	cos := utils.CosDeg(degrees)
	sin := utils.SinDeg(degrees)

	self.m00 = cos
	self.m01 = -sin
	self.m02 = 0
	self.m10 = sin
	self.m11 = cos
	self.m12 = 0
	return self
}

// Sets this matrix to a rotation matrix that will rotate any vector in counter-clockwise direction around the z-axis.
// param radians The angle in radians.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToRotationRad(radians float32) *Affine2 {
	cos := utils.Cos(radians)
	sin := utils.Sin(radians)

	self.m00 = cos
	self.m01 = -sin
	self.m02 = 0
	self.m10 = sin
	self.m11 = cos
	self.m12 = 0
	return self
}

// Sets this matrix to a rotation matrix that will rotate any vector in counter-clockwise direction around the z-axis.
// param cos The angle cosine.
// param sin The angle sine.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToRotation(cos, sin float32) *Affine2 {
	self.m00 = cos
	self.m01 = -sin
	self.m02 = 0
	self.m10 = sin
	self.m11 = cos
	self.m12 = 0
	return self
}

// Sets this matrix to a shearing matrix.
// param shearX The shear in x direction.
// param shearY The shear in y direction.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToShearing(shearX, shearY float32) *Affine2 {
	self.m00 = 1
	self.m01 = shearX
	self.m02 = 0
	self.m10 = shearY
	self.m11 = 1
	self.m12 = 0
	return self
}

// Sets this matrix to a shearing matrix.
// param shear The shear vector.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToShearingV(shear *Vector2) *Affine2 {
	return self.SetToShearing(shear.X, shear.Y)
}

// Sets this matrix to a concatenation of translation, rotation and scale. It is a more efficient form for:
// <code>idt().translate(x, y).rotate(degrees).scale(scaleX, scaleY)</code>
// param x The translation in x.
// param y The translation in y.
// param degrees The angle in degrees.
// param scaleX The scale in y.
// param scaleY The scale in x.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToTrnRotScl(x, y, degrees, scaleX, scaleY float32) *Affine2 {
	self.m02 = x
	self.m12 = y

	if degrees == 0 {
		self.m00 = scaleX
		self.m01 = 0
		self.m10 = 0
		self.m11 = scaleY
	} else {
		sin := utils.SinDeg(degrees)
		cos := utils.CosDeg(degrees)

		self.m00 = cos * scaleX
		self.m01 = -sin * scaleY
		self.m10 = sin * scaleX
		self.m11 = cos * scaleY
	}
	return self
}

// Sets this matrix to a concatenation of translation, rotation and scale. It is a more efficient form for:
// <code>idt().translate(trn).rotate(degrees).scale(scale)</code>
// param trn The translation vector.
// param degrees The angle in degrees.
// param scale The scale vector.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToTrnRotSclV(trn *Vector2, degrees float32, scale *Vector2) *Affine2 {
	return self.SetToTrnRotScl(trn.X, trn.Y, degrees, scale.X, scale.Y)
}

// Sets this matrix to a concatenation of translation, rotation and scale. It is a more efficient form for:
// <code>idt().translate(x, y).rotateRad(radians).scale(scaleX, scaleY)</code>
// param x The translation in x.
// param y The translation in y.
// param radians The angle in radians.
// param scaleX The scale in y.
// param scaleY The scale in x.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToTrnRotRadScl(x, y, radians, scaleX, scaleY float32) *Affine2 {
	self.m02 = x
	self.m12 = y

	if radians == 0 {
		self.m00 = scaleX
		self.m01 = 0
		self.m10 = 0
		self.m11 = scaleY
	} else {
		sin := utils.Sin(radians)
		cos := utils.Cos(radians)

		self.m00 = cos * scaleX
		self.m01 = -sin * scaleY
		self.m10 = sin * scaleX
		self.m11 = cos * scaleY
	}
	return self
}

// Sets this matrix to a concatenation of translation, rotation and scale. It is a more efficient form for:
// <code>idt().translate(trn).rotateRad(radians).scale(scale)</code>
// param trn The translation vector.
// param radians The angle in radians.
// param scale The scale vector.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToTrnRotRadSclV(trn *Vector2, radians float32, scale *Vector2) *Affine2 {
	return self.SetToTrnRotRadScl(trn.X, trn.Y, radians, scale.X, scale.Y)
}

// Sets this matrix to a concatenation of translation and scale. It is a more efficient form for:
// <code>idt().translate(x, y).scale(scaleX, scaleY)</code>
// param x The translation in x.
// param y The translation in y.
// param scaleX The scale in y.
// param scaleY The scale in x.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToTrnScl(x, y, scaleX, scaleY float32) *Affine2 {
	self.m00 = scaleX
	self.m01 = 0
	self.m02 = x
	self.m10 = 0
	self.m11 = scaleY
	self.m12 = y
	return self
}

// Sets this matrix to a concatenation of translation and scale. It is a more efficient form for:
// idt().translate(trn).scale(scale)</code>
// param trn The translation vector.
// param scale The scale vector.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToTrnSclV(trn, scale *Vector2) *Affine2 {
	return self.SetToTrnScl(trn.X, trn.Y, scale.X, scale.Y)
}

// Sets this matrix to the product of two matrices.
// param l Left matrix.
// param r Right matrix.
// return This matrix for the purpose of chaining operations.
func (self *Affine2) SetToProduct(l, r *Affine2) *Affine2 {
	self.m00 = l.m00*r.m00 + l.m01*r.m10
	self.m01 = l.m00*r.m01 + l.m01*r.m11
	self.m02 = l.m00*r.m02 + l.m01*r.m12 + l.m02
	self.m10 = l.m10*r.m00 + l.m11*r.m10
	self.m11 = l.m10*r.m01 + l.m11*r.m11
	self.m12 = l.m10*r.m02 + l.m11*r.m12 + l.m12
	return self
}

// Inverts this matrix given that the determinant is != 0.
// return This matrix for the purpose of chaining operations.
// @throws GdxRuntimeException if the matrix is singular (not invertible)
func (self *Affine2) Inv() *Affine2 {
	det := self.Det()
	if det == 0 {
		panic("Can't invert a singular affine matrix")
	}

	invDet := 1.0 / det

	tmp00 := self.m11
	tmp01 := -self.m01
	tmp02 := self.m01*self.m12 - self.m11*self.m02
	tmp10 := -self.m10
	tmp11 := self.m00
	tmp12 := self.m10*self.m02 - self.m00*self.m12

	self.m00 = invDet * tmp00
	self.m01 = invDet * tmp01
	self.m02 = invDet * tmp02
	self.m10 = invDet * tmp10
	self.m11 = invDet * tmp11
	self.m12 = invDet * tmp12
	return self
}

// Postmultiplies this matrix with the provided matrix and stores the result in this matrix. For example:
// A.mul(B) results in A := AB
// param other Matrix to multiply by.
// return This matrix for the purpose of chaining operations together.
func (self *Affine2) Mul(other *Affine2) *Affine2 {
	tmp00 := self.m00*other.m00 + self.m01*other.m10
	tmp01 := self.m00*other.m01 + self.m01*other.m11
	tmp02 := self.m00*other.m02 + self.m01*other.m12 + self.m02
	tmp10 := self.m10*other.m00 + self.m11*other.m10
	tmp11 := self.m10*other.m01 + self.m11*other.m11
	tmp12 := self.m10*other.m02 + self.m11*other.m12 + self.m12

	self.m00 = tmp00
	self.m01 = tmp01
	self.m02 = tmp02
	self.m10 = tmp10
	self.m11 = tmp11
	self.m12 = tmp12
	return self
}

// Premultiplies this matrix with the provided matrix and stores the result in this matrix. For example:
// A.preMul(B) results in A := BA
// param other The other Matrix to multiply by
// return This matrix for the purpose of chaining operations.
func (self *Affine2) PreMulA(other *Affine2) *Affine2 {
	tmp00 := other.m00*self.m00 + other.m01*self.m10
	tmp01 := other.m00*self.m01 + other.m01*self.m11
	tmp02 := other.m00*self.m02 + other.m01*self.m12 + other.m02
	tmp10 := other.m10*self.m00 + other.m11*self.m10
	tmp11 := other.m10*self.m01 + other.m11*self.m11
	tmp12 := other.m10*self.m02 + other.m11*self.m12 + other.m12

	self.m00 = tmp00
	self.m01 = tmp01
	self.m02 = tmp02
	self.m10 = tmp10
	self.m11 = tmp11
	self.m12 = tmp12
	return self
}

// Postmultiplies this matrix by a translation matrix.
// param x The x-component of the translation vector.
// param y The y-component of the translation vector.
// return This matrix for the purpose of chaining.
func (self *Affine2) Translate(x, y float32) *Affine2 {
	self.m02 += self.m00*x + self.m01*y
	self.m12 += self.m10*x + self.m11*y
	return self
}

// Postmultiplies this matrix by a translation matrix.
// param trn The translation vector.
// return This matrix for the purpose of chaining.
func (self *Affine2) TranslateV(trn *Vector2) *Affine2 {
	return self.Translate(trn.X, trn.Y)
}

// Premultiplies this matrix by a translation matrix.
// param x The x-component of the translation vector.
// param y The y-component of the translation vector.
// return This matrix for the purpose of chaining.
func (self *Affine2) PreTranslate(x, y float32) *Affine2 {
	self.m02 += x
	self.m12 += y
	return self
}

// Premultiplies this matrix by a translation matrix.
// param trn The translation vector.
// return This matrix for the purpose of chaining.
func (self *Affine2) PreTranslateV(trn *Vector2) *Affine2 {
	return self.PreTranslate(trn.X, trn.Y)
}

// Postmultiplies this matrix with a scale matrix.
// param scaleX The scale in the x-axis.
// param scaleY The scale in the y-axis.
// return This matrix for the purpose of chaining.
func (self *Affine2) Scale(scaleX, scaleY float32) *Affine2 {
	self.m00 *= scaleX
	self.m01 *= scaleY
	self.m10 *= scaleX
	self.m11 *= scaleY
	return self
}

// Postmultiplies this matrix with a scale matrix.
// param scale The scale vector.
// return This matrix for the purpose of chaining.
func (self *Affine2) ScaleV(scale *Vector2) *Affine2 {
	return self.Scale(scale.X, scale.Y)
}

// Premultiplies this matrix with a scale matrix.
// param scaleX The scale in the x-axis.
// param scaleY The scale in the y-axis.
// return This matrix for the purpose of chaining.
func (self *Affine2) PreScale(scaleX, scaleY float32) *Affine2 {
	self.m00 *= scaleX
	self.m01 *= scaleX
	self.m02 *= scaleX
	self.m10 *= scaleY
	self.m11 *= scaleY
	self.m12 *= scaleY
	return self
}

// Premultiplies this matrix with a scale matrix.
// param scale The scale vector.
// return This matrix for the purpose of chaining.
func (self *Affine2) PreScaleV(scale *Vector2) *Affine2 {
	return self.PreScale(scale.X, scale.Y)
}

// Postmultiplies this matrix with a (counter-clockwise) rotation matrix.
// param degrees The angle in degrees
// return This matrix for the purpose of chaining.
func (self *Affine2) Rotate(degrees float32) *Affine2 {
	if degrees == 0 {
		return self
	}

	cos := utils.CosDeg(degrees)
	sin := utils.SinDeg(degrees)

	tmp00 := self.m00*cos + self.m01*sin
	tmp01 := self.m00*-sin + self.m01*cos
	tmp10 := self.m10*cos + self.m11*sin
	tmp11 := self.m10*-sin + self.m11*cos

	self.m00 = tmp00
	self.m01 = tmp01
	self.m10 = tmp10
	self.m11 = tmp11
	return self
}

// Postmultiplies this matrix with a (counter-clockwise) rotation matrix.
// param radians The angle in radians
// return This matrix for the purpose of chaining.
func (self *Affine2) RotateRad(radians float32) *Affine2 {
	if radians == 0 {
		return self
	}

	cos := utils.Cos(radians)
	sin := utils.Sin(radians)

	tmp00 := self.m00*cos + self.m01*sin
	tmp01 := self.m00*-sin + self.m01*cos
	tmp10 := self.m10*cos + self.m11*sin
	tmp11 := self.m10*-sin + self.m11*cos

	self.m00 = tmp00
	self.m01 = tmp01
	self.m10 = tmp10
	self.m11 = tmp11
	return self
}

// Premultiplies this matrix with a (counter-clockwise) rotation matrix.
// param degrees The angle in degrees
// return This matrix for the purpose of chaining.
func (self *Affine2) PreRotate(degrees float32) *Affine2 {
	if degrees == 0 {
		return self
	}

	cos := utils.CosDeg(degrees)
	sin := utils.SinDeg(degrees)

	tmp00 := cos*self.m00 - sin*self.m10
	tmp01 := cos*self.m01 - sin*self.m11
	tmp02 := cos*self.m02 - sin*self.m12
	tmp10 := sin*self.m00 + cos*self.m10
	tmp11 := sin*self.m01 + cos*self.m11
	tmp12 := sin*self.m02 + cos*self.m12

	self.m00 = tmp00
	self.m01 = tmp01
	self.m02 = tmp02
	self.m10 = tmp10
	self.m11 = tmp11
	self.m12 = tmp12
	return self
}

// Premultiplies this matrix with a (counter-clockwise) rotation matrix.
// param radians The angle in radians
// return This matrix for the purpose of chaining.
func (self *Affine2) PreRotateRad(radians float32) *Affine2 {
	if radians == 0 {
		return self
	}

	cos := utils.Cos(radians)
	sin := utils.Sin(radians)

	tmp00 := cos*self.m00 - sin*self.m10
	tmp01 := cos*self.m01 - sin*self.m11
	tmp02 := cos*self.m02 - sin*self.m12
	tmp10 := sin*self.m00 + cos*self.m10
	tmp11 := sin*self.m01 + cos*self.m11
	tmp12 := sin*self.m02 + cos*self.m12

	self.m00 = tmp00
	self.m01 = tmp01
	self.m02 = tmp02
	self.m10 = tmp10
	self.m11 = tmp11
	self.m12 = tmp12
	return self
}

// Postmultiplies this matrix by a shear matrix.
// param shearX The shear in x direction.
// param shearY The shear in y direction.
// return This matrix for the purpose of chaining.
func (self *Affine2) Shear(shearX, shearY float32) *Affine2 {
	tmp0 := self.m00 + shearY*self.m01
	tmp1 := self.m01 + shearX*self.m00
	self.m00 = tmp0
	self.m01 = tmp1

	tmp0 = self.m10 + shearY*self.m11
	tmp1 = self.m11 + shearX*self.m10
	self.m10 = tmp0
	self.m11 = tmp1
	return self
}

// Postmultiplies this matrix by a shear matrix.
// param shear The shear vector.
// return This matrix for the purpose of chaining.
func (self *Affine2) ShearV(shear *Vector2) *Affine2 {
	return self.Shear(shear.X, shear.Y)
}

// Premultiplies this matrix by a shear matrix.
// param shearX The shear in x direction.
// param shearY The shear in y direction.
// return This matrix for the purpose of chaining.
func (self *Affine2) PreShear(shearX, shearY float32) *Affine2 {
	tmp00 := self.m00 + shearX*self.m10
	tmp01 := self.m01 + shearX*self.m11
	tmp02 := self.m02 + shearX*self.m12
	tmp10 := self.m10 + shearY*self.m00
	tmp11 := self.m11 + shearY*self.m01
	tmp12 := self.m12 + shearY*self.m02

	self.m00 = tmp00
	self.m01 = tmp01
	self.m02 = tmp02
	self.m10 = tmp10
	self.m11 = tmp11
	self.m12 = tmp12
	return self
}

// Premultiplies this matrix by a shear matrix.
// param shear The shear vector.
// return This matrix for the purpose of chaining.
func (self *Affine2) PreShearV(shear *Vector2) *Affine2 {
	return self.PreShear(shear.X, shear.Y)
}

// Calculates the determinant of the matrix.
// return The determinant of this matrix.
func (self *Affine2) Det() float32 {
	return self.m00*self.m11 - self.m01*self.m10
}

// Get the x-y translation component of the matrix.
// param position Output vector.
// return Filled position.
func (self *Affine2) GetTranslation(position *Vector2) *Vector2 {
	position.X = self.m02
	position.Y = self.m12
	return position
}

// Check if the this is a plain translation matrix.
// return True if scale is 1 and rotation is 0.
func (self *Affine2) IsTranslation() bool {
	return self.m00 == 1 && self.m11 == 1 && self.m01 == 0 && self.m10 == 0
}

// Check if this is an indentity matrix.
// return True if scale is 1 and rotation is 0.
func (self *Affine2) IsIdt() bool {
	return self.m00 == 1 && self.m02 == 0 && self.m12 == 0 && self.m11 == 1 && self.m01 == 0 && self.m10 == 0
}

// Applies the affine transformation on a vector.
func (self *Affine2) ApplyTo(point *Vector2) {
	x := point.X
	y := point.Y
	point.X = self.m00*x + self.m01*y + self.m02
	point.Y = self.m10*x + self.m11*y + self.m12
}

func (self *Affine2) String() string {
	return ""
	//return "[" + self.m00 + "|" + self.m01 + "|" + self.m02 + "]\n[" + self.m10 + "|" + self.m11 + "|" + self.m12 + "]\n[0.0|0.0|0.1]"
}
