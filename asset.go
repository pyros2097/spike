// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/audio"
)

type (
	AssetConfig struct {
		images     []string
		sounds     []string
		musics     []string
		fonts      []string
		animations []string
		tmx        []string
	}
)

var (
	imagesMap     map[string]int
	soundsMap     map[string]int
	musicsMap     map[string]int
	fontsMap      map[string]int
	animationsMap map[string]int
	musicPlayer   *audio.Player
	soundsPlayer  *audio.Player
)

func InitAssets(config *AssetConfig) {
}

func PlaySound(name string) {
	println("Playing Sound: " + name)
	rc, err := asset.Open("sound/" + name + ".wav")
	if err != nil {
		panic(err)
	}
	soundsPlayer, err = audio.NewPlayer(rc, 0, 0)
	if err != nil {
		panic(err)
	}
	soundsPlayer.Seek(0)
	soundsPlayer.Play()
	// player.Close()
}

func StopSound(name string) {
	soundsPlayer.Stop()
	println("Stopping Sound: " + name)
}

func PlayMusic() {
}

func StopMusic() {
}

func Font() {
}

func Tex() {

}

func LoadTmx() {

}

func UnloadTmx() {

}

func UnloadAll() {
	// for k, v := range imagesMap {
	// }
	// for k, v := range fontsMap {
	// }
	// for k, v := range musicsMap {
	// }
	// for k, v := range soundsMap {
	// }
	// for k, v := range animationsMap {
	// }
}
