package language

import (
	"github.com/RenatoGeh/gospn/spn"
)

// Vector represents a K-dimensional vector, where K is the size of the vocabulary.
type Vector struct {
	spn.Node
	// Index
	i int
}

// NewVector creates a new vector that represents the index-th previous word.
func NewVector(index int) *Vector {
	// Leave Node with all fields set to default values.
	return &Vector{i: index}
}

// Value returns the word index (which entry is set to 1).
func (v *Vector) Value(val spn.VarSet) float64 { return v.Soft(val, "soft") }

// Max returns the word index (which entry is set to 1).
func (v *Vector) Max(val spn.VarSet) float64 { return float64(val[v.i]) }

// Type returns the type of this node ("vector").
func (v *Vector) Type() string { return "vector" }

// Soft is a common base for all soft inference methods.
func (v *Vector) Soft(val spn.VarSet, key string) float64 {
	w := float64(val[v.i])
	v.Store(key, w)
	return w
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (v *Vector) ArgMax(val spn.VarSet) (spn.VarSet, float64) {
	retval := make(spn.VarSet)
	retval[v.i] = val[v.i]
	return retval, float64(val[v.i])
}

// AddParent does nothing.
func (v *Vector) AddParent(p spn.SPN) {}
