package spn

import (
	"fmt"
	"math"

	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/sys"
)

// Some of the following functions are non-recursive versions of equivalent spn.SPN methods. They
// are done using a Queue or Stack to perform the graph search instead of the recursion call stack.
// When the SPN is dense, running the recursive versions can take exponential time (as we do not
// account for already visited vertices). In these static function versions, all searches are done
// in time linear to the graphs. For this reason, unless the SPN is a tree (or the graph sparse
// enough), the preferred method is using the static function version. When the SPN is a tree, the
// best method is the recursive version, as it takes less memory and same time usage in average
// when compared to the static versions.

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

	// Get topological order.
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

	// Get topological order.
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

	// Find max values.
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

	// Find MAP states.
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

func norm(v []float64) {
	var n float64
	for i := range v {
		n += v[i]
	}
	for i := range v {
		v[i] /= n
	}
}

// NormalizeSPN recursively normalizes the SPN S.
func NormalizeSPN(S SPN) SPN {
	Q := common.Queue{}
	V := make(map[SPN]bool)

	Q.Enqueue(S)
	V[S] = true

	for !Q.Empty() {
		s := Q.Dequeue().(SPN)
		if s.Type() == "sum" {
			z := s.(*Sum)
			W := z.Weights()
			norm(W)
		}
		ch := s.Ch()
		for _, c := range ch {
			if !V[c] {
				Q.Enqueue(c)
				V[c] = true
			}
		}
	}

	return S
}

// ComputeHeight computes the height of a certain SPN S.
func ComputeHeight(S SPN) int {
	T := common.Stack{}
	V := make(map[SPN]int)

	T.Push(S)
	V[S] = 0

	var h int
	for !T.Empty() {
		s := T.Pop().(SPN)
		if s.Type() == "leaf" && V[s] > h {
			h = V[s]
		}
		ch := s.Ch()
		for _, c := range ch {
			if _, e := V[c]; !e {
				T.Push(c)
				V[c] = V[s] + 1
			}
		}
	}

	return h
}

// ComputeScope computes the scope of a certain SPN S.
func ComputeScope(S SPN) []int {
	if _sc := S.rawSc(); _sc != nil {
		return _sc
	}

	Q, T := common.Queue{}, common.Stack{}
	V := make(map[SPN]bool)

	Q.Enqueue(S)
	T.Push(S)
	V[S] = true

	for !Q.Empty() {
		s := Q.Dequeue().(SPN)
		ch := s.Ch()
		for _, c := range ch {
			if !V[c] && c.Type() != "leaf" {
				Q.Enqueue(c)
				T.Push(c)
				V[c] = true
			}
		}
	}

	for !T.Empty() {
		s := T.Pop().(SPN)
		ch := s.Ch()
		sc := make(map[int]bool)
		for _, c := range ch {
			csc := c.rawSc()
			for _, v := range csc {
				sc[v] = true
			}
		}
		nsc := make([]int, len(sc))
		var i int
		for k, _ := range sc {
			nsc[i] = k
			i++
		}
		s.setRawSc(nsc)
	}

	return S.rawSc()
}

// Complete returns whether the SPN is complete.
func Complete(S SPN) bool {
	ComputeScope(S)
	Q := common.Queue{}
	V := make(map[SPN]bool)

	Q.Enqueue(S)
	V[S] = true

	for !Q.Empty() {
		s := Q.Dequeue().(SPN)
		ch := s.Ch()
		if s.Type() == "sum" {
			sc := s.rawSc()
			v := make(map[int]int)
			for _, u := range sc {
				v[u]++
			}
			for _, c := range ch {
				csc := c.rawSc()
				// Invariant: ComputeScope guarantees that there will be no duplicates.
				if len(csc) != len(sc) {
					sys.Printf("len(csc)=%d != len(sc)=%d\n", len(csc), len(sc))
					sys.Printf("%v\n%v\n", csc, sc)
					return false
				}
				for _, u := range csc {
					_, e := v[u]
					if !e {
						sys.Printf("v[%d] does not exist\n", u)
						return false
					}
					v[u]++
				}
			}
			k := len(ch) + 1
			for _, u := range v {
				if u != k {
					sys.Printf("u=%d != k=%d\n", u, k)
					return false
				}
			}
		}
		for _, c := range ch {
			if !V[S] && c.Type() != "leaf" {
				Q.Enqueue(c)
				V[c] = true
			}
		}
	}

	return true
}

// Decomposable returns whether the SPN is decomposable.
func Decomposable(S SPN) bool {
	ComputeScope(S)
	Q := common.Queue{}
	V := make(map[SPN]bool)

	Q.Enqueue(S)
	V[S] = true

	for !Q.Empty() {
		s := Q.Dequeue().(SPN)
		ch := s.Ch()
		if s.Type() == "product" {
			sc := s.rawSc()
			v := make(map[int]int)
			for _, u := range sc {
				v[u]++
			}
			for _, c := range ch {
				csc := c.rawSc()
				// Invariant: ComputeScope guarantees that there will be no duplicates.
				for _, u := range csc {
					_, e := v[u]
					if !e {
						return false
					}
					v[u]++
				}
			}
			for _, u := range v {
				if u != 2 {
					return false
				}
			}
		}
		for _, c := range ch {
			if !V[S] && c.Type() != "leaf" {
				Q.Enqueue(c)
				V[c] = true
			}
		}
	}

	return true
}
