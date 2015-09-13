// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Gdx is a Game Framework built on top of golang mobile with all the batteries included.
// It has all configuration for assets, sound, music, textures, animations already setup.
// You can directly start coding your game without initializing complicated things.
// It follows the rails design of convention over configuration.
// It consists of Scenes and has a SceneManager to control which scenes are to displayed when.
// Be sure to see the exmaples source as it can help you a lot in understanding the framework.
//
// Assets
// ------
// Important All asset files must be lowercase only.. otherwise it causes problems with android.
// All Assets are to be stored in the assets directory.
// Automatic Asset Loading the directory Structure should be like this
//
// assets/icons/icon.png - your game icon which is loaded by the framework
//
// assets/atlas/ --- all your Texture Atlas files .atlas and .png go here
//
// assets/fonts/ --- all your BitmapFont files .fnt and .png go here
//
// assets/musics/ --- all your Music files .mp3 go here
//
// assets/sounds/ --- all your Music files .mp3 go here
//
// assets/particles/ --- all your Particle files .part go here
//
// assets/maps/ --- all your TMX map files .tmx go here
package gdx
