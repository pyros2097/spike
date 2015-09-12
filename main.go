// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gdx

import (
	"encoding/binary"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pyros2097/gdx/graphics"

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
	targetWidth, targetHeight float32
	pauseState                bool
	camera2d                  *graphics.Camera

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
	log.Info("Initializing Gdx")
	targetWidth = width
	targetHeight = height
	camera2d = &graphics.Camera{}
	allScenes = make(map[string]*Scene)
	fpsTicker = time.NewTicker(1000 / 30 * time.Millisecond)
	running = true
	initConfig(title)
}

func Run() {
	lastTime := time.Now()
	app.Main(func(a app.App) {
		for now := range fpsTicker.C {
			print("timing")
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
						Start()
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
					touchX = e.X
					touchY = e.Y
				}
			}
		}
	})
}

func onStart() {
	program, err = glutil.CreateProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
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
}

func onStop() {
	gl.DeleteProgram(program)
	gl.DeleteBuffer(buf)
}

func onPaint(sz size.Event) {
	gl.ClearColor(1, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

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

func Start() {
	log.Info("Starting")
	SetScene(currentScene.Name)
}

func Load(sceneName string) {

}

func Save() {

}

func Pause() {
	log.Info("Pausing")
	pauseState = true
	if currentScene.onPause != nil {
		currentScene.onPause()
	}
}

func Resume() {
	log.Info("Resuming")
	pauseState = false
	if currentScene.onResume != nil {
		currentScene.onResume()
	}
}

/*
 * Use this to exit your game safely
 * It will automatically unload all your assets and dispose the stage
 */
/** Schedule an exit from the application. On android, this will cause a call to pause() and dispose() some time in the future,
 * it will not immediately finish your application.
 * On iOS this should be avoided in production as it breaks Apples guidelines*/
func Exit() {
	log.Info("Exiting")
	running = false
	if currentScene.onExit != nil {
		currentScene.onExit()
		// currentScene.onPause()
	}
}

// public Preferences getPreferences (String name);

// public Clipboard getClipboard ();
/** <p>
 * An <code>Application</code> is the main entry point of your project. It sets up a window and rendering surface and manages the
 * different aspects of your application, namely {@link Graphics}, {@link Audio}, {@link Input} and {@link Files}. Think of an
 * Application being equivalent to Swing's <code>JFrame</code> or Android's <code>Activity</code>.
 * </p>*/

// * <p>
// * Note that all {@link Music} instances will be automatically paused when the {@link ApplicationListener#pause()} method is
// * called, and automatically resumed when the {@link ApplicationListener#resume()} method is called.
// * </p>
