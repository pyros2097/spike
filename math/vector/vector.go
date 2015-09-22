// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package vector

import (
	. "github.com/pyros2097/spike/math/interpolation"
)

/** Encapsulates a general vector. Allows chaining operations by returning a reference to itself in all modification methods. See
 * {@link Vector2} and {@link Vector3} for specific implementations.
 * @author Xoppa */
type IVector interface {
	/** @return a copy of this vector */
	Copy() interface{}

	/** @return The euclidean length */
	Len() float32

	/** This method is faster than {@link Vector#len()} because it avoids calculating a square root. It is useful for comparisons,
	 * but not for getting exact lengths, as the return value is the square of the actual length.
	 * @return The squared euclidean length */
	Len2() float32

	/** Limits the length of this vector, based on the desired maximum length.
	 * @param limit desired maximum length for this vector
	 * @return this vector for chaining */
	Limit(limit float32) interface{}

	/** Limits the length of this vector, based on the desired maximum length squared.
	 * <p />
	 * This method is slightly faster than limit().
	 * @param limit2 squared desired maximum length for this vector
	 * @return this vector for chaining
	 * @see #len2() */
	Limit2(limit2 float32) interface{}

	/** Sets the length of this vector. Does nothing is this vector is zero.
	 * @param len desired length for this vector
	 * @return this vector for chaining */
	SetLength(length float32) interface{}

	/** Sets the length of this vector, based on the square of the desired length. Does nothing is this vector is zero.
	 * <p />
	 * This method is slightly faster than setLength().
	 * @param len2 desired square of the length for this vector
	 * @return this vector for chaining
	 * @see #len2() */
	SetLength2(length2 float32) interface{}

	/** Clamps this vector's length to given min and max values
	 * @param min Min length
	 * @param max Max length
	 * @return This vector for chaining */
	Clamp(min, max float32) interface{}

	/** Sets this vector from the given vector
	 * @param v The vector
	 * @return This vector for chaining */
	SetV(v interface{}) interface{}

	/** Subtracts the given vector from this vector.
	 * @param v The vector
	 * @return This vector for chaining */
	SubV(v interface{}) interface{}

	/** Normalizes this vector. Does nothing if it is zero.
	 * @return This vector for chaining */
	Nor() interface{}

	/** Adds the given vector to this vector
	 * @param v The vector
	 * @return This vector for chaining */
	AddV(v interface{}) interface{}

	/** @param v The other vector
	 * @return The dot product between this and the other vector */
	DotV(v interface{}) float32

	/** Scales this vector by a scalar
	 * @param scalar The scalar
	 * @return This vector for chaining */
	SclScalar(scalar float32) interface{}

	/** Scales this vector by another vector
	 * @return This vector for chaining */
	SclV(v interface{}) interface{}

	/** @param v The other vector
	 * @return the distance between this and the other vector */
	DstV(v interface{}) float32

	/** This method is faster than {@link Vector#dst(Vector)} because it avoids calculating a square root. It is useful for
	 * comparisons, but not for getting accurate distances, as the return value is the square of the actual distance.
	 * @param v The other vector
	 * @return the squared distance between this and the other vector */
	Dst2V(v interface{}) float32

	/** Linearly interpolates between this vector and the target vector by alpha which is in the range [0,1]. The result is stored
	 * in this vector.
	 * @param target The target vector
	 * @param alpha The interpolation coefficient
	 * @return This vector for chaining. */
	Lerp(target interface{}, alpha float32) interface{}

	/** Interpolates between this vector and the given target vector by alpha (within range [0,1]) using the given Interpolation
	 * method. the result is stored in this vector.
	 * @param target The target vector
	 * @param alpha The interpolation coefficient
	 * @param interpolator An Interpolation object describing the used interpolation method
	 * @return This vector for chaining. */
	Interpolate(target interface{}, alpha float32, interpolation Interpolation) interface{}

	/** @return Whether this vector is a unit length vector */
	IsUnit() bool

	/** @return Whether this vector is a unit length vector within the given margin. */
	IsUnitMargin(margin float32) bool

	/** @return Whether this vector is a zero vector */
	IsZero() bool

	/** @return Whether the length of this vector is smaller than the given margin */
	IsZeroMargin(margin float32) bool

	/** @return true if this vector is in line with the other vector (either in the same or the opposite direction) */
	IsOnLine(other interface{}) bool

	/** @return true if this vector is in line with the other vector (either in the same or the opposite direction) */
	IsOnLineEpsilon(other interface{}, epsilon float32)

	/** @return true if this vector is collinear with the other vector ({@link #isOnLine(Vector, float)} &&
	 *         {@link #hasSameDirection(Vector)}). */
	IsCollinearEpsilon(other interface{}, epsilon float32) bool

	/** @return true if this vector is collinear with the other vector ({@link #isOnLine(Vector)} &&
	 *         {@link #hasSameDirection(Vector)}). */
	IsCollinear(other interface{}) bool

	/** @return true if this vector is opposite collinear with the other vector ({@link #isOnLine(Vector, float)} &&
	 *         {@link #hasOppositeDirection(Vector)}). */
	IsCollinearOppositeEpsilon(other interface{}, epsilon float32) bool

	/** @return true if this vector is opposite collinear with the other vector ({@link #isOnLine(Vector)} &&
	 *         {@link #hasOppositeDirection(Vector)}). */
	IsCollinearOpposite(other interface{}) bool

	/** @return Whether this vector is perpendicular with the other vector. True if the dot product is 0. */
	IsPerpendicular(v interface{}) bool

	/** @return Whether this vector is perpendicular with the other vector. True if the dot product is 0.
	 * @param epsilon a positive small number close to zero */
	IsPerpendicularEpsilon(vector interface{}, epsilon float32) bool

	/** @return Whether this vector has similar direction compared to the other vector. True if the normalized dot product is > 0. */
	HasSameDirection(v interface{}) bool

	/** @return Whether this vector has opposite direction compared to the other vector. True if the normalized dot product is < 0. */
	HasOppositeDirection(v interface{}) bool

	/** Compares this vector with the other vector, using the supplied epsilon for fuzzy equality testing.
	 * @param other
	 * @param epsilon
	 * @return whether the vectors have fuzzy equality. */
	EpsilonEqualsV(other interface{}, epsilon float32) bool

	/** First scale a supplied vector, then add it to this vector.
	 * @param v addition vector
	 * @param scalar for scaling the addition vector */
	MulAddScalar(v interface{}, scalar float32) interface{}

	/** First scale a supplied vector, then add it to this vector.
	 * @param v addition vector
	 * @param mulVec vector by whose values the addition vector will be scaled */
	MulAdd(v interface{}, mulVec interface{}) interface{}

	/** Sets the components of this vector to 0
	 * @return This vector for chaining */
	SetZero() interface{}
}
