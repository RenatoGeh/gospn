package spn

// Sum represents an SPN sum node.
type Sum struct {
	// Children nodes.
	ch []Node
	// Edge weights.
	w []float64
	// Parent node.
	pa Node
	// Node scope.
	sc []int
}

// NewSum returns an empty Sum node with given parent.
func NewSum() *Sum {
	s := &Sum{}
	s.pa, s.sc = nil, nil
	return s
}

// Adds a child without adding weight. After this function call you must call AddWeight (or call
// AddChildW instead.
func (s *Sum) AddChild(c Node) {
	s.ch = append(s.ch, c)
	c.SetParent(s)
	s.sc = nil
}

func (s *Sum) AddWeight(w float64) {
	s.w = append(s.w, w)
}

// AddChild adds a new child to this sum node with a weight w.
func (s *Sum) AddChildW(c Node, w float64) {
	s.ch = append(s.ch, c)
	s.w = append(s.w, w)
	c.SetParent(s)
	s.sc = nil
}

// Sets the parent node.
func (s *Sum) SetParent(pa Node) { s.pa = pa }

// Ch returns the set of children nodes.
func (s *Sum) Ch() []Node { return s.ch }

// Pa returns the parent node.
func (s *Sum) Pa() Node { return s.pa }

// Type returns the type of this node: 'sum'.
func (s *Sum) Type() string { return "sum" }

// Sc returns the scope of this node.
func (s *Sum) Sc() []int {
	if s.sc == nil {
		// Since all sum nodes are complete, then all children must have the same scope (we consider
		// the SPN definition by Gens and Domingos).
		copy(s.sc, s.ch[0].Sc())
	}
	return s.sc
}

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

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (s *Sum) ArgMax(valuation VarSet) (VarSet, float64) {
	n, max := len(s.ch), 0.0
	var mch Node = nil

	// Since a sum node must be complete, there can never be a leaf adjacent to a sum node, as that
	// would imply that all its children would also have to be leaves with the same scope as each
	// other. Since leaves are univariate distributions, this would mean a clustering over the same
	// variable, which would annul the clustering done in learning and leave us with either a
	// supercluster or the full distribution. And that makes no sense. Therefore all children from a
	// sum node must not be leaves. For this reason we seek only the max edge instead and delegate
	// to its children.
	for i := 0; i < n; i++ {
		ch := s.ch[i]
		m := s.w[i] * ch.Max(valuation)
		if m > max {
			max = m
			mch = ch
		}
	}

	return mch.ArgMax(valuation)
}
