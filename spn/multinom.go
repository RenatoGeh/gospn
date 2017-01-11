package spn

import (
	"math"
)

// Mode of a univariate distribution.
type Mode struct {
	// Value of variable when it is the highest.
	index int
	// Highest value of variable.
	val float64
}

// Multinomial represents a multinomial distribution.
type Multinomial struct {
	Node
	// Variable ID
	varid int
	// Discrete probability distribution
	pr []float64
	// Mode of pr. We pre-compute this to save time.
	mode Mode
}

// NewMultinomial constructs a new Multinomial.
func NewMultinomial(varid int, dist []float64) *Multinomial {
	n := len(dist)
	var m float64
	var mi int

	for i := 0; i < n; i++ {
		if dist[i] > m {
			m = dist[i]
			mi = i
		}
	}

	return &Multinomial{NewNode(varid), varid, dist, Mode{mi, m}}
}

// NewCountingMultinomial constructs a new Multinomial from a count slice.
func NewCountingMultinomial(varid int, counts []int) *Multinomial {
	n := len(counts)

	pr := make([]float64, n)
	s := 0.0
	for i := 0; i < n; i++ {
		s += 1.0 + float64(counts[i])
		pr[i] = float64(1 + counts[i])
	}

	for i := 0; i < n; i++ {
		pr[i] /= float64(s)
	}

	var m float64
	var mi int

	for i := 0; i < n; i++ {
		if pr[i] > m {
			m = pr[i]
			mi = i
		}
	}

	return &Multinomial{NewNode(varid), varid, pr, Mode{mi, m}}
}

// NewEmptyMultinomial constructs a new empty Multinomial for learning.  Argument m is the
// cardinality of varid.
func NewEmptyMultinomial(varid, m int) *Multinomial {
	pr := make([]float64, m)

	for i := 0; i < m; i++ {
		pr[i] = 1.0 / float64(m)
	}

	lsc := make(map[int]int)
	lsc[varid] = varid
	return &Multinomial{NewNode(varid), varid, pr, Mode{0, pr[0]}}
}

// Type returns the type of this node.
func (m *Multinomial) Type() string { return "leaf" }

// Soft is a common base for all soft inference methods.
func (m *Multinomial) Soft(val VarSet, key string) float64 {
	if _lv, ok := m.Stored(key); ok && m.stores {
		return _lv
	}

	v, ok := val[m.varid]
	var l float64
	if ok {
		l = m.pr[v]
	} else {
		l = 1.0
	}
	m.Store(key, l)

	return l
}

// LSoft is Soft in logspace.
func (m *Multinomial) LSoft(val VarSet, key string) float64 {
	_lv := m.s[key]
	if _lv >= 0 {
		return _lv
	}

	v, ok := val[m.varid]
	var l float64
	if ok {
		l = math.Log(m.pr[v])
	} else {
		l = 0.0 // ln(1.0) = 0.0
	}
	m.Store(key, l)

	return l
}

// Value returns the probability of a certain valuation. That is Pr(X=val[varid]), where
// Pr is a probability function over a Multinomial distribution.
func (m *Multinomial) Value(val VarSet, key string) float64 {
	return m.Soft(val, key)
}

// Max returns the MAP state given a valuation.
func (m *Multinomial) Max(val VarSet) float64 {
	v, ok := val[m.varid]
	if ok {
		return math.Log(m.pr[v])
	}
	return math.Log(m.mode.val)
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (m *Multinomial) ArgMax(val VarSet) (VarSet, float64) {
	retval := make(VarSet)
	v, ok := val[m.varid]

	if ok {
		retval[m.varid] = v
		return retval, math.Log(m.pr[v])
	}

	retval[m.varid] = m.mode.index
	return retval, math.Log(m.mode.val)
}
