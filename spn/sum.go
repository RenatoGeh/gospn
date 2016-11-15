package spn

import (
	//"fmt"
	"math"
	"sort"

	utils "github.com/RenatoGeh/gospn/utils"
)

// Sum represents an SPN sum node.
type Sum struct {
	// Children nodes.
	ch []Node
	// Edge weights.
	w []float64
	// Parent node.
	pa Node
	// Node scope.
	sc []int
}

// NewSum returns an empty Sum node with given parent.
func NewSum() *Sum {
	s := &Sum{}
	s.pa, s.sc = nil, nil
	return s
}

// AddChild adds a child without adding weight. After this function call you must call AddWeight (or call
// AddChildW instead.
func (s *Sum) AddChild(c Node) {
	s.ch = append(s.ch, c)
	c.SetParent(s)
	s.sc = nil
}

// AddWeight adds a new weight to the latest added edge.
func (s *Sum) AddWeight(w float64) {
	s.w = append(s.w, w)
}

// AddChildW adds a new child to this sum node with a weight w.
func (s *Sum) AddChildW(c Node, w float64) {
	s.ch = append(s.ch, c)
	s.w = append(s.w, w)
	c.SetParent(s)
	s.sc = nil
}

// SetParent sets the parent node.
func (s *Sum) SetParent(pa Node) { s.pa = pa }

// Ch returns the set of children nodes.
func (s *Sum) Ch() []Node { return s.ch }

// Pa returns the parent node.
func (s *Sum) Pa() Node { return s.pa }

// Type returns the type of this node: 'sum'.
func (s *Sum) Type() string { return "sum" }

// Sc returns the scope of this node.
func (s *Sum) Sc() []int {
	if s.sc == nil {
		// Since all sum nodes are complete, then all children must have the same scope (we consider
		// the SPN definition by Gens and Domingos).
		copy(s.sc, s.ch[0].Sc())
	}
	return s.sc
}

// Weights returns weights.
func (s *Sum) Weights() []float64 { return s.w }

// Value returns the value of this SPN given a set of valuations.
func (s *Sum) Value(valuation VarSet) float64 {
	n := len(s.ch)
	ch := s.Ch()

	vals := make([]float64, n)
	for i := 0; i < n; i++ {
		v, w := ch[i].Value(valuation), math.Log(s.w[i])
		vals[i] = v + w
	}
	sort.Float64s(vals)
	p, r := vals[n-1], 0.0

	for i := 0; i < n-1; i++ {
		//fmt.Printf("vals[i] = %f vs p = %f\n", vals[i], p)
		r += math.Exp(vals[i] - p)
	}

	r = p + math.Log1p(r)

	//p1 := utils.Log(s.w[0]) + ch[0].Value(valuation)
	//for i := 1; i < n; i++ {
	//p2 := utils.Log(s.w[i]) + ch[i].Value(valuation)
	//l := utils.LogSumPair(utils.AntiLog(p1), utils.AntiLog(p2))
	//p1 = l
	//}

	//for i := 0; i < n; i++ {
	//vch := (s.ch[i]).Value(valuation)
	//v += utils.Log(s.w[i]) * vch
	//fmt.Printf("ch[%d] of type \"%s\" pa \"%s\": %f\n", i, s.ch[i].Type(), "sum", vch)
	//v += s.w[i] * vch
	//}

	//fmt.Printf("Value of sum node: antiln(%f)=%f\n", r, utils.AntiLog(r))
	return r
}

// Max returns the MAP state given a valuation.
func (s *Sum) Max(valuation VarSet) float64 {
	var max = utils.Inf(-1)
	n := len(s.ch)

	for i := 0; i < n; i++ {
		cv := utils.Log(s.w[i]) + s.ch[i].Max(valuation)
		if cv > max {
			max = cv
		}
	}

	return max
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (s *Sum) ArgMax(valuation VarSet) (VarSet, float64) {
	n, max := len(s.ch), utils.Inf(-1)
	var mch Node

	// Since a sum node must be complete, there can never be a leaf adjacent to a sum node, as that
	// would imply that all its children would also have to be leaves with the same scope as each
	// other. Since leaves are univariate distributions, this would mean a clustering over the same
	// variable, which would annul the clustering done in learning and leave us with either a
	// supercluster or the full distribution. And that makes no sense. Therefore all children from a
	// sum node must not be leaves. For this reason we seek only the max edge instead and delegate
	// to its children.
	for i := 0; i < n; i++ {
		ch := s.ch[i]
		m := utils.Log(s.w[i]) + ch.Max(valuation)
		if m > max {
			max, mch = m, ch
		}
	}

	set, _ := mch.ArgMax(valuation)
	//fmt.Printf("Sum: %f, %v\n", utils.AntiLog(max), set)
	return set, max
}

// Data for a sum node is the number of children and weights.
func (s *Sum) Data() (int, []float64) {
	return len(s.ch), s.w
}