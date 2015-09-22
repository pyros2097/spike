// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package graphics

import (
	"math"

	. "github.com/pyros2097/spike/math/collision"
	. "github.com/pyros2097/spike/math/vector"
)

// Base class for OrthographicCamera and PerspectiveCamera
type Camera struct {
	// the position of the camera
	Position *Vector3
	// the unit length direction vector of the camera
	Direction *Vector3
	// the unit length up vector of the camera
	Up *Vector3

	// the projection matrix
	Projection *Matrix4
	// the view matrix
	View *Matrix4
	// the combined projection and view matrix
	Combined *Matrix4
	// the inverse combined projection and view matrix
	InvProjectionView *Matrix4

	// the near clipping plane distance, has to be positive
	Near float32
	// the far clipping plane distance, has to be positive
	Far float32

	// the viewport width
	ViewportWidth float32
	// the viewport height
	ViewportHeight float32

	// the frustum
	// public final Frustum frustum = new Frustum();

	tmpVec *Vector3
	Ray    *Ray

	// the field of view of the height, in degrees for Perspective
	FieldOfView float32
	// the zoom of the camera for Orthographic
	Zoom float32
	tmp  *Vector3
}

// Constructs a new OrthographicCamera, using the given viewport width and height. For pixel perfect 2D rendering just supply
// the screen size, for other unit scales (e.g. meters for box2d) proceed accordingly. The camera will show the region
// [-viewportWidth/2, -(viewportHeight/2-1)] - [(viewportWidth/2-1), viewportHeight/2]
// func NewOrthographicCamera(viewportWidth, viewportHeight float32) *OrthographicCamera {
//   cam := &OrthographicCamera{Zoom:1, Near: 0, ViewportWidth: viewportWidth, ViewportHeight: viewportHeight}
//   cam.Update()
// }
//
// func NewPerspectiveCameraEmpty() *PerspectiveCamera {
//   return &PerspectiveCamera{FieldOfView: 67}
// }
// Constructs a new {@link PerspectiveCamera} with the given field of view and viewport size. The aspect ratio is derived from
//the viewport size.
//@param fieldOfViewY the field of view of the height, in degrees, the field of view for the width will be calculated
//          according to the aspect ratio.
//@param viewportWidth the viewport width
//@param viewportHeight the viewport height
func NewPerspectiveCamera(fieldOfViewY, viewportWidth, viewportHeight float32) *PerspectiveCamera {
	self.fieldOfView = fieldOfViewY
	self.viewportWidth = viewportWidth
	self.viewportHeight = viewportHeight
	self.Update()
}
func NewCamera(viewportWidth, viewportHeight float32) *Camera {
	return &Camera{
		Position:          NewVector3Empty(),
		Direction:         NewVector3(0, 0, -1),
		Up:                NewVector3(0, 1, 0),
		Projection:        NewMatrix4Empty(),
		View:              NewMatrix4Empty(),
		Combined:          NewMatrix4Empty(),
		InvProjectionView: NewMatrix4Empty(),
		Near:              1, //Near = 0 if orthographic camera
		Far:               100,
		ViewportWidth:     0,
		ViewportHeight:    0,
		tmpVec:            NewVector3Empty(),
		tmp:               NewVector3Empty(),
		Ray:               NewRay(NewVector3Empty(), NewVector3Empty()),
		Zoom:              1,
		FieldOfView:       67,
	}
}

// Recalculates the projection and view matrix of this camera and the {@link Frustum} planes. Use this after you've manipulated
// any of the attributes of the camera.
func (self *Camera) Update() {
	self.UpdateFrumstum(true)
}

// Recalculates the projection and view matrix of this camera and the {@link Frustum} planes if <code>updateFrustum</code> is
// true. Use this after you've manipulated any of the attributes of the camera.
func (self *Camera) UpdateFrumstum(updateFrustum bool) {}

func (self *PerspectiveCamera) UpdateFrustumPerspective(updateFrustum bool) {
	aspect = viewportWidth / viewportHeight
	projection.setToProjection(Math.abs(near), Math.abs(far), fieldOfView, aspect)
	view.setToLookAt(position, tmp.set(position).add(direction), up)
	combined.set(projection)
	Matrix4.mul(combined.val, view.val)

	if updateFrustum {
		invProjectionView.set(combined)
		Matrix4.inv(invProjectionView.val)
		frustum.update(invProjectionView)
	}
}

func (self *OrthographicCamera) UpdateFrustumOrtho(updateFrustum bool) {
	projection.setToOrtho(zoom*-viewportWidth/2, zoom*(viewportWidth/2), zoom*-(viewportHeight/2),
		zoom*viewportHeight/2, near, far)
	view.SetToLookAt(position, tmp.set(position).add(direction), up)
	combined.SetV(projection)
	Matrix4.MulM4(combined.val, view.val)

	if updateFrustum {
		invProjectionView.set(combined)
		Matrix4.inv(invProjectionView.val)
		frustum.update(invProjectionView)
	}
}

// Recalculates the direction of the camera to look at the point (x, y, z). This function assumes the up vector is normalized.
// param x the x-coordinate of the point to look at
// param y the x-coordinate of the point to look at
// param z the x-coordinate of the point to look at
func (self *Camera) LookAt(x, y, z float32) {
	tmpVec.Set(x, y, z).SubV3(position).Nor()
	if !tmpVec.IsZero() {
		dot := tmpVec.DotV(up) // up and direction must ALWAYS be orthonormal vectors
		if Math.abs(dot-1) < 0.000000001 {
			// Collinear
			up.set(self.Direction).SclScalar(-1)
		} else if Math.abs(dot+1) < 0.000000001 {
			// Collinear opposite
			up.setV(self.Direction)
		}
		direction.set(tmpVec)
		self.NormalizeUp()
	}
}

// Recalculates the direction of the camera to look at the point (x, y, z).
// @param target the point to look at
func (self *Camera) LookAtV3(target *Vector3) {
	self.LookAt(target.X, target.Y, target.Z)
}

// Normalizes the up vector by first calculating the right vector via a cross product between direction and up, and then
// recalculating the up vector via a cross product between right and direction.
func (self *Camera) NormalizeUp() {
	tmpVec.SetV3(self.Direction).CrsV(self.Up).Nor()
	self.Up.SetV3(tmpVec).CrsV(self.Direction).Nor()
}

// Rotates the direction and up vector of this camera by the given angle around the given axis. The direction and up vector
// will not be orthogonalized.
// param angle the angle
// param axisX the x-component of the axis
// param axisY the y-component of the axis
// param axisZ the z-component of the axis
func (self *Camera) Rotate(angle, axisX, axisY, axisZ float32) {
	self.Direction.Rotate(angle, axisX, axisY, axisZ)
	self.Up.Rotate(angle, axisX, axisY, axisZ)
}

// Rotates the direction and up vector of this camera by the given angle around the given axis. The direction and up vector
// will not be orthogonalized.
// param axis the axis to rotate around
// param angle the angle
func (self *Camera) RotateV3(axis *Vector3, angle float32) {
	self.Direction.Rotate(axis, angle)
	self.Up.Rotate(axis, angle)
}

// Rotates the direction and up vector of this camera by the given rotation matrix. The direction and up vector will not be
// orthogonalized.
// param transform The rotation matrix
func (self *Camera) RotateM4(transform *Matrix4) {
	self.Direction.Rot(transform)
	self.Up.Rot(transform)
}

// Rotates the direction and up vector of this camera by the given {@link Quaternion}. The direction and up vector will not be
// orthogonalized.
// param quat The quaternion
func (self *Camera) RotateQ(quat *Quaternion) {
	quat.Transform(direction)
	quat.Transform(up)
}

// Rotates the direction and up vector of this camera by the given angle around the given axis, with the axis attached to given
// point. The direction and up vector will not be orthogonalized.
// param point the point to attach the axis to
// param axis the axis to rotate around
// param angle the angle
func (self *Camera) RotateAround(point, axis *Vector3, angle float32) {
	tmpVec.Set(point)
	tmpVec.Sub(position)
	self.Translate(tmpVec)
	self.Rotate(axis, angle)
	tmpVec.Rotate(axis, angle)
	self.Translate(-tmpVec.x, -tmpVec.y, -tmpVec.z)
}

// Transform the position, direction and up vector by the given matrix
// param transform The transform matrix
func (self *Camera) Transform(transform *Matrix4) {
	self.Position.Mul(transform)
	self.Rotate(transform)
}

// Moves the camera by the given amount on each axis.
// param x the displacement on the x-axis
// param y the displacement on the y-axis
// param z the displacement on the z-axis
func (self *Camera) Translate(x, y, z float32) {
	self.Position.Add(x, y, z)
}

// Moves the camera by the given vector.
// param vec the displacement vector
func (self *Camera) TranslateV3(vec *Vector3) {
	self.Position.Add(vec)
}

// Orthographic
// Rotates the camera by the given angle around the direction vector. The direction and up vector will not be orthogonalized.
// param angle
// Sets this camera to an orthographic projection using a viewport fitting the screen resolution, centered at
// (Gdx.graphics.getWidth()/2, Gdx.graphics.getHeight()/2), with the y-axis pointing up or down.
// param yDown whether y should be pointing down
func (self *OrthographicCamera) SetToOrtho(yDown bool) {
	self.SetToOrthoVW(yDown, Gdx.graphics.getWidth(), Gdx.graphics.getHeight())
}

// Sets this camera to an orthographic projection, centered at (viewportWidth/2, viewportHeight/2), with the y-axis pointing up
// or down.
// param yDown whether y should be pointing down.
// param viewportWidth
// param viewportHeight
func (self *OrthographicCamera) SetToOrthoVW(yDown bool, viewportWidth, viewportHeight float32) {
	if yDown {
		up.set(0, -1, 0)
		direction.set(0, 0, 1)
	} else {
		up.set(0, 1, 0)
		direction.set(0, 0, -1)
	}
	position.Set(zoom*viewportWidth/2.0, zoom*viewportHeight/2.0, 0)
	this.viewportWidth = viewportWidth
	this.viewportHeight = viewportHeight
	update()
}

func (self *Camera) RotateAngle(angle float32) {
	self.Rotate(self.Direction, angle)
}

// Moves the camera by the given amount on each axis.
// param x the displacement on the x-axis
// param y the displacement on the y-axis
func (self *Camera) TranslateXY(x, y float32) {
	self.Translate(x, y, 0)
}

// Moves the camera by the given vector.
// param vec the displacement vector
func (self *Camera) TranslateV2(vec *Vector2) {
	self.Translate(vec.X, vec.Y, 0)
}

// Function to translate a point given in screen coordinates to world space. It's the same as GLU gluUnProject, but does not
// rely on OpenGL. The x- and y-coordinate of vec are assumed to be in screen coordinates (origin is the top left corner, y
// pointing down, x pointing to the right) as reported by the touch methods in {@link Input}. A z-coordinate of 0 will return a
// point on the near plane, a z-coordinate of 1 will return a point on the far plane. This method allows you to specify the
// viewport position and dimensions in the coordinate system expected by {@link GL20#glViewport(int, int, int, int)}, with the
// origin in the bottom left corner of the screen.
// param screenCoords the point in screen coordinates (origin top left)
// param viewportX the coordinate of the bottom left corner of the viewport in glViewport coordinates.
// param viewportY the coordinate of the bottom left corner of the viewport in glViewport coordinates.
// param viewportWidth the width of the viewport in pixels
// param viewportHeight the height of the viewport in pixels
func (self *Camera) Unproject(screenCoords *Vector3, viewportX, viewportY, viewportWidth, viewportHeight float32) *Vector3 {
	x := screenCoords.X
	y := screenCoords.Y
	x = x - viewportX
	y = Gdx.graphics.getHeight() - y - 1
	y = y - viewportY
	screenCoords.X = (2*x)/viewportWidth - 1
	screenCoords.Y = (2*y)/viewportHeight - 1
	screenCoords.Z = 2*screenCoords.Z - 1
	screenCoords.Prj(invProjectionView)
	return screenCoords
}

// Function to translate a point given in screen coordinates to world space. It's the same as GLU gluUnProject but does not
// rely on OpenGL. The viewport is assumed to span the whole screen and is fetched from {@link Graphics#getWidth()} and
// {@link Graphics#getHeight()}. The x- and y-coordinate of vec are assumed to be in screen coordinates (origin is the top left
// corner, y pointing down, x pointing to the right) as reported by the touch methods in {@link Input}. A z-coordinate of 0
// will return a point on the near plane, a z-coordinate of 1 will return a point on the far plane.
// param screenCoords the point in screen coordinates
func (self *Camera) UnprojectV3(screenCoords *Vector3) *Vector3 {
	self.Unproject(screenCoords, 0, 0, Gdx.graphics.getWidth(), Gdx.graphics.getHeight())
	return screenCoords
}

// Projects the {@link Vector3} given in world space to screen coordinates. It's the same as GLU gluProject with one small
// deviation: The viewport is assumed to span the whole screen. The screen coordinate system has its origin in the
// <b>bottom</b> left, with the y-axis pointing <b>upwards</b> and the x-axis pointing to the right. This makes it easily
// useable in conjunction with {@link Batch} and similar classes.
func (self *Camera) ProjectV3(worldCoords *Vector3) *Vector3 {
	self.Project(worldCoords, 0, 0, Gdx.graphics.getWidth(), Gdx.graphics.getHeight())
	return worldCoords
}

// Projects the {@link Vector3} given in world space to screen coordinates. It's the same as GLU gluProject with one small
// deviation: The viewport is assumed to span the whole screen. The screen coordinate system has its origin in the
// <b>bottom</b> left, with the y-axis pointing <b>upwards</b> and the x-axis pointing to the right. This makes it easily
// useable in conjunction with {@link Batch} and similar classes. This method allows you to specify the viewport position and
// dimensions in the coordinate system expected by {@link GL20#glViewport(int, int, int, int)}, with the origin in the bottom
// left corner of the screen.
// param viewportX the coordinate of the bottom left corner of the viewport in glViewport coordinates.
// param viewportY the coordinate of the bottom left corner of the viewport in glViewport coordinates.
// param viewportWidth the width of the viewport in pixels
// param viewportHeight the height of the viewport in pixels
func (self *Camera) Project(worldCoords *Vector3, viewportX, viewportY, viewportWidth, viewportHeight float32) *Vector3 {
	worldCoords.Prj(combined)
	worldCoords.X = viewportWidth*(worldCoords.X+1)/2 + viewportX
	worldCoords.Y = viewportHeight*(worldCoords.Y+1)/2 + viewportY
	worldCoords.Z = (worldCoords.Z + 1) / 2
	return worldCoords
}

// Creates a picking {@link Ray} from the coordinates given in screen coordinates. It is assumed that the viewport spans the
// whole screen. The screen coordinates origin is assumed to be in the top left corner, its y-axis pointing down, the x-axis
// pointing to the right. The returned instance is not a new instance but an internal member only accessible via this function.
// param viewportX the coordinate of the bottom left corner of the viewport in glViewport coordinates.
// param viewportY the coordinate of the bottom left corner of the viewport in glViewport coordinates.
// param viewportWidth the width of the viewport in pixels
// param viewportHeight the height of the viewport in pixels
// return the picking Ray.
func (self *Camera) GetPickRay(screenX, screenY, viewportX, viewportY, viewportWidth, viewportHeight float32) *Ray {
	self.Unproject(ray.origin.Set(screenX, screenY, 0), viewportX, viewportY, viewportWidth, viewportHeight)
	self.Unproject(ray.direction.Set(screenX, screenY, 1), viewportX, viewportY, viewportWidth, viewportHeight)
	ray.direction.SubV(ray.origin).Nor()
	return ray
}

// Creates a picking {@link Ray} from the coordinates given in screen coordinates. It is assumed that the viewport spans the
// whole screen. The screen coordinates origin is assumed to be in the top left corner, its y-axis pointing down, the x-axis
// pointing to the right. The returned instance is not a new instance but an internal member only accessible via this function.
// return the picking Ray.
func (self *Camera) GetPickRay(screenX, screenY float32) *Ray {
	return self.GetPickRay(screenX, screenY, 0, 0, Gdx.graphics.getWidth(), Gdx.graphics.getHeight())
}
