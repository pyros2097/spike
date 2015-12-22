// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package spike

import (
	"fmt"
)

/** A color class, holding the r, g, b and alpha component as floats in the range [0,1]. All methods perform clamping on the
 * internal values after execution.
 */
var (
	CLEAR = NewColor(0, 0, 0, 0)
	BLACK = NewColor(0, 0, 0, 1)

	WHITE      = NewColorHex(0xffffffff)
	LIGHT_GRAY = NewColorHex(0xbfbfbfff)
	GRAY       = NewColorHex(0x7f7f7fff)
	DARK_GRAY  = NewColorHex(0x3f3f3fff)
	SLATE      = NewColorHex(0x708090ff)

	BLUE  = NewColor(0, 0, 1, 1)
	NAVY  = NewColor(0, 0, 0.5, 1)
	ROYAL = NewColorHex(0x4169e1ff)
	SKY   = NewColorHex(0x87ceebff)
	CYAN  = NewColor(0, 1, 1, 1)
	TEAL  = NewColor(0, 0.5, 0.5, 1)

	GREEN      = NewColorHex(0x00ff00ff)
	CHARTREUSE = NewColorHex(0x7fff00ff)
	LIME       = NewColorHex(0x32cd32ff)
	FOREST     = NewColorHex(0x228b22ff)
	OLIVE      = NewColorHex(0x6b8e23ff)

	YELLOW    = NewColorHex(0xffff00ff)
	GOLD      = NewColorHex(0xffd700ff)
	GOLDENROD = NewColorHex(0xdaa520ff)

	BROWN     = NewColorHex(0x8b4513ff)
	TAN       = NewColorHex(0xd2b48cff)
	FIREBRICK = NewColorHex(0xb22222ff)

	RED     = NewColorHex(0xff0000ff)
	CORAL   = NewColorHex(0xff7f50ff)
	ORANGE  = NewColorHex(0xffa500ff)
	SALMON  = NewColorHex(0xfa8072ff)
	PINK    = NewColorHex(0xff69b4ff)
	MAGENTA = NewColor(1, 0, 1, 1)

	PURPLE = NewColorHex(0xa020f0ff)
	VIOLET = NewColorHex(0xee82eeff)
	MAROON = NewColorHex(0xb03060ff)
)

type Color struct {
	// the red, green, blue and alpha components
	R, G, B, A float32
}

// Constructor, sets the components of the color
func NewColor(r, g, b, a float32) *Color {
	color := &Color{r, g, b, a}
	color.Clamp()
	return color
}

// Constructs a new Color with all components set to 0
func NewColorEmpty() *Color {
	return &Color{}
}

//@see #rgba8888ToColor(Color, int)
func NewColorHex(rgba8888 uint32) *Color {
	color := &Color{}
	RGBA8888ToColor(color, rgba8888)
	return color
}

func NewColorCopy(color *Color) *Color {
	return NewColor(color.R, color.G, color.B, color.A)
}

// Sets this color to the given color
func (self *Color) SetColor(color *Color) {
	self.R = color.R
	self.G = color.G
	self.B = color.B
	self.A = color.A
}

// Sets this Color's component values.
func (self *Color) Set(r, g, b, a float32) *Color {
	self.R = r
	self.G = g
	self.B = b
	self.A = a
	return self.Clamp()
}

// Sets this color's component values through an integer representation.
func (self *Color) SetHex(rgba uint32) *Color {
	RGBA8888ToColor(self, rgba)
	return self
}

// Multiplies the this color and the given color
func (self *Color) MulColor(color *Color) *Color {
	self.R *= color.R
	self.G *= color.G
	self.B *= color.B
	self.A *= color.A
	return self.Clamp()
}

// Multiplies all components of this Color with the given value
func (self *Color) MulValue(value float32) *Color {
	self.R *= value
	self.G *= value
	self.B *= value
	self.A *= value
	return self.Clamp()
}

// Multiplies this Color's color components by the given ones
func (self *Color) Mul(r, g, b, a float32) *Color {
	self.R *= r
	self.G *= g
	self.B *= b
	self.A *= a
	return self.Clamp()
}

// Adds the given color to this color
func (self *Color) AddColor(color *Color) *Color {
	self.R += color.R
	self.G += color.G
	self.B += color.B
	self.A += color.A
	return self.Clamp()
}

// Adds the given color component values to this Color's values.
func (self *Color) Add(r, g, b, a float32) *Color {
	self.R += r
	self.G += g
	self.B += b
	self.A += a
	return self.Clamp()
}

// Subtracts the given color from this color
func (self *Color) SubColor(color *Color) *Color {
	self.R -= color.R
	self.G -= color.G
	self.B -= color.B
	self.A -= color.A
	return self.Clamp()
}

// Subtracts the given values from this Color's component values
func (self *Color) Sub(r, g, b, a float32) *Color {
	self.R -= r
	self.G -= g
	self.B -= b
	self.A -= a
	return self.Clamp()
}

// Clamps this Color's components to a valid range [0 - 1]
func (self *Color) Clamp() *Color {
	switch {
	case self.R < 0:
		self.R = 0
	case self.R > 1:
		self.R = 1
	case self.G < 0:
		self.G = 0
	case self.G > 1:
		self.G = 1
	case self.B < 0:
		self.B = 0
	case self.B > 1:
		self.B = 1
	case self.A < 0:
		self.A = 0
	case self.A > 1:
		self.A = 1
	}
	return self
}

/** Linearly interpolates between this color and the target color by t which is in the range [0,1]. The result is stored in this
 * color.
 * @param target The target color
 * @param t The interpolation coefficient
 * @return This color for chaining. */
func (self *Color) LerpColor(target *Color, t float32) *Color {
	self.R += t * (target.R - self.R)
	self.G += t * (target.G - self.G)
	self.B += t * (target.B - self.B)
	self.A += t * (target.A - self.A)
	return self.Clamp()
}

/** Linearly interpolates between this color and the target color by t which is in the range [0,1]. The result is stored in this
 * color.
 * @param r The red component of the target color
 * @param g The green component of the target color
 * @param b The blue component of the target color
 * @param a The alpha component of the target color
 * @param t The interpolation coefficient
 * @return This color for chaining. */
func (self *Color) Lerp(r, g, b, a, t float32) *Color {
	self.R += t * (r - self.R)
	self.G += t * (g - self.G)
	self.B += t * (b - self.B)
	self.A += t * (a - self.A)
	return self.Clamp()
}

// Multiplies the RGB values by the alpha.
func (self *Color) PremultiplyAlpha() *Color {
	self.R *= self.A
	self.G *= self.A
	self.B *= self.A
	return self
}

func (self *Color) Equals(other *Color) bool {
	if self == other {
		return true
	}
	if other == nil {
		return false
	}
	return self.ToIntBits() == other.ToIntBits()
}

func (self *Color) HashCode() uint32 {
	// int result = (r != +0.0f ? NumberUtils.FloatToIntBits(r) : 0);
	// result = 31 * result + (g != +0.0f ? NumberUtils.FloatToIntBits(g) : 0);
	// result = 31 * result + (b != +0.0f ? NumberUtils.FloatToIntBits(b) : 0);
	// result = 31 * result + (a != +0.0f ? NumberUtils.FloatToIntBits(a) : 0);
	// return result;
	return 0
}

/** Packs the color components into a 32-bit integer with the format ABGR and then converts it to a float.
 * @return the packed color as a 32-bit float
 * @see NumberUtils#intToFloatColoruint32 */
func (self *Color) ToFloatBits() float32 {
	return 0
	// return NumberUtils.IntToFloatColor((uint32(255*a) << 24) | (uint32(255*b) << 16) | (uint32(255*g) << 8) | (uint32(255 * r)))
}

/** Packs the color components into a 32-bit integer with the format ABGR.
 * @return the packed color as a 32-bit int. */
func (self *Color) ToIntBits() uint32 {
	return (uint32(255*self.A) << 24) | (uint32(255*self.B) << 16) | (uint32(255*self.G) << 8) | (uint32(255 * self.R))
}

// Returns the color encoded as hex string with the format RRGGBBAA.
func (self *Color) ToString() string {
	value := (uint32(255*self.R) << 24) | (uint32(255*self.G) << 16) | (uint32(255*self.B) << 8) | (uint32(255 * self.A))
	return fmt.Sprintf("%x", value)
}

// @return a copy of this color
func (self *Color) Copy() *Color {
	return NewColorCopy(self)
}

// TODO: Must be static
/** Returns a new color from a hex string with the format RRGGBBAA.
 * @see #toString() */
func ValueOf(hex string) {
	// int r = Integer.valueOf(hex.substring(0, 2), 16);
	// int g = Integer.valueOf(hex.substring(2, 4), 16);
	// int b = Integer.valueOf(hex.substring(4, 6), 16);
	// int a = hex.length() != 8 ? 255 : Integer.valueOf(hex.substring(6, 8), 16);
	// return NewColor(r / 255f, g / 255f, b / 255f, a / 255f);
}

/** Packs the color components into a 32-bit integer with the format ABGR and then converts it to a float. Note that no range
 * checking is performed for higher performance.
 * @param r the red component, 0 - 255
 * @param g the green component, 0 - 255
 * @param b the blue component, 0 - 255
 * @param a the alpha component, 0 - 255
 * @return the packed color as a float
 * @see NumberUtils#intToFloatColoruint32 */
func ToFloatBits(r, g, b, a int) float32 {
	return 0
	// return NumberUtils.IntToFloatColor((a << 24) | (b << 16) | (g << 8) | r)
}

/** Packs the color components into a 32-bit integer with the format ABGR and then converts it to a float.
 * @return the packed color as a 32-bit float
 * @see NumberUtils#intToFloatColoruint32 */
func ToFloatBitsF(r, g, b, a float32) float32 {
	return 0
	// return NumberUtils.IntToFloatColor((uint32(255*a) << 24) | (uint32(255*b) << 16) | (uint32(255*g) << 8) | (uint32(255 * r)))
}

/** Packs the color components into a 32-bit integer with the format ABGR. Note that no range checking is performed for higher
 * performance.
 * @param r the red component, 0 - 255
 * @param g the green component, 0 - 255
 * @param b the blue component, 0 - 255
 * @param a the alpha component, 0 - 255
 * @return the packed color as a 32-bit int */
func ToIntBits(r, g, b, a uint32) uint32 {
	return (a << 24) | (b << 16) | (g << 8) | r
}

func Alpha(alpha float32) uint32 {
	return uint32(alpha * 255.0)
}

func LuminanceAlpha(luminance, alpha float32) uint32 {
	return (uint32(luminance*255.0) << 8) | uint32(alpha*255)
}

func RGB565(r, g, b float32) uint32 {
	return (uint32(r*31) << 11) | (uint32(g*63) << 5) | uint32(b*31)
}

func RGBA4444(r, g, b, a float32) uint32 {
	return (uint32(r*15) << 12) | (uint32(g*15) << 8) | (uint32(b*15) << 4) | uint32(a*15)
}

func RGB888(r, g, b float32) uint32 {
	return (uint32(r*255) << 16) | (uint32(g*255) << 8) | uint32(b*255)
}

func RGBA8888(r, g, b, a float32) uint32 {
	return (uint32(r*255) << 24) | (uint32(g*255) << 16) | (uint32(b*255) << 8) | uint32(a*255)
}

func ARGB8888(a, r, g, b float32) uint32 {
	return (uint32(a*255) << 24) | (uint32(r*255) << 16) | (uint32(g*255) << 8) | uint32(b*255)
}

func RGB565Color(color *Color) uint32 {
	return (uint32(color.R*31) << 11) | (uint32(color.G*63) << 5) | uint32(color.B*31)
}

func RGBA4444Color(color *Color) uint32 {
	return (uint32(color.R*15) << 12) | (uint32(color.G*15) << 8) | (uint32(color.B*15) << 4) | uint32(color.A*15)
}

func RGB888Colot(color *Color) uint32 {
	return (uint32(color.R*255) << 16) | (uint32(color.G*255) << 8) | uint32(color.B*255)
}

func RGBA8888Color(color *Color) uint32 {
	return (uint32(color.R*255) << 24) | (uint32(color.G*255) << 16) | (uint32(color.B*255) << 8) | uint32(color.A*255)
}

func ARGB8888Color(color *Color) uint32 {
	return (uint32(color.A*255) << 24) | (uint32(color.R*255) << 16) | (uint32(color.G*255) << 8) | uint32(color.B*255)
}

// TODO: Check >>> 3riple shift

/** Sets the Color components using the specified integer value in the format RGB565. This is inverse to the RGB565(r, g, b)
 * method.
 *
 * @param color The Color to be modified.
 * @param value An integer color value in RGB565 format. */
func RGB565ToColor(color *Color, value uint32) {
	color.R = float32(((value & 0x0000F800) >> 11) / 31)
	color.G = float32(((value & 0x000007E0) >> 5) / 63)
	color.B = float32(((value & 0x0000001F) >> 0) / 31)
}

/** Sets the Color components using the specified integer value in the format RGBA4444. This is inverse to the RGBa4444(r, g, b,
 * a) method.
 *
 * @param color The Color to be modified.
 * @param value An integer color value in RGBA4444 format. */
func RGBA4444ToColor(color *Color, value uint32) {
	color.R = float32(((value & 0x0000f000) >> 12) / 15)
	color.G = float32(((value & 0x00000f00) >> 8) / 15)
	color.B = float32(((value & 0x000000f0) >> 4) / 15)
	color.A = float32((value & 0x0000000f) / 15)
}

/** Sets the Color components using the specified integer value in the format RGB888. This is inverse to the RGB888(r, g, b)
 * method.
 *
 * @param color The Color to be modified.
 * @param value An integer color value in RGB888 format. */
func RGB888ToColor(color *Color, value uint32) {
	color.R = float32(((value & 0x00ff0000) >> 16) / 255)
	color.G = float32(((value & 0x0000ff00) >> 8) / 255)
	color.B = float32((value & 0x000000ff) / 255)
}

/** Sets the Color components using the specified integer value in the format RGBA8888. This is inverse to the RGBa8888(r, g, b,
 * a) method.
 *
 * @param color The Color to be modified.
 * @param value An integer color value in RGBA8888 format. */
func RGBA8888ToColor(color *Color, value uint32) {
	color.R = float32(((value & 0xff000000) >> 24) / 255)
	color.G = float32(((value & 0x00ff0000) >> 16) / 255)
	color.B = float32(((value & 0x0000ff00) >> 8) / 255)
	color.A = float32((value & 0x000000ff) / 255)
}

/** Sets the Color components using the specified integer value in the format ARGB8888. This is the inverse to the aRGB8888(a,
 * r, g, b) method
 *
 * @param color The Color to be modified.
 * @param value An integer color value in ARGB8888 format. */
func ARGB8888ToColor(color *Color, value uint32) {
	color.A = float32(((value & 0xff000000) >> 24) / 255)
	color.R = float32(((value & 0x00ff0000) >> 16) / 255)
	color.G = float32(((value & 0x0000ff00) >> 8) / 255)
	color.B = float32((value & 0x000000ff) / 255)
}

/** A general purpose class containing named colors that can be changed at will. For example, the markup language defined by the
 * {@code BitmapFontCache} class uses this class to retrieve colors and the user can define his own colors.
 */
var colorsMap map[string]*Color

func init() {
	ResetColor()
}

// Returns the color map.
func GetColors() map[string]*Color {
	return colorsMap
}

/** Convenience method to lookup a color by {@code name}. The invocation of this method is equivalent to the expression
 * {@code Colors.getColors().get(name)}
 *
 * @param name the name of the color
 * @return the color to which the specified {@code name} is mapped, or {@code null} if there was no mapping for {@code name}. */
func GetColor(name string) *Color {
	return colorsMap[name]
}

/** Convenience method to add a {@code color} with its {@code name}. The invocation of this method is equivalent to the
 * expression {@code Colors.getColors().put(name, color)}
 *
 * @param name the name of the color
 * @param color the color
 * @return the previous {@code color} associated with {@code name}, or {@code null} if there was no mapping for {@code name}. */
func PutColor(name string, color *Color) {
	colorsMap[name] = color
}

// Resets the color map to the predefined colors.
func ResetColor() {
	colorsMap = make(map[string]*Color)
	colorsMap["CLEAR"] = CLEAR
	colorsMap["BLACK"] = BLACK

	colorsMap["WHITE"] = WHITE
	colorsMap["LIGHT_GRAY"] = LIGHT_GRAY
	colorsMap["GRAY"] = GRAY
	colorsMap["DARK_GRAY"] = DARK_GRAY
	colorsMap["SLATE"] = SLATE

	colorsMap["BLUE"] = BLUE
	colorsMap["NAVY"] = NAVY
	colorsMap["ROYAL"] = ROYAL
	colorsMap["SKY"] = SKY
	colorsMap["CYAN"] = CYAN
	colorsMap["TEAL"] = TEAL

	colorsMap["GREEN"] = GREEN
	colorsMap["CHARTREUSE"] = CHARTREUSE
	colorsMap["LIME"] = LIME
	colorsMap["FOREST"] = FOREST
	colorsMap["OLIVE"] = OLIVE

	colorsMap["YELLOW"] = YELLOW
	colorsMap["GOLD"] = GOLD
	colorsMap["GOLDENROD"] = GOLDENROD

	colorsMap["BROWN"] = BROWN
	colorsMap["TAN"] = TAN
	colorsMap["FIREBRICK"] = FIREBRICK

	colorsMap["RED"] = RED
	colorsMap["CORAL"] = CORAL
	colorsMap["ORANGE"] = ORANGE
	colorsMap["SALMON"] = SALMON
	colorsMap["PINK"] = PINK
	colorsMap["MAGENTA"] = MAGENTA

	colorsMap["PURPLE"] = PURPLE
	colorsMap["VIOLET"] = VIOLET
	colorsMap["MAROON"] = MAROON
}
