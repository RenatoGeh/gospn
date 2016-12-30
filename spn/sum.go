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
	pweights map[string][]float64
}

// NewSum creates a new Sum node.
func NewSum() *Sum {
	return &Sum{Node: NewNode(), pweights: make(map[string][]float64)}
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
func (s *Sum) Sc() map[int]int {
	if s.sc == nil {
		csc := s.ch[0].Sc()
		for k := range csc {
			s.sc[k] = k
		}
	}
	return s.sc
}

// Bsoft is a common base for all soft inference methods.
func (s *Sum) Bsoft(val VarSet, key string) float64 {
	n := len(s.ch)

	vals := make([]float64, n)
	for i := 0; i < n; i++ {
		v, w := s.ch[i].Bsoft(val, key), math.Log(s.w[i])
		vals[i] = v + w
	}
	sort.Float64s(vals)
	p, r := vals[n-1], 0.0

	for i := 0; i < n-1; i++ {
		r += math.Exp(vals[i] - p)
	}

	r = p + math.Log1p(r)
	s.Store(key, r)
	return r
}

// Value returns the value of this node given an instantiation.
func (s *Sum) Value(val VarSet) float64 {
	return s.Bsoft(val, "soft")
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
func (s *Sum) PWeights(key string) []float64 {
	return s.pweights[key]
}

// Derive recursively derives this node and its children based on the last inference value.
func (s *Sum) Derive(wkey, nkey, ikey string) {
	n := len(s.ch)
	if s.pweights[wkey] == nil {
		s.pweights[wkey] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		st := s.ch[i].Storer()
		st[nkey] += math.Log1p(math.Exp(math.Log(s.w[i]) + s.Stored(nkey) - st[nkey]))
	}

	for i := 0; i < n; i++ {
		s.pweights[wkey][i] = s.ch[i].Stored(ikey) + s.Stored(nkey)
	}

	for i := 0; i < n; i++ {
		s.ch[i].Derive(wkey, nkey, ikey)
	}
}

// GenUpdate generatively updates weights given an eta learning rate.
func (s *Sum) GenUpdate(eta float64, wkey string) {
	n := len(s.ch)
	t := 0.0

	for i := 0; i < n; i++ {
		s.w[i] += eta + math.Exp(s.pweights[wkey][i])
		t += s.w[i]
	}

	// Normalize weights.
	for i := 0; i < n; i++ {
		s.w[i] /= t
	}
}

// Normalize normalizes weights.
func (s *Sum) Normalize() {
	n, t := len(s.ch), 0.0
	for i := 0; i < n; i++ {
		t += s.w[i]
	}
	for i := 0; i < n; i++ {
		s.w[i] /= t
	}
	for i := 0; i < n; i++ {
		s.ch[i].Normalize()
	}
}

// DiscUpdate discriminatively updates weights given an eta learning rate.
func (s *Sum) DiscUpdate(eta float64) {
	//n := len(s.ch)

	//root.RResetDP("disc_correct")
	//root.RResetDP("disc_expected")
	//correct := root.Bsoft(T, "disc_correct")
	//expected := root.Bsoft(X, "disc_expected")

	//for i := 0; i < n; i++ {
	//s.w[i] += eta*(math.Exp(s.pweights[i] - correct) - math.Exp(s.pweights
	//}
}
