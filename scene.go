// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"log"
	"time"

	"github.com/pyros2097/spike/input/gesture"
	"github.com/pyros2097/spike/scene2d"
)

var (
	allScenes    map[string]*Scene
	currentScene *Scene
)

type InputEvent struct {
}

// An InputProcessor is used to receive input events from the keyboard and the touch screen (mouse on the desktop). For this it
// has to be registered with the {@link Input#setInputProcessor(InputProcessor)} method. It will be called each frame before the
// call to {@link ApplicationListener#render()}. Each method returns a boolean in case you want to use this with the
// {@link InputMultiplexer} to chain input processors.
type Scene struct {
	scene2d.Actor
	Name     string
	OnPause  func()
	OnResume func()
	OnClick  func(x, y float32)

	// Called when the screen was touched or a mouse button was pressed. The button parameter will be {@link Buttons#LEFT} on iOS.
	// param screenX The x coordinate, origin is in the upper left corner
	// param screenY The y coordinate, origin is in the upper left corner
	// param pointer the pointer for the event.
	// param button the button
	// return whether the input was processed
	/** Called when a mouse button or a finger touch goes down on the actor. If true is returned, this listener will receive all
	 * touchDragged and touchUp events, even those not over this actor, until touchUp is received. Also when true is returned, the
	 * event is {@link Event#handle() handled}.
	 * @see InputEvent */
	OnTouchDown func(x, y float32, pointer, button int)

	// Called when a finger was lifted or a mouse button was released. The button parameter will be {@link Buttons#LEFT} on iOS.
	// param pointer the pointer for the event.
	// param button the button
	// return whether the input was processed

	/** Called when a mouse button or a finger touch goes up anywhere, but only if touchDown previously returned true for the mouse
	 * button or touch. The touchUp event is always {@link Event#handle() handled}.
	 * @see InputEvent */
	OnTouchUp func(x, y float32, pointer, button int)

	// Called when a finger or the mouse was dragged.
	// param pointer the pointer for the event.
	// @return whether the input was processed
	/** Called when a mouse button or a finger touch is moved anywhere, but only if touchDown previously returned true for the mouse
	 * button or touch. The touchDragged event is always {@link Event#handle() handled}.
	 * @see InputEvent */
	OnTouchDragged func(x, y float32, pointer int)

	// Called when a swipe gesture occurs
	OnGesture func(gtype gesture.Type)

	// Called when a key was typed
	// param character The character
	// return whether the input was processed
	OnKeyTyped func(key uint8)

	// Called when a key was released
	// param keycode one of the constants in {@link Input.Keys}
	// return whether the input was processed
	// When true is returned, the event is {@link Event#handle() handled
	OnKeyUp func(keycode uint8)

	// Called when a key was pressed
	// param keycode one of the constants in {@link Input.Keys}
	// return whether the input was processed
	OnKeyDown func(event *InputEvent, keycode uint8)

	// Called when the mouse was moved without any buttons being pressed. Will not be called on iOS.
	// @return whether the input was processed
	// This event only occurs on the desktop
	// public boolean mouseMoved (int screenX, int screenY);

	// Called when the mouse wheel was scrolled. Will not be called on iOS.
	// param amount the scroll amount, -1 or 1 depending on the direction the wheel was scrolled.
	// @return whether the input was processed.
	// public boolean scrolled (int amount);

	/** Called any time the mouse cursor or a finger touch is moved over an actor. On the desktop, this event occurs even when no
	 * mouse buttons are pressed (pointer will be -1).
	 * @param fromActor May be null.
	 * @see InputEvent */
	// public void enter (InputEvent event, float x, float y, int pointer, Actor fromActor) {
	// }

	/** Called any time the mouse cursor or a finger touch is moved out of an actor. On the desktop, this event occurs even when no
	 * mouse buttons are pressed (pointer will be -1).
	 * @param toActor May be null.
	 * @see InputEvent */
	// public void exit (InputEvent event, float x, float y, int pointer, Actor toActor) {
	// }

	/** Called when the mouse wheel has been scrolled. When true is returned, the event is {@link Event#handle() handled}. */
	// public boolean scrolled (InputEvent event, float x, float y, int amount) {
	// 	return false;
	// }

	BeforeShow    func()
	BeforeHide    func()
	AfterShow     func()
	AfterHide     func()
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
	log.Println("Adding Scene: " + scene.Name)
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
	log.Println("Setting Scene: " + name)
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
