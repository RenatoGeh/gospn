package utils

// VarData is a wrapper struct that contains a variable ID and its observed data.
// Observed data is data that is to be used in learning. Each data instance i in data (i.e.
// data[i]) is a variable instantiation.
type VarData struct {
	// Variable ID.
	Varid int
	// Number of possible instantiations (levels/categories) of Varid.
	Categories int
	// Observed data.
	Data []int
}

// NewVarData constructs a new VarData. Equivalent to &VarData{varid, categories, data}.
func NewVarData(varid, categories int, data []int) *VarData {
	return &VarData{varid, categories, data}
}
