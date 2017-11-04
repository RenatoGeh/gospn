package spn

import (
	"github.com/RenatoGeh/gospn/common"
)

// Topological sorting algorithms.

func visit(S SPN, Q *common.Queue, V map[SPN]bool) {
	if V[S] {
		return
	}
	V[S] = true
	ch := S.Ch()
	for _, c := range ch {
		visit(c, Q, V)
	}
	Q.Enqueue(S)
}

func TopSortTarjanRec(G SPN) *common.Queue {
	Q := &common.Queue{}
	V := make(map[SPN]bool)
	visit(G, Q, V)
	return Q
}

// TopSortTarjan returns the topological sorting of a graph G. It follows the version described in
// [Tarjan, 1974] but in a non-recursive fashion.
func TopSortTarjan(G SPN) *common.Queue {
	Q := &common.Queue{}
	S := common.Stack{}
	P := make(map[SPN]bool)
	V := make(map[SPN]bool)
	S.Push(G)
	V[G] = true
	for !S.Empty() {
		u := S.Pop().(SPN)
		if P[u] {
			Q.Enqueue(u)
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
	return Q
}
