// Package spn contains the structure of an SPN.
package spn

// Sum represents an SPN sum node.
type Sum struct {
	// Children nodes.
	ch []Node
	// Edge weights.
	w []float64
	// Parent node.
	pa Node
}

// NewSum returns an empty Sum node with given parent.
func NewSum(pa Node) *Sum {
	s := &Sum{}
	s.pa = pa
	return s
}

// AddChild adds a new child to this sum node with a weight w.
func (s *Sum) AddChild(c Node, w float64) {
	s.ch = append(s.ch, c)
	s.w = append(s.w, w)
}

// Ch returns the set of children nodes.
func (s *Sum) Ch() []Node { return s.ch }

// Pa returns the parent node.
func (s *Sum) Pa() Node { return s.pa }

// Type returns the type of this node: 'sum'.
func (s *Sum) Type() string { return "sum" }

// Value returns the value of this SPN given a set of valuations.
func (s *Sum) Value(valuation VarSet) float64 {
	var v float64 = 0
	n := len(s.ch)

	for i := 0; i < n; i++ {
		v += s.w[i] * (s.ch[i]).Value(valuation)
	}

	return v
}

// Max returns the MAP state given a valuation.
func (s *Sum) Max(valuation VarSet) float64 {
	var max float64 = 0
	n := len(s.ch)

	for i := 0; i < n; i++ {
		cv := s.ch[i].Max(valuation)
		if cv > max {
			max = cv
		}
	}

	return max
}
