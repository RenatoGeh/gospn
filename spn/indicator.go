package spn

import (
	"math"
)

// Indicator is an indicator variable (has value 1 or 0 according to evidence).
type Indicator struct {
	Node
	// Variable ID
	varid int
	// Variable instance setting (for bool, only 0 or 1). Indicates which instance this leaf
	// indicates.
	setting int
}

// NewIndicator constructs a new Indicator.
func NewIndicator(varid int, set int) *Indicator {
	return &Indicator{NewNode(varid), varid, set}
}

// Type returns the type of this node.
func (ind *Indicator) Type() string { return "leaf" }

// Soft is a common base for all soft inference methods.
func (ind *Indicator) Soft(val VarSet, key string) float64 {
	if _lv, ok := ind.Stored(key); ok && ind.stores {
		return _lv
	}

	v, ok := val[ind.varid]
	var l float64
	if !ok || v == ind.setting {
		l = 1.0
	}

	ind.Store(key, l)
	return l
}

// LSoft is Soft in logspace.
func (ind *Indicator) LSoft(val VarSet, key string) float64 {
	if _lv, ok := ind.Stored(key); ok && ind.stores {
		return _lv
	}

	v, ok := val[ind.varid]
	l := math.Inf(-1)
	if !ok || v == ind.setting {
		l = 0.0
	}

	ind.Store(key, l)
	return l
}

// Value is the value of this node.
func (ind *Indicator) Value(val VarSet) float64 {
	return ind.Soft(val, "soft")
}

// Max returns the MAP state given a valuation.
func (ind *Indicator) Max(val VarSet) float64 {
	v, ok := val[ind.varid]
	if !ok || v == ind.setting {
		return 0.0
	}
	return math.Inf(-1)
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (ind *Indicator) ArgMax(val VarSet) (VarSet, float64) {
	retval := make(VarSet)
	v, ok := val[ind.varid]

	if !ok || v == ind.setting {
		retval[ind.varid] = v
		return retval, 0.0
	}

	retval[ind.varid] = ind.setting
	return retval, math.Inf(-1)
}
