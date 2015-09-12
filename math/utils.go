// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Utility and fast math functions.
package math

import (
	"math"
	"math/rand"
)

var (
	NanoToSec float32 = 1 / 1000000000

	FLOAT_ROUNDING_ERROR float32 = 0.000001 // 32 bits
	PI                   float32 = 3.1415927
	PI2                  float32 = PI * 2
	E                    float32 = 2.7182818

	SIN_BITS  uint = 14 // 16KB. Adjust for accuracy.
	SIN_MASK  int  = ^(-1 << SIN_BITS)
	SIN_COUNT int  = SIN_MASK + 1

	RadFull    float32 = PI * 2
	DegFull    float32 = 360
	RadToIndex float32 = float32(SIN_COUNT) / RadFull
	DegToIndex float32 = float32(SIN_COUNT) / DegFull

	// multiply by this to convert from radians to degrees
	RadiansToDegrees float32 = 180 / PI
	radDeg           float32 = RadiansToDegrees

	// multiply by this to convert from degrees to radians
	DegreesToRadians float32 = PI / 180
	DegRad           float32 = DegreesToRadians

	BIG_ENOUGH_INT   float32 = 16 * 1024
	BIG_ENOUGH_FLOOR float32 = BIG_ENOUGH_INT
	CEIL             float32 = 0.9999999
	BIG_ENOUGH_CEIL  float32 = 16384.999999999996
	BIG_ENOUGH_ROUND float32 = BIG_ENOUGH_INT + 0.5
)

// var SinTable [SIN_COUNT]float32

// func init() {
// 	for i := 0; i < SIN_COUNT; i++ {
// 		SinTable[i] = float32(math.Sin(float64((i + 0.5) / SIN_COUNT * RadFull)))
// 	}
// 	for i := 0; i < 360; i += 90 {
// 		SinTable[int(int(i*DegToIndex)&SIN_MASK)] = float32(math.Sin(float64(i * DegreesToRadians)))
// 	}
// }

// // Returns the sine in radians from a lookup table.
// func Sin(radians float32) float32 {
// 	return SinTable[int((radians*RadToIndex)&SIN_MASK)]
// }

// // Returns the cosine in radians from a lookup table.
// func Cos(radians float32) float32 {
// 	return SinTable[int(((radians+PI/2)*RadToIndex)&SIN_MASK)]
// }

// // Returns the sine in radians from a lookup table.
// func SinDeg(degrees float32) float32 {
// 	return SinTable[int((degrees*DegToIndex)&SIN_MASK)]
// }

// // Returns the cosine in radians from a lookup table.
// func CosDeg(degrees float32) float32 {
// 	return SinTable[int(((degrees+90)*DegToIndex)&SIN_MASK)]
// }

/** Returns atan2 in radians, faster but less accurate than math.atan2. Average error of 0.00231 radians (0.1323 degrees),
 * largest error of 0.00488 radians (0.2796 degrees). */
func Atan2(y, x float32) float32 {
	if x == 0 {
		if y > 0 {
			return PI / 2
		}
		if y == 0 {
			return 0
		}
		return -PI / 2
	}
	atan := y / x
	z := atan
	if math.Abs(float64(z)) < 1 {
		atan = z / (1 + 0.28*z*z)
		if x < 0 {
			if y < 0 {
				return atan - PI
			}
			return atan + PI
		}
		return atan
	}
	atan = PI/2 - z/(z*z+0.28)
	if y < 0 {
		return atan - PI
	}
	return atan
}

// Returns the next power of two. Returns the specified value if the value is already a power of two.
func NextPowerOfTwo(value int) int {
	if value == 0 {
		return 1
	}
	value--
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	return value + 1
}

func IsPowerOfTwo(value int) bool {
	return value != 0 && (value&value-1) == 0
}

func ClampInt8(value int8, min int8, max int8) int8 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ClampInt(value int, min int, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ClampInt64(value int64, min int64, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ClampFloat32(value float32, min float32, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ClampFloat64(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Linearly interpolates between fromValue to toValue on progress position.
func Lerp(fromValue, toValue, progress float32) float32 {
	return fromValue + (toValue-fromValue)*progress
}

/** Returns the largest integer less than or equal to the specified float. This method will only properly floor floats from
 * -(2^14) to (Float.MAX_VALUE - 2^14). */
func Floor(value float32) int {
	return int((value + BIG_ENOUGH_FLOOR) - BIG_ENOUGH_INT)
}

/** Returns the largest integer less than or equal to the specified float. This method will only properly floor floats that are
 * positive. Note this method simply casts the float32 to int. */
func FloorPositive(value float32) int {
	return int(value)
}

/** Returns the smallest integer greater than or equal to the specified float. This method will only properly ceil floats from
 * -(2^14) to (Float.MAX_VALUE - 2^14). */
func Ceil(value float32) int {
	return int((value + BIG_ENOUGH_CEIL) - BIG_ENOUGH_INT)
}

/** Returns the smallest integer greater than or equal to the specified float. This method will only properly ceil floats that
 * are positive. */
func CeilPositive(value float32) int {
	return int((value + CEIL))
}

/** Returns the closest integer to the specified float. This method will only properly round floats from -(2^14) to
 * (Float.MAX_VALUE - 2^14). */
func Round(value float32) int {
	return int((value + BIG_ENOUGH_ROUND) - BIG_ENOUGH_INT)
}

// Returns the closest integer to the specified float. This method will only properly round floats that are positive.
func RoundPositive(value float32) int {
	return int(value + 0.5)
}

// Returns true if the value is zero (using the default tolerance as upper bound)
func IsZero(value float32) bool {
	return float32(math.Abs(float64(value))) <= FLOAT_ROUNDING_ERROR
}

/** Returns true if the value is zero.
 * @param tolerance represent an upper bound below which the value is considered zero. */
func isZeroTolerance(value, tolerance float32) bool {
	return math.Abs(float64(value)) <= float64(tolerance)
}

// Returns true if a is nearly equal to b. The function uses the default floating error tolerance.
func IsEqual(a, b float32) bool {
	return float32(math.Abs(float64(a-b))) <= FLOAT_ROUNDING_ERROR
}

/** Returns true if a is nearly equal to b.
 * @param a the first value.
 * @param b the second value.
 * @param tolerance represent an upper bound below which the two values are considered equal. */
func IsEqualTolerance(a, b, tolerance float32) bool {
	return math.Abs(float64(a-b)) <= float64(tolerance)
}

// @return the logarithm of value with base a
func Log(a, value float32) float32 {
	return float32((math.Log(float64(value)) / math.Log(float64(a))))
}

// @return the logarithm of value with base 2
func Log2(value float32) float32 {
	return float32(math.Log2(float64(value)))
}

// TODO int64 and float

/** Returns a random number between 0 (inclusive) and the specified value (inclusive). */
func Random(n int) int {
	return rand.Intn(n)
}

// Returns a random number between start (inclusive) and end (inclusive).
func RandomRange(start, end int) int {
	return rand.Intn(end - start + 1)
}

func RandomFloat() float32 {
	return rand.Float32()
}

//   /** Returns a random boolean value. */
//   func boolean randomBoolean () {
//     return random.nextBoolean();
//   }

//   /** Returns true if a random value between 0 and 1 is less than the specified value. */
//   func boolean randomBoolean (float32 chance) {
//     return MathUtils.random() < chance;
//   }

//   /** Returns -1 or 1, randomly. */
//   func int randomSign () {
//     return 1 | (random.nextInt() >> 31);
//   }

/** Returns a triangularly distributed random number between -1.0 (exclusive) and 1.0 (exclusive), where values around zero are
 * more likely.
 * <p>
 * This is an optimized version of {@link #randomTriangular(float, float, float) randomTriangular(-1, 1, 0)} */
func RandomTriangular() float32 {
	return rand.Float32() - rand.Float32()
}

/* Returns a triangularly distributed random number between {@code -max} (exclusive) and {@code max} (exclusive), where values
 * around zero are more likely.
 * <p>
 * This is an optimized version of {@link #randomTriangular(float, float, float) randomTriangular(-max, max, 0)}
 * @param max the upper limit */
func RandomTriangularMax(max float32) float32 {
	return (rand.Float32() - rand.Float32()) * max
}

/** Returns a triangularly distributed random number between {@code min} (inclusive) and {@code max} (exclusive), where the
 * {@code mode} argument defaults to the midpoint between the bounds, giving a symmetric distribution.
 * <p>
 * This method is equivalent of {@link #randomTriangular(float, float, float) randomTriangular(min, max, (max - min) * .5f)}
 * @param min the lower limit
 * @param max the upper limit */
func RandomTriangularMinMax(min, max float32) float32 {
	return RandomTriangularMinMaxMode(min, max, min+(max-min)*0.5)
}

/** Returns a triangularly distributed random number between {@code min} (inclusive) and {@code max} (exclusive), where values
 * around {@code mode} are more likely.
 * @param min the lower limit
 * @param max the upper limit
 * @param mode the point around which the values are more likely */
func RandomTriangularMinMaxMode(min, max, mode float32) float32 {
	u := rand.Float32()
	d := max - min
	if u <= (mode-min)/d {
		return min + float32(math.Sqrt(float64(u*d*(mode-min))))
	}
	return max - float32(math.Sqrt(float64((1-u)*d*(max-mode))))
}
