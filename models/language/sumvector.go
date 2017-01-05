package language

import (
	"github.com/RenatoGeh/gospn/spn"
	"math"
)

// SumVector represents the H layer of the structure described in LMSPN. A layer
// H_11,...,H_1D,H_21,...,H_2D,...,H_N1,...,H_ND of sum vectors is a compression of the N
// K-dimensional vectors into a single continuous D-dimensional feature vector.
type SumVector struct {
	spn.Node
	// Weights
	w []float64
	// Length of w
	n int
	// Store partial deriatives wrt weights.
	cpw []float64
	epw []float64
}

// NewSumVector creates a new SumVector node.
func NewSumVector(waddr, cpweight, epweight []float64) *SumVector {
	return &SumVector{spn.NewNode(), waddr, len(waddr), cpweight, epweight}
}

// Soft is a common base for all soft inference methods.
func (s *SumVector) Soft(val spn.VarSet, key string) float64 {
	_lv := s.Stored(key)
	if _lv >= 0 {
		return _lv
	}

	// By definition, a SumVector contains only one child: a Vector node.
	// Note to self: don't forget in this case we are using VarSet as a slice (and as such they are
	// (not really) ordered by index).
	ch := s.Ch()
	v := math.Log(s.w[int(ch[0].Value(val))])
	s.Store(key, v)
	return v
}

// Value returns the value of this node given an instantiation.
func (s *SumVector) Value(val spn.VarSet) float64 {
	return s.Soft(val, "soft")
}

// Max returns the MAP value of this node given an evidence.
func (s *SumVector) Max(val spn.VarSet) float64 {
	return s.Soft(val, "max")
}

// Type returns the type of this node.
func (s *SumVector) Type() string { return "sum_vector" }

// Derive recursively derives this node and its children based on the last inference value.
func (s *SumVector) Derive(wkey, nkey, ikey string) {
	ch := s.Ch()[0]
	var pweight []float64

	if wkey == "cpweight" {
		pweight = s.cpw
	} else {
		pweight = s.epw
	}

	pweight[int(ch.Stored(ikey))] = s.Stored(nkey)
}

// GenUpdate generatively updates weights given an eta learning rate.
func (s *SumVector) GenUpdate(eta float64, wkey string) {
	k := int(s.Ch()[0].Stored("soft"))
	s.w[k] += eta + math.Exp(s.epw[k])

	// Normalize
	t := 0.0
	for i := 0; i < s.n; i++ {
		t += s.w[i]
	}
	for i := 0; i < s.n; i++ {
		s.w[i] /= t
	}
}

// DiscUpdate discriminatively updates weights given an eta learning rate.
func (s *SumVector) DiscUpdate(eta, correct, expected float64, wckey, wekey string) {
	k := int(s.Ch()[0].Stored("soft"))
	s.w[k] += eta * (math.Exp(s.cpw[k]-correct) - math.Exp(s.epw[k]-expected))

	// Normalize
	t := 0.0
	for i := 0; i < s.n; i++ {
		t += s.w[i]
	}
	for i := 0; i < s.n; i++ {
		s.w[i] /= t
	}
}
