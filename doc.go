// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Spike is a Game Framework built on top of golang mobile with all the batteries included. It is a lot similar to
// libgdx but has lesser abstractions and easier integration and setup. It favors composition over inheritance and
// also has a declarative interface something like React.
// It has all configuration for assets, sound, music, textures, animations already setup.
// You can directly start coding your game without initializing complicated things.
// It follows the Rails design of convention over configuration.
// It consists of Scenes and you can control which scenes are to displayed when.
// Be sure to see the example source code as it can help you a lot in understanding the framework.
//
// # Requires
//
// 1.Go >= 1.5
//
// 2.OpenGL >= 2.0
//
// 3.libopenal-dev
//
// 4.Linux
//
// # Assets
//
// Note: All asset files must be lowercase only.. otherwise it causes problems with android.
// All Assets are to be stored in the assets directory.
// For Automatic Asset Loading the directory Structure should be like this
//
//   assets/icons/icon.png - your game icon which is loaded by the framework
//   assets/atlas/ --- all your Texture Atlas files .atlas and .png go here
//   assets/fonts/ --- all your BitmapFont files .fnt and .png go here
//   assets/musics/ --- all your Music files .mp3 go here
//   assets/sounds/ --- all your Music files .mp3 go here
//   assets/particles/ --- all your Particle files .part go here
//   assets/maps/ --- all your TMX map files .tmx go here
//
// # Usage
//
// menu := &spike.Scene{
//   Name:    "Menu",
//   BGColor: graphics.Color{0, 0, 1, 1},
//   Children: []*spike.Actor{
//     &spike.Actor{
//       X: 43,
//       OnInput: func(self *spike.Actor, event spike.InputEvent) {
//         println(event.Type.String())
//       },
//       OnAct: func(self *spike.Actor, delta float32) {
//         self.X = 111
//       },
//     },
//   },
// }
// options := &spike.Scene{Name: "Options", BGColor: graphics.Color{0, 0, 0, 1}}
// spike.Init("example", 800, 480)
// spike.AddScene(menu)
// spike.AddScene(options)
// spike.PlaySound("boing")
// spike.Run()
package spike
