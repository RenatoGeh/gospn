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

	sys.Println("Finding topological sort...")
	// Get topological order.
	O := common.Queue{}
	TopSortTarjan(S, &O)
	sys.Println("Found topological sort.")
	sys.Free()

	sys.Println("Starting inference computations...")
	table, _ := storage.Table(tk)
	for !O.Empty() {
		s := O.Dequeue().(SPN)
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

	// Get topological order.
	T := common.Queue{}
	TopSortTarjan(S, &T)

	tab, _ := storage.Table(tk)
	// Find max values.
	for !T.Empty() {
		s := T.Dequeue().(SPN)
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

	Q := common.Queue{}
	V := make(map[SPN]bool)
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
			var mv []SPN
			for i, c := range ch {
				v, _ := tab.Single(c)
				u := math.Log(W[i]) + v
				if u > m {
					mv, m = []SPN{c}, u
				} else if u == m {
					mv = append(mv, c)
				}
			}
			// Randomly break tie.
			mvc := mv[sys.Random.Intn(len(mv))]
			if mvc != nil && !V[mvc] {
				Q.Enqueue(mvc)
				V[mvc] = true
			}
		case "product":
			ch := s.Ch()
			for _, c := range ch {
				if !V[c] {
					Q.Enqueue(c)
					V[c] = true
				}
			}
		}
	}

	return S, tk, M
}

func norm(v []float64) {
	var min float64
	for _, u := range v {
		if u < min {
			min = u
		}
	}
	if min < 0 {
		for i := range v {
			v[i] += min
		}
	}

	var norm float64
	for i := range v {
		norm += v[i]
	}
	for i := range v {
		v[i] /= norm
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
	T := common.Queue{}
	TopSortTarjan(S, &T)

	for !T.Empty() {
		s := T.Dequeue().(SPN)
		ch := s.Ch()
		if len(ch) > 0 {
			msc := make(map[int]bool)
			for _, c := range ch {
				csc := c.rawSc()
				for _, v := range csc {
					msc[v] = true
				}
			}
			var sc []int
			for k, _ := range msc {
				sc = append(sc, k)
			}
			s.setRawSc(sc)
		}
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
			n := len(sc)
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
				n -= len(csc)
			}
			for _, u := range v {
				if u != 2 {
					return false
				}
			}
			if n != 0 {
				return false
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

// TraceMAP returns the max child index of each sum node in a map. We assume decomposability and
// completeness. When this condition is not met, one arbitrary MAP state is chosen. When the SPN is
// both decomposable and complete, it is easy to see that the induced MAP trace of the SPN's graph
// is a tree, and thus no two paths from the root to a leaf intersect. For the negative case, there
// can be two paths that do intersect, and thus we could have randomly chosen different max
// children in case of ties. In this case, TraceMAP chooses the first child it finds to meet the
// criteria.
func TraceMAP(S SPN, I VarSet) map[SPN]int {
	st := NewStorer()
	_, itk, _ := StoreMAP(S, I, -1, st)
	T, _ := st.Table(itk)
	trace := make(map[SPN]int)

	V := make(map[SPN]bool)
	Q := common.Queue{}
	Q.Enqueue(S)
	V[S] = true

	for !Q.Empty() {
		s := Q.Dequeue().(SPN)
		ch := s.Ch()
		if s.Type() == "sum" {
			mi := -1
			for i, c := range ch {
				if T.ExistsSPN(c) && !V[c] {
					mi = i
				}
			}
			//if mi < 0 {
			//sys.Printf("mi=%d, SPN={\n%v\n}\n", mi, s)
			//} else {
			if mi >= 0 {
				trace[s] = mi
				c := ch[mi]
				Q.Enqueue(c)
				V[c] = true
			}
		}
		for _, c := range ch {
			if !V[c] && c.Type() != "leaf" {
				Q.Enqueue(c)
				V[c] = true
			}
		}
	}

	return trace
}
