// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux

// An app that demonstrates the 2d scenegraph
// Get the basic example and use gomobile to build or install it on your device.
//
//   $ go get -d github.com/pyros2097/spike/example
//   $ gomobile build github.com/pyros2097/spike/example # will build an APK
//
//   # plug your Android device to your computer or start an Android emulator.
//   # if you have adb installed on your machine, use gomobile install to
//   # build and deploy the APK to an Android target.
//   $ gomobile install github.com/pyros2097/spike/example
//
// Switch to your device or emulator to start the Basic application from
// the launcher.
// You can also run the application on your desktop by running the command
// below. (Note: It currently doesn't work on Windows.)
//   $ go install github.com/pyros2097/spike/example && example
package main

import (
	"github.com/pyros2097/spike"
)

func main() {
	menu := &spike.Scene{
		Name:    "Menu",
		BGColor: spike.Color{0, 0, 1, 1},
		Actor: spike.Actor{
			Children: []*spike.Actor{
				{
					X:     43,
					Y:     99,
					Color: &spike.Color{0, 0, 0, 0},
					Init: func(self *spike.Actor) {
						self.ActionAlpha(0.5, 300000000, nil)
						// self.Delay(5000).Run(func() {
						// 	println("estings")
						// }).AMoveTo(22, 10, 0, nil)
					},
					Input: func(self *spike.Actor, event spike.InputEvent) {
						// println(event.Type.String())
					},
					Act: func(self *spike.Actor, delta float32) {
						// println(self.X)
						println(self.Color.A)
						self.X = 111
					},
				},
			},
		},
	}
	options := &spike.Scene{Name: "Options", BGColor: spike.Color{0, 0, 0, 1}}
	spike.Init("example", 800, 480)
	spike.AddScene(menu)
	spike.AddScene(options)
	spike.PlaySound("boing")
	spike.Run()
}
