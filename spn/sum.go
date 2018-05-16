package spn

import (
	"github.com/RenatoGeh/gospn/learn/parameters"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/utils"
	"math"
)

// Sum represents a sum node in an SPN.
type Sum struct {
	Node
	w []float64
}

// NewSum creates a new Sum node.
func NewSum() *Sum {
	return &Sum{}
}

// AddWeight adds a new weight to the sum node.
func (s *Sum) AddWeight(w float64) {
	s.w = append(s.w, w)
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

	l := s.Compute(vals)
	sys.Printf("Sum value: %f = ln(%f)\n", l, math.Exp(l))
	return l
}

// Compute returns the soft value of this node's type given the children's values.
func (s *Sum) Compute(cv []float64) float64 {
	return utils.LogSumExp(cv)
}

// ComputeHard returns the soft value using hard weights (unit weights). Expects only children
// values, and not weighted child value. Parameter scount is the smooth sum count constant.
func (s *Sum) ComputeHard(cv []float64, scount float64) float64 {
	imax, max := -1, -1.0
	for i := range s.ch {
		w := s.w[i]
		if imax < 0 || max < w {
			imax, max = i, w
		}
	}
	if imax < 0 {
		return utils.LogZero
	}
	var v float64
	for i := range cv {
		v += s.w[i] * math.Exp(cv[i]-max)
	}
	return math.Log(v/scount) + max
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
		m := math.Log(s.w[i]) + ch.Max(val)
		if m > max {
			max, mch = m, ch
		}
	}

	amax, _ := mch.ArgMax(val)
	return amax, max
}

// Type returns the type of this node.
func (s *Sum) Type() string { return "sum" }

// Weights returns weights if sum node. Returns nil otherwise.
func (s *Sum) Weights() []float64 {
	return s.w
}

// AddChild adds a child.
func (s *Sum) AddChild(c SPN) {
	s.ch = append(s.ch, c)
}

// Parameters returns the parameters of this object. If no bound parameter is found, binds default
// parameter values and returns.
func (s *Sum) Parameters() *parameters.P {
	p, e := parameters.Retrieve(s)
	if !e {
		p = parameters.Default()
		parameters.Bind(s, p)
	}
	return p
}
