package learn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/learn/parameters"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/test"
	"math"
	"testing"
)

func initSimpleSPN() (spn.SPN, []spn.SPN) {
	R := spn.NewSum()
	parameters.Bind(R, parameters.New(true, false, 0.0, parameters.SoftGD, 0.1, 1, 0, 0, 1))
	P1, P2 := spn.NewProduct(), spn.NewProduct()
	S1, S2, S3, S4 := spn.NewSum(), spn.NewSum(), spn.NewSum(), spn.NewSum()
	X11, X12 := spn.NewMultinomial(0, []float64{0.3, 0.7}), spn.NewMultinomial(0, []float64{0.6, 0.4})
	X21, X22 := spn.NewMultinomial(1, []float64{0.1, 0.9}), spn.NewMultinomial(1, []float64{0.5, 0.5})
	R.AddChildW(P1, 0.4)
	R.AddChildW(P2, 0.6)
	P1.AddChild(S1)
	P1.AddChild(S2)
	P2.AddChild(S3)
	P2.AddChild(S4)
	S1.AddChildW(X11, 0.2)
	S1.AddChildW(X12, 0.8)
	S2.AddChildW(X11, 0.6)
	S2.AddChildW(X12, 0.4)
	S3.AddChildW(X21, 0.5)
	S3.AddChildW(X22, 0.5)
	S4.AddChildW(X21, 0.7)
	S4.AddChildW(X22, 0.3)
	label := []spn.SPN{R, P1, P2, S1, S2, S3, S4, X11, X12, X21, X22}
	return R, label
}

func approxEqual(p, q, eps float64) bool {
	return p == q || math.Abs(p-q) < eps
}

func TestSimpleSPN(t *testing.T) {
	R, L := initSimpleSPN()
	n := len(L)
	const (
		EPS = 1e-15
	)
	fmt.Println("Printing sample SPN...")
	for i := range L {
		s := L[i]
		fmt.Printf("%d:\n  Type: %s\n", i, s.Type())
		if s.Type() == "sum" {
			w := s.(*spn.Sum).Weights()
			fmt.Printf("  Weights: { %.5f", w[0])
			for j := 1; j < len(w); j++ {
				fmt.Printf(", %.5f", w[j])
			}
			fmt.Printf(" }\n")
		} else if s.Type() == "leaf" {
			fmt.Printf("  Var: %d\n", s.(*spn.Multinomial).Sc()[0])
			continue
		}
		ch := s.Ch()
		fmt.Printf("  Ch: { ")
		for _, c := range ch {
			var k int
			for l, cc := range L {
				if cc == c {
					k = l
				}
			}
			fmt.Printf("%d ", k)
		}
		fmt.Printf(" }\n")
	}
	I := make(spn.VarSet)
	I[0], I[1] = 1, 0
	st := spn.NewStorer()
	itk, stk, wtk := st.NewTicket(), st.NewTicket(), st.NewTicket()
	spn.StoreInference(R, I, itk, st)
	// Correct inference values
	fmt.Println("Testing inference values...")
	ci := []float64{0.14632, 0.2668, 0.066, 0.46, 0.58, 0.3, 0.22, 0.7, 0.4, 0.1, 0.5}
	for i, s := range L {
		v, _ := st.Single(itk, s)
		if !approxEqual(v, math.Log(ci[i]), EPS) {
			t.Errorf("INF Label %d: Expected %.5f, got %.5f", i, ci[i], math.Exp(v))
		}
	}
	Q := common.Queue{}
	DeriveSPN(R, st, stk, itk, &Q)
	// Correct SPN derivatives
	fmt.Println("Testing SPN derivatives...")
	cs := []float64{1.0, 0.4, 0.6, 0.232, 0.184, 0.132, 0.18}
	for i := 0; i < n-4; i++ {
		s := L[i]
		v, _ := st.Single(stk, s)
		if !approxEqual(v, math.Log(cs[i]), EPS) {
			t.Errorf("DS Label %d: Expected %.5f, got %.5f", i, cs[i], math.Exp(v))
		}
	}
	DeriveWeights(R, st, wtk, stk, itk, &Q)
	// Correct weight derivatives
	fmt.Println("Testing weight derivatives...")
	cw := [][]float64{[]float64{0.2668, 0.066}, nil, nil, []float64{0.1624, 0.0928},
		[]float64{0.1288, 0.0736}, []float64{0.0132, 0.066}, []float64{0.018, 0.09}}
	for i := 0; i < n-4; i++ {
		p := cw[i]
		if p == nil {
			continue
		}
		s := L[i]
		v, _ := st.Value(wtk, s)
		for j := range p {
			if !approxEqual(math.Log(p[j]), v[j], EPS) {
				t.Errorf("DW Label %d-%d: Expected %.5f, got %.5f", i, j, p[j], math.Exp(v[j]))
			}
		}
	}
	applyGD(R, 0.1, wtk, st, &Q, false)
	// Correct learned weights
	fmt.Println("Testing gradient descent application...")
	clw := [][]float64{[]float64{0.42668, 0.6066}, nil, nil, []float64{0.21624, 0.80928},
		[]float64{0.61288, 0.40736}, []float64{0.50132, 0.5066}, []float64{0.7018, 0.309}}
	for i := 0; i < n-4; i++ {
		p := clw[i]
		if p == nil {
			continue
		}
		w := L[i].(*spn.Sum).Weights()
		for j := range w {
			if !approxEqual(p[j], w[j], EPS) {
				t.Errorf("GD Label %d-%d: Expected %.5f, got %.5f", i, j, p[j], w[j])
			}
		}
	}
}

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

	fmt.Println("\n=== GenerativeGD Test ===")
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
