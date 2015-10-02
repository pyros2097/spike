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
	"github.com/pyros2097/spike/graphics"
	"github.com/pyros2097/spike/scene2d"
)

var (
	hello *scene2d.Actor
)

func main() {

	menu := &spike.Scene{
		Name:    "menu",
		BGColor: &graphics.Color{0, 0, 1, 1},
		Children: []*scene2d.Actor{
			&scene2d.Actor{
				X: 43,
				OnTouchDown: func(self *scene2d.Actor, x, y float32, p, b int) {
					println(x)
					println(y)
				},
				OnTap: func(self *scene2d.Actor, x, y float32, p, b int) {
					println("MEFAE ERAPPP")
				},
				OnAct: func(self *scene2d.Actor, delta float32) {
					// print(self.X)
					self.X = 111
				},
			},
		},
	}
	options := &spike.Scene{Name: "options", BGColor: &graphics.Color{0, 0, 0, 1}}
	spike.Init("example", 800, 480)
	spike.AddScene(menu)
	spike.AddScene(options)
	spike.PlaySound("boing")
	spike.Run()
}

// type Actor struct {
// 	x, y int
// }

// var children []*Actor

// func main() {
// 	child := &Actor{4, 5}
// 	children = append(children, child)
// 	for _, c := range children {
// 		println(c.x)
// 		c.x = 444
// 	}
// 	for _, c := range children {
// 		println(c.x)
// 		c.x = 12
// 	}
// 	println(child.x)
// }
