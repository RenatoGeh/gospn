package spn

import (
	"github.com/RenatoGeh/gospn/utils"
)

// Indicator is an indicator node of a variable value X=x. Its value is 1 if X=x or is not set, and
// 0 otherwise.
type Indicator struct {
	Node
	// Variable ID
	varid int
	// Value of variable
	v int
}

// NewIndicator constructs a new indicator node.
func NewIndicator(varid int, v int) *Indicator {
	return &Indicator{Node{nil, []int{varid}}, varid, v}
}

// Type returns the type of this node.
func (i *Indicator) Type() string { return "leaf" }

// Value returns the probability of a certain valuation. In the case of an indicator node, 1 if X=x
// or is not set and 0 otherwise.
func (i *Indicator) Value(val VarSet) float64 {
	u, e := val[i.varid]
	if !e || u == i.v {
		return 0
	}
	return utils.LogZero
}

// Max returns the MAP.
func (i *Indicator) Max(val VarSet) float64 {
	u, e := val[i.varid]
	if !e || u == i.v {
		return 0
	}
	return utils.LogZero
}

// ArgMax returns the MAP and the MAP state.
func (i *Indicator) ArgMax(val VarSet) (VarSet, float64) {
	retval := make(VarSet)
	u, e := val[i.varid]
	if !e || u == i.v {
		retval[i.varid] = u
		return retval, 0
	}
	return retval, utils.LogZero
}

// Sc returns the scope of this node.
func (i *Indicator) Sc() []int {
	if len(i.sc) == 0 {
		i.sc = []int{i.varid}
	}
	return i.sc
}
