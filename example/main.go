// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux

// An app that demonstrates the 2d scenegraph
// Get the basic example and use gomobile to build or install it on your device.
//
//   $ go get -d github.com/pyros2097/gdx/example
//   $ gomobile build github.com/pyros2097/gdx/example # will build an APK
//
//   # plug your Android device to your computer or start an Android emulator.
//   # if you have adb installed on your machine, use gomobile install to
//   # build and deploy the APK to an Android target.
//   $ gomobile install github.com/pyros2097/gdx/example
//
// Switch to your device or emulator to start the Basic application from
// the launcher.
// You can also run the application on your desktop by running the command
// below. (Note: It currently doesn't work on Windows.)
//   $ go install github.com/pyros2097/gdx/example && example
package main

import (
	"github.com/pyros2097/gdx"
	"github.com/pyros2097/gdx/scene2d"
)

var (
	hello *scene2d.Actor
)

func main() {
	menu := &gdx.Scene{Name: "menu"}
	gdx.Init("example", 800, 480)
	gdx.AddScene(menu)
	gdx.Run()
}
