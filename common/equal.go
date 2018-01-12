package common

import "math"

func ApproxEqual(a, b, eps float64) bool {
	return a == b || math.Abs(a-b) < eps
}
