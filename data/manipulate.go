package data

import (
	"github.com/RenatoGeh/gospn/learn"
	"math"
	"math/rand"
)

func cascadeRounding(p []float64) []int {
	n := len(p)
	q := make([]int, n)
	var s float64
	var t int
	for i := range p {
		s += p[i]
		q[i] = int(math.Round(s - float64(t)))
		t += q[i]
	}
	return q
}

func copyMap(d map[int]int, s map[int]int) {
	for k, v := range s {
		d[k] = v
	}
}

// ExtractLabels attempts to separate the real variable values and the labels from a dataset. A
// label is always the last variable in a .data file. The converse is not true, since a dataset may
// not contain labels if it's not a classification job. In this case, the ExtractLabels function
// still tries to extract the last real variable values as labels. It is up to the user to only use
// ExtractLabels when the dataset is known to have classification labels.
// Return values are the original scope unaltered, the dataset with label values taken out from the
// matrix, the label variable, and a slice where each value in index i contains the classification
// value of the i-th element of the design matrix.
func ExtractLabels(S map[int]*learn.Variable, D []map[int]int) (map[int]*learn.Variable, []map[int]int, *learn.Variable, []int) {
	n, m := len(S), len(D)
	lv := S[n-1]
	L := make([]int, m)
	M := make([]map[int]int, m)

	for i, I := range D {
		k := 0
		M[i] = make(map[int]int)
		for p, v := range I {
			if p > k {
				k = p
			}
			M[i][p] = v
		}
		l := I[k]
		L[i] = l
		delete(M[i], k)
	}

	return S, M, lv, L
}

// Partition partitions dataset D into random subdatasets following the proportions given by p. For
// example, if p=(0.3, 0.7), Partition will return a slice P of size |p| where |P[0]|=0.3*|D| and
// |P[1]|=0.7*|D|. This function assumes D has no labels. For a balanced uniformly partitioning wrt
// the labels of the dataset, use PartitionByLabels.
func Partition(D []map[int]int, p []float64) [][]map[int]int {
	n, m := len(D), len(p)
	R := rand.Perm(n)
	q := make([]float64, m)
	for i := range q {
		q[i] = float64(n) * p[i]
	}
	S := cascadeRounding(q)
	M := make([][]map[int]int, m)
	var k int
	for i, s := range S {
		M[i] = make([]map[int]int, s)
		for j := 0; j < s; j++ {
			M[i][j] = make(map[int]int)
			copyMap(M[i][j], D[R[k]])
			k++
		}
	}
	return M
}

// PartitionByLabels partitions the dataset D in a similar fashion to Partition. However,
// PartitionByLabels tries to keep the same proportion of labels for each subdataset. If the
// result of the proportions multiplied by |D| is an integer, then PartitionByLabels returns an
// exact partitioning following given proportions. Otherwise, the function tries to best
// approximate the given proportions.
// Arguments are the original dataset D, slice L of true labels of each instance, the number of
// classes c, and p the proportions.
func PartitionByLabels(D []map[int]int, L []int, c int, p []float64) ([][]map[int]int, [][]int) {
	C := make([][]int, c) // Indices of instances divided by their class.
	for i := range D {
		l := L[i]
		C[l] = append(C[l], i)
	}
	N := make([]int, c) // Number of indices in each class.
	for i := range C {
		N[i] = len(C[i])
	}
	n, m := len(D), len(p)
	S := make([][]int, m) // Size of each partition (element is a c-tuple).
	s := make([]int, m)
	for i := 0; i < m; i++ {
		q := make([]float64, c) // Tuple of size c, with each q_i the number of elements of class i.
		for j := 0; j < c; j++ {
			q[j] = math.Floor(float64(N[j]) * p[i])
		}
		S[i] = cascadeRounding(q)
		for _, v := range S[i] {
			s[i] += v
		}
	}
	R := make([][]int, c) // Random index permutations for each class.
	for i := 0; i < c; i++ {
		R[i] = rand.Perm(N[i])
	}
	K := make([]int, c)           // Counters for each class.
	M := make([][]map[int]int, m) // Resulting maps.
	U := make([][]int, m)         // Resulting label slices.
	for i := 0; i < m; i++ {
		M[i] = make([]map[int]int, s[i])
		U[i] = make([]int, s[i])
		var t int
		for j := 0; j < c; j++ {
			z := S[i][j]
			for u := 0; u < z; u++ {
				l := R[j][K[j]]
				u := C[j][l]
				M[i][t] = make(map[int]int)
				copyMap(M[i][t], D[u])
				U[i][t] = j
				t++
				K[j]++
			}
		}
	}
	// Give leftovers to partition that is furthest from its ideal state.
	var d float64
	var di int
	for i := 0; i < m; i++ {
		u := float64(n) / float64(len(M[i]))
		du := math.Abs(u - p[i])
		if du > d {
			d, di = du, i
		}
	}
	for i := 0; i < c; i++ {
		for K[i] < len(R[i]) {
			l := R[i][K[i]]
			V := make(map[int]int)
			copyMap(V, D[l])
			M[di] = append(M[di], V)
			U[di] = append(U[di], i)
			K[i]++
		}
	}

	return M, U
}

// Split takes a dataset D, the number of classes c and label assignments L and returns the dataset
// split by labels. That is, create c subdatasets where for each of these subdatasets, all items in
// the i-th dataset belongs to the class i, and such that the union of all subdatasets is D.
func Split(D []map[int]int, c int, L []int) [][]map[int]int {
	S := make([][]map[int]int, c)
	for i, l := range L {
		M := make(map[int]int)
		copyMap(M, D[i])
		S[l] = append(S[l], M)
	}
	return S
}

// Copy copies the dataset and labels. If labels does not exist, returns only the dataset.
func Copy(D []map[int]int, L []int) ([]map[int]int, []int) {
	if len(D) != len(L) && L != nil {
		panic("Length of dataset does not match length of labels.")
	}
	n := len(D)
	nD := make([]map[int]int, n)
	var nL []int
	if L != nil {
		nL = make([]int, n)
	}
	for i := 0; i < n; i++ {
		nD[i] = make(map[int]int)
		copyMap(nD[i], D[i])
		if L != nil {
			nL[i] = L[i]
		}
	}
	return nD, nL
}

// Shuffle shuffles a dataset and sets its labels accordingly (if it exists). It is an in-place shuffle.
func Shuffle(D []map[int]int, L []int) {
	if len(D) != len(L) && L != nil {
		panic("Length of dataset does not match length of labels.")
	}
	rand.Shuffle(len(D), func(i, j int) {
		D[i], D[j] = D[j], D[i]
		if L != nil {
			L[i], L[j] = L[j], L[i]
		}
	})
}

// Join returns the concatenation of two datasets and their labels (if both exist).
func Join(D []map[int]int, E []map[int]int, L []int, M []int) ([]map[int]int, []int) {
	n, m := len(D), len(E)
	o := n + m
	c := L != nil && M != nil
	if c && (len(L) != n || len(M) != m) {
		panic("Length of dataset does not match length of labels.")
	}
	F := make([]map[int]int, o)
	var P []int
	if c {
		P = make([]int, o)
	}
	for i := 0; i < n; i++ {
		F[i] = make(map[int]int)
		copyMap(F[i], D[i])
		if c {
			P[i] = L[i]
		}
	}
	for i := n; i < o; i++ {
		F[i] = make(map[int]int)
		copyMap(F[i], E[i-n])
		if c {
			P[i] = M[i-n]
		}
	}
	return F, P
}

// SubtractLabel takes a dataset D, its label array and a label l and returns the result of
// subtracting D with every instance of label l. It also returns the set of instances of label l.
// It does not modify D. Instead, it copies D and returns the result of the subtraction.
func SubtractLabel(D []map[int]int, L []int, l int) ([]map[int]int, []int, []map[int]int, []int) {
	n := len(L)
	if n != len(D) {
		panic("Length of dataset does not match length of labels.")
	}
	var S, T []map[int]int
	var P, Q []int
	for i := 0; i < n; i++ {
		if L[i] == l {
			t := make(map[int]int)
			copyMap(t, D[i])
			T = append(T, t)
			Q = append(Q, l)
		} else {
			s := make(map[int]int)
			copyMap(s, D[i])
			S = append(S, s)
			P = append(P, L[i])
		}
	}
	return S, P, T, Q
}

// SubtractVariable takes a dataset D, a variable v and returns the result of subtracting all
// entries that belong to variable v. It does not modify D. Instead, it copies D and returns the
// result of the subtraction.
func SubtractVariable(D []map[int]int, v *learn.Variable) []map[int]int {
	u := v.Varid
	E := make([]map[int]int, len(D))
	for i, I := range D {
		E[i] = make(map[int]int)
		for k, v := range I {
			if k != u {
				E[i][k] = v
			}
		}
	}
	return E
}
