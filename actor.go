// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"math"
	"sync"

	"github.com/pyros2097/spike/g2d"
	. "github.com/pyros2097/spike/interpolation"
	"github.com/pyros2097/spike/math/shape"
	"github.com/pyros2097/spike/math/vector"
	"github.com/pyros2097/spike/utils"
)

// Touchable state of Actors
// Determines how touch input events are distributed to an actor and any children.
type Touchable int

const (
	// TouchableEnabled All touch input events will be received by the actor and any children.
	TouchableEnabled Touchable = iota

	// TouchableDisabled No touch input events will be received by the actor or any children.
	TouchableDisabled

	// TouchableChildrenOnly No touch input events will be received by the actor, but children will still receive events. Note that events on the
	// children will still bubble to the parent.
	TouchableChildrenOnly
)

// Action is used to modify the actor's parameters over time
type Action struct {
	// The actor this action is attached to, or nil if it is not attached.
	Actor *Actor

	// The actor this action targets, or nil if a target has not been set.
	// Sets the actor this action will manipulate. If no target actor is set, {@link #setActor(Actor)} will set the target actor
	// when the action is added to an actor.
	Target *Actor

	// Updates the action based on time. Typically this is called each frame by {@link Actor#act(float)}.
	// delta Time in seconds since the last frame.
	// returns true if the action is done. This method may continue to be called after the action is done.
	Act func(a *Action, delta float32) bool

	// Sets the state of the action so it can be run again.
	Restart func(a *Action)

	// Resets the optional state of this action to as if it were newly created, allowing the action to be pooled and reused. State
	// required to be set for every usage of this action or computed during the action does not need to be reset.
	// The default implementation calls {@link #restart()}.
	// If a subclass has optional state, it must override this method, call super, and reset the optional state. */
	Reset func(a *Action)
}

var actionsPool = sync.Pool{
	New: func() interface{} {
		return &Action{}
	},
}

/** Sets the actor this action is attached to. This also sets the {@link #setTarget(Actor) target} actor if it is nil. This
 * method is called automatically when an action is added to an actor. This method is also called with nil when an action is
 * removed from an actor.
 * <p>
 * When set to nil, if the action has a {@link #setPool(Pool) pool} then the action is {@link Pool#free(Object) returned} to
 * the pool (which calls {@link #reset()}) and the pool is set to nil. If the action does not have a pool, {@link #reset()} is
 * not called.
 * <p>
 * This method is not typically a good place for an action subclass to query the actor's state because the action may not be
 * executed for some time, eg it may be {@link DelayAction delayed}. The actor's state is best queried in the first call to
 * {@link #act(float)}. For a {@link TemporalAction}, use TemporalAction#begin(). */

// NewAction creates a new empty action
func NewAction() *Action {
	return &Action{
		Act: func(a *Action, delta float32) bool {
			return true
		},
		Restart: func(a *Action) {},
		Reset: func(a *Action) {
			a.Actor = nil
			a.Target = nil
			// a.pool = nil
			a.Restart(a)
		},
	}
}

// Actor is a 2D scene graph node that may contain other actors.
// Actors have a z-order equal to the order they were inserted.
// Actors inserted later will be drawn on top of actors added earlier.
// Touch events that hit more than one actor are distributed to topmost actors first.
// An actor has a position, rectangular size, origin, scale, rotation, Z index, and color.
// The position corresponds to the unrotated, unscaled bottom left corner of the actor.
// The position is relative to the actor's parent.
// The origin is relative to the position and is used for scale and rotation.
//
// An actor has a list of in progress {@link Action actions} that are applied to the actor (often over time).
// These are generally used to change the presentation of the actor (moving it, resizing it, etc).
//
// An actor has two kinds of listeners associated with it: "capture" and regular.
// The listeners are notified of events the actor or its children receive.
// The regular listeners are designed to allow an actor to respond to events that have been delivered.
// The capture listeners are designed to allow a parent or container actor to handle events before child actors. See {@link #fire}
// for more details.
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

	// Determines how touch events are distributed to a actor. Default is {@link Touchable#enabled}.
	TouchState Touchable

	// An application specific object for convenience.
	UserObject interface{}

	// The color the actor will be tinted when drawn.
	Color *Color // make this 1 1 1 1

	Init func(a *Actor)

	// Called on the update cycle
	Act func(a *Actor, delta float32)

	// Draws the actor. The batch is configured to draw in the parent's coordinate system.
	// {@link Batch#draw(com.badlogic.gdx.graphics.g2d.TextureRegion, float, float, float, float, float, float, float, float, float)
	// This draw method} is convenient to draw a rotated and scaled TextureRegion. {@link Batch#begin()} has already been called on
	// the batch. If {@link Batch#end()} is called to draw without the batch then {@link Batch#begin()} must be called before the
	// method returns.
	// <p>
	// The default implementation does nothing.
	// param parentAlpha Should be multiplied with the actor's alpha, allowing a parent's alpha to affect all children.
	Draw func(a *Actor, batch g2d.Batch, parentAlpha float32)

	// Called when an input event occurrs
	Input func(a *Actor, event InputEvent)

	Children                        []*Actor
	worldTransform                  *vector.Affine2
	computedTransform, oldTransform *vector.Matrix4
	transform                       bool
	cullingArea                     shape.Rectangle
	initialized                     bool
}

var (
	tmp = vector.NewVector2Empty()
)

// Draws the group and its children. The default implementation calls {@link #applyTransform(Batch, Matrix4)} if needed, then
// {@link #drawChildren(Batch, float)}, then {@link #resetTransform(Batch)} if needed.
// func (a *Group) draw(batch g2d.Batch, parentAlpha float32) {
// if (transform) applyTransform(batch, computeTransform());
// drawChildren(batch, parentAlpha);
// if (transform) resetTransform(batch);
// }

//  // Draws all children. {@link #applyTransform(Batch, Matrix4)} should be called before and {@link #resetTransform(Batch)} after
//   * this method if {@link #setTransform(boolean) transform} is true. If {@link #setTransform(boolean) transform} is false these
//   * methods don't need to be called, children positions are temporarily offset by the group position when drawn. This method
//   * avoids drawing children completely outside the {@link #setCullingArea(Rectangle) culling area}, if set.
//  protected func (a *Group) drawChildren (Batch batch, float parentAlpha) {
//    parentAlpha *= this.color.a;
//    SnapshotArray<Actor> children = this.children;
//    Actor[] actors = children.begin();
//    Rectangle cullingArea = this.cullingArea;
//    if (cullingArea != null) {
//      // Draw children only if inside culling area.
//      float cullLeft = cullingArea.x;
//      float cullRight = cullLeft + cullingArea.width;
//      float cullBottom = cullingArea.y;
//      float cullTop = cullBottom + cullingArea.height;
//      if (transform) {
//        for (int i = 0, n = children.size; i < n; i++) {
//          Actor child = actors[i];
//          if (!child.isVisible()) continue;
//          float cx = child.x, cy = child.y;
//          if (cx <= cullRight && cy <= cullTop && cx + child.width >= cullLeft && cy + child.height >= cullBottom)
//            child.draw(batch, parentAlpha);
//        }
//      } else {
//        // No transform for this group, offset each child.
//        float offsetX = x, offsetY = y;
//        x = 0;
//        y = 0;
//        for (int i = 0, n = children.size; i < n; i++) {
//          Actor child = actors[i];
//          if (!child.isVisible()) continue;
//          float cx = child.x, cy = child.y;
//          if (cx <= cullRight && cy <= cullTop && cx + child.width >= cullLeft && cy + child.height >= cullBottom) {
//            child.x = cx + offsetX;
//            child.y = cy + offsetY;
//            child.draw(batch, parentAlpha);
//            child.x = cx;
//            child.y = cy;
//          }
//        }
//        x = offsetX;
//        y = offsetY;
//      }
//    } else {
//      // No culling, draw all children.
//      if (transform) {
//        for (int i = 0, n = children.size; i < n; i++) {
//          Actor child = actors[i];
//          if (!child.isVisible()) continue;
//          child.draw(batch, parentAlpha);
//        }
//      } else {
//        // No transform for this group, offset each child.
//        float offsetX = x, offsetY = y;
//        x = 0;
//        y = 0;
//        for (int i = 0, n = children.size; i < n; i++) {
//          Actor child = actors[i];
//          if (!child.isVisible()) continue;
//          float cx = child.x, cy = child.y;
//          child.x = cx + offsetX;
//          child.y = cy + offsetY;
//          child.draw(batch, parentAlpha);
//          child.x = cx;
//          child.y = cy;
//        }
//        x = offsetX;
//        y = offsetY;
//      }
//    }
//    children.end();
//  }

//  // Draws this actor's debug lines if {@link #getDebug()} is true and, regardless of {@link #getDebug()}, calls
//   * {@link Actor#drawDebug(ShapeRenderer)} on each child.
// func (a *Group) drawDebug (ShapeRenderer shapes) {
//    drawDebugBounds(shapes);
//    if (transform) applyTransform(shapes, computeTransform());
//    drawDebugChildren(shapes);
//    if (transform) resetTransform(shapes);
//  }

//  // Draws all children. {@link #applyTransform(Batch, Matrix4)} should be called before and {@link #resetTransform(Batch)} after
//   * this method if {@link #setTransform(boolean) transform} is true. If {@link #setTransform(boolean) transform} is false these
//   * methods don't need to be called, children positions are temporarily offset by the group position when drawn. This method
//   * avoids drawing children completely outside the {@link #setCullingArea(Rectangle) culling area}, if set.
//  protected void drawDebugChildren (ShapeRenderer shapes) {
//    SnapshotArray<Actor> children = this.children;
//    Actor[] actors = children.begin();
//    // No culling, draw all children.
//    if (transform) {
//      for (int i = 0, n = children.size; i < n; i++) {
//        Actor child = actors[i];
//        if (!child.isVisible()) continue;
//        if (!child.getDebug() && !(child instanceof Group)) continue;
//        child.drawDebug(shapes);
//      }
//      shapes.flush();
//    } else {
//      // No transform for this group, offset each child.
//      float offsetX = x, offsetY = y;
//      x = 0;
//      y = 0;
//      for (int i = 0, n = children.size; i < n; i++) {
//        Actor child = actors[i];
//        if (!child.isVisible()) continue;
//        if (!child.getDebug() && !(child instanceof Group)) continue;
//        float cx = child.x, cy = child.y;
//        child.x = cx + offsetX;
//        child.y = cy + offsetY;
//        child.drawDebug(shapes);
//        child.x = cx;
//        child.y = cy;
//      }
//      x = offsetX;
//      y = offsetY;
//    }
//    children.end();
//  }

//  // Returns the transform for this group's coordinate system.
//  protected Matrix4 computeTransform () {
//    Affine2 worldTransform = this.worldTransform;

//    float originX = this.originX;
//    float originY = this.originY;
//    float rotation = this.rotation;
//    float scaleX = this.scaleX;
//    float scaleY = this.scaleY;

//    worldTransform.setToTrnRotScl(x + originX, y + originY, rotation, scaleX, scaleY);
//    if (originX != 0 || originY != 0) worldTransform.translate(-originX, -originY);

//    // Find the first parent that transforms.
//    Group parentGroup = parent;
//    while (parentGroup != null) {
//      if (parentGroup.transform) break;
//      parentGroup = parentGroup.parent;
//    }
//    if (parentGroup != null) worldTransform.preMul(parentGroup.worldTransform);

//    computedTransform.set(worldTransform);
//    return computedTransform;
//  }

//  // Set the batch's transformation matrix, often with the result of {@link #computeTransform()}. Note this causes the batch to
//   * be flushed. {@link #resetTransform(Batch)} will restore the transform to what it was before this call.
//  protected void applyTransform (Batch batch, Matrix4 transform) {
//    oldTransform.set(batch.getTransformMatrix());
//    batch.setTransformMatrix(transform);
//  }

//  // Restores the batch transform to what it was before {@link #applyTransform(Batch, Matrix4)}. Note this causes the batch to be
//   * flushed.
//  protected void resetTransform (Batch batch) {
//    batch.setTransformMatrix(oldTransform);
//  }

//  // Set the shape renderer transformation matrix, often with the result of {@link #computeTransform()}. Note this causes the
//   * shape renderer to be flushed. {@link #resetTransform(ShapeRenderer)} will restore the transform to what it was before this
//   * call.
//  protected void applyTransform (ShapeRenderer shapes, Matrix4 transform) {
//    oldTransform.set(shapes.getTransformMatrix());
//    shapes.setTransformMatrix(transform);
//  }

//  // Restores the shape renderer transform to what it was before {@link #applyTransform(Batch, Matrix4)}. Note this causes the
//   * shape renderer to be flushed.
//  protected void resetTransform (ShapeRenderer shapes) {
//    shapes.setTransformMatrix(oldTransform);
//  }

//  // Children completely outside of this rectangle will not be drawn. This is only valid for use with unrotated and unscaled
//   * actors!
// func (a *Group) setCullingArea (Rectangle cullingArea) {
//    this.cullingArea = cullingArea;
//  }

//  // @see #setCullingArea(Rectangle)
//  public Rectangle getCullingArea () {
//    return cullingArea;
//  }

//  // Returns the first actor found with the specified name. Note this recursively compares the name of every actor in the group.
//  public <T extends Actor> T findActor (String name) {
//    Array<Actor> children = this.children;
//    for (int i = 0, n = children.size; i < n; i++)
//      if (name.equals(children.get(i).getName())) return (T)children.get(i);
//    for (int i = 0, n = children.size; i < n; i++) {
//      Actor child = children.get(i);
//      if (child instanceof Group) {
//        Actor actor = ((Group)child).findActor(name);
//        if (actor != null) return (T)actor;
//      }
//    }
//    return null;
//  }

//  // Swaps two actors by index. Returns false if the swap did not occur because the indexes were out of bounds.
//  public boolean swapActor (int first, int second) {
//    int maxIndex = children.size;
//    if (first < 0 || first >= maxIndex) return false;
//    if (second < 0 || second >= maxIndex) return false;
//    children.swap(first, second);
//    return true;
//  }

//  // Swaps two actors. Returns false if the swap did not occur because the actors are not children of this group.
//  public boolean swapActor (Actor first, Actor second) {
//    int firstIndex = children.indexOf(first, true);
//    int secondIndex = children.indexOf(second, true);
//    if (firstIndex == -1 || secondIndex == -1) return false;
//    children.swap(firstIndex, secondIndex);
//    return true;
//  }

//  // Returns an ordered list of child actors in this group.
//  public SnapshotArray<Actor> getChildren () {
//    return children;
//  }

//  public boolean hasChildren () {
//    return children.size > 0;
//  }

//  // When true (the default), the Batch is transformed so children are drawn in their parent's coordinate system. This has a
//   * performance impact because {@link Batch#flush()} must be done before and after the transform. If the actors in a group are
//   * not rotated or scaled, then the transform for the group can be set to false. In this case, each child's position will be
//   * offset by the group's position for drawing, causing the children to appear in the correct location even though the Batch has
//   * not been transformed.
// func (a *Group) setTransform (boolean transform) {
//    this.transform = transform;
//  }

//  public boolean isTransform () {
//    return transform;
//  }

//  // Converts coordinates for this group to those of a descendant actor. The descendant does not need to be a direct child.
//   * @throws IllegalArgumentException if the specified actor is not a descendant of this group.
//  public Vector2 localToDescendantCoordinates (Actor descendant, Vector2 localCoords) {
//    Group parent = descendant.parent;
//    if (parent == null) throw new IllegalArgumentException("Child is not a descendant: " + descendant);
//    // First convert to the actor's parent coordinates.
//    if (parent != this) localToDescendantCoordinates(parent, localCoords);
//    // Then from each parent down to the descendant.
//    descendant.parentToLocalCoordinates(localCoords);
//    return localCoords;
//  }

//  // If true, {@link #drawDebug(ShapeRenderer)} will be called for this group and, optionally, all children recursively.
// func (a *Group) setDebug (boolean enabled, boolean recursively) {
//    setDebug(enabled);
//    if (recursively) {
//      for (Actor child : children) {
//        if (child instanceof Group) {
//          ((Group)child).setDebug(enabled, recursively);
//        } else {
//          child.setDebug(enabled);
//        }
//      }
//    }
//  }

//  // Calls {@link #setDebug(boolean, boolean)} with {@code true, true}.
//  public Group debugAll () {
//    setDebug(true, true);
//    return this;
//  }

//  // Returns a description of the actor hierarchy, recursively.
//  func (a *Group) toString() string {
//    StringBuilder buffer = new StringBuilder(128);
//    toString(buffer, 1);
//    buffer.setLength(buffer.length() - 1);
//    return buffer.toString();
//  }

//  void toString (StringBuilder buffer, int indent) {
//    buffer.append(super.toString());
//    buffer.append('\n');

//    Actor[] actors = children.begin();
//    for (int i = 0, n = children.size; i < n; i++) {
//      for (int ii = 0; ii < indent; ii++)
//        buffer.append("|  ");
//      Actor actor = actors[i];
//      if (actor instanceof Group)
//        ((Group)actor).toString(buffer, indent + 1);
//      else {
//        buffer.append(actor);
//        buffer.append('\n');
//      }
//    }
//    children.end();
//  }

// transform:         true,
// worldTransform:    vector.NewAffine2Empty(),
// oldTransform:      vector.NewMatrix4Empty(),
// computedTransform: vector.NewMatrix4Empty(),

// Updates the actor based on time. Typically this is called each frame by {@link Stage#act(float)}.
// The default implementation calls {@link Action#act(float)} on each action and removes actions that are complete.
// param delta Time in seconds since the last frame.
func (a *Actor) act(delta float32) {
	if !a.initialized && a.Init != nil {
		a.Init(a)
		a.initialized = true
	}
	if len(a.Actions) > 0 {
		action := a.Actions[0]
		if action.Act(action, delta) {
			action.Actor = nil
			action.Act = nil
			actionsPool.Put(action)
			a.Actions = a.Actions[1:len(a.Actions)]
		}
	}
	if a.Act != nil {
		a.Act(a, delta)
	}
	if a.Children != nil {
		for _, child := range a.Children {
			child.act(delta)
		}
	}
}

func (a *Actor) draw(batch g2d.Batch, parentAlpha float32) {
	if a.Draw != nil {
		a.Draw(a, batch, parentAlpha)
	}
}

func (a *Actor) AddAction(action *Action) {
	action.Actor = a
	a.Actions = append(a.Actions, action)
	// if (stage != null && stage.getActionsRequestRendering()) Gdx.graphics.requestRendering();
}

func (a *Actor) RemoveAction(action *Action) {
	for i, action := range a.Actions {
		if action == action {
			a.Actions, a.Actions[len(a.Actions)-1] = append(a.Actions[:i], a.Actions[i+1:]...), nil
			action.Actor = nil
			actionsPool.Put(action)
			break
		}
	}
}

// Returns true if the actor has one or more
func (a *Actor) HasActions() bool {
	return len(a.Actions) > 0
}

// Removes all actions on a actor.
func (a *Actor) ClearActions() {
	for _, action := range a.Actions {
		action.Actor = nil
		actionsPool.Put(action)
	}
	a.Actions = make([]*Action, 0)
}

// Called when actors are added to or removed from the group.
func (a *Actor) ChildrenChanged() {
}

// Adds an actor as a child of this group. The actor is first removed from its parent group, if any.
func (a *Actor) AddActor(actor *Actor) {
	if actor.Parent != nil {
		actor.Parent.RemoveActor(actor)
	}
	a.Children = append(a.Children, actor)
	actor.Parent = a
	a.ChildrenChanged()
}

// Adds an actor as a child of this group, at a specific index. The actor is first removed from its parent group, if any.
// @param index May be greater than the number of children.
func (a *Actor) AddActorAt(index int, actor *Actor) {
	if actor.Parent != nil {
		actor.Parent.RemoveActor(actor)
	}
	if index >= len(a.Children) {
		a.Children = append(a.Children, actor)
	} else {
		// children.insert(index, actor)
	}
	actor.Parent = a
	a.ChildrenChanged()
}

//  // Adds an actor as a child of this group, immediately before another child actor. The actor is first removed from its parent
//  // group, if any.
// func (a *Actor) addActorBefore (actorBefore, actor *Actor) {
//    if actor.Parent != nil {
// actor.Parent.RemoveActor(actor, false)
// }
//    int index = children.indexOf(actorBefore, true);
//    children.insert(index, actor);
//    actor.setParent(this);
//    actor.setStage(getStage());
//    a.ChildrenChanged()
//  }

//  // Adds an actor as a child of this group, immediately after another child actor. The actor is first removed from its parent
//   * group, if any.
// func (a *Group) addActorAfter (Actor actorAfter, Actor actor) {
//    if (actor.parent != null) actor.parent.removeActor(actor, false);
//    int index = children.indexOf(actorAfter, true);
//    if (index == children.size)
//      children.add(actor);
//    else
//      children.insert(index + 1, actor);
//    actor.setParent(this);
//    actor.setStage(getStage());
//    a.ChildrenChanged()
//  }

// Removes an actor from this group.
// If the actor will not be used again and has actions, they should be
// {@link Actor#clearActions() cleared} so the actions will be returned to their
// {@link Action#setPool(com.badlogic.gdx.utils.Pool) pool}, if any. This is not done automatically.
// returns true if the actor was removed from this group.
func (a *Actor) RemoveActor(actor *Actor) bool {
	// if (!children.removeValue(actor, true)) return false;
	a.Parent = nil
	a.ChildrenChanged()
	return true
}

// Removes all actors from this group.
func (a *Actor) ClearChildren() {
	for _, child := range a.Children {
		child.Parent = nil
	}
	a.Children = make([]*Actor, 0)
	a.ChildrenChanged()
}

// Removes all children, actions, and listeners from this group.
func (a *Actor) Clear() {
	a.ClearActions()
	a.ClearChildren()
}

// Returns true if input events are processed by a actor.
func (a *Actor) IsTouchable() bool {
	return a.TouchState == TouchableEnabled
}

// Returns the X position of the specified {@link Align alignment}
func (a *Actor) GetXAlign(alignment utils.Alignment) float32 {
	x := a.X
	if (alignment & utils.AlignmentRight) != 0 {
		x += a.W
	} else if (alignment & utils.AlignmentLeft) == 0 {
		x += a.W / 2
	}
	return x
}

func (a *Actor) SetX(x float32) {
	if a.X != x {
		a.X = x
		a.positionChanged()
	}
}

// Returns the Y position of the specified {@link Align alignment}
func (a *Actor) GetYAlign(alignment utils.Alignment) float32 {
	y := a.Y
	if (alignment & utils.AlignmentTop) != 0 {
		y += a.H
	} else if (alignment & utils.AlignmentBottom) == 0 {
		y += a.H / 2
	}
	return y
}

func (a *Actor) SetY(y float32) {
	if a.Y != y {
		a.Y = y
		a.positionChanged()
	}
}

// Sets the position of the actor's bottom left corner.
func (a *Actor) SetPosition(x, y float32) {
	if a.X != x || a.Y != y {
		a.X = x
		a.Y = y
		a.positionChanged()
	}
}

// Sets the position using the specified alignment.
// Note this may set the position to non-integer coordinates.
func (a *Actor) SetPositionAlign(x, y float32, alignment utils.Alignment) {
	if (alignment & utils.AlignmentRight) != 0 {
		a.X -= a.W
	} else if (alignment & utils.AlignmentLeft) == 0 {
		a.X -= a.W / 2
	}
	if (alignment & utils.AlignmentTop) != 0 {
		a.Y -= a.H
	} else if (alignment & utils.AlignmentBottom) == 0 {
		a.Y -= a.H / 2
	}
	if a.X != x || a.Y != y {
		a.X = x
		a.Y = y
		a.positionChanged()
	}
}

// Add x and y to current position.
func (a *Actor) MoveBy(x, y float32) {
	if x != 0 || y != 0 {
		a.X += x
		a.Y += y
		a.positionChanged()
	}
}

func (a *Actor) SetWidth(width float32) {
	oldWidth := a.W
	a.W = width
	if width != oldWidth {
		a.sizeChanged()
	}
}

func (a *Actor) SetHeight(height float32) {
	oldHeight := a.H
	a.H = height
	if height != oldHeight {
		a.sizeChanged()
	}
}

// Returns y plus height.
func (a *Actor) GetTop() float32 {
	return a.Y + a.H
}

//  Returns x plus width.
func (a *Actor) GetRight() float32 {
	return a.X + a.W
}

// Called when the actor's position has been changed.
func (a *Actor) positionChanged() {
}

// Called when the actor's size has been changed.
func (a *Actor) sizeChanged() {
}

// Sets the width and height.
func (a *Actor) SetSize(w, h float32) {
	oldWidth := a.W
	oldHeight := a.H
	a.W = w
	a.H = h
	if w != oldWidth || h != oldHeight {
		a.sizeChanged()
	}
}

// Adds the specified size to the current size.
func (a *Actor) SizeBy(w, h float32) {
	a.W += w
	a.H += h
	a.sizeChanged()
}

// Set bounds the x, y, width, and height.
func (a *Actor) SetBounds(x, y, w, h float32) {
	if a.X != x || a.Y != y {
		a.X = x
		a.Y = y
		a.positionChanged()
	}
	if a.W != w || a.H != h {
		a.W = h
		a.H = h
		a.sizeChanged()
	}
}

//Sets the origin position which is relative to the actor's bottom left corner.
func (a *Actor) setOrigin(originX, originY float32) {
	a.OX = originX
	a.OY = originY
}

// Sets the origin position to the specified {@link Align alignment}.
func (a *Actor) SetOriginAlign(alignment utils.Alignment) {
	if (alignment & utils.AlignmentLeft) != 0 {
		a.OX = 0
	} else if (alignment & utils.AlignmentRight) != 0 {
		a.OX = a.W
	} else {
		a.OX = a.W / 2
	}

	if (alignment & utils.AlignmentBottom) != 0 {
		a.OY = 0
	} else if (alignment & utils.AlignmentTop) != 0 {
		a.OY = a.H
	} else {
		a.OY = a.H / 2
	}
}

// Sets the scale X and scale Y.
func (a *Actor) SetScale(scaleX, scaleY float32) {
	a.SX = scaleX
	a.SY = scaleY
}

// Adds the specified scale to the current scale.
func (a *Actor) ScaleBy(scaleX, scaleY float32) {
	a.SX += scaleX
	a.SY += scaleY
}

func (a *Actor) GetRotation() float32 {
	return a.Rotation
}

func (a *Actor) SetRotation(degrees float32) {
	a.Rotation = degrees
}

// Adds the specified rotation to the current rotation.
func (a *Actor) RotateBy(amountInDegrees float32) {
	a.Rotation += amountInDegrees
}

// Changes the z-order for a actor so it is in front of all siblings.
func (a *Actor) ToFront() {
	a.SetZIndex(math.MaxInt32)
}

// Changes the z-order for a actor so it is in back of all siblings.
func (a *Actor) ToBack() {
	a.SetZIndex(0)
}

// Sets the z-index of a actor. The z-index is the index into the parent's {@link Group#getChildren() children}, where a
// lower index is below a higher index. Setting a z-index higher than the number of children will move the child to the front.
// Setting a z-index less than zero is invalid.
func (a *Actor) SetZIndex(index uint32) {
	if a.Parent != nil {
		if len(a.Parent.Children) == 1 {
			return
		}
		// if !a.Parent.RemoveActor(a) {
		// 	return
		// }
		// if index > len(a.Parent.Children) {
		// 	a.Parent.Children = append(a.Parent.Children, a)
		// } else {
		//   a.Parent.Children.insert(index, a);
		// }
	}
}

// Returns the z-index of a actor.
func (a *Actor) GetZIndex() int32 {
	if a.Parent == nil {
		return -1
	}
	return 0
	// return parent.children.indexOf(a, true);
}

//  public Actor hit (float x, float y, boolean touchable) {
//    if (touchable && getTouchable() == Touchable.disabled) return null;
//    Vector2 point = tmp;
//    Actor[] childrenArray = children.items;
//    for (int i = children.size - 1; i >= 0; i--) {
//      Actor child = childrenArray[i];
//      if (!child.isVisible()) continue;
//      child.parentToLocalCoordinates(point.set(x, y));
//      Actor hit = child.hit(point.x, point.y, touchable);
//      if (hit != null) return hit;
//    }
//    return super.hit(x, y, touchable);
//  }

//   // Sets a actor as the event {@link Event#setTarget(Actor) target} and propagates the event to a actor and ancestor
//    * actors as necessary. If a actor is not in the stage, the stage must be set before calling a method.
//    * <p>
//    * Events are fired in 2 phases.
//    * <ol>
//    * <li>The first phase (the "capture" phase) notifies listeners on each actor starting at the root and propagating downward to
//    * (and including) a actor.</li>
//    * <li>The second phase notifies listeners on each actor starting at a actor and, if {@link Event#getBubbles()} is true,
//    * propagating upward to the root.</li>
//    * </ol>
//    * If the event is {@link Event#stop() stopped} at any time, it will not propagate to the next actor.
//    * @return true if the event was {@link Event#cancel() cancelled}.
//   public boolean fire (Event event) {
//     if (event.getStage() == null) event.setStage(getStage());
//     event.setTarget(a);

//     // Collect ancestors so event propagation is unaffected by hierarchy changes.
//     Array<Group> ancestors = Pools.obtain(Array.class);
//     Group parent = a.parent;
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

//   // Notifies a actor's listeners of the event. The event is not propagated to any parents. Before notifying the listeners,
//    * a actor is set as the {@link Event#getListenerActor() listener actor}. The event {@link Event#setTarget(Actor) target}
//    * must be set before calling a method. If a actor is not in the stage, the stage must be set before calling a method.
//    * param capture If true, the capture listeners will be notified instead of the regular listeners.
//    * @return true of the event was {@link Event#cancel() cancelled}.
//   public boolean notify (Event event, boolean capture) {
//     if (event.getTarget() == null) throw new IllegalArgumentException("The event target cannot be null.");

//     DelayedRemovalArray<EventListener> listeners = capture ? captureListeners : a.listeners;
//     if (listeners.size == 0) return event.isCancelled();

//     event.setListenerActor(a);
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
//             event.getStage().addTouchFocus(listener, a, inputEvent.getTarget(), inputEvent.getPointer(),
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
//    * This method is used to delegate touchDown, mouse, and enter/exit events. If a method returns null, those events will not
//    * occur on a Actor.
//    * <p>
//    * The default implementation returns a actor if the point is within a actor's bounds.
//    *
//    * param touchable If true, the hit detection will respect the {@link #setTouchable(Touchable) touchability}.
//    * @see Touchable
//   public Actor hit (float x, float y, boolean touchable) {
//     if (touchable && a.touchable != Touchable.enabled) return null;
//     return x >= 0 && x < width && y >= 0 && y < height ? a : null;
//   }

//   // Returns true if a actor is the same as or is the descendant of the specified actor.
//   public boolean isDescendantOf (Actor actor) {
//     if (actor == null) throw new IllegalArgumentException("actor cannot be null.");
//     Actor parent = a;
//     while (true) {
//       if (parent == null) return false;
//       if (parent == actor) return true;
//       parent = parent.parent;
//     }
//   }

//   // Returns true if a actor is the same as or is the ascendant of the specified actor.
//   public boolean isAscendantOf (Actor actor) {
//     if (actor == null) throw new IllegalArgumentException("actor cannot be null.");
//     while (true) {
//       if (actor == null) return false;
//       if (actor == a) return true;
//       actor = actor.parent;
//     }
//   }

//   // Calls {@link #clipBegin(float, float, float, float)} to clip a actor's bounds.
//   public boolean clipBegin () {
//     return clipBegin(x, y, width, height);
//   }

//   // Clips the specified screen aligned rectangle, specified relative to the transform matrix of the stage's Batch. The transform
//    * matrix and the stage's camera must not have rotational components. Calling a method must be followed by a call to
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
//     Stage stage = a.stage;
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
//     Stage stage = a.stage;
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
//     final float rotation = -a.rotation;
//     final float scaleX = a.scaleX;
//     final float scaleY = a.scaleY;
//     final float x = a.x;
//     final float y = a.y;
//     if (rotation == 0) {
//       if (scaleX == 1 && scaleY == 1) {
//         localCoords.x += x;
//         localCoords.y += y;
//       } else {
//         final float originX = a.originX;
//         final float originY = a.originY;
//         localCoords.x = (localCoords.x - originX) * scaleX + originX + x;
//         localCoords.y = (localCoords.y - originY) * scaleY + originY + y;
//       }
//     } else {
//       final float cos = (float)Math.cos(rotation * MathUtils.degreesToRadians);
//       final float sin = (float)Math.sin(rotation * MathUtils.degreesToRadians);
//       final float originX = a.originX;
//       final float originY = a.originY;
//       final float tox = (localCoords.x - originX) * scaleX;
//       final float toy = (localCoords.y - originY) * scaleY;
//       localCoords.x = (tox * cos + toy * sin) + originX + x;
//       localCoords.y = (tox * -sin + toy * cos) + originY + y;
//     }
//     return localCoords;
//   }

//   // Converts coordinates for a actor to those of a parent actor. The ascendant does not need to be a direct parent.
//   public Vector2 localToAscendantCoordinates (Actor ascendant, Vector2 localCoords) {
//     Actor actor = a;
//     while (actor != null) {
//       actor.localToParentCoordinates(localCoords);
//       actor = actor.parent;
//       if (actor == ascendant) break;
//     }
//     return localCoords;
//   }

//   // Converts the coordinates given in the parent's coordinate system to a actor's coordinate system.
//   public Vector2 parentToLocalCoordinates (Vector2 parentCoords) {
//     final float rotation = a.rotation;
//     final float scaleX = a.scaleX;
//     final float scaleY = a.scaleY;
//     final float childX = x;
//     final float childY = y;
//     if (rotation == 0) {
//       if (scaleX == 1 && scaleY == 1) {
//         parentCoords.x -= childX;
//         parentCoords.y -= childY;
//       } else {
//         final float originX = a.originX;
//         final float originY = a.originY;
//         parentCoords.x = (parentCoords.x - childX - originX) / scaleX + originX;
//         parentCoords.y = (parentCoords.y - childY - originY) / scaleY + originY;
//       }
//     } else {
//       final float cos = (float)Math.cos(rotation * MathUtils.degreesToRadians);
//       final float sin = (float)Math.sin(rotation * MathUtils.degreesToRadians);
//       final float originX = a.originX;
//       final float originY = a.originY;
//       final float tox = parentCoords.x - childX - originX;
//       final float toy = parentCoords.y - childY - originY;
//       parentCoords.x = (tox * cos + toy * sin) / scaleX + originX;
//       parentCoords.y = (tox * -sin + toy * cos) / scaleY + originY;
//     }
//     return parentCoords;
//   }

//   // Draws a actor's debug lines if {@link #getDebug()} is true.
//   public void drawDebug (ShapeRenderer shapes) {
//     drawDebugBounds(shapes);
//   }

//   // Draws a rectange for the bounds of a actor if {@link #getDebug()} is true.
//   protected void drawDebugBounds (ShapeRenderer shapes) {
//     if (!debug) return;
//     shapes.set(ShapeType.Line);
//     shapes.setColor(stage.getDebugColor());
//     shapes.rect(x, y, originX, originY, width, height, scaleX, scaleY, rotation);
//   }

func (a *Actor) GetBounds() *shape.Rectangle {
	return shape.NewRectangle(a.X, a.Y, a.W, a.H)
}

func (a *Actor) CollidesXY(x, y float32) bool {
	if a.GetBounds().Overlaps(shape.NewRectangle(x, y, 5, 5)) {
		return true
	}
	return false
}

func (a *Actor) Collides(other *Actor) bool {
	if a.GetBounds().Overlaps(other.GetBounds()) {
		return true
	}
	return false
}

func (a *Actor) String() string {
	return a.Name
}

/*****************************************************************************\
|													  Actions                                           |
|                                                                             |
\*****************************************************************************/

// Sets the actor's Actor#setTouchable(Touchable) touchability.
func (a *Actor) ActionTouchable(touchable Touchable) *Actor {
	action := actionsPool.Get().(*Action)
	action.Act = func(action *Action, delta float32) bool {
		a.TouchState = touchable
		return true
	}
	return a
}

func (a *Actor) ActionVisible(visible bool) *Actor {
	action := actionsPool.Get().(*Action)
	action.Act = func(action *Action, delta float32) bool {
		a.Visible = visible
		return true
	}
	return a
}

func (a *Actor) Run(callback func()) *Actor {
	action := actionsPool.Get().(*Action)
	action.Act = func(a *Action, delta float32) bool {
		callback()
		return true
	}
	a.AddAction(action)
	return a
}

func (a *Actor) Delay(duration float32) *Actor {
	currenTime := float32(0)
	action := actionsPool.Get().(*Action)
	action.Act = func(a *Action, delta float32) bool {
		// println(delta)
		println(currenTime)
		if currenTime < duration {
			currenTime += 11
			if currenTime < duration {
				return false
			}
			delta = currenTime - duration
		}
		return true
	}
	a.AddAction(action)
	return a
}

// Repeat an action n number of times.
func (a *Actor) Repeat(n int, callback func()) {
	action := actionsPool.Get().(*Action)
	action.Act = func(action *Action, delta float32) bool {
		if n == 0 {
			return true
		}
		n--
		callback()
		return false
	}
	a.AddAction(action)
}

// Forever Executes an action always unless removed
func (a *Actor) Forever(callback func()) {
	action := actionsPool.Get().(*Action)
	action.Act = func(action *Action, delta float32) bool {
		callback()
		return false
	}
	a.AddAction(action)
}

// Parallel Executes a number of actions at the same time.
func (a *Actor) Parallel(actions ...*Action) *Actor {
	complete := false
	action := actionsPool.Get().(*Action)
	action.Reset = func(a *Action) {
		a.Actor = nil
		a.Target = nil
		a.Restart(a)
		actions = make([]*Action, 0)
	}
	action.Restart = func(a *Action) {
		complete = false
		for _, action := range actions {
			action.Restart(action)
		}
	}
	action.Act = func(action *Action, delta float32) bool {
		if complete {
			return true
		}
		complete = true
		for i := 0; i < len(actions) && action.Actor != nil; i++ {
			currentAction := actions[i]
			if currentAction.Actor != nil && !currentAction.Act(currentAction, delta) {
				complete = false
			}
			if action.Actor == nil {
				return true
			}
		}
		return complete
	}
	return a
}

// Sequence executes a number of actions one at a time.
func (a *Actor) Sequence(actions ...*Action) *Actor {
	index := 0
	action := actionsPool.Get().(*Action)
	action.Reset = func(a *Action) {
		actions = make([]*Action, 0)
	}
	action.Restart = func(a *Action) {
		for _, action := range actions {
			action.Restart(action)
		}
		index = 0
	}
	action.Act = func(action *Action, delta float32) bool {
		if index >= len(actions) {
			return true
		}
		if actions[index].Act(actions[index], delta) {
			if action.Actor == nil {
				return true
			}
			index++
			if index >= len(actions) {
				return true
			}
		}
		return false
	}
	return a
}

// Base class for actions that transition over time using the percent complete.
func generateTemporalAction(duration float32, interp Interpolation, update func(percent float32)) *Action {
	reverse, began, complete := false, false, false
	percent := float32(0)
	time := float32(0)
	action := actionsPool.Get().(*Action)
	action.Act = func(self *Action, delta float32) bool {
		if !began {
			//begin()
			began = true
		}
		time += delta
		complete = time >= duration
		if complete {
			percent = 1
		} else {
			percent = time / duration
			if interp != nil {
				percent = interp(percent)
			}
		}
		if reverse {
			update(1 - percent)
		} else {
			update(percent)
		}
		// if complete {
		//	end()
		//}
		return complete
	}
	return action
}

// Base class for actions that tgenerateTemporalActionransition over time using
// the percent complete since the last frame.
func generateRelativeTemporalAction(duration float32, interp Interpolation, update func(percentDelta float32)) *Action {
	lastPercent := float32(0)
	return generateTemporalAction(duration, interp, func(percent float32) {
		update(percent - lastPercent)
		lastPercent = percent
	})
}

// ActionAlpha sets the alpha for an actor's color (or a specified color), from the current
// alpha to the new alpha. Note this action transitions from the alpha at the
// time the action starts to the specified alpha.
// Transitions from the alpha at the time this action starts to the specified alpha.
func (a *Actor) ActionAlpha(alphaValue, duration float32, interp Interpolation) *Actor {
	start := a.Color.A
	action := generateTemporalAction(duration, interp, func(percent float32) {
		a.Color.A = start + (alphaValue-start)*percent
	})
	a.AddAction(action)
	return a
}

// ActionColor sets the actor's color (or a specified color), from the current
// to the new color. Note this action transitions from the color at the time the action
// starts to the specified color.
func (a *Actor) ActionColor(end *Color, duration float32, interp Interpolation) *Actor {
	startR := a.Color.R
	startG := a.Color.G
	startB := a.Color.B
	startA := a.Color.A
	action := generateTemporalAction(duration, interp, func(percent float32) {
		dr := startR + (end.R-startR)*percent
		dg := startG + (end.G-startG)*percent
		db := startB + (end.B-startB)*percent
		da := startA + (end.A-startA)*percent
		a.Color.Set(dr, dg, db, da)
	})
	a.AddAction(action)
	return a
}

// ActionMoveBy moves an actor to a relative position.
func (a *Actor) ActionMoveBy(amountX, amountY, duration float32, interp Interpolation) *Actor {
	action := generateRelativeTemporalAction(duration, interp, func(percentDelta float32) {
		a.MoveBy(amountX*percentDelta, amountY*percentDelta)
	})
	a.AddAction(action)
	return a
}

// ActionMoveTo moves an actor from its current position to a specific position.
func (a *Actor) ActionMoveTo(endX, endY, duration float32, interp Interpolation) *Actor {
	startX := a.X
	startY := a.Y
	// Align.bottomLeftaction
	action := generateTemporalAction(duration, interp, func(percent float32) {
		a.SetPosition(startX+(endX-startX)*percent, startY+(endY-startY)*percent)
	})
	a.AddAction(action)
	return a
}

// ActionMoveToAligned moves an actor from its current position to a specific position
// with a alignment.
func (a *Actor) ActionMoveToAligned(endX, endY float32, alignment utils.Alignment, duration float32, interp Interpolation) *Actor {
	startX := a.X
	startY := a.Y
	action := generateTemporalAction(duration, interp, func(percent float32) {
		a.SetPositionAlign(startX+(endX-startX)*percent, startY+(endY-startY)*percent, alignment)
	})
	a.AddAction(action)
	return a
}

// ActionSizeTo moves an actor from its current size to a specific size.
func (a *Actor) ActionSizeTo(endW, endH, duration float32, interp Interpolation) *Actor {
	startW := a.W
	startH := a.H
	action := generateTemporalAction(duration, interp, func(percent float32) {
		a.SetSize(startW+(endW-startW)*percent, startH+(endH-startH)*percent)
	})
	a.AddAction(action)
	return a
}

// ActionSizeBy moves an actor from its current size to a relative size.
func (a *Actor) ActionSizeBy(amountW, amountH, duration float32, interp Interpolation) *Actor {
	action := generateRelativeTemporalAction(duration, interp, func(percentDelta float32) {
		a.SizeBy(amountW*percentDelta, amountH*percentDelta)
	})
	a.AddAction(action)
	return a
}

// ActionScaleTo sets the actor's scale from its current value to a specific value.
func (a *Actor) ActionScaleTo(endX, endY, duration float32, interp Interpolation) *Actor {
	startX := a.X
	startY := a.Y
	action := generateTemporalAction(duration, interp, func(percent float32) {
		a.SetScale(startX+(endX-startX)*percent, startY+(endY-startY)*percent)
	})
	a.AddAction(action)
	return a
}

// ActionScaleBy scales an actor's scale to a relative size.
func (a *Actor) ActionScaleBy(amountX, amountY, duration float32, interp Interpolation) *Actor {
	action := generateRelativeTemporalAction(duration, interp, func(percentDelta float32) {
		a.ScaleBy(amountX*percentDelta, amountY*percentDelta)
	})
	a.AddAction(action)
	return a
}

// ActionRotateTo sets the actor's rotation from its current value to a specific value.
func (a *Actor) ActionRotateTo(end, duration float32, interp Interpolation) *Actor {
	start := a.Rotation
	action := generateTemporalAction(duration, interp, func(percent float32) {
		a.SetRotation(start + (end-start)*percent)
	})
	a.AddAction(action)
	return a
}

// ActionRotateBy Sets the actor's rotation from its current value to a relative value.
func (a *Actor) ActionRotateBy(amount, duration float32, interp Interpolation) *Actor {
	action := generateRelativeTemporalAction(duration, interp, func(percentDelta float32) {
		a.RotateBy(amount * percentDelta)
	})
	a.AddAction(action)
	return a
}

/*****************************************************************************\
|													  Effects                                           |
|                                                                             |
\*****************************************************************************/

// enum EffectDuration {
//   Once,
//   OnceToAndBack,
//   Looping,
//   LoopingToAndBack
// }

// EffectFadeIn transitions from the alpha at the time this action starts to an alpha of 1.
func (a *Actor) EffectFadeIn(duration float32, interp Interpolation) *Actor {
	a.Color.A = 0
	return a.ActionAlpha(1, duration, interp)
}

// EffectFadeOut transitions from the alpha at the time this action starts to an alpha of 0.
func (a *Actor) EffectFadeOut(duration float32, interp Interpolation) *Actor {
	a.Color.A = 1
	return a.ActionAlpha(0, duration, interp)
}

//EffectScaleIn sets the actors scale to 0 and
func (a *Actor) EffectScaleIn(duration float32, interp Interpolation) *Actor {
	a.SetScale(0, 0)
	a.ActionScaleTo(1, 1, duration, interp)
	return a
}

//EffectScaleOut sets the actors scale to 0 and
func (a *Actor) EffectScaleOut(duration float32, interp Interpolation) *Actor {
	a.SetScale(1, 1)
	a.ActionScaleTo(0, 0, duration, interp)
	return a
}

//EffectShakeInOut sets the actors scale to 0 and
func (a *Actor) EffectShakeInOut(value, duration float32, interp Interpolation) *Actor {
	return a.ActionRotateTo(value, duration, interp).
		ActionRotateTo(-value, duration, interp).
		ActionRotateTo(0, duration, interp)
}

//EffectScaleAndShake sets the actors scale to 0 and
func (a *Actor) EffectScaleAndShake(scaleRatioX, scaleRatioY, shakeAngle, duration float32) *Actor {
	originalAngle := a.Rotation
	a.ActionScaleTo(scaleRatioX, scaleRatioY, duration, nil)
	a.ActionRotateTo(shakeAngle, duration, nil)
	a.ActionRotateTo(-shakeAngle, duration, nil)
	a.ActionRotateTo(originalAngle, duration, nil)
	a.ActionScaleTo(1, 1, duration, nil)
	return a
}

// EffectScaleAndShakeFadeOut Scale effect, Back To Normal, Fade Out (SC, BTN, FO)
func (a *Actor) EffectScaleAndShakeFadeOut(scaleRatioX, scaleRatioY, fadeBeforeDuration, duration float32) *Actor {
	a.ActionScaleTo(scaleRatioX, scaleRatioY, duration, nil)
	a.ActionScaleTo(1, 1, duration, nil)
	a.Delay(fadeBeforeDuration)
	a.EffectFadeOut(duration, nil)
	return a
}

//     switch(effectType){
//       case PatrolX:
//         actor.addAction(Actions.forever(Actions.sequence(Actions.moveBy(value, 0, duration, interp),
//             Actions.moveBy(-value, 0, duration, interp))));
//         break;
//       case PatrolY:
//         actor.addAction(Actions.forever(Actions.sequence(Actions.moveBy(0, value, duration, interp),
//             Actions.moveBy(0, -value, duration, interp))));
//         break;
//       case SlideLeft:
//         actor.setPosition(999, y);
//         actor.addAction(Actions.moveTo(x, y, duration, interp));
//         break;
//       case SlideRight:
//         actor.setPosition(-999, y);
//         actor.addAction(Actions.moveTo(x, y, duration, interp));
//         break;
//       case SlideUp:
//         actor.setPosition(x, -999);
//         actor.addAction(Actions.moveTo(x, y, duration, interp));
//         break;
//       case SlideDown:
//         actor.setPosition(x, 999);
//         actor.addAction(Actions.moveTo(x, y, duration, interp));
//         break;
