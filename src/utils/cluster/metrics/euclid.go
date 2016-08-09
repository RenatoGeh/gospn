package metrics

import "math"

func Euclidean(p1 []int, p2 []int) float64 {
	// By definition, len(p1)=len(p2).
	n, s := len(p1), 0.0
	for i := 0; i < n; i++ {
		sqr := float64(p1[i] - p2[i])
		s += sqr * sqr
	}
	return math.Sqrt(s)
}
