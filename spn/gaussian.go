package spn

import (
	//"fmt"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
)

const (
	// GaussMax is the maximum value of a standard gaussian, namely 1/sqrt(2*pi).
	GaussMax = 0.398942280 // 1/sqrt(2*pi) := max value of a standard Gaussian
)

// Gaussian represents a gaussian distribution.
type Gaussian struct {
	Node
	// Variable ID
	varid int
	// GoNum's normal distribution
	dist distuv.Normal
}

// NewGaussianParams constructs a new Gaussian from a mean and variance.
func NewGaussianParams(varid int, mu float64, sigma float64) *Gaussian {
	sd := math.Sqrt(sigma)
	return &Gaussian{Node{sc: []int{varid}}, varid, distuv.Normal{mu, sd, nil}}
}

// NewGaussianRaw constructs a new Gaussian from a slice of values.
func NewGaussianRaw(varid int, vals []float64) *Gaussian {
	var mean, sd float64
	n := len(vals)

	for i := 0; i < n; i++ {
		mean += vals[i]
	}
	mean /= float64(n)

	for i := 0; i < n; i++ {
		d := vals[i] - mean
		sd += d * d
	}

	sd = math.Sqrt(sd / float64(n))

	return &Gaussian{Node{sc: []int{varid}}, varid, distuv.Normal{mean, sd, nil}}
}

// NewGaussian constructs a new Gaussian from a counting slice.
func NewGaussian(varid int, counts []int) *Gaussian {
	var mean, sd float64
	var N int
	n := len(counts)

	for i := range counts {
		N += counts[i]
	}

	for i := 0; i < n; i++ {
		mean += float64(counts[i]) / float64(N) * float64(i)
	}

	for i := 0; i < n; i++ {
		d := float64(i) - mean
		sd += (float64(counts[i]) / float64(N)) * d * d
	}
	sd = math.Sqrt(sd)
	if sd == 0 {
		sd = 1
	}

	return &Gaussian{Node{sc: []int{varid}}, varid, distuv.Normal{mean, sd, nil}}
}

// Type returns the type of this node.
func (g *Gaussian) Type() string { return "leaf" }

// Value returns the probability of a certain valuation. That is Pr(X=val[varid]), where
// Pr is a probability function over a gaussian distribution.
func (g *Gaussian) Value(val VarSet) float64 {
	v, ok := val[g.varid]
	if ok {
		//fmt.Println("Yelloooo")
		return g.dist.LogProb(float64(v))
	}
	return 0.0 // ln(1.0) = 0.0
}

// Max returns the MAP state given a valuation.
func (g *Gaussian) Max(val VarSet) float64 {
	v, ok := val[g.varid]
	if ok {
		return g.dist.LogProb(float64(v))
	}
	return math.Log(GaussMax)
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (g *Gaussian) ArgMax(val VarSet) (VarSet, float64) {
	retval := make(VarSet)
	v, ok := val[g.varid]

	if ok {
		retval[g.varid] = v
		z := g.dist.LogProb(float64(v))
		return retval, z
	}

	retval[g.varid] = int(g.dist.Mu)
	return retval, math.Log(GaussMax)
}
