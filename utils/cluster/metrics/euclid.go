package metrics

import "math"

// Euclidean computes the euclidean distance between two ordered sets of instances.
func Euclidean(p1 []int, p2 []int) float64 {
	// By definition, len(p1)=len(p2).
	var s float64
	for i, p := range p1 {
		l := float64(p - p2[i])
		s += l * l
	}
	return math.Sqrt(s)
}

func EuclideanF(p []float64, q []float64) float64 {
	var s float64
	for i, u := range p {
		l := u - q[i]
		s += l * l
	}
	return math.Sqrt(s)
}
