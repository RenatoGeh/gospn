package spn

import (
	//"github.com/RenatoGeh/gospn/sys"
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
	return &Gaussian{Node{nil, []int{varid}}, varid, distuv.Normal{mu, sigma, nil}}
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

	//sys.Printf("Created new gaussian with Mu: %f and StdDev: %f\n", mean, sd)
	return &Gaussian{Node{sc: []int{varid}}, varid, distuv.Normal{mean, sd, nil}}
}

// NewGaussianFit constructs a new Gaussian from GoNum's Fit function.
func NewGaussianFit(varid int, counts []float64) *Gaussian {
	N := distuv.Normal{}
	sample := make([]float64, len(counts))
	for i := range sample {
		sample[i] = float64(i)
	}
	N.Fit(sample, counts)
	return &Gaussian{Node{sc: []int{varid}}, varid, N}
}

// Type returns the type of this node.
func (g *Gaussian) Type() string { return "leaf" }

func zeroSigma(v int, ok bool, mu float64) float64 {
	if ok && v == int(mu) {
		return 0
	}
	return math.Inf(-1)
}

// Value returns the probability of a certain valuation. That is Pr(X=val[varid]), where
// Pr is a probability function over a gaussian distribution.
func (g *Gaussian) Value(val VarSet) float64 {
	v, ok := val[g.varid]
	var l float64
	if g.dist.Sigma == 0 {
		return zeroSigma(v, ok, g.dist.Mu)
	} else if ok {
		l = g.dist.LogProb(float64(v))
	} else {
		l = 0.0 // ln(1.0) = 0.0
	}
	//sys.Printf("Gaussian value (mu=%f, sigma=%f) for value %d (pixel %d): %f = ln(%f)\n", g.dist.Mu, g.dist.Sigma, v, g.varid, l, math.Exp(l))
	return l
}

// Max returns the MAP state given a valuation.
func (g *Gaussian) Max(val VarSet) float64 {
	v, ok := val[g.varid]
	if g.dist.Sigma == 0 {
		return zeroSigma(v, ok, g.dist.Mu)
	} else if ok {
		return g.dist.LogProb(float64(v))
	}
	return g.dist.LogProb(g.dist.Mu)
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (g *Gaussian) ArgMax(val VarSet) (VarSet, float64) {
	retval := make(VarSet)
	v, ok := val[g.varid]

	if g.dist.Sigma == 0 {
		retval[g.varid] = int(g.dist.Mu)
		return retval, 0.0
	} else if ok {
		retval[g.varid] = v
		z := g.dist.LogProb(float64(v))
		return retval, z
	}

	retval[g.varid] = int(g.dist.Mu)
	return retval, g.dist.LogProb(g.dist.Mu)
}

// Params returns mean and standard deviation.
func (g *Gaussian) Params() (float64, float64) {
	return g.dist.Mu, g.dist.Sigma
}

// Sc returns the scope of this node.
func (g *Gaussian) Sc() []int {
	if len(g.sc) == 0 {
		g.sc = []int{g.varid}
	}
	return g.sc
}
