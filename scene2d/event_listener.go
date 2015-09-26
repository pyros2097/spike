// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package scene2d

// Low level interface for receiving events. Typically there is a listener class for each specific event class.
// @see InputListener
// @see InputEvent
type EventListener interface {
	// Try to handle the given event, if it is applicable.
	// @return true if the event should be considered {@link Event#handle() handled} by scene2d.
	Handle(event *Event) bool
}
