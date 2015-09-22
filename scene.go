// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gdx

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pyros2097/spike/scene2d"
)

type GestureType int

const (
	SwipeLeft GestureType = iota
	SwipeRight
	SwipeUp
	SwipeDown
	SwipeNone
)

type Scene struct {
	scene2d.Actor
	Name           string
	OnPause        func()
	OnResume       func()
	OnClick        func(x, y float32)
	OnTouchDown    func(x, y float32, pointer, button int)
	OnTouchUp      func(x, y float32, pointer, button int)
	OnTouchDragged func(x, y float32, pointer int)
	OnGesture      func(gtype GestureType)
	OnKeyTyped     func(key int8)
	OnKeyUp        func(keycode int)
	OnKeyDown      func(keycode int)
	BeforeShow     func()
	BeforeHide     func()
	AfterShow      func()
	AfterHide      func()
	TransitionIn   func(scene *Scene)
	TransitionOut  func(scene *Scene)
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

// Set the current scene to be displayed
func SetScene(name string) {
	log.Info("Setting Scene: " + name)
	var ok bool
	if currentScene, ok = allScenes[name]; ok {
		// currentScene.Transition()
	}
}

// Set the current scene to be displayed with an amount of delay
func SetSceneWithDelay(name string, duration time.Duration) {
	time.AfterFunc(duration, func() {
		SetScene(name)
	})
}

// Returns the current scene being Displayed on stage
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

// private static Image imgbg = null;
// 	public void setBackground(String texName) {
// 		if(imgbg != null)
// 			removeBackground();
// 		if(Asset.tex(texName) != null){
// 			imgbg = new Image(new TextureRegionDrawable(Asset.tex(texName)), Scaling.stretch);
// 			imgbg.setFillParent(true);
// 			stage2d.addActor(imgbg);
// 			imgbg.toBack();
// 		}
// 	}

// 	public void removeBackground() {
// 		getRoot().removeActor(imgbg);
// 	}
