package spn

import (
	"github.com/RenatoGeh/gospn/common"
)

// Graph traversal/search algorithms.

func searchFunc(G SPN, f func(SPN) bool, C common.Collection) {
	V := make(map[SPN]bool)
	C.Give(G)
	V[G] = true
	for !C.Empty() {
		u := C.Take().(SPN)
		if !f(u) {
			return
		}
		ch := u.Ch()
		for _, c := range ch {
			if !V[c] {
				C.Give(c)
				V[c] = true
			}
		}
	}
}

// BreadthFirst applies a function f to each node of the graph G. The graph traversal is node using
// a breadth-first search approach. If f returns false, then the search ends. Else, it continues.
func BreadthFirst(G SPN, f func(SPN) bool) { searchFunc(G, f, &common.Queue{}) }

// DepthFirst applies a function f to each node of the graph G. The graph traversal is node using
// a depth-first search approach. If f returns false, then the search ends. Else, it continues.
func DepthFirst(G SPN, f func(SPN) bool) { searchFunc(G, f, &common.Stack{}) }

// CountNodes counts the number of nodes, returning the number of sum, product and leaf nodes in
// this order.
func CountNodes(G SPN) (int, int, int) {
	var s, p, l int
	BreadthFirst(G, func(S SPN) bool {
		if t := S.Type(); t == "sum" {
			s++
		} else if t == "product" {
			p++
		} else {
			l++
		}
		return true
	})
	return s, p, l
}
