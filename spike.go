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

	"github.com/pyros2097/spike/input/gesture"
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
	lastTouchEvent touch.Type
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
						gesture.TouchDown(touchX, touchY, 0, 0)
					case touch.TypeEnd: // touch up
						println("End")
						gesture.TouchUp(touchX, touchY, 0, 0)
					case touch.TypeMove: // touch drag
						println("Moving")
						gesture.TouchDragged(touchX, touchY, 0)
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
	// TODO(crawshaw): the debug package needs to put GL state init here
	// Can this be an app.RegisterFilter call now??
	SetScene(currentScene.Name)
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
		for event := gesture.GetEvent(); event != nil; {
			if event.Type == gesture.Tap && child.OnTap != nil {
				child.OnTap(child, touchX, touchY, 1, 0) //tapCount, 0)
			}
			event = nil
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
