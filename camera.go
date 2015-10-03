// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"math"

	. "github.com/pyros2097/spike/math"
	. "github.com/pyros2097/spike/math/collision"
	. "github.com/pyros2097/spike/math/vector"
	// "github.com/pyros2097/spike/utils/scaling"
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
	frustum *Frustum

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
		ViewportWidth:     viewportWidth,
		ViewportHeight:    viewportHeight,
		tmpVec:            NewVector3Empty(),
		tmp:               NewVector3Empty(),
		Ray:               NewRay(NewVector3Empty(), NewVector3Empty()),
		Zoom:              1,
		FieldOfView:       67,
		frustum:           NewFrustumEmpty(),
	}
}

func NewPerspectiveCamera(fieldOfViewY, viewportWidth, viewportHeight float32) *Camera {
	self := NewCamera(viewportWidth, viewportHeight)
	self.FieldOfView = fieldOfViewY
	self.Update()
	return self
}

// Recalculates the projection and view matrix of this camera and the {@link Frustum} planes. Use this after you've manipulated
// any of the attributes of the camera.
func (self *Camera) Update() {
	self.UpdateFrumstum(true)
}

// Recalculates the projection and view matrix of this camera and the {@link Frustum} planes if <code>updateFrustum</code> is
// true. Use this after you've manipulated any of the attributes of the camera.
func (self *Camera) UpdateFrumstum(updateFrustum bool) {}

func (self *Camera) UpdateFrustumPerspective(updateFrustum bool) {
	aspect := self.ViewportWidth / self.ViewportHeight
	self.Projection.SetToProjectionNear(float32(math.Abs(float64(self.Near))), float32(math.Abs(float64(self.Far))),
		self.FieldOfView, aspect) // TODO: check this call overload
	self.View.SetToLookAtPos(self.Position, self.tmp.SetV(self.Position).AddV(self.Direction), self.Up)
	self.Combined.SetM4(self.Projection)
	// Matrix4.mul(combined.val, view.val)

	if updateFrustum {
		self.InvProjectionView.SetM4(self.Combined)
		// Matrix4.inv(invProjectionView.val)
		self.frustum.Update(self.InvProjectionView)
	}
}

func (self *Camera) UpdateFrustumOrtho(updateFrustum bool) {
	self.Projection.SetToOrtho2DNear(self.Zoom*-self.ViewportWidth/2, self.Zoom*(self.ViewportWidth/2), self.Zoom*-(self.ViewportHeight/2),
		self.Zoom*self.ViewportHeight/2, self.Near, self.Far)
	self.View.SetToLookAtPos(self.Position, self.tmp.SetV(self.Position).AddV(self.Direction), self.Up)
	self.Combined.SetM4(self.Projection)
	// Matrix4.MulM4(combined.val, view.val)

	if updateFrustum {
		self.InvProjectionView.SetM4(self.Combined)
		// Matrix4.inv(invProjectionView.val)
		self.frustum.Update(self.InvProjectionView)
	}
}

// Recalculates the direction of the camera to look at the point (x, y, z). This function assumes the up vector is normalized.
// param x the x-coordinate of the point to look at
// param y the x-coordinate of the point to look at
// param z the x-coordinate of the point to look at
func (self *Camera) LookAt(x, y, z float32) {
	self.tmpVec.Set(x, y, z).SubV(self.Position).Nor()
	if !self.tmpVec.IsZero() {
		dot := self.tmpVec.DotV(self.Up) // up and direction must ALWAYS be orthonormal vectors
		if math.Abs(float64(dot-1)) < 0.000000001 {
			// Collinear
			self.Up.SetV(self.Direction).SclScalar(-1)
		} else if math.Abs(float64(dot+1)) < 0.000000001 {
			// Collinear opposite
			self.Up.SetV(self.Direction)
		}
		self.Direction.SetV(self.tmpVec)
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
	self.tmpVec.SetV(self.Direction).CrsV(self.Up).Nor()
	self.Up.SetV(self.tmpVec).CrsV(self.Direction).Nor()
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
	self.Direction.RotateV(axis, angle)
	self.Up.RotateV(axis, angle)
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
	quat.Transform(self.Direction)
	quat.Transform(self.Up)
}

// Rotates the direction and up vector of this camera by the given angle around the given axis, with the axis attached to given
// point. The direction and up vector will not be orthogonalized.
// param point the point to attach the axis to
// param axis the axis to rotate around
// param angle the angle
func (self *Camera) RotateAround(point, axis *Vector3, angle float32) {
	self.tmpVec.SetV(point)
	self.tmpVec.SubV(self.Position)
	self.TranslateV3(self.tmpVec)
	self.RotateV3(axis, angle)
	self.tmpVec.RotateV(axis, angle)
	self.Translate(-self.tmpVec.X, -self.tmpVec.Y, -self.tmpVec.Z)
}

// Transform the position, direction and up vector by the given matrix
// param transform The transform matrix
func (self *Camera) Transform(transform *Matrix4) {
	self.Position.Mul(transform)
	self.RotateM4(transform)
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
	self.Position.AddV(vec)
}

// Orthographic
// Rotates the camera by the given angle around the direction vector. The direction and up vector will not be orthogonalized.
// param angle
// Sets this camera to an orthographic projection using a viewport fitting the screen resolution, centered at
// (Gdx.graphics.getWidth()/2, Gdx.graphics.getHeight()/2), with the y-axis pointing up or down.
// param yDown whether y should be pointing down
func (self *Camera) SetToOrtho(yDown bool) {
	// Gdx.graphics.getWidth() TODO:
	// Gdx.graphics.getHeight()
	self.SetToOrthoVW(yDown, 0, 0)
}

// Sets this camera to an orthographic projection, centered at (viewportWidth/2, viewportHeight/2), with the y-axis pointing up
// or down.
// param yDown whether y should be pointing down.
// param viewportWidth
// param viewportHeight
func (self *Camera) SetToOrthoVW(yDown bool, viewportWidth, viewportHeight float32) {
	if yDown {
		self.Up.Set(0, -1, 0)
		self.Direction.Set(0, 0, 1)
	} else {
		self.Up.Set(0, 1, 0)
		self.Direction.Set(0, 0, -1)
	}
	self.Position.Set(self.Zoom*viewportWidth/2.0, self.Zoom*viewportHeight/2.0, 0)
	self.ViewportWidth = viewportWidth
	self.ViewportHeight = viewportHeight
	self.Update()
}

func (self *Camera) RotateAngle(angle float32) {
	self.RotateV3(self.Direction, angle)
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
	y = 0 //Gdx.graphics.getHeight() - y - 1 TODO:
	y = y - viewportY
	screenCoords.X = (2*x)/viewportWidth - 1
	screenCoords.Y = (2*y)/viewportHeight - 1
	screenCoords.Z = 2*screenCoords.Z - 1
	screenCoords.Prj(self.InvProjectionView)
	return screenCoords
}

// Function to translate a point given in screen coordinates to world space. It's the same as GLU gluUnProject but does not
// rely on OpenGL. The viewport is assumed to span the whole screen and is fetched from {@link Graphics#getWidth()} and
// {@link Graphics#getHeight()}. The x- and y-coordinate of vec are assumed to be in screen coordinates (origin is the top left
// corner, y pointing down, x pointing to the right) as reported by the touch methods in {@link Input}. A z-coordinate of 0
// will return a point on the near plane, a z-coordinate of 1 will return a point on the far plane.
// param screenCoords the point in screen coordinates
func (self *Camera) UnprojectV3(screenCoords *Vector3) *Vector3 {
	// self.Unproject(screenCoords, 0, 0, Gdx.graphics.getWidth(), Gdx.graphics.getHeight()) TODO:
	return screenCoords
}

// Projects the {@link Vector3} given in world space to screen coordinates. It's the same as GLU gluProject with one small
// deviation: The viewport is assumed to span the whole screen. The screen coordinate system has its origin in the
// <b>bottom</b> left, with the y-axis pointing <b>upwards</b> and the x-axis pointing to the right. This makes it easily
// useable in conjunction with {@link Batch} and similar classes.
func (self *Camera) ProjectV3(worldCoords *Vector3) *Vector3 {
	// self.Project(worldCoords, 0, 0, Gdx.graphics.getWidth(), Gdx.graphics.getHeight()) TODO:
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
	worldCoords.Prj(self.Combined)
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
	self.Unproject(self.Ray.Origin.Set(screenX, screenY, 0), viewportX, viewportY, viewportWidth, viewportHeight)
	self.Unproject(self.Ray.Direction.Set(screenX, screenY, 1), viewportX, viewportY, viewportWidth, viewportHeight)
	self.Ray.Direction.SubV(self.Ray.Origin).Nor()
	return self.Ray
}

// Creates a picking {@link Ray} from the coordinates given in screen coordinates. It is assumed that the viewport spans the
// whole screen. The screen coordinates origin is assumed to be in the top left corner, its y-axis pointing down, the x-axis
// pointing to the right. The returned instance is not a new instance but an internal member only accessible via this function.
// return the picking Ray.
func (self *Camera) GetPickRayXY(screenX, screenY float32) *Ray {
	//Gdx.graphics.getWidth(), Gdx.graphics.getHeight()
	return self.GetPickRay(screenX, screenY, 0, 0, 0, 0) // TODO
}

// A ScalingViewport that uses {@link Scaling#fit} so it keeps the aspect ratio by scaling the world up to fit the screen, adding
// black bars (letterboxing) for the remaining space.
// public class FitViewport extends ScalingViewport {
// 	/** Creates a new viewport using a new {@link OrthographicCamera}. */
// 	public FitViewport (float worldWidth, float worldHeight) {
// 		super(Scaling.fit, worldWidth, worldHeight);
// 	}

// 	public FitViewport (float worldWidth, float worldHeight, Camera camera) {
// 		super(Scaling.fit, worldWidth, worldHeight, camera);
// 	}
// }

// A viewport where the world size is based on the size of the screen. By default 1 world unit == 1 screen pixel, but this ratio
// can be {@link #setUnitsPerPixel(float) changed}.
// public class ScreenViewport extends Viewport {
// 	private float unitsPerPixel = 1;

// 	/** Creates a new viewport using a new {@link OrthographicCamera}. */
// 	public ScreenViewport () {
// 		this(new OrthographicCamera());
// 	}

// 	public ScreenViewport (Camera camera) {
// 		setCamera(camera);
// 	}

// 	@Override
// 	public void update (int screenWidth, int screenHeight, boolean centerCamera) {
// 		setScreenBounds(0, 0, screenWidth, screenHeight);
// 		setWorldSize(screenWidth * unitsPerPixel, screenHeight * unitsPerPixel);
// 		apply(centerCamera);
// 	}

// 	public float getUnitsPerPixel () {
// 		return unitsPerPixel;
// 	}

// 	/** Sets the number of pixels for each world unit. Eg, a scale of 2.5 means there are 2.5 world units for every 1 screen pixel.
// 	 * Default is 1. */
// 	public void setUnitsPerPixel (float unitsPerPixel) {
// 		this.unitsPerPixel = unitsPerPixel;
// 	}
// }

// A ScalingViewport that uses {@link Scaling#stretch} so it does not keep the aspect ratio, the world is scaled to take the whole
// screen.
// public class StretchViewport extends ScalingViewport {
// 	/** Creates a new viewport using a new {@link OrthographicCamera}. */
// 	public StretchViewport (float worldWidth, float worldHeight) {
// 		super(Scaling.stretch, worldWidth, worldHeight);
// 	}

// 	public StretchViewport (float worldWidth, float worldHeight, Camera camera) {
// 		super(Scaling.stretch, worldWidth, worldHeight, camera);
// 	}
// }

// A viewport that scales the world using {@link Scaling}.
// {@link Scaling#fit} keeps the aspect ratio by scaling the world up to fit the screen, adding black bars (letterboxing) for the
// remaining space.
// {@link Scaling#fill} keeps the aspect ratio by scaling the world up to take the whole screen (some of the world may be off
// screen).
// {@link Scaling#stretch} does not keep the aspect ratio, the world is scaled to take the whole screen.
// {@link Scaling#none} keeps the aspect ratio by using a fixed size world (the world may not fill the screen or some of the world
//may be off screen).
// public class ScalingViewport extends Viewport {
// 	private Scaling scaling;

// 	/** Creates a new viewport using a new {@link OrthographicCamera}. */
// 	public ScalingViewport (Scaling scaling, float worldWidth, float worldHeight) {
// 		this(scaling, worldWidth, worldHeight, new OrthographicCamera());
// 	}

// 	public ScalingViewport (Scaling scaling, float worldWidth, float worldHeight, Camera camera) {
// 		this.scaling = scaling;
// 		setWorldSize(worldWidth, worldHeight);
// 		setCamera(camera);
// 	}

// 	@Override
// 	public void update (int screenWidth, int screenHeight, boolean centerCamera) {
// 		Vector2 scaled = scaling.apply(getWorldWidth(), getWorldHeight(), screenWidth, screenHeight);
// 		int viewportWidth = Math.round(scaled.x);
// 		int viewportHeight = Math.round(scaled.y);

// 		// Center.
// 		setScreenBounds((screenWidth - viewportWidth) / 2, (screenHeight - viewportHeight) / 2, viewportWidth, viewportHeight);

// 		apply(centerCamera);
// 	}

// 	public Scaling getScaling () {
// 		return scaling;
// 	}

// 	public void setScaling (Scaling scaling) {
// 		this.scaling = scaling;
// 	}
// }

// A ScalingViewport that uses {@link Scaling#fill} so it keeps the aspect ratio by scaling the world up to take the whole screen
// (some of the world may be off screen).
// public class FillViewport extends ScalingViewport {
// 	/** Creates a new viewport using a new {@link OrthographicCamera}. */
// 	public FillViewport (float worldWidth, float worldHeight) {
// 		super(Scaling.fill, worldWidth, worldHeight);
// 	}

// 	public FillViewport (float worldWidth, float worldHeight, Camera camera) {
// 		super(Scaling.fill, worldWidth, worldHeight, camera);
// 	}
// }

// A viewport that keeps the world aspect ratio by extending the world in one direction. The world is first scaled to fit within
// the viewport, then the shorter dimension is lengthened to fill the viewport. A maximum size can be specified to limit how much
// the world is extended and black bars (letterboxing) are used for any remaining space.
// public class ExtendViewport extends Viewport {
// 	private float minWorldWidth, minWorldHeight;
// 	private float maxWorldWidth, maxWorldHeight;

// 	/** Creates a new viewport using a new {@link OrthographicCamera} with no maximum world size. */
// 	public ExtendViewport (float minWorldWidth, float minWorldHeight) {
// 		this(minWorldWidth, minWorldHeight, 0, 0, new OrthographicCamera());
// 	}

// 	/** Creates a new viewport with no maximum world size. */
// 	public ExtendViewport (float minWorldWidth, float minWorldHeight, Camera camera) {
// 		this(minWorldWidth, minWorldHeight, 0, 0, camera);
// 	}

// 	/** Creates a new viewport using a new {@link OrthographicCamera} and a maximum world size.
// 	 * @see ExtendViewport#ExtendViewport(float, float, float, float, Camera) */
// 	public ExtendViewport (float minWorldWidth, float minWorldHeight, float maxWorldWidth, float maxWorldHeight) {
// 		this(minWorldWidth, minWorldHeight, maxWorldWidth, maxWorldHeight, new OrthographicCamera());
// 	}

// 	/** Creates a new viewport with a maximum world size.
// 	 * @param maxWorldWidth User 0 for no maximum width.
// 	 * @param maxWorldHeight User 0 for no maximum height. */
// 	public ExtendViewport (float minWorldWidth, float minWorldHeight, float maxWorldWidth, float maxWorldHeight, Camera camera) {
// 		this.minWorldWidth = minWorldWidth;
// 		this.minWorldHeight = minWorldHeight;
// 		this.maxWorldWidth = maxWorldWidth;
// 		this.maxWorldHeight = maxWorldHeight;
// 		setCamera(camera);
// 	}

// 	@Override
// 	public void update (int screenWidth, int screenHeight, boolean centerCamera) {
// 		// Fit min size to the screen.
// 		float worldWidth = minWorldWidth;
// 		float worldHeight = minWorldHeight;
// 		Vector2 scaled = Scaling.fit.apply(worldWidth, worldHeight, screenWidth, screenHeight);

// 		// Extend in the short direction.
// 		int viewportWidth = Math.round(scaled.x);
// 		int viewportHeight = Math.round(scaled.y);
// 		if (viewportWidth < screenWidth) {
// 			float toViewportSpace = viewportHeight / worldHeight;
// 			float toWorldSpace = worldHeight / viewportHeight;
// 			float lengthen = (screenWidth - viewportWidth) * toWorldSpace;
// 			if (maxWorldWidth > 0) lengthen = Math.min(lengthen, maxWorldWidth - minWorldWidth);
// 			worldWidth += lengthen;
// 			viewportWidth += Math.round(lengthen * toViewportSpace);
// 		} else if (viewportHeight < screenHeight) {
// 			float toViewportSpace = viewportWidth / worldWidth;
// 			float toWorldSpace = worldWidth / viewportWidth;
// 			float lengthen = (screenHeight - viewportHeight) * toWorldSpace;
// 			if (maxWorldHeight > 0) lengthen = Math.min(lengthen, maxWorldHeight - minWorldHeight);
// 			worldHeight += lengthen;
// 			viewportHeight += Math.round(lengthen * toViewportSpace);
// 		}

// 		setWorldSize(worldWidth, worldHeight);

// 		// Center.
// 		setScreenBounds((screenWidth - viewportWidth) / 2, (screenHeight - viewportHeight) / 2, viewportWidth, viewportHeight);

// 		apply(centerCamera);
// 	}

// 	public float getMinWorldWidth () {
// 		return minWorldWidth;
// 	}

// 	public void setMinWorldWidth (float minWorldWidth) {
// 		this.minWorldWidth = minWorldWidth;
// 	}

// 	public float getMinWorldHeight () {
// 		return minWorldHeight;
// 	}

// 	public void setMinWorldHeight (float minWorldHeight) {
// 		this.minWorldHeight = minWorldHeight;
// 	}

// 	public float getMaxWorldWidth () {
// 		return maxWorldWidth;
// 	}

// 	public void setMaxWorldWidth (float maxWorldWidth) {
// 		this.maxWorldWidth = maxWorldWidth;
// 	}

// 	public float getMaxWorldHeight () {
// 		return maxWorldHeight;
// 	}

// 	public void setMaxWorldHeight (float maxWorldHeight) {
// 		this.maxWorldHeight = maxWorldHeight;
// 	}
// }

// Manages a {@link Camera} and determines how world coordinates are mapped to and from the screen.
// public abstract class Viewport {
// 	private Camera camera;
// 	private float worldWidth, worldHeight;
// 	private int screenX, screenY, screenWidth, screenHeight;

// 	private final Vector3 tmp = new Vector3();

// 	/** Calls {@link #apply(boolean)} with false. */
// 	public void apply () {
// 		apply(false);
// 	}

// 	/** Applies the viewport to the camera and sets the glViewport.
// 	 * @param centerCamera If true, the camera position is set to the center of the world. */
// 	public void apply (boolean centerCamera) {
// 		Gdx.gl.glViewport(screenX, screenY, screenWidth, screenHeight);
// 		camera.viewportWidth = worldWidth;
// 		camera.viewportHeight = worldHeight;
// 		if (centerCamera) camera.position.set(worldWidth / 2, worldHeight / 2, 0);
// 		camera.update();
// 	}

// 	/** Calls {@link #update(int, int, boolean)} with false. */
// 	public final void update (int screenWidth, int screenHeight) {
// 		update(screenWidth, screenHeight, false);
// 	}

// 	/** Configures this viewport's screen bounds using the specified screen size and calls {@link #apply(boolean)}. Typically called
// 	 * from {@link ApplicationListener#resize(int, int)} or {@link Screen#resize(int, int)}.
// 	 * <p>
// 	 * The default implementation only calls {@link #apply(boolean)}. */
// 	public void update (int screenWidth, int screenHeight, boolean centerCamera) {
// 		apply(centerCamera);
// 	}

// 	/** Transforms the specified screen coordinate to world coordinates.
// 	 * @return The vector that was passed in, transformed to world coordinates.
// 	 * @see Camera#unproject(Vector3) */
// 	public Vector2 unproject (Vector2 screenCoords) {
// 		tmp.set(screenCoords.x, screenCoords.y, 1);
// 		camera.unproject(tmp, screenX, screenY, screenWidth, screenHeight);
// 		screenCoords.set(tmp.x, tmp.y);
// 		return screenCoords;
// 	}

// 	/** Transforms the specified world coordinate to screen coordinates.
// 	 * @return The vector that was passed in, transformed to screen coordinates.
// 	 * @see Camera#project(Vector3) */
// 	public Vector2 project (Vector2 worldCoords) {
// 		tmp.set(worldCoords.x, worldCoords.y, 1);
// 		camera.project(tmp, screenX, screenY, screenWidth, screenHeight);
// 		worldCoords.set(tmp.x, tmp.y);
// 		return worldCoords;
// 	}

// 	/** Transforms the specified screen coordinate to world coordinates.
// 	 * @return The vector that was passed in, transformed to world coordinates.
// 	 * @see Camera#unproject(Vector3) */
// 	public Vector3 unproject (Vector3 screenCoords) {
// 		camera.unproject(screenCoords, screenX, screenY, screenWidth, screenHeight);
// 		return screenCoords;
// 	}

// 	/** Transforms the specified world coordinate to screen coordinates.
// 	 * @return The vector that was passed in, transformed to screen coordinates.
// 	 * @see Camera#project(Vector3) */
// 	public Vector3 project (Vector3 worldCoords) {
// 		camera.project(worldCoords, screenX, screenY, screenWidth, screenHeight);
// 		return worldCoords;
// 	}

// 	/** @see Camera#getPickRay(float, float, float, float, float, float) */
// 	public Ray getPickRay (float screenX, float screenY) {
// 		return camera.getPickRay(screenX, screenY, this.screenX, this.screenY, screenWidth, screenHeight);
// 	}

// 	/** @see ScissorStack#calculateScissors(Camera, float, float, float, float, Matrix4, Rectangle, Rectangle) */
// 	public void calculateScissors (Matrix4 batchTransform, Rectangle area, Rectangle scissor) {
// 		ScissorStack.calculateScissors(camera, screenX, screenY, screenWidth, screenHeight, batchTransform, area, scissor);
// 	}

// 	/** Transforms a point to real screen coordinates (as opposed to OpenGL ES window coordinates), where the origin is in the top
// 	 * left and the the y-axis is pointing downwards. */
// 	public Vector2 toScreenCoordinates (Vector2 worldCoords, Matrix4 transformMatrix) {
// 		tmp.set(worldCoords.x, worldCoords.y, 0);
// 		tmp.mul(transformMatrix);
// 		camera.project(tmp);
// 		tmp.y = Gdx.graphics.getHeight() - tmp.y;
// 		worldCoords.x = tmp.x;
// 		worldCoords.y = tmp.y;
// 		return worldCoords;
// 	}

// 	public Camera getCamera () {
// 		return camera;
// 	}

// 	public void setCamera (Camera camera) {
// 		this.camera = camera;
// 	}

// 	public float getWorldWidth () {
// 		return worldWidth;
// 	}

// 	/** The virtual width of this viewport in world coordinates. This width is scaled to the viewport's screen width. */
// 	public void setWorldWidth (float worldWidth) {
// 		this.worldWidth = worldWidth;
// 	}

// 	public float getWorldHeight () {
// 		return worldHeight;
// 	}

// 	/** The virtual height of this viewport in world coordinates. This height is scaled to the viewport's screen height. */
// 	public void setWorldHeight (float worldHeight) {
// 		this.worldHeight = worldHeight;
// 	}

// 	public void setWorldSize (float worldWidth, float worldHeight) {
// 		this.worldWidth = worldWidth;
// 		this.worldHeight = worldHeight;
// 	}

// 	public int getScreenX () {
// 		return screenX;
// 	}

// 	/** Sets the viewport's offset from the left edge of the screen. This is typically set by {@link #update(int, int, boolean)}. */
// 	public void setScreenX (int screenX) {
// 		this.screenX = screenX;
// 	}

// 	public int getScreenY () {
// 		return screenY;
// 	}

// 	/** Sets the viewport's offset from the bottom edge of the screen. This is typically set by {@link #update(int, int, boolean)}. */
// 	public void setScreenY (int screenY) {
// 		this.screenY = screenY;
// 	}

// 	public int getScreenWidth () {
// 		return screenWidth;
// 	}

// 	/** Sets the viewport's width in screen coordinates. This is typically set by {@link #update(int, int, boolean)}. */
// 	public void setScreenWidth (int screenWidth) {
// 		this.screenWidth = screenWidth;
// 	}

// 	public int getScreenHeight () {
// 		return screenHeight;
// 	}

// 	/** Sets the viewport's height in screen coordinates. This is typically set by {@link #update(int, int, boolean)}. */
// 	public void setScreenHeight (int screenHeight) {
// 		this.screenHeight = screenHeight;
// 	}

// 	/** Sets the viewport's position in screen coordinates. This is typically set by {@link #update(int, int, boolean)}. */
// 	public void setScreenPosition (int screenX, int screenY) {
// 		this.screenX = screenX;
// 		this.screenY = screenY;
// 	}

// 	/** Sets the viewport's size in screen coordinates. This is typically set by {@link #update(int, int, boolean)}. */
// 	public void setScreenSize (int screenWidth, int screenHeight) {
// 		this.screenWidth = screenWidth;
// 		this.screenHeight = screenHeight;
// 	}

// 	/** Sets the viewport's bounds in screen coordinates. This is typically set by {@link #update(int, int, boolean)}. */
// 	public void setScreenBounds (int screenX, int screenY, int screenWidth, int screenHeight) {
// 		this.screenX = screenX;
// 		this.screenY = screenY;
// 		this.screenWidth = screenWidth;
// 		this.screenHeight = screenHeight;
// 	}

// 	/** Returns the left gutter (black bar) width in screen coordinates. */
// 	public int getLeftGutterWidth () {
// 		return screenX;
// 	}

// 	/** Returns the right gutter (black bar) x in screen coordinates. */
// 	public int getRightGutterX () {
// 		return screenX + screenWidth;
// 	}

// 	/** Returns the right gutter (black bar) width in screen coordinates. */
// 	public int getRightGutterWidth () {
// 		return Gdx.graphics.getWidth() - (screenX + screenWidth);
// 	}

// 	/** Returns the bottom gutter (black bar) height in screen coordinates. */
// 	public int getBottomGutterHeight () {
// 		return screenY;
// 	}

// 	/** Returns the top gutter (black bar) y in screen coordinates. */
// 	public int getTopGutterY () {
// 		return screenY + screenHeight;
// 	}

// 	/** Returns the top gutter (black bar) height in screen coordinates. */
// 	public int getTopGutterHeight () {
// 		return Gdx.graphics.getHeight() - (screenY + screenHeight);
// 	}
// }

// package scene2d;
// import com.badlogic.gdx.Gdx;
// import com.badlogic.gdx.graphics.OrthographicCamera;
// import com.badlogic.gdx.math.Interpolation;
// import com.badlogic.gdx.math.Vector3;
// import com.badlogic.gdx.scenes.scene2d.Actor;
// import com.badlogic.gdx.utils.Array;
// import com.badlogic.gdx.math.Rectangle;

// public class Camera extends OrthographicCamera {
// 	private static Camera instance;
// 	private static float duration;
// 	private static float time;
// 	private static Interpolation interpolation;
// 	private static boolean moveCompleted;
// 	private static float lastPercent;
// 	private static float percentDelta;
// 	private static float panSpeedX, panSpeedY;
// 	private static final Array<Actor> hudActors = new Array<Actor>();

// 	private static Actor followedActor = null;
// 	private static float followSpeed = 1f;
// 	/*
// 	 *  This is to set the offsets of camera position when following the actor
// 	 *  When the camera follows the actor its (x,y) position is set to actor's (x,y) position
// 	 *  based on followSpeed. The offsets are used to position the camera in such a way that the actor
// 	 *  doesn't need to be at the center of the camera always
// 	 */
// 	public static Rectangle followOffset = new Rectangle(10,70,10,60);
// 	private static boolean followContinous = false;

// 	public static boolean usePan = false;
//     public static boolean useDrag = false;
// 	/*
// 	 * Sets the speed at which the camera pans. By default it moves 1px for a duration a 1f
// 	 * so its speed is 1px/f. So reduce the duration to increase its speed.
// 	 * ex: setPanSpeed(0.5) will change its speed to 2px/f
// 	 * Here: f can/maybe also indicate seconds
// 	 */
// 	public static float panSpeed = 5f;
// 	public static Rectangle panSize;

// 	/*
// 	 *  This sets the boundary of the camera till what position can it move or pan in the
// 	 *  directions left, right, top, down. This is to prevent is from panning overboard the game area.
// 	 *  Usually the bounds of the camera is like a rectangle. This must be calculated carefully
// 	 *  as the camera's position is based on its center.
// 	*/
// 	public static Rectangle bounds = new Rectangle(0,0,999,999);

// 	Camera(){
// 		setToOrtho(false, Scene.targetWidth, Scene.targetHeight);
// 		position.set(Scene.targetWidth/2, Scene.targetHeight/2, 0);
// 		instance = this;
// 		panSize = new Rectangle(10, 10, Scene.targetWidth-10, Scene.targetHeight - 10);
// 	}

// 	/*
// 	 * Moves the camera to x,y over a time duration
// 	 */
// 	public void moveTo(float x, float y, float duration) {
// 		moveBy(x-position.x, y-position.y, duration);
// 	}

// 	/*
// 	 * Moves the camera by amountX, amountY over a time duration
// 	 */
// 	public static void moveBy (float amountX, float amountY, float duration) {
// 		moveBy(amountX, amountY, duration, null);
// 	}

// 	/*
// 	 * Moves the camera by amountX, amountY over a time duration and interpolation interp
// 	 */
// 	public static void moveBy (float amountX, float amountY, float dur, Interpolation interp) {
// 		duration = dur;
// 		interpolation = interp;
// 		panSpeedX = amountX;
// 		panSpeedY = amountY;
// 		lastPercent = 0;
// 		time = 0;
// 		moveCompleted = false;
// 	}

// 	private static Rectangle cullRect = new Rectangle();
// 	private void moveByAction(float delta){
// 		time += delta;
// 		moveCompleted = time >= duration;
// 		float percent;
// 		if (moveCompleted)
// 			percent = 1;
// 		else {
// 			percent = time / duration;
// 			if (interpolation != null) percent = interpolation.apply(percent);
// 		}
// 		percentDelta = percent - lastPercent;
// 		if(Scene.cullingEnabled){
// 			cullRect.set(getXLeft(), getYBot(), Scene.targetWidth, Scene.targetHeight);
// 			Scene.getCurrentScene().setCullingArea(cullRect);
// 		}
// 		translate(panSpeedX * percentDelta, panSpeedY * percentDelta, 0);
// 		for(Actor actor: hudActors)
// 			actor.setPosition(actor.getX()+panSpeedX * percentDelta, actor.getY()+panSpeedY * percentDelta);
// 		lastPercent = percent;
// 		if (moveCompleted) interpolation = null;
// 	}

// 	public void resetCamera(){
// 		position.set(Scene.targetWidth/2, Scene.targetHeight/2, 0);
// 	}

// 	/*
// 	 * This makes the camera follow the actor once and only once. Once the camera reaches its
// 	 * target, it stops following the actor.
// 	 */
// 	public static void followActor(Actor actor){
// 		followedActor = actor;
// 		followContinous = false;
// 	}

// 	/*
// 	 * This makes the camera follow the actor continuously, even after the camera reaches its
// 	 * target, it keeps following the if the actor changes its position.
// 	 */
// 	public static void followActorContinuously(Actor actor){
// 		followedActor = actor;
// 		followContinous = true;
// 	}

// 	 * Sets the speed at which the camera follows the actor. By default it moves 1px for a duration of 1f
// 	 * so its speed is 1px/f. So reduce the duration to increase its speed.
// 	 * ex: setPanSpeed(0.5) will change its speed to 2px/f
// 	 * Here: f can/maybe also indicate seconds

// 	public static void setFollowSpeed(float duration){
// 		followSpeed = duration;
// 	}

// 	private void follow(){
// 		//if(camera.position.x == followedActor.getX()+followLeftOffset &&
// 		//	camera.position.y == followedActor.getY()+followTopOffset)
// 		//return;
// 		//moveTo(followedActor.getX()+followLeftOffset, followedActor.getY()+followTopOffset, 100f);
// 		if(position.x < followedActor.getX() - followOffset.x) moveBy(1f, 0, followSpeed);
// 		else if(position.x > followedActor.getX() + followOffset.width) moveBy(-1f, 0, followSpeed);
// 		else if(position.y < followedActor.getY() - followOffset.y) moveBy(0, 1f, followSpeed);
// 		else if(position.y > followedActor.getY() - followOffset.height) moveBy(0, -1f, followSpeed);
// 		else {
// 			if(!followContinous)
// 				followedActor = null;
// 		}
// 	}

// 	@Override
// 	public void update(){
// 		super.update();
// 		Scene.mouse.x = Gdx.input.getX();
// 		Scene.mouse.y = Gdx.graphics.getHeight() - Gdx.input.getY();
// 		if(!moveCompleted)
// 			moveByAction(Scene.stateTime);//FIXME
// 		if(usePan)
// 			panCameraWithMouse();
// 		if(followedActor != null)
// 			follow();
// 	}

// 	private void panCameraWithMouse(){
// 		 if(Scene.mouse.x > panSize.width && getXLeft() < bounds.width) moveBy(1f, 0, panSpeed);
// 		 else if(Scene.mouse.x < panSize.x  && getXLeft() > bounds.x)  moveBy(-1f, 0, panSpeed);
// 		 else if(Scene.mouse.y < panSize.y && getYBot() > bounds.y) moveBy(0, -1f, panSpeed);
// 		 else if(Scene.mouse.y > panSize.height && getYBot() < bounds.height) moveBy(0, 1f, panSpeed);
// 	}

// 	private final static Vector3 curr = new Vector3();
// 	private final static Vector3 last = new Vector3(-1, -1, -1);
// 	private final static Vector3 deltaDrag = new Vector3();
// 	private static float deltaCamX = 0;
// 	private static float deltaCamY = 0;
// 	public static void dragCam(int x, int y){
// 		instance.unproject(curr.set(x, y, 0));
//     	if (!(last.x == -1 && last.y == -1 && last.z == -1)) {
//     		instance.unproject(deltaDrag.set(last.x, last.y, 0));
//     		deltaDrag.sub(curr);
//     		deltaCamX = deltaDrag.x + instance.position.x;
//     		deltaCamY = deltaDrag.y + instance.position.y;
//     		if(deltaCamX > bounds.x && deltaCamX < bounds.x+bounds.width)
//     			moveBy(deltaDrag.x, 0f, 0f);
//     		if(deltaCamY > bounds.y && deltaCamY < bounds.y+bounds.height)
//     			moveBy(0f, deltaDrag.y, 0f);
//     	}
//     	last.set(x, y, 0);
//     }

// 	public static void resetDrag(){
// 		last.set(-1, -1, -1);
// 	}

// 	public static void reset(){
// 		usePan = false;
// 		followActor(null);
// 		clearAllHud();
// 		instance.position.set(Scene.targetWidth/2, Scene.targetHeight/2, 0);
// 	}

// 	/* If you want to make any elements/actors to move along with the camera
// 	 * like HUD's add them using this method */
// 	public static void addHud(Actor actor){
// 		if(actor != null){
// 			Scene.getCurrentScene().addActor(actor);
// 			hudActors.add(actor);
// 		}
// 	}

// 	/* If you want to make any elements/actors to move along with the camera
// 	 * like HUD's add them using this method */
// 	public static void addHud(String actorName){
// 		if(actorName != null && !actorName.isEmpty()){
// 			Actor actor = Scene.getCurrentScene().findActor(actorName);
// 			Scene.getCurrentScene().addActor(actor);
// 			hudActors.add(actor);
// 		}
// 	}

// 	/* If you want to any elements/actors which was a Hud the use this */
// 	public static void removeHud(Actor actor){
// 		Scene.getCurrentScene().removeActor(actor);
// 		hudActors.removeValue(actor, true);
// 	}

// 	/* If you want to any elements/actors which was a Hud the use this */
// 	public static void removeHud(String actorName){
// 		if(actorName != null && !actorName.isEmpty()){
// 			Actor actor = Scene.getCurrentScene().findActor(actorName);
// 			Scene.getCurrentScene().removeActor(actor);
// 			hudActors.removeValue(actor, true);
// 		}
// 	}

// 	/*
// 	 * Clears all hud's registered with the camera
// 	 */
// 	public static void clearAllHud(){
// 		hudActors.clear();
// 	}

// 	/*
// 	 * Returns the x position of the camera
// 	 */
// 	public static float getX(){
// 		return instance.position.x;
// 	}

// 	/*
// 	 * Returns the y position of the camera
// 	 */
// 	public static float getY(){
// 		return instance.position.y;
// 	}

// 	public static float getXLeft(){
// 		return instance.position.x - Scene.targetWidth/2;
// 	}

// 	public static float getYBot(){
// 		return instance.position.y - Scene.targetHeight/2;
// 	}

// 	public static float getWidth(){
// 		return instance.viewportWidth;
// 	}

// 	public static float getHeight(){
// 		return instance.viewportHeight;
// 	}
// }
