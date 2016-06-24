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
	pr []float64
	// Mode of pr. We pre-compute this to save time.
	mode float64
}

// Constructs a new UnivDist.
func NewUnivDist(pa Node, varid int, dist []float64) *UnivDist {
	n := len(dist)
	var m float64 = 0

	for i := 0; i < n; i++ {
		if dist[i] > m {
			m = dist[i]
		}
	}

	return &UnivDist{pa, varid, dist, m}
}

// Constructs a new empty UnivDist for learning. We initialize pr to a uniform distribution.
// Argument m is the cardinality of varid.
func NewEmptyUnivDist(pa Node, varid, m int) *UnivDist {
	pr := make([]float64, m)

	for i := 0; i < m; i++ {
		pr[i] = 1.0 / float64(m)
	}

	return &UnivDist{pa, varid, pr, pr[0]}
}

// Ch returns the set of childre nodes. Since a node is a UnivDist iff it is a leaf, Ch=\emptyset.
func (ud *UnivDist) Ch() []Node { return nil }

// Pa returns the parent node.
func (ud *UnivDist) Pa() Node { return ud.pa }

// Type return this node's type: 'leaf'.
func (ud *UnivDist) Type() string { return "leaf" }

// Returns the probability of a certain valuation. That is Pr(X=valuation[varid]), where
// Pr=UnivDist.
func (ud *UnivDist) Value(valuation VarSet) float64 {
	val, ok := valuation[ud.varid]
	if ok {
		return ud.pr[val]
	}
	return 1.0
}

// Max returns the MAP state given a valuation.
func (ud *UnivDist) Max(valuation VarSet) float64 {
	val, ok := valuation[ud.varid]
	if ok {
		return ud.pr[val]
	}
	return ud.mode
}
