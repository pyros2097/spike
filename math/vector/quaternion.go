// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package vector

import (
	"math"

	"github.com/pyros2097/gdx/math/utils"
)

var (
	tmp1 = NewQuaternion(0, 0, 0, 0)
	tmp2 = NewQuaternion(0, 0, 0, 0)
)

type Quaternion struct {
	x, y, z, w float32
}

// Constructor, sets the four components of the quaternion.
func NewQuaternion(x, y, z, w float32) *Quaternion {
	return &Quaternion{x, y, z, w}
}

func NewQuaternionEmpty() *Quaternion {
	q := &Quaternion{}
	return q.Idt()
}

// Constructor, sets the quaternion components from the given quaternion.
// param quaternion The quaternion to copy.
func NewQuaternionCopy(quaternion *Quaternion) *Quaternion {
	q := &Quaternion{}
	return q.SetQ(quaternion)
}

// Constructor, sets the quaternion from the given axis vector and the angle around that axis in degrees.
// param axis The axis
// param angle The angle in degrees.
func NewQuaternionAxis(axis *Vector3, angle float32) *Quaternion {
	q := &Quaternion{}
	return q.SetV3(axis, angle)
}

// Sets the components of the quaternion
// param x The x-component
// param y The y-component
// param z The z-component
// param w The w-component
func (self *Quaternion) Set(x, y, z, w float32) *Quaternion {
	self.x = x
	self.y = y
	self.z = z
	self.w = w
	return self
}

// Sets the quaternion components from the given quaternion.
// param quaternion The quaternion.
func (self *Quaternion) SetQ(quaternion *Quaternion) *Quaternion {
	return self.Set(quaternion.x, quaternion.y, quaternion.z, quaternion.w)
}

// Sets the quaternion components from the given axis and angle around that axis.
// param axis The axis
// param angle The angle in degrees
func (self *Quaternion) SetV3(axis *Vector3, angle float32) *Quaternion {
	return self.SetFromAxis(axis.X, axis.Y, axis.Z, angle)
}

// return a copy of this quaternion
func (self *Quaternion) Copy() *Quaternion {
	return NewQuaternionCopy(self)
}

// return the euclidean length of the specified quaternion
func LenQ(x, y, z, w float32) float32 {
	return float32(math.Sqrt(float64(x*x + y*y + z*z + w*w)))
}

// return the euclidean length of this quaternion
func (self *Quaternion) Len() float32 {
	return float32(math.Sqrt(float64(self.x*self.x + self.y*self.y + self.z*self.z + self.w*self.w)))
}

// Sets the quaternion to the given euler angles in degrees.
// param yaw the rotation around the y axis in degrees
// param pitch the rotation around the x axis in degrees
// param roll the rotation around the z axis degrees
// return this quaternion
func (self *Quaternion) SetEulerAngles(yaw, pitch, roll float32) *Quaternion {
	return self.SetEulerAnglesRad(yaw*utils.DegreesToRadians, pitch*utils.DegreesToRadians, roll*utils.DegreesToRadians)
}

// Sets the quaternion to the given euler angles in radians.
// param yaw the rotation around the y axis in radians
// param pitch the rotation around the x axis in radians
// param roll the rotation around the z axis in radians
// return this quaternion
func (self *Quaternion) SetEulerAnglesRad(yaw, pitch, roll float32) *Quaternion {
	hr := roll * 0.5
	shr := float32(math.Sin(float64(hr)))
	chr := float32(math.Cos(float64(hr)))
	hp := pitch * 0.5
	shp := float32(math.Sin(float64(hp)))
	chp := float32(math.Cos(float64(hp)))
	hy := yaw * 0.5
	shy := float32(math.Sin(float64(hy)))
	chy := float32(math.Cos(float64(hy)))
	chy_shp := chy * shp
	shy_chp := shy * chp
	chy_chp := chy * chp
	shy_shp := shy * shp

	self.x = (chy_shp * chr) + (shy_chp * shr) // cos(yaw/2) * sin(pitch/2) * cos(roll/2) + sin(yaw/2) * cos(pitch/2) * sin(roll/2)
	self.y = (shy_chp * chr) - (chy_shp * shr) // sin(yaw/2) * cos(pitch/2) * cos(roll/2) - cos(yaw/2) * sin(pitch/2) * sin(roll/2)
	self.z = (chy_chp * shr) - (shy_shp * chr) // cos(yaw/2) * cos(pitch/2) * sin(roll/2) - sin(yaw/2) * sin(pitch/2) * cos(roll/2)
	self.w = (chy_chp * chr) + (shy_shp * shr) // cos(yaw/2) * cos(pitch/2) * cos(roll/2) + sin(yaw/2) * sin(pitch/2) * sin(roll/2)
	return self
}

// Get the pole of the gimbal lock, if any.
// return positive (+1) for north pole, negative (-1) for south pole, zero (0) when no gimbal lock
func (self *Quaternion) GetGimbalPole() int {
	t := self.y*self.x + self.z*self.w
	if t > 0.499 {
		return 1
	} else if t < -0.499 {
		return -1
	} else {
		return 0
	}
}

// Get the roll euler angle in radians, which is the rotation around the z axis. Requires that this quaternion is normalized.
// return the rotation around the z axis in radians (between -PI and +PI)
func (self *Quaternion) GetRollRad() float32 {
	pole := self.GetGimbalPole()
	if pole == 0 {
		return utils.Atan2(2*(self.w*self.z+self.y*self.x), 1-2*(self.x*self.x+self.z*self.z))
	}
	return float32(pole) * 2 * utils.Atan2(self.y, self.w)
}

// Get the roll euler angle in degrees, which is the rotation around the z axis. Requires that this quaternion is normalized.
// return the rotation around the z axis in degrees (between -180 and +180)
func (self *Quaternion) GetRoll() float32 {
	return self.GetRollRad() * utils.RadiansToDegrees
}

// Get the pitch euler angle in radians, which is the rotation around the x axis. Requires that this quaternion is normalized.
// return the rotation around the x axis in radians (between -(PI/2) and +(PI/2))
func (self *Quaternion) GetPitchRad() float32 {
	pole := self.GetGimbalPole()
	if pole == 0 {
		return float32(math.Asin(float64(utils.ClampFloat32(2*(self.w*self.x-self.z*self.y), -1, 1))))
	}
	return float32(pole) * utils.PI * 0.5
}

// Get the pitch euler angle in degrees, which is the rotation around the x axis. Requires that this quaternion is normalized.
// return the rotation around the x axis in degrees (between -90 and +90)
func (self *Quaternion) GetPitch() float32 {
	return self.GetPitchRad() * utils.RadiansToDegrees
}

// Get the yaw euler angle in radians, which is the rotation around the y axis. Requires that this quaternion is normalized.
// return the rotation around the y axis in radians (between -PI and +PI)
func (self *Quaternion) GetYawRad() float32 {
	if self.GetGimbalPole() == 0 {
		return utils.Atan2(2*(self.y*self.w+self.x*self.z), 1-2*(self.y*self.y+self.x*self.x))
	}
	return 0
}

// Get the yaw euler angle in degrees, which is the rotation around the y axis. Requires that this quaternion is normalized.
// return the rotation around the y axis in degrees (between -180 and +180)
func (self *Quaternion) GetYaw() float32 {
	return self.GetYawRad() * utils.RadiansToDegrees
}

func Len2Q(x, y, z, w float32) float32 {
	return x*x + y*y + z*z + w*w
}

// return the length of this quaternion without square root
func (self *Quaternion) Len2() float32 {
	return self.x*self.x + self.y*self.y + self.z*self.z + self.w*self.w
}

// Normalizes this quaternion to unit length
func (self *Quaternion) Nor() *Quaternion {
	length := self.Len2()
	if length != 0 && utils.IsEqual(length, 1) {
		length = float32(math.Sqrt(float64(length)))
		self.w /= length
		self.x /= length
		self.y /= length
		self.z /= length
	}
	return self
}

// Conjugate the quaternion.
func (self *Quaternion) Conjugate() *Quaternion {
	self.x = -self.x
	self.y = -self.y
	self.z = -self.z
	return self
}

// TODO : this would better fit into the vector3 class
// Transforms the given vector using this quaternion
// param v Vector to transform
func (self *Quaternion) Transform(v *Vector3) *Vector3 {
	tmp2.SetQ(self)
	tmp2.Conjugate()
	tmp2.MulLeftQ(tmp1.Set(v.X, v.Y, v.Z, 0)).MulLeftQ(self)

	v.X = tmp2.x
	v.Y = tmp2.y
	v.Z = tmp2.z
	return v
}

// Multiplies this quaternion with another one in the form of this = this * other
// param other Quaternion to multiply with
func (self *Quaternion) MulQ(other *Quaternion) *Quaternion {
	newX := self.w*other.x + self.x*other.w + self.y*other.z - self.z*other.y
	newY := self.w*other.y + self.y*other.w + self.z*other.x - self.x*other.z
	newZ := self.w*other.z + self.z*other.w + self.x*other.y - self.y*other.x
	newW := self.w*other.w - self.x*other.x - self.y*other.y - self.z*other.z
	self.x = newX
	self.y = newY
	self.z = newZ
	self.w = newW
	return self
}

// Multiplies this quaternion with another one in the form of this = this * other
// param x the x component of the other quaternion to multiply with
// param y the y component of the other quaternion to multiply with
// param z the z component of the other quaternion to multiply with
// param w the w component of the other quaternion to multiply with
func (self *Quaternion) Mul(x, y, z, w float32) *Quaternion {
	newX := self.w*x + self.x*w + self.y*z - self.z*y
	newY := self.w*y + self.y*w + self.z*x - self.x*z
	newZ := self.w*z + self.z*w + self.x*y - self.y*x
	newW := self.w*w - self.x*x - self.y*y - self.z*z
	self.x = newX
	self.y = newY
	self.z = newZ
	self.w = newW
	return self
}

// Multiplies this quaternion with another one in the form of this = other * this
// param other Quaternion to multiply with
func (self *Quaternion) MulLeftQ(other *Quaternion) *Quaternion {
	newX := other.w*self.x + other.x*self.w + other.y*self.z - other.z*self.y
	newY := other.w*self.y + other.y*self.w + other.z*self.x - other.x*self.z
	newZ := other.w*self.z + other.z*self.w + other.x*self.y - other.y*self.x
	newW := other.w*self.w - other.x*self.x - other.y*self.y - other.z*self.z
	self.x = newX
	self.y = newY
	self.z = newZ
	self.w = newW
	return self
}

// Multiplies this quaternion with another one in the form of this = other * this
// param x the x component of the other quaternion to multiply with
// param y the y component of the other quaternion to multiply with
// param z the z component of the other quaternion to multiply with
// param w the w component of the other quaternion to multiply with
func (self *Quaternion) MulLeft(x, y, z, w float32) *Quaternion {
	newX := w*self.x + x*self.w + y*self.z - z*y
	newY := w*self.y + y*self.w + z*self.x - x*z
	newZ := w*self.z + z*self.w + x*self.y - y*x
	newW := w*self.w - x*self.x - y*self.y - z*z
	self.x = newX
	self.y = newY
	self.z = newZ
	self.w = newW
	return self
}

// Add the x,y,z,w components of the passed in quaternion to the ones of this quaternion
func (self *Quaternion) AddQ(quaternion *Quaternion) *Quaternion {
	self.x += quaternion.x
	self.y += quaternion.y
	self.z += quaternion.z
	self.w += quaternion.w
	return self
}

// Add the x,y,z,w components of the passed in quaternion to the ones of this quaternion
func (self *Quaternion) Add(qx, qy, qz, qw float32) *Quaternion {
	self.x += qx
	self.y += qy
	self.z += qz
	self.w += qw
	return self
}

// TODO : the matrix4 set(quaternion) doesnt set the last row+col of the matrix to 0,0,0,1 so... that's why there is this
// method
// Fills a 4x4 matrix with the rotation matrix represented by this quaternion.
//
// param matrix Matrix to fill
func (self *Quaternion) ToMatrix(matrix [16]float32) {
	xx := self.x * self.x
	xy := self.x * self.y
	xz := self.x * self.z
	xw := self.x * self.w
	yy := self.y * self.y
	yz := self.y * self.z
	yw := self.y * self.w
	zz := self.z * self.z
	zw := self.z * self.w
	// Set matrix from quaternion
	matrix[M4_00] = 1 - 2*(yy+zz)
	matrix[M4_01] = 2 * (xy - zw)
	matrix[M4_02] = 2 * (xz + yw)
	matrix[M4_03] = 0
	matrix[M4_10] = 2 * (xy + zw)
	matrix[M4_11] = 1 - 2*(xx+zz)
	matrix[M4_12] = 2 * (yz - xw)
	matrix[M4_13] = 0
	matrix[M4_20] = 2 * (xz - yw)
	matrix[M4_21] = 2 * (yz + xw)
	matrix[M4_22] = 1 - 2*(xx+yy)
	matrix[M4_23] = 0
	matrix[M4_30] = 0
	matrix[M4_31] = 0
	matrix[M4_32] = 0
	matrix[M4_33] = 1
}

// Sets the quaternion to an identity Quaternio
func (self *Quaternion) Idt() *Quaternion {
	return self.Set(0, 0, 0, 1)
}

// return If this quaternion is an identity Quaternion
func (self *Quaternion) IsIdentity() bool {
	return utils.IsZero(self.x) && utils.IsZero(self.y) && utils.IsZero(self.z) && utils.IsEqual(self.w, 1)
}

// return If this quaternion is an identity Quaternion
func (self *Quaternion) IsIdentityTolerance(tolerance float32) bool {
	return utils.IsZeroTolerance(self.x, tolerance) && utils.IsZeroTolerance(self.y, tolerance) &&
		utils.IsZeroTolerance(self.z, tolerance) && utils.IsEqualTolerance(self.w, 1, tolerance)
}

// todo : the setFromAxis(v3,float) method should replace the set(v3,float) method
// Sets the quaternion components from the given axis and angle around that axis.
//
// param axis The axis
// param degrees The angle in degrees
func (self *Quaternion) SetFromAxisV3(axis *Vector3, degrees float32) *Quaternion {
	return self.SetFromAxis(axis.X, axis.Y, axis.Z, degrees)
}

// Sets the quaternion components from the given axis and angle around that axis.
// param axis The axis
// param radians The angle in radians
func (self *Quaternion) SetFromAxisRadV3(axis *Vector3, radians float32) *Quaternion {
	return self.SetFromAxisRad(axis.X, axis.Y, axis.Z, radians)
}

// Sets the quaternion components from the given axis and angle around that axis.
// param x X direction of the axis
// param y Y direction of the axis
// param z Z direction of the axis
// param degrees The angle in degrees
func (self *Quaternion) SetFromAxis(x, y, z, degrees float32) *Quaternion {
	return self.SetFromAxisRad(x, y, z, degrees*utils.DegreesToRadians)
}

// Sets the quaternion components from the given axis and angle around that axis.
// param x X direction of the axis
// param y Y direction of the axis
// param z Z direction of the axis
// param radians The angle in radians
func (self *Quaternion) SetFromAxisRad(x, y, z, radians float32) *Quaternion {
	d := LenV3(x, y, z)
	if d == 0 {
		return self.Idt()
	}
	d = 1 / d
	var l_ang float32
	if radians < 0 {
		l_ang = utils.PI2 // - (-radians % utils.PI2) // TODO FIX THIS
	} else {
		l_ang = radians // % utils.PI2 // TODO FIX THIS
	}
	l_sin := float32(math.Sin(float64(l_ang / 2)))
	l_cos := float32(math.Cos(float64(l_ang / 2)))
	return self.Set(d*x*l_sin, d*y*l_sin, d*z*l_sin, l_cos).Nor()
}

// Sets the Quaternion from the given matrix, optionally removing any scaling.
func (self *Quaternion) SetFromMatrixM4Normalize(normalizeAxes bool, matrix *Matrix4) *Quaternion {
	return self.SetFromAxesNormalize(normalizeAxes, matrix.val[M4_00], matrix.val[M4_01], matrix.val[M4_02],
		matrix.val[M4_10], matrix.val[M4_11], matrix.val[M4_12], matrix.val[M4_20],
		matrix.val[M4_21], matrix.val[M4_22])
}

// Sets the Quaternion from the given rotation matrix, which must not contain scaling.
func (self *Quaternion) SetFromMatrixM4(matrix *Matrix4) *Quaternion {
	return self.SetFromMatrixM4Normalize(false, matrix)
}

// Sets the Quaternion from the given matrix, optionally removing any scaling.
func (self *Quaternion) SetFromMatrixM3Normalize(normalizeAxes bool, matrix *Matrix3) *Quaternion {
	return self.SetFromAxesNormalize(normalizeAxes, matrix.val[M3_00], matrix.val[M3_01], matrix.val[M3_02],
		matrix.val[M3_10], matrix.val[M3_11], matrix.val[M3_12], matrix.val[M3_20],
		matrix.val[M3_21], matrix.val[M3_22])
}

// Sets the Quaternion from the given rotation matrix, which must not contain scaling.
func (self *Quaternion) SetFromMatrixM3(matrix *Matrix3) *Quaternion {
	return self.SetFromMatrixM3Normalize(false, matrix)
}

// Sets the Quaternion from the given x-, y- and z-axis which have to be orthonormal.
// Taken from Bones framework for JPCT, see http://www.aptalkarga.com/bones/ which in turn took it from Graphics Gem code at
// ftp://ftp.cis.upenn.edu/pub/graphics/shoemake/quatut.ps.Z.
//
// param xx x-axis x-coordinate
// param xy x-axis y-coordinate
// param xz x-axis z-coordinate
// param yx y-axis x-coordinate
// param yy y-axis y-coordinate
// param yz y-axis z-coordinate
// param zx z-axis x-coordinate
// param zy z-axis y-coordinate
// param zz z-axis z-coordinate
func (self *Quaternion) SetFromAxes(xx, xy, xz, yx, yy, yz, zx, zy, zz float32) *Quaternion {
	return self.SetFromAxesNormalize(false, xx, xy, xz, yx, yy, yz, zx, zy, zz)
}

// Sets the Quaternion from the given x-, y- and z-axis.
// Taken from Bones framework for JPCT, see http://www.aptalkarga.com/bones/ which in turn took it from Graphics Gem code at
// ftp://ftp.cis.upenn.edu/pub/graphics/shoemake/quatut.ps.Z.
//
// param normalizeAxes whether to normalize the axes (necessary when they contain scaling)
// param xx x-axis x-coordinate
// param xy x-axis y-coordinate
// param xz x-axis z-coordinate
// param yx y-axis x-coordinate
// param yy y-axis y-coordinate
// param yz y-axis z-coordinate
// param zx z-axis x-coordinate
// param zy z-axis y-coordinate
// param zz z-axis z-coordinate
func (self *Quaternion) SetFromAxesNormalize(normalizeAxes bool, xx, xy, xz, yx, yy, yz, zx, zy, zz float32) *Quaternion {
	if normalizeAxes {
		lx := 1 / LenV3(xx, xy, xz)
		ly := 1 / LenV3(yx, yy, yz)
		lz := 1 / LenV3(zx, zy, zz)
		xx *= lx
		xy *= lx
		xz *= lx
		yx *= ly
		yy *= ly
		yz *= ly
		zx *= lz
		zy *= lz
		zz *= lz
	}
	// the trace is the sum of the diagonal elements; see
	// http://mathworld.wolfram.com/MatrixTrace.html
	t := xx + yy + zz

	// we protect the division by s by ensuring that s>=1
	if t >= 0 { // |w| >= .5
		s := float32(math.Sqrt(float64(t + 1))) // |s|>=1 ...
		self.w = 0.5 * s
		s = 0.5 / s // so this division isn't bad
		self.x = (zy - yz) * s
		self.y = (xz - zx) * s
		self.z = (yx - xy) * s
	} else if (xx > yy) && (xx > zz) {
		s := float32(math.Sqrt(float64(1.0 + xx - yy - zz))) // |s|>=1
		self.x = s * 0.5                                     // |x| >= .5
		s = 0.5 / s
		self.y = (yx + xy) * s
		self.z = (xz + zx) * s
		self.w = (zy - yz) * s
	} else if yy > zz {
		s := float32(math.Sqrt(float64(1.0 + yy - xx - zz))) // |s|>=1
		self.y = s * 0.5                                     // |y| >= .5
		s = 0.5 / s
		self.x = (yx + xy) * s
		self.z = (zy + yz) * s
		self.w = (xz - zx) * s
	} else {
		s := float32(math.Sqrt(float64(1.0 + zz - xx - yy))) // |s|>=1
		self.z = s * 0.5                                     // |z| >= .5
		s = 0.5 / s
		self.x = (xz + zx) * s
		self.y = (zy + yz) * s
		self.w = (yx - xy) * s
	}

	return self
}

// Set this quaternion to the rotation between two vectors.
// param v1 The base vector, which should be normalized.
// param v2 The target vector, which should be normalized.
func (self *Quaternion) SetFromCrossV3(v1, v2 *Vector3) *Quaternion {
	dot := utils.ClampFloat32(v1.DotV(v2), -1, 1)
	angle := float32(math.Acos(float64(dot)))
	return self.SetFromAxisRad(v1.Y*v2.Z-v1.Z*v2.Y, v1.Z*v2.X-v1.X*v2.Z, v1.X*v2.Y-v1.Y*v2.X, angle)
}

// Set this quaternion to the rotation between two vectors.
// param x1 The base vectors x value, which should be normalized.
// param y1 The base vectors y value, which should be normalized.
// param z1 The base vectors z value, which should be normalized.
// param x2 The target vector x value, which should be normalized.
// param y2 The target vector y value, which should be normalized.
// param z2 The target vector z value, which should be normalized.
func (self *Quaternion) SetFromCross(x1, y1, z1, x2, y2, z2 float32) *Quaternion {
	dot := utils.ClampFloat32(DotV3(x1, y1, z1, x2, y2, z2), -1, 1)
	angle := float32(math.Acos(float64(dot)))
	return self.SetFromAxisRad(y1*z2-z1*y2, z1*x2-x1*z2, x1*y2-y1*x2, angle)
}

// Spherical linear interpolation between this quaternion and the other quaternion, based on the alpha value in the range
// [0,1]. Taken from. Taken from Bones framework for JPCT, see http://www.aptalkarga.com/bones/
// param end the end quaternion
// param alpha alpha in the range [0,1
func (self *Quaternion) Slerp(end *Quaternion, alpha float32) *Quaternion {
	d := self.x*end.x + self.y*end.y + self.z*end.z + self.w*end.w
	var absDot float32
	if d < 0 {
		absDot = -d
	} else {
		absDot = d
	}
	// Set the first and second scale for the interpolation
	scale0 := 1 - alpha
	scale1 := alpha

	// Check if the angle between the 2 quaternions was big enough to
	// warrant such calculations
	if (1 - absDot) > 0.1 { // Get the angle between the 2 quaternions,
		// and then store the sin() of that angle
		angle := float32(math.Acos(float64(absDot)))
		invSinTheta := 1 / float32(math.Sin(float64(angle)))

		// Calculate the scale for q1 and q2, according to the angle and
		// it's sine value
		scale0 = float32(math.Sin(float64(((1 - alpha) * angle) * invSinTheta)))
		scale1 = float32(math.Sin(float64((alpha * angle) * invSinTheta)))
	}

	if d < 0 {
		scale1 = -scale1
	}

	// Calculate the x, y, z and w values for the quaternion by using a
	// special form of linear interpolation for quaternions.
	self.x = (scale0 * self.x) + (scale1 * end.x)
	self.y = (scale0 * self.y) + (scale1 * end.y)
	self.z = (scale0 * self.z) + (scale1 * end.z)
	self.w = (scale0 * self.w) + (scale1 * end.w)

	// Return the interpolated quaternion
	return self
}

// Spherical linearly interpolates multiple quaternions and stores the result in this Quaternion.
// Will not destroy the data previously inside the elements of q.
// result = (q_1^w_1)*(q_2^w_2)* ... *(q_n^w_n) where w_i=1/n.
// param q List of quaternions
func (self *Quaternion) SlerpQ(q []*Quaternion) *Quaternion {

	//Calculate exponents and multiply everything from left to right
	w := float32(1.0 / len(q))
	self.SetQ(q[0]).Exp(w)
	for i := 1; i < len(q); i++ {
		self.MulQ(tmp1.SetQ(q[i]).Exp(w))
	}
	self.Nor()
	return self
}

// Spherical linearly interpolates multiple quaternions by the given weights and stores the result in this Quaternion.
// Will not destroy the data previously inside the elements of q or w.
// result = (q_1^w_1)*(q_2^w_2)* ... *(q_n^w_n) where the sum of w_i is 1.
// Lists must be equal in length.
// param q List of quaternions
// param w List of weights
func (self *Quaternion) SlerpW(q []*Quaternion, w []float32) *Quaternion {

	//Calculate exponents and multiply everything from left to right
	self.SetQ(q[0]).Exp(w[0])
	for i := 1; i < len(q); i++ {
		self.MulQ(tmp1.SetQ(q[i]).Exp(w[i]))
	}
	self.Nor()
	return self
}

// Calculates (this quaternion)^alpha where alpha is a real number and stores the result in this quaternion.
// See http://en.wikipedia.org/wiki/Quaternion#Exponential.2C_logarithm.2C_and_power
// param alpha Exponent
func (self *Quaternion) Exp(alpha float32) *Quaternion {

	//Calculate |q|^alpha
	norm := self.Len()
	normExp := float32(math.Pow(float64(norm), float64(alpha)))

	//Calculate theta
	theta := float32(math.Acos(float64(self.w / norm)))

	//Calculate coefficient of basis elements
	//If theta is small enough, use the limit of sin(alpha*theta) / sin(theta) instead of actual value
	var coeff float32
	if math.Abs(float64(theta)) < 0.001 {
		coeff = normExp * alpha / norm
	} else {
		coeff = float32(float64(normExp) * math.Sin(float64(alpha*theta)/(float64(norm)*math.Sin(float64(theta)))))
	}

	//Write results
	self.w = float32((float64(normExp) * math.Cos(float64(alpha*theta))))
	self.x *= coeff
	self.y *= coeff
	self.z *= coeff

	//Fix any possible discrepancies
	self.Nor()

	return self
}

// func (self *Quaternion) HashCode() int {
// 	int prime = 31;
// 	int result = 1;
// 	result = prime * result + NumberUtils.floatToRawIntBits(w);
// 	result = prime * result + NumberUtils.floatToRawIntBits(x);
// 	result = prime * result + NumberUtils.floatToRawIntBits(y);
// 	result = prime * result + NumberUtils.floatToRawIntBits(z);
// 	return result;
// }

// func (self *Quaternion) Equals(Object obj) {
// 	if (this == obj) {
// 		return true;
// 	}
// 	if (obj == null) {
// 		return false;
// 	}
// 	if (!(obj instanceof Quaternion)) {
// 		return false;
// 	}
// 	other *Quaternion = (Quaternion)obj;
// 	return (NumberUtils.floatToRawIntBits(w) == NumberUtils.floatToRawIntBits(other.w))
// 		&& (NumberUtils.floatToRawIntBits(x) == NumberUtils.floatToRawIntBits(other.x))
// 		&& (NumberUtils.floatToRawIntBits(y) == NumberUtils.floatToRawIntBits(other.y))
// 		&& (NumberUtils.floatToRawIntBits(z) == NumberUtils.floatToRawIntBits(other.z));
// }

// Get the dot product between the two quaternions (commutative).
// param x1 the x component of the first quaternion
// param y1 the y component of the first quaternion
// param z1 the z component of the first quaternion
// param w1 the w component of the first quaternion
// param x2 the x component of the second quaternion
// param y2 the y component of the second quaternion
// param z2 the z component of the second quaternion
// param w2 the w component of the second quaternion
// return the dot product between the first and second quaternion.
func Dot(x1, y1, z1, w1, x2, y2, z2, w2 float32) float32 {
	return x1*x2 + y1*y2 + z1*z2 + w1*w2
}

// Get the dot product between this and the other quaternion (commutative).
// param other the other quaternion.
// return the dot product of this and the other quaternion.
func (self *Quaternion) DotQ(other *Quaternion) float32 {
	return self.x*other.x + self.y*other.y + self.z*other.z + self.w*other.w
}

// Get the dot product between this and the other quaternion (commutative).
// param x the x component of the other quaternion
// param y the y component of the other quaternion
// param z the z component of the other quaternion
// param w the w component of the other quaternion
//  return the dot product of this and the other quaternion.
func (self *Quaternion) Dot(x, y, z, w float32) float32 {
	return self.x*x + self.y*y + self.z*z + self.w*w
}

// Multiplies the components of this quaternion with the given scalar.
// param scalar the scalar.
func (self *Quaternion) MulScalar(scalar float32) *Quaternion {
	self.x *= scalar
	self.y *= scalar
	self.z *= scalar
	self.w *= scalar
	return self
}

// Get the axis angle representation of the rotation in degrees. The supplied vector will receive the axis (x, y and z values)
// of the rotation and the value returned is the angle in degrees around that axis. Note that this method will alter the
// supplied vector, the existing value of the vector is ignored. </p> This will normalize this quaternion if needed. The
// received axis is a unit vector. However, if this is an identity quaternion (no rotation), then the length of the axis may be
// zero.
//
// param axis vector which will receive the axis
// return the angle in degrees
// @see <a href="http://en.wikipedia.org/wiki/Axis%E2%80%93angle_representation">wikipedia</a>
// @see <a href="http://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToAngle">calculation</a>
func (self *Quaternion) GetAxisAngle(axis *Vector3) float32 {
	return self.GetAxisAngleRad(axis) * utils.RadiansToDegrees
}

// Get the axis-angle representation of the rotation in radians. The supplied vector will receive the axis (x, y and z values)
// of the rotation and the value returned is the angle in radians around that axis. Note that this method will alter the
// supplied vector, the existing value of the vector is ignored. </p> This will normalize this quaternion if needed. The
// received axis is a unit vector. However, if this is an identity quaternion (no rotation), then the length of the axis may be
// zero.
//
// param axis vector which will receive the axis
// return the angle in radians
// @see <a href="http://en.wikipedia.org/wiki/Axis%E2%80%93angle_representation">wikipedia</a>
// @see <a href="http://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToAngle">calculation</a>
func (self *Quaternion) GetAxisAngleRad(axis *Vector3) float32 {
	// if w>1 acos and sqrt will produce errors, this cant happen if quaternion is normalised
	if self.w > 1 {
		self.Nor()
	}
	angle := float32(2.0 * math.Acos(float64(self.w)))
	s := float32(math.Sqrt(float64(1 - self.w*self.w))) // assuming quaternion normalised then w is less than 1, so term always positive.
	if s < utils.FLOAT_ROUNDING_ERROR {                 // test to avoid divide by zero, s is always positive due to sqrt
		// if s close to zero then direction of axis not important
		axis.X = self.x // if it is important that axis is normalised then replace with x=1; y=z=0;
		axis.Y = self.y
		axis.Z = self.z
	} else {
		axis.X = float32(self.x / s) // normalise axis
		axis.Y = float32(self.y / s)
		axis.Z = float32(self.z / s)
	}

	return angle
}

// Get the angle in radians of the rotation this quaternion represents. Does not normalize the quaternion. Use
// {@link #getAxisAngleRad(Vector3)} to get both the axis and the angle of this rotation. Use
// {@link #getAngleAroundRad(Vector3)} to get the angle around a specific axis.
// return the angle in radians of the rotation
func (self *Quaternion) GetAngleRad() float32 {
	if self.w > 1 {
		return float32(2.0 * math.Acos(float64(self.w/self.Len())))
	} else {
		return float32(2.0 * math.Acos(float64(self.w)))
	}
}

// Get the angle in degrees of the rotation this quaternion represents. Use {@link #getAxisAngle(Vector3)} to get both the axis
// and the angle of this rotation. Use {@link #getAngleAround(Vector3)} to get the angle around a specific axis.
// return the angle in degrees of the rotation
func (self *Quaternion) GetAngle() float32 {
	return self.GetAngleRad() * utils.RadiansToDegrees
}

// Get the swing rotation and twist rotation for the specified axis. The twist rotation represents the rotation around the
// specified axis. The swing rotation represents the rotation of the specified axis itself, which is the rotation around an
// axis perpendicular to the specified axis.
//
// The swing and twist rotation can be used to reconstruct the original quaternion: this = swing * twist
//
// param axisX the X component of the normalized axis for which to get the swing and twist rotation
// param axisY the Y component of the normalized axis for which to get the swing and twist rotation
// param axisZ the Z component of the normalized axis for which to get the swing and twist rotation
// param swing will receive the swing rotation: the rotation around an axis perpendicular to the specified axis
// param twist will receive the twist rotation: the rotation around the specified axis
// @see <a href="http://www.euclideanspace.com/maths/geometry/rotations/for/decomposition">calculation</a>
func (self *Quaternion) GetSwingTwist(axisX, axisY, axisZ float32, swing, twist *Quaternion) {
	d := DotV3(self.x, self.y, self.z, axisX, axisY, axisZ)
	twist.Set(axisX*d, axisY*d, axisZ*d, self.w).Nor()
	swing.SetQ(twist).Conjugate().MulLeftQ(self)
}

// Get the swing rotation and twist rotation for the specified axis. The twist rotation represents the rotation around the
// specified axis. The swing rotation represents the rotation of the specified axis itself, which is the rotation around an
// axis perpendicular to the specified axis.
//
// The swing and twist rotation can be used to reconstruct the original quaternion: this = swing * twist
//
// param axis the normalized axis for which to get the swing and twist rotation
// param swing will receive the swing rotation: the rotation around an axis perpendicular to the specified axis
// param twist will receive the twist rotation: the rotation around the specified axis
// @see <a href="http://www.euclideanspace.com/maths/geometry/rotations/for/decomposition">calculation</a>
func (self *Quaternion) GetSwingTwistV(axis *Vector3, swing, twist *Quaternion) {
	self.GetSwingTwist(axis.X, axis.Y, axis.Z, swing, twist)
}

// Get the angle in radians of the rotation around the specified axis. The axis must be normalized.
// param axisX the x component of the normalized axis for which to get the angle
// param axisY the y component of the normalized axis for which to get the angle
// param axisZ the z component of the normalized axis for which to get the angle
// return the angle in radians of the rotation around the specified axis
func (self *Quaternion) GetAngleAroundRad(axisX, axisY, axisZ float32) float32 {
	d := DotV3(self.x, self.y, self.z, axisX, axisY, axisZ)
	l2 := Len2Q(axisX*d, axisY*d, axisZ*d, self.w)
	if utils.IsZero(l2) {
		return 0
	} else {
		return float32((2.0 * math.Acos(float64(utils.ClampFloat32((self.w/float32(math.Sqrt(float64(l2)))), -1, 1)))))
	}
}

// Get the angle in radians of the rotation around the specified axis. The axis must be normalized.
// param axis the normalized axis for which to get the angle
// return the angle in radians of the rotation around the specified axis
func (self *Quaternion) GetAngleAroundRadV(axis *Vector3) float32 {
	return self.GetAngleAroundRad(axis.X, axis.Y, axis.Z)
}

// Get the angle in degrees of the rotation around the specified axis. The axis must be normalized.
// param axisX the x component of the normalized axis for which to get the angle
// param axisY the y component of the normalized axis for which to get the angle
// param axisZ the z component of the normalized axis for which to get the angle
// return the angle in degrees of the rotation around the specified axis
func (self *Quaternion) GetAngleAround(axisX, axisY, axisZ float32) float32 {
	return self.GetAngleAroundRad(axisX, axisY, axisZ) * utils.RadiansToDegrees
}

// Get the angle in degrees of the rotation around the specified axis. The axis must be normalized.
// param axis the normalized axis for which to get the angle
// return the angle in degrees of the rotation around the specified axis
func (self *Quaternion) GetAngleAroundV(axis *Vector3) float32 {
	return self.GetAngleAround(axis.X, axis.Y, axis.Z)
}

func (self *Quaternion) String() string {
	return ""
	// return "[" + self.x + "|" + self.y + "|" + self.z + "|" + self.w + "]"
}
