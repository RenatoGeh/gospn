// Package spn contains the structure of an SPN.
package spn

// Sum represents an SPN sum node.
type Sum struct {
	// Children nodes.
	ch []Node
	// Edge weights.
	w []float32
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
func (s *Sum) AddChild(c Node, w float32) {
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
func (s *Sum) Value(valuation VarSet) float32 {
	var v float32 = 0
	n := len(s.ch)

	for i := 0; i < n; i++ {
		v += s.w[i] * (s.ch[i]).Value(valuation)
	}

	return v
}
