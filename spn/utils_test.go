package spn

import (
	"github.com/RenatoGeh/gospn/common"
	"reflect"
	"testing"
)

func sampleSPN() SPN {
	R := NewSum()
	P1, P2 := NewProduct(), NewProduct()
	f1, f2 := NewMultinomial(0, []float64{0.9, 0.1}), NewMultinomial(1, []float64{0.6, 0.4})
	S1, S2, S3, S4 := NewSum(), NewSum(), NewSum(), NewSum()
	Y11, Y12 := NewMultinomial(2, []float64{0.8, 0.2}), NewMultinomial(2, []float64{0.3, 0.7})
	Y21, Y22 := NewMultinomial(3, []float64{0.4, 0.6}), NewMultinomial(3, []float64{0.9, 0.1})

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

	return R
}

func doBFS(S SPN, f func(SPN) bool) {
	c := common.Queue{}
	visited := make(map[SPN]bool)
	c.Enqueue(S)
	visited[S] = true

	for !c.Empty() {
		s := c.Dequeue().(SPN)
		if !f(s) {
			break
		}
		ch := s.Ch()
		for _, cs := range ch {
			if !visited[cs] {
				c.Enqueue(cs)
				visited[cs] = true
			}
		}
	}
	visited = nil
}

// TestStoreInference tests StoreInference.
func TestStoreInference(t *testing.T) {
	S := sampleSPN()
	I := make(VarSet)
	I[0], I[2] = 1, 0
	st := NewStorer()
	_, tk := StoreInference(S, I, -1, st)
	doBFS(S, func(s SPN) bool {
		v := s.Value(I)
		u, _ := st.Single(tk, s)
		if v != u {
			t.Errorf("Expected (Value) %f == %f (StoreInference) ? true, got false.", v, u)
			return false
		}
		return true
	})
}

// TestStoreMap tests StoreMAP.
func TestStoreMAP(t *testing.T) {
	S := sampleSPN()
	I := make(VarSet)
	I[0], I[2] = 1, 0
	st := NewStorer()
	_, tk, M := StoreMAP(S, I, -1, st)
	doBFS(S, func(s SPN) bool {
		v := s.Max(I)
		u, _ := st.Single(tk, s)
		if v != u {
			t.Errorf("Expected (Max) %f == %f (StoreMAP) ? true, got false.", v, u)
			return false
		}
		return true
	})
	N, _ := S.ArgMax(I)
	if !reflect.DeepEqual(M, N) {
		t.Errorf("Expected MAP states to be equal, got different.\n  M=%v\n  N=%v", M, N)
	}
}
