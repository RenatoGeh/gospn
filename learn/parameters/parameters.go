package parameters

import (
	"sync"
)

// Constants to be used for P.LearningType.
const (
	HardGD = iota // Hard gradient descent key.
	SoftGD        // Soft gradient descent key.
	HardEM        // Hard expectation-maximization key.
	SoftEM        // Soft expectation-maximization key.
)

// Constants to be used with Method(P.LearningType).
const (
	GD = iota
	EM
)

// Constants to be used with Hardness(P.LearningType).
const (
	Hard = iota
	Soft
)

var (
	mu sync.Mutex
)

// P is a collection of available parameters for learning algorithms.
//
// Disclaimer: Parameters do not work on inline methods (e.g. S.Value(E)) since that would require
// GoSPN storing a P pointer in each Node.
type P struct {
	Normalize    bool    // Normalize on weight update.
	HardWeight   bool    // Hard weights (true) or soft weights (false).
	SmoothSum    float64 // Constant for smoothing sum counts when hard weights is true.
	LearningType int     // Soft or hard EM or GD (only applies to weight learning functions).
	Eta          float64 // Learning rate.
	Epsilon      float64 // Epsilon convergence criterion (in logspace).
	BatchSize    int     // Batch size if mini-batch. If bs <= 1, then no batching.
	Lambda       float64 // Regularization constant.
	Iterations   int     // Number of iterations for gradient descent.
}

// Default returns a P instance with the following default options:
//  Normalize    = true
//  HardWeight   = false
//  SmoothSum    = 0.01
//  HardLearning = parameters.SoftGD
//  Eta          = 0.1
//  Epsilon      = 1.0
//  BatchSize    = 0
func Default() *P {
	return &P{true, false, 0.01, SoftGD, 0.1, 1.0, 0, 0.01, 4}
}

// New returns a P instance with the given parameters as option values.
func New(norm, hw bool, sm float64, t int, eta, eps float64, bs int, l float64, i int) *P {
	return &P{norm, hw, sm, t, eta, eps, bs, l, i}
}

// Method returns what super-type of learning method P.LearningType is.
func Method(t int) int {
	if t <= 1 {
		return GD
	}
	return EM
}

// Hardness returns whether P.LearningType uses soft or hard inference for learning.
func Hardness(t int) int {
	if t%2 == 0 {
		return Hard
	}
	return Soft
}

// Parametriable defines a type that has parameters.
type Parametrizable interface {
	// Parameters returns the parameters of this object.
	Parameters() *P
}

var bindings map[Parametrizable]*P

func init() {
	bindings = make(map[Parametrizable]*P)
}

func Bind(e Parametrizable, p *P) {
	mu.Lock()
	bindings[e] = p
	mu.Unlock()
}

func Unbind(e Parametrizable) {
	mu.Lock()
	delete(bindings, e)
	mu.Unlock()
}

func Exists(p Parametrizable) bool {
	mu.Lock()
	_, e := bindings[p]
	mu.Unlock()
	return e
}

func Retrieve(e Parametrizable) (*P, bool) {
	mu.Lock()
	p, q := bindings[e]
	mu.Unlock()
	return p, q
}
