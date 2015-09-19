// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

// Interface that specifies a path of type T within the window 0.0<=t<=1.0.
type Path interface {
	DerivativeAt(out interface{}, t float32) interface{}

	// return The value of the path at t where 0<=t<=1
	ValueAt(out interface{}, t float32) interface{}

	// return The approximated value (between 0 and 1) on the path which is closest to the specified value. Note that the
	// implementation of this method might be optimized for speed against precision, see {@link #locate(Object)} for a more
	// precise (but more intensive) method.
	Approximate(v interface{}) float32

	// return The precise location (between 0 and 1) on the path which is closest to the specified value. Note that the
	// implementation of this method might be CPU intensive, see {@link #approximate(Object)} for a faster (but less
	// recise) method.
	Locate(v interface{}) float32

	// param samples The amount of divisions used to approximate length. Higher values will produce more precise results,
	// but will be more CPU intensive.
	// return An approximated length of the spline through sampling the curve and accumulating the euclidean distances between
	// the sample points.

	ApproxLength(samples int) float32
}
