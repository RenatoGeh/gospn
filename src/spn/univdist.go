// Package spn contains the structure of an SPN.
package spn

// A univariate distribution is a probability distribution with unary scope.
// UnivDist actually represents a multinomial distribution.
type UnivDist struct {
	// Parent node.
	pa Node
	// Variable ID
	varid int
	// Discrete probability distribution
	pr []float32
}

// Constructs a new UnivDist. Same as &UnivDist{pa, varid, dist}.
func NewUnivDist(pa Node, varid int, dist []float32) *UnivDist {
	return &UnivDist{pa, varid, dist}
}

// Constructs a new empty UnivDist for learning. We initialize pr to a uniform distribution.
// Argument m is the cardinality of varid.
func NewEmptyUnivDist(pa Node, varid, m int) *UnivDist {
	pr := make([]float32, m)
	for i := 0; i < m; i++ {
		pr[i] = 1 / float32(m)
	}
	return &UnivDist{pa, varid, pr}
}

// Ch returns the set of childre nodes. Since a node is a UnivDist iff it is a leaf, Ch=\emptyset.
func (self *UnivDist) Ch() []Node { return nil }

// Pa returns the parent node.
func (self *UnivDist) Pa() Node { return self.pa }

// Type return this node's type: 'leaf'.
func (self *UnivDist) Type() string { return "leaf" }

// Returns the probability of a certain valuation. That is Pr(X=valuation[varid]), where
// Pr=UnivDist.
func (self *UnivDist) Value(valuation VarSet) float32 {
	return self.pr[valuation[self.varid]]
}
