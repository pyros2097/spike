// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"encoding/binary"
	"math"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"

	"github.com/pyros2097/spike/math/vector"
)

const (
	coordsPerVertex = 3
	vertexCount     = 3

	vertexShader = `#version 100
	uniform vec2 offset;

	attribute vec4 position;
	void main() {
	 // offset comes in with x/y values between 0 and 1.
	 // position bounds are -1 to 1.
	 vec4 offset4 = vec4(1.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	 gl_Position = position + offset4;
	}`

	fragmentShader = `#version 100
	precision mediump float;
	uniform vec4 color;
	void main() {
	 gl_FragColor = color;
	}`
)

var (
	targetWidth  float32
	targetHeight float32
	pauseState   bool

	Camera2d *Camera
	Camera3d *Camera

	running   bool
	fpsTicker *time.Ticker
	lastTime  time.Duration
	delta     time.Duration

	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer
	err      error

	green  float32
	touchX float32
	touchY float32

	triangleData = f32.Bytes(binary.LittleEndian,
		0.0, 0.4, 0.0, // top left
		0.0, 0.0, 0.0, // bottom left
		0.4, 0.0, 0.0, // bottom right
	)
)

type SBatch struct {
}

func (b SBatch) Begin() {}
func (b SBatch) End()   {}

var tempBatch *SBatch

// 	/** Creates a new GestureDetector with default values: halfTapSquareSize=20, tapCountInterval=0.4f, longPressDuration=1.1f,
// 	 * maxFlingDelay=0.15f. */
//final VelocityTracker tracker = new VelocityTracker();
//final Task longPressTask = new Task() {
// 		@Override
// 		public void run () {
// 			if (!longPressFired) longPressFired = listener.longPress(pointer1.x, pointer1.y);
// 		}
// 	};

/** {@link InputProcessor} implementation that detects gestures (tap, long press, fling, pan, zoom, pinch) and hands them to a
// public class GestureDetector extends InputAdapter {
// 	/** @param halfTapSquareSize half width in pixels of the square around an initial touch event, see
// 	 *           {@link GestureListener#tap(float, float, int, int)}.
// 	 * @param tapCountInterval time in seconds that must pass for two touch down/up sequences to be detected as consecutive taps.
// 	 * @param longPressDuration time in seconds that must pass for the detector to fire a
// 	 *           {@link GestureListener#longPress(float, float)} event.
// 	 * @param maxFlingDelay time in seconds the finger must have been dragged for a fling event to be fired, see
// 	 *           {@link GestureListener#fling(float, float, int)}
// 	 * @param listener May be null if the listener will be set later. */
type GestureEvent uint8

const (
	None GestureEvent = iota
	Tap
	Fling
	Pan
	LongPress
)

var (
	tapSquareSize                      float32 = 20
	tapCountInterval                   int64   = 0 //0.4f
	longPressSeconds                   float32 = 1.1
	maxFlingDelay                      float64 = 0.15
	inTapSquare                                = false
	tapCount                           int     = 0
	lastTapTime                        int64   = 0
	lastTapX, lastTapY                 float32
	lastTapButton, lastTapPointer      int
	longPressFired, pinching, panning  bool
	tapSquareCenterX, tapSquareCenterY float32
	gestureStartTime                   int64
	lastTouchEvent                     touch.Type

	pointer1         = vector.NewVector2Empty()
	pointer2         = vector.NewVector2Empty()
	initialPointer1  = vector.NewVector2Empty()
	initialPointer2  = vector.NewVector2Empty()
	lastGestureEvent GestureEvent
	gestureQueue     chan GestureEvent
)

func isWithinTapSquare(x, y, centerX, centerY float32) bool {
	return float32(math.Abs(float64(x-centerX))) < tapSquareSize && float32(math.Abs(float64(y-centerY))) < tapSquareSize
}

func touchDown(x, y float32, pointer, button int) bool {
	if pointer > 1 {
		return false
	}
	if pointer == 0 {
		pointer1.Set(x, y)
		gestureStartTime = time.Now().UnixNano() //Gdx.input.getCurrentEventTime()
		// 	tracker.start(x, y, gestureStartTime);
		inTapSquare = true
		pinching = false
		longPressFired = false
		tapSquareCenterX = x
		tapSquareCenterY = y
		// 		if (!longPressTask.isScheduled()) Timer.schedule(longPressTask, longPressSeconds);
		// 	}
	} else {
		// 	// Start pinch.
		// 	pointer2.set(x, y);
		// 	inTapSquare = false;
		// 	pinching = true;
		// 	initialPointer1.set(pointer1);
		// 	initialPointer2.set(pointer2);
		// 	longPressTask.cancel();
	}
	return false
}

func touchDragged(x, y float32, pointer int) bool {
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
	// if (pinching) {
	// 	if (listener != null) {
	// 		boolean result = listener.pinch(initialPointer1, initialPointer2, pointer1, pointer2);
	// 		return listener.zoom(initialPointer1.dst(initialPointer2), pointer1.dst(pointer2)) || result;
	// 	}
	// 	return false;
	// }

	// // update tracker
	// tracker.update(x, y, Gdx.input.getCurrentEventTime());

	// // check if we are still tapping.
	if inTapSquare && !isWithinTapSquare(x, y, tapSquareCenterX, tapSquareCenterY) {
		// 	longPressTask.cancel();
		inTapSquare = false
	}

	// if we have left the tap square, we are panning
	if !inTapSquare {
		panning = true
		// 	return listener.pan(x, y, tracker.deltaX, tracker.deltaY);
	}

	return false
}

func touchUp(x, y float32, pointer, button int) bool {
	if pointer > 1 {
		return false
	}
	// check if we are still tapping.
	if inTapSquare && !isWithinTapSquare(x, y, tapSquareCenterX, tapSquareCenterY) {
		inTapSquare = false
	}

	wasPanning := panning
	panning = false

	// 		longPressTask.cancel();
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
		gestureQueue <- Tap
		return true
		// return listener.tap(x, y, tapCount, button);
	}

	// 		if (pinching) {
	// 			// handle pinch end
	// 			pinching = false;
	// 			panning = true;
	// 			// we are in pan mode again, reset velocity tracker
	// 			if (pointer == 0) {
	// 				// first pointer has lifted off, set up panning to use the second pointer...
	// 				tracker.start(pointer2.x, pointer2.y, Gdx.input.getCurrentEventTime());
	// 			} else {
	// 				// second pointer has lifted off, set up panning to use the first pointer...
	// 				tracker.start(pointer1.x, pointer1.y, Gdx.input.getCurrentEventTime());
	// 			}
	// 			return false;
	// 		}

	// handle no longer panning
	handled := false
	if wasPanning && !panning {
		handled = false //listener.panStop(x, y, pointer, button);
	}

	// handle fling
	gestureStartTime = 0
	// 		long time = Gdx.input.getCurrentEventTime();
	// 		if (time - tracker.lastTime < maxFlingDelay) {
	// 			tracker.update(x, y, time);
	// 			handled = listener.fling(tracker.getVelocityX(), tracker.getVelocityY(), button) || handled;
	// 		}
	return handled
}

/*Important:
 *  The Target Width  and Target Height refer to the nominal width and height of the game for the
 *  graphics which are created  for this width and height, this allows for the Stage to scale this
 *  graphics for all screen width and height. Therefore your game will work on all screen sizes
 *  but maybe blurred or look awkward on some devices.
 *  ex:
 *  My Game targetWidth = 800 targetHeight = 480
 *  Then my game works perfectly for SCREEN_WIDTH = 800 SCREEN_HEIGHT = 480
 *  and on others screen sizes it is just zoomed/scaled but works fine thats all
 */
func Init(title string, width, height float32) {
	println("Initializing Gdx")
	targetWidth = width
	targetHeight = height
	// camera2d = &graphics.Camera{}
	allScenes = make(map[string]*Scene)
	fpsTicker = time.NewTicker(1000 / 30 * time.Millisecond)
	running = true
	initConfig(title)
}

// It sets up a window and rendering surface and manages the
// different aspects of your application, namely {@link Graphics}, {@link Audio}, {@link Input} and {@link Files}.
// This is the main entry point of your project.
// Note that all Music instances will be automatically paused when the current scene's OnPause() method is
// called, and automatically resumed when the OnResume() method is called.
func Run() {
	lastTime := time.Now()
	app.Main(func(a app.App) {
		for now := range fpsTicker.C {
			if !running {
				break
			}
			delta = now.Sub(lastTime)
			lastTime = now
			var sz size.Event
			// e := <-a.Events()
			for e := range a.Events() {
				switch e := app.Filter(e).(type) {
				case lifecycle.Event:
					switch e.Crosses(lifecycle.StageVisible) {
					case lifecycle.CrossOn:
						Start()
					case lifecycle.CrossOff:
						Exit()
					}
				case size.Event: // resize event
					sz = e
					touchX = float32(sz.WidthPx / 2)
					touchY = float32(sz.HeightPx / 2)
				case paint.Event:
					// print("Painting")
					// Do update here
					onPaint(sz, float32(delta))
					a.EndPaint(e)
				case touch.Event:
					// print("Touching")
					// send input events here or before paint just store the last state
					touchX = e.X
					touchY = e.Y
					lastTouchEvent = e.Type
					switch e.Type {
					case touch.TypeBegin: // touch down
						println("Begin")
						touchDown(touchX, touchY, 0, 0)
					case touch.TypeEnd: // touch up
						println("End")
						touchUp(touchX, touchY, 0, 0)
					case touch.TypeMove: // touch drag
						println("Moving")
						touchDragged(touchX, touchY, 0)
						// println(e.Sequence)
						// log.Printf("%d", e.Sequence)
					}
				}
			}
		}
	})
}

func GetTouchX() float32 {
	return touchX
}

func GetTouchY() float32 {
	return touchY
}

func Start() {
	println("Starting")
	program, err = glutil.CreateProgram(vertexShader, fragmentShader)
	if err != nil {
		panic("error creating GL program: " + err.Error())
		return
	}

	buf = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = gl.GetAttribLocation(program, "position")
	color = gl.GetUniformLocation(program, "color")
	offset = gl.GetUniformLocation(program, "offset")

	lastTouchEvent = 111
	gestureQueue = make(chan GestureEvent, 100)
	// TODO(crawshaw): the debug package needs to put GL state init here
	// Can this be an app.RegisterFilter call now??
	SetScene(currentScene.Name)
	go func() {
		for {
			lastGestureEvent = <-gestureQueue
		}
	}()
}

// This is the main rendering call that updates the time, updates the stage,
// loads assets asynchronously, updates the camera and FPS text.
func onPaint(sz size.Event, delta float32) {
	// asset.Load() Async
	gl.ClearColor(currentScene.BGColor.R, currentScene.BGColor.G, currentScene.BGColor.B, currentScene.BGColor.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for _, child := range currentScene.Children {
		if child.OnAct != nil {
			child.OnAct(child, delta)
		}
		if child.OnDraw != nil {
			child.OnDraw(child, tempBatch, 1.0)
		}
		if lastTouchEvent == touch.TypeBegin && child.OnTouchDown != nil {
			child.OnTouchDown(child, touchX, touchY, 1, 1)
		}
		if lastTouchEvent == touch.TypeEnd && child.OnTouchUp != nil {
			child.OnTouchUp(child, touchX, touchY, 1, 1)
		}
		if lastTouchEvent == touch.TypeMove && child.OnTouchDragged != nil {
			child.OnTouchDragged(child, touchX, touchY, 1)
		}
		if lastGestureEvent == Tap && child.OnTap != nil {
			child.OnTap(child, touchX, touchY, tapCount, 0)
		}
	}
	lastGestureEvent = None

	gl.UseProgram(program)

	green += 0.01
	if green > 1 {
		green = 0
	}
	gl.Uniform4f(color, 0, green, 0, 1)

	gl.Uniform2f(offset, touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))

	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.EnableVertexAttribArray(position)
	gl.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	gl.DisableVertexAttribArray(position)

	debug.DrawFPS(sz)
}

func Pause() {
	println("Pausing")
	pauseState = true
	musicPlayer.Pause()
	soundsPlayer.Pause()
	// unloadAll()
	if currentScene.OnPause != nil {
		currentScene.OnPause(currentScene)

	}
}

func Resume() {
	println("Resuming")
	pauseState = false
	musicPlayer.Pause()
	soundsPlayer.Play()
	if currentScene.OnResume != nil {
		currentScene.OnResume(currentScene)
	}
}

// Use this to exit your game safely
// It will automatically unload all your assets and dispose the stage
// Schedule an exit from the application. On android, this will cause a call to pause() and dispose() some time in the future,
// it will not immediately finish your application.
// On iOS this should be avoided in production as it breaks Apples guidelines
func Exit() {
	println("Exiting")
	running = false
	if currentScene.OnPause != nil {
		currentScene.OnPause(currentScene)
		soundsPlayer.Close()
	}
	gl.DeleteProgram(program)
	gl.DeleteBuffer(buf)
}

// 	/** No further gesture events will be triggered for the current touch, if any. */
// 	public void cancel () {
// 		longPressTask.cancel();
// 		longPressFired = true;
// 	}

// 	/** @return whether the user touched the screen long enough to trigger a long press event. */
// 	public boolean isLongPressed () {
// 		return isLongPressed(longPressSeconds);
// 	}

// 	* @param duration
// 	 * @return whether the user touched the screen for as much or more than the given duration.
// 	public boolean isLongPressed (float duration) {
// 		if (gestureStartTime == 0) return false;
// 		return TimeUtils.nanoTime() - gestureStartTime > (long)(duration * 1000000000l);
// 	}

// 	public boolean isPanning () {
// 		return panning;
// 	}

// 	public void reset () {
// 		gestureStartTime = 0;
// 		panning = false;
// 		inTapSquare = false;
// 	}

// 	/** The tap square will not longer be used for the current touch. */
// 	public void invalidateTapSquare () {
// 		inTapSquare = false;
// 	}

// 	public void setTapSquareSize (float halfTapSquareSize) {
// 		this.tapSquareSize = halfTapSquareSize;
// 	}

// 	/** @param tapCountInterval time in seconds that must pass for two touch down/up sequences to be detected as consecutive taps. */
// 	public void setTapCountInterval (float tapCountInterval) {
// 		this.tapCountInterval = (long)(tapCountInterval * 1000000000l);
// 	}

// 	public void setLongPressSeconds (float longPressSeconds) {
// 		this.longPressSeconds = longPressSeconds;
// 	}

// 	public void setMaxFlingDelay (long maxFlingDelay) {
// 		this.maxFlingDelay = maxFlingDelay;
// 	}

// 	static class VelocityTracker {
// 		int sampleSize = 10;
// 		float lastX, lastY;
// 		float deltaX, deltaY;
// 		long lastTime;
// 		int numSamples;
// 		float[] meanX = new float[sampleSize];
// 		float[] meanY = new float[sampleSize];
// 		long[] meanTime = new long[sampleSize];

// 		public void start (float x, float y, long timeStamp) {
// 			lastX = x;
// 			lastY = y;
// 			deltaX = 0;
// 			deltaY = 0;
// 			numSamples = 0;
// 			for (int i = 0; i < sampleSize; i++) {
// 				meanX[i] = 0;
// 				meanY[i] = 0;
// 				meanTime[i] = 0;
// 			}
// 			lastTime = timeStamp;
// 		}

// 		public void update (float x, float y, long timeStamp) {
// 			long currTime = timeStamp;
// 			deltaX = x - lastX;
// 			deltaY = y - lastY;
// 			lastX = x;
// 			lastY = y;
// 			long deltaTime = currTime - lastTime;
// 			lastTime = currTime;
// 			int index = numSamples % sampleSize;
// 			meanX[index] = deltaX;
// 			meanY[index] = deltaY;
// 			meanTime[index] = deltaTime;
// 			numSamples++;
// 		}

// 		public float getVelocityX () {
// 			float meanX = getAverage(this.meanX, numSamples);
// 			float meanTime = getAverage(this.meanTime, numSamples) / 1000000000.0f;
// 			if (meanTime == 0) return 0;
// 			return meanX / meanTime;
// 		}

// 		public float getVelocityY () {
// 			float meanY = getAverage(this.meanY, numSamples);
// 			float meanTime = getAverage(this.meanTime, numSamples) / 1000000000.0f;
// 			if (meanTime == 0) return 0;
// 			return meanY / meanTime;
// 		}

// 		private float getAverage (float[] values, int numSamples) {
// 			numSamples = Math.min(sampleSize, numSamples);
// 			float sum = 0;
// 			for (int i = 0; i < numSamples; i++) {
// 				sum += values[i];
// 			}
// 			return sum / numSamples;
// 		}

// 		private long getAverage (long[] values, int numSamples) {
// 			numSamples = Math.min(sampleSize, numSamples);
// 			long sum = 0;
// 			for (int i = 0; i < numSamples; i++) {
// 				sum += values[i];
// 			}
// 			if (numSamples == 0) return 0;
// 			return sum / numSamples;
// 		}

// 		private float getSum (float[] values, int numSamples) {
// 			numSamples = Math.min(sampleSize, numSamples);
// 			float sum = 0;
// 			for (int i = 0; i < numSamples; i++) {
// 				sum += values[i];
// 			}
// 			if (numSamples == 0) return 0;
// 			return sum;
// 		}
// 	}
// }
