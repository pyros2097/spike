// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package geom

import (
	"math"

	. "github.com/pyros2097/spike/math/utils"
	. "github.com/pyros2097/spike/math/vector"
)

var (
	tmp1 = NewVector2Empty()
	tmp2 = NewVector2Empty()
	tmp3 = NewVector2Empty()
)

// Computes the barycentric coordinates v,w for the specified point in the triangle.
// If barycentric.X >= 0 && barycentric.Y >= 0 && barycentric.X + barycentric.Y <= 1 then the point is inside the triangle.
// If vertices a,b,c have values aa,bb,cc then to get an interpolated value at point p:
//
// Barycentric(p, a, b, c, barycentric)
// u := 1 - barycentric.X - barycentric.Y
// x := u * aa.x + barycentric.X * bb.x + barycentric.Y * cc.x
// y := u * aa.y + barycentric.X * bb.y + barycentric.Y * cc.y
func ToBarycoord(p, a, b, c, barycentricOut *Vector2) *Vector2 {
	v0 := tmp1.SetV(b).SubV(a)
	v1 := tmp2.SetV(c).SubV(a)
	v2 := tmp3.SetV(p).SubV(a)
	d00 := v0.DotV(v0)
	d01 := v0.DotV(v1)
	d11 := v1.DotV(v1)
	d20 := v2.DotV(v0)
	d21 := v2.DotV(v1)
	denom := d00*d11 - d01*d01
	barycentricOut.X = (d11*d20 - d01*d21) / denom
	barycentricOut.Y = (d00*d21 - d01*d20) / denom
	return barycentricOut
}

// Returns true if the barycentric coordinates are inside the triangle.
func BarycoordInsideTriangle(barycentric *Vector2) bool {
	return barycentric.X >= 0 && barycentric.Y >= 0 && barycentric.X+barycentric.Y <= 1
}

// Returns interpolated values given the barycentric coordinates of a point in a triangle and the values at each vertex.
// @return interpolatedOut
func FromBarycoord(barycentric, a, b, c, interpolatedOut *Vector2) *Vector2 {
	u := 1 - barycentric.X - barycentric.Y
	interpolatedOut.X = u*a.X + barycentric.X*b.X + barycentric.Y*c.X
	interpolatedOut.Y = u*a.Y + barycentric.X*b.Y + barycentric.Y*c.Y
	return interpolatedOut
}

// Returns an interpolated value given the barycentric coordinates of a point in a triangle and the values at each vertex.
// @return interpolatedOut
func FromBarycoordF(barycentric *Vector2, a, b, c float32) float32 {
	u := 1 - barycentric.X - barycentric.Y
	return u*a + barycentric.X*b + barycentric.Y*c
}

/** Returns the lowest positive root of the quadric equation given by a* x * x + b * x + c = 0. If no solution is given
 * Float.Nan is returned.
 * @param a the first coefficient of the quadric equation
 * @param b the second coefficient of the quadric equation
 * @param c the third coefficient of the quadric equation
 * @return the lowest positive root or Float.Nan */
func LowestPositiveRoot(a, b, c float32) float32 {
	det := b*b - 4*a*c
	if det < 0 {
		return float32(math.NaN())
	}

	sqrtD := float32(math.Sqrt(float64(det)))
	invA := 1 / (2 * a)
	r1 := (-b - sqrtD) * invA
	r2 := (-b + sqrtD) * invA

	if r1 > r2 {
		tmp := r2
		r2 = r1
		r1 = tmp
	}

	if r1 > 0 {
		return r1
	}
	if r2 > 0 {
		return r2
	}
	return float32(math.NaN())
}

func Colinear(x1, y1, x2, y2, x3, y3 float32) bool {
	dx21 := x2 - x1
	dy21 := y2 - y1
	dx32 := x3 - x2
	dy32 := y3 - y2
	det := dx32*dy21 - dx21*dy32
	return float32(math.Abs(float64(det))) < FLOAT_ROUNDING_ERROR
}

func TriangleCentroid(x1, y1, x2, y2, x3, y3 float32, centroid *Vector2) *Vector2 {
	centroid.X = (x1 + x2 + x3) / 3
	centroid.Y = (y1 + y2 + y3) / 3
	return centroid
}

// Returns the circumcenter of the triangle. The input points must not be colinear
func TriangleCircumcenter(x1, y1, x2, y2, x3, y3 float32, circumcenter *Vector2) *Vector2 {
	dx21 := x2 - x1
	dy21 := y2 - y1
	dx32 := x3 - x2
	dy32 := y3 - y2
	dx13 := x1 - x3
	dy13 := y1 - y3
	det := dx32*dy21 - dx21*dy32
	if float32(math.Abs(float64(det))) < FLOAT_ROUNDING_ERROR {
		panic("Triangle points must not be colinear.")
	}
	det *= 2
	sqr1 := x1*x1 + y1*y1
	sqr2 := x2*x2 + y2*y2
	sqr3 := x3*x3 + y3*y3
	circumcenter.Set((sqr1*dy32+sqr2*dy13+sqr3*dy21)/det, -(sqr1*dx32+sqr2*dx13+sqr3*dx21)/det)
	return circumcenter
}

func TriangleArea(x1, y1, x2, y2, x3, y3 float32) float32 {
	return float32(math.Abs(float64((x1-x3)*(y2-y1)-(x1-x2)*(y3-y1))) * 0.5)
}

func QuadrilateralCentroid(x1, y1, x2, y2, x3, y3, x4, y4 float32, centroid *Vector2) *Vector2 {
	avgX1 := (x1 + x2 + x3) / 3
	avgY1 := (y1 + y2 + y3) / 3
	avgX2 := (x1 + x4 + x3) / 3
	avgY2 := (y1 + y4 + y3) / 3
	centroid.X = avgX1 - (avgX1-avgX2)/2
	centroid.Y = avgY1 - (avgY1-avgY2)/2
	return centroid
}

// Returns the centroid for the specified non-self-intersecting polygon.
func PolygonCentroid(polygon []float32, offset, count int, centroid *Vector2) *Vector2 {
	if count < 6 {
		panic("A polygon must have 3 or more coordinate pairs.")
	}
	var x float32 = 0
	var y float32 = 0
	var signedArea float32 = 0
	i := offset
	for n := offset + count - 2; i < n; i += 2 {
		x0 := polygon[i]
		y0 := polygon[i+1]
		x1 := polygon[i+2]
		y1 := polygon[i+3]
		a := x0*y1 - x1*y0
		signedArea += a
		x += (x0 + x1) * a
		y += (y0 + y1) * a
	}

	x0 := polygon[i]
	y0 := polygon[i+1]
	x1 := polygon[offset]
	y1 := polygon[offset+1]
	a := x0*y1 - x1*y0
	signedArea += a
	x += (x0 + x1) * a
	y += (y0 + y1) * a

	if signedArea == 0 {
		centroid.X = 0
		centroid.Y = 0
	} else {
		signedArea *= 0.5
		centroid.X = x / (6 * signedArea)
		centroid.Y = y / (6 * signedArea)
	}
	return centroid
}

// Computes the area for a convex polygon.
func PolygonArea(polygon []float32, offset, count int) float32 {
	var area float32 = 0
	i := offset
	for n := offset + count; i < n; i += 2 {
		x1 := i
		y1 := i + 1
		x2 := (i + 2) % n
		if x2 < offset {
			x2 += offset
		}
		y2 := (i + 3) % n
		if y2 < offset {
			y2 += offset
		}
		area += polygon[x1] * polygon[y2]
		area -= polygon[x2] * polygon[y1]
	}
	area *= 0.5
	return area
}

func EnsureCCW(polygon []float32) {
	if !AreVerticesClockwise(polygon, 0, len(polygon)) {
		return
	}
	lastX := len(polygon) - 2
	i := 0
	for n := len(polygon) / 2; i < n; i += 2 {
		other := lastX - i
		x := polygon[i]
		y := polygon[i+1]
		polygon[i] = polygon[other]
		polygon[i+1] = polygon[other+1]
		polygon[other] = x
		polygon[other+1] = y
	}
}

func AreVerticesClockwise(polygon []float32, offset, count int) bool {
	if count <= 2 {
		return false
	}
	var area float32 = 0
	var p1x, p1y, p2x, p2y float32
	i := offset
	for n := offset + count - 3; i < n; i += 2 {
		p1x = polygon[i]
		p1y = polygon[i+1]
		p2x = polygon[i+2]
		p2y = polygon[i+3]
		area += p1x*p2y - p2x*p1y
	}
	p1x = polygon[count-2]
	p1y = polygon[count-1]
	p2x = polygon[0]
	p2y = polygon[1]
	return area+p1x*p2y-p2x*p1y < 0
}
