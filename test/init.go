package test

import (
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
)

// SampleSPN returns a sample SPN for testing.
func SampleSPN() (spn.SPN, spn.VarSet) {
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

	return R, val
}

// DoBFS takes an SPN and does a graph search on it. At every node, it calls a function f, passing
// the current node as its argument. Collection c determines what kind of graph search DoBFS is to
// perform.
func DoBFS(S spn.SPN, f func(spn.SPN) bool, c common.Collection) {
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
	visited = nil
	sys.Free()
}
