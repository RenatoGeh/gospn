package parameters

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

// P is a collection of available parameters for learning algorithms.
//
// Disclaimer: Parameters do not work on inline methods (e.g. S.Value(E)) since that would require
// GoSPN storing a P pointer in each Node.
type P struct {
	Normalize    bool    // Normalize on weight update.
	HardWeight   bool    // Hard weights (true) or soft weights (false).
	LearningType int     // Soft or hard EM or GD (only applies to weight learning functions).
	Eta          float64 // Learning rate.
	Epsilon      float64 // Epsilon convergence criterion (in logspace).
	BatchSize    int     // Batch size if mini-batch. If bs <= 1, then no batching.
}

// Default returns a P instance with the following default options:
//  Normalize    = true
//  HardWeight   = false
//  HardLearning = parameters.SoftGD
//  Eta          = 0.1
//  Epsilon      = 1.0
//  BatchSize    = 0
func Default() *P {
	return &P{true, false, SoftGD, 0.1, 1.0, 0}
}

// New returns a P instance with the given parameters as option values.
func New(norm, hw bool, t int, eta, eps float64, bs int) *P {
	return &P{norm, hw, t, eta, eps, bs}
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
