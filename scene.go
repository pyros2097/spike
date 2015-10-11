// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"time"

	"github.com/pyros2097/spike/graphics"
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
// private String sceneName = "";
// public String sceneBackground = "None";
// public String sceneMusic = "None";
// public String sceneTransition = "None";
// public float sceneDuration = 0;
// public InterpolationType sceneInterpolationType = InterpolationType.Linear;
// public static float splashDuration = 0f;
// public static boolean pauseState = false;

// 	public EffectType effectType = EffectType.None;
// 	public float effectValue = 0f;
// 	public float effectDuration = 0f;
// 	public InterpolationType interpolationType = InterpolationType.Linear;
// 	public float addActorDelay = 0f;
// 	public float addEffectDelay = 0f;

// 	public EventType evtType = EventType.None;
// 	public Scene.OnEventType onEvtType = Scene.OnEventType.DoEffect;
// 	public String eventScene = "";
type Scene struct {
	Group
	Name    string
	BGColor graphics.Color

	OnPause  func(self *Scene)
	OnResume func(self *Scene)

	BeforeShow    func(self *Scene)
	BeforeHide    func(self *Scene)
	AfterShow     func(self *Scene)
	AfterHide     func(self *Scene)
	TransitionIn  func(scene *Scene)
	TransitionOut func(scene *Scene)
}

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

func (self *Scene) AddActor(actor *Actor) {
	self.Children = append(self.Children, actor)
}

func (self *Scene) AddActorWithDelay(actor *Actor, duration time.Duration) {
	time.AfterFunc(duration, func() {
		self.AddActor(actor)
	})
}

func (self *Scene) RemoveActor(actor *Actor) {
	i := actor.Z
	self.Children, self.Children[len(self.Children)-1] = append(self.Children[:i], self.Children[i+1:]...), nil
}

func (self *Scene) RemoveActorWithDelay(actor *Actor, duration time.Duration) {
	// addAction(Actions.sequence(Actions.delay(delay), Actions.removeActor(actor)));
}

func (self *Scene) RemoveActorWithName(name string) {
	// return removeActor(findActor(actorName));
}

// func (self *Scene) Hit(x, y float32) *Actor {
// }

// public Actor hit(float x, float y){
// 	return hit(x, y, true);
// }

// func AddActor3d() {
// }

// 	public void showToast(String message, float duration){
// 		Table table = new Table(Asset.skin);
// 		table.add("   "+message+"   ");
// 		table.setBackground(Asset.skin.getDrawable("dialogDim"));
// 		table.pack();
// 		table.setPosition(Scene.targetWidth/2 - table.getWidth(), Scene.targetHeight/2 - table.getHeight());
// 		addActor(table);
// 		table.addAction(Actions.sequence(Actions.delay(duration), Actions.removeActor(table)));
// 	}

// 	public void showMessageDialog(String title, String message){
// 		Dialog dialog = new Dialog(title, Asset.skin);
// 		dialog.getContentTable().add(message);
// 		dialog.button("OK", "OK");
// 		dialog.pack();
// 		dialog.show(getStage());
// 	}

// 	public boolean showConfirmDialog(String title, String message){
// 		Dialog dialog = new Dialog(title, Asset.skin);
// 		dialog.button("Yes", "Yes");
// 		dialog.button("No", "No");
// 		dialog.pack();
// 		dialog.show(getStage());
// 		//if(dialog.result().equals("Yes")) FIXME update Gdx
// 		//	return true;
// 		return false;
// 	}

// 	public void outline(Actor actor){
// 		selectionBox.setPosition(actor.getX(), actor.getY());
// 		selectionBox.setSize(actor.getWidth(), actor.getHeight());
// 		stage2d.addActor(selectionBox);
// 	}
