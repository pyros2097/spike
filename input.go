// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"math"
	"time"

	"github.com/pyros2097/spike/math/vector"
)

// Interface to the input facilities. This allows polling the state of the keyboard, the touch screen and the accelerometer.
// On some backends (desktop) the touch screen is replaced by mouse input. The accelerometer is of course not available on
// all backends.
// Keyboard keys are translated to the constants in {@link Keys} transparently on all systems. Do not use system specific key
// constants.
// This also offers methods to use (and test for the presence of) other input systems like vibration, compass, on-screen
// keyboards, and cursor capture. Support for simple input dialogs is also provided.

// The Types of Input events that can occurr.
type InputType uint8

const (
	None InputType = iota

	// A new touch for a pointer on the stage was detected
	// Called when the screen was touched or a mouse button was pressed. The button parameter will be {@link Buttons#LEFT} on iOS.
	// param screenX The x coordinate, origin is in the upper left corner
	// param screenY The y coordinate, origin is in the upper left corner
	// param pointer the pointer for the event.
	// param button the button
	// return whether the input was processed
	// Called when a mouse button or a finger touch goes down on the actor. If true is returned, this listener will receive all
	// touchDragged and touchUp events, even those not over this actor, until touchUp is received. Also when true is returned, the
	// event is {@link Event#handle() handled}.
	TouchDown

	// A pointer has stopped touching the stage.
	// Called when a finger was lifted or a mouse button was released. The button parameter will be {@link Buttons#LEFT} on iOS.
	// param pointer the pointer for the event.
	// param button the button
	// return whether the input was processed
	// Called when a mouse button or a finger touch goes up anywhere, but only if touchDown previously returned true for the mouse
	// button or touch. The touchUp event is always {@link Event#handle() handled}.
	TouchUp

	// A pointer that is touching the stage has moved.
	// Called when a finger or the mouse was dragged.
	// param pointer the pointer for the event.
	// @return whether the input was processed
	// Called when a mouse button or a finger touch is moved anywhere, but only if touchDown previously returned true for the mouse
	// button or touch. The touchDragged event is always {@link Event#handle() handled}.
	TouchDragged

	// The mouse pointer has moved (without a mouse button being active).
	MouseMoved

	// The mouse pointer or an active touch have entered (i.e., {@link Actor#hit(float, float, boolean) hit}) an actor.
	Enter

	// The mouse pointer or an active touch have exited an actor.
	Exit

	// The mouse scroll wheel has changed.
	Scrolled

	// Called when a key was pressed
	// param keycode one of the constants in {@link Input.Keys}
	// return whether the input was processed
	KeyDown

	// A keyboard key has been released.
	// A keyboard key has been pressed.
	// Called when a key was released
	// param keycode one of the constants in {@link Input.Keys}
	// return whether the input was processed
	// When true is returned, the event is {@link Event#handle() handled
	KeyUp

	// A keyboard key has been pressed and released.
	// Called when a key was typed
	// param character The character
	// return whether the input was processed
	KeyTyped

	// Called when a tap occured. A tap happens if a touch went down on the screen and was lifted again without moving outside
	// of the tap square. The tap square is a rectangular area around the initial touch position as specified on construction
	// time of the {@link GestureDetector}.
	// @param count the number of taps.
	Tap

	// Called when the user dragged a finger over the screen and lifted it. Reports the last known velocity of the finger in
	// pixels per second.
	// param velocityX velocity on x in seconds
	// param velocityY velocity on y in seconds
	Fling

	// Called when the user drags a finger over the screen.
	// param deltaX the difference in pixels to the last drag event on x.
	// param deltaY the difference in pixels to the last drag event on y.
	Pan

	// Called when no longer panning.
	PanStop

	// Called when a long press was detected
	LongPress

	// Called when the user performs a pinch zoom gesture. The original distance is the distance in pixels when the gesture
	// started.
	// param initialDistance distance between fingers when the gesture started.
	// param distance current distance between fingers.
	Zoom

	// Called when a user performs a pinch zoom gesture. Reports the initial positions of the two involved fingers and their
	// current positions.
	// param initialPointer1
	// param initialPointer2
	// param pointer1
	// param pointer2
	Pinch

	// Called when a swipe gesture occurs
	SwipeLeft
	SwipeRight
	SwipeUp
	SwipeDown
)

// The type of Mouse buttons
type Button uint8

const (
	ButtonLeft Button = iota
	ButtonRight
	ButtonMiddle
	ButtonBack
	ButtonForward
)

// Register an instance of this class with a {@link GestureDetector} to receive gestures such as taps, long presses, flings,
// panning or pinch zooming. Each method returns a boolean indicating if the event should be handed to the next listener (false
// to hand it to the next listener, true otherwise).

// Called when the mouse was moved without any buttons being pressed. Will not be called on iOS.
// @return whether the input was processed
// This event only occurs on the desktop
// public boolean mouseMoved (int screenX, int screenY);

// Called when the mouse wheel was scrolled. Will not be called on iOS.
// param amount the scroll amount, -1 or 1 depending on the direction the wheel was scrolled.
// @return whether the input was processed.
// public boolean scrolled (int amount);

/** Called any time the mouse cursor or a finger touch is moved over an actor. On the desktop, this event occurs even when no
 * mouse buttons are pressed (pointer will be -1).
 * @param fromActor May be null.
 * @see InputEvent */
// public void enter (InputEvent event, float x, float y, int pointer, Actor fromActor) {
// }

/** Called any time the mouse cursor or a finger touch is moved out of an actor. On the desktop, this event occurs even when no
 * mouse buttons are pressed (pointer will be -1).
 * @param toActor May be null.
 * @see InputEvent */
// public void exit (InputEvent event, float x, float y, int pointer, Actor toActor) {
// }

/** Called when the mouse wheel has been scrolled. When true is returned, the event is {@link Event#handle() handled}. */
// public boolean scrolled (InputEvent event, float x, float y, int amount) {
//  return false;
// }

//  Event for actor input: touch, mouse, keyboard, and scroll.
type InputEvent struct {
	// The type of input event.
	Type InputType

	// The stage x coordinate where the event occurred. Valid for: touchDown, touchDragged, touchUp, mouseMoved, enter, and exit.
	X float32

	// The stage x coordinate where the event occurred. Valid for: touchDown, touchDragged, touchUp, mouseMoved, enter, and exit.
	Y float32

	// The pointer index for the event. The first touch is index 0, second touch is index 1, etc. Always -1 on desktop. Valid for:
	// touchDown, touchDragged, touchUp, enter, and exit.
	Pointer uint8

	// The index for the mouse button pressed. Always 0 on Android. Valid for: touchDown and touchUp.
	Button uint8

	// The key code of the key that was pressed. Valid for: keyDown and keyUp.
	KeyCode int

	// The character for the key that was type. Valid for: keyTyped.
	Character uint8

	// The amount the mouse was scrolled. Valid for: scrolled.
	ScrollAmount int

	// The actor related to the event. Valid for: enter and exit. For enter, this is the actor being exited, or null.
	// For exit, this is the actor being entered, or null.
	// RelatedActor *scene2d.Actor
}

func (self *InputEvent) reset() {
	//    super.reset();
	// RelatedActor = nil
	// self.Button = 255
}

func (input InputType) String() string {
	switch input {
	case TouchDown:
		return "TouchDown"
	case TouchUp:
		return "TouchUp"
	case TouchDragged:
		return "TouchDragged"
	case KeyDown:
		return "KeyDown"
	case KeyUp:
		return "KeyUp"
	case KeyTyped:
		return "KeyTyped"
	case Tap:
		return "Tap"
	case LongPress:
		return "LongPress"
	case Fling:
		return "Fling"
	case Zoom:
		return "Zoom"
	case Pan:
		return "Pan"
	case PanStop:
		return "PanStop"
	case SwipeUp:
		return "SwipeUp"
	case SwipeDown:
		return "SwipeDown"
	case SwipeLeft:
		return "SwipeLeft"
	case SwipeRight:
		return "SwipeRight"
	case None:
		return "None"
	default:
		return "Input"
	}
}

//  /** Sets actorCoords to this event's coordinates relative to the specified actor.
//   * @param actorCoords Output for resulting coordinates. */
//  public Vector2 toCoordinates (Actor actor, Vector2 actorCoords) {
//    actorCoords.set(stageX, stageY);
//    actor.stageToLocalCoordinates(actorCoords);
//    return actorCoords;
//  }

// Returns true of this event is a touchUp triggered by {@link Stage#cancelTouchFocus()}. */
// func (self *GestureEvent) IsTouchFocusCancel() bool {
//  return self.X == math.MinInt32 || self.Y == math.MinInt32
// }

var (
	TapSquareSize                      float32 = 20
	tapCountInterval                   int64   = 0 //0.4f
	LongPressSeconds                   float32 = 1.1
	MaxFlingDelay                      int64   = 0 //0.15
	inTapSquare                                = false
	tapCount                           int     = 0
	lastTapTime                        int64   = 0
	lastTapX, lastTapY                 float32
	lastTapButton, lastTapPointer      int
	longPressFired, pinching, panning  bool
	tapSquareCenterX, tapSquareCenterY float32
	gestureStartTime                   int64
	longPressScheduled                 = false
	longPressTask                      = time.AfterFunc(time.Duration(float32(time.Second)*LongPressSeconds), fireLongPress)
	fireLongPress                      = func() {
		if !longPressFired {
			longPressFired = true
			InputChannel <- InputEvent{
				Type: LongPress,
				X:    pointer1.X,
				Y:    pointer1.Y,
			}
		}
	}
	pointer1        = vector.NewVector2Empty()
	pointer2        = vector.NewVector2Empty()
	initialPointer1 = vector.NewVector2Empty()
	initialPointer2 = vector.NewVector2Empty()
	event           InputEvent
	queue           []InputEvent
	InputChannel    = make(chan InputEvent, 100)

	// Gesture related
	gestureStarted                       = false
	touchDragIntervalRange       float32 = 200 // 200px drag for a gesture event
	touchInitialX, touchInitialY float32
	touchCurrentX, touchCurrentY float32
	difX, difY                   float32
	prevDifX, prevDifY           float32
	touchDifX, touchDifY         float32
)

// Detects gestures (tap, long press, fling, pan, zoom, pinch) and hands them to a listener
// param halfTapSquareSize half width in pixels of the square around an initial touch event, see
//          {@link GestureListener#tap(float, float, int, int)}.
// param tapCountInterval time in seconds that must pass for two touch down/up sequences to be detected as consecutive taps.
// param longPressDuration time in seconds that must pass for the detector to fire a
//          {@link GestureListener#longPress(float, float)} event.
// param MaxFlingDelay time in seconds the finger must have been dragged for a fling event to be fired, see
//          {@link GestureListener#fling(float, float, int)}
// param listener May be null if the listener will be set later.
// tapCountInterval=0.4f

func isWithinTapSquare(x, y, centerX, centerY float32) bool {
	return float32(math.Abs(float64(x-centerX))) < TapSquareSize && float32(math.Abs(float64(y-centerY))) < TapSquareSize
}

// Todo: use pools, and check if pointer is ok
func GetEvent() *InputEvent {
	if len(queue) > 0 {
		event = queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		return &event
	}
	return nil
}

// /* Dragging Camera Related */
// static Actor validActor = null;
// private final ClickListener touchInput = new ClickListener(){
//  @Override
//  public void clicked(InputEvent event, float x, float y){
//    super.clicked(event, x, y);
//    mouse.set(x, y);
//    validActor = hit(x,y);
//    if(validActor != null && validActor.getName() != null)
//      Scene.getCurrentScene().onClick(validActor);
//  }

func doTouchDown(x, y float32, pointer, button int) bool {
	InputChannel <- InputEvent{
		Type:    TouchDown,
		X:       x,
		Y:       y,
		Pointer: 0,
		Button:  0,
	}
	if pointer == 0 {
		pointer1.Set(x, y)
		gestureStartTime = time.Now().UnixNano() //Gdx.input.getCurrentEventTime()
		velocityStart(x, y, gestureStartTime)
		inTapSquare = true
		pinching = false
		longPressFired = false
		tapSquareCenterX = x
		tapSquareCenterY = y
		if !longPressScheduled {
			longPressTask = time.AfterFunc(time.Duration(float32(time.Second)*LongPressSeconds), fireLongPress)
		}
	} else {
		// Start pinch.
		pointer2.Set(x, y)
		inTapSquare = false
		pinching = true
		initialPointer1.SetV(pointer1)
		initialPointer2.SetV(pointer2)
		longPressTask.Stop()
		longPressScheduled = false
	}
	//    mouse.set(x, y);
	//    mousePointer = pointer;
	//    mouseButton = button;
	touchInitialX = x
	touchInitialY = y
	gestureStarted = true
	//    validActor = hit(x,y);
	//    if(validActor != null && validActor.getName() != null)
	//      Scene.getCurrentScene().onTouchDown(validActor);
	//    return super.touchDown(event, x, y, pointer, button);
	return false
}

func doTouchUp(x, y float32, pointer, button int) bool {
	InputChannel <- InputEvent{
		Type: TouchUp,
		X:    x,
		Y:    y,
	}
	// check if we are still tapping.
	if inTapSquare && !isWithinTapSquare(x, y, tapSquareCenterX, tapSquareCenterY) {
		inTapSquare = false
	}

	wasPanning := panning
	panning = false

	longPressTask.Stop()
	longPressScheduled = false
	if longPressFired {
		return false
	}

	if inTapSquare {
		// handle taps
		if lastTapButton != button || lastTapPointer != pointer || time.Now().UnixNano()-lastTapTime > tapCountInterval ||
			!isWithinTapSquare(x, y, lastTapX, lastTapY) {
			tapCount = 0
		}
		tapCount++
		lastTapTime = time.Now().UnixNano()
		lastTapX = x
		lastTapY = y
		lastTapButton = button
		lastTapPointer = pointer
		gestureStartTime = 0
		// queue = append(queue, InputEvent{
		// 	Type:    Tap,
		// 	X:       x,
		// 	Y:       y,
		// 	Pointer: 0,
		// 	Button:  0,
		// })
		InputChannel <- InputEvent{
			Type: Tap,
			X:    x,
			Y:    y,
		}
		return true
	}

	if pinching {
		// handle pinch end
		pinching = false
		panning = true
		// we are in pan mode again, reset velocity tracker
		if pointer == 0 {
			// first pointer has lifted off, set up panning to use the second pointer...
			velocityStart(pointer2.X, pointer2.Y, time.Now().UnixNano()) //Gdx.input.getCurrentEventTime())
		} else {
			// second pointer has lifted off, set up panning to use the first pointer...
			velocityStart(pointer1.X, pointer1.Y, time.Now().UnixNano()) //Gdx.input.getCurrentEventTime())
		}
		return false
	}

	// handle no longer panning
	handled := false
	if wasPanning && !panning {
		InputChannel <- InputEvent{
			Type:    PanStop,
			X:       x,
			Y:       y,
			Pointer: 0,
			Button:  0,
		}
		handled = false
	}

	// handle fling
	gestureStartTime = 0
	time := time.Now().UnixNano() //Gdx.input.getCurrentEventTime();
	if time-LastTime < MaxFlingDelay {
		velocityUpdate(x, y, time)
		InputChannel <- InputEvent{
			Type:   Fling,
			X:      getVelocityX(),
			Y:      getVelocityY(),
			Button: 0,
		}
	}
	// reset Gesture
	difX = 0.0
	difY = 0.0
	touchInitialX = 0.0
	touchInitialY = 0.0
	touchCurrentX = 0.0
	touchCurrentY = 0.0
	touchDifX = 0.0
	touchDifY = 0.0
	prevDifX = 0.0
	prevDifY = 0.0
	gestureStarted = false
	// reset Gesture
	// Camera.resetDrag();

	return handled
}

func doTouchDragged(x, y float32, pointer int) bool {
	InputChannel <- InputEvent{
		Type: TouchDragged,
		X:    x,
		Y:    y,
	}
	if pointer > 1 {
		return false
	}
	if longPressFired {
		return false
	}

	if pointer == 0 {
		pointer1.Set(x, y)
	} else {
		pointer2.Set(x, y)
	}

	// // handle pinch zoom
	if pinching {
		//  if (listener != null) {
		//    boolean result = listener.pinch(initialPointer1, initialPointer2, pointer1, pointer2);
		//    return listener.zoom(initialPointer1.dst(initialPointer2), pointer1.dst(pointer2)) || result;
		//  }
		//  return false;
	}

	// update tracker
	velocityUpdate(x, y, time.Now().UnixNano()) //Gdx.input.getCurrentEventTime())

	// // check if we are still tapping.
	if inTapSquare && !isWithinTapSquare(x, y, tapSquareCenterX, tapSquareCenterY) {
		longPressTask.Stop()
		longPressScheduled = false
		inTapSquare = false
	}

	// if we have left the tap square, we are panning
	if !inTapSquare {
		panning = true
		InputChannel <- InputEvent{
			Type: Pan,
			X:    deltaX,
			Y:    deltaY,
		}
	}
	if gestureStarted == true {
		touchCurrentX = x
		touchCurrentY = y
		if dir := getGestureDirection(); dir != None {
			gestureStarted = false
			InputChannel <- InputEvent{
				Type: dir,
				X:    x,
				Y:    y,
			}
		}
	}
	//    if(Camera.useDrag)
	//      Camera.dragCam((int)x, (int)-y);

	return false
}

// No further gesture events will be triggered for the current touch, if any.
func Cancel() {
	longPressTask.Stop()
	longPressFired = true
}

// return whether the user touched the screen long enough to trigger a long press event.
func IsLongPressed() bool {
	return isLongPressedDuration(LongPressSeconds)
}

// param duration
// return whether the user touched the screen for as much or more than the given duration.
func isLongPressedDuration(duration float32) bool {
	if gestureStartTime == 0 {
		return false
	}
	return time.Now().UnixNano()-gestureStartTime > int64(duration*1000000000)
}

func IsPanning() bool {
	return panning
}

func Reset() {
	gestureStartTime = 0
	panning = false
	inTapSquare = false
}

// The tap square will not longer be used for the current touch
func InvalidateTapSquare() {
	inTapSquare = false
}

// tapCountInterval time in seconds that must pass for two touch down/up sequences to be detected as consecutive taps
func SetTapCountInterval(tapCountInterval float32) {
	tapCountInterval = tapCountInterval * 1000000000
}

var (
	sampleSize     int = 10
	lastX, lastY   float32
	deltaX, deltaY float32
	LastTime       int64
	numSamples     int
	meanX          = make([]float32, 10)
	meanY          = make([]float32, 10)
	meanTime       = make([]int64, 10)
)

func velocityStart(x, y float32, timeStamp int64) {
	lastX = x
	lastY = y
	deltaX = 0
	deltaY = 0
	numSamples = 0
	for i := 0; i < sampleSize; i++ {
		meanX[i] = 0
		meanY[i] = 0
		meanTime[i] = 0
	}
	LastTime = timeStamp
}

func velocityUpdate(x, y float32, timeStamp int64) {
	currTime := timeStamp
	deltaX = x - lastX
	deltaY = y - lastY
	lastX = x
	lastY = y
	deltaTime := currTime - LastTime
	LastTime = currTime
	index := numSamples % sampleSize
	meanX[index] = deltaX
	meanY[index] = deltaY
	meanTime[index] = deltaTime
	numSamples++
}

func getVelocityX() float32 {
	meanX := getAverage(meanX, numSamples)
	meanTime := getAverageInt(meanTime, numSamples) / 1000000000.0
	if meanTime == 0 {
		return 0
	}
	return meanX / float32(meanTime)
}

func getVelocityY() float32 {
	meanY := getAverage(meanY, numSamples)
	meanTime := getAverageInt(meanTime, numSamples) / 1000000000.0
	if meanTime == 0 {
		return 0
	}
	return meanY / float32(meanTime)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func getAverage(values []float32, numSamples int) float32 {
	numSamples = min(sampleSize, numSamples)
	var sum float32
	for i := 0; i < numSamples; i++ {
		sum += values[i]
	}
	return sum / float32(numSamples)
}

func getAverageInt(values []int64, numSamples int) int64 {
	numSamples = min(sampleSize, numSamples)
	var sum int64
	for i := 0; i < numSamples; i++ {
		sum += values[i]
	}
	if numSamples == 0 {
		return 0
	}
	return sum / int64(numSamples)
}

func getSum(values []float32, numSamples int) float32 {
	numSamples = min(sampleSize, numSamples)
	var sum float32
	for i := 0; i < numSamples; i++ {
		sum += values[i]
	}
	if numSamples == 0 {
		return 0
	}
	return sum
}

func isTouchDragInterval() bool {
	difX = float32(math.Abs(float64(touchInitialX - touchCurrentX)))
	difY = float32(math.Abs(float64(touchInitialY - touchCurrentY)))
	if difX > touchDragIntervalRange || difY > touchDragIntervalRange {
		return true
	}
	return false
}

func getGestureDirection() InputType {
	prevDifX = difX
	prevDifY = difY
	difX = float32(math.Abs(float64(touchInitialX - touchCurrentX)))
	difY = float32(math.Abs(float64(touchInitialY - touchCurrentY)))
	//
	// Get minimal changes on drag
	// <p> checkMomentumChanges
	// EXAMPLE<br>
	// User drags finger to left, suddenly dragging to right without removing
	// his finger from the screen
	//
	if prevDifX > difX || prevDifY > difY {
		touchInitialX = touchCurrentX
		touchInitialY = touchCurrentY
		//
		difX = 0.0
		difY = 0.0
		prevDifX = 0.0
		prevDifY = 0.0

		// Set touch differences, optional, if you need amount of change from
		// initial touch to drag, USE THIS: on touchDrag, pan or similar mthods
		touchDifX = float32(math.Abs(float64(touchInitialX - touchCurrentX)))
		touchDifY = float32(math.Abs(float64(touchInitialY - touchCurrentY)))
	}
	switch {
	case touchInitialY < touchCurrentY && difY > difX:
		return SwipeDown
	case touchInitialY > touchCurrentY && difY > difX:
		return SwipeUp
	case touchInitialX < touchCurrentX && difY < difX:
		return SwipeRight
	case touchInitialX > touchCurrentX && difY < difX:
		return SwipeLeft
	default:
		return None
	}
}

// type Orientation int

// const (
// 	Landscape Orientation = 0
// 	Portrait  Orientation = 1
// )

// // Callback interface for {@link Input#getTextInput(TextInputListener, String, String, String)}
// func TextInputListener(text string) {}
// func TextInputCanceled()            {}

// // Enumeration of potentially available peripherals. Use with {@link Input#isPeripheralAvailable(Peripheral)}.
// // public enum Peripheral {
// //   HardwareKeyboard, OnscreenKeyboard, MultitouchScreen, Accelerometer, Compass, Vibrator
// // }

// // return The value of the accelerometer on its x-axis. ranges between [-10,10].
// func GetAccelerometerX() float32 {
// }

// // return The value of the accelerometer on its y-axis. ranges between [-10,10].
// func GetAccelerometerY() float32 {
// }

// // return The value of the accelerometer on its y-axis. ranges between [-10,10].
// func GetAccelerometerZ() float32 {
// }

// // return The x coordinate of the last touch on touch screen devices and the current mouse position on desktop for the first
// // pointer in screen coordinates. The screen origin is the top left corner.
// func GetX() int {
// }

// // Returns the x coordinate in screen coordinates of the given pointer. Pointers are indexed from 0 to n. The pointer id
// // identifies the order in which the fingers went down on the screen, e.g. 0 is the first finger, 1 is the second and so on.
// // When two fingers are touched down and the first one is lifted the second one keeps its index. If another finger is placed on
// // the touch screen the first free index will be used.
// //
// // param pointer the pointer id.
// // return the x coordinate
// func GetXPointer(pointer int) int {
// }

// // return the different between the current pointer location and the last pointer location on the x-axis.
// func GetDeltaX() int {
// }

// // return the different between the current pointer location and the last pointer location on the x-axis.
// func GetDeltaXPointer(pointer int) int {
// }

// // return The y coordinate of the last touch on touch screen devices and the current mouse position on desktop for the first
// // pointer in screen coordinates. The screen origin is the top left corner.
// func GetY() int {
// }

// // Returns the y coordinate in screen coordinates of the given pointer. Pointers are indexed from 0 to n. The pointer id
// // identifies the order in which the fingers went down on the screen, e.g. 0 is the first finger, 1 is the second and so on.
// // When two fingers are touched down and the first one is lifted the second one keeps its index. If another finger is placed on
// // the touch screen the first free index will be used.
// //
// // param pointer the pointer id.
// // return the y coordinate
// func GetYPointer(pointer int) int {
// }

// // return the different between the current pointer location and the last pointer location on the y-axis.
// func GetDeltaY() int {
// }

// // return the different between the current pointer location and the last pointer location on the y-axis.
// func GetDeltaYPointer(pointer int) int {
// }

// // return whether the screen is currently touched.
// func IsTouched() bool {
// }

// // return whether a new touch down event just occurred.
// func JustTouched() bool {
// }

// // Whether the screen is currently touched by the pointer with the given index. Pointers are indexed from 0 to n. The pointer
// // id identifies the order in which the fingers went down on the screen, e.g. 0 is the first finger, 1 is the second and so on.
// // When two fingers are touched down and the first one is lifted the second one keeps its index. If another finger is placed on
// // the touch screen the first free index will be used.
// //
// // param pointer the pointer
// // return whether the screen is touched by the pointer
// func IsTouchedPointer(pointer int) bool {
// }

// // Whether a given button is pressed or not. Button constants can be found in {@link Buttons}. On Android only the Button#LEFT
// // constant is meaningful before version 4.0.
// // param button the button to check.
// // return whether the button is down or not.
// func IsButtonPressed(button int) bool {
// }

// // Returns whether the key is pressed.
// //
// // param key The key code as found in {@link Input.Keys}.
// // return true or false.
// func IsKeyPressed(key int) bool {
// }

// // Returns whether the key has just been pressed.
// //
// // param key The key code as found in {@link Input.Keys}.
// // return true or false.
// func IsKeyJustPressed(key int) bool {
// }

// // System dependent method to input a string of text. A dialog box will be created with the given title and the given text as a
// // message for the user. Once the dialog has been closed the provided {@link TextInputListener} will be called on the rendering
// // thread.
// //
// // param listener The TextInputListener.
// // param title The title of the text input dialog.
// // param text The message presented to the user.
// // func GetTextInput(TextInputListener listener, title, text, hint string)

// // Sets the on-screen keyboard visible if available.
// //
// // param visible visible or not
// func SetOnscreenKeyboardVisible(visible bool) {
// }

// // Vibrates for the given amount of time. Note that you'll need the permission
// // <code> <uses-permission android:name="android.permission.VIBRATE" /></code> in your manifest file in order for this to work.
// //
// // param milliseconds the number of milliseconds to vibrate.
// func Vibrate(milliseconds int) {
// }

// // Vibrate with a given pattern. Pass in an array of ints that are the times at which to turn on or off the vibrator. The first
// // one is how long to wait before turning it on, and then after that it alternates. If you want to repeat, pass the index into
// // the pattern at which to start the repeat.
// // param pattern an array of longs of times to turn the vibrator on or off.
// // param repeat the index into pattern at which to repeat, or -1 if you don't want to repeat.
// func VibrateRepeat(pattern []int64, repeat int) {
// }

// // Stops the vibrator
// func CancelVibrate() {
// }

// // The azimuth is the angle of the device's orientation around the z-axis. The positive z-axis points towards the earths
// // center.
// //
// // @see <a
// //      href="http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])">http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])</a>
// // return the azimuth in degrees
// func GetAzimuth() float32 {
// }

// // The pitch is the angle of the device's orientation around the x-axis. The positive x-axis roughly points to the west and is
// // orthogonal to the z- and y-axis.
// // @see <a
// //      href="http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])">http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])</a>
// // return the pitch in degrees
// func GetPitch() float32 {
// }

// // The roll is the angle of the device's orientation around the y-axis. The positive y-axis points to the magnetic north pole
// // of the earth.
// // @see <a
// //      href="http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])">http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])</a>
// // return the roll in degrees
// func GetRoll() float32 {
// }

// // Returns the rotation matrix describing the devices rotation as per <a href=
// // "http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])"
// // >SensorManager#getRotationMatrix(float[], float[], float[], float[])</a>. Does not manipulate the matrix if the platform
// // does not have an accelerometer.
// // param matrix
// func GetRotationMatrix(matrix []float32) {
// }

// // return the time of the event currently reported to the {@link InputProcessor}.
// func GetCurrentEventTime() int64 {
// }

// // Sets whether the BACK button on Android should be caught. This will prevent the app from being paused. Will have no effect
// // on the desktop.
// //
// // param catchBack whether to catch the back button
// func SetCatchBackKey(catchBack bool) {
// }

// // return whether the back button is currently being caught
// func IsCatchBackKey() bool {
// }

// // Sets whether the MENU button on Android should be caught. This will prevent the onscreen keyboard to show up. Will have no
// // effect on the desktop.
// //
// // param catchMenu whether to catch the menu button
// func SetCatchMenuKey(catchMenu bool) {
// }

// // Sets the {@link InputProcessor} that will receive all touch and key input events. It will be called before the
// // {@link ApplicationListener#render()} method each frame.
// //
// // param processor the InputProcessor
// // func SetInputProcessor(InputProcessor processor);

// // return the currently set {@link InputProcessor} or null.
// // public InputProcessor getInputProcessor ();

// // Queries whether a {@link Peripheral} is currently available. In case of Android and the {@link Peripheral#HardwareKeyboard}
// // this returns the whether the keyboard is currently slid out or not.
// //
// // param peripheral the {@link Peripheral}
// // return whether the peripheral is available or not.
// func IsPeripheralAvailable(peripheral Peripheral) bool {
// }

// // return the rotation of the device with respect to its native orientation.
// func GetRotation() int {
// }

// // return the native orientation of the device.
// func GetNativeOrientation() Orientation {
// }

// // Only viable on the desktop. Will confine the mouse cursor location to the window and hide the mouse cursor. X and y
// // coordinates are still reported as if the mouse was not catched.
// // param catched whether to catch or not to catch the mouse cursor
// func SetCursorCatched(catched bool) {
// }

// // return whether the mouse cursor is catched.
// func IsCursorCatched() bool {
// }

// // Only viable on the desktop. Will set the mouse cursor location to the given window coordinates (origin top-left corner).
// // param x the x-position
// // param y the y-position
// func SetCursorPosition(x, y int) {
// }

// The type of Keyboard Keys
type KeyCode uint8

// const (
// 	ANY_KEY             KeyCode = -1
// 	NUM_0                       = 7
// 	NUM_1                       = 8
// 	NUM_2                       = 9
// 	NUM_3                       = 10
// 	NUM_4                       = 11
// 	NUM_5                       = 12
// 	NUM_6                       = 13
// 	NUM_7                       = 14
// 	NUM_8                       = 15
// 	NUM_9                       = 16
// 	A                           = 29
// 	ALT_LEFT                    = 57
// 	ALT_RIGHT                   = 58
// 	APOSTROPHE                  = 75
// 	AT                          = 77
// 	B                           = 30
// 	BACK                        = 4
// 	BACKSLASH                   = 73
// 	C                           = 31
// 	CALL                        = 5
// 	CAMERA                      = 27
// 	CLEAR                       = 28
// 	COMMA                       = 55
// 	D                           = 32
// 	DEL                         = 67
// 	BACKSPACE                   = 67
// 	FORWARD_DEL                 = 112
// 	DPAD_CENTER                 = 23
// 	DPAD_DOWN                   = 20
// 	DPAD_LEFT                   = 21
// 	DPAD_RIGHT                  = 22
// 	DPAD_UP                     = 19
// 	CENTER                      = 23
// 	DOWN                        = 20
// 	LEFT                        = 21
// 	RIGHT                       = 22
// 	UP                          = 19
// 	E                           = 33
// 	ENDCALL                     = 6
// 	ENTER                       = 66
// 	ENVELOPE                    = 65
// 	EQUALS                      = 70
// 	EXPLORER                    = 64
// 	F                           = 34
// 	FOCUS                       = 80
// 	G                           = 35
// 	GRAVE                       = 68
// 	H                           = 36
// 	HEADSETHOOK                 = 79
// 	HOME                        = 3
// 	I                           = 37
// 	J                           = 38
// 	K                           = 39
// 	L                           = 40
// 	LEFT_BRACKET                = 71
// 	M                           = 41
// 	MEDIA_FAST_FORWARD          = 90
// 	MEDIA_NEXT                  = 87
// 	MEDIA_PLAY_PAUSE            = 85
// 	MEDIA_PREVIOUS              = 88
// 	MEDIA_REWIND                = 89
// 	MEDIA_STOP                  = 86
// 	MENU                        = 82
// 	MINUS                       = 69
// 	MUTE                        = 91
// 	N                           = 42
// 	NOTIFICATION                = 83
// 	NUM                         = 78
// 	O                           = 43
// 	P                           = 44
// 	PERIOD                      = 56
// 	PLUS                        = 81
// 	POUND                       = 18
// 	POWER                       = 26
// 	Q                           = 45
// 	R                           = 46
// 	RIGHT_BRACKET               = 72
// 	S                           = 47
// 	SEARCH                      = 84
// 	SEMICOLON                   = 74
// 	SHIFT_LEFT                  = 59
// 	SHIFT_RIGHT                 = 60
// 	SLASH                       = 76
// 	SOFT_LEFT                   = 1
// 	SOFT_RIGHT                  = 2
// 	SPACE                       = 62
// 	STAR                        = 17
// 	SYM                         = 63
// 	T                           = 48
// 	TAB                         = 61
// 	U                           = 49
// 	UNKNOWN                     = 0
// 	V                           = 50
// 	VOLUME_DOWN                 = 25
// 	VOLUME_UP                   = 24
// 	W                           = 51
// 	X                           = 52
// 	Y                           = 53
// 	Z                           = 54
// 	META_ALT_LEFT_ON            = 16
// 	META_ALT_ON                 = 2
// 	META_ALT_RIGHT_ON           = 32
// 	META_SHIFT_LEFT_ON          = 64
// 	META_SHIFT_ON               = 1
// 	META_SHIFT_RIGHT_ON         = 128
// 	META_SYM_ON                 = 4
// 	CONTROL_LEFT                = 129
// 	CONTROL_RIGHT               = 130
// 	ESCAPE                      = 131
// 	END                         = 132
// 	INSERT                      = 133
// 	PAGE_UP                     = 92
// 	PAGE_DOWN                   = 93
// 	PICTSYMBOLS                 = 94
// 	SWITCH_CHARSET              = 95
// 	BUTTON_CIRCLE               = 255
// 	BUTTON_A                    = 96
// 	BUTTON_B                    = 97
// 	BUTTON_C                    = 98
// 	BUTTON_X                    = 99
// 	BUTTON_Y                    = 100
// 	BUTTON_Z                    = 101
// 	BUTTON_L1                   = 102
// 	BUTTON_R1                   = 103
// 	BUTTON_L2                   = 104
// 	BUTTON_R2                   = 105
// 	BUTTON_THUMBL               = 106
// 	BUTTON_THUMBR               = 107
// 	BUTTON_START                = 108
// 	BUTTON_SELECT               = 109
// 	BUTTON_MODE                 = 110

// 	NUMPAD_0 = 144
// 	NUMPAD_1 = 145
// 	NUMPAD_2 = 146
// 	NUMPAD_3 = 147
// 	NUMPAD_4 = 148
// 	NUMPAD_5 = 149
// 	NUMPAD_6 = 150
// 	NUMPAD_7 = 151
// 	NUMPAD_8 = 152
// 	NUMPAD_9 = 153

// 	COLON = 243
// 	F1    = 244
// 	F2    = 245
// 	F3    = 246
// 	F4    = 247
// 	F5    = 248
// 	F6    = 249
// 	F7    = 250
// 	F8    = 251
// 	F9    = 252
// 	F10   = 253
// 	F11   = 254
// 	F12   = 255
// )

// // return a human readable representation of the keycode. The returned value can be used in
// // {@link Input.Keys#valueOf(String)}
// func ToString(keycode uint8) string {
// 	switch keycode {
// 	// META* variables should not be used with this method.
// 	case UNKNOWN:
// 		return "Unknown"
// 	case SOFT_LEFT:
// 		return "Soft Left"
// 	case SOFT_RIGHT:
// 		return "Soft Right"
// 	case HOME:
// 		return "Home"
// 	case BACK:
// 		return "Back"
// 	case CALL:
// 		return "Call"
// 	case ENDCALL:
// 		return "End Call"
// 	case NUM_0:
// 		return "0"
// 	case NUM_1:
// 		return "1"
// 	case NUM_2:
// 		return "2"
// 	case NUM_3:
// 		return "3"
// 	case NUM_4:
// 		return "4"
// 	case NUM_5:
// 		return "5"
// 	case NUM_6:
// 		return "6"
// 	case NUM_7:
// 		return "7"
// 	case NUM_8:
// 		return "8"
// 	case NUM_9:
// 		return "9"
// 	case STAR:
// 		return "*"
// 	case POUND:
// 		return "#"
// 	case UP:
// 		return "Up"
// 	case DOWN:
// 		return "Down"
// 	case LEFT:
// 		return "Left"
// 	case RIGHT:
// 		return "Right"
// 	case CENTER:
// 		return "Center"
// 	case VOLUME_UP:
// 		return "Volume Up"
// 	case VOLUME_DOWN:
// 		return "Volume Down"
// 	case POWER:
// 		return "Power"
// 	case CAMERA:
// 		return "Camera"
// 	case CLEAR:
// 		return "Clear"
// 	case A:
// 		return "A"
// 	case B:
// 		return "B"
// 	case C:
// 		return "C"
// 	case D:
// 		return "D"
// 	case E:
// 		return "E"
// 	case F:
// 		return "F"
// 	case G:
// 		return "G"
// 	case H:
// 		return "H"
// 	case I:
// 		return "I"
// 	case J:
// 		return "J"
// 	case K:
// 		return "K"
// 	case L:
// 		return "L"
// 	case M:
// 		return "M"
// 	case N:
// 		return "N"
// 	case O:
// 		return "O"
// 	case P:
// 		return "P"
// 	case Q:
// 		return "Q"
// 	case R:
// 		return "R"
// 	case S:
// 		return "S"
// 	case T:
// 		return "T"
// 	case U:
// 		return "U"
// 	case V:
// 		return "V"
// 	case W:
// 		return "W"
// 	case X:
// 		return "X"
// 	case Y:
// 		return "Y"
// 	case Z:
// 		return "Z"
// 	case COMMA:
// 		return ","
// 	case PERIOD:
// 		return "."
// 	case ALT_LEFT:
// 		return "L-Alt"
// 	case ALT_RIGHT:
// 		return "R-Alt"
// 	case SHIFT_LEFT:
// 		return "L-Shift"
// 	case SHIFT_RIGHT:
// 		return "R-Shift"
// 	case TAB:
// 		return "Tab"
// 	case SPACE:
// 		return "Space"
// 	case SYM:
// 		return "SYM"
// 	case EXPLORER:
// 		return "Explorer"
// 	case ENVELOPE:
// 		return "Envelope"
// 	case ENTER:
// 		return "Enter"
// 	case DEL:
// 		return "Delete" // also BACKSPACE
// 	case GRAVE:
// 		return "`"
// 	case MINUS:
// 		return "-"
// 	case EQUALS:
// 		return "="
// 	case LEFT_BRACKET:
// 		return "["
// 	case RIGHT_BRACKET:
// 		return "]"
// 	case BACKSLASH:
// 		return "\\"
// 	case SEMICOLON:
// 		return ";"
// 	case APOSTROPHE:
// 		return "'"
// 	case SLASH:
// 		return "/"
// 	case AT:
// 		return "@"
// 	case NUM:
// 		return "Num"
// 	case HEADSETHOOK:
// 		return "Headset Hook"
// 	case FOCUS:
// 		return "Focus"
// 	case PLUS:
// 		return "Plus"
// 	case MENU:
// 		return "Menu"
// 	case NOTIFICATION:
// 		return "Notification"
// 	case SEARCH:
// 		return "Search"
// 	case MEDIA_PLAY_PAUSE:
// 		return "Play/Pause"
// 	case MEDIA_STOP:
// 		return "Stop Media"
// 	case MEDIA_NEXT:
// 		return "Next Media"
// 	case MEDIA_PREVIOUS:
// 		return "Prev Media"
// 	case MEDIA_REWIND:
// 		return "Rewind"
// 	case MEDIA_FAST_FORWARD:
// 		return "Fast Forward"
// 	case MUTE:
// 		return "Mute"
// 	case PAGE_UP:
// 		return "Page Up"
// 	case PAGE_DOWN:
// 		return "Page Down"
// 	case PICTSYMBOLS:
// 		return "PICTSYMBOLS"
// 	case SWITCH_CHARSET:
// 		return "SWITCH_CHARSET"
// 	case BUTTON_A:
// 		return "A Button"
// 	case BUTTON_B:
// 		return "B Button"
// 	case BUTTON_C:
// 		return "C Button"
// 	case BUTTON_X:
// 		return "X Button"
// 	case BUTTON_Y:
// 		return "Y Button"
// 	case BUTTON_Z:
// 		return "Z Button"
// 	case BUTTON_L1:
// 		return "L1 Button"
// 	case BUTTON_R1:
// 		return "R1 Button"
// 	case BUTTON_L2:
// 		return "L2 Button"
// 	case BUTTON_R2:
// 		return "R2 Button"
// 	case BUTTON_THUMBL:
// 		return "Left Thumb"
// 	case BUTTON_THUMBR:
// 		return "Right Thumb"
// 	case BUTTON_START:
// 		return "Start"
// 	case BUTTON_SELECT:
// 		return "Select"
// 	case BUTTON_MODE:
// 		return "Button Mode"
// 	case FORWARD_DEL:
// 		return "Forward Delete"
// 	case CONTROL_LEFT:
// 		return "L-Ctrl"
// 	case CONTROL_RIGHT:
// 		return "R-Ctrl"
// 	case ESCAPE:
// 		return "Escape"
// 	case END:
// 		return "End"
// 	case INSERT:
// 		return "Insert"
// 	case NUMPAD_0:
// 		return "Numpad 0"
// 	case NUMPAD_1:
// 		return "Numpad 1"
// 	case NUMPAD_2:
// 		return "Numpad 2"
// 	case NUMPAD_3:
// 		return "Numpad 3"
// 	case NUMPAD_4:
// 		return "Numpad 4"
// 	case NUMPAD_5:
// 		return "Numpad 5"
// 	case NUMPAD_6:
// 		return "Numpad 6"
// 	case NUMPAD_7:
// 		return "Numpad 7"
// 	case NUMPAD_8:
// 		return "Numpad 8"
// 	case NUMPAD_9:
// 		return "Numpad 9"
// 	case COLON:
// 		return ":"
// 	case F1:
// 		return "F1"
// 	case F2:
// 		return "F2"
// 	case F3:
// 		return "F3"
// 	case F4:
// 		return "F4"
// 	case F5:
// 		return "F5"
// 	case F6:
// 		return "F6"
// 	case F7:
// 		return "F7"
// 	case F8:
// 		return "F8"
// 	case F9:
// 		return "F9"
// 	case F10:
// 		return "F10"
// 	case F11:
// 		return "F11"
// 	case F12:
// 		return "F12"
// 		// BUTTON_CIRCLE unhandled, as it conflicts with the more likely to be pressed F12
// 	default:
// 		// key name not found
// 		return ""
// 	}
// }
