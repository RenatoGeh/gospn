package sspn

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_randist.h>
*/
import "C"

import (
	"math"
)

// Gaussian represents a gaussian distribution.
type Gaussian struct {
	Node
	// Variable ID
	varid int
	// Mean
	mean float64
	// Standard deviation
	sd float64
}

// NewGaussian constructs a new Gaussian from a counting slice.
func NewGaussian(varid int, counts []int) *Gaussian {
	var mean, sd float64
	var N int
	n := len(counts)

	// Standardizing gaussian from N(mean, sd) to N(0, 1).

	for i := 0; i < n; i++ {
		mean += float64(counts[i] * i)
		N += counts[i]
	}

	mean /= float64(N)

	for i := 0; i < n; i++ {
		d := float64(i) - mean
		sd += float64(counts[i]) * d * d
	}
	sd = math.Sqrt(sd)

	return &Gaussian{Node{sc: []int{varid}}, varid, mean, sd}
}

// Type returns the type of this node.
func (g *Gaussian) Type() string { return "leaf" }

// Value returns the probability of a certain valuation. That is Pr(X=val[varid]), where
// Pr is a probability function over a gaussian distribution.
func (g *Gaussian) Value(val VarSet) float64 {
	v, ok := val[g.varid]
	if ok {
		return math.Log(float64(C.gsl_ran_ugaussian_pdf(C.double((float64(v) - g.mean) / g.sd))))
	}
	return 0.0 // ln(1.0) = 0.0
}

// Max returns the MAP state given a valuation.
func (g *Gaussian) Max(val VarSet) float64 {
	v, ok := val[g.varid]
	if ok {
		return math.Log(float64(C.gsl_ran_ugaussian_pdf(C.double((float64(v) - g.mean) / g.sd))))
	}
	return math.Log(g.mean)
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (g *Gaussian) ArgMax(val VarSet) (VarSet, float64) {
	retval := make(VarSet)
	v, ok := val[g.varid]

	if ok {
		retval[g.varid] = v
		return retval,
			math.Log(float64(C.gsl_ran_ugaussian_pdf(C.double((float64(v) - g.mean) / g.sd))))
	}

	retval[g.varid] = int(g.mean)
	return retval, math.Log(g.mean)
}
