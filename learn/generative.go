package learn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/spn"
	"math"
)

// Generative learning with learning rate eta.
func Generative(S spn.SPN, data []map[int]int, eta float64) spn.SPN {
	n := len(data)
	conv := 1.0
	last := 0.0
	klast := 0.0

	fmt.Println("Running generative learning.")

	// Set root's partial derivative to 1.
	S.Rootify()
	for math.Abs(conv) > 0.01 {
		// Reset SPN's DP table.
		S.ResetDP()
		s := 0.0
		for i := 0; i < n; i++ {
			// Call soft inference to store values.
			k := S.Value(data[i])
			s += k - klast
			klast = k
			fmt.Println("Computed and stored soft inference values.")
			// Backpropagate through SPN.
			S.Derive()
			fmt.Println("Backpropagated through SPN computing derivatives.")
			// Update weights.
			S.GenUpdate(eta)
			fmt.Printf("Weight Updated according to instance %d.\n", i)
		}
		d := s - last
		last = s
		conv = d
		fmt.Printf("Generative Learning diff: %.5f\n", math.Abs(conv))
	}
	return S
}
