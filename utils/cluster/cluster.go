package cluster

func copyMatrixF(A [][]int) [][]float64 {
	D := make([][]float64, len(A))
	for i, I := range A {
		J := make([]float64, len(I))
		for j, v := range I {
			J[j] = float64(v)
		}
		D[i] = J
	}
	return D
}

func copyMatrix(A [][]float64) [][]float64 {
	B := make([][]float64, len(A))
	for i, a := range A {
		B[i] = make([]float64, len(a))
		copy(B[i], a)
	}
	return B
}

func toCluster(k int, D [][]int, G []int) []map[int][]int {
	C := make([]map[int][]int, k)
	for i := range C {
		C[i] = make(map[int][]int)
	}
	for i, j := range G {
		I := D[i]
		J := make([]int, len(I))
		for l, v := range I {
			J[l] = int(v)
		}
		C[j-1][i] = J
	}
	return C
}

func toClusterF(k int, D [][]float64, G []int) []map[int][]float64 {
	C := make([]map[int][]float64, k)
	for i := range C {
		C[i] = make(map[int][]float64)
	}
	for i, j := range G {
		I := D[i]
		J := make([]float64, len(I))
		copy(J, I)
		C[j-1][i] = J
	}
	return C
}

func copyMap(src map[int]int) map[int]int {
	dst := make(map[int]int)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func guessToData(k int, G []int, M [][]float64, Sc map[int]int) [][]map[int]int {
	C := make([][]map[int]int, k)
	for i, g := range G {
		I := make(map[int]int)
		for j := range M[i] {
			I[Sc[j]] = int(M[i][j])
		}
		C[g] = append(C[g], I)
	}
	return C
}
