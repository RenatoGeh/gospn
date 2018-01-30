package spn

import (
	"github.com/RenatoGeh/gospn/common"
)

// Topological sorting algorithms.

func visit(S SPN, C common.Collection, V map[SPN]bool) {
	if V[S] {
		return
	}
	V[S] = true
	ch := S.Ch()
	for _, c := range ch {
		visit(c, C, V)
	}
	C.Give(S)
}

// TopSortTarjan returns the topological sorting of a graph G. It follows the version described in
// [Tarjan, 1974]. The argument C indicates how the topological sorting should be ordered (it C is
// a queue, the function returns an inversed topological sort (dependency ordering); if C is a
// stack, the function returns the topological sorting).
func TopSortTarjanRec(G SPN, C common.Collection) common.Collection {
	if C == nil {
		C = &common.Queue{}
	}
	V := make(map[SPN]bool)
	visit(G, C, V)
	return C
}

// TopSortTarjan returns the topological sorting of a graph G. It follows the version described in
// [Tarjan, 1974] but in a non-recursive fashion. The argument C indicates how the topological
// sorting should be ordered (it C is a queue, the function returns an inversed topological sort
// (dependency ordering); if C is a stack, the function returns the topological sorting).
func TopSortTarjan(G SPN, C common.Collection) common.Collection {
	if C == nil {
		C = &common.Queue{}
	}
	S := common.Stack{}
	P := make(map[SPN]bool)
	V := make(map[SPN]bool)
	S.Push(G)
	V[G] = true
	for !S.Empty() {
		u := S.Pop().(SPN)
		if P[u] {
			C.Give(u)
			continue
		}
		S.Push(u)
		P[u] = true
		ch := u.Ch()
		for _, c := range ch {
			if !V[c] {
				S.Push(c)
				V[c] = true
			}
		}
	}
	return C
}

// TopSortTarjanFunc traverses the graph G following TopSortTarjan, but at each step we also
// perform a function f. Useful for computing inline operations at each topological sort insertion.
// If f returns false, then the topological sort halts immediately, preserving the Queue at the
// moment of falsehood. The argument C indicates how the topological sorting should be ordered (it
// C is a queue, the function returns an inversed topological sort (dependency ordering); if C is a
// stack, the function returns the topological sorting).
func TopSortTarjanFunc(G SPN, C common.Collection, f func(SPN) bool) common.Collection {
	if C == nil {
		C = &common.Queue{}
	}
	S := common.Stack{}
	P := make(map[SPN]bool)
	V := make(map[SPN]bool)
	S.Push(G)
	V[G] = true
	for !S.Empty() {
		u := S.Pop().(SPN)
		if P[u] {
			C.Give(u)
			if !f(u) {
				return C
			}
			continue
		}
		S.Push(u)
		P[u] = true
		ch := u.Ch()
		for _, c := range ch {
			if !V[c] {
				S.Push(c)
				V[c] = true
			}
		}
	}
	return C
}

// TopSortDFS finds a topological sort using a DFS. The argument C indicates how the topological
// sorting should be ordered (it C is a queue, the function returns an inversed topological sort
// (dependency ordering); if C is a stack, the function returns the topological sorting).
func TopSortDFS(G SPN, C common.Collection) common.Collection {
	S := common.Stack{}
	V := make(map[SPN]bool)
	if C == nil {
		C = &common.Queue{}
	}

	S.Push(G)
	V[G] = true

	for !S.Empty() {
		s := S.Pop().(SPN)
		ch := s.Ch()
		if len(ch) == 0 {
			C.Give(s)
		} else {
			var m []SPN
			for _, c := range ch {
				if !V[c] {
					m = append(m, c)
				}
			}
			if len(m) == 0 {
				C.Give(s)
			} else {
				S.Push(s)
				for _, c := range m {
					S.Push(c)
					V[c] = true
				}
			}
		}
	}

	return C
}
