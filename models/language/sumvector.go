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
	w []*float64
	// Length of w
	n int
	// Store partial deriatives wrt weights.
	pweights map[string][]float64
}

// NewSumVector creates a new SumVector node.
func NewSumVector(waddr []*float64) *SumVector {
	return &SumVector{spn.NewNode(), waddr, len(waddr), make(map[string][]float64)}
}

// Soft is a common base for all soft inference methods.
func (s *SumVector) Soft(val spn.VarSet, key string) float64 {
	// By definition, a SumVector contains only one child: a Vector node.
	// Note to self: don't forget in this case we are using VarSet as a slice (and as such they are
	// (not really) ordered by index).
	ch := s.Ch()
	v := math.Log(*s.w[int(ch[0].Value(val))])
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
	if s.pweights[wkey] == nil {
		s.pweights[wkey] = make([]float64, s.n)
	}

}

// ResetDP resets a key on the DP table. If key is nil, resets everything.
func (s *SumVector) ResetDP(key string) {
	s.Node.ResetDP(key)
	if key == "" {
		for k := range s.pweights {
			s.pweights[k] = nil
		}
	} else {
		s.pweights[key] = nil
	}
}
