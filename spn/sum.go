package spn

import (
	"math"
	"sort"
)

// Sum represents a sum node in an SPN.
type Sum struct {
	Node
	w []float64
	// Store partial derivatives wrt weights.
	pweights []float64
}

// NewSum creates a new Sum node.
func NewSum() *Sum {
	return &Sum{}
}

// AddWeight adds a new weight to the sum node.
func (s *Sum) AddWeight(w float64) {
	s.w = append(s.w, w)
	s.pweights = append(s.pweights, 0)
}

// AddChildW adds a new child to this sum node with a weight w.
func (s *Sum) AddChildW(c SPN, w float64) {
	s.AddChild(c)
	s.AddWeight(w)
}

// Sc returns the scope of this node.
func (s *Sum) Sc() []int {
	if s.sc == nil {
		copy(s.sc, s.ch[0].Sc())
	}
	return s.sc
}

// Value returns the value of this node given an instantiation.
func (s *Sum) Value(val VarSet) float64 {
	n := len(s.ch)

	vals := make([]float64, n)
	for i := 0; i < n; i++ {
		v, w := s.ch[i].Value(val), math.Log(s.w[i])
		vals[i] = v + w
	}
	sort.Float64s(vals)
	p, r := vals[n-1], 0.0

	for i := 0; i < n-1; i++ {
		r += math.Exp(vals[i] - p)
	}

	r = p + math.Log1p(r)
	s.s = r

	return r
}

// Max returns the MAP value of this node given an evidence.
func (s *Sum) Max(val VarSet) float64 {
	max := math.Inf(-1)
	n := len(s.ch)

	for i := 0; i < n; i++ {
		cv := math.Log(s.w[i]) + s.ch[i].Max(val)
		if cv > max {
			max = cv
		}
	}

	return max
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (s *Sum) ArgMax(val VarSet) (VarSet, float64) {
	n, max := len(s.ch), math.Inf(-1)
	var mch SPN

	for i := 0; i < n; i++ {
		ch := s.ch[i]
		// Note to future self: use DP to avoid recomputations.
		m := math.Log(s.w[i]) + ch.Max(val)
		if m > max {
			max, mch = m, ch
		}
	}

	amax, mval := mch.ArgMax(val)
	return amax, mval
}

// Type returns the type of this node.
func (s *Sum) Type() string { return "sum" }

// Weights returns weights if sum product. Returns nil otherwise.
func (s *Sum) Weights() []float64 {
	return s.w
}

// PWeights returns the partial derivatives wrt this node's weights.
func (s *Sum) PWeights() []float64 {
	return s.pweights
}

// Derive recursively derives this node and its children based on the last inference value.
func (s *Sum) Derive() {
	n := len(s.ch)

	for i := 0; i < n; i++ {
		da := s.ch[i].DrvtAddr()
		*da += math.Log1p(math.Exp(math.Log(s.w[i]) + s.pnode - *da))
	}

	for i := 0; i < n; i++ {
		s.pweights[i] = s.ch[i].Stored() + s.pnode
	}

	for i := 0; i < n; i++ {
		s.ch[i].Derive()
	}
}

// GenUpdate generatively updates weights given an eta learning rate.
func (s *Sum) GenUpdate(eta float64) {
	n := len(s.ch)

	for i := 0; i < n; i++ {
		s.w[i] += eta + math.Exp(s.pweights[i])
	}
}
