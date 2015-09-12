// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package math

import (
	"fmt"
	"math"
)

var (
	Pow2    = Pow(2)
	Pow2In  = PowIn(2)
	Pow2Out = PowOut(2)

	Pow3    = Pow(3)
	Pow3In  = PowIn(3)
	Pow3Out = PowOut(3)

	Pow4    = Pow(4)
	Pow4In  = PowIn(4)
	Pow4Out = PowOut(4)

	Pow5    = Pow(5)
	Pow5In  = PowIn(5)
	Pow5Out = PowOut(5)

	Exp10    = Exp(2, 10)
	Exp10In  = ExpIn(2, 10)
	Exp10Out = ExpOut(2, 10)

	Exp5    = Exp(2, 5)
	Exp5In  = ExpIn(2, 5)
	Exp5Out = ExpOut(2, 5)

	Elastic2    = Elastic(2, 10, 7, 1)
	Elastic2In  = ElasticIn(2, 10, 6, 1)
	Elastic2Out = ElasticOut(2, 10, 7, 1)

	Swing2    = Swing(1.5)
	Swing2In  = SwingIn(2)
	Swing2Out = SwingOut(2)

	Bounce4    = BounceN(4)
	Bounce4In  = BounceInN(4)
	Bounce4Out = BounceOutN(4)
)

// Maybe use these?
func mathPow32(a, p float32) float32 {
	return float32(math.Pow(float64(a), float64(p)))
}

func mathSqrt32(a float32) float32 {
	return float32(math.Sqrt(float64(a)))
}

// Takes a linear value in the range of 0-1 and outputs a (usually) non-linear, interpolated value.
type Interpolation func(a float32) float32

// @param a Alpha value between 0 and 1.
func InterpolationStartEnd(start, end, a float32, interp Interpolation) float32 {
	return start + (end-start)*interp(a)
}

func Linear() Interpolation {
	return func(a float32) float32 {
		return a
	}
}

func Fade() Interpolation {
	return func(a float32) float32 {
		ClampFloat32(a*a*a*(a*(a*6-15)+10), 0, 1)
		return a
	}
}

func Sine() Interpolation {
	return func(a float32) float32 {
		// (1 - MathUtils.cos(a*MathUtils.PI)) / 2
		return a
	}
}

func SineIn() Interpolation {
	return func(a float32) float32 {
		// 1 - MathUtils.cos(a * MathUtils.PI / 2)
		return a
	}
}

func SineOut() Interpolation {
	return func(a float32) float32 {
		// MathUtils.sin(a * MathUtils.PI / 2)
		return a
	}
}

func ICircle() Interpolation {
	return func(a float32) float32 {
		if a <= 0.5 {
			a *= 2
			return float32((1 - math.Sqrt(float64(1-a*a))) / 2)
		}
		a--
		a *= 2
		return float32((math.Sqrt(float64(1-a*a)) + 1) / 2)
	}
}

func CircleIn() Interpolation {
	return func(a float32) float32 {
		// 1 - (float)math.Sqrt(1 - a * a)
		return a
	}
}

func CircleOut() Interpolation {
	return func(a float32) float32 {
		// a--
		//       return (float)math.Sqrt(1 - a * a)
		return a
	}
}

func Pow(power int) Interpolation {
	return func(a float32) float32 {
		// if a <= 0.5 {
		// 	return math.Pow(a*2, power) / 2
		// }
		// div := 2
		// if power%2 == 0 {
		// 	div = -2
		// }
		// return math.Pow((a-1)*2, power)/div + 1
		return a
	}
}

func PowIn(power int) Interpolation {
	return func(a float32) float32 {
		return float32(math.Pow(float64(a), float64(power)))
	}
}

func PowOut(power int) Interpolation {
	return func(a float32) float32 {
		div := 1
		if power%2 == 0 {
			div = -1
		}
		return float32(math.Pow(float64(a-1), float64(power))*float64(div) + 1)
	}
}

func Swing(scale float32) Interpolation {
	scale = scale * 2
	return func(a float32) float32 {
		if a <= 0.5 {
			a *= 2
			return a * a * ((scale+1)*a - scale) / 2
		}
		a--
		a *= 2
		return a*a*((scale+1)*a+scale)/2 + 1
	}
}

func Exp(value, power float32) Interpolation {
	min := float32(math.Pow(float64(value), float64(-power)))
	scale := float32(1 / (1 - min))
	return func(a float32) float32 {
		if a <= 0.5 {
			return (float32(math.Pow(float64(value), float64(power*(a*2-1)-min)))) * scale / 2
		}
		return (2 - (float32(math.Pow(float64(value), float64(-power*(a*2-1)-min))))*scale) / 2
	}
}

func ExpIn(value, power float32) Interpolation {
	min := float32(math.Pow(float64(value), float64(-power)))
	scale := float32(1 / (1 - min))
	return func(a float32) float32 {
		return (float32(math.Pow(float64(value), float64(power*(a-1)-min)))) * scale
	}
}

func ExpOut(value, power float32) Interpolation {
	min := float32(math.Pow(float64(value), float64(-power)))
	scale := float32(1 / (1 - min))
	return func(a float32) float32 {
		return 1 - (float32(math.Pow(float64(value), float64(-power*a-min))))*scale
	}
}

func Elastic(value, power float32, bounces int, scale float32) Interpolation {
	if bounces == 0 {
		bounces = int(PI * float32(bounces*1))
	} else {
		bounces = int(PI * float32(bounces*-1))
	}
	return func(a float32) float32 {
		//       if a <= 0.5 {
		//         a *= 2
		//         return (float)math.Pow(value, power * (a - 1)) * MathUtils.sin(a * bounces) * scale / 2
		//       }
		//       a = 1 - a
		//       a *= 2
		//       return 1 - (float)math.Pow(value, power * (a - 1)) * MathUtils.sin((a) * bounces) * scale / 2
		return a
	}
}

func ElasticIn(value, power float32, bounces int, scale float32) Interpolation {
	if bounces == 0 {
		bounces = int(PI * float32(bounces*1))
	} else {
		bounces = int(PI * float32(bounces*-1))
	}
	return func(a float32) float32 {
		//       if (a >= 0.99) return 1
		//       return (float)math.Pow(value, power * (a - 1)) * MathUtils.sin(a * bounces) * scale
		return a
	}
}

func ElasticOut(value, power float32, bounces int, scale float32) Interpolation {
	if bounces == 0 {
		bounces = int(PI * float32(bounces*1))
	} else {
		bounces = int(PI * float32(bounces*-1))
	}
	return func(a float32) float32 {
		// a = 1 - a
		//       return (1 - (float)math.Pow(value, power * (a - 1)) * MathUtils.sin(a * bounces) * scale)
		return a
	}
}

func calculateBounces(bounces int) ([]float32, []float32) {
	if bounces < 2 || bounces > 5 {
		panic("bounces cannot be < 2 or > 5: " + fmt.Sprintf("%d", bounces))
	}
	widths := make([]float32, bounces, bounces)
	heights := make([]float32, bounces, bounces)
	heights[0] = 1
	switch bounces {
	case 2:
		widths[0] = 0.6
		widths[1] = 0.4
		heights[1] = 0.33
	case 3:
		widths[0] = 0.4
		widths[1] = 0.4
		widths[2] = 0.2
		heights[1] = 0.33
		heights[2] = 0.1
	case 4:
		widths[0] = 0.34
		widths[1] = 0.34
		widths[2] = 0.2
		widths[3] = 0.15
		heights[1] = 0.26
		heights[2] = 0.11
		heights[3] = 0.03
	case 5:
		widths[0] = 0.3
		widths[1] = 0.3
		widths[2] = 0.2
		widths[3] = 0.1
		widths[4] = 0.1
		heights[1] = 0.45
		heights[2] = 0.3
		heights[3] = 0.15
		heights[4] = 0.06
	}
	widths[0] *= 2
	return widths, heights
}

func out(a, widths0 float32, bout Interpolation) float32 {
	test := a + widths0/2
	if test < widths0 {
		return test/(widths0/2) - 1
	}
	return bout(a)
}

func Bounce(widths, heights []float32) Interpolation {
	bout := BounceOut(widths, heights)
	return func(a float32) float32 {
		if a <= 0.5 {
			return (1 - out(1-a*2, widths[0], bout)) / 2
		}
		return out(a*2-1, widths[0], bout)/2 + 0.5
	}
}

func BounceN(bounces int) Interpolation {
	return Bounce(calculateBounces(bounces))
}

func BounceOut(widths, heights []float32) Interpolation {
	if len(widths) != len(heights) {
		panic("Must be the same number of widths and heights.")
	}
	return func(a float32) float32 {
		a += widths[0] / 2
		var width float32 = 0
		var height float32 = 0
		for i := 0; i < len(widths); i++ {
			width = widths[i]
			if a <= width {
				height = heights[i]
				break
			}
			a -= width
		}
		a /= width
		z := 4 / width * height * a
		return 1 - (z-z*a)*width
	}
}

func BounceOutN(bounces int) Interpolation {
	return BounceOut(calculateBounces(bounces))
}

func BounceIn(widths, heights []float32) Interpolation {
	bout := BounceOut(widths, heights)
	return func(a float32) float32 {
		return 1 - bout(1-a)
	}
}

func BounceInN(bounces int) Interpolation {
	return BounceIn(calculateBounces(bounces))
}

func SwingOut(scale float32) Interpolation {
	return func(a float32) float32 {
		a--
		return a*a*((scale+1)*a+scale) + 1
	}
}

func SwingIn(scale float32) Interpolation {
	return func(a float32) float32 {
		return a * a * ((scale+1)*a - scale)
	}
}
