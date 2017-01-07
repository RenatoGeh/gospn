package spn

import (
	//"fmt"
	"math"
	"sort"
)

// Sum represents a sum node in an SPN.
type Sum struct {
	Node
	w []float64
	// Store partial derivatives wrt weights.
	pweights map[string][]float64
	// Auto-normalizes on weight updating.
	norm bool
}

// NewSum creates a new Sum node.
func NewSum() *Sum {
	return &Sum{Node: NewNode(), pweights: make(map[string][]float64), norm: true}
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

// AutoNormalize sets whether this sum node should auto normalize on weight update.
func (s *Sum) AutoNormalize(norm bool) { s.norm = norm }

// Soft is a common base for all soft inference methods.
func (s *Sum) Soft(val VarSet, key string) float64 {
	if _lv, ok := s.Stored(key); ok {
		return _lv
	}

	n := len(s.ch)

	vals := make([]float64, n)
	for i := 0; i < n; i++ {
		v, w := s.ch[i].Soft(val, key), math.Log(s.w[i])
		vals[i] = v + w
		//if s.root {
		//fmt.Printf("Root v+w=%f+(log(%f)=%f)=%f\n", v, s.w[i], w, vals[i])
		//}
	}
	sort.Float64s(vals)
	p, r := vals[n-1], 0.0

	for i := 0; i < n-1; i++ {
		r += math.Exp(vals[i] - p)
	}

	r = p + math.Log1p(r)

	//if key == "soft" {
	//fmt.Printf("Sum %f\n", r)
	//}
	s.Store(key, r)
	return r
}

// Value returns the value of this node given an instantiation.
func (s *Sum) Value(val VarSet) float64 {
	return s.Soft(val, "soft")
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
		v, _ := s.Stored(nkey)
		st[nkey] += math.Log1p(math.Exp(math.Log(s.w[i]) + v - st[nkey]))
	}

	for i := 0; i < n; i++ {
		v, _ := s.Stored(nkey)
		u, _ := s.ch[i].Stored(ikey)
		s.pweights[wkey][i] = u + v
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
func (s *Sum) DiscUpdate(eta float64, ds *DiscStorer, wckey, wekey string) {
	n := len(s.ch)
	t, min := 0.0, s.w[0]

	correct, expected := ds.Correct(), ds.Expected()
	for i := 0; i < n; i++ {
		s.w[i] += eta * ((s.pweights[wckey][i] / correct) - (s.pweights[wekey][i] / expected))
		t += s.w[i]
		if s.w[i] < min {
			min = s.w[i]
		}
	}

	if s.norm {
		min = math.Abs(min)
		t += float64(n) * min
		for i := 0; i < n; i++ {
			s.w[i] = (s.w[i] + min) / t
		}
	}

	for i := 0; i < n; i++ {
		s.ch[i].DiscUpdate(eta, ds, wckey, wekey)
	}
}

// ResetDP resets a key on the DP table. If key is nil, resets everything.
func (s *Sum) ResetDP(key string) {
	s.Node.ResetDP(key)
	if key == "" {
		for k := range s.pweights {
			s.pweights[k] = nil
		}
	} else {
		s.pweights[key] = nil
	}
}
