package goutil

import "math"

func MinInt(x, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

func MaxInt(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}
