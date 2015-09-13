// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gdx

import (
	log "github.com/Sirupsen/logrus"
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
)

func InitAssets(config *AssetConfig) {
}

func SoundPlay(name string) {
	log.Info("Playing Sound: " + name)
}

func SoundPause() {

}

func SoundResume() {

}

func SoundStop(name string) {
	log.Infof("Stopping Sound: " + name)
}

func SoundDispose() {

}

func Font() {

}

func Tex() {

}

func loadTmx() {

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
