// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package actions

import (
	. "github.com/pyros2097/spike"
	"github.com/pyros2097/spike/graphics"
	. "github.com/pyros2097/spike/interpolation"
	"sync"
)

func superReset(self *Action) {
	self.Actor = nil
	self.Target = nil
	// self.pool = nil
	self.Restart(self)
}

// Adds an action to an actor.
func Add(base *Action) *Action {
	action := NewAction()
	action.Act = func(self *Action, delta float32) bool {
		action.Target.AddAction(base)
		return true
	}
	action.Restart = func(self *Action) {
		if base != nil {
			base.Restart(base)
		}
	}
	action.Reset = func(self *Action) {
		superReset(self)
		base = nil
	}
	return action
}

// Removes an action from an actor.
func Remove(base *Action) *Action {
	action := NewAction()
	action.Act = func(self *Action, delta float32) bool {
		action.Target.RemoveAction(base)
		return true
	}
	action.Restart = func(self *Action) {
		if base != nil {
			base.Restart(base)
		}
	}
	action.Reset = func(self *Action) {
		superReset(self)
		base = nil
	}
	return action
}

// Removes an actor from the stage.
func RemoveActor() *Action {
	removed := false
	action := NewAction()
	action.Act = func(self *Action, delta float32) bool {
		if !removed {
			removed = true
			action.Target.RemoveActor()
		}
		return true
	}
	action.Restart = func(self *Action) {
		removed = false
	}
	return action
}

var pool = sync.Pool{}

// Executes a number of actions at the same time.
func Parallel(actions ...*Action) *Action {
	complete := false
	action := NewAction()
	action.Reset = func(self *Action) {
		superReset(self)
		actions = make([]*Action, 0)
	}
	action.Restart = func(self *Action) {
		complete = false
		for _, action := range actions {
			action.Restart(action)
		}
	}
	action.Act = func(self *Action, delta float32) bool {
		if complete {
			return true
		}
		complete = true
		// Pool pool = getPool();
		// setPool(null); // Ensure this action can't be returned to the pool while executing.
		// try {
		for i := 0; i < len(actions) && action.Actor != nil; i++ {
			currentAction := actions[i]
			if currentAction.Actor != nil && !currentAction.Act(currentAction, delta) {
				complete = false
			}
			if action.Actor == nil {
				return true // This action was removed.
			}
		}
		return complete
		// } finally {
		// setPool(pool);
		// }
	}
	return action
}

// Executes a number of actions one at a time.
func Sequence(actions ...*Action) *Action {
	index := 0
	action := NewAction()
	action.Reset = func(self *Action) {
		actions = make([]*Action, 0)
	}
	action.Restart = func(self *Action) {
		for _, action := range actions {
			action.Restart(action)
		}
		index = 0
	}
	action.Act = func(self *Action, delta float32) bool {
		if index >= len(actions) {
			return true
		}
		// Pool pool = getPool();
		// setPool(null); // Ensure this action can't be returned to the pool while executings.
		// try {
		if actions[index].Act(actions[index], delta) {
			// if action.Actor == nil {
			// return true // This action was removed.
			// }
			index++
			if index >= len(actions) {
				return true
			}
		}
		return false
		// } finally {
		// 	setPool(pool);
		// }
	}
	return action
}

// TODO: create pools for each type of action and obtain from them
// dont need to add a reference the pool in each action
// action := pool.Get().(*VisibleAction)

// Sets the actor's visibility.
func Visible(visible bool) *Action {
	action := NewAction() // TODO: get from pool
	action.Act = func(self *Action, delta float32) bool {
		self.Target.Visible = visible
		return true
	}
	return action
}

// Sets an actor's Layout#setLayoutEnabled(boolean) layout to enabled or disabled. The actor must implements ILayout.
func Layout(enabled bool) *Action {
	action := NewAction()
	action.Act = func(self *Action, delta float32) bool {
		// ((Layout)target).setLayoutEnabled(enabled)
		return true
	}
	return action
}

// Base class for an action that wraps another action.
func Delegate(base *Action) (delegated *Action) {
	action := NewAction()
	action.Act = func(self *Action, delta float32) bool {
		base.Act(base, delta)
		return true
	}
	action.Restart = func(self *Action) {
		if base != nil {
			base.Restart(base)
		}
	}
	action.Reset = func(self *Action) {
		base.Reset(base)
		base = nil
	}
	return action
}

// 	public void setActor (Actor actor) {
// 		if (action != null) action.setActor(actor);
// 		super.setActor(actor);
// 	}

// 	public void setTarget (Actor target) {
// 		if (action != null) action.setTarget(target);
// 		super.setTarget(target);
// 	}

// Delays execution of an action or inserts a pause
func Delay(duration float32, delayedAction *Action) *Action {
	action := Delegate(delayedAction)
	time := float32(0)
	action.Act = func(self *Action, delta float32) bool {
		// println(delta)
		println(time)
		if time < duration {
			time += 11
			if time < duration {
				return false
			}
			delta = time - duration
		}
		if delayedAction == nil {
			return true
		}
		return delayedAction.Act(delayedAction, delta)
	}
	return action
}

// Executes an action only after all other actions on the actor at the time this action's target was set have finished.
func After(base *Action) {
	// waitForActions.addAll(target.getActions());
	action := Delegate(base)
	action.Restart = func(self *Action) {
		if base != nil {
			base.Restart(base)
		}
		// waitForActions.clear()
	}
	action.Act = func(self *Action, delta float32) bool {
		// Array<Action> currentActions = target.getActions();
		// if (currentActions.size == 1) waitForActions.clear();
		// for (int i = waitForActions.size - 1; i >= 0; i--) {
		// 	Action action = waitForActions.get(i);
		// 	int index = currentActions.indexOf(action, true);
		// 	if (index == -1) waitForActions.removeIndex(i);
		// }
		// if (waitForActions.size > 0) return false;
		// return base.act(base, delta);
		return true
	}
}

// Sets the actor's Actor#setTouchable(Touchable) touchability.
func TouchableSet(touchable Touchable) *Action {
	action := NewAction()
	action.Act = func(self *Action, delta float32) bool {
		action.Target.TouchState = touchable
		return true
	}
	return action
}

// Multiplies the delta of an action.
func TimeScale(scale float32, scaledAction *Action) *Action {
	action := NewAction()
	action.Act = func(self *Action, delta float32) bool {
		return scaledAction.Act(scaledAction, delta*scale)
	}
	return action
}

// Base class for actions that transition over time using the percent complete.
func Temporal(duration float32, interp Interpolation) *Action {
	duration, time := float32(0), float32(0)
	reverse, began, complete := false, false, false
	action := NewAction()
	action.Act = func(self *Action, delta float32) bool {
		if complete {
			return true
		}
		// 		Pool pool = getPool();
		// 		setPool(null); // Ensure this action can't be returned to the pool while executing.
		// 		try {
		if !began {
			//begin()
			began = true
		}
		time += delta
		complete = time >= duration
		// percent := float32(0)
		// if complete {
		// 	percent = 1
		// } else {
		// 	percent = time / duration
		// 	if interpolation != null {
		// 		percent = interpolation.apply(percent)
		// 	}
		// }
		if reverse {
			// 	update(1 - percent)
		} else {
			// 	update(percent)
		}
		// 			if (complete) end();
		// 			return complete;
		// 		} finally {
		// 			setPool(pool);
		// 		}
		return true
	}
	return action
}

// 	/** Called the first time {@link #act(float)} is called. This is a good place to query the {@link #actor actor's} starting
// 	 * state. */
// 	protected void begin () {
// 	}

// 	/** Called the last time {@link #act(float)} is called. */
// 	protected void end () {
// 	}

// 	* Called each frame.
// 	 * @param percent The percentage of completion for this action, growing from 0 to 1 over the duration. If
// 	 *           {@link #setReverse(boolean) reversed}, this will shrink from 1 to 0.
// 	abstract protected void update (float percent);

// 	/** Skips to the end of the transition. */
// 	public void finish () {
// 		time = duration;
// 	}

// 	public void restart () {
// 		time = 0;
// 		began = false;
// 		complete = false;
// 	}

// 	public void reset () {
// 		super.reset();
// 		reverse = false;
// 		interpolation = null;
// 	}

// 	/** Gets the transition time so far. */
// 	public float getTime () {
// 		return time;
// 	}

// 	/** Sets the transition time so far. */
// 	public void setTime (float time) {
// 		this.time = time;
// 	}

// 	public float getDuration () {
// 		return duration;
// 	}

// 	/** Sets the length of the transition in seconds. */
// 	public void setDuration (float duration) {
// 		this.duration = duration;
// 	}

// 	public Interpolation getInterpolation () {
// 		return interpolation;
// 	}

// 	public void setInterpolation (Interpolation interpolation) {
// 		this.interpolation = interpolation;
// 	}

// 	public boolean isReverse () {
// 		return reverse;
// 	}

// 	/** When true, the action's progress will go from 100% to 0%. */
// 	public void setReverse (boolean reverse) {
// 		this.reverse = reverse;
// 	}
// }

// Moves the actor instantly.
func MoveTo(x, y, duration float32, interp Interpolation) *Action {
	action := NewAction() // TODO: get from pool
	if duration == 0 {
		action.Act = func(self *Action, delta float32) bool {
			self.Target.SetPosition(x, y)
			return true
		}
	} else {
		// action.setDuration(duration)
		// action.setInterpolation(interpolation)
	}
	return action
}

func MoveToAligned(x, y float32, alignment int, duration float32, interp Interpolation) {
	// 		MoveToAction action = action(MoveToAction.class);
	// 		action.setPosition(x, y, alignment);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Moves the actor instantly.
func MoveBy(amountX, amountY, duration float32, interp Interpolation) {
	// 		MoveByAction action = action(MoveByAction.class);
	// 		action.setAmount(amountX, amountY);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Sizes the actor instantly.
func SizeToAction(x, y, duration float32, interp Interpolation) {
	// 		SizeToAction action = action(SizeToAction.class);
	// 		action.setSize(x, y);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Sizes the actor instantly.
func SizeBy(amountX, amountY, duration float32, interp Interpolation) {
	// 		SizeByAction action = action(SizeByAction.class);
	// 		action.setAmount(amountX, amountY);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Scales the actor instantly.
func ScaleTo(x, y, duration float32, interp Interpolation) {
	// 		ScaleToAction action = action(ScaleToAction.class);
	// 		action.setScale(x, y);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Scales the actor instantly.
func ScaleBy(amountX, amountY, duration float32, interp Interpolation) {
	// 		ScaleByAction action = action(ScaleByAction.class);
	// 		action.setAmount(amountX, amountY);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Rotates the actor instantly.
func RotateTo(rotation, duration float32, interp Interpolation) {
	// 		RotateToAction action = action(RotateToAction.class);
	// 		action.setRotation(rotation);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Rotates the actor instantly.
func RotateBy(rotationAmount, duration float32, interp Interpolation) {
	// 		RotateByAction action = action(RotateByAction.class);
	// 		action.setAmount(rotationAmount);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Transitions from the color at the time this action starts to the specified color.
func Color(color graphics.Color, duration float32, interp Interpolation) {
	// 		ColorAction action = action(ColorAction.class);
	// 		action.setEndColor(color);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Transitions from the alpha at the time this action starts to the specified alpha.
func Alpha(a, duration float32, interp Interpolation) {
	// 		AlphaAction action = action(AlphaAction.class);
	// 		action.setAlpha(a);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Transitions from the alpha at the time this action starts to an alpha of 0.
func FadeOut(duration float32, interp Interpolation) {
	// 		AlphaAction action = action(AlphaAction.class);
	// 		action.setAlpha(0);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

// Transitions from the alpha at the time this action starts to an alpha of 1.
func FadeIn(duration float32, interp Interpolation) {
	// 		AlphaAction action = action(AlphaAction.class);
	// 		action.setAlpha(1);
	// 		action.setDuration(duration);
	// 		action.setInterpolation(interpolation);
	// 		return action;
}

func Repeat(count int, repeatedAction *Action) {
	// 		RepeatAction action = action(RepeatAction.class);
	// 		action.setCount(count);
	// 		action.setAction(repeatedAction);
	// 		return action;
}

func Forever(repeatedAction *Action) {
	// 		RepeatAction action = action(RepeatAction.class);
	// 		action.setCount(RepeatAction.FOREVER);
	// 		action.setAction(repeatedAction);
	// 		return action;
}

/** Sets the alpha for an actor's color (or a specified color), from the current alpha to the new alpha. Note this action
 * transitions from the alpha at the time the action starts to the specified alpha.
 * @author Nathan Sweet */
// public class AlphaAction extends TemporalAction {
// 	private float start, end;
// 	private Color color;

// 	protected void begin () {
// 		if (color == null) color = target.getColor();
// 		start = color.a;
// 	}

// 	protected void update (float percent) {
// 		color.a = start + (end - start) * percent;
// 	}

// 	public void reset () {
// 		super.reset();
// 		color = null;
// 	}

// 	public Color getColor () {
// 		return color;
// 	}

// 	/** Sets the color to modify. If null (the default), the {@link #getActor() actor's} {@link Actor#getColor() color} will be
// 	 * used. */
// 	public void setColor (Color color) {
// 		this.color = color;
// 	}

// 	public float getAlpha () {
// 		return end;
// 	}

// 	public void setAlpha (float alpha) {
// 		this.end = alpha;
// 	}
// }

/** Sets the actor's color (or a specified color), from the current to the new color. Note this action transitions from the color
 * at the time the action starts to the specified color.
 * @author Nathan Sweet */
// public class ColorAction extends TemporalAction {
// 	private float startR, startG, startB, startA;
// 	private Color color;
// 	private final Color end = new Color();

// 	protected void begin () {
// 		if (color == null) color = target.getColor();
// 		startR = color.r;
// 		startG = color.g;
// 		startB = color.b;
// 		startA = color.a;
// 	}

// 	protected void update (float percent) {
// 		float r = startR + (end.r - startR) * percent;
// 		float g = startG + (end.g - startG) * percent;
// 		float b = startB + (end.b - startB) * percent;
// 		float a = startA + (end.a - startA) * percent;
// 		color.set(r, g, b, a);
// 	}

// 	public void reset () {
// 		super.reset();
// 		color = null;
// 	}

// 	public Color getColor () {
// 		return color;
// 	}

// 	/** Sets the color to modify. If null (the default), the {@link #getActor() actor's} {@link Actor#getColor() color} will be
// 	 * used. */
// 	public void setColor (Color color) {
// 		this.color = color;
// 	}

// 	public Color getEndColor () {
// 		return end;
// 	}

// 	/** Sets the color to transition to. Required. */
// 	public void setEndColor (Color color) {
// 		end.set(color);
// 	}
// }

/** An action that has a float, whose value is transitioned over time.
 * @author Nathan Sweet */
// public class FloatAction extends TemporalAction {
// 	private float start, end;
// 	private float value;

// 	/** Creates a FloatAction that transitions from 0 to 1. */
// 	public FloatAction () {
// 		start = 0;
// 		end = 1;
// 	}

// 	/** Creates a FloatAction that transitions from start to end. */
// 	public FloatAction (float start, float end) {
// 		this.start = start;
// 		this.end = end;
// 	}

// 	protected void begin () {
// 		value = start;
// 	}

// 	protected void update (float percent) {
// 		value = start + (end - start) * percent;
// 	}

// 	/** Gets the current float value. */
// 	public float getValue () {
// 		return value;
// 	}

// 	/** Sets the current float value. */
// 	public void setValue (float value) {
// 		this.value = value;
// 	}

// 	public float getStart () {
// 		return start;
// 	}

// 	/** Sets the value to transition from. */
// 	public void setStart (float start) {
// 		this.start = start;
// 	}

// 	public float getEnd () {
// 		return end;
// 	}

// 	/** Sets the value to transition to. */
// 	public void setEnd (float end) {
// 		this.end = end;
// 	}
// }

/** An action that has an int, whose value is transitioned over time.
 * @author Nathan Sweet */
// public class IntAction extends TemporalAction {
// 	private int start, end;
// 	private int value;

// 	/** Creates an IntAction that transitions from 0 to 1. */
// 	public IntAction () {
// 		start = 0;
// 		end = 1;
// 	}

// 	/** Creates an IntAction that transitions from start to end. */
// 	public IntAction (int start, int end) {
// 		this.start = start;
// 		this.end = end;
// 	}

// 	protected void begin () {
// 		value = start;
// 	}

// 	protected void update (float percent) {
// 		value = (int)(start + (end - start) * percent);
// 	}

// 	/** Gets the current int value. */
// 	public int getValue () {
// 		return value;
// 	}

// 	/** Sets the current int value. */
// 	public void setValue (int value) {
// 		this.value = value;
// 	}

// 	public int getStart () {
// 		return start;
// 	}

// 	/** Sets the value to transition from. */
// 	public void setStart (int start) {
// 		this.start = start;
// 	}

// 	public int getEnd () {
// 		return end;
// 	}

// 	/** Sets the value to transition to. */
// 	public void setEnd (int end) {
// 		this.end = end;
// 	}
// }

/** Moves an actor to a relative position.
 * @author Nathan Sweet */
// public class MoveByAction extends RelativeTemporalAction {
// 	private float amountX, amountY;

// 	protected void updateRelative (float percentDelta) {
// 		target.moveBy(amountX * percentDelta, amountY * percentDelta);
// 	}

// 	public void setAmount (float x, float y) {
// 		amountX = x;
// 		amountY = y;
// 	}

// 	public float getAmountX () {
// 		return amountX;
// 	}

// 	public void setAmountX (float x) {
// 		amountX = x;
// 	}

// 	public float getAmountY () {
// 		return amountY;
// 	}

// 	public void setAmountY (float y) {
// 		amountY = y;
// 	}
// }

/** Moves an actor from its current position to a specific position.
 * @author Nathan Sweet */
// public class MoveToAction extends TemporalAction {
// 	private float startX, startY;
// 	private float endX, endY;
// 	private int alignment = Align.bottomLeft;

// 	protected void begin () {
// 		startX = target.getX(alignment);
// 		startY = target.getY(alignment);
// 	}

// 	protected void update (float percent) {
// 		target.setPosition(startX + (endX - startX) * percent, startY + (endY - startY) * percent, alignment);
// 	}

// 	public void reset () {
// 		super.reset();
// 		alignment = Align.bottomLeft;
// 	}

// 	public void setPosition (float x, float y) {
// 		endX = x;
// 		endY = y;
// 	}

// 	public void setPosition (float x, float y, int alignment) {
// 		endX = x;
// 		endY = y;
// 		this.alignment = alignment;
// 	}

// 	public float getX () {
// 		return endX;
// 	}

// 	public void setX (float x) {
// 		endX = x;
// 	}

// 	public float getY () {
// 		return endY;
// 	}

// 	public void setY (float y) {
// 		endY = y;
// 	}

// 	public int getAlignment () {
// 		return alignment;
// 	}

// 	public void setAlignment (int alignment) {
// 		this.alignment = alignment;
// 	}
// }

/** Base class for actions that transition over time using the percent complete since the last frame.
 * @author Nathan Sweet */
// abstract public class RelativeTemporalAction extends TemporalAction {
// 	private float lastPercent;

// 	protected void begin () {
// 		lastPercent = 0;
// 	}

// 	protected void update (float percent) {
// 		updateRelative(percent - lastPercent);
// 		lastPercent = percent;
// 	}

// 	abstract protected void updateRelative (float percentDelta);
// }

/** Repeats an action a number of times or forever.
 * @author Nathan Sweet */
// public class RepeatAction extends DelegateAction {
// 	static public final int FOREVER = -1;

// 	private int repeatCount, executedCount;
// 	private boolean finished;

// 	protected boolean delegate (float delta) {
// 		if (executedCount == repeatCount) return true;
// 		if (action.act(delta)) {
// 			if (finished) return true;
// 			if (repeatCount > 0) executedCount++;
// 			if (executedCount == repeatCount) return true;
// 			if (action != null) action.restart();
// 		}
// 		return false;
// 	}

// 	/** Causes the action to not repeat again. */
// 	public void finish () {
// 		finished = true;
// 	}

// 	public void restart () {
// 		super.restart();
// 		executedCount = 0;
// 		finished = false;
// 	}

// 	/** Sets the number of times to repeat. Can be set to {@link #FOREVER}. */
// 	public void setCount (int count) {
// 		this.repeatCount = count;
// 	}

// 	public int getCount () {
// 		return repeatCount;
// 	}
// }

/** Sets the actor's rotation from its current value to a relative value.
 * @author Nathan Sweet */
// public class RotateByAction extends RelativeTemporalAction {
// 	private float amount;

// 	protected void updateRelative (float percentDelta) {
// 		target.rotateBy(amount * percentDelta);
// 	}

// 	public float getAmount () {
// 		return amount;
// 	}

// 	public void setAmount (float rotationAmount) {
// 		amount = rotationAmount;
// 	}
// }

/** Sets the actor's rotation from its current value to a specific value.
 * @author Nathan Sweet */
// public class RotateToAction extends TemporalAction {
// 	private float start, end;

// 	protected void begin () {
// 		start = target.getRotation();
// 	}

// 	protected void update (float percent) {
// 		target.setRotation(start + (end - start) * percent);
// 	}

// 	public float getRotation () {
// 		return end;
// 	}

// 	public void setRotation (float rotation) {
// 		this.end = rotation;
// 	}
// }

/** Scales an actor's scale to a relative size.
 * @author Nathan Sweet */
// public class ScaleByAction extends RelativeTemporalAction {
// 	private float amountX, amountY;

// 	protected void updateRelative (float percentDelta) {
// 		target.scaleBy(amountX * percentDelta, amountY * percentDelta);
// 	}

// 	public void setAmount (float x, float y) {
// 		amountX = x;
// 		amountY = y;
// 	}

// 	public void setAmount (float scale) {
// 		amountX = scale;
// 		amountY = scale;
// 	}

// 	public float getAmountX () {
// 		return amountX;
// 	}

// 	public void setAmountX (float x) {
// 		this.amountX = x;
// 	}

// 	public float getAmountY () {
// 		return amountY;
// 	}

// 	public void setAmountY (float y) {
// 		this.amountY = y;
// 	}

// }

/** Sets the actor's scale from its current value to a specific value.
 * @author Nathan Sweet */
// public class ScaleToAction extends TemporalAction {
// 	private float startX, startY;
// 	private float endX, endY;

// 	protected void begin () {
// 		startX = target.getScaleX();
// 		startY = target.getScaleY();
// 	}

// 	protected void update (float percent) {
// 		target.setScale(startX + (endX - startX) * percent, startY + (endY - startY) * percent);
// 	}

// 	public void setScale (float x, float y) {
// 		endX = x;
// 		endY = y;
// 	}

// 	public void setScale (float scale) {
// 		endX = scale;
// 		endY = scale;
// 	}

// 	public float getX () {
// 		return endX;
// 	}

// 	public void setX (float x) {
// 		this.endX = x;
// 	}

// 	public float getY () {
// 		return endY;
// 	}

// 	public void setY (float y) {
// 		this.endY = y;
// 	}
// }

/** Moves an actor from its current size to a relative size.
 * @author Nathan Sweet */
// public class SizeByAction extends RelativeTemporalAction {
// 	private float amountWidth, amountHeight;

// 	protected void updateRelative (float percentDelta) {
// 		target.sizeBy(amountWidth * percentDelta, amountHeight * percentDelta);
// 	}

// 	public void setAmount (float width, float height) {
// 		amountWidth = width;
// 		amountHeight = height;
// 	}

// 	public float getAmountWidth () {
// 		return amountWidth;
// 	}

// 	public void setAmountWidth (float width) {
// 		amountWidth = width;
// 	}

// 	public float getAmountHeight () {
// 		return amountHeight;
// 	}

// 	public void setAmountHeight (float height) {
// 		amountHeight = height;
// 	}
// }

/** Moves an actor from its current size to a specific size.
 * @author Nathan Sweet */
// public class SizeToAction extends TemporalAction {
// 	private float startWidth, startHeight;
// 	private float endWidth, endHeight;

// 	protected void begin () {
// 		startWidth = target.getWidth();
// 		startHeight = target.getHeight();
// 	}

// 	protected void update (float percent) {
// 		target.setSize(startWidth + (endWidth - startWidth) * percent, startHeight + (endHeight - startHeight) * percent);
// 	}

// 	public void setSize (float width, float height) {
// 		endWidth = width;
// 		endHeight = height;
// 	}

// 	public float getWidth () {
// 		return endWidth;
// 	}

// 	public void setWidth (float width) {
// 		endWidth = width;
// 	}

// 	public float getHeight () {
// 		return endHeight;
// 	}

// 	public void setHeight (float height) {
// 		endHeight = height;
// 	}
// }
