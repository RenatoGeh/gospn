package learn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/test"
	"math"
	"testing"
)

func initSPN() (spn.SPN, spn.VarSet, *spn.Storer) {
	R, val := test.SampleSPN()
	return R, val, spn.NewStorer()
}

func testSPN(R spn.SPN, c common.Collection) {
	test.DoBFS(R, func(s spn.SPN) bool {
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

func testStoreInference(R spn.SPN, val spn.VarSet, storage *spn.Storer, c common.Collection) int {
	itk := storage.NewTicket()
	spn.StoreInference(R, val, itk, storage)

	fmt.Println("\n=== StoreInference Test ===\n")
	itab, _ := storage.Table(itk)
	test.DoBFS(R, func(s spn.SPN) bool {
		v, _ := itab.Single(s)
		fmt.Printf("SPN type: %s, Sc: %v\n  S(X): %.3f\n", s.Type(), s.Sc(), math.Exp(v))
		fmt.Printf("  Cmp: %.3f\n", math.Exp(s.Value(val)))
		return true
	}, c)

	return itk
}

func testDeriveSPN(R spn.SPN, storage *spn.Storer, itk int, c common.Collection) int {
	dtk := storage.NewTicket()
	DeriveSPN(R, storage, dtk, itk, c)

	fmt.Println("\n=== DeriveSPN Test ===\n")
	dtab, _ := storage.Table(dtk)
	test.DoBFS(R, func(s spn.SPN) bool {
		v, _ := dtab.Single(s)
		fmt.Printf("SPN type: %s, Sc: %v\n  dS(X)/dS: %.3f\n", s.Type(), s.Sc(), math.Exp(v))
		return true
	}, c)

	return dtk
}

func testDeriveWeights(R spn.SPN, storage *spn.Storer, itk, dtk int, c common.Collection) {
	wtk := storage.NewTicket()
	DeriveWeights(R, storage, wtk, dtk, itk, c)

	fmt.Println("\n=== DeriveWeights Test ===\n")
	dtab, _ := storage.Table(dtk)
	itab, _ := storage.Table(itk)
	wtab, _ := storage.Table(wtk)
	test.DoBFS(R, func(s spn.SPN) bool {
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
