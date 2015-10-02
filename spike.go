// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"encoding/binary"
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
						onStart()
					case lifecycle.CrossOff:
						onStop()
						Exit()
					}
				case size.Event:
					sz = e
					touchX = float32(sz.WidthPx / 2)
					touchY = float32(sz.HeightPx / 2)
				case paint.Event:
					onPaint(sz)
					a.EndPaint(e)
				case touch.Event:
					switch e.Type {
					case touch.TypeBegin: // touch down
						println("Begin")
					case touch.TypeEnd: // touch up
						println("End")
					case touch.TypeMove:
						println("Moving")
						println(e.Sequence)
						// log.Printf("%d", e.Sequence)
					}
					touchX = e.X
					touchY = e.Y
					// if touchX > 50 {
					// SetScene("options")
					// }
				}
			}
		}
	})
}

func onStart() {
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

	// Config.setup();
	// Serializer.setup();
	// stage2d = new com.badlogic.gdx.scenes.scene2d.Stage(Gdx.graphics.getWidth(), Gdx.graphics.getHeight(),
	// 		Scene.configJson.getBoolean("keepAspectRatio"));

	// setTouchable(Touchable.childrenOnly);
	// Camera.reset();
	// stage2d.clear();
	// stage3d.clear();
	// stage2d.addListener(touchInput);
	// setSize(targetWidth, targetHeight);
	// setBounds(0,0, targetWidth, targetHeight);
	// setColor(1f, 1f, 1f, 1f);
	// setVisible(true);
	// stage2d.getRoot().setPosition(0, 0);
	// stage2d.getRoot().setVisible(true);
	// stage3d.getRoot().setPosition(0, 0, 0);
	// stage3d.getRoot().setVisible(true);
	// sceneName = this.getClass().getName();
	// setName(sceneName);
	// Scene.log("Current Scene: "+sceneName);
	// currentScene = this;
	// load(sceneName);
	// cullingEnabled = true;

	// stage2d.getRoot().setName("Root");
	// stage2d.getRoot().setTouchable(Touchable.childrenOnly);
	// stage2d.setCamera(new Camera());
	// inputMux = new InputMultiplexer();
	// inputMux.addProcessor(stage2d);
	// stage3d = new Stage3d();
	// //camController = new CameraInputController(stage3d.getCamera());
	// //inputMux.addProcessor(stage3d);
	// //inputMux.addProcessor(camController);
	// shapeRenderer = new ShapeRenderer();
	// selectionBox.setTouchable(Touchable.disabled);
	// selectionBox.setName("Shape");
	// JsonValue sv = Scene.jsonReader.parse(Gdx.files.internal(Asset.basePath+"scene"));
	// for(JsonValue jv: sv.iterator())
	// 	scenesMap.put(jv.name, jv.asString());
	// setScene(scenesMap.firstKey());
	// Gdx.input.setCatchBackKey(true);
	// Gdx.input.setCatchMenuKey(true);
	// Gdx.input.setInputProcessor(inputMux);
	// xlines = (int)Gdx.graphics.getWidth()/dots;
	// 	ylines = (int)Gdx.graphics.getHeight()/dots;

	// TODO(crawshaw): the debug package needs to put GL state init here
	// Can this be an app.RegisterFilter call now??
	SetScene(currentScene.Name)
}

func onStop() {
	gl.DeleteProgram(program)
	gl.DeleteBuffer(buf)
}

// This is the main rendering call that updates the time, updates the stage,
// loads assets asynchronously, updates the camera and FPS text.
func onPaint(sz size.Event) {
	// asset.Load() Async
	gl.ClearColor(currentScene.BGColor.R, currentScene.BGColor.G, currentScene.BGColor.B, currentScene.BGColor.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// for _, child := range children {
	//   if child.OnTouchUp != nil {}
	//   if child.OnTouchDown != nil {}
	// }
	//stage3d.act();
	// stage3d.draw();
	//camController.update();
	// Scene.stage2d.act();//Gdx.graphics.getDeltaTime();
	// Scene.stage2d.draw();

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

func Load(sceneName string) {

}

func Save() {

}

func Pause() {
	println("Pausing")
	pauseState = true
	musicPlayer.Pause()
	soundsPlayer.Pause()
	// unloadAll()
	if currentScene.OnPause != nil {
		currentScene.OnPause()

	}
}

func Resume() {
	println("Resuming")
	pauseState = false
	musicPlayer.Pause()
	soundsPlayer.Play()
	if currentScene.OnResume != nil {
		currentScene.OnResume()
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
		currentScene.OnPause()
		soundsPlayer.Close()
	}
}

/** {@link InputProcessor} implementation that detects gestures (tap, long press, fling, pan, zoom, pinch) and hands them to a
 * {@link GestureListener}.
 * @author mzechner */
// public class GestureDetector extends InputAdapter {
// 	final GestureListener listener;
// 	private float tapSquareSize;
// 	private long tapCountInterval;
// 	private float longPressSeconds;
// 	private long maxFlingDelay;

// 	private boolean inTapSquare;
// 	private int tapCount;
// 	private long lastTapTime;
// 	private float lastTapX, lastTapY;
// 	private int lastTapButton, lastTapPointer;
// 	boolean longPressFired;
// 	private boolean pinching;
// 	private boolean panning;

// 	private final VelocityTracker tracker = new VelocityTracker();
// 	private float tapSquareCenterX, tapSquareCenterY;
// 	private long gestureStartTime;
// 	Vector2 pointer1 = new Vector2();
// 	private final Vector2 pointer2 = new Vector2();
// 	private final Vector2 initialPointer1 = new Vector2();
// 	private final Vector2 initialPointer2 = new Vector2();

// 	private final Task longPressTask = new Task() {
// 		@Override
// 		public void run () {
// 			if (!longPressFired) longPressFired = listener.longPress(pointer1.x, pointer1.y);
// 		}
// 	};

// 	/** Creates a new GestureDetector with default values: halfTapSquareSize=20, tapCountInterval=0.4f, longPressDuration=1.1f,
// 	 * maxFlingDelay=0.15f. */
// 	public GestureDetector (GestureListener listener) {
// 		this(20, 0.4f, 1.1f, 0.15f, listener);
// 	}

// 	/** @param halfTapSquareSize half width in pixels of the square around an initial touch event, see
// 	 *           {@link GestureListener#tap(float, float, int, int)}.
// 	 * @param tapCountInterval time in seconds that must pass for two touch down/up sequences to be detected as consecutive taps.
// 	 * @param longPressDuration time in seconds that must pass for the detector to fire a
// 	 *           {@link GestureListener#longPress(float, float)} event.
// 	 * @param maxFlingDelay time in seconds the finger must have been dragged for a fling event to be fired, see
// 	 *           {@link GestureListener#fling(float, float, int)}
// 	 * @param listener May be null if the listener will be set later. */
// 	public GestureDetector (float halfTapSquareSize, float tapCountInterval, float longPressDuration, float maxFlingDelay,
// 		GestureListener listener) {
// 		this.tapSquareSize = halfTapSquareSize;
// 		this.tapCountInterval = (long)(tapCountInterval * 1000000000l);
// 		this.longPressSeconds = longPressDuration;
// 		this.maxFlingDelay = (long)(maxFlingDelay * 1000000000l);
// 		this.listener = listener;
// 	}

// 	@Override
// 	public boolean touchDown (int x, int y, int pointer, int button) {
// 		return touchDown((float)x, (float)y, pointer, button);
// 	}

// 	public boolean touchDown (float x, float y, int pointer, int button) {
// 		if (pointer > 1) return false;

// 		if (pointer == 0) {
// 			pointer1.set(x, y);
// 			gestureStartTime = Gdx.input.getCurrentEventTime();
// 			tracker.start(x, y, gestureStartTime);
// 			if (Gdx.input.isTouched(1)) {
// 				// Start pinch.
// 				inTapSquare = false;
// 				pinching = true;
// 				initialPointer1.set(pointer1);
// 				initialPointer2.set(pointer2);
// 				longPressTask.cancel();
// 			} else {
// 				// Normal touch down.
// 				inTapSquare = true;
// 				pinching = false;
// 				longPressFired = false;
// 				tapSquareCenterX = x;
// 				tapSquareCenterY = y;
// 				if (!longPressTask.isScheduled()) Timer.schedule(longPressTask, longPressSeconds);
// 			}
// 		} else {
// 			// Start pinch.
// 			pointer2.set(x, y);
// 			inTapSquare = false;
// 			pinching = true;
// 			initialPointer1.set(pointer1);
// 			initialPointer2.set(pointer2);
// 			longPressTask.cancel();
// 		}
// 		return listener.touchDown(x, y, pointer, button);
// 	}

// 	@Override
// 	public boolean touchDragged (int x, int y, int pointer) {
// 		return touchDragged((float)x, (float)y, pointer);
// 	}

// 	public boolean touchDragged (float x, float y, int pointer) {
// 		if (pointer > 1) return false;
// 		if (longPressFired) return false;

// 		if (pointer == 0)
// 			pointer1.set(x, y);
// 		else
// 			pointer2.set(x, y);

// 		// handle pinch zoom
// 		if (pinching) {
// 			if (listener != null) {
// 				boolean result = listener.pinch(initialPointer1, initialPointer2, pointer1, pointer2);
// 				return listener.zoom(initialPointer1.dst(initialPointer2), pointer1.dst(pointer2)) || result;
// 			}
// 			return false;
// 		}

// 		// update tracker
// 		tracker.update(x, y, Gdx.input.getCurrentEventTime());

// 		// check if we are still tapping.
// 		if (inTapSquare && !isWithinTapSquare(x, y, tapSquareCenterX, tapSquareCenterY)) {
// 			longPressTask.cancel();
// 			inTapSquare = false;
// 		}

// 		// if we have left the tap square, we are panning
// 		if (!inTapSquare) {
// 			panning = true;
// 			return listener.pan(x, y, tracker.deltaX, tracker.deltaY);
// 		}

// 		return false;
// 	}

// 	@Override
// 	public boolean touchUp (int x, int y, int pointer, int button) {
// 		return touchUp((float)x, (float)y, pointer, button);
// 	}

// 	public boolean touchUp (float x, float y, int pointer, int button) {
// 		if (pointer > 1) return false;

// 		// check if we are still tapping.
// 		if (inTapSquare && !isWithinTapSquare(x, y, tapSquareCenterX, tapSquareCenterY)) inTapSquare = false;

// 		boolean wasPanning = panning;
// 		panning = false;

// 		longPressTask.cancel();
// 		if (longPressFired) return false;

// 		if (inTapSquare) {
// 			// handle taps
// 			if (lastTapButton != button || lastTapPointer != pointer || TimeUtils.nanoTime() - lastTapTime > tapCountInterval
// 				|| !isWithinTapSquare(x, y, lastTapX, lastTapY)) tapCount = 0;
// 			tapCount++;
// 			lastTapTime = TimeUtils.nanoTime();
// 			lastTapX = x;
// 			lastTapY = y;
// 			lastTapButton = button;
// 			lastTapPointer = pointer;
// 			gestureStartTime = 0;
// 			return listener.tap(x, y, tapCount, button);
// 		}

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

// 		// handle no longer panning
// 		boolean handled = false;
// 		if (wasPanning && !panning) handled = listener.panStop(x, y, pointer, button);

// 		// handle fling
// 		gestureStartTime = 0;
// 		long time = Gdx.input.getCurrentEventTime();
// 		if (time - tracker.lastTime < maxFlingDelay) {
// 			tracker.update(x, y, time);
// 			handled = listener.fling(tracker.getVelocityX(), tracker.getVelocityY(), button) || handled;
// 		}
// 		return handled;
// 	}

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

// 	private boolean isWithinTapSquare (float x, float y, float centerX, float centerY) {
// 		return Math.abs(x - centerX) < tapSquareSize && Math.abs(y - centerY) < tapSquareSize;
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

// 	/** Register an instance of this class with a {@link GestureDetector} to receive gestures such as taps, long presses, flings,
// 	 * panning or pinch zooming. Each method returns a boolean indicating if the event should be handed to the next listener (false
// 	 * to hand it to the next listener, true otherwise).
// 	 * @author mzechner */
// 	public static interface GestureListener {
// 		/** @see InputProcessor#touchDown(int, int, int, int) */
// 		public boolean touchDown (float x, float y, int pointer, int button);

// 		/** Called when a tap occured. A tap happens if a touch went down on the screen and was lifted again without moving outside
// 		 * of the tap square. The tap square is a rectangular area around the initial touch position as specified on construction
// 		 * time of the {@link GestureDetector}.
// 		 * @param count the number of taps. */
// 		public boolean tap (float x, float y, int count, int button);

// 		public boolean longPress (float x, float y);

// 		/** Called when the user dragged a finger over the screen and lifted it. Reports the last known velocity of the finger in
// 		 * pixels per second.
// 		 * @param velocityX velocity on x in seconds
// 		 * @param velocityY velocity on y in seconds */
// 		public boolean fling (float velocityX, float velocityY, int button);

// 		/** Called when the user drags a finger over the screen.
// 		 * @param deltaX the difference in pixels to the last drag event on x.
// 		 * @param deltaY the difference in pixels to the last drag event on y. */
// 		public boolean pan (float x, float y, float deltaX, float deltaY);

// 		/** Called when no longer panning. */
// 		public boolean panStop (float x, float y, int pointer, int button);

// 		/** Called when the user performs a pinch zoom gesture. The original distance is the distance in pixels when the gesture
// 		 * started.
// 		 * @param initialDistance distance between fingers when the gesture started.
// 		 * @param distance current distance between fingers. */
// 		public boolean zoom (float initialDistance, float distance);

// 		/** Called when a user performs a pinch zoom gesture. Reports the initial positions of the two involved fingers and their
// 		 * current positions.
// 		 * @param initialPointer1
// 		 * @param initialPointer2
// 		 * @param pointer1
// 		 * @param pointer2 */
// 		public boolean pinch (Vector2 initialPointer1, Vector2 initialPointer2, Vector2 pointer1, Vector2 pointer2);
// 	}

// 	/** Derrive from this if you only want to implement a subset of {@link GestureListener}.
// 	 * @author mzechner */
// 	public static class GestureAdapter implements GestureListener {
// 		@Override
// 		public boolean touchDown (float x, float y, int pointer, int button) {
// 			return false;
// 		}

// 		@Override
// 		public boolean tap (float x, float y, int count, int button) {
// 			return false;
// 		}

// 		@Override
// 		public boolean longPress (float x, float y) {
// 			return false;
// 		}

// 		@Override
// 		public boolean fling (float velocityX, float velocityY, int button) {
// 			return false;
// 		}

// 		@Override
// 		public boolean pan (float x, float y, float deltaX, float deltaY) {
// 			return false;
// 		}

// 		@Override
// 		public boolean panStop (float x, float y, int pointer, int button) {
// 			return false;
// 		}

// 		@Override
// 		public boolean zoom (float initialDistance, float distance) {
// 			return false;
// 		}

// 		@Override
// 		public boolean pinch (Vector2 initialPointer1, Vector2 initialPointer2, Vector2 pointer1, Vector2 pointer2) {
// 			return false;
// 		}
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
