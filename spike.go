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
	PauseState   bool

	Camera2d *Camera
	Camera3d *Camera

	running   bool
	fpsTicker *time.Ticker
	lastTime  time.Duration
	deltaTime time.Duration

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
			deltaTime = now.Sub(lastTime)
			lastTime = now
			var sz size.Event
			for e := range a.Events() {
				switch e := app.Filter(e).(type) {
				case lifecycle.Event:
					switch e.Crosses(lifecycle.StageVisible) {
					case lifecycle.CrossOn:
						appStart()
					case lifecycle.CrossOff:
						appExit()
					}
				case size.Event: // resize event
					sz = e
					touchX = float32(sz.WidthPx / 2)
					touchY = float32(sz.HeightPx / 2)
				case paint.Event:
					// print("Painting")
					// Do update here
					appPaint(sz, float32(deltaTime))
					a.EndPaint(e)
				case touch.Event:
					// print("Touching")
					// send input events here or before paint just store the last state
					touchX = e.X
					touchY = e.Y
					switch e.Type {
					case touch.TypeBegin:
						// println("Begin")
						doTouchDown(touchX, touchY, 0, 0)
					case touch.TypeEnd:
						// println("End")
						doTouchUp(touchX, touchY, 0, 0)
					case touch.TypeMove:
						// println("Moving")
						doTouchDragged(touchX, touchY, 0)
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

func appStart() {
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

	// TODO(crawshaw): the debug package needs to put GL state init here
	// Can this be an app.RegisterFilter call now??
	SetScene(currentScene.Name)
}

// This is the main rendering call that updates the time, updates the stage,
// loads assets asynchronously, updates the camera and FPS text.
func appPaint(sz size.Event, delta float32) {
	gl.ClearColor(currentScene.BGColor.R, currentScene.BGColor.G, currentScene.BGColor.B, currentScene.BGColor.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for _, child := range currentScene.Children {
		if child.OnAct != nil {
			child.OnAct(child, delta)
		}
		if child.OnDraw != nil {
			child.OnDraw(child, tempBatch, 1.0)
		}
		if len(InputChannel) > 0 {
			for e := range InputChannel {
				if child.OnInput != nil {
					child.OnInput(child, e)
				}
				if len(InputChannel) == 0 {
					break
				}
			}
		}
	}

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

func appPause() {
	println("Pausing")
	PauseState = true
	musicPlayer.Pause()
	soundsPlayer.Pause()
	// unloadAll()
	if currentScene.OnPause != nil {
		currentScene.OnPause(currentScene)

	}
}

func appResume() {
	println("Resuming")
	PauseState = false
	musicPlayer.Pause()
	soundsPlayer.Play()
	// reloadAll()
	if currentScene.OnResume != nil {
		currentScene.OnResume(currentScene)
	}
}

// Use this to exit your game safely
// It will automatically unload all your assets and dispose the stage
// Schedule an exit from the application. On android, this will cause a call to pause() and dispose() some time in the future,
// it will not immediately finish your application.
// On iOS this should be avoided in production as it breaks Apples guidelines
func appExit() {
	println("Exiting")
	running = false
	if currentScene.OnPause != nil {
		currentScene.OnPause(currentScene)
		soundsPlayer.Close()
	}
	gl.DeleteProgram(program)
	gl.DeleteBuffer(buf)
}

// /***********************************************************************************************************
// * 					Utilities Related Functions												   	       		   *
// ************************************************************************************************************/
// 	/*
// 	 * The the angle in degrees of the inclination of a line
// 	 * @param cx, cy The center point x, y
// 	 * @param tx, ty The target point x, y
// 	 */
// 	public static float getAngle(float cx, float cy, float tx, float ty) {
// 		float angle = (float) Math.toDegrees(MathUtils.atan2(tx - cx, ty - cy));
// 		//if(angle < 0){
// 		//    angle += 360;
// 		//}
// 		return angle;
// 	}

// 	private static Vector2 distVector = new Vector2();
// 	public static final float getDistance(Actor a, Actor b){
// 		distVector.set(a.getX(), a.getY());
// 		return distVector.dst(b.getX(), b.getY());
// 	}

// 	public static final float getDistance(float x1, float y1, float x2, float y2){
// 		distVector.set(x1, y1);
// 		return distVector.dst(x2, y2);
// 	}

// 	/*
// 	 * Capitalizes the First Letter of a String
// 	 */
// 	public static String capitalize(String text){
// 		if(text != null && text != "")
// 			return (text.substring(0, 1)).toUpperCase() + text.substring(1);
// 		else
// 			return "";
// 	}

// 	/*
// 	 * UnCapitalizes the First Letter of a String
// 	 */
// 	public static String uncapitalize(String text){
// 		return text.substring(0, 1).toLowerCase()+text.substring(1);
// 	}

// 	public static Rectangle getBounds(Actor actor) {
// 		float posX = actor.getX();
// 		float posY = actor.getY();
// 		float width = actor.getWidth();
// 		float height = actor.getHeight();
// 		return new Rectangle(posX, posY, width, height);
// 	}

// 	public static boolean collides(Actor actor, float x, float y) {
// 		Rectangle rectA1 = getBounds(actor);
// 		Rectangle rectA2 = new Rectangle(x, y, 5, 5);
// 		// Check if rectangles collides
// 		if (Intersector.overlaps(rectA1, rectA2)) {
// 			return true;
// 		} else {
// 			return false;
// 		}
// 	}

// 	public static boolean collides(Actor actor1, Actor actor2) {
// 		Rectangle rectA1 = getBounds(actor1);
// 		Rectangle rectA2 = getBounds(actor2);
// 		// Check if rectangles collides
// 		if (Intersector.overlaps(rectA1, rectA2)) {
// 			return true;
// 		} else {
// 			return false;
// 		}
// 	}

// 	/**
// 	 * Get screen time from start in format of HH:MM:SS. It is calculated from
// 	 * "secondsTime" parameter.
// 	 * */
// 	public static String toScreenTime(float secondstime) {
// 		int seconds = (int)(secondstime % 60);
// 		int minutes = (int)((secondstime / 60) % 60);
// 		int hours =  (int)((secondstime / 3600) % 24);
// 		return new String(addZero(hours) + ":" + addZero(minutes) + ":" + addZero(seconds));
// 	}

// 	private static String addZero(int value){
// 		String str = "";
// 		if(value < 10)
// 			 str = "0" + value;
// 		else
// 			str = "" + value;
// 		return str;
// 	}

// 	private static ShapeRenderer shapeRenderer;
// 	private final static Actor selectionBox = new Actor(){
// 		@Override
// 		public void draw(Batch batch, float alpha){
// 			batch.end();
// 			shapeRenderer.begin(ShapeType.Line);
// 			shapeRenderer.setColor(Color.GREEN);
// 			shapeRenderer.rect(getX(), getY(),
// 						getWidth(),getHeight());
// 			shapeRenderer.end();
// 			batch.begin();
// 		}
// 	};
