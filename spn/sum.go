package spn

import (
	//"fmt"
	"github.com/RenatoGeh/gospn/common"
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
	// L2 regularization weight penalty.
	l float64
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
	if _lv, ok := s.Stored(key); ok && s.stores {
		return _lv
	}

	v, n := 0.0, len(s.ch)
	for i := 0; i < n; i++ {
		p := s.ch[i].Soft(val, key)
		v += s.w[i] * p
		//if s.root {
		//fmt.Printf("Root %f * %f = %f\n", s.w[i], p, s.w[i]*p)
		//}
	}

	s.Store(key, v)
	return v
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
		cv := s.w[i] * s.ch[i].Max(val)
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
		m := s.w[i] * ch.Max(val)
		if m > max {
			max, mch = m, ch
		}
	}

	amax, _ := mch.ArgMax(val)
	return amax, max
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

// Derive derives this node only.
func (s *Sum) Derive(wkey, nkey, ikey string, mode InfType) int {
	n := len(s.ch)
	if s.pweights[wkey] == nil {
		s.pweights[wkey] = make([]float64, n)
	}

	if mode == SOFT {
		v, _ := s.Stored(nkey)
		for i := 0; i < n; i++ {
			st := s.ch[i].Storer()
			st[nkey] += s.w[i] * v
			s.pweights[wkey][i] = st[ikey] * v
		}

		//for i := 0; i < n; i++ {
		//s.ch[i].Derive(wkey, nkey, ikey)
		//}

		return -1
	}

	max, imax := math.Inf(-1), -1
	for i := 0; i < n; i++ {
		v, _ := s.ch[i].Stored(ikey)
		if v > max {
			max, imax = v, i
		}
	}
	s.pweights[wkey][imax]++

	return imax
}

// LSoft is Soft in logspace.
func (s *Sum) LSoft(val VarSet, key string) float64 {
	if _lv, ok := s.Stored(key); ok {
		return _lv
	}

	n := len(s.ch)

	vals := make([]float64, n)
	for i := 0; i < n; i++ {
		v, w := s.ch[i].LSoft(val, key), math.Log(s.w[i])
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

// GenUpdate generatively updates weights given an eta learning rate.
func (s *Sum) GenUpdate(eta float64, wkey string) {
	n := len(s.ch)
	t := 0.0

	for i := 0; i < n; i++ {
		s.w[i] += eta * s.pweights[wkey][i]
		t += s.w[i]
	}

	// Normalize weights.
	for i := 0; i < n; i++ {
		s.w[i] /= t
	}
}

// NormalizeThis normalizes only this node's weights.
func (s *Sum) NormalizeThis() {
	n := len(s.w)
	min := s.w[0]
	for i := 1; i < n; i++ {
		if s.w[i] < min {
			min = s.w[i]
		}
	}
	if min < 0 {
		min = math.Abs(min)
		for i := 0; i < n; i++ {
			s.w[i] += min
		}
	}
	var norm float64
	for i := 0; i < n; i++ {
		norm += s.w[i]
	}
	for i := 0; i < n; i++ {
		s.w[i] = s.w[i] / norm
	}
}

// Normalize normalizes the SPN's weights.
func (s *Sum) Normalize() {
	s.NormalizeThis()
	n := len(s.ch)

	for i := 0; i < n; i++ {
		s.ch[i].Normalize()
	}
}

// DiscUpdate discriminatively updates weights given an eta learning rate.
func (s *Sum) DiscUpdate(eta float64, ds *DiscStorer, wckey, wekey string, mode InfType) {
	n := len(s.ch)

	if mode == SOFT {
		//ds.ResetSPN("")
		correct, expected := ds.Correct(), ds.Expected()
		//ds.DeriveExpected(s)
		//ds.DeriveCorrect(s)
		for i := 0; i < n; i++ {
			cc := s.pweights[wckey][i] / (correct + 0.000001)
			ce := s.pweights[wekey][i] / (expected + 0.000001)
			//cpn, _ := s.Stored("cpnode")
			//epn, _ := s.Stored("epnode")
			//fmt.Printf("Sum -> cpnode = %.10f, epnode = %.10f\n", cpn, epn)
			//fmt.Printf("Sum -> cc = %.10f / (%.10f + 0.000001) = %.10f\n", s.pweights[wckey][i], correct, cc)
			//fmt.Printf("Sum -> ce = %.10f / (%.10f + 0.000001) = %.10f\n", s.pweights[wekey][i], expected, ce)
			//fmt.Printf("Sum -> w: %f += (cc: %.10f - ce: %.10f = %.10f)\n", s.w[i], cc, ce, cc-ce)
			s.w[i] += eta * (cc - ce - 2*s.l*s.w[i])
			//s.w[i] += eta * ((s.pweights[wckey][i] / correct) - (s.pweights[wekey][i] / expected))
		}
	} else {
		ds.ResetSPN("")
		ds.DeriveExpected(s)
		ds.DeriveCorrect(s)
		if s.pweights[wckey] == nil {
			s.pweights[wckey] = make([]float64, n)
		}
		if s.pweights[wekey] == nil {
			s.pweights[wekey] = make([]float64, n)
		}
		for i := 0; i < n; i++ {
			//fmt.Printf("[%d-%d]: (%s)->%v, (%s)->%v\n", i, n, wckey, s.pweights[wckey], wekey, s.pweights[wekey])
			c, e := s.pweights[wckey][i], s.pweights[wekey][i]
			s.w[i] += eta * ((c-e)/(s.w[i]+0.01) - 2*s.l*s.w[i])
			if c > 0 || e > 0 {
				s.ch[i].DiscUpdate(eta, ds, wckey, wekey, mode)
			}
		}
	}

	if s.norm {
		s.NormalizeThis()
	}

	if mode == SOFT {
		for i := 0; i < n; i++ {
			s.ch[i].DiscUpdate(eta, ds, wckey, wekey, mode)
		}
	}
}

// RResetDP recursively ResetDPs all children.
func (s *Sum) RResetDP(key string) {
	n := len(s.ch)

	s.ResetDP(key)
	for i := 0; i < n; i++ {
		s.ch[i].RResetDP(key)
	}
}

// ResetDP resets a key on the DP table. If key is nil, resets everything.
func (s *Sum) ResetDP(key string) {
	s.Node.ResetDP(key)
	if key == "" {
		s.pweights = make(map[string][]float64)
	} else {
		delete(s.pweights, key)
	}
}

// RootDerive derives all nodes in a BFS fashion.
func (s *Sum) RootDerive(wkey, nkey, ikey string, mode InfType) {
	q := common.Queue{}

	q.Enqueue(s)

	for !q.Empty() {
		t := q.Dequeue().(SPN)
		ch := t.Ch()

		r := t.Derive(wkey, nkey, ikey, mode)

		if ch != nil && r != 0 {
			if r < 0 {
				n := len(ch)
				for i := 0; i < n; i++ {
					q.Enqueue(ch[i])
				}
			} else {
				q.Enqueue(ch[r])
			}
		}
	}
}

// L2 regularization weight penalty.
func (s *Sum) L2() float64 { return s.l }

// SetL2 changes the L2 regularization weight penalty throughout all SPN.
func (s *Sum) SetL2(l float64) {
	n := len(s.ch)
	s.l = l
	for i := 0; i < n; i++ {
		s.ch[i].SetL2(l)
	}
}
