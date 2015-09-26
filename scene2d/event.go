// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package scene2d

/** The base class for all events.
 * <p>
 * By default an event will "bubble" up through an actor's parent's handlers (see {@link #setBubbles(boolean)}).
 * <p>
 * An actor's capture listeners can {@link #stop()} an event to prevent child actors from seeing it.
 * <p>
 * An Event may be marked as "handled" which will end its propagation outside of the Stage (see {@link #handle()}). The default
 * {@link Actor#fire(Event)} will mark events handled if an {@link EventListener} returns true.
 * <p>
 * A cancelled event will be stopped and handled. Additionally, many actors will undo the side-effects of a canceled event. (See
 * {@link #cancel()}.)
 *
 * @see InputEvent
 * @see Actor#fire(Event) */
type Event struct {
	targetActor   *Actor
	listenerActor *Actor

	// true means event occurred during the capture phase
	capture bool

	// true means propagate to target's parents
	bubbles bool

	// true means the event was handled (the stage will eat the input)
	handled bool

	// true means event propagation was stopped
	stopped bool

	// true means propagation was stopped and any action that this event would cause should not happen
	cancelled bool
}

func NewEvent() *Event {
	return &Event{
		capture:   false,
		bubbles:   true,
		handled:   false,
		stopped:   false,
		cancelled: false,
	}
}

/** Marks this event as handled. This does not affect event propagation inside scene2d, but causes the {@link Stage} event
 * methods to return true, which will eat the event so it is not passed on to the application under the stage. */
func (self *Event) handle() {
	self.handled = true
}

/** Marks this event cancelled. This {@link #handle() handles} the event and {@link #stop() stops} the event propagation. It
 * also cancels any default action that would have been taken by the code that fired the event. Eg, if the event is for a
 * checkbox being checked, cancelling the event could uncheck the checkbox. */
func (self *Event) cancel() {
	self.cancelled = true
	self.stopped = true
	self.handled = true
}

/** Marks this event has being stopped. This halts event propagation. Any other listeners on the {@link #getListenerActor()
 * listener actor} are notified, but after that no other listeners are notified. */
func (self *Event) stop() {
	self.stopped = true
}

func (self *Event) reset() {
	// self.stage = nil
	self.targetActor = nil
	self.listenerActor = nil
	self.capture = false
	self.bubbles = true
	self.handled = false
	self.stopped = false
	self.cancelled = false
}

/** Returns the actor that the event originated from. */
func (self *Event) GetTarget() *Actor {
	return self.targetActor
}

func (self *Event) SetTarget(targetActor *Actor) {
	self.targetActor = targetActor
}

/** Returns the actor that this listener is attached to. */
func (self *Event) GetListenerActor() *Actor {
	return self.listenerActor
}

func (self *Event) SetListenerActor(listenerActor *Actor) {
	self.listenerActor = listenerActor
}

func (self *Event) GetBubbles() bool {
	return self.bubbles
}

/** If true, after the event is fired on the target actor, it will also be fired on each of the parent actors, all the way to
 * the root. */
func (self *Event) SetBubbles(bubbles bool) {
	self.bubbles = bubbles
}

/** {@link #handle()} */
func (self *Event) isHandled() bool {
	return self.handled
}

/** @see #stop() */
func (self *Event) isStopped() bool {
	return self.stopped
}

/** @see #cancel() */
func (self *Event) isCancelled() bool {
	return self.cancelled
}

func (self *Event) setCapture(capture bool) {
	self.capture = capture
}

/** If true, the event was fired during the capture phase.
 * @see Actor#fire(Event) */
func (self *Event) isCapture() bool {
	return self.capture
}
