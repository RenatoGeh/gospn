package spn

// Mode of a univariate distribution.
type Mode struct {
	// Value of variable when it is the highest.
	index int
	// Highest value of variable.
	val float64
}

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
	mode Mode
	// Scope of this univariate distribution. We store this to avoid new creations of slices.
	sc []int
}

// Constructs a new UnivDist.
func NewUnivDist(pa Node, varid int, dist []float64) *UnivDist {
	n := len(dist)
	var m float64 = 0
	var mi int = 0

	for i := 0; i < n; i++ {
		if dist[i] > m {
			m = dist[i]
			mi = i
		}
	}

	return &UnivDist{pa, varid, dist, Mode{mi, m}, []int{varid}}
}

// Constructs a new empty UnivDist for learning. We initialize pr to a uniform distribution.
// Argument m is the cardinality of varid.
func NewEmptyUnivDist(pa Node, varid, m int) *UnivDist {
	pr := make([]float64, m)

	for i := 0; i < m; i++ {
		pr[i] = 1.0 / float64(m)
	}

	return &UnivDist{pa, varid, pr, Mode{0, pr[0]}, []int{varid}}
}

// Ch returns the set of childre nodes. Since a node is a UnivDist iff it is a leaf, Ch=\emptyset.
func (ud *UnivDist) Ch() []Node { return nil }

// Pa returns the parent node.
func (ud *UnivDist) Pa() Node { return ud.pa }

// Type return this node's type: 'leaf'.
func (ud *UnivDist) Type() string { return "leaf" }

func (ud *UnivDist) Sc() []int {
	// A univariate distribution has unary scope by definition.
	return ud.sc
}

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
	return ud.mode.val
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (ud *UnivDist) ArgMax(valuation VarSet) (VarSet, float64) {
	retval := make(VarSet)
	val, ok := valuation[ud.varid]

	if ok {
		retval[ud.varid] = val
		return retval, ud.pr[val]
	}

	retval[ud.varid] = ud.mode.index
	return retval, ud.mode.val
}
