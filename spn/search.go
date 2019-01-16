package spn

import (
	"github.com/RenatoGeh/gospn/common"
)

// Graph traversal/search algorithms.

// f return value can either be 0: continue searching; 1: stop search; or -1: skip this branch.
func searchFunc(G SPN, f func(SPN) int, C common.Collection) {
	V := make(map[SPN]bool)
	C.Give(G)
	V[G] = true
	for !C.Empty() {
		u := C.Take().(SPN)
		r := f(u)
		if r > 0 {
			return
		}
		ch := u.Ch()
		for _, c := range ch {
			if !V[c] && r == 0 {
				C.Give(c)
				V[c] = true
			}
		}
	}
}

// BreadthFirst applies a function f to each node of the graph G. The graph traversal is node using
// a breadth-first search approach. If f returns false, then the search ends. Else, it continues.
func BreadthFirst(G SPN, f func(SPN) int) { searchFunc(G, f, &common.Queue{}) }

// DepthFirst applies a function f to each node of the graph G. The graph traversal is node using
// a depth-first search approach. If f returns false, then the search ends. Else, it continues.
func DepthFirst(G SPN, f func(SPN) int) { searchFunc(G, f, &common.Stack{}) }

// CountNodes counts the number of nodes, returning the number of sum, product and leaf nodes in
// this order.
func CountNodes(G SPN) (int, int, int) {
	var s, p, l int
	BreadthFirst(G, func(S SPN) int {
		if t := S.Type(); t == "sum" {
			s++
		} else if t == "product" {
			p++
		} else {
			l++
		}
		return 0
	})
	return s, p, l
}
