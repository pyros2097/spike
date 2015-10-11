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

var (
	targetWidth  float32
	targetHeight float32
	PauseState   bool
	// DebugColor Color{0, 1, 1, 1}

	running   bool
	fpsTicker *time.Ticker
	lastTime  time.Duration
	deltaTime time.Duration

	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green  float32
	touchX float32
	touchY float32

	err error
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
		var glctx gl.Context
		visible, sz := false, size.Event{}
		for now := range fpsTicker.C {
			if !running {
				break
			}
			deltaTime = now.Sub(lastTime)
			lastTime = now
			for e := range a.Events() {
				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					switch e.Crosses(lifecycle.StageVisible) {
					case lifecycle.CrossOn:
						visible = true
						glctx, _ = e.DrawContext.(gl.Context)
						appStart(glctx)
					case lifecycle.CrossOff:
						appStop(glctx)
					}
				case size.Event: // resize event
					sz = e
					touchX = float32(sz.WidthPx / 2)
					touchY = float32(sz.HeightPx / 2)
				case paint.Event:
					if visible {
						appPaint(glctx, sz, float32(deltaTime))
						a.Publish()
						// Keep animating.
						a.Send(paint.Event{})
					}
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

func appStart(glctx gl.Context) {
	println("Starting")
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		panic("error creating GL program: " + err.Error())
		return
	}

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = glctx.GetAttribLocation(program, "position")
	color = glctx.GetUniformLocation(program, "color")
	offset = glctx.GetUniformLocation(program, "offset")
	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)

	SetScene(currentScene.Name)
}

// Use this to exit your game safely
// It will automatically unload all your assets and dispose the stage
// Schedule an exit from the application. On android, this will cause a call to pause() and dispose() some time in the future,
// it will not immediately finish your application.
// On iOS this should be avoided in production as it breaks Apples guidelines
func appStop(glctx gl.Context) {
	println("Exiting")
	running = false
	if currentScene.OnPause != nil {
		currentScene.OnPause(currentScene)
		soundsPlayer.Close()
	}
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
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

// This is the main rendering call that updates the current scene and all children in the scene
func appPaint(glctx gl.Context, sz size.Event, delta float32) {
	glctx.ClearColor(currentScene.BGColor.R, currentScene.BGColor.G, currentScene.BGColor.B, currentScene.BGColor.A)
	glctx.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for _, child := range currentScene.Children {
		child.act(delta)
		child.draw(tempBatch, 1.0)
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

	glctx.UseProgram(program)

	green += 0.01
	if green > 1 {
		green = 0
	}
	glctx.Uniform4f(color, 0, green, 0, 1)

	glctx.Uniform2f(offset, touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))

	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.EnableVertexAttribArray(position)
	glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	glctx.DisableVertexAttribArray(position)

	fps.Draw(sz)
}

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, 0.4, 0.0, // top left
	0.0, 0.0, 0.0, // bottom left
	0.4, 0.0, 0.0, // bottom right
)

const (
	coordsPerVertex = 3
	vertexCount     = 3
)

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
