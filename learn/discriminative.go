package learn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/spn"
	"math"
)

// Discriminative learning with learning rate eta.
func Discriminative(S spn.SPN, data []map[int]int, Y []int, X []int, eta float64) spn.SPN {
	n := len(data)
	conv := 1.0
	last := 0.0

	fmt.Printf("Discriminative learning with %d instances.\n", n)

	S.SetStore(true)
	for math.Abs(conv) > 0.01 {
		s := 0.0
		klast := 0.0
		for i := 0; i < n; i++ {
			S.RResetDP("")
			S.Rootify("cpnode")
			S.Rootify("epnode")
			C, E := make(spn.VarSet), make(spn.VarSet)
			ny, nx := len(Y), len(X)
			for j := 0; j < ny; j++ {
				C[Y[j]] = data[i][Y[j]]
			}
			for j := 0; j < nx; j++ {
				_v := data[i][X[j]]
				C[X[j]] = _v
				E[X[j]] = _v
			}
			ds := spn.NewDiscStorer(S, C, E)
			// Stores correct/guess values.
			fmt.Println("Storing correct/guess soft inference values...")
			S.Soft(C, "correct")
			// Derive correct/guess nodes.
			fmt.Println("Derivating correct/guess nodes...")
			S.Derive("cpweight", "cpnode", "correct")
			// Stores expected values.
			fmt.Println("Storing expected soft inference values...")
			S.Soft(E, "expected")
			// Derive expected nodes.
			fmt.Println("Derivating expected nodes...")
			S.Derive("epweight", "epnode", "expected")
			// Update weights.
			fmt.Println("Updating weights...")
			S.DiscUpdate(eta, ds, "cpweight", "epweight")

			fmt.Printf("Adding convergence diff for instance %d...\n", i)
			k := S.Value(data[i])
			s += k - klast
			klast = k

			C = nil
			E = nil
		}
		d := s - last
		last = s
		conv = d
		fmt.Printf("Discriminative Learning diff: %.5f\n", math.Abs(conv))
	}
	return S
}
