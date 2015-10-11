// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package offering various methods for intersection testing between different geometric objects.
package intersector

import (
	"math"

	. "github.com/pyros2097/spike/math"
	. "github.com/pyros2097/spike/math/collision"
	. "github.com/pyros2097/spike/math/shape"
	"github.com/pyros2097/spike/math/utils"
	. "github.com/pyros2097/spike/math/vector"
)

var (
	v0 = NewVector3Empty()
	v1 = NewVector3Empty()
	v2 = NewVector3Empty()
)

// Returns whether the given point is inside the triangle. This assumes that the point is on the plane of the triangle. No
// check is performed that this is the case.
//
// param point the point
// param t1 the first vertex of the triangle
// param t2 the second vertex of the triangle
// param t3 the third vertex of the triangle
// return whether the point is in the triangle
func IsPointInTriangleV3(point, t1, t2, t3 *Vector3) bool {
	v0.SetV(t1).SubV(point)
	v1.SetV(t2).SubV(point)
	v2.SetV(t3).SubV(point)

	ab := v0.DotV(v1)
	ac := v0.DotV(v2)
	bc := v1.DotV(v2)
	cc := v2.DotV(v2)

	if bc*ac-cc*ab < 0 {
		return false
	}
	bb := v1.DotV(v1)
	if ab*bc-ac*bb < 0 {
		return false
	}
	return true
}

// Returns true if the given point is inside the triangle.
func IsPointInTriangleV2(p, a, b, c *Vector2) bool {
	return IsPointInTriangle(p.X, p.Y, a.X, a.Y, b.X, b.Y, c.X, c.Y)
}

// Returns true if the given point is inside the triangle.
func IsPointInTriangle(px, py, ax, ay, bx, by, cx, cy float32) bool {
	px1 := px - ax
	py1 := py - ay
	side12 := (bx-ax)*py1-(by-ay)*px1 > 0
	if (cx-ax)*py1-(cy-ay)*px1 > 0 == side12 {
		return false
	}
	if (cx-bx)*(py-by)-(cy-by)*(px-bx) > 0 != side12 {
		return false
	}
	return true
}

func IntersectSegmentPlane(start, end *Vector3, plane *Plane, intersection *Vector3) bool {
	dir := v0.SetV(end).SubV(start)
	denom := dir.DotV(plane.GetNormal())
	t := -(start.DotV(plane.GetNormal()) + plane.GetD()) / denom
	if t < 0 || t > 1 {
		return false
	}
	intersection.SetV(start).AddV(dir.SclScalar(t))
	return true
}

// Taken from https://en.wikipedia.org/wiki/Sign_function
func signum(x float32) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

// Determines on which side of the given line the point is. Returns -1 if the point is on the left side of the line, 0 if the
// point is on the line and 1 if the point is on the right side of the line. Left and right are relative to the lines direction
// which is linePoint1 to linePoint2.
func PointLineSideV2(linePoint1, linePoint2, point *Vector2) int {
	return signum((linePoint2.X-linePoint1.X)*(point.Y-linePoint1.Y) -
		(linePoint2.Y-linePoint1.Y)*(point.X-linePoint1.X))
}

func PointLineSide(linePoint1X, linePoint1Y, linePoint2X, linePoint2Y, pointX, pointY float32) int {
	return signum((linePoint2X-linePoint1X)*(pointY-linePoint1Y) -
		(linePoint2Y-linePoint1Y)*(pointX-linePoint1X))
}

// Checks whether the given point is in the polygon.
// param polygon The polygon vertices passed as an array
// param point The point
// return true if the point is in the polygon
func IsPointInPolygonV2(polygon []*Vector2, point *Vector2) bool {
	lastVertice := polygon[len(polygon)-1]
	oddNodes := false
	for i := 0; i < len(polygon); i++ {
		vertice := polygon[i]
		if (vertice.Y < point.Y && lastVertice.Y >= point.Y) || (lastVertice.Y < point.Y && vertice.Y >= point.Y) {
			if vertice.Y+(point.Y-vertice.Y)/(lastVertice.Y-vertice.Y)*(lastVertice.Y-vertice.Y) < point.Y {
				oddNodes = !oddNodes
			}
		}
		lastVertice = vertice
	}
	return oddNodes
}

// Returns true if the specified point is in the polygon.
// param offset Starting polygon index.
// param count Number of array indices to use after offset.
func isPointInPolygon(polygon []float32, offset, count int, x, y float32) bool {
	oddNodes := false
	j := offset + count - 2
	i := offset
	for n := j; i <= n; i += 2 {
		yi := polygon[i+1]
		yj := polygon[j+1]
		if (yi < y && yj >= y) || (yj < y && yi >= y) {
			xi := polygon[i]
			if xi+(y-yi)/(yj-yi)*(polygon[j]-xi) < x {
				oddNodes = !oddNodes
			}
		}
		j = i
	}
	return oddNodes
}

// Returns the distance between the given line and point. Note the specified line is not a line segment.
func DistanceLinePoint(startX, startY, endX, endY, pointX, pointY float32) float32 {
	normalLength := math.Sqrt(float64((endX-startX)*(endX-startX) + (endY-startY)*(endY-startY)))
	return float32(math.Abs(float64((pointX-startX)*(endY-startY)-(pointY-startY)*(endX-startX))) / normalLength)
}

// Returns the distance between the given segment and point.
func DistanceSegmentPoint(startX, startY, endX, endY, pointX, pointY float32) float32 {
	return NearestSegmentPoint(startX, startY, endX, endY, pointX, pointY, v2tmp).Dst(pointX, pointY)
}

// Returns the distance between the given segment and point.
func DistanceSegmentPointV2(start, end, point *Vector2) float32 {
	return NearestSegmentPointV2(start, end, point, v2tmp).DstV(point)
}

// Returns a point on the segment nearest to the specified point.
func NearestSegmentPointV2(start, end, point, nearest *Vector2) *Vector2 {
	length2 := start.Dst2V(end)
	if length2 == 0 {
		return nearest.SetV(start)
	}
	t := ((point.X-start.X)*(end.X-start.X) + (point.Y-start.Y)*(end.Y-start.Y)) / length2
	if t < 0 {
		return nearest.SetV(start)
	}
	if t > 1 {
		return nearest.SetV(end)
	}
	return nearest.Set(start.X+t*(end.X-start.X), start.Y+t*(end.Y-start.Y))
}

// Returns a point on the segment nearest to the specified point.
func NearestSegmentPoint(startX, startY, endX, endY, pointX, pointY float32, nearest *Vector2) *Vector2 {
	xDiff := endX - startX
	yDiff := endY - startY
	length2 := xDiff*xDiff + yDiff*yDiff
	if length2 == 0 {
		return nearest.Set(startX, startY)
	}
	t := ((pointX-startX)*(endX-startX) + (pointY-startY)*(endY-startY)) / length2
	if t < 0 {
		return nearest.Set(startX, startY)
	}
	if t > 1 {
		return nearest.Set(endX, endY)
	}
	return nearest.Set(startX+t*(endX-startX), startY+t*(endY-startY))
}

// Returns whether the given line segment intersects the given circle.
// param start The start point of the line segment
// param end The end point of the line segment
// param center The center of the circle
// param squareRadius The squared radius of the circle
// return Whether the line segment and the circle intersect
func IntersectSegmentCircle(start, end, center *Vector2, squareRadius float32) bool {
	tmp.Set(end.X-start.X, end.Y-start.Y, 0)
	tmp1.Set(center.X-start.X, center.Y-start.Y, 0)
	l := tmp.Len()
	u := tmp1.DotV(tmp.Nor())
	if u <= 0 {
		tmp2.Set(start.X, start.Y, 0)
	} else if u >= l {
		tmp2.Set(end.X, end.Y, 0)
	} else {
		tmp3.SetV(tmp.SclScalar(u)) // remember tmp is already normalized
		tmp2.Set(tmp3.X+start.X, tmp3.Y+start.Y, 0)
	}

	x := center.X - tmp2.X
	y := center.Y - tmp2.Y

	return x*x+y*y <= squareRadius
}

// Checks whether the line segment and the circle intersect and returns by how much and in what direction the line has to move
// away from the circle to not intersect.
//
// param start The line segment starting point
// param end The line segment end point
// param point The center of the circle
// param radius The radius of the circle
// param displacement The displacement vector set by the method having unit length
// return The displacement or math.MaxFloat32 if no intersection is present
func IntersectSegmentCircleDisplace(start, end, point *Vector2, radius float32, displacement *Vector2) float32 {
	u := (point.X-start.X)*(end.X-start.X) + (point.Y-start.Y)*(end.Y-start.Y)
	d := start.DstV(end)
	u /= d * d
	if u < 0 || u > 1 {
		return math.MaxFloat32
	}
	tmp.Set(end.X, end.Y, 0).Sub(start.X, start.Y, 0)
	tmp2.Set(start.X, start.Y, 0).AddV(tmp.SclScalar(u))
	d = tmp2.Dst(point.X, point.Y, 0)
	if d < radius {
		displacement.SetV(point).Sub(tmp2.X, tmp2.Y).Nor()
		return d
	} else {
		return math.MaxFloat32
	}
}

// Intersect two 2D Rays and return the scalar parameter of the first ray at the intersection point.
// You can get the intersection point by: Vector2 point(direction1).scl(scalar).add(start1);
// For more information, check: http://stackoverflow.com/a/565282/1091440
// param start1 Where the first ray start
// param direction1 The direction the first ray is pointing
// param start2 Where the second ray start
// param direction2 The direction the second ray is pointing
// return scalar parameter on the first ray describing the point where the intersection happens. May be negative.
// In case the rays are collinear, math will be returned.
func IntersectRayRay(start1, direction1, start2, direction2 *Vector2) float32 {
	difx := start2.X - start1.X
	dify := start2.Y - start1.Y
	d1xd2 := direction1.X*direction2.Y - direction1.Y*direction2.X
	if d1xd2 == 0.0 {
		return math.MaxFloat32 //collinear
	}
	d2sx := direction2.X / d1xd2
	d2sy := direction2.Y / d1xd2
	return difx*d2sy - dify*d2sx
}

// Intersects a {@link Ray} and a {@link Plane}. The intersection point is stored in intersection in case an intersection is
// present.
// param intersection The vector the intersection point is written to (optional)
// return Whether an intersection is present.
func IntersectRayPlane(ray *Ray, plane *Plane, intersection *Vector3) bool {
	denom := ray.Direction.DotV(plane.GetNormal())
	if denom != 0 {
		t := -(ray.Origin.DotV(plane.GetNormal()) + plane.GetD()) / denom
		if t < 0 {
			return false
		}
		if intersection != nil {
			intersection.SetV(ray.Origin).AddV(v0.SetV(ray.Direction).SclScalar(t))
		}
		return true
	} else if plane.TestPointV3(ray.Origin) == PlaneOn {
		if intersection != nil {
			intersection.SetV(ray.Origin)
		}
		return true
	} else {
		return false
	}
}

// Intersects a line and a plane. The intersection is returned as the distance from the first point to the plane. In case an
// intersection happened, the return value is in the range [0,1]. The intersection point can be recovered by point1 + t *
// (point2 - point1) where t is the return value of this method.
func IntersectLinePlane(x, y, z, x2, y2, z2 float32, plane *Plane, intersection *Vector3) float32 {
	direction := tmp.Set(x2, y2, z2).Sub(x, y, z)
	origin := tmp2.Set(x, y, z)
	denom := direction.DotV(plane.GetNormal())
	if denom != 0 {
		t := -(origin.DotV(plane.GetNormal()) + plane.GetD()) / denom
		if intersection != nil {
			intersection.SetV(origin).AddV(direction.SclScalar(t))
		}
		return t
	} else if plane.TestPointV3(origin) == PlaneOn {
		if intersection != nil {
			intersection.SetV(origin)
		}
		return 0
	}
	return -1
}

var (
	p = NewPlaneEmpty()
	i = NewVector3Empty()
)

// Intersect a {@link Ray} and a triangle, returning the intersection point in intersection.
// param ray The ray
// param t1 The first vertex of the triangle
// param t2 The second vertex of the triangle
// param t3 The third vertex of the triangle
// param intersection The intersection point (optional)
// return True in case an intersection is present.
func IntersectRayTriangle(ray *Ray, t1, t2, t3, intersection *Vector3) bool {
	edge1 := v0.SetV(t2).SubV(t1)
	edge2 := v1.SetV(t3).SubV(t1)
	pvec := v2.SetV(ray.Direction).CrsV(edge2)
	det := edge1.DotV(pvec)
	if utils.IsZero(det) {
		p.SetP3(t1, t2, t3)
		if p.TestPointV3(ray.Origin) == PlaneOn && IsPointInTriangleV3(ray.Origin, t1, t2, t3) {
			if intersection != nil {
				intersection.SetV(ray.Origin)
			}
			return true
		}
		return false
	}

	det = 1.0 / det

	tvec := i.SetV(ray.Origin).SubV(t1)
	u := tvec.DotV(pvec) * det
	if u < 0.0 || u > 1.0 {
		return false
	}

	qvec := tvec.CrsV(edge1)
	v := ray.Direction.DotV(qvec) * det
	if v < 0.0 || u+v > 1.0 {
		return false
	}

	t := edge2.DotV(qvec) * det
	if t < 0 {
		return false
	}

	if intersection != nil {
		if t <= utils.FLOAT_ROUNDING_ERROR {
			intersection.SetV(ray.Origin)
		} else {
			ray.GetEndPoint(intersection, t)
		}
	}

	return true
}

var (
	dir   = NewVector3Empty()
	start = NewVector3Empty()
)

// Intersects a {@link Ray} and a sphere, returning the intersection point in intersection.
// param ray The ray, the direction component must be normalized before calling this method
// param center The center of the sphere
// param radius The radius of the sphere
// param intersection The intersection point (optional, can be null)
// return Whether an intersection is present.
func IntersectRaySphere(ray *Ray, center *Vector3, radius float32, intersection *Vector3) bool {
	len := ray.Direction.Dot(center.X-ray.Origin.X, center.Y-ray.Origin.Y, center.Z-ray.Origin.Z)
	if len < 0 { // behind the ray
		return false
	}
	dst2 := center.Dst2V3(ray.Origin.X+ray.Direction.X*len, ray.Origin.Y+ray.Direction.Y*len, ray.Origin.Z+
		ray.Direction.Z*len)
	r2 := radius * radius
	if dst2 > r2 {
		return false
	}
	if intersection != nil {
		intersection.SetV(ray.Direction).SclScalar(len - float32(math.Sqrt(float64(r2-dst2)))).AddV(ray.Origin)
	}
	return true
}

// Intersects a {@link Ray} and a {@link BoundingBox}, returning the intersection point in intersection.
// This intersection is defined as the point on the ray closest to the origin which is within the specified
// bounds.
//
// The returned intersection (if any) is guaranteed to be within the bounds of the bounding box, but
// it can occasionally diverge slightly from ray, due to small floating-point errors.</p>
//
// If the origin of the ray is inside the box, this method returns true and the intersection point is
// set to the origin of the ray, accordingly to the definition above.</p>
//
// param ray The ray
// param box The box
// param intersection The intersection point (optional)
// return Whether an intersection is present.
func IntersectRayBounds(ray *Ray, box *BoundingBox, intersection *Vector3) bool {
	if box.ContainsV3(ray.Origin) {
		if intersection != nil {
			intersection.SetV(ray.Origin)
		}
		return true
	}
	lowest := float32(0)
	hit := false
	var t float32

	// min x
	if ray.Origin.X <= box.Min.X && ray.Direction.X > 0 {
		t = (box.Min.X - ray.Origin.X) / ray.Direction.X
		if t >= 0 {
			v2.SetV(ray.Direction).SclScalar(t).AddV(ray.Origin)
			if v2.Y >= box.Min.X && v2.Y <= box.Max.Y && v2.Z >= box.Min.X && v2.Z <= box.Max.Z && (!hit || t < lowest) {
				hit = true
				lowest = t
			}
		}
	}
	// max x
	if ray.Origin.X >= box.Max.X && ray.Direction.X < 0 {
		t = (box.Max.X - ray.Origin.X) / ray.Direction.X
		if t >= 0 {
			v2.SetV(ray.Direction).SclScalar(t).AddV(ray.Origin)
			if v2.Y >= box.Min.X && v2.Y <= box.Max.Y && v2.Z >= box.Min.X && v2.Z <= box.Max.Z && (!hit || t < lowest) {
				hit = true
				lowest = t
			}
		}
	}
	// min y
	if ray.Origin.Y <= box.Min.X && ray.Direction.Y > 0 {
		t = (box.Min.X - ray.Origin.Y) / ray.Direction.Y
		if t >= 0 {
			v2.SetV(ray.Direction).SclScalar(t).AddV(ray.Origin)
			if v2.X >= box.Min.X && v2.X <= box.Max.X && v2.Z >= box.Min.X && v2.Z <= box.Max.Z && (!hit || t < lowest) {
				hit = true
				lowest = t
			}
		}
	}
	// max y
	if ray.Origin.Y >= box.Max.Y && ray.Direction.Y < 0 {
		t = (box.Max.Y - ray.Origin.Y) / ray.Direction.Y
		if t >= 0 {
			v2.SetV(ray.Direction).SclScalar(t).AddV(ray.Origin)
			if v2.X >= box.Min.X && v2.X <= box.Max.X && v2.Z >= box.Min.X && v2.Z <= box.Max.Z && (!hit || t < lowest) {
				hit = true
				lowest = t
			}
		}
	}
	// min z
	if ray.Origin.Z <= box.Min.X && ray.Direction.Z > 0 {
		t = (box.Min.X - ray.Origin.Z) / ray.Direction.Z
		if t >= 0 {
			v2.SetV(ray.Direction).SclScalar(t).AddV(ray.Origin)
			if v2.X >= box.Min.X && v2.X <= box.Max.X && v2.Y >= box.Min.X && v2.Y <= box.Max.Y && (!hit || t < lowest) {
				hit = true
				lowest = t
			}
		}
	}
	// max y
	if ray.Origin.Z >= box.Max.Z && ray.Direction.Z < 0 {
		t = (box.Max.Z - ray.Origin.Z) / ray.Direction.Z
		if t >= 0 {
			v2.SetV(ray.Direction).SclScalar(t).AddV(ray.Origin)
			if v2.X >= box.Min.X && v2.X <= box.Max.X && v2.Y >= box.Min.X && v2.Y <= box.Max.Y && (!hit || t < lowest) {
				hit = true
				lowest = t
			}
		}
	}
	if hit && intersection != nil {
		intersection.SetV(ray.Direction).SclScalar(lowest).AddV(ray.Origin)
		if intersection.X < box.Min.X {
			intersection.X = box.Min.X
		} else if intersection.X > box.Max.X {
			intersection.X = box.Max.X
		}
		if intersection.Y < box.Min.Y {
			intersection.Y = box.Min.Y
		} else if intersection.Y > box.Max.Y {
			intersection.Y = box.Max.Y
		}
		if intersection.Z < box.Min.Z {
			intersection.Z = box.Min.Z
		} else if intersection.Z > box.Max.Z {
			intersection.Z = box.Max.Z
		}
	}
	return hit
}

// Quick check whether the given {@link Ray} and {@link BoundingBox} intersect.
// param ray The ray
// param box The bounding box
// return Whether the ray and the bounding box intersect.
func IntersectRayBoundsFast(ray *Ray, box *BoundingBox) bool {
	return IntersectRayBoundsFastV3(ray, box.GetCenter(tmp1), box.GetDimensions(tmp2))
}

// Quick check whether the given {@link Ray} and {@link BoundingBox} intersect.
// param ray The ray
// param center The center of the bounding box
// param dimensions The dimensions (width, height and depth) of the bounding box
// return Whether the ray and the bounding box intersect.
func IntersectRayBoundsFastV3(ray *Ray, center, dimensions *Vector3) bool {
	divX := 1 / ray.Direction.X
	divY := 1 / ray.Direction.Y
	divZ := 1 / ray.Direction.Z

	minx := float64(((center.X - dimensions.X*0.5) - ray.Origin.X) * divX)
	maxx := float64(((center.X + dimensions.X*0.5) - ray.Origin.X) * divX)
	if minx > maxx {
		minx, maxx = maxx, minx
	}

	miny := float64(((center.Y - dimensions.Y*0.5) - ray.Origin.Y) * divY)
	maxy := float64(((center.Y + dimensions.Y*0.5) - ray.Origin.Y) * divY)
	if miny > maxy {
		miny, maxy = maxy, miny
	}

	minz := float64(((center.Z - dimensions.Z*0.5) - ray.Origin.Z) * divZ)
	maxz := float64(((center.Z + dimensions.Z*0.5) - ray.Origin.Z) * divZ)
	if minz > maxz {
		minz, maxz = maxz, minz
	}

	min := math.Max(math.Max(minx, miny), minz)
	max := math.Min(math.Min(maxx, maxy), maxz)

	return max >= 0 && max >= min
}

var (
	best  = NewVector3Empty()
	tmp   = NewVector3Empty()
	tmp1  = NewVector3Empty()
	tmp2  = NewVector3Empty()
	tmp3  = NewVector3Empty()
	v2tmp = NewVector2Empty()
)

// Intersects the given ray with list of triangles. Returns the nearest intersection point in intersection
// param ray The ray
// param triangles The triangles, each successive 3 elements from a vertex
// param intersection The nearest intersection point (optional)
// return Whether the ray and the triangles intersect.
func IntersectRayTriangles(ray *Ray, triangles []float32, intersection *Vector3) bool {
	min_dist := float32(math.MaxFloat32)
	hit := false

	if len(triangles)/3%3 != 0 {
		panic("triangle list size is not a multiple of 3")
	}

	for i := 0; i < len(triangles)-6; i += 9 {
		result := IntersectRayTriangle(ray, tmp1.Set(triangles[i], triangles[i+1], triangles[i+2]),
			tmp2.Set(triangles[i+3], triangles[i+4], triangles[i+5]),
			tmp3.Set(triangles[i+6], triangles[i+7], triangles[i+8]), tmp)

		if result == true {
			dist := ray.Origin.Dst2V(tmp)
			if dist < min_dist {
				min_dist = dist
				best.SetV(tmp)
				hit = true
			}
		}
	}

	if hit == false {
		return false
	} else {
		if intersection != nil {
			intersection.SetV(best)
		}
		return true
	}
}

// Intersects the given ray with list of triangles. Returns the nearest intersection point in intersection
// param ray The ray
// param vertices the vertices
// param indices the indices, each successive 3 shorts index the 3 vertices of a triangle
// param vertexSize the size of a vertex in floats
// param intersection The nearest intersection point (optional)
// return Whether the ray and the triangles intersect.
func IntersectRayTrianglesIndex(ray *Ray, vertices []float32, indices []uint8, vertexSize uint8, intersection *Vector3) bool {
	min_dist := float32(math.MaxFloat32)
	hit := false

	if len(indices)%3 != 0 {
		panic("triangle list size is not a multiple of 3")
	}

	for i := 0; i < len(indices); i += 3 {
		i1 := indices[i] * vertexSize
		i2 := indices[i+1] * vertexSize
		i3 := indices[i+2] * vertexSize

		result := IntersectRayTriangle(ray, tmp1.Set(vertices[i1], vertices[i1+1], vertices[i1+2]),
			tmp2.Set(vertices[i2], vertices[i2+1], vertices[i2+2]),
			tmp3.Set(vertices[i3], vertices[i3+1], vertices[i3+2]), tmp)

		if result == true {
			dist := ray.Origin.Dst2V(tmp)
			if dist < min_dist {
				min_dist = dist
				best.SetV(tmp)
				hit = true
			}
		}
	}

	if hit == false {
		return false
	} else {
		if intersection != nil {
			intersection.SetV(best)
		}
		return true
	}
}

// Intersects the given ray with list of triangles. Returns the nearest intersection point in intersection
// param ray The ray
// param triangles The triangles
// param intersection The nearest intersection point (optional)
// return Whether the ray and the triangles intersect.
func IntersectRayTrianglesV3(ray *Ray, triangles []*Vector3, intersection *Vector3) bool {
	min_dist := float32(math.MaxFloat32)
	hit := false

	if len(triangles)%3 != 0 {
		panic("triangle list size is not a multiple of 3")
	}

	for i := 0; i < len(triangles)-2; i += 3 {
		result := IntersectRayTriangle(ray, triangles[i], triangles[i+1], triangles[i+2], tmp)

		if result == true {
			dist := ray.Origin.Dst2V(tmp)
			if dist < min_dist {
				min_dist = dist
				best.SetV(tmp)
				hit = true
			}
		}
	}

	if !hit {
		return false
	} else {
		if intersection != nil {
			intersection.SetV(best)
		}
		return true
	}
}

// Intersects the two lines and returns the intersection point in intersection.
// param p1 The first point of the first line
// param p2 The second point of the first line
// param p3 The first point of the second line
// param p4 The second point of the second line
// param intersection The intersection point
// return Whether the two lines intersect
func IntersectLinesV3(p1, p2, p3, p4, intersection *Vector2) bool {
	x1 := p1.X
	y1 := p1.Y
	x2 := p2.X
	y2 := p2.Y
	x3 := p3.X
	y3 := p3.Y
	x4 := p4.X
	y4 := p4.Y

	d := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)
	if d == 0 {
		return false
	}

	if intersection != nil {
		ua := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / d
		intersection.Set(x1+(x2-x1)*ua, y1+(y2-y1)*ua)
	}
	return true
}

// Intersects the two lines and returns the intersection point in intersection.
// param intersection The intersection point, or null.
// return Whether the two lines intersect
func IntersectLines(x1, y1, x2, y2, x3, y3, x4, y4 float32, intersection *Vector2) bool {
	d := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)
	if d == 0 {
		return false
	}
	if intersection != nil {
		ua := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / d
		intersection.Set(x1+(x2-x1)*ua, y1+(y2-y1)*ua)
	}
	return true
}

// Check whether the given line and {@link Polygon} intersect.
// param p1 The first point of the line
// param p2 The second point of the line
// param polygon The polygon
// return Whether polygon and line intersects
func IntersectLinePolygon(p1, p2 *Vector2, polygon *Polygon) bool {
	vertices := polygon.GetTransformedVertices()
	x1 := p1.X
	y1 := p1.Y
	x2 := p2.X
	y2 := p2.Y
	n := len(vertices)
	x3 := vertices[n-2]
	y3 := vertices[n-1]
	for i := 0; i < n; i += 2 {
		x4 := vertices[i]
		y4 := vertices[i+1]
		d := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)
		if d != 0 {
			yd := y1 - y3
			xd := x1 - x3
			ua := ((x4-x3)*yd - (y4-y3)*xd) / d
			if ua >= 0 && ua <= 1 {
				return true
			}
		}
		x3 = x4
		y3 = y4
	}
	return false
}

// Determines whether the given rectangles intersect and, if they do, sets the supplied {@code intersection} rectangle to the
// area of overlap.
// return Whether the rectangles intersect
func IntersectRectangles(rectangle1, rectangle2, intersection *Rectangle) bool {
	if rectangle1.Overlaps(rectangle2) {
		intersection.X = float32(math.Max(float64(rectangle1.X), float64(rectangle2.X)))
		intersection.W = float32(math.Min(float64(rectangle1.X+rectangle1.W), float64(rectangle2.X+rectangle2.W)) - float64(intersection.X))
		intersection.Y = float32(math.Max(float64(rectangle1.Y), float64(rectangle2.Y)))
		intersection.H = float32(math.Min(float64(rectangle1.Y+rectangle1.H), float64(rectangle2.Y+rectangle2.H)) - float64(intersection.Y))
		return true
	}
	return false
}

// Check whether the given line segment and {@link Polygon} intersect.
// param p1 The first point of the segment
// param p2 The second point of the segment
// return Whether polygon and segment intersect
func IntersectSegmentPolygon(p1, p2 *Vector2, polygon *Polygon) bool {
	vertices := polygon.GetTransformedVertices()
	x1 := p1.X
	y1 := p1.Y
	x2 := p2.X
	y2 := p2.Y
	n := len(vertices)
	x3 := vertices[n-2]
	y3 := vertices[n-1]
	for i := 0; i < n; i += 2 {
		x4 := vertices[i]
		y4 := vertices[i+1]
		d := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)
		if d != 0 {
			yd := y1 - y3
			xd := x1 - x3
			ua := ((x4-x3)*yd - (y4-y3)*xd) / d
			if ua >= 0 && ua <= 1 {
				ub := ((x2-x1)*yd - (y2-y1)*xd) / d
				if ub >= 0 && ub <= 1 {
					return true
				}
			}
		}
		x3 = x4
		y3 = y4
	}
	return false
}

// Intersects the two line segments and returns the intersection point in intersection.
// param p1 The first point of the first line segment
// param p2 The second point of the first line segment
// param p3 The first point of the second line segment
// param p4 The second point of the second line segment
// param intersection The intersection point (optional)
// return Whether the two line segments intersect
func IntersectSegmentsV2(p1, p2, p3, p4, intersection *Vector2) bool {
	x1 := p1.X
	y1 := p1.Y
	x2 := p2.X
	y2 := p2.Y
	x3 := p3.X
	y3 := p3.Y
	x4 := p4.X
	y4 := p4.Y

	d := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)
	if d == 0 {
		return false
	}

	yd := y1 - y3
	xd := x1 - x3
	ua := ((x4-x3)*yd - (y4-y3)*xd) / d
	if ua < 0 || ua > 1 {
		return false
	}

	ub := ((x2-x1)*yd - (y2-y1)*xd) / d
	if ub < 0 || ub > 1 {
		return false
	}

	if intersection != nil {
		intersection.Set(x1+(x2-x1)*ua, y1+(y2-y1)*ua)
	}
	return true
}

func IntersectSegments(x1, y1, x2, y2, x3, y3, x4, y4 float32, intersection *Vector2) bool {
	d := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)
	if d == 0 {
		return false
	}

	yd := y1 - y3
	xd := x1 - x3
	ua := ((x4-x3)*yd - (y4-y3)*xd) / d
	if ua < 0 || ua > 1 {
		return false
	}

	ub := ((x2-x1)*yd - (y2-y1)*xd) / d
	if ub < 0 || ub > 1 {
		return false
	}

	if intersection != nil {
		intersection.Set(x1+(x2-x1)*ua, y1+(y2-y1)*ua)
	}
	return true
}

func Det(a, b, c, d float32) float32 {
	return a*d - b*c
}

func Det64(a, b, c, d float64) float64 {
	return a*d - b*c
}

func OverlapsC(c1, c2 *Circle) bool {
	return c1.Overlaps(c2)
}

func OverlapsR(r1, r2 *Rectangle) bool {
	return r1.Overlaps(r2)
}

func OverlapsCR(c *Circle, r *Rectangle) bool {
	closestX := c.X
	closestY := c.Y

	if c.X < r.X {
		closestX = r.X
	} else if c.X > r.X+r.W {
		closestX = r.X + r.W
	}

	if c.Y < r.Y {
		closestY = r.Y
	} else if c.Y > r.Y+r.H {
		closestY = r.Y + r.H
	}

	closestX = closestX - c.X
	closestX *= closestX
	closestY = closestY - c.Y
	closestY *= closestY

	return closestX+closestY < c.Radius*c.Radius
}

// Check whether specified counter-clockwise wound convex polygons overlap. If they do, optionally obtain a Minimum Translation Vector indicating the
// minimum magnitude vector required to push the polygon p1 out of collision with polygon p2.
// param p1 The first polygon.
// param p2 The second polygon.
// param mtv A Minimum Translation Vector to fill in the case of a collision, or null (optional).
// return Whether polygons overlap.
func OverlapConvexPolygonsP(p1, p2 *Polygon, mtv *MinimumTranslationVector) bool {
	return OverlapConvexPolygonsV(p1.GetTransformedVertices(), p2.GetTransformedVertices(), mtv)
}

// @see #overlapConvexPolygons(float[], int, int, float[], int, int, MinimumTranslationVector)
func OverlapConvexPolygonsV(verts1, verts2 []float32, mtv *MinimumTranslationVector) bool {
	return OverlapConvexPolygons(verts1, 0, len(verts1), verts2, 0, len(verts2), mtv)
}

// Check whether polygons defined by the given counter-clockwise wound vertex arrays overlap. If they do, optionally obtain a Minimum Translation
// Vector indicating the minimum magnitude vector required to push the polygon defined by verts1 out of the collision with the polygon defined by verts2.
// param verts1 Vertices of the first polygon.
// param verts2 Vertices of the second polygon.
// param mtv A Minimum Translation Vector to fill in the case of a collision, or null (optional).
// return Whether polygons overlap.
func OverlapConvexPolygons(verts1 []float32, offset1 int, count1 int, verts2 []float32, offset2 int, count2 int,
	mtv *MinimumTranslationVector) bool {
	overlap := math.MaxFloat32
	smallestAxisX := float32(0)
	smallestAxisY := float32(0)
	var numInNormalDir int

	end1 := offset1 + count1
	end2 := offset2 + count2

	// Get polygon1 axes
	for i := offset1; i < end1; i += 2 {
		x1 := verts1[i]
		y1 := verts1[i+1]
		x2 := verts1[(i+2)%count1]
		y2 := verts1[(i+3)%count1]

		axisX := y1 - y2
		axisY := -(x1 - x2)

		length := float32(math.Sqrt(float64(axisX*axisX + axisY*axisY)))
		axisX /= length
		axisY /= length

		// -- Begin check for separation on this axis --//

		// Project polygon1 onto this axis
		min1 := axisX*verts1[0] + axisY*verts1[1]
		max1 := min1
		for j := offset1; j < end1; j += 2 {
			p := axisX*verts1[j] + axisY*verts1[j+1]
			if p < min1 {
				min1 = p
			} else if p > max1 {
				max1 = p
			}
		}

		// Project polygon2 onto this axis
		numInNormalDir = 0
		min2 := axisX*verts2[0] + axisY*verts2[1]
		max2 := min2
		for j := offset2; j < end2; j += 2 {
			// Counts the number of points that are within the projected area.
			numInNormalDir -= PointLineSide(x1, y1, x2, y2, verts2[j], verts2[j+1])
			p := axisX*verts2[j] + axisY*verts2[j+1]
			if p < min2 {
				min2 = p
			} else if p > max2 {
				max2 = p
			}
		}

		if !(min1 <= min2 && max1 >= min2 || min2 <= min1 && max2 >= min1) {
			return false
		} else {
			o := math.Min(float64(max1), float64(max2)) - math.Max(float64(min1), float64(min2))
			if min1 < min2 && max1 > max2 || min2 < min1 && max2 > max1 {
				mins := math.Abs(float64(min1 - min2))
				maxs := math.Abs(float64(max1 - max2))
				if mins < maxs {
					o += mins
				} else {
					o += maxs
				}
			}
			if o < overlap {
				overlap = o
				// Adjusts the direction based on the number of points found
				if numInNormalDir >= 0 {
					smallestAxisX = axisX
					smallestAxisY = axisY
				} else {
					smallestAxisX = -axisX
					smallestAxisY = -axisY
				}
			}
		}
		// -- End check for separation on this axis --//
	}

	// Get polygon2 axes
	for i := offset2; i < end2; i += 2 {
		x1 := verts2[i]
		y1 := verts2[i+1]
		x2 := verts2[(i+2)%count2]
		y2 := verts2[(i+3)%count2]

		axisX := y1 - y2
		axisY := -(x1 - x2)

		length := float32(math.Sqrt(float64(axisX*axisX + axisY*axisY)))
		axisX /= length
		axisY /= length

		// -- Begin check for separation on this axis --//
		numInNormalDir = 0

		// Project polygon1 onto this axis
		min1 := axisX*verts1[0] + axisY*verts1[1]
		max1 := min1
		for j := offset1; j < end1; j += 2 {
			p := axisX*verts1[j] + axisY*verts1[j+1]
			// Counts the number of points that are within the projected area.
			numInNormalDir -= PointLineSide(x1, y1, x2, y2, verts1[j], verts1[j+1])
			if p < min1 {
				min1 = p
			} else if p > max1 {
				max1 = p
			}
		}

		// Project polygon2 onto this axis
		min2 := axisX*verts2[0] + axisY*verts2[1]
		max2 := min2
		for j := offset2; j < end2; j += 2 {
			p := axisX*verts2[j] + axisY*verts2[j+1]
			if p < min2 {
				min2 = p
			} else if p > max2 {
				max2 = p
			}
		}

		if !(min1 <= min2 && max1 >= min2 || min2 <= min1 && max2 >= min1) {
			return false
		} else {
			o := math.Min(float64(max1), float64(max2)) - math.Max(float64(min1), float64(min2))

			if min1 < min2 && max1 > max2 || min2 < min1 && max2 > max1 {
				mins := math.Abs(float64(min1 - min2))
				maxs := math.Abs(float64(max1 - max2))
				if mins < maxs {
					o += mins
				} else {
					o += maxs
				}
			}

			if o < overlap {
				overlap = o
				// Adjusts the direction based on the number of points found
				if numInNormalDir < 0 {
					smallestAxisX = axisX
					smallestAxisY = axisY
				} else {
					smallestAxisX = -axisX
					smallestAxisY = -axisY
				}

			}
		}
		// -- End check for separation on this axis --//
	}
	if mtv != nil {
		mtv.normal.Set(smallestAxisX, smallestAxisY)
		mtv.depth = float32(overlap)
	}
	return true
}

// Splits the triangle by the plane. The result is stored in the SplitTriangle instance. Depending on where the triangle is
// relative to the plane, the result can be:
// <li>Triangle is fully in front/behind: {@link SplitTriangle#front} or {@link SplitTriangle#back} will contain the original
// triangle, {@link SplitTriangle#total} will be one.</li>
// <li>Triangle has two vertices in front, one behind: {@link SplitTriangle#front} contains 2 triangles,
// {@link SplitTriangle#back} contains 1 triangles, {@link SplitTriangle#total} will be 3.</li>
// <li>Triangle has one vertex in front, two behind: {@link SplitTriangle#front} contains 1 triangle,
// {@link SplitTriangle#back} contains 2 triangles, {@link SplitTriangle#total} will be 3.</li>
// The input triangle should have the form: x, y, z, x2, y2, z2, x3, y3, z3. One can add additional attributes per vertex which
// will be interpolated if split, such as texture coordinates or normals. Note that these additional attributes won't be
// normalized, as might be necessary in case of normals.
// param triangle
// param plane
// param split output SplitTriangle
func SplitTriangleF(triangle []float32, plane *Plane, split SplitTriangle) {
	stride := len(triangle) / 3
	r1 := plane.TestPoint(triangle[0], triangle[1], triangle[2]) == PlaneBack
	r2 := plane.TestPoint(triangle[0+stride], triangle[1+stride], triangle[2+stride]) == PlaneBack
	r3 := plane.TestPoint(triangle[0+stride*2], triangle[1+stride*2], triangle[2+stride*2]) == PlaneBack

	split.Reset()

	// easy case, triangle is on one side (point on plane means front).
	if r1 == r2 && r2 == r3 {
		split.total = 1
		if r1 {
			split.numBack = 1
			// System.arraycopy(triangle, 0, split.back, 0, len(triangle));
		} else {
			split.numFront = 1
			// System.arraycopy(triangle, 0, split.front, 0, len(triangle));
		}
		return
	}

	// set number of triangles
	split.total = 3
	if r1 {
		split.numFront = 0
	} else {
		split.numFront = 1
	}
	if r2 {
		split.numFront += 0
	} else {
		split.numFront += 1
	}
	if r3 {
		split.numFront += 0
	} else {
		split.numFront += 1
	}
	split.numBack = split.total - split.numFront

	// hard case, split the three edges on the plane
	// determine which array to fill first, front or back, flip if we
	// cross the plane
	split.FrontCurrent = !r1

	// split first edge
	first := 0
	second := stride
	if r1 != r2 {
		// split the edge
		SplitEdge(triangle, first, second, stride, plane, split.edgeSplit, 0)

		// add first edge vertex and new vertex to current side
		split.Add(triangle, first, stride)
		split.Add(split.edgeSplit, 0, stride)

		// flip side and add new vertex and second edge vertex to current side
		split.FrontCurrent = !split.FrontCurrent
		split.Add(split.edgeSplit, 0, stride)
	} else {
		// add both vertices
		split.Add(triangle, first, stride)
	}

	// split second edge
	first = stride
	second = stride + stride
	if r2 != r3 {
		// split the edge
		SplitEdge(triangle, first, second, stride, plane, split.edgeSplit, 0)

		// add first edge vertex and new vertex to current side
		split.Add(triangle, first, stride)
		split.Add(split.edgeSplit, 0, stride)

		// flip side and add new vertex and second edge vertex to current side
		split.FrontCurrent = !split.FrontCurrent
		split.Add(split.edgeSplit, 0, stride)
	} else {
		// add both vertices
		split.Add(triangle, first, stride)
	}

	// split third edge
	first = stride + stride
	second = 0
	if r3 != r1 {
		// split the edge
		SplitEdge(triangle, first, second, stride, plane, split.edgeSplit, 0)

		// add first edge vertex and new vertex to current side
		split.Add(triangle, first, stride)
		split.Add(split.edgeSplit, 0, stride)

		// flip side and add new vertex and second edge vertex to current side
		split.FrontCurrent = !split.FrontCurrent
		split.Add(split.edgeSplit, 0, stride)
	} else {
		// add both vertices
		split.Add(triangle, first, stride)
	}

	// triangulate the side with 2 triangles
	if split.numFront == 2 {
		// System.arraycopy(split.front, stride * 2, split.front, stride * 3, stride * 2);
		// System.arraycopy(split.front, 0, split.front, stride * 5, stride);
	} else {
		// System.arraycopy(split.back, stride * 2, split.back, stride * 3, stride * 2);
		// System.arraycopy(split.back, 0, split.back, stride * 5, stride);
	}
}

var sintersection = NewVector3Empty()

func SplitEdge(vertices []float32, s, e, stride int, plane *Plane, split []float32, offset int) {
	t := IntersectLinePlane(vertices[s], vertices[s+1], vertices[s+2], vertices[e], vertices[e+1],
		vertices[e+2], plane, sintersection)
	split[offset+0] = sintersection.X
	split[offset+1] = sintersection.Y
	split[offset+2] = sintersection.Z
	for i := 3; i < stride; i++ {
		a := vertices[s+i]
		b := vertices[e+i]
		split[offset+i] = a + t*(b-a)
	}
}

type SplitTriangle struct {
	front        []float32
	back         []float32
	edgeSplit    []float32
	numFront     int
	numBack      int
	total        int
	FrontCurrent bool
	frontOffset  int
	backOffset   int
}

// Creates a new instance, assuming numAttributes attributes per triangle vertex.
// param numAttributes must be >= 3
func NewSplitTriangle(numAttributes int) *SplitTriangle {
	return &SplitTriangle{
		front:     make([]float32, numAttributes*3*2),
		back:      make([]float32, numAttributes*3*2),
		edgeSplit: make([]float32, numAttributes),
	}
}

func (self *SplitTriangle) Add(vertex []float32, offset, stride int) {
	if self.FrontCurrent {
		// System.arraycopy(vertex, offset, front, frontOffset, stride);
		self.frontOffset += stride
	} else {
		// System.arraycopy(vertex, offset, back, backOffset, stride);
		self.backOffset += stride
	}
}

func (self *SplitTriangle) Reset() {
	self.FrontCurrent = false
	self.frontOffset = 0
	self.backOffset = 0
	self.numFront = 0
	self.numBack = 0
	self.total = 0
}

type MinimumTranslationVector struct {
	normal Vector2
	depth  float32
}
