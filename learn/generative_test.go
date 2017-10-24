package learn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/test"
	"math"
	"testing"
)

// Copy of applyGD, but for testing.
func applyGDTest(S spn.SPN, eta float64, wtk int, storage *spn.Storer, c common.Collection) {
	visited := make(map[spn.SPN]bool)
	wt, _ := storage.Table(wtk)
	c.Give(S)
	visited[S] = true
	fmt.Println("Gradient descent:")
	for !c.Empty() {
		s := c.Take().(spn.SPN)
		ch := s.Ch()
		fmt.Printf("  Type: %s, Sc: %v\n", s.Type(), s.Sc())
		if s.Type() == "sum" {
			sum := s.(*spn.Sum)
			W := sum.Weights()
			dW, _ := wt.Value(s)
			for i := range W {
				delta := eta * math.Exp(dW[i])
				fmt.Printf("    delta = eta * exp(dW[%d]) = %.3f * %.3f = %.3f\n", i, eta, math.Exp(dW[i]), delta)
				fmt.Printf("    W[%d] += delta -> W[%d] = %.3f + %.3f = ", i, i, W[i], delta)
				W[i] += delta
				fmt.Printf("%.3f\n    Sc(Ch[%d])=%v\n    ==\n", W[i], i, ch[i].Sc())
			}
			Normalize(W)
		}
		for _, cs := range ch {
			if cs.Type() != "leaf" && !visited[cs] {
				c.Give(cs)
				visited[cs] = true
			}
		}
	}
}

func TestGenerativeGD(t *testing.T) {
	R, _ := test.SampleSPN()
	storage := spn.NewStorer()

	fmt.Println("\n=== GenerativeGD Test ===\n")
	data := make([]spn.VarSet, 16)
	for i := range data {
		data[i] = make(spn.VarSet)
	}
	d := 0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				for l := 0; l < 2; l++ {
					data[d][0], data[d][1], data[d][2], data[d][3] = i, j, k, l
					d++
				}
			}
		}
	}
	fmt.Println("Data:")
	for i, I := range data {
		fmt.Printf("  data[%d] = {\n", i)
		for k, v := range I {
			fmt.Printf("    [%d]: %d,\n", k, v)
		}
		fmt.Println("}")
	}

	c := &common.Queue{}
	eta := 0.1
	wtk, stk, itk := storage.NewTicket(), storage.NewTicket(), storage.NewTicket()
	for _, I := range data {
		// Store inference values under T[itk].
		spn.StoreInference(R, I, itk, storage)
		// Store SPN derivatives under T[stk].
		DeriveSPN(R, storage, stk, itk, c)
		// Store weights derivatives under T[wtk].
		DeriveWeights(R, storage, wtk, stk, itk, c)
		// Apply gradient descent.
		applyGDTest(R, eta, wtk, storage, c)
	}

	test.DoBFS(R, func(s spn.SPN) bool {
		t := s.Type()
		fmt.Printf("Type: %s, Sc: %v\n", t, s.Sc())
		if t == "sum" {
			W := s.(*spn.Sum).Weights()
			for i, w := range W {
				fmt.Printf("  W[%d] = %.3f\n", i, w)
			}
		}
		return true
	}, c)
}
