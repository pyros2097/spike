// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

import (
	. "github.com/pyros2097/spike/math/collision"
	. "github.com/pyros2097/spike/math/vector"
)

var (
	clipSpacePlanePoints = []*Vector3{NewVector3(-1, -1, -1), NewVector3(1, -1, -1),
		NewVector3(1, 1, -1), NewVector3(-1, 1, -1), // near clip
		NewVector3(-1, -1, 1), NewVector3(1, -1, 1), NewVector3(1, 1, 1), NewVector3(-1, 1, 1)} // far clip

	clipSpacePlanePointsArray [8 * 3]float32

	tmpV = NewVector3Empty()
)

func init() {
	j := 0
	for _, v := range clipSpacePlanePoints {
		clipSpacePlanePointsArray[j] = v.X
		j++
		clipSpacePlanePointsArray[j] = v.Y
		j++
		clipSpacePlanePointsArray[j] = v.Z
		j++
	}
}

// A truncated rectangular pyramid. Used to define the viewable region and its projection onto the screen.
type Frustum struct {
	// the six clipping planes, near, far, left, right, top, bottom *
	Planes [6]*Plane
	// eight points making up the near and far clipping "rectangles". order is counter clockwise, starting at bottom left *
	PlanePoints      []*Vector3
	PlanePointsArray [8 * 3]float32
}

func NewFrustumEmpty() *Frustum {
	self := &Frustum{}
	self.PlanePoints = []*Vector3{
		NewVector3Empty(), NewVector3Empty(), NewVector3Empty(), NewVector3Empty(),
		NewVector3Empty(), NewVector3Empty(), NewVector3Empty(),
	}
	for i := 0; i < 6; i++ {
		self.Planes[i] = NewPlane(NewVector3Empty(), 0)
	}
	return self
}

// Updates the clipping plane's based on the given inverse combined projection and view matrix, e.g. from an
// OrthographicCamera or PerspectiveCamera
// inverseProjectionView the combined projection and view matrices.
func (self *Frustum) Update(inverseProjectionView *Matrix4) {
	// System.arraycopy(clipSpacePlanePointsArray, 0, planePointsArray, 0, clipSpacePlanePointsArray.length);
	// Matrix4.prj(inverseProjectionView.val, planePointsArray, 0, 8, 3);
	j := 0
	for i := 0; i < 8; i++ {
		self.PlanePoints[i].X = self.PlanePointsArray[j]
		j++
		self.PlanePoints[i].Y = self.PlanePointsArray[j]
		j++
		self.PlanePoints[i].Z = self.PlanePointsArray[j]
		j++
	}
	self.Planes[0].SetP3(self.PlanePoints[1], self.PlanePoints[0], self.PlanePoints[2])
	self.Planes[1].SetP3(self.PlanePoints[4], self.PlanePoints[5], self.PlanePoints[7])
	self.Planes[2].SetP3(self.PlanePoints[0], self.PlanePoints[4], self.PlanePoints[3])
	self.Planes[3].SetP3(self.PlanePoints[5], self.PlanePoints[1], self.PlanePoints[6])
	self.Planes[4].SetP3(self.PlanePoints[2], self.PlanePoints[3], self.PlanePoints[6])
	self.Planes[5].SetP3(self.PlanePoints[4], self.PlanePoints[0], self.PlanePoints[1])
}

// Returns whether the point is in the frustum.
// point The point
// return Whether the point is in the frustum.
func (self *Frustum) PointInFrustumV3(point *Vector3) bool {
	for i := 0; i < len(self.Planes); i++ {
		result := self.Planes[i].TestPointV3(point)
		if result == PlaneBack {
			return false
		}
	}
	return true
}

// Returns whether the point is in the frustum.
// x The X coordinate of the point
// y The Y coordinate of the point
// z The Z coordinate of the point
// return Whether the point is in the frustum.
func (self *Frustum) PointInFrustum(x, y, z float32) bool {
	for i := 0; i < len(self.Planes); i++ {
		result := self.Planes[i].TestPoint(x, y, z)
		if result == PlaneBack {
			return false
		}
	}
	return true
}

// Returns whether the given sphere is in the frustum.
// center The center of the sphere
// radius The radius of the sphere
// return Whether the sphere is in the frustum
func (self *Frustum) sphereInFrustumV3(center *Vector3, radius float32) bool {
	for i := 2; i < 6; i++ {
		if (self.Planes[i].normal.X*center.X + self.Planes[i].normal.Y*center.Y +
			self.Planes[i].normal.Z*center.Z) < (-radius - self.Planes[i].d) {
			return false
		}
	}
	return true
}

// Returns whether the given sphere is in the frustum.
// x The X coordinate of the center of the sphere
// y The Y coordinate of the center of the sphere
// z The Z coordinate of the center of the sphere
// radius The radius of the sphere
// return Whether the sphere is in the frustum
func (self *Frustum) SphereInFrustum(x, y, z, radius float32) bool {
	for i := 2; i < 6; i++ {
		if (self.Planes[i].normal.X*x + self.Planes[i].normal.Y*y +
			self.Planes[i].normal.Z*z) < (-radius - self.Planes[i].d) {
			return false
		}
	}
	return true
}

// Returns whether the given sphere is in the frustum not checking whether it is behind the near and far clipping plane.
// center The center of the sphere
// radius The radius of the sphere
// return Whether the sphere is in the frustum
func (self *Frustum) SphereInFrustumWithoutNearFarV3(center *Vector3, radius float32) bool {
	for i := 2; i < 6; i++ {
		if (self.Planes[i].normal.X*center.X + self.Planes[i].normal.Y*center.Y +
			self.Planes[i].normal.Z*center.Z) < (-radius - self.Planes[i].d) {
			return false
		}
	}
	return true
}

// Returns whether the given sphere is in the frustum not checking whether it is behind the near and far clipping plane.
// x The X coordinate of the center of the sphere
// y The Y coordinate of the center of the sphere
// z The Z coordinate of the center of the sphere
// radius The radius of the sphere
// return Whether the sphere is in the frustum
func (self *Frustum) SphereInFrustumWithoutNearFar(x, y, z, radius float32) bool {
	for i := 2; i < 6; i++ {
		if (self.Planes[i].normal.X*x + self.Planes[i].normal.Y*y +
			self.Planes[i].normal.Z*z) < (-radius - self.Planes[i].d) {
			return false
		}
	}
	return true
}

// Returns whether the given {@link BoundingBox} is in the frustum.
// bounds The bounding box
// return Whether the bounding box is in the frustum
func (self *Frustum) BoundsInFrustumBox(bounds *BoundingBox) bool {
	for i := 0; i < len(self.Planes); i++ {
		if self.Planes[i].TestPointV3(bounds.GetCorner000(tmpV)) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPointV3(bounds.GetCorner001(tmpV)) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPointV3(bounds.GetCorner010(tmpV)) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPointV3(bounds.GetCorner011(tmpV)) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPointV3(bounds.GetCorner100(tmpV)) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPointV3(bounds.GetCorner101(tmpV)) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPointV3(bounds.GetCorner110(tmpV)) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPointV3(bounds.GetCorner111(tmpV)) != PlaneBack {
			continue
		}
		return false
	}

	return true
}

// Returns whether the given bounding box is in the frustum.
// return Whether the bounding box is in the frustum
func (self *Frustum) BoundsInFrustumV3(center, dimensions *Vector3) bool {
	return self.BoundsInFrustum(center.X, center.Y, center.Z, dimensions.X/2, dimensions.Y/2, dimensions.Z/2)
}

// Returns whether the given bounding box is in the frustum.
// return Whether the bounding box is in the frustum
func (self *Frustum) BoundsInFrustum(x, y, z, halfWidth, halfHeight, halfDepth float32) bool {
	for i := 0; i < len(self.Planes); i++ {
		if self.Planes[i].TestPoint(x+halfWidth, y+halfHeight, z+halfDepth) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPoint(x+halfWidth, y+halfHeight, z-halfDepth) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPoint(x+halfWidth, y-halfHeight, z+halfDepth) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPoint(x+halfWidth, y-halfHeight, z-halfDepth) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPoint(x-halfWidth, y+halfHeight, z+halfDepth) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPoint(x-halfWidth, y+halfHeight, z-halfDepth) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPoint(x-halfWidth, y-halfHeight, z+halfDepth) != PlaneBack {
			continue
		}
		if self.Planes[i].TestPoint(x-halfWidth, y-halfHeight, z-halfDepth) != PlaneBack {
			continue
		}
		return false
	}

	return true
}

// COMMENTED
// Calculates the pick ray for the given window coordinates. Assumes the window coordinate system has it's y downwards. The
// returned Ray is a member of this instance so don't reuse it outside this class.
//
// param screen_width The window width in pixels
// param screen_height The window height in pixels
// param mouse_x The window x-coordinate
// param mouse_y The window y-coordinate
// param pos The camera position
// param dir The camera direction, having unit length
// param up The camera up vector, having unit length
// return the picking ray.
// func (self *Frustum) CalculatePickRay(screen_width, screen_height, mouse_x, mouse_y float32,  pos,  dir, up *Vector3) {
//   n_x := mouse_x - screen_width / 2.0
//   n_y := mouse_y - screen_height / 2.0
//   n_x /= screen_width / 2.0
//   n_y /= screen_height / 2.0

//   Z.set(dir.tmp().mul(-1)).nor();
//   X.set(up.tmp().crs(Z)).nor();
//   Y.set(Z.tmp().crs(X)).nor();
//   near_center.set(pos.tmp3().sub(Z.tmp2().mul(near)));
//   Vector3 near_point = X.tmp3().mul(near_width).mul(n_x).add(Y.tmp2().mul(near_height).mul(n_y));
//   near_point.add(near_center);

//   return ray.set(near_point.tmp(), near_point.sub(pos).nor());
// }
