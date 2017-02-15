package testspn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/spn"
)

const (
	nU = 4
	nT = 3
)

// TestSPN tests discriminative learning and backpropagation.
func TestSPN() {
	// Creating structure.
	Y := []*spn.Indicator{spn.NewIndicator(0, 1), spn.NewIndicator(0, 0)}
	X := []*spn.Indicator{spn.NewIndicator(1, 1), spn.NewIndicator(1, 0)}
	U := make([]*spn.Sum, nU)
	for i := 0; i < nU; i++ {
		U[i] = spn.NewSum()
	}
	wU := [][]float64{{0.3, 0.7},
		{0.6, 0.4},
		{0.8, 0.2},
		{0.5, 0.5}}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			U[i].AddChildW(Y[j], wU[i][j])
		}
	}
	for i := 2; i < 4; i++ {
		for j := 0; j < 2; j++ {
			U[i].AddChildW(X[j], wU[i][j])
		}
	}

	T := make([]*spn.Product, nT)
	for i := 0; i < nT; i++ {
		T[i] = spn.NewProduct()
	}

	T[0].AddChild(U[0])
	T[0].AddChild(U[2])
	T[1].AddChild(U[0])
	T[1].AddChild(U[3])
	T[2].AddChild(U[1])
	T[2].AddChild(U[3])

	S := spn.NewSum()
	wS := []float64{0.2, 0.5, 0.3}
	for i := 0; i < nT; i++ {
		S.AddChildW(T[i], wS[i])
	}

	// Testing values.
	fmt.Println("Testing soft inference values...")
	S.SetStore(false)
	val := make(spn.VarSet)
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			val[0] = i
			val[1] = j
			fmt.Printf("Pr(Y=%d, X=%d)=", val[0], val[1])
			fmt.Printf("%.10f\n", S.Value(val))
			S.RResetDP("")
		}
	}

	// Testing backprop.
	fmt.Println("\nTesting derivatives (backpropagation)...")
	S.SetStore(true)
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if y < 2 {
				val[0] = y
			} else {
				delete(val, 0)
			}
			if x < 2 {
				val[1] = x
			} else {
				delete(val, 1)
			}
			S.RResetDP("")
			fmt.Printf("Pr( ")
			if y < 2 {
				fmt.Printf("Y=%d ", val[0])
			}
			if x < 2 {
				fmt.Printf("X=%d ", val[1])
			}
			fmt.Printf(")=")
			fmt.Printf("%.10f\n", S.Soft(val, "inference"))
			S.Rootify("node")
			S.RootDerive("weights", "node", "inference", spn.SOFT)
			dn, _ := S.Stored("node")
			dW := S.PWeights("weights")
			fmt.Printf("  dS/dS = %.10f\n", dn)
			for i := 0; i < len(dW); i++ {
				fmt.Printf("    dS/dw_%d = %.10f\n", i, dW[i])
			}
			fmt.Println("===========")
			for i := 0; i < nT; i++ {
				dn, _ := T[i].Stored("node")
				soft, _ := T[i].Stored("inference")
				fmt.Printf("  dT_%d/dS = %.10f\n    T_%d = %.10f\n", i, dn, i, soft)
			}
			fmt.Println("===========")
			for i := 0; i < nU; i++ {
				dn, _ := U[i].Stored("node")
				dW = U[i].PWeights("weights")
				soft, _ := U[i].Stored("inference")
				fmt.Printf("  dU_%d/dT = %.10f\n    U_%d = %.10f\n", i, dn, i, soft)
				for j := 0; j < len(dW); j++ {
					fmt.Printf("    dS/dw_%d = %.10f\n", j, dW[j])
				}
			}
			fmt.Println("===========")
		}
	}
}
