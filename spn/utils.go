package spn

import (
	"fmt"
	"math"

	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/sys"
)

// StoreInference takes an SPN S and stores the values for an instance I on a DP table storage
// at the position designated by the ticket tk. Returns S and the ticket used (if tk < 0,
// StoreInference creates a new ticket).
func StoreInference(S SPN, I VarSet, tk int, storage *Storer) (SPN, int) {
	if tk < 0 {
		tk = storage.NewTicket()
	}

	visited := make(map[SPN]bool)
	c, _c := &common.Stack{}, &common.Queue{}
	c.Give(S)
	visited[S] = true
	_c.Give(S)

	for !_c.Empty() {
		s := _c.Take().(SPN)
		ch := s.Ch()
		for _, cs := range ch {
			if !visited[cs] {
				_c.Give(cs)
				c.Give(cs)
				visited[cs] = true
			}
		}
	}

	_c, visited = nil, nil // free memory as soon as soon as the garbage collector allows
	sys.Free()
	table, _ := storage.Table(tk)
	for !c.Empty() {
		s := c.Take().(SPN)
		switch t := s.Type(); t {
		case "leaf":
			table.StoreSingle(s, s.Value(I))
		case "sum":
			sum := s.(*Sum)
			ch := sum.Ch()
			W := sum.Weights()
			n := len(ch)
			vals := make([]float64, n)
			for i, cs := range ch {
				v, e := table.Single(cs)
				if !e {
					// Should never occur. Just in case what I thought of is flawed.
					fmt.Println("Something terrible has just happened. (StoreInference:learn/derive.go)")
				}
				vals[i] = v + math.Log(W[i])
			}
			table.StoreSingle(s, sum.Compute(vals))
		case "product":
			prod := s.(*Product)
			ch := prod.Ch()
			n := len(ch)
			vals := make([]float64, n)
			for i, cs := range ch {
				vals[i], _ = table.Single(cs)
			}
			table.StoreSingle(s, prod.Compute(vals))
		}
	}
	c = nil
	sys.Free()
	return S, tk
}

// StoreMAP takes an SPN S and stores the MAP values for an instance I on a DP table storage
// at the position designated by the ticket tk. Returns S and the ticket used (if tk < 0,
// StoreMAP creates a new ticket).
func StoreMAP(S SPN, I VarSet, tk int, storage *Storer) (SPN, int, VarSet) {
	if tk < 0 {
		tk = storage.NewTicket()
	}

	Q, T := common.Queue{}, common.Stack{}
	V := make(map[SPN]bool)
	tab, _ := storage.Table(tk)

	Q.Enqueue(S)
	T.Push(S)
	V[S] = true

	for !Q.Empty() {
		s := Q.Dequeue().(SPN)
		ch := s.Ch()
		for _, c := range ch {
			if !V[c] {
				Q.Enqueue(c)
				T.Push(c)
				V[c] = true
			}
		}
	}

	for !T.Empty() {
		s := T.Pop().(SPN)
		switch t := s.Type(); t {
		case "leaf":
			m := s.Max(I)
			tab.StoreSingle(s, m)
		case "sum":
			sum := s.(*Sum)
			W := sum.Weights()
			ch := s.Ch()
			mv := math.Inf(-1)
			for i, c := range ch {
				v, _ := tab.Single(c)
				u := math.Log(W[i]) + v
				if u > mv {
					mv = u
				}
			}
			tab.StoreSingle(s, mv)
		case "product":
			ch := s.Ch()
			var v float64
			for _, c := range ch {
				cv, _ := tab.Single(c)
				v += cv
			}
			tab.StoreSingle(s, v)
		}
	}

	V = make(map[SPN]bool)
	Q.Enqueue(S)
	V[S] = true
	M := make(VarSet)

	for !Q.Empty() {
		s := Q.Dequeue().(SPN)
		switch t := s.Type(); t {
		case "leaf":
			N, _ := s.ArgMax(I)
			for k, v := range N {
				M[k] = v
			}
		case "sum":
			sum := s.(*Sum)
			W := sum.Weights()
			ch := s.Ch()
			m := math.Inf(-1)
			var mvc SPN
			for i, c := range ch {
				v, _ := tab.Single(c)
				u := math.Log(W[i]) + v
				if u > m && !V[c] {
					m, mvc = u, c
				}
			}
			Q.Enqueue(mvc)
		case "product":
			ch := s.Ch()
			for _, c := range ch {
				if !V[c] {
					Q.Enqueue(c)
				}
			}
		}
		V[s] = true
	}

	return S, tk, M
}
