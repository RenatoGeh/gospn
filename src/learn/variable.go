package learn

// Variable is a wrapper struct that contains a variable ID and its observed distribution.
type Variable struct {
	// Variable ID.
	varid int
	// Distribution.
	pr []float64
}
