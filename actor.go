// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"github.com/pyros2097/spike/g2d"
	"github.com/pyros2097/spike/graphics"
	"github.com/pyros2097/spike/math/shape"
	"github.com/pyros2097/spike/utils"
)

// Touchable state of Actors
// Determines how touch input events are distributed to an actor and any children.
type Touchable int

const (
	// All touch input events will be received by the actor and any children.
	TouchableEnabled Touchable = iota

	// No touch input events will be received by the actor or any children.
	TouchableDisabled

	// No touch input events will be received by the actor, but children will still receive events. Note that events on the
	// children will still bubble to the parent.
	TouchableChildrenOnly
)

// 2D scene graph node. An actor has a position, rectangular size, origin, scale, rotation, Z index, and color. The position
// corresponds to the unrotated, unscaled bottom left corner of the actor. The position is relative to the actor's parent. The
// origin is relative to the position and is used for scale and rotation.
//
// An actor has a list of in progress {@link Action actions} that are applied to the actor (often over time). These are generally
// used to change the presentation of the actor (moving it, resizing it, etc). See {@link #act(float)}, {@link Action} and its
// many subclasses.
//
// An actor has two kinds of listeners associated with it: "capture" and regular. The listeners are notified of events the actor
// or its children receive. The regular listeners are designed to allow an actor to respond to events that have been delivered.
// The capture listeners are designed to allow a parent or container actor to handle events before child actors. See {@link #fire}
// for more details.
//
// An {@link InputListener} can receive all the basic input events. More complex listeners (like {@link ClickListener} and
// {@link ActorGestureListener}) can listen for and combine primitive events and recognize complex interactions like multi-touch
// or pinch.
type IActor interface {
	Draw(batch *g2d.Batch, parentAlpha float32)
	Act(delta float32)
	RemoveActor()
	AddAction(action Action)
	RemoveAction(action Action)
	GetActions() []*Action
	HasActions() bool
	ClearActions()
}

type Actor struct {
	Name     string  // Name of the Actor
	X, Y     float32 // X, Y position of the actor's left edge.
	W, H     float32 // Width and Height
	Z        uint32  // zindex
	OX, OY   float32 // origin
	SX, SY   float32 // scale
	Rotation float32

	// If false, the actor will not be drawn and will not receive touch events. Default is true.
	Visible bool

	Debug bool

	Actions []*Action

	// Called by the framework when an actor is added to or removed from a group.
	// param parent May be null if the actor has been removed from the parent.
	Parent *Actor

	// Parent, FirstChild, LastChild, PrevSibling, NextSibling *Actor

	// Determines how touch events are distributed to self actor. Default is {@link Touchable#enabled}.
	TouchState Touchable

	// An application specific object for convenience.
	UserObject interface{}

	// The color the actor will be tinted when drawn.
	Color *graphics.Color // make this 1 1 1 1

	// Called on the update cycle
	OnAct func(self *Actor, delta float32)

	// Draws the actor. The batch is configured to draw in the parent's coordinate system.
	// {@link Batch#draw(com.badlogic.gdx.graphics.g2d.TextureRegion, float, float, float, float, float, float, float, float, float)
	// This draw method} is convenient to draw a rotated and scaled TextureRegion. {@link Batch#begin()} has already been called on
	// the batch. If {@link Batch#end()} is called to draw without the batch then {@link Batch#begin()} must be called before the
	// method returns.
	// <p>
	// The default implementation does nothing.
	// param parentAlpha Should be multiplied with the actor's alpha, allowing a parent's alpha to affect all children.
	OnDraw func(self *Actor, batch g2d.Batch, parentAlpha float32)

	// Called when an input event occurrs
	OnInput func(self *Actor, event InputEvent)
}

// Updates the actor based on time. Typically this is called each frame by {@link Stage#act(float)}.
// The default implementation calls {@link Action#act(float)} on each action and removes actions that are complete.
// param delta Time in seconds since the last frame.
func (self *Actor) act(delta float32) {
	if len(self.Actions) > 0 {
		for i := 0; i < len(self.Actions); i++ {
			action := self.Actions[i]
			if action.Act(delta) && i < len(self.Actions) {
				current := self.Actions[i] // implement fast Array
				actionIndex := i
				if current != action {
					// actionIndex = actions.indexOf(action, true)
				}
				if actionIndex != -1 {
					// self.actions.removeIndex(actionIndex)
					action.SetActor(nil)
					i--
				}
			}
		}
	}
	if self.OnAct != nil {
		self.OnAct(self, delta)
	}
}

func (self *Actor) draw(batch g2d.Batch, parentAlpha float32) {
	if self.OnDraw != nil {
		self.OnDraw(self, batch, parentAlpha)
	}
}

// Removes self actor from its parent, if it has a parent.
func (self *Actor) RemoveActor() bool {
	if self.Parent == nil {
		// self.parent.removeActor(self)
		return false
	}
	return true
}

func (self *Actor) AddAction(action Action) {
	action.SetActor(self)
	// actions.add(action)
	// if (stage != null && stage.getActionsRequestRendering()) Gdx.graphics.requestRendering();
}

func (self *Actor) RemoveAction(action Action) {
	// if (actions.removeValue(action, true)) action.setActor(null)
}

// Returns true if the actor has one or more actions.
func (self *Actor) HasActions() bool {
	return len(self.Actions) > 0
}

// Removes all actions on self actor.
func (self *Actor) ClearActions() {
	for i := 0; i < len(self.Actions); i++ {
		// actions.get(i).setActor(null);
	}
	// actions.clear()
}

// Returns true if the actor's parent is not null.
func (self *Actor) HasParent() bool {
	return self.Parent != nil
}

// Returns true if input events are processed by self actor.
func (self *Actor) IsTouchable() bool {
	return self.TouchState == TouchableEnabled
}

// Returns the X position of the specified {@link Align alignment}
func (self *Actor) GetXAlign(alignment utils.Alignment) float32 {
	x := self.X
	if (alignment & utils.AlignmentRight) != 0 {
		x += self.W
	} else if (alignment & utils.AlignmentLeft) == 0 {
		x += self.W / 2
	}
	return x
}

func (self *Actor) SetX(x float32) {
	if self.X != x {
		self.X = x
		self.positionChanged()
	}
}

// Returns the Y position of the specified {@link Align alignment}
func (self *Actor) GetYAlign(alignment utils.Alignment) float32 {
	y := self.Y
	if (alignment & utils.AlignmentTop) != 0 {
		y += self.H
	} else if (alignment & utils.AlignmentBottom) == 0 {
		y += self.H / 2
	}
	return y
}

func (self *Actor) SetY(y float32) {
	if self.Y != y {
		self.Y = y
		self.positionChanged()
	}
}

// Sets the position of the actor's bottom left corner.
func (self *Actor) SetPosition(x, y float32) {
	if self.X != x || self.Y != y {
		self.X = x
		self.Y = y
		self.positionChanged()
	}
}

// Sets the position using the specified {@link Align alignment}.
// Note this may set the position to non-integer coordinates.
func (self *Actor) SetPositionAlign(x, y float32, alignment utils.Alignment) {
	if (alignment & utils.AlignmentRight) != 0 {
		self.X -= self.W
	} else if (alignment & utils.AlignmentLeft) == 0 {
		self.X -= self.W / 2
	}
	if (alignment & utils.AlignmentTop) != 0 {
		self.Y -= self.H
	} else if (alignment & utils.AlignmentBottom) == 0 {
		self.Y -= self.H / 2
	}
	if self.X != x || self.Y != y {
		self.X = x
		self.Y = y
		self.positionChanged()
	}
}

// Add x and y to current position.
func (self *Actor) MoveBy(x, y float32) {
	if x != 0 || y != 0 {
		self.X += x
		self.Y += y
		self.positionChanged()
	}
}

func (self *Actor) SetWidth(width float32) {
	oldWidth := self.W
	self.W = width
	if width != oldWidth {
		self.sizeChanged()
	}
}

func (self *Actor) SetHeight(height float32) {
	oldHeight := self.H
	self.H = height
	if height != oldHeight {
		self.sizeChanged()
	}
}

// Returns y plus height.
func (self *Actor) GetTop() float32 {
	return self.Y + self.H
}

//  Returns x plus width.
func (self *Actor) GetRight() float32 {
	return self.X + self.W
}

// Called when the actor's position has been changed.
func (self *Actor) positionChanged() {
}

// Called when the actor's size has been changed.
func (self *Actor) sizeChanged() {
}

// Sets the width and height.
func (self *Actor) setSize(w, h float32) {
	oldWidth := self.W
	oldHeight := self.H
	self.W = w
	self.H = h
	if w != oldWidth || h != oldHeight {
		self.sizeChanged()
	}
}

// Adds the specified size to the current size.
func (self *Actor) SizeBy(w, h float32) {
	self.W += w
	self.H += h
	self.sizeChanged()
}

// Set bounds the x, y, width, and height.
func (self *Actor) SetBounds(x, y, w, h float32) {
	if self.X != x || self.Y != y {
		self.X = x
		self.Y = y
		self.positionChanged()
	}
	if self.W != w || self.H != h {
		self.W = h
		self.H = h
		self.sizeChanged()
	}
}

//Sets the origin position which is relative to the actor's bottom left corner.
func (self *Actor) setOrigin(originX, originY float32) {
	self.OX = originX
	self.OY = originY
}

// Sets the origin position to the specified {@link Align alignment}.
func (self *Actor) SetOriginAlign(alignment utils.Alignment) {
	if (alignment & utils.AlignmentLeft) != 0 {
		self.OX = 0
	} else if (alignment & utils.AlignmentRight) != 0 {
		self.OX = self.W
	} else {
		self.OX = self.W / 2
	}

	if (alignment & utils.AlignmentBottom) != 0 {
		self.OY = 0
	} else if (alignment & utils.AlignmentTop) != 0 {
		self.OY = self.H
	} else {
		self.OY = self.H / 2
	}
}

// Sets the scale X and scale Y.
func (self *Actor) SetScale(scaleX, scaleY float32) {
	self.SX = scaleX
	self.SY = scaleY
}

// Adds the specified scale to the current scale.
func (self *Actor) ScaleBy(scaleX, scaleY float32) {
	self.SX += scaleX
	self.SY += scaleY
}

func (self *Actor) GetRotation() float32 {
	return self.Rotation
}

func (self *Actor) SetRotation(degrees float32) {
	self.Rotation = degrees
}

// Adds the specified rotation to the current rotation.
func (self *Actor) RotateBy(amountInDegrees float32) {
	self.Rotation += amountInDegrees
}

// Changes the z-order for self actor so it is in front of all siblings.
func (self *Actor) ToFront() {
	// setZIndex(Integer.MAX_VALUE)
}

// Changes the z-order for self actor so it is in back of all siblings.
func (self *Actor) ToBack() {
	// setZIndex(0)
}

// Sets the z-index of self actor. The z-index is the index into the parent's {@link Group#getChildren() children}, where a
// lower index is below a higher index. Setting a z-index higher than the number of children will move the child to the front.
// Setting a z-index less than zero is invalid.
func (self *Actor) SetZIndex(index uint32) {
	// Group parent = self.parent;
	// if (parent == null) return;
	// Array<Actor> children = parent.children;
	// if (children.size == 1) return;
	// if (!children.removeValue(self, true)) return;
	// if (index >= children.size)
	//   children.add(self);
	// else
	//   children.insert(index, self);
}

// Returns the z-index of self actor.
func (self *Actor) GetZIndex() uint32 {
	return 0
	// Group parent = self.parent;
	// if (parent == null) return -1;
	// return parent.children.indexOf(self, true);
}

//   // Sets self actor as the event {@link Event#setTarget(Actor) target} and propagates the event to self actor and ancestor
//    * actors as necessary. If self actor is not in the stage, the stage must be set before calling self method.
//    * <p>
//    * Events are fired in 2 phases.
//    * <ol>
//    * <li>The first phase (the "capture" phase) notifies listeners on each actor starting at the root and propagating downward to
//    * (and including) self actor.</li>
//    * <li>The second phase notifies listeners on each actor starting at self actor and, if {@link Event#getBubbles()} is true,
//    * propagating upward to the root.</li>
//    * </ol>
//    * If the event is {@link Event#stop() stopped} at any time, it will not propagate to the next actor.
//    * @return true if the event was {@link Event#cancel() cancelled}.
//   public boolean fire (Event event) {
//     if (event.getStage() == null) event.setStage(getStage());
//     event.setTarget(self);

//     // Collect ancestors so event propagation is unaffected by hierarchy changes.
//     Array<Group> ancestors = Pools.obtain(Array.class);
//     Group parent = self.parent;
//     while (parent != null) {
//       ancestors.add(parent);
//       parent = parent.parent;
//     }

//     try {
//       // Notify all parent capture listeners, starting at the root. Ancestors may stop an event before children receive it.
//       Object[] ancestorsArray = ancestors.items;
//       for (int i = ancestors.size - 1; i >= 0; i--) {
//         Group currentTarget = (Group)ancestorsArray[i];
//         currentTarget.notify(event, true);
//         if (event.isStopped()) return event.isCancelled();
//       }

//       // Notify the target capture listeners.
//       notify(event, true);
//       if (event.isStopped()) return event.isCancelled();

//       // Notify the target listeners.
//       notify(event, false);
//       if (!event.getBubbles()) return event.isCancelled();
//       if (event.isStopped()) return event.isCancelled();

//       // Notify all parent listeners, starting at the target. Children may stop an event before ancestors receive it.
//       for (int i = 0, n = ancestors.size; i < n; i++) {
//         ((Group)ancestorsArray[i]).notify(event, false);
//         if (event.isStopped()) return event.isCancelled();
//       }

//       return event.isCancelled();
//     } finally {
//       ancestors.clear();
//       Pools.free(ancestors);
//     }
//   }

//   // Notifies self actor's listeners of the event. The event is not propagated to any parents. Before notifying the listeners,
//    * self actor is set as the {@link Event#getListenerActor() listener actor}. The event {@link Event#setTarget(Actor) target}
//    * must be set before calling self method. If self actor is not in the stage, the stage must be set before calling self method.
//    * param capture If true, the capture listeners will be notified instead of the regular listeners.
//    * @return true of the event was {@link Event#cancel() cancelled}.
//   public boolean notify (Event event, boolean capture) {
//     if (event.getTarget() == null) throw new IllegalArgumentException("The event target cannot be null.");

//     DelayedRemovalArray<EventListener> listeners = capture ? captureListeners : self.listeners;
//     if (listeners.size == 0) return event.isCancelled();

//     event.setListenerActor(self);
//     event.setCapture(capture);
//     if (event.getStage() == null) event.setStage(stage);

//     listeners.begin();
//     for (int i = 0, n = listeners.size; i < n; i++) {
//       EventListener listener = listeners.get(i);
//       if (listener.handle(event)) {
//         event.handle();
//         if (event instanceof InputEvent) {
//           InputEvent inputEvent = (InputEvent)event;
//           if (inputEvent.getType() == Type.touchDown) {
//             event.getStage().addTouchFocus(listener, self, inputEvent.getTarget(), inputEvent.getPointer(),
//               inputEvent.getButton());
//           }
//         }
//       }
//     }
//     listeners.end();

//     return event.isCancelled();
//   }

//   // Returns the deepest actor that contains the specified point and is {@link #getTouchable() touchable} and
//    * {@link #isVisible() visible}, or null if no actor was hit. The point is specified in the actor's local coordinate system (0,0
//    * is the bottom left of the actor and width,height is the upper right).
//    * <p>
//    * This method is used to delegate touchDown, mouse, and enter/exit events. If self method returns null, those events will not
//    * occur on self Actor.
//    * <p>
//    * The default implementation returns self actor if the point is within self actor's bounds.
//    *
//    * param touchable If true, the hit detection will respect the {@link #setTouchable(Touchable) touchability}.
//    * @see Touchable
//   public Actor hit (float x, float y, boolean touchable) {
//     if (touchable && self.touchable != Touchable.enabled) return null;
//     return x >= 0 && x < width && y >= 0 && y < height ? self : null;
//   }

//   // Add a listener to receive events that {@link #hit(float, float, boolean) hit} self actor. See {@link #fire(Event)}.
//    *
//    * @see InputListener
//    * @see ClickListener
//   public boolean addListener (EventListener listener) {
//     if (!listeners.contains(listener, true)) {
//       listeners.add(listener);
//       return true;
//     }
//     return false;
//   }

//   public boolean removeListener (EventListener listener) {
//     return listeners.removeValue(listener, true);
//   }

//   public Array<EventListener> getListeners () {
//     return listeners;
//   }

//   // Adds a listener that is only notified during the capture phase.
//    * @see #fire(Event)
//   public boolean addCaptureListener (EventListener listener) {
//     if (!captureListeners.contains(listener, true)) captureListeners.add(listener);
//     return true;
//   }

//   public boolean removeCaptureListener (EventListener listener) {
//     return captureListeners.removeValue(listener, true);
//   }

//   public Array<EventListener> getCaptureListeners () {
//     return captureListeners;
//   }

//   // Returns true if self actor is the same as or is the descendant of the specified actor.
//   public boolean isDescendantOf (Actor actor) {
//     if (actor == null) throw new IllegalArgumentException("actor cannot be null.");
//     Actor parent = self;
//     while (true) {
//       if (parent == null) return false;
//       if (parent == actor) return true;
//       parent = parent.parent;
//     }
//   }

//   // Returns true if self actor is the same as or is the ascendant of the specified actor.
//   public boolean isAscendantOf (Actor actor) {
//     if (actor == null) throw new IllegalArgumentException("actor cannot be null.");
//     while (true) {
//       if (actor == null) return false;
//       if (actor == self) return true;
//       actor = actor.parent;
//     }
//   }

//   // Calls {@link #clipBegin(float, float, float, float)} to clip self actor's bounds.
//   public boolean clipBegin () {
//     return clipBegin(x, y, width, height);
//   }

//   // Clips the specified screen aligned rectangle, specified relative to the transform matrix of the stage's Batch. The transform
//    * matrix and the stage's camera must not have rotational components. Calling self method must be followed by a call to
//    * {@link #clipEnd()} if true is returned.
//    * @return false if the clipping area is zero and no drawing should occur.
//    * @see ScissorStack
//   public boolean clipBegin (float x, float y, float width, float height) {
//     if (width <= 0 || height <= 0) return false;
//     Rectangle tableBounds = Rectangle.tmp;
//     tableBounds.x = x;
//     tableBounds.y = y;
//     tableBounds.width = width;
//     tableBounds.height = height;
//     Stage stage = self.stage;
//     Rectangle scissorBounds = Pools.obtain(Rectangle.class);
//     stage.calculateScissors(tableBounds, scissorBounds);
//     if (ScissorStack.pushScissors(scissorBounds)) return true;
//     Pools.free(scissorBounds);
//     return false;
//   }

//   // Ends clipping begun by {@link #clipBegin(float, float, float, float)}.
//   public void clipEnd () {
//     Pools.free(ScissorStack.popScissors());
//   }

//   // Transforms the specified point in screen coordinates to the actor's local coordinate system.
//   public Vector2 screenToLocalCoordinates (Vector2 screenCoords) {
//     Stage stage = self.stage;
//     if (stage == null) return screenCoords;
//     return stageToLocalCoordinates(stage.screenToStageCoordinates(screenCoords));
//   }

//   // Transforms the specified point in the stage's coordinates to the actor's local coordinate system.
//   public Vector2 stageToLocalCoordinates (Vector2 stageCoords) {
//     if (parent != null) parent.stageToLocalCoordinates(stageCoords);
//     parentToLocalCoordinates(stageCoords);
//     return stageCoords;
//   }

//   // Transforms the specified point in the actor's coordinates to be in the stage's coordinates.
//    * @see Stage#toScreenCoordinates(Vector2, com.badlogic.gdx.math.Matrix4)
//   public Vector2 localToStageCoordinates (Vector2 localCoords) {
//     return localToAscendantCoordinates(null, localCoords);
//   }

//   // Transforms the specified point in the actor's coordinates to be in the parent's coordinates.
//   public Vector2 localToParentCoordinates (Vector2 localCoords) {
//     final float rotation = -self.rotation;
//     final float scaleX = self.scaleX;
//     final float scaleY = self.scaleY;
//     final float x = self.x;
//     final float y = self.y;
//     if (rotation == 0) {
//       if (scaleX == 1 && scaleY == 1) {
//         localCoords.x += x;
//         localCoords.y += y;
//       } else {
//         final float originX = self.originX;
//         final float originY = self.originY;
//         localCoords.x = (localCoords.x - originX) * scaleX + originX + x;
//         localCoords.y = (localCoords.y - originY) * scaleY + originY + y;
//       }
//     } else {
//       final float cos = (float)Math.cos(rotation * MathUtils.degreesToRadians);
//       final float sin = (float)Math.sin(rotation * MathUtils.degreesToRadians);
//       final float originX = self.originX;
//       final float originY = self.originY;
//       final float tox = (localCoords.x - originX) * scaleX;
//       final float toy = (localCoords.y - originY) * scaleY;
//       localCoords.x = (tox * cos + toy * sin) + originX + x;
//       localCoords.y = (tox * -sin + toy * cos) + originY + y;
//     }
//     return localCoords;
//   }

//   // Converts coordinates for self actor to those of a parent actor. The ascendant does not need to be a direct parent.
//   public Vector2 localToAscendantCoordinates (Actor ascendant, Vector2 localCoords) {
//     Actor actor = self;
//     while (actor != null) {
//       actor.localToParentCoordinates(localCoords);
//       actor = actor.parent;
//       if (actor == ascendant) break;
//     }
//     return localCoords;
//   }

//   // Converts the coordinates given in the parent's coordinate system to self actor's coordinate system.
//   public Vector2 parentToLocalCoordinates (Vector2 parentCoords) {
//     final float rotation = self.rotation;
//     final float scaleX = self.scaleX;
//     final float scaleY = self.scaleY;
//     final float childX = x;
//     final float childY = y;
//     if (rotation == 0) {
//       if (scaleX == 1 && scaleY == 1) {
//         parentCoords.x -= childX;
//         parentCoords.y -= childY;
//       } else {
//         final float originX = self.originX;
//         final float originY = self.originY;
//         parentCoords.x = (parentCoords.x - childX - originX) / scaleX + originX;
//         parentCoords.y = (parentCoords.y - childY - originY) / scaleY + originY;
//       }
//     } else {
//       final float cos = (float)Math.cos(rotation * MathUtils.degreesToRadians);
//       final float sin = (float)Math.sin(rotation * MathUtils.degreesToRadians);
//       final float originX = self.originX;
//       final float originY = self.originY;
//       final float tox = parentCoords.x - childX - originX;
//       final float toy = parentCoords.y - childY - originY;
//       parentCoords.x = (tox * cos + toy * sin) / scaleX + originX;
//       parentCoords.y = (tox * -sin + toy * cos) / scaleY + originY;
//     }
//     return parentCoords;
//   }

//   // Draws self actor's debug lines if {@link #getDebug()} is true.
//   public void drawDebug (ShapeRenderer shapes) {
//     drawDebugBounds(shapes);
//   }

//   // Draws a rectange for the bounds of self actor if {@link #getDebug()} is true.
//   protected void drawDebugBounds (ShapeRenderer shapes) {
//     if (!debug) return;
//     shapes.set(ShapeType.Line);
//     shapes.setColor(stage.getDebugColor());
//     shapes.rect(x, y, originX, originY, width, height, scaleX, scaleY, rotation);
//   }

func (self *Actor) GetBounds() *shape.Rectangle {
	return NewRectangle(self.X, self.Y, self.W, self.H)
}

func (self *Actor) CollidesXY(x, y float32) bool {
	if self.GetBounds().Overlaps(NewRectangle(x, y, 5, 5)) {
		return true
	}
	return false
}

func (self *Actor) Collides(other *Actor) bool {
	if self.GetBounds().Overlaps(other.GetBounds()) {
		return true
	}
	return false
}

func (self *Actor) String() string {
	return self.Name
}
