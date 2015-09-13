// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package scene2d

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

func (self *Action) toString() string {
	return "Action"
}
