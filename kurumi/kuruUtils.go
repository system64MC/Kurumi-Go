package kurumi

import "math"

func Clamp(low, val, high int) int {
	if val < low {
		return low
	}
	if val > high {
		return high
	}
	return val
}

func ClampF64(low, val, high float64) float64 {
	if val < low {
		return low
	}
	if val > high {
		return high
	}
	return val
}

func moduloFix(a, b int) int {
	return ((a % b) + b) % b
}

func moduloF64(a float64, b float64) float64 {
	tmp := math.Mod(a, b) + b
	return math.Mod(tmp, b)
}

func moduloInt(a int, b int) int {
	return ((a % b) + b) % b
}

func minandmax(values []uint8) (uint8, uint8) {
	min := values[0] //assign the first element equal to min
	max := values[0] //assign the first element equal to max
	for _, number := range values {
		if number < min {
			min = number
		}
		if number > max {
			max = number
		}
	}
	return min, max
}

func minandmaxFloat(values []float64) (float64, float64) {
	min := values[0] //assign the first element equal to min
	max := values[0] //assign the first element equal to max
	for _, number := range values {
		if number < min {
			min = number
		}
		if number > max {
			max = number
		}
	}
	min = math.Abs(min)
	return min, max
}

func lerp(x, y, a float64) float64 {
	return x*(1-a) + y*a
}

func LinearInterpolation(x1 float64, y1 float64, x2 float64, y2 float64, x float64) float64 {
	slope := (y2 - y1) / (x2 - x1)
	return y1 + (slope * (x - x1))
}
