package dennis

import (
	"github.com/RenatoGeh/gospn/learn"
	"math"
)

type scope map[int]learn.Variable
type scopeSlice []scope

func (S scopeSlice) contains(s scope) bool {
	for _, r := range S {
		if s.equal(r) {
			return true
		}
	}
	return false
}

// subsetOf returns whether scope s is a subset of scope t. The function is well defined for s
// empty (returns true).
func (s scope) subsetOf(t scope) bool {
	for q, _ := range s {
		if _, e := t[q]; !e {
			return false
		}
	}
	return true
}

// lenUnion returns the length of the union of sets s and t.
func (s scope) lenUnion(t scope) int {
	m := len(t)
	n := len(s)
	var l int
	var T, S scope
	if m > n {
		l, T, S = m, s, t
	} else {
		l, T, S = n, t, s
	}
	for k, _ := range T {
		if _, e := S[k]; !e {
			l++
		}
	}
	return l
}

// lenIntersect returns the length of the intersection of sets s and t.
func (s scope) lenIntersect(t scope) int {
	var n int
	for k, _ := range t {
		if _, e := s[k]; e {
			n++
		}
	}
	return n
}

// minus returns the set minus result of s and t.
func (s scope) minus(t scope) scope {
	r := make(scope)
	for k, v := range s {
		if _, e := t[k]; !e {
			r[k] = v
		}
	}
	return r
}

// equals returns whether the two scopes are identical.
func (s scope) equal(t scope) bool {
	n, m := len(s), len(t)
	if n != m {
		return false
	}
	var S, T scope
	if n > m {
		S, T = s, t
	} else {
		S, T = t, s
	}
	for k, _ := range S {
		if _, e := T[k]; !e {
			return false
		}
	}
	return true
}

// similarTo returns whether scope s is "similar" to scope "t". We next define similarity according
// to Aaron Dennis and Dan Ventura in
//  Learning the Architecture of Sum-Product Networks Using Clustering on Variables
// Definition. A set S is "similar" to another set T if
//  (p = |S union T|/|S intersect T|) > 1-epsilon
// where 1 > epsilon >= 0. Usually epsilon is "small enough".
// This definition follows from the fact that p=1 if and only if S=T.
// In this function, we let argument d=1-epsilon.
func (s scope) similarTo(r, t scope, d float64) bool {
	p1 := float64(s.lenIntersect(r)) / float64(s.lenUnion(r))
	p2 := float64(s.lenIntersect(t)) / float64(s.lenUnion(t))
	return math.Max(p1, p2) > d
}
