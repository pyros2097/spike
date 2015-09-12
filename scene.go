// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gdx

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pyros2097/gdx/scene2d"
)

type Scene struct {
	scene2d.Actor
	Name           string
	onPause        func()
	onResume       func()
	onExit         func()
	onClick        func(x, y float32)
	onTouchDown    func(x, y float32, pointer, button int)
	onTouchUp      func(x, y float32, pointer, button int)
	onTouchDragged func(x, y float32, pointer int)
	onSwipeLeft    func()
	onSwipeRight   func()
	onSwipeDown    func()
	onSwipeUp      func()
	beforeShow     func()
	beforeHide     func()
	afterShow      func()
	afterHide      func()
	transitionIn   func(scene *Scene)
	transitionOut  func(scene *Scene)
}

var (
	allScenes    map[string]*Scene
	currentScene *Scene
)

func init() {
	allScenes = make(map[string]*Scene)
}

func AddScene(scene *Scene) {
	if len(scene.Name) == 0 {
		panic("You must specify a scene name")
	}
	log.Info("Adding Scene: " + scene.Name)
	allScenes[scene.Name] = scene
	if currentScene == nil {
		currentScene = scene
	}
}

func RemoveScene(name string) {
	delete(allScenes, name)
}

/**
 * Set the current scene to be displayed
 * @param className The scene's Class Name
 **/
func SetScene(name string) {
	log.Info("Setting Scene: " + name)
	var ok bool
	if currentScene, ok = allScenes[name]; ok {
		// currentScene.Transition()
	}
}

func SetSceneWithDelay(name string, duration time.Duration) {
	time.AfterFunc(duration, func() {
		SetScene(name)
	})
}

func GetCurrentScene() *Scene {
	return currentScene
}

func GetScene(name string) *Scene {
	return allScenes[name]
}

func AddHud() {

}

func DrawGrid() {
}

func DrawSelection() {
}
