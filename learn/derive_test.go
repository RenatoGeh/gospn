package learn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/spn"
	"math"
	"testing"
)

func initSPN() (spn.SPN, spn.VarSet, *Storer) {
	R := spn.NewSum()
	P1, P2 := spn.NewProduct(), spn.NewProduct()
	f1, f2 := spn.NewMultinomial(0, []float64{0.9, 0.1}), spn.NewMultinomial(1, []float64{0.6, 0.4})
	S1, S2, S3, S4 := spn.NewSum(), spn.NewSum(), spn.NewSum(), spn.NewSum()
	Y11, Y12 := spn.NewMultinomial(2, []float64{0.8, 0.2}), spn.NewMultinomial(2, []float64{0.3, 0.7})
	Y21, Y22 := spn.NewMultinomial(3, []float64{0.4, 0.6}), spn.NewMultinomial(3, []float64{0.9, 0.1})

	R.AddChildW(P1, 0.3)
	R.AddChildW(P2, 0.7)
	P1.AddChild(f1)
	P1.AddChild(S1)
	P1.AddChild(S3)
	P2.AddChild(S2)
	P2.AddChild(S4)
	P2.AddChild(f2)
	S1.AddChildW(Y11, 0.6)
	S1.AddChildW(Y12, 0.4)
	S2.AddChildW(Y11, 0.2)
	S2.AddChildW(Y12, 0.8)
	S3.AddChildW(Y21, 0.9)
	S3.AddChildW(Y22, 0.1)
	S4.AddChildW(Y21, 0.5)
	S4.AddChildW(Y22, 0.5)

	val := make(spn.VarSet)
	val[0] = 0
	val[1] = 0
	val[2] = 0
	val[3] = 0

	return R, val, NewStorer()
}

func doBFS(S spn.SPN, f func(spn.SPN) bool, c common.Collection) {
	if c == nil {
		c = &common.Queue{}
	}
	visited := make(map[spn.SPN]bool)
	c.Give(S)
	visited[S] = true

	for !c.Empty() {
		s := c.Take().(spn.SPN)
		if !f(s) {
			break
		}
		ch := s.Ch()
		for _, cs := range ch {
			if !visited[cs] {
				c.Give(cs)
				visited[cs] = true
			}
		}
	}
}

func testSPN(R spn.SPN, c common.Collection) {
	doBFS(R, func(s spn.SPN) bool {
		t := s.Type()
		fmt.Printf("SPN type: %s\n", t)
		if t == "sum" {
			W := (s.(*spn.Sum)).Weights()
			for i, w := range W {
				fmt.Printf("  W[%d] = %.3f\n", i, w)
			}
		} else if t == "leaf" {
			d := (s.(*spn.Multinomial)).Pr()
			for i, p := range d {
				fmt.Printf("  Pr[%d] = %.3f\n", i, p)
			}
		}
		return true
	}, c)
}

func testStoreInference(R spn.SPN, val spn.VarSet, storage *Storer, c common.Collection) int {
	itk := storage.NewTicket()
	StoreInference(R, val, itk, storage)

	fmt.Println("\n=== StoreInference Test ===\n")
	itab, _ := storage.Table(itk)
	doBFS(R, func(s spn.SPN) bool {
		v, _ := itab.Single(s)
		fmt.Printf("SPN type: %s, Sc: %v\n  S(X): %.3f\n", s.Type(), s.Sc(), math.Exp(v))
		fmt.Printf("  Cmp: %.3f\n", math.Exp(s.Value(val)))
		return true
	}, c)

	return itk
}

func testDeriveSPN(R spn.SPN, storage *Storer, itk int, c common.Collection) int {
	dtk := storage.NewTicket()
	DeriveSPN(R, storage, dtk, itk, c)

	fmt.Println("\n=== DeriveSPN Test ===\n")
	dtab, _ := storage.Table(dtk)
	doBFS(R, func(s spn.SPN) bool {
		v, _ := dtab.Single(s)
		fmt.Printf("SPN type: %s, Sc: %v\n  dS(X)/dS: %.3f\n", s.Type(), s.Sc(), math.Exp(v))
		return true
	}, c)

	return dtk
}

func testDeriveWeights(R spn.SPN, storage *Storer, itk, dtk int, c common.Collection) {
	wtk := storage.NewTicket()
	DeriveWeights(R, storage, wtk, dtk, itk, c)

	fmt.Println("\n=== DeriveWeights Test ===\n")
	dtab, _ := storage.Table(dtk)
	itab, _ := storage.Table(itk)
	wtab, _ := storage.Table(wtk)
	doBFS(R, func(s spn.SPN) bool {
		v, _ := wtab.Value(s)
		p, _ := dtab.Single(s)
		fmt.Printf("SPN type: %s, Sc: %v\n", s.Type(), s.Sc())
		ch := s.Ch()
		for i, l := range v {
			q, _ := itab.Single(ch[i])
			fmt.Printf("  (dS(X)/dW)[%d]: %.3f == %.3f * %.3f == %.3f == dS(X)/dS * S(X)\n",
				i, math.Exp(l), math.Exp(p), math.Exp(q), math.Exp(p+q))
		}
		return true
	}, c)
}

func TestDerive(t *testing.T) {
	R, val, storage := initSPN()

	c := &common.Queue{}
	testSPN(R, c)
	itk := testStoreInference(R, val, storage, c)
	dtk := testDeriveSPN(R, storage, itk, c)
	testDeriveWeights(R, storage, itk, dtk, c)
}
