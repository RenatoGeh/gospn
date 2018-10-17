package learn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/utils"
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

// CompleteDataToMatrix returns a complete dataset's matrix form. A dataset is complete if every
// variable's varid in D's scope is equivalent to its map position and the scope's maximum varid is
// equal to the number of variables in D plus one (i.e. there are no "holes" in the scope).
func CompleteDataToMatrix(D spn.Dataset) [][]int {
	if D == nil || D[0] == nil {
		return nil
	}
	n := len(D)
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

// DataToMatrix returns a Dataset's matrix form. Assumes a consistent dataset.
func DataToMatrix(D spn.Dataset) ([][]int, map[int]int) {
	if D == nil || D[0] == nil {
		return nil, nil
	}
	n := len(D)
	// Assumption: D is consistent. A dataset D is consistent if for every pair of instances (I, J)
	// of D, Sc(I)=Sc(J). A direct implication of consistency is that len(D[i])=len(D[j]) for any
	// i!=j.
	m := len(D[0])
	M := make([][]int, n)
	S := make(map[int]int) // map[newid]varid
	Z := make(map[int]int) // map[varid]newid
	var u int
	for i, I := range D {
		M[i] = make([]int, m)
		for k, v := range I {
			if t, e := Z[k]; !e {
				Z[k] = u
				S[u] = k
				M[i][u] = v
				u++
			} else {
				M[i][t] = v
			}
		}
	}
	return M, S
}

// CompleteDataToMatrixF returns a complete dataset's matrix form. A dataset is complete if every
// variable's varid in D's scope is equivalent to its map position and the scope's maximum varid is
// equal to the number of variables in D plus one (i.e. there are no "holes" in the scope).
func CompleteDataToMatrixF(D spn.Dataset) [][]float64 {
	if D == nil || D[0] == nil {
		return nil
	}
	n := len(D)
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

// DataToMatrixF returns a Dataset's matrix form. Assumes a consistent dataset.
func DataToMatrixF(D spn.Dataset) ([][]float64, map[int]int) {
	if D == nil || D[0] == nil {
		return nil, nil
	}
	n := len(D)
	// Assumption: D is consistent. A dataset D is consistent if for every pair of instances (I, J)
	// of D, Sc(I)=Sc(J). A direct implication of consistency is that len(D[i])=len(D[j]) for any
	// i!=j.
	m := len(D[0])
	M := make([][]float64, n)
	Z := make(map[int]int) // map[varid]newid
	S := make(map[int]int) // map[newid]varid
	var u int
	fmt.Printf("n: %d, m: %d\n", n, m)
	for i, I := range D {
		M[i] = make([]float64, m)
		//fmt.Printf("  i: %d, len(I): %d\n", i, len(I))
		for k, v := range I {
			if t, e := Z[k]; !e {
				Z[k] = u // S and Z define a one-to-one and onto mapping.
				S[u] = k
				//fmt.Printf("    !e, k: %d, u: %d, v: %v\n", k, u, v)
				M[i][u] = float64(v)
				u++
			} else {
				//fmt.Printf("    e, k: %d, u: %d, v: %v\n", k, u, v)
				M[i][t] = float64(v)
			}
		}
	}
	return M, S
}

// MatrixToData returns a dataset from matrix M.
func MatrixToData(M [][]int) spn.Dataset {
	D := make(spn.Dataset, len(M))
	for i, L := range M {
		D[i] = make(map[int]int)
		for j, e := range L {
			D[i][j] = e
		}
	}
	return D
}

func ReflectScope(Sc map[int]*Variable) map[int]*Variable {
	nsc := make(map[int]*Variable)
	for u, v := range Sc {
		nsc[u] = v
	}
	return nsc
}

func CopyScope(Sc map[int]*Variable) map[int]*Variable {
	nsc := make(map[int]*Variable)
	for u, v := range Sc {
		nsc[u] = &Variable{Varid: v.Varid, Categories: v.Categories, Name: v.Name}
	}
	return nsc
}

func DataToVarData(D []map[int]int, Sc map[int]*Variable) []*utils.VarData {
	n := len(Sc)
	vdata, l := make([]*utils.VarData, n), 0
	for _, v := range Sc {
		tn := len(D)
		tdata := make([]int, tn)
		for j := 0; j < tn; j++ {
			tdata[j] = D[j][v.Varid]
		}
		vdata[l] = utils.NewVarData(v.Varid, v.Categories, tdata)
		l++
	}
	return vdata
}
