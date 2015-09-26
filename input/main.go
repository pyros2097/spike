// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Interface to the input facilities. This allows polling the state of the keyboard, the touch screen and the accelerometer.
// On some backends (desktop) the touch screen is replaced by mouse input. The accelerometer is of course not available on
// all backends.
//
// Keyboard keys are translated to the constants in {@link Keys} transparently on all systems. Do not use system specific key
// constants.
//
// This also offers methods to use (and test for the presence of) other input systems like vibration, compass, on-screen
// keyboards, and cursor capture. Support for simple input dialogs is also provided.
package input

type Orientation int

const (
	Landscape Orientation = 0
	Portrait  Orientation = 1
)

// Callback interface for {@link Input#getTextInput(TextInputListener, String, String, String)}
func TextInputListener(text string) {}
func TextInputCanceled()            {}

// Enumeration of potentially available peripherals. Use with {@link Input#isPeripheralAvailable(Peripheral)}.
// public enum Peripheral {
//   HardwareKeyboard, OnscreenKeyboard, MultitouchScreen, Accelerometer, Compass, Vibrator
// }

// return The value of the accelerometer on its x-axis. ranges between [-10,10].
func GetAccelerometerX() float32 {
}

// return The value of the accelerometer on its y-axis. ranges between [-10,10].
func GetAccelerometerY() float32 {
}

// return The value of the accelerometer on its y-axis. ranges between [-10,10].
func GetAccelerometerZ() float32 {
}

// return The x coordinate of the last touch on touch screen devices and the current mouse position on desktop for the first
// pointer in screen coordinates. The screen origin is the top left corner.
func GetX() int {
}

// Returns the x coordinate in screen coordinates of the given pointer. Pointers are indexed from 0 to n. The pointer id
// identifies the order in which the fingers went down on the screen, e.g. 0 is the first finger, 1 is the second and so on.
// When two fingers are touched down and the first one is lifted the second one keeps its index. If another finger is placed on
// the touch screen the first free index will be used.
//
// param pointer the pointer id.
// return the x coordinate
func GetXPointer(pointer int) int {
}

// return the different between the current pointer location and the last pointer location on the x-axis.
func GetDeltaX() int {
}

// return the different between the current pointer location and the last pointer location on the x-axis.
func GetDeltaXPointer(pointer int) int {
}

// return The y coordinate of the last touch on touch screen devices and the current mouse position on desktop for the first
// pointer in screen coordinates. The screen origin is the top left corner.
func GetY() int {
}

// Returns the y coordinate in screen coordinates of the given pointer. Pointers are indexed from 0 to n. The pointer id
// identifies the order in which the fingers went down on the screen, e.g. 0 is the first finger, 1 is the second and so on.
// When two fingers are touched down and the first one is lifted the second one keeps its index. If another finger is placed on
// the touch screen the first free index will be used.
//
// param pointer the pointer id.
// return the y coordinate
func GetYPointer(pointer int) int {
}

// return the different between the current pointer location and the last pointer location on the y-axis.
func GetDeltaY() int {
}

// return the different between the current pointer location and the last pointer location on the y-axis.
func GetDeltaYPointer(pointer int) int {
}

// return whether the screen is currently touched.
func IsTouched() bool {
}

// return whether a new touch down event just occurred.
func JustTouched() bool {
}

// Whether the screen is currently touched by the pointer with the given index. Pointers are indexed from 0 to n. The pointer
// id identifies the order in which the fingers went down on the screen, e.g. 0 is the first finger, 1 is the second and so on.
// When two fingers are touched down and the first one is lifted the second one keeps its index. If another finger is placed on
// the touch screen the first free index will be used.
//
// param pointer the pointer
// return whether the screen is touched by the pointer
func IsTouchedPointer(pointer int) bool {
}

// Whether a given button is pressed or not. Button constants can be found in {@link Buttons}. On Android only the Button#LEFT
// constant is meaningful before version 4.0.
// param button the button to check.
// return whether the button is down or not.
func IsButtonPressed(button int) bool {
}

// Returns whether the key is pressed.
//
// param key The key code as found in {@link Input.Keys}.
// return true or false.
func IsKeyPressed(key int) bool {
}

// Returns whether the key has just been pressed.
//
// param key The key code as found in {@link Input.Keys}.
// return true or false.
func IsKeyJustPressed(key int) bool {
}

// System dependent method to input a string of text. A dialog box will be created with the given title and the given text as a
// message for the user. Once the dialog has been closed the provided {@link TextInputListener} will be called on the rendering
// thread.
//
// param listener The TextInputListener.
// param title The title of the text input dialog.
// param text The message presented to the user.
// func GetTextInput(TextInputListener listener, title, text, hint string)

// Sets the on-screen keyboard visible if available.
//
// param visible visible or not
func SetOnscreenKeyboardVisible(visible bool) {
}

// Vibrates for the given amount of time. Note that you'll need the permission
// <code> <uses-permission android:name="android.permission.VIBRATE" /></code> in your manifest file in order for this to work.
//
// param milliseconds the number of milliseconds to vibrate.
func Vibrate(milliseconds int) {
}

// Vibrate with a given pattern. Pass in an array of ints that are the times at which to turn on or off the vibrator. The first
// one is how long to wait before turning it on, and then after that it alternates. If you want to repeat, pass the index into
// the pattern at which to start the repeat.
// param pattern an array of longs of times to turn the vibrator on or off.
// param repeat the index into pattern at which to repeat, or -1 if you don't want to repeat.
func VibrateRepeat(pattern []int64, repeat int) {
}

// Stops the vibrator
func CancelVibrate() {
}

// The azimuth is the angle of the device's orientation around the z-axis. The positive z-axis points towards the earths
// center.
//
// @see <a
//      href="http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])">http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])</a>
// return the azimuth in degrees
func GetAzimuth() float32 {
}

// The pitch is the angle of the device's orientation around the x-axis. The positive x-axis roughly points to the west and is
// orthogonal to the z- and y-axis.
// @see <a
//      href="http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])">http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])</a>
// return the pitch in degrees
func GetPitch() float32 {
}

// The roll is the angle of the device's orientation around the y-axis. The positive y-axis points to the magnetic north pole
// of the earth.
// @see <a
//      href="http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])">http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])</a>
// return the roll in degrees
func GetRoll() float32 {
}

// Returns the rotation matrix describing the devices rotation as per <a href=
// "http://developer.android.com/reference/android/hardware/SensorManager.html#getRotationMatrix(float[], float[], float[], float[])"
// >SensorManager#getRotationMatrix(float[], float[], float[], float[])</a>. Does not manipulate the matrix if the platform
// does not have an accelerometer.
// param matrix
func GetRotationMatrix(matrix []float32) {
}

// return the time of the event currently reported to the {@link InputProcessor}.
func GetCurrentEventTime() int64 {
}

// Sets whether the BACK button on Android should be caught. This will prevent the app from being paused. Will have no effect
// on the desktop.
//
// param catchBack whether to catch the back button
func SetCatchBackKey(catchBack bool) {
}

// return whether the back button is currently being caught
func IsCatchBackKey() bool {
}

// Sets whether the MENU button on Android should be caught. This will prevent the onscreen keyboard to show up. Will have no
// effect on the desktop.
//
// param catchMenu whether to catch the menu button
func SetCatchMenuKey(catchMenu bool) {
}

// Sets the {@link InputProcessor} that will receive all touch and key input events. It will be called before the
// {@link ApplicationListener#render()} method each frame.
//
// param processor the InputProcessor
// func SetInputProcessor(InputProcessor processor);

// return the currently set {@link InputProcessor} or null.
// public InputProcessor getInputProcessor ();

// Queries whether a {@link Peripheral} is currently available. In case of Android and the {@link Peripheral#HardwareKeyboard}
// this returns the whether the keyboard is currently slid out or not.
//
// param peripheral the {@link Peripheral}
// return whether the peripheral is available or not.
func IsPeripheralAvailable(peripheral Peripheral) bool {
}

// return the rotation of the device with respect to its native orientation.
func GetRotation() int {
}

// return the native orientation of the device.
func GetNativeOrientation() Orientation {
}

// Only viable on the desktop. Will confine the mouse cursor location to the window and hide the mouse cursor. X and y
// coordinates are still reported as if the mouse was not catched.
// param catched whether to catch or not to catch the mouse cursor
func SetCursorCatched(catched bool) {
}

// return whether the mouse cursor is catched.
func IsCursorCatched() bool {
}

// Only viable on the desktop. Will set the mouse cursor location to the given window coordinates (origin top-left corner).
// param x the x-position
// param y the y-position
func SetCursorPosition(x, y int) {
}
