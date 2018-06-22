package learn

import (
	"github.com/RenatoGeh/gospn/spn"
)

// Variable is a wrapper struct that contains the variable ID and its number of categories.
type Variable struct {
	// Variable ID.
	Varid int
	// Number of categories.
	Categories int
	// Variable name.
	Name string
}

type Scope map[int]*Variable

// ExtractInstance extracts all instances of variable v from dataset D and joins them into a single
// slice.
func ExtractInstance(v int, D spn.Dataset) []int {
	p := make([]int, len(D))
	for i, I := range D {
		p[i] = I[v]
	}
	return p
}

// DataToMatrix returns a Dataset's matrix form. Assumes a consistent dataset.
func DataToMatrix(D spn.Dataset) [][]int {
	if D == nil || D[0] == nil {
		return nil
	}
	n := len(D)
	// Assumption: D is consistent. A dataset D is consistent if for every pair of instances (I, J)
	// of D, Sc(I)=Sc(J). A direct implication of consistency is that len(D[i])=len(D[j]) for any
	// i!=j.
	m := len(D[0])
	M := make([][]int, n)
	for i, I := range D {
		M[i] = make([]int, m)
		for k, v := range I {
			M[i][k] = v
		}
	}
	return M
}

// DataToMatrixF returns a Dataset's matrix form. Assumes a consistent dataset.
func DataToMatrixF(D spn.Dataset) [][]float64 {
	if D == nil || D[0] == nil {
		return nil
	}
	n := len(D)
	// Assumption: D is consistent. A dataset D is consistent if for every pair of instances (I, J)
	// of D, Sc(I)=Sc(J). A direct implication of consistency is that len(D[i])=len(D[j]) for any
	// i!=j.
	m := len(D[0])
	M := make([][]float64, n)
	for i, I := range D {
		M[i] = make([]float64, m)
		for k, v := range I {
			M[i][k] = float64(v)
		}
	}
	return M
}
