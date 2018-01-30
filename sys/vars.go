package sys

import "math"

// Global config
var (
	Verbose = false
)

// Image vars
var (
	Width  = 46
	Height = 56
	Max    = 8
)

// Gens parameters
var (
	Pval = 0.0001
	Eps  = 4.0
	Mp   = 4
)

// Math
var (
	EqualEpsilon    = 1e-5
	LogEqualEpsilon float64
)

// Memory variables
var (
	MemLowBoundShrink = 1024
	MemConservative   = true
)

func init() {
	LogEqualEpsilon = math.Log(EqualEpsilon)
}
