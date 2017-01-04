package language

import (
	"github.com/RenatoGeh/gospn/spn"
	"math/rand"
)

// Language is a language modelling SPN structure based on the article
// 	Language Modelling with Sum-Product Networks
// 	Wei-Chen Cheng, Stanley Kok, Hoai Vu Pham, Hai Leong Chieu, Kian Ming A. Chai
// 	INTERSPEECH 2014
// We shall refer to this article via the codename LMSPN.
// vfile is the vocabulary filename.
func Language(vfile string) {

}

// Structure returns the SPN structure as described in LMSPN.
func Structure(K, D, N int) spn.SPN {
	// K-dimensional vectors layer (V layer).

	// The V layer of N vectors
	V := make([]*Vector, N)
	for i := 0; i < N; i++ {
		// We give each vector an index i+1 (index 0 is reserved for query variable).
		V[i] = NewVector(i + 1)
	}

	// H_ij sum layer

	// wmatrix is the weight matrix
	// | w_11 w_12 w_13 ... w_1K |
	// | w_21 w_22 w_23 ... w_2K |
	// | ...  ...  ..   ... ...  |
	// | w_D1 w_D2 w_D3 ... w_DK |
	// cpmatrix and epmatrix are wmatrix's respective derivative slices
	wmatrix, cpmatrix, epmatrix := make([][]float64, D), make([][]float64, D), make([][]float64, D)
	for i := 0; i < D; i++ {
		wmatrix[i] = make([]float64, K)
		cpmatrix[i] = make([]float64, K)
		epmatrix[i] = make([]float64, K)
	}
	// hmatrix is the H sum nodes matrix
	// | H_11 H_12 H_13 ... H_1D |
	// | H_21 H_22 H_23 ... H_2D |
	// | ...  ...  ..   ... ...  |
	// | H_N1 H_N2 H_N3 ... H_ND |
	hmatrix := make([][]*SumVector, N)
	for i := 0; i < D; i++ {
		hmatrix[i] = make([]*SumVector, D)
	}

	// Create each H_ij node.
	for i := 0; i < N; i++ {
		for j := 0; j < D; j++ {
			hmatrix[i][j] = NewSumVector(wmatrix[i], cpmatrix[i], epmatrix[i])
			// Connect sum node H_ij to vector node V_i.
			hmatrix[i][j].AddChild(V[i])
		}
	}

	// M sum nodes layer

	M := make([]*spn.Sum, K)
	for i := 0; i < K; i++ {
		// Create each M_i sum node.
		M[i] = spn.NewSum()
		// Connect each H_pq sum node to M_i.
		for p := 0; p < N; p++ {
			for q := 0; q < D; q++ {
				// Give it a random [0,1) weight.
				M[i].AddChildW(hmatrix[p][q], rand.Float64())
			}
		}
	}

	// G product nodes layer

	G := make([]*spn.Product, K)
	for i := 0; i < K; i++ {
		G[i] = spn.NewProduct()
		// Add each M_i sum node twice as child to square it (simulating covariance).
		G[i].AddChild(M[i])
		G[i].AddChild(M[i])
	}

	// B sum nodes layer

	B := make([]*spn.Sum, K)
	for i := 0; i < K; i++ {
		B[i] = spn.NewSum()
		// Add both M_i and G_i as children of B_i.
		B[i].AddChildW(M[i], rand.Float64())
		B[i].AddChildW(G[i], rand.Float64())
	}

	// S product nodes layer

	S := make([]*ProductIndicator, K)
	for i := 0; i < K; i++ {
		S[i] = NewProductIndicator(i)
		// Add B_i as child of S_i.
		S[i].AddChild(B[i])
	}

	// Root node.
	R := spn.NewSum()
	for i := 0; i < K; i++ {
		// Add each S_i node to the root node.
		R.AddChildW(S[i], rand.Float64())
	}

	return R
}
