// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// The type of Gestures that can occur
package gesture

import (
	"math"
	"time"

	"github.com/pyros2097/spike/math/vector"
)

type GestureType uint8

const (
	None GestureType = iota
	Tap
	Fling
	Pan
	LongPress
	SwipeLeft
	SwipeRight
	SwipeUp
	SwipeDown
)

type GestureEvent struct {
	Type            GestureType
	X, Y            float32
	Pointer, Button int
}

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
	tracker                            = NewVelocityTracker()
	longPressScheduled                 = false
	longPressTask                      = time.AfterFunc(time.Duration(float32(time.Second)*LongPressSeconds), fireLongPress)
	fireLongPress                      = func() {
		if !longPressFired {
			longPressFired = true //listener.longPress(pointer1.x, pointer1.y);
		}
	}
	pointer1        = vector.NewVector2Empty()
	pointer2        = vector.NewVector2Empty()
	initialPointer1 = vector.NewVector2Empty()
	initialPointer2 = vector.NewVector2Empty()
	event           GestureEvent
	queue           []GestureEvent
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
func GetEvent() *GestureEvent {
	if len(queue) > 0 {
		event = queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		return &event
	}
	return nil
}

func TouchDown(x, y float32, pointer, button int) bool {
	if pointer > 1 {
		return false
	}
	if pointer == 0 {
		pointer1.Set(x, y)
		gestureStartTime = time.Now().UnixNano() //Gdx.input.getCurrentEventTime()
		tracker.Start(x, y, gestureStartTime)
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
	return false
}

func TouchUp(x, y float32, pointer, button int) bool {
	if pointer > 1 {
		return false
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
		queue = append(queue, GestureEvent{
			Type:    Tap,
			X:       x,
			Y:       y,
			Pointer: 0,
			Button:  0,
		})
		return true
	}

	if pinching {
		// handle pinch end
		pinching = false
		panning = true
		// we are in pan mode again, reset velocity tracker
		if pointer == 0 {
			// first pointer has lifted off, set up panning to use the second pointer...
			tracker.Start(pointer2.X, pointer2.Y, time.Now().UnixNano()) //Gdx.input.getCurrentEventTime())
		} else {
			// second pointer has lifted off, set up panning to use the first pointer...
			tracker.Start(pointer1.X, pointer1.Y, time.Now().UnixNano()) //Gdx.input.getCurrentEventTime())
		}
		return false
	}

	// handle no longer panning
	handled := false
	if wasPanning && !panning {
		//listener.panStop(x, y, pointer, button);
		handled = false
	}

	// handle fling
	gestureStartTime = 0
	time := time.Now().UnixNano() //Gdx.input.getCurrentEventTime();
	if time-tracker.LastTime < MaxFlingDelay {
		tracker.Update(x, y, time)
		// lastGestureEvent = Fling
		// handled = listener.fling(tracker.getVelocityX(), tracker.getVelocityY(), button) || handled
	}
	return handled
}

func TouchDragged(x, y float32, pointer int) bool {
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
		// 	if (listener != null) {
		// 		boolean result = listener.pinch(initialPointer1, initialPointer2, pointer1, pointer2);
		// 		return listener.zoom(initialPointer1.dst(initialPointer2), pointer1.dst(pointer2)) || result;
		// 	}
		// 	return false;
	}

	// update tracker
	tracker.Update(x, y, time.Now().UnixNano()) //Gdx.input.getCurrentEventTime())

	// // check if we are still tapping.
	if inTapSquare && !isWithinTapSquare(x, y, tapSquareCenterX, tapSquareCenterY) {
		longPressTask.Stop()
		longPressScheduled = false
		inTapSquare = false
	}

	// if we have left the tap square, we are panning
	if !inTapSquare {
		panning = true
		// 	return listener.pan(x, y, tracker.deltaX, tracker.deltaY);
	}

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

type VelocityTracker struct {
	sampleSize     int
	lastX, lastY   float32
	deltaX, deltaY float32
	LastTime       int64
	numSamples     int
	meanX          []float32
	meanY          []float32
	meanTime       []int64
}

func NewVelocityTracker() *VelocityTracker {
	return &VelocityTracker{
		sampleSize: 10,
		meanX:      make([]float32, 10),
		meanY:      make([]float32, 10),
		meanTime:   make([]int64, 10),
	}
}

func (self *VelocityTracker) Start(x, y float32, timeStamp int64) {
	self.lastX = x
	self.lastY = y
	self.deltaX = 0
	self.deltaY = 0
	self.numSamples = 0
	for i := 0; i < self.sampleSize; i++ {
		self.meanX[i] = 0
		self.meanY[i] = 0
		self.meanTime[i] = 0
	}
	self.LastTime = timeStamp
}

func (self *VelocityTracker) Update(x, y float32, timeStamp int64) {
	currTime := timeStamp
	self.deltaX = x - self.lastX
	self.deltaY = y - self.lastY
	self.lastX = x
	self.lastY = y
	deltaTime := currTime - self.LastTime
	self.LastTime = currTime
	index := self.numSamples % self.sampleSize
	self.meanX[index] = self.deltaX
	self.meanY[index] = self.deltaY
	self.meanTime[index] = deltaTime
	self.numSamples++
}

func (self *VelocityTracker) GetVelocityX() float32 {
	meanX := self.GetAverage(self.meanX, self.numSamples)
	meanTime := self.GetAverageInt(self.meanTime, self.numSamples) / 1000000000.0
	if meanTime == 0 {
		return 0
	}
	return meanX / float32(meanTime)
}

func (self *VelocityTracker) GetVelocityY() float32 {
	meanY := self.GetAverage(self.meanY, self.numSamples)
	meanTime := self.GetAverageInt(self.meanTime, self.numSamples) / 1000000000.0
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

func (self *VelocityTracker) GetAverage(values []float32, numSamples int) float32 {
	numSamples = min(self.sampleSize, numSamples)
	var sum float32
	for i := 0; i < numSamples; i++ {
		sum += values[i]
	}
	return sum / float32(numSamples)
}

func (self *VelocityTracker) GetAverageInt(values []int64, numSamples int) int64 {
	numSamples = min(self.sampleSize, numSamples)
	var sum int64
	for i := 0; i < numSamples; i++ {
		sum += values[i]
	}
	if numSamples == 0 {
		return 0
	}
	return sum / int64(numSamples)
}

func (self *VelocityTracker) GetSum(values []float32, numSamples int) float32 {
	numSamples = min(self.sampleSize, numSamples)
	var sum float32
	for i := 0; i < numSamples; i++ {
		sum += values[i]
	}
	if numSamples == 0 {
		return 0
	}
	return sum
}
