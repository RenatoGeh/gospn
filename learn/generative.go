package learn

import (
	"github.com/RenatoGeh/gospn/spn"
	"math"
)

// Generative learning with learning rate eta.
func Generative(S spn.SPN, data []map[int]int, eta float64) spn.SPN {
	n := len(data)
	conv := 0.0
	last := 0.0

	// Set root's partial derivative to 1.
	S.Rootify()
	for math.Abs(conv) > 0.0001 {
		s := 0.0
		for i := 0; i < n; i++ {
			// Call soft inference to store values.
			s += S.Value(data[i])
			// Backpropagate through SPN.
			S.Derive()
			// Update weights.
			S.GenUpdate(eta)
		}
		d := s - last
		last = s
		conv = d
	}
	return S
}
