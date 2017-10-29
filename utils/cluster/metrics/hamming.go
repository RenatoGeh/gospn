package metrics

// Euclidean computes the euclidean distance between two ordered sets of instances.
func Hamming(p1 []int, p2 []int) int {
	// By definition, len(p1)=len(p2).
	n, s := len(p1), 0
	for i := 0; i < n; i++ {
		if p1[i] != p2[i] {
			s++
		}
	}
	return s
}

func HammingF(p1 []int, p2 []int) float64 {
	// By definition, len(p1)=len(p2).
	n, s := len(p1), 0.0
	for i := 0; i < n; i++ {
		if p1[i] != p2[i] {
			s++
		}
	}
	return s
}
