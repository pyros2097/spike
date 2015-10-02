// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

/** <p>
 * A Preference instance is a hash map holding different values. It is stored alongside your application (SharedPreferences on
 * Android, LocalStorage on GWT, on the desktop a Java Preferences file in a ".prefs" directory will be created, and on iOS an
 * NSMutableDictionary will be written to the given file). CAUTION: On the desktop platform, all libgdx applications share the same
 * ".prefs" directory. To avoid collisions use specific names like "com.myname.game1.settings" instead of "settings"
 * </p>
 *
 * <p>
 * Changes to a preferences instance will be cached in memory until {@link #flush()} is invoked.
 * </p>
 *
 * <p>
 * Use {@link Application#getPreferences(String)} to look up a specific preferences instance. Note that on several backends the
 * preferences name will be used as the filename, so make sure the name is valid for a filename.
 * </p>
 */
/** Returns the {@link Preferences} instance of this Application. It can be used to store application settings across runs.
 * @param name the name of the preferences, must be useable as a file name.
 * @return the preferences. */
type Preferences interface {
	PutBoolean(key string, val bool) Preferences

	PutInteger(key string, val int) Preferences

	PutLong(key string, val int64) Preferences

	PutFloat(key string, val float32) Preferences

	PutString(key, val string) Preferences

	Put(vals map[string]interface{}) Preferences

	GetBoolean(key string, defValue bool) bool

	GetInteger(key string, defValue int) int

	GetLong(key string, defValue int64) int64

	GetFloat(key string, defValue float32) float32

	GetString(key, defValue string) string

	// Returns a read only map[string]interface{} with all the key, objects of the preferences.
	Get() map[string]interface{}

	Contains(key string) bool

	Clear()

	Remove(key string)

	// Makes sure the preferences are persisted.
	Flush()
}

/** The Configuration/Settings for the Game
 * <p>
 * The Config class contains all the necessary options for all different platforms into one class.
 * Here you can save all the data of the game that is required to be persistent.<br>
 */

const (
	MUSIC            = "music"
	SOUND            = "sound"
	VOLUME_MUSIC     = "volumeMusic"
	VOLUME_SOUND     = "volumeSOUND"
	VIBRATION        = "vibration"
	BATTLE_ANIMATION = "battleanimation"
	SEMI_AUTOMATIC   = "semiautomatic"
	FIRST_LAUNCH     = "firstLaunch"
	LEVELS           = "levels"
	CURRENT_LEVEL    = "currentlevel"
	SAVEDATA         = "savedata"
	TOTAL_TIME       = "totaltime"
	PANSPEED         = "panspeed"
	DRAGSPEED        = "dragspeed"
	KEYBOARD         = "keyboard"
	SCORE            = "score"
)

var (
	prefs        Preferences
	HasMusic     bool
	HasSound     bool
	UseKeyboard  bool
	HasVibration bool
	VolMusic     float32
	VolSound     float32
	SpeedPan     float32
	SpeedDrag    float32
	Score        int
)

func initConfig(configName string) {
	println("Initializing Config")
	// prefs = Gdx.app.getPreferences(Scene.configJson.getString("title"))
	// HasMusic = prefs.GetBoolean(MUSIC, true)
	// HasSound = prefs.GetBoolean(SOUND, true)
	// HasSound = prefs.GetBoolean(VIBRATION, false)

	// VolMusic = prefs.GetFloat(VOLUME_MUSIC, 1)
	// VolSound = prefs.GetFloat(VOLUME_SOUND, 1)

	// UseKeyboard = prefs.GetBoolean(KEYBOARD, true)

	// SpeedPan = prefs.GetFloat(PANSPEED, 5)
	// SpeedDrag = prefs.GetFloat(DRAGSPEED, 5)

	// Score = prefs.GetInteger(SCORE, 0)
}

func LoadSaveData() string {
	return prefs.GetString(SAVEDATA, "")
}

func WriteSaveData(data string) {
	prefs.PutString(SAVEDATA, data)
	prefs.Flush()
}

func ReadTotalTime() float32 {
	return prefs.GetFloat(TOTAL_TIME, 0.0)
}

func WriteTotalTime(secondstime float32) {
	prefs.PutFloat(TOTAL_TIME, secondstime+ReadTotalTime())
	prefs.Flush()
}

//     public static int levels(){
//         return prefs.getInteger(LEVELS, 20);
//     }

//     func setLevels(int ue){
//         prefs.putInteger(LEVELS, ue);
//         prefs.flush();
//     }

//     public static int currentLevel(){
//         return prefs.getInteger(CURRENT_LEVEL, 0);
//     }

//     func setCurrentLevel(int ue){
//         prefs.putInteger(CURRENT_LEVEL, ue);
//         prefs.flush();
//     }

//     public static boolean isBattleEnabled(){
//         return prefs.getBoolean(BATTLE_ANIMATION, true);
//     }

//     func setBattle(boolean ue){
//         prefs.putBoolean(BATTLE_ANIMATION, ue);
//         prefs.flush();
//     }

//     public static boolean isSemiAutomatic(){
//         return prefs.getBoolean(SEMI_AUTOMATIC, false);
//     }

//     func setSemiAutomatic(boolean ue){
//         prefs.putBoolean(SEMI_AUTOMATIC, ue);
//         prefs.flush();
//     }

//     func setPanSpeed(float ue){
//         prefs.putFloat(PANSPEED, ue);
//         prefs.flush();
//         speedPan = ue;
//     }

//     func setDragSpeed(float ue){
//         prefs.putFloat(VOLUME_SOUND, ue);
//         prefs.flush();
//         speedDrag = ue;
//     }

func EnableKeyboard(enable bool) {
	prefs.PutBoolean(KEYBOARD, enable)
	prefs.Flush()
	UseKeyboard = enable
}

func EnableSound(enable bool) {
	prefs.PutBoolean(SOUND, enable)
	prefs.Flush()
	HasSound = enable

}

func EnableMusic(enable bool) {
	prefs.PutBoolean(MUSIC, enable)
	prefs.Flush()
	HasMusic = enable
}

func SetMusicVolume(volume float32) {
	prefs.PutFloat(VOLUME_MUSIC, volume)
	prefs.Flush()
	VolMusic = volume
	// Asset.musicVolume()
}

func SetSoundVolume(volume float32) {
	prefs.PutFloat(VOLUME_SOUND, volume)
	prefs.Flush()
	VolSound = volume
	// Asset.musicVolume()
}

func setScore(value int) {
	prefs.PutInteger(SCORE, value)
	prefs.Flush()
	Score = value
}

func setVibration(enable bool) {
	prefs.PutBoolean(VIBRATION, enable)
	prefs.Flush()
	HasVibration = enable
}

func setFirstLaunchDone() {
	prefs.PutBoolean(FIRST_LAUNCH, false)
	prefs.Flush()

}

func isFirstLaunch() bool {
	return prefs.GetBoolean(FIRST_LAUNCH, true)
}
