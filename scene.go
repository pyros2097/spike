// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"time"

	"github.com/pyros2097/spike/graphics"
	"github.com/pyros2097/spike/scene2d"
)

var (
	allScenes    map[string]*Scene
	currentScene *Scene
)

// TODO: ADD props Validation for Name, BGColor etc
// An InputProcessor is used to receive input events from the keyboard and the touch screen (mouse on the desktop). For this it
// has to be registered with the {@link Input#setInputProcessor(InputProcessor)} method. It will be called each frame before the
// call to {@link ApplicationListener#render()}. Each method returns a boolean in case you want to use this with the
// {@link InputMultiplexer} to chain input processors.
type Scene struct {
	scene2d.Actor
	Name    string
	BGColor *graphics.Color

	Children []*scene2d.Actor

	OnPause  func(self *Scene)
	OnResume func(self *Scene)

	BeforeShow    func(self *Scene)
	BeforeHide    func(self *Scene)
	AfterShow     func(self *Scene)
	AfterHide     func(self *Scene)
	TransitionIn  func(scene *Scene)
	TransitionOut func(scene *Scene)
}

/** EventListener for low-level input events. Unpacks {@link InputEvent}s and calls the appropriate method. By default the methods
 * here do nothing with the event. Users are expected to override the methods they are interested in, like this:
 *
 * <pre>
 * actor.addListener(new InputListener() {
 * 	public boolean touchDown (InputEvent event, float x, float y, int pointer, int button) {
 * 		Gdx.app.log(&quot;Example&quot;, &quot;touch started at (&quot; + x + &quot;, &quot; + y + &quot;)&quot;);
 * 		return false;
 * 	}
 *
 * 	public void touchUp (InputEvent event, float x, float y, int pointer, int button) {
 * 		Gdx.app.log(&quot;Example&quot;, &quot;touch done at (&quot; + x + &quot;, &quot; + y + &quot;)&quot;);
 * 	}
 * });
 * </pre> */
// public class InputListener implements EventListener {
// 	static private final Vector2 tmpCoords = new Vector2();

// 	public boolean handle (Event e) {
// 		if (!(e instanceof InputEvent)) return false;
// 		InputEvent event = (InputEvent)e;

// 		switch (event.getType()) {
// 		case keyDown:
// 			return keyDown(event, event.getKeyCode());
// 		case keyUp:
// 			return keyUp(event, event.getKeyCode());
// 		case keyTyped:
// 			return keyTyped(event, event.getCharacter());
// 		}

// 		event.toCoordinates(event.getListenerActor(), tmpCoords);

// 		switch (event.getType()) {
// 		case touchDown:
// 			return touchDown(event, tmpCoords.x, tmpCoords.y, event.getPointer(), event.getButton());
// 		case touchUp:
// 			touchUp(event, tmpCoords.x, tmpCoords.y, event.getPointer(), event.getButton());
// 			return true;
// 		case touchDragged:
// 			touchDragged(event, tmpCoords.x, tmpCoords.y, event.getPointer());
// 			return true;
// 		case mouseMoved:
// 			return mouseMoved(event, tmpCoords.x, tmpCoords.y);
// 		case scrolled:
// 			return scrolled(event, tmpCoords.x, tmpCoords.y, event.getScrollAmount());
// 		case enter:
// 			enter(event, tmpCoords.x, tmpCoords.y, event.getPointer(), event.getRelatedActor());
// 			return false;
// 		case exit:
// 			exit(event, tmpCoords.x, tmpCoords.y, event.getPointer(), event.getRelatedActor());
// 			return false;
// 		}
// 		return false;
// 	}

func (self *Scene) SetBackground(texName string) {
	// 		if(Asset.tex(texName) != null){
	// 			imgbg = new Image(new TextureRegionDrawable(Asset.tex(texName)), Scaling.stretch);
	// 			imgbg.setFillParent(true);
	// 			stage2d.addActor(imgbg);
	// 			imgbg.toBack();
	// 		}
}

func (self *Scene) RemoveBackground() {
	// self.RemoveActor()
}

func (self *Scene) AddHud() {
}

func (self *Scene) DrawGrid() {
}

func (self *Scene) DrawSelection() {
}

func init() {
	allScenes = make(map[string]*Scene)
}

func AddScene(scene *Scene) {
	println("Adding Scene: " + scene.Name)
	allScenes[scene.Name] = scene
	if currentScene == nil {
		currentScene = scene
	}
}

func RemoveScene(name string) {
	delete(allScenes, name)
}

// Sets the current scene to be displayed
func SetScene(name string) {
	println("Setting Scene: " + name)
	if currentScene != nil {
		if currentScene.BeforeHide != nil {
			currentScene.BeforeHide(currentScene)
		}
		if currentScene.AfterHide != nil {
			currentScene.AfterHide(currentScene)
		}
	}
	currentScene = allScenes[name]
	if currentScene.BeforeShow != nil {
		currentScene.BeforeShow(currentScene)
	}
	if currentScene.AfterShow != nil {
		currentScene.AfterShow(currentScene)
	}
	// setTouchable(Touchable.childrenOnly);
	// Camera.reset();
	// stage2d.clear();
	// stage3d.clear();
	// setSize(targetWidth, targetHeight);
	// setBounds(0,0, targetWidth, targetHeight);
	// setColor(1f, 1f, 1f, 1f);
	// setVisible(true);
	// stage2d.getRoot().setPosition(0, 0);
	// stage2d.getRoot().setVisible(true);
	// stage3d.getRoot().setPosition(0, 0, 0);
	// stage3d.getRoot().setVisible(true);
	// cullingEnabled = true;

	// stage2d.getRoot().setName("Root");
	// stage2d.getRoot().setTouchable(Touchable.childrenOnly);
	// stage2d.setCamera(new Camera());
	// stage3d = new Stage3d();
	// //camController = new CameraInputController(stage3d.getCamera());
	// Gdx.input.setCatchBackKey(true);
	// Gdx.input.setCatchMenuKey(true);
	// Gdx.input.setInputProcessor(inputMux);
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

func (self *Scene) AddActor(actor *scene2d.Actor) {
	self.Children = append(self.Children, actor)
}

func (self *Scene) AddActorWithDelay(actor *scene2d.Actor, duration time.Duration) {
	time.AfterFunc(duration, func() {
		self.AddActor(actor)
	})
}

func (self *Scene) RemoveActor(actor *scene2d.Actor) {
	i := actor.Z
	self.Children, self.Children[len(self.Children)-1] = append(self.Children[:i], self.Children[i+1:]...), nil
}

func (self *Scene) RemoveActorWithDelay(actor *scene2d.Actor, duration time.Duration) {
	// addAction(Actions.sequence(Actions.delay(delay), Actions.removeActor(actor)));
}

func (self *Scene) RemoveActorWithName(name string) {
	// return removeActor(findActor(actorName));
}

// func (self *Scene) Hit(x, y float32) *scene2d.Actor {
// }

// public Actor hit(float x, float y){
// 	return hit(x, y, true);
// }

// func AddActor3d() {
// }
