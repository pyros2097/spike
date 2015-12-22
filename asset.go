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

// 	public static void musicPlay(String filename){
// 		if(Config.isMusic){
// 			if(currentMusic != null)
// 				if(currentMusic.isPlaying())
// 					if(currentMusicName == filename)
// 						return;
// 					else
// 						musicStop();
// 			if(musicMap.containsKey(filename)){
// 				Scene.log("Music: playing "+filename);
// 				currentMusic = musicMap.get(filename);//Gdx.audio.newMusic(Gdx.files.internal("music/"+filename));
// 				currentMusic.setVolume(Config.volMusic);
// 				currentMusic.setLooping(true);
// 				currentMusic.play();
// 				currentMusicName = filename;
// 			}
// 			else{
// 				Scene.log("Music File Not Found: "+filename);
// 			}
// 		}
// 	}

// 	/** Pauses the current music file being played */
// 	public static void musicPause(){
// 		if(currentMusic != null)
// 			if(currentMusic.isPlaying()){
// 				Scene.log("Music: pausing "+currentMusicName);
// 				currentMusic.pause();
// 			}
// 	}

// 	/** Resumes the current music file being played */
// 	public static void musicResume(){
// 		if(currentMusic != null)
// 			if(!currentMusic.isPlaying()){
// 				Scene.log("Music: resuming "+currentMusicName);
// 				currentMusic.play();
// 			}
// 		else
// 			musicPlay(currentMusicName);
// 	}

// 	/** Stops the current music file being played */
// 	public static void musicStop(){
// 		if(currentMusic != null){
// 			Scene.log("Music: stoping "+currentMusicName);
// 			currentMusic.stop();
// 			currentMusic = null;
// 		}
// 	}

// 	/** Sets the volume music file being played */
// 	public static void musicVolume(){
// 		if(currentMusic != null);
// 			currentMusic.setVolume(Config.volMusic);
// 	}

// 	/** Disoposes the current music file being played */
// 	public static void musicDispose(){
// 		if(currentMusic != null);
// 			currentMusic.dispose();
// 	}

// /***********************************************************************************************************
// * 								Sound Related Global Functions							   				   *
// ************************************************************************************************************/
// 	/** Plays the sound file which was dynamically loaded if it is present otherwise logs the name
// 	 *  @param filename The Sound File name only
// 	 *  @ex <code>soundPlay("bang")</code>
// 	 *  */
// 	public static void soundPlay(String filename){
// 		if(Config.isSound){
// 			if(soundMap.containsKey(filename)){
// 				currentSound = soundMap.get(filename);
// 				long id = currentSound.play(Config.volSound);
// 				currentSound.setLooping(id, false);
// 				currentSound.setPriority(id, 99);
// 				Scene.log("Sound:"+"Play "+ filename);
// 			}
// 			else{
// 				Scene.log("Music File Not Found: "+filename);
// 			}
// 		}
// 	}

// 	/** Plays the sound file "click" */
// 	public static void soundClick(){
// 		if(Config.isSound && soundMap.containsKey("click")){
// 	        currentSound = soundMap.get("click");
// 			long id = currentSound.play(Config.volSound);
// 			currentSound.setLooping(id, false);
// 			currentSound.setPriority(id, 99);
// 			Scene.log("Sound:"+"Play "+ "click");
// 		}
// 	}

// 	/** Pauses the current sound file being played */
// 	public static void soundPause(){
// 		Scene.log("Sound:"+"Pausing");
// 		if(currentSound != null)
// 			currentSound.pause();
// 	}

// 	/** Resumes the current sound file being played */
// 	public static void soundResume(){
// 		Scene.log("Sound:"+"Resuming");
// 		if(currentSound != null)
// 			currentSound.resume();
// 	}

// 	/** Stops the current sound file being played */
// 	public static void soundStop(){
// 		Scene.log("Sound:"+"Stopping");
// 		if(currentSound != null)
// 			currentSound.stop();
// 	}

// 	/** Disposes the current sound file being played */
// 	public static void soundDispose(){
// 		Scene.log("Sound:"+"Disposing Sound");
// 		if(currentSound != null)
// 			currentSound.dispose();
// 	}

// /***********************************************************************************************************
// * 								BitmapFont Related Functions							   				   *
// ************************************************************************************************************/
// 	/** If key is present returns the BitmapFont that was dynamically loaded
// 	 *  else returns null
// 	 *  @param fontname The BitmapFont name
// 	 *  @return BitmapFont or null
// 	 *  @ex font("font1") or font("arial")
// 	 *  */
// 	public static BitmapFont font(String fontname){
// 		if(fontMap.containsKey(fontname)){
// 			return fontMap.get(fontname);
// 		}
// 		else{
// 			Scene.log("Font File Not Found: "+fontname);
// 			return null;
// 		}
// 	}

// /***********************************************************************************************************
// * 								Texture Related Functions							   				   	   *
// ************************************************************************************************************/
// 	/** If key is present returns the TextureRegion that was loaded from all the atlases
// 	 *  else returns null
// 	 *  @param textureregionName The TextureRegion name
// 	 *  @return TextureRegion or null
// 	 *  */
// 	public static TextureRegion tex(String textureregionName){
// 		if(texMap.containsKey(textureregionName)){
// 			return texMap.get(textureregionName);
// 		}
// 		else{
// 			Scene.log("TextureRegion Not Found: "+textureregionName);
// 			return null;
// 		}
// 	}

// /***********************************************************************************************************
// * 								Animation Related Functions							   				   	   *
// ************************************************************************************************************/
// 	private static Animation anim(String texName, int numberOfFrames, int hOffset) {
// 		// Key frames list
// 		TextureRegion[] keyFrames = new TextureRegion[numberOfFrames];
// 		TextureRegion texture = Asset.tex(texName);
// 		int width = texture.getRegionWidth() / numberOfFrames;
// 		int height = texture.getRegionHeight();
// 		// Set key frames (each comes from the single texture)
// 		for (int i = 0; i < numberOfFrames; i++)
// 			keyFrames[i] = new TextureRegion(texture, width * i, hOffset, width, height);
// 		Animation animation = new Animation(1f/numberOfFrames, keyFrames);
// 		return animation;
// 	}

// 	private static Animation anim(String texName, int numberOfFrames, float duration, int hOffset) {
// 		// Key frames list
// 		TextureRegion[] keyFrames = new TextureRegion[numberOfFrames];
// 		TextureRegion texture = Asset.tex(texName);
// 		int width = texture.getRegionWidth() / numberOfFrames;
// 		int height = texture.getRegionHeight();
// 		// Set key frames (each comes from the single texture)
// 		for (int i = 0; i < numberOfFrames; i++)
// 			keyFrames[i] = new TextureRegion(texture, width * i, hOffset, width, height);
// 		Animation animation = new Animation(duration, keyFrames);
// 		return animation;
// 	}

// 	/**
// 	 * Get animation from single textureRegion which contains all frames
// 	 * (It is like a single png which has all the frames). Each frames' width should be same.
// 	 * <p>
// 	 * @param texName
// 	 * 			  the name of the texture region
// 	 * @param numberOfFrames
// 	 *            number of frames of the texture Region
// 	 * @return animation created
// 	 *
// 	 * */
// 	public static Animation anim(String texName, int numberOfFrames) {
// 		return anim(texName, numberOfFrames, 0);
// 	}

// 	/**
// 	 * Get animation from single textureRegion which contains all frames
// 	 * (It is like a single png which has all the frames). Each frames' width should be same.
// 	 * <p>
// 	 * @param texName
// 	 * 			  the name of the texture region
// 	 * @param numberOfFrames
// 	 *            number of frames of the texture Region
// 	 * @param duration
// 	 *            each frame duration on play
// 	 * @return animation created
// 	 *
// 	 * */
// 	public static Animation anim(String texName, int numberOfFrames, float duration) {
// 		return anim(texName, numberOfFrames, duration, 0);
// 	}

// 	/**
// 	 * There is only single texture region which contains all frames
// 	 * (It is like a single png which has all the frames).
// 	 * The texture regions consists of both rows and columns
// 	 * <p>
// 	 *
// 	 * @param textureAtlas
// 	 *            atlas which contains the single animation texture
// 	 * @param animationName
// 	 *            name of the animation in atlas
// 	 * @param numberOfFrames
// 	 *            number of frames of the animation
// 	 * @param numberOfMaximumFramesInTheSheet
// 	 *            maximum number of frame in a row in the sheet
// 	 * @param numberOfRows
// 	 *            number of rows that the sheet contains
// 	 * @param indexOfAnimationRow
// 	 *            the row index (starts from 0) that desired animation exists
// 	 * @param frameDuration
// 	 *            each frame duration on play
// 	 * @return animation created
// 	 *
// 	 * */
// 	public static Animation[] anim(String texName, int rows, int cols, float duration) {
// 		TextureRegion texture = Asset.tex(texName);
// 		// Set key frames (each comes from the single texture)
// 		/*for (int i = 0; i < numberOfFrames; i++) {
// 			keyFrames[i] = new TextureRegion(
// 					textureRegion,
// 					(textureRegion.getRegionWidth() / numberOfMaximumFramesInTheSheet)
// 							* i, textureRegion.getRegionHeight() / numberOfRows
// 							* indexOfAnimationRow,
// 					textureRegion.getRegionWidth()
// 							/ numberOfMaximumFramesInTheSheet,
// 					textureRegion.getRegionHeight() / numberOfRows);
// 		}*/
// 		int height = texture.getRegionHeight()/rows;
// 		Animation[] animations = new Animation[rows];
// 		for(int i=0;i<rows;i++)
// 			animations[i] = anim(texName, cols, i*height);
// 		return animations;
// 	}

// /***********************************************************************************************************
// 	* 								TMX MAP Related Functions							   				   *
// ************************************************************************************************************/
// 	/*
// 	 * Loads a Tmx map by specifying the map/level no
// 	 * eg: loadTmx(4) -> returns the TiledMap "map/level4.tmx"
// 	 *
// 	 * Note: Tmx Maps must be loaded and unloaded manually as they may take a lot of time to load
// 	 */
// 	public static TiledMap map(int i){
// 		assetMan.setLoader(TiledMap.class, new TmxMapLoader(new InternalFileHandleResolver()));
// 		assetMan.load(basePath+"map/level"+i+".tmx", TiledMap.class);
// 		assetMan.finishLoading();
// 		return assetMan.get(basePath+"map/level"+i+".tmx", TiledMap.class);
// 	}

// 	/*
// 	 * unloads a Tmx map by specifying the map/level no
// 	 * eg: unloadTmx(4) -> unloads the TiledMap "map/level4.tmx"
// 	 *
// 	 * Note: Tmx Maps must be unloaded manually
// 	 */
// 	public static void unloadmap(int i){
// 		assetMan.unload(basePath+"map/level"+i+".tmx");
// 	}

// 	/*
// 	 * Load a G3db model from the model directory
// 	 * @param modelName The name of the modelFile
// 	 * @example loadModel("ship");
// 	 */
// 	public static Actor3d loadModel(String modelName){
// 		assetMan.load(basePath+"model/"+modelName+".g3db", Model.class);
// 		assetMan.finishLoading();
// 		return new Actor3d(assetMan.get(basePath+"model/"+modelName+".g3db", Model.class));
// 	}

// 	/*
// 	 * Unload a G3db model which was previously loaded in the assetManager
// 	 * @param modelName The name of the modelFile that was loaded
// 	 * @example unloadModel("ship");
// 	 */
// 	public static void unloadModel(String modelName){
// 		assetMan.unload(basePath+"model/"+modelName+".g3db");
// 	}

// 	/*
// 	 * Load a Obj model from the model directory
// 	 * @param modelName The name of the modelFile
// 	 * @example loadModel("ship");
// 	 */
// 	public static Actor3d loadModelObj(String modelName){
// 		assetMan.load(basePath+"model/"+modelName+".obj", Model.class);
// 		assetMan.finishLoading();
// 		return new Actor3d(assetMan.get(basePath+"model/"+modelName+".obj", Model.class));
// 	}

// 	/*
// 	 * Unload a Obj model which was previously loaded in the assetManager
// 	 * @param modelName The name of the modelFile that was loaded
// 	 * @example unloadModel("ship");
// 	 */
// 	public static void unloadModelObj(String modelName){
// 		assetMan.unload(basePath+"model/"+modelName+".obj");
// 	}
// 	/*