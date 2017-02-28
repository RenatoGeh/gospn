package language

import (
	//"fmt"
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
	// L2 regularization weight penalty.
	l float64
}

// NewSumVector creates a new SumVector node.
func NewSumVector(waddr []float64) *SumVector {
	n := len(waddr)
	return &SumVector{spn.NewNode(), waddr, n, make([]float64, n), make([]float64, n), 0}
}

// Soft is a common base for all soft inference methods.
func (s *SumVector) Soft(val spn.VarSet, key string) float64 {
	if _lv, ok := s.Stored(key); ok && s.Stores() {
		return _lv
	}
	// By definition, a SumVector contains only one child: a Vector node.
	// Note to self: don't forget in this case we are using VarSet as a slice (and as such they are
	// (not really) ordered by index).
	ch := s.Ch()
	v := s.w[int(ch[0].Soft(val, key))]
	//fmt.Printf("SumVector: %.10f\n", v)

	//if key == "soft" {
	//fmt.Printf("SumVector (%p) weights (k=%d):\n", s, int(ch[0].Soft(val, key)))
	//for i := 0; i < len(s.w); i++ {
	//fmt.Printf("w[%d]=%f ", i, s.w[i])
	//}

	//fmt.Printf("SumVector %f\n", v)
	//}

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

// NormalizeThis normalizes only this node's weights.
func (s *SumVector) NormalizeThis() {
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
			s.w[i] += 2 * min
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
func (s *SumVector) Normalize() {
	s.NormalizeThis()
}

// Derive derives this node only.
func (s *SumVector) Derive(wkey, nkey, ikey string, mode spn.InfType) int {
	ch := s.Ch()[0]
	var pweight []float64

	if wkey == "cpweight" {
		pweight = s.cpw
	} else {
		pweight = s.epw
	}

	if mode == spn.SOFT {
		v, _ := ch.Stored(ikey)
		u, _ := s.Stored(nkey)
		k, n := int(v), len(pweight)
		pweight[k] = u * s.w[k]
		for i := 0; i < n; i++ {
			if i != k {
				pweight[i] = 0.0
			}
		}
	} else {
		v, _ := ch.Stored(ikey)
		pweight[int(v)]++
	}

	return 0
}

// GenUpdate generatively updates weights given an eta learning rate.
func (s *SumVector) GenUpdate(eta float64, wkey string) {
	v, _ := s.Ch()[0].Stored("correct")
	k := int(v)
	s.w[k] += eta * s.epw[k]

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
func (s *SumVector) DiscUpdate(eta float64, ds *spn.DiscStorer, wckey, wekey string, mode spn.InfType) {
	if v, _ := s.Stored("visited"); v == 0 {
		s.Store("visited", 1)
	} else {
		return
	}

	n := s.n
	if mode == spn.SOFT {
		//ds.ResetSPN("")
		correct, expected := ds.Correct(), ds.Expected()
		for i := 0; i < n; i++ {
			//ds.DeriveExpected(s)
			//ds.DeriveCorrect(s)
			//v, _ := s.Ch()[0].Stored("correct")
			//k := int(v)
			cc := s.cpw[i] / correct
			ce := s.epw[i] / expected
			//if s.w[k] < 0 || s.epw[k] >= expected {
			//fmt.Printf("s.epw: %.10f expected: %.10f\n", s.epw[k], expected)
			//fmt.Printf("s.cpw: %.10f correct: %.10f\n", s.cpw[k], correct)
			//fmt.Printf("s.w[k]: %.10f\n", s.w[k])
			//}
			//cpn, _ := s.Stored("cpnode")
			//epn, _ := s.Stored("epnode")
			//fmt.Printf("SumVector -> cpnode = %.10f, epnode = %.10f\n", cpn, epn)
			//fmt.Printf("SumVector -> cc = %.10f / (%.10f + 0.000001) = %.10f\n", s.cpw[k], correct, cc)
			//fmt.Printf("SumVector -> ce = %.10f / (%.10f + 0.000001) = %.10f\n", s.epw[k], expected, ce)
			//fmt.Printf("SumVector -> w: %f += (cc: %.10f - ce: %.10f = %.10f)\n", s.w[k], cc, ce, cc-ce)
			//if cc-ce != 0 {
			//fmt.Printf("%s -> dw[%d] = %.2f * (%.5f / %.5f - %.5f / %.5f) = %.2f * (%.5f - %.5f) = "+
			//"%.2f * %.5f = %.5f\n", s.ID(), i, eta, s.cpw[i], correct, s.epw[i], expected, eta, cc,
			//ce, eta, cc-ce, eta*(cc-ce))
			//}
			if s.l == 0 {
				s.w[i] += eta * (cc - ce)
			} else {
				s.w[i] += eta * (cc - ce - 2*s.l*s.w[i])
			}
		}

		// Normalize
		//s.NormalizeThis()
		return
	}
	//ds.ResetSPN("")
	//ds.DeriveExpected(s)
	//ds.DeriveCorrect(s)
	for i := 0; i < n; i++ {
		s.w[i] += eta * ((s.cpw[i]-s.epw[i])/(s.w[i]+0.01) - 2*s.l*s.w[i])
	}

	// Normalize
	s.NormalizeThis()
}

// RResetDP recursively ResetDPs all children.
func (s *SumVector) RResetDP(key string) {
	s.Ch()[0].ResetDP(key)
	s.ResetDP(key)
}

// ResetDP resets a key on the DP table. If key is nil, resets everything.
func (s *SumVector) ResetDP(key string) {
	s.Node.ResetDP(key)
	if key == "" {
		// Compiler will optimize to memclr (as of gc 1.5+).
		for i := range s.epw {
			s.epw[i] = 0.0
		}
		// Compiler will optimize to memclr (as of gc 1.5+).
		for i := range s.epw {
			s.cpw[i] = 0.0
		}
	}
}

// L2 regularization weight penalty.
func (s *SumVector) L2() float64 { return s.l }

// SetL2 changes the L2 regularization weight penalty throughout all SPN.
func (s *SumVector) SetL2(l float64) {
	ch := s.Ch()
	n := len(ch)
	s.l = l
	for i := 0; i < n; i++ {
		ch[i].SetL2(l)
	}
}
