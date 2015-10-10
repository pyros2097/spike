// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

type Pool struct {
}

func (self *Pool) Free(actor *Action) {

}

type Action struct {
	// The actor this action is attached to, or nil if it is not attached.
	actor *Actor

	// The actor this action targets, or nil if a target has not been set.
	target *Actor

	pool *Pool
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
func (self *Action) SetActor(actor *Actor) {
	self.actor = actor
	if self.target == nil {
		self.target = actor
	}
	if self.actor == nil {
		if self.pool != nil {
			self.pool.Free(self)
			self.pool = nil
		}
	}
}

// @return null if the action is not attached to an actor.
func (self *Action) GetActor() *Actor {
	return self.actor
}

/** Sets the actor this action will manipulate. If no target actor is set, {@link #setActor(Actor)} will set the target actor
 * when the action is added to an actor. */
func (self *Action) SetTarget(target *Actor) {
	self.target = target
}

// @return null if the action has no target.
func (self *Action) GetTarget() *Actor {
	return self.target
}

func (self *Action) GetPool() *Pool {
	return self.pool
}

/** Sets the pool that the action will be returned to when removed from the actor.
 * @param pool May be null.
 * @see #setActor(Actor) */
func (self *Action) SetPool(pool *Pool) {
	self.pool = pool
}

/** Updates the action based on time. Typically this is called each frame by {@link Actor#act(float)}.
 * @param delta Time in seconds since the last frame.
 * @return true if the action is done. This method may continue to be called after the action is done. */
func (self *Action) Act(delta float32) bool {
	return false
}

// Sets the state of the action so it can be run again.
func (self *Action) Restart() {
}

/** Resets the optional state of this action to as if it were newly created, allowing the action to be pooled and reused. State
 * required to be set for every usage of this action or computed during the action does not need to be reset.
 * <p>
 * The default implementation calls {@link #restart()}.
 * <p>
 * If a subclass has optional state, it must override this method, call super, and reset the optional state. */
func (self *Action) Reset() {
	self.actor = nil
	self.target = nil
	self.pool = nil
	self.Restart()
}

func (self *Action) String() string {
	return "Action"
}

// Adds an action to an actor.
type AddAction struct {
	Action
	base *Action
}

func (self *AddAction) Act(delta float32) bool {
	// self.target.AddAction(self.base)
	return true
}

func (self *AddAction) Restart() {
	if self.base != nil {
		self.base.Restart()
	}
}

func (self *AddAction) Reset() {
	self.Action.Reset()
	self.base = nil
}

// Static convenience methods for using pooled actions, intended for static import.
// public class Actions {
// 	/** Returns a new or pooled action of the specified type. */
// 	static public <T extends Action> T action (Class<T> type) {
// 		Pool<T> pool = Pools.get(type);
// 		T action = pool.obtain();
// 		action.setPool(pool);
// 		return action;
// 	}

// 	static public AddAction addAction (Action action) {
// 		AddAction addAction = action(AddAction.class);
// 		addAction.setAction(action);
// 		return addAction;
// 	}

// 	static public AddAction addAction (Action action, Actor targetActor) {
// 		AddAction addAction = action(AddAction.class);
// 		addAction.setTarget(targetActor);
// 		addAction.setAction(action);
// 		return addAction;
// 	}

// 	static public RemoveAction removeAction (Action action) {
// 		RemoveAction removeAction = action(RemoveAction.class);
// 		removeAction.setAction(action);
// 		return removeAction;
// 	}

// 	static public RemoveAction removeAction (Action action, Actor targetActor) {
// 		RemoveAction removeAction = action(RemoveAction.class);
// 		removeAction.setTarget(targetActor);
// 		removeAction.setAction(action);
// 		return removeAction;
// 	}

// 	/** Moves the actor instantly. */
// 	static public MoveToAction moveTo (float x, float y) {
// 		return moveTo(x, y, 0, null);
// 	}

// 	static public MoveToAction moveTo (float x, float y, float duration) {
// 		return moveTo(x, y, duration, null);
// 	}

// 	static public MoveToAction moveTo (float x, float y, float duration, Interpolation interpolation) {
// 		MoveToAction action = action(MoveToAction.class);
// 		action.setPosition(x, y);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	static public MoveToAction moveToAligned (float x, float y, int alignment) {
// 		return moveToAligned(x, y, alignment, 0, null);
// 	}

// 	static public MoveToAction moveToAligned (float x, float y, int alignment, float duration) {
// 		return moveToAligned(x, y, alignment, duration, null);
// 	}

// 	static public MoveToAction moveToAligned (float x, float y, int alignment, float duration, Interpolation interpolation) {
// 		MoveToAction action = action(MoveToAction.class);
// 		action.setPosition(x, y, alignment);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Moves the actor instantly. */
// 	static public MoveByAction moveBy (float amountX, float amountY) {
// 		return moveBy(amountX, amountY, 0, null);
// 	}

// 	static public MoveByAction moveBy (float amountX, float amountY, float duration) {
// 		return moveBy(amountX, amountY, duration, null);
// 	}

// 	static public MoveByAction moveBy (float amountX, float amountY, float duration, Interpolation interpolation) {
// 		MoveByAction action = action(MoveByAction.class);
// 		action.setAmount(amountX, amountY);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Sizes the actor instantly. */
// 	static public SizeToAction sizeTo (float x, float y) {
// 		return sizeTo(x, y, 0, null);
// 	}

// 	static public SizeToAction sizeTo (float x, float y, float duration) {
// 		return sizeTo(x, y, duration, null);
// 	}

// 	static public SizeToAction sizeTo (float x, float y, float duration, Interpolation interpolation) {
// 		SizeToAction action = action(SizeToAction.class);
// 		action.setSize(x, y);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Sizes the actor instantly. */
// 	static public SizeByAction sizeBy (float amountX, float amountY) {
// 		return sizeBy(amountX, amountY, 0, null);
// 	}

// 	static public SizeByAction sizeBy (float amountX, float amountY, float duration) {
// 		return sizeBy(amountX, amountY, duration, null);
// 	}

// 	static public SizeByAction sizeBy (float amountX, float amountY, float duration, Interpolation interpolation) {
// 		SizeByAction action = action(SizeByAction.class);
// 		action.setAmount(amountX, amountY);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Scales the actor instantly. */
// 	static public ScaleToAction scaleTo (float x, float y) {
// 		return scaleTo(x, y, 0, null);
// 	}

// 	static public ScaleToAction scaleTo (float x, float y, float duration) {
// 		return scaleTo(x, y, duration, null);
// 	}

// 	static public ScaleToAction scaleTo (float x, float y, float duration, Interpolation interpolation) {
// 		ScaleToAction action = action(ScaleToAction.class);
// 		action.setScale(x, y);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Scales the actor instantly. */
// 	static public ScaleByAction scaleBy (float amountX, float amountY) {
// 		return scaleBy(amountX, amountY, 0, null);
// 	}

// 	static public ScaleByAction scaleBy (float amountX, float amountY, float duration) {
// 		return scaleBy(amountX, amountY, duration, null);
// 	}

// 	static public ScaleByAction scaleBy (float amountX, float amountY, float duration, Interpolation interpolation) {
// 		ScaleByAction action = action(ScaleByAction.class);
// 		action.setAmount(amountX, amountY);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Rotates the actor instantly. */
// 	static public RotateToAction rotateTo (float rotation) {
// 		return rotateTo(rotation, 0, null);
// 	}

// 	static public RotateToAction rotateTo (float rotation, float duration) {
// 		return rotateTo(rotation, duration, null);
// 	}

// 	static public RotateToAction rotateTo (float rotation, float duration, Interpolation interpolation) {
// 		RotateToAction action = action(RotateToAction.class);
// 		action.setRotation(rotation);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Rotates the actor instantly. */
// 	static public RotateByAction rotateBy (float rotationAmount) {
// 		return rotateBy(rotationAmount, 0, null);
// 	}

// 	static public RotateByAction rotateBy (float rotationAmount, float duration) {
// 		return rotateBy(rotationAmount, duration, null);
// 	}

// 	static public RotateByAction rotateBy (float rotationAmount, float duration, Interpolation interpolation) {
// 		RotateByAction action = action(RotateByAction.class);
// 		action.setAmount(rotationAmount);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Sets the actor's color instantly. */
// 	static public ColorAction color (Color color) {
// 		return color(color, 0, null);
// 	}

// 	/** Transitions from the color at the time this action starts to the specified color. */
// 	static public ColorAction color (Color color, float duration) {
// 		return color(color, duration, null);
// 	}

// 	/** Transitions from the color at the time this action starts to the specified color. */
// 	static public ColorAction color (Color color, float duration, Interpolation interpolation) {
// 		ColorAction action = action(ColorAction.class);
// 		action.setEndColor(color);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Sets the actor's alpha instantly. */
// 	static public AlphaAction alpha (float a) {
// 		return alpha(a, 0, null);
// 	}

// 	/** Transitions from the alpha at the time this action starts to the specified alpha. */
// 	static public AlphaAction alpha (float a, float duration) {
// 		return alpha(a, duration, null);
// 	}

// 	/** Transitions from the alpha at the time this action starts to the specified alpha. */
// 	static public AlphaAction alpha (float a, float duration, Interpolation interpolation) {
// 		AlphaAction action = action(AlphaAction.class);
// 		action.setAlpha(a);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Transitions from the alpha at the time this action starts to an alpha of 0. */
// 	static public AlphaAction fadeOut (float duration) {
// 		return alpha(0, duration, null);
// 	}

// 	/** Transitions from the alpha at the time this action starts to an alpha of 0. */
// 	static public AlphaAction fadeOut (float duration, Interpolation interpolation) {
// 		AlphaAction action = action(AlphaAction.class);
// 		action.setAlpha(0);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	/** Transitions from the alpha at the time this action starts to an alpha of 1. */
// 	static public AlphaAction fadeIn (float duration) {
// 		return alpha(1, duration, null);
// 	}

// 	/** Transitions from the alpha at the time this action starts to an alpha of 1. */
// 	static public AlphaAction fadeIn (float duration, Interpolation interpolation) {
// 		AlphaAction action = action(AlphaAction.class);
// 		action.setAlpha(1);
// 		action.setDuration(duration);
// 		action.setInterpolation(interpolation);
// 		return action;
// 	}

// 	static public VisibleAction show () {
// 		return visible(true);
// 	}

// 	static public VisibleAction hide () {
// 		return visible(false);
// 	}

// 	static public VisibleAction visible (boolean visible) {
// 		VisibleAction action = action(VisibleAction.class);
// 		action.setVisible(visible);
// 		return action;
// 	}

// 	static public TouchableAction touchable (Touchable touchable) {
// 		TouchableAction action = action(TouchableAction.class);
// 		action.setTouchable(touchable);
// 		return action;
// 	}

// 	static public RemoveActorAction removeActor () {
// 		return action(RemoveActorAction.class);
// 	}

// 	static public RemoveActorAction removeActor (Actor removeActor) {
// 		RemoveActorAction action = action(RemoveActorAction.class);
// 		action.setTarget(removeActor);
// 		return action;
// 	}

// 	static public DelayAction delay (float duration) {
// 		DelayAction action = action(DelayAction.class);
// 		action.setDuration(duration);
// 		return action;
// 	}

// 	static public DelayAction delay (float duration, Action delayedAction) {
// 		DelayAction action = action(DelayAction.class);
// 		action.setDuration(duration);
// 		action.setAction(delayedAction);
// 		return action;
// 	}

// 	static public TimeScaleAction timeScale (float scale, Action scaledAction) {
// 		TimeScaleAction action = action(TimeScaleAction.class);
// 		action.setScale(scale);
// 		action.setAction(scaledAction);
// 		return action;
// 	}

// 	static public SequenceAction sequence (Action action1) {
// 		SequenceAction action = action(SequenceAction.class);
// 		action.addAction(action1);
// 		return action;
// 	}

// 	static public SequenceAction sequence (Action action1, Action action2) {
// 		SequenceAction action = action(SequenceAction.class);
// 		action.addAction(action1);
// 		action.addAction(action2);
// 		return action;
// 	}

// 	static public SequenceAction sequence (Action action1, Action action2, Action action3) {
// 		SequenceAction action = action(SequenceAction.class);
// 		action.addAction(action1);
// 		action.addAction(action2);
// 		action.addAction(action3);
// 		return action;
// 	}

// 	static public SequenceAction sequence (Action action1, Action action2, Action action3, Action action4) {
// 		SequenceAction action = action(SequenceAction.class);
// 		action.addAction(action1);
// 		action.addAction(action2);
// 		action.addAction(action3);
// 		action.addAction(action4);
// 		return action;
// 	}

// 	static public SequenceAction sequence (Action action1, Action action2, Action action3, Action action4, Action action5) {
// 		SequenceAction action = action(SequenceAction.class);
// 		action.addAction(action1);
// 		action.addAction(action2);
// 		action.addAction(action3);
// 		action.addAction(action4);
// 		action.addAction(action5);
// 		return action;
// 	}

// 	static public SequenceAction sequence (Action... actions) {
// 		SequenceAction action = action(SequenceAction.class);
// 		for (int i = 0, n = actions.length; i < n; i++)
// 			action.addAction(actions[i]);
// 		return action;
// 	}

// 	static public SequenceAction sequence () {
// 		return action(SequenceAction.class);
// 	}

// 	static public ParallelAction parallel (Action action1) {
// 		ParallelAction action = action(ParallelAction.class);
// 		action.addAction(action1);
// 		return action;
// 	}

// 	static public ParallelAction parallel (Action action1, Action action2) {
// 		ParallelAction action = action(ParallelAction.class);
// 		action.addAction(action1);
// 		action.addAction(action2);
// 		return action;
// 	}

// 	static public ParallelAction parallel (Action action1, Action action2, Action action3) {
// 		ParallelAction action = action(ParallelAction.class);
// 		action.addAction(action1);
// 		action.addAction(action2);
// 		action.addAction(action3);
// 		return action;
// 	}

// 	static public ParallelAction parallel (Action action1, Action action2, Action action3, Action action4) {
// 		ParallelAction action = action(ParallelAction.class);
// 		action.addAction(action1);
// 		action.addAction(action2);
// 		action.addAction(action3);
// 		action.addAction(action4);
// 		return action;
// 	}

// 	static public ParallelAction parallel (Action action1, Action action2, Action action3, Action action4, Action action5) {
// 		ParallelAction action = action(ParallelAction.class);
// 		action.addAction(action1);
// 		action.addAction(action2);
// 		action.addAction(action3);
// 		action.addAction(action4);
// 		action.addAction(action5);
// 		return action;
// 	}

// 	static public ParallelAction parallel (Action... actions) {
// 		ParallelAction action = action(ParallelAction.class);
// 		for (int i = 0, n = actions.length; i < n; i++)
// 			action.addAction(actions[i]);
// 		return action;
// 	}

// 	static public ParallelAction parallel () {
// 		return action(ParallelAction.class);
// 	}

// 	static public RepeatAction repeat (int count, Action repeatedAction) {
// 		RepeatAction action = action(RepeatAction.class);
// 		action.setCount(count);
// 		action.setAction(repeatedAction);
// 		return action;
// 	}

// 	static public RepeatAction forever (Action repeatedAction) {
// 		RepeatAction action = action(RepeatAction.class);
// 		action.setCount(RepeatAction.FOREVER);
// 		action.setAction(repeatedAction);
// 		return action;
// 	}

// 	static public RunnableAction run (Runnable runnable) {
// 		RunnableAction action = action(RunnableAction.class);
// 		action.setRunnable(runnable);
// 		return action;
// 	}

// 	static public LayoutAction layout (boolean enabled) {
// 		LayoutAction action = action(LayoutAction.class);
// 		action.setLayoutEnabled(enabled);
// 		return action;
// 	}

// 	static public AfterAction after (Action action) {
// 		AfterAction afterAction = action(AfterAction.class);
// 		afterAction.setAction(action);
// 		return afterAction;
// 	}

// 	static public AddListenerAction addListener (EventListener listener, boolean capture) {
// 		AddListenerAction addAction = action(AddListenerAction.class);
// 		addAction.setListener(listener);
// 		addAction.setCapture(capture);
// 		return addAction;
// 	}

// 	static public AddListenerAction addListener (EventListener listener, boolean capture, Actor targetActor) {
// 		AddListenerAction addAction = action(AddListenerAction.class);
// 		addAction.setTarget(targetActor);
// 		addAction.setListener(listener);
// 		addAction.setCapture(capture);
// 		return addAction;
// 	}

// 	static public RemoveListenerAction removeListener (EventListener listener, boolean capture) {
// 		RemoveListenerAction addAction = action(RemoveListenerAction.class);
// 		addAction.setListener(listener);
// 		addAction.setCapture(capture);
// 		return addAction;
// 	}

// 	static public RemoveListenerAction removeListener (EventListener listener, boolean capture, Actor targetActor) {
// 		RemoveListenerAction addAction = action(RemoveListenerAction.class);
// 		addAction.setTarget(targetActor);
// 		addAction.setListener(listener);
// 		addAction.setCapture(capture);
// 		return addAction;
// 	}
// }

/** Returns a new or pooled action of the specified type. */
// static public <T extends Action3d> T action3d (Class<T> type) {
//         Pool<T> pool = Pools.get(type);
//         T action = pool.obtain();
//         action.setPool(pool);
//         return action;
// }

// static public AddAction addAction (Action3d action) {
//     AddAction addAction = action3d(AddAction.class);
//     addAction.setAction(action);
//     return addAction;
// }

// static public AddAction addAction (Action3d action, Actor3d targetActor) {
//     AddAction addAction = action3d(AddAction.class);
//     addAction.setTargetActor(targetActor);
//     addAction.setAction(action);
//     return addAction;
// }

// static public RemoveAction removeAction (Action3d action) {
//     RemoveAction removeAction = action3d(RemoveAction.class);
//     removeAction.setAction(action);
//     return removeAction;
// }

// static public RemoveAction removeAction (Action3d action, Actor3d targetActor) {
//     RemoveAction removeAction = action3d(RemoveAction.class);
//     removeAction.setTargetActor(targetActor);
//     removeAction.setAction(action);
//     return removeAction;
// }

// /** Moves the actor instantly. */
// static public MoveToAction moveTo (float x, float y, float z) {
//         return moveTo(x, y, z, 0, null);
// }

// static public MoveToAction moveTo (float x, float y, float z, float duration) {
//         return moveTo(x, y, z, duration, null);
// }

// static public MoveToAction moveTo (float x, float y, float z ,float duration, Interpolation interpolation) {
//         MoveToAction action = action3d(MoveToAction.class);
//         action.setPosition(x, y, z);
//         action.setDuration(duration);
//         action.setInterpolation(interpolation);
//         return action;
// }

// /** Moves the actor instantly. */
// static public MoveByAction moveBy (float amountX, float amountY, float amountZ) {
//         return moveBy(amountX, amountY, amountZ, 0, null);
// }

// static public MoveByAction moveBy (float amountX, float amountY, float amountZ, float duration) {
//         return moveBy(amountX, amountY, amountZ, duration, null);
// }

// static public MoveByAction moveBy (float amountX, float amountY, float amountZ, float duration, Interpolation interpolation) {
//         MoveByAction action = action3d(MoveByAction.class);
//         action.setAmount(amountX, amountY, amountZ);
//         action.setDuration(duration);
//         action.setInterpolation(interpolation);
//         return action;
// }

// /** Scales the actor instantly. */
// static public ScaleToAction scaleTo (float x, float y, float z) {
//         return scaleTo(x, y, z, 0, null);
// }

// static public ScaleToAction scaleTo (float x, float y, float z, float duration) {
//         return scaleTo(x, y, z, duration, null);
// }

// static public ScaleToAction scaleTo (float x, float y, float z, float duration, Interpolation interpolation) {
//         ScaleToAction action = action3d(ScaleToAction.class);
//         action.setScale(x, y, z);
//         action.setDuration(duration);
//         action.setInterpolation(interpolation);
//         return action;
// }

// /** Scales the actor instantly. */
// static public ScaleByAction scaleBy (float amountX, float amountY, float amountZ) {
//         return scaleBy(amountX, amountY, amountZ, 0, null);
// }

// static public ScaleByAction scaleBy (float amountX, float amountY, float amountZ, float duration) {
//         return scaleBy(amountX, amountY, amountZ, duration, null);
// }

// static public ScaleByAction scaleBy (float amountX, float amountY, float amountZ, float duration, Interpolation interpolation) {
//         ScaleByAction action = action3d(ScaleByAction.class);
//         action.setAmount(amountX, amountY, amountZ);
//         action.setDuration(duration);
//         action.setInterpolation(interpolation);
//         return action;
// }

// /** Rotates the actor instantly. */
// static public RotateToAction rotateTo (float yaw, float pitch, float roll) {
//         return rotateTo(yaw, pitch, roll, 0, null);
// }

// static public RotateToAction rotateTo (float yaw, float pitch, float roll, float duration) {
//     return rotateTo(yaw, pitch, roll, duration, null);
// }

// static public RotateToAction rotateTo (float yaw, float pitch, float roll, float duration, Interpolation interpolation) {
//     RotateToAction action = action3d(RotateToAction.class);
//     action.setRotation(yaw, pitch, roll);
//     action.setDuration(duration);
//     action.setInterpolation(interpolation);
//     return action;
// }

// /** Rotates the actor instantly. */
// static public RotateByAction rotateBy (float yaw, float pitch, float roll) {
//         return rotateBy(yaw, pitch, roll, 0, null);
// }

// static public RotateByAction rotateBy (float yaw, float pitch, float roll, float duration) {
//     return rotateBy(yaw, pitch, roll, duration, null);
// }

// static public RotateByAction rotateBy (float yaw, float pitch, float roll, float duration, Interpolation interpolation) {
//     RotateByAction action = action3d(RotateByAction.class);
//     action.setAmount(yaw, pitch, roll);
//     action.setDuration(duration);
//     action.setInterpolation(interpolation);
//     return action;
// }

// static public VisibleAction show () {
//     return visible(true);
// }

// static public VisibleAction hide () {
//     return visible(false);
// }

// static public VisibleAction visible (boolean visible) {
//     VisibleAction action = action3d(VisibleAction.class);
//     action.setVisible(visible);
//     return action;
// }

// static public RemoveActorAction removeActor (Actor3d removeActor) {
//     RemoveActorAction action = action3d(RemoveActorAction.class);
//     action.setRemoveActor(removeActor);
//     return action;
// }

// static public DelayAction delay (float duration) {
//     DelayAction action = action3d(DelayAction.class);
//     action.setDuration(duration);
//     return action;
// }

// static public DelayAction delay (float duration, Action3d delayedAction) {
//     DelayAction action = action3d(DelayAction.class);
//     action.setDuration(duration);
//     action.setAction(delayedAction);
//     return action;
// }

// static public TimeScaleAction timeScale (float scale, Action3d scaledAction) {
//     TimeScaleAction action = action3d(TimeScaleAction.class);
//     action.setScale(scale);
//     action.setAction(scaledAction);
//     return action;
// }

// static public SequenceAction sequence (Action3d action1) {
//     SequenceAction action = action3d(SequenceAction.class);
//     action.addAction(action1);
//     return action;
// }

// static public SequenceAction sequence (Action3d action1, Action3d action2) {
//     SequenceAction action = action3d(SequenceAction.class);
//     action.addAction(action1);
//     action.addAction(action2);
//     return action;
// }

// static public SequenceAction sequence (Action3d action1, Action3d action2, Action3d action3) {
//     SequenceAction action = action3d(SequenceAction.class);
//     action.addAction(action1);
//     action.addAction(action2);
//     action.addAction(action3);
//     return action;
// }

// static public SequenceAction sequence (Action3d action1, Action3d action2, Action3d action3, Action3d action4) {
//     SequenceAction action = action3d(SequenceAction.class);
//     action.addAction(action1);
//     action.addAction(action2);
//     action.addAction(action3);
//     action.addAction(action4);
//     return action;
// }

// static public SequenceAction sequence (Action3d action1, Action3d action2, Action3d action3, Action3d action4, Action3d action5) {
//     SequenceAction action = action3d(SequenceAction.class);
//     action.addAction(action1);
//     action.addAction(action2);
//     action.addAction(action3);
//     action.addAction(action4);
//     action.addAction(action5);
//     return action;
// }

// static public SequenceAction sequence (Action3d... actions) {
//     SequenceAction action = action3d(SequenceAction.class);
//     for (int i = 0, n = actions.length; i < n; i++)
//             action.addAction(actions[i]);
//     return action;
// }

// static public SequenceAction sequence () {
//     return action3d(SequenceAction.class);
// }

// static public ParallelAction parallel (Action3d action1) {
//     ParallelAction action = action3d(ParallelAction.class);
//     action.addAction(action1);
//     return action;
// }

// static public ParallelAction parallel (Action3d action1, Action3d action2) {
//     ParallelAction action = action3d(ParallelAction.class);
//     action.addAction(action1);
//     action.addAction(action2);
//     return action;
// }

// static public ParallelAction parallel (Action3d action1, Action3d action2, Action3d action3) {
//     ParallelAction action = action3d(ParallelAction.class);
//     action.addAction(action1);
//     action.addAction(action2);
//     action.addAction(action3);
//     return action;
// }

// static public ParallelAction parallel (Action3d action1, Action3d action2, Action3d action3, Action3d action4) {
//     ParallelAction action = action3d(ParallelAction.class);
//     action.addAction(action1);
//     action.addAction(action2);
//     action.addAction(action3);
//     action.addAction(action4);
//     return action;
// }

// static public ParallelAction parallel (Action3d action1, Action3d action2, Action3d action3, Action3d action4, Action3d action5) {
//     ParallelAction action = action3d(ParallelAction.class);
//     action.addAction(action1);
//     action.addAction(action2);
//     action.addAction(action3);
//     action.addAction(action4);
//     action.addAction(action5);
//     return action;
// }

// static public ParallelAction parallel (Action3d... actions) {
//     ParallelAction action = action3d(ParallelAction.class);
//     for (int i = 0, n = actions.length; i < n; i++)
//             action.addAction(actions[i]);
//     return action;
// }

// static public ParallelAction parallel () {
//     return action3d(ParallelAction.class);
// }

// static public RepeatAction repeat (int count, Action3d repeatedAction) {
//     RepeatAction action = action3d(RepeatAction.class);
//     action.setCount(count);
//     action.setAction(repeatedAction);
//     return action;
// }

// static public RepeatAction forever (Action3d repeatedAction) {
//     RepeatAction action = action3d(RepeatAction.class);
//     action.setCount(RepeatAction.FOREVER);
//     action.setAction(repeatedAction);
//     return action;
// }

// static public RunnableAction run (Runnable runnable) {
//     RunnableAction action = action3d(RunnableAction.class);
//     action.setRunnable(runnable);
//     return action;
// }

// static public AfterAction after (Action3d action) {
//     AfterAction afterAction = action3d(AfterAction.class);
//     afterAction.setAction(action);
//     return afterAction;
// }
