// Package learn contains the structural learning algorithm as well as a k-means clustering
// and a independence test.
package learn

// Variable is a wrapper struct that contains a variable ID and its observed distribution.
type Variable struct {
	// Variable ID.
	varid int
	// Distribution.
	pr []float64
}
