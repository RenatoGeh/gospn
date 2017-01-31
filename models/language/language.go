package language

import (
	"fmt"
	"github.com/RenatoGeh/gospn/spn"
	"math"
	"math/rand"
)

// Language is a language modelling SPN structure based on the article
// 	Language Modelling with Sum-Product Networks
// 	Wei-Chen Cheng, Stanley Kok, Hoai Vu Pham, Hai Leong Chieu, Kian Ming A. Chai
// 	INTERSPEECH 2014
// We shall refer to this article via the codename LMSPN.
// vfile is the vocabulary filename, D is the dimension of the feature vectors and N is number of
// previous words to be set as evidence.
func Language(vfile string, D, N int) {
	fmt.Println("Parsing voc file...")
	voc := Parse(vfile)
	K := voc.Size()
	fmt.Println("Creating SPN structure...")
	S := Structure(K, D, N)
	//fmt.Println("Learning...")
	Learn(S, voc, D, N, spn.SOFT)

	pre := make(spn.VarSet)
	fmt.Printf("Generated the following first %d words from vocabulary:\n ", N)
	//for i := 1; i <= N; i++ {
	//w, id := voc.Rand()
	//pre[i] = id
	//fmt.Printf(" %s", w)
	//}
	pre[1] = 0
	pre[2] = 1
	pre[3] = 2

	const M = 100
	S.SetStore(false)
	fmt.Printf("\nInferring the next %d words...\n ", M)
	for i := 0; i < M; i++ {
		imax, max := -1, -1.0
		for j := 0; j < K; j++ {
			pre[0] = j
			S.RResetDP("")
			v := S.Value(pre)
			fmt.Printf("\nPr(X=%d|%d", j, pre[1])
			for l := 2; l < len(pre); l++ {
				fmt.Printf(",%d", pre[l])
			}
			fmt.Printf(")=%.10f", v)
			if v > max {
				max, imax = v, j
			}
		}
		fmt.Printf(" %s", voc.Translate(imax))
		for j := N; j >= 2; j-- {
			pre[j] = pre[j-1]
		}
		pre[1] = imax
	}
}

// Learn learns weights according to LMSPN.
func Learn(S spn.SPN, voc *Vocabulary, D, N int, mode spn.InfType) spn.SPN {
	const eta = 0.01

	conv := 1.0
	last := 0.0

	S.SetStore(true)
	voc.Set(N)
	S.Normalize()
	S.SetL2(0.0001)
	combs := voc.Combinations()
	for _l := 0; _l < 2; _l++ {
		s := 0.0
		klast := 0.0
		for i := 0; i < combs; i++ {
			S.RResetDP("")
			S.Rootify("cpnode")
			S.Rootify("epnode")
			C, E := voc.Next(), make(spn.VarSet)
			m := len(C)
			fmt.Printf("Learning with words: %s", voc.Translate(C[0]))
			for i := 1; i < m; i++ {
				E[i] = C[i]
				fmt.Printf(" %s", voc.Translate(C[i]))
			}
			fmt.Printf("\n")
			ds := spn.NewDiscStorer(S, C, E)
			//ds.Store(false)
			// Stores correct/guess values.
			fmt.Println("Storing correct/guess soft inference values...")
			fmt.Printf("Correct = %f\n", S.Soft(C, "correct"))
			// Derive correct/guess nodes.
			fmt.Println("Derivating correct/guess nodes...")
			S.RootDerive("cpweight", "cpnode", "correct", mode)
			// Stores expected values.
			fmt.Println("Storing expected soft inference values...")
			fmt.Printf("Expected = %f\n", S.Soft(E, "expected"))
			// Derive expected nodes.
			fmt.Println("Derivating expected nodes...")
			S.RootDerive("epweight", "epnode", "expected", mode)
			// Update weights.
			fmt.Println("Updating weights...")
			S.DiscUpdate(eta, ds, "cpweight", "epweight", mode)

			fmt.Printf("Adding convergence diff for instance %d...\n", i)
			S.RResetDP("")
			k := S.Value(C)
			fmt.Printf("Diff component %d: %f\n", i, k)
			s += k - klast
			klast = k

			//T := make(spn.VarSet)
			//T[1] = 1
			//T[2] = 2
			//T[3] = 3
			//for j := 0; j < voc.Size(); j++ {
			//T[0] = j
			//S.RResetDP("")
			//v := S.Value(T)
			//fmt.Printf("\nPr(X=%d|%d", j, T[1])
			//for l := 2; l < len(T); l++ {
			//fmt.Printf(",%d", T[l])
			//}
			//fmt.Printf(")=%.10f", v)
			//}

			C = nil
			E = nil
		}
		d := s - last
		last = s
		conv = d
		fmt.Printf("Discriminative Learning diff: %.5f\n", math.Abs(conv))
		voc.Set(N)
	}

	return S
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
	wmatrix := make([][]float64, D)
	for i := 0; i < D; i++ {
		wmatrix[i] = make([]float64, K)
		for j := 0; j < K; j++ {
			wmatrix[i][j] = rand.Float64()
		}
	}
	// hmatrix is the H sum nodes matrix
	// | H_11 H_12 H_13 ... H_1D |
	// | H_21 H_22 H_23 ... H_2D |
	// | ...  ...  ..   ... ...  |
	// | H_N1 H_N2 H_N3 ... H_ND |
	hmatrix := make([][]*SumVector, N)
	for i := 0; i < N; i++ {
		hmatrix[i] = make([]*SumVector, D)
		// Create each H_ij node.
		for j := 0; j < D; j++ {
			hmatrix[i][j] = NewSumVector(wmatrix[j])
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
		M[i].AutoNormalize(true)
	}

	// G product nodes layer

	G := make([]*SquareProduct, K)
	for i := 0; i < K; i++ {
		G[i] = NewSquareProduct()
		// Add each M_i sum node as child to square it (simulating covariance).
		G[i].AddChild(M[i])
	}

	// B sum nodes layer

	B := make([]*spn.Sum, K)
	for i := 0; i < K; i++ {
		B[i] = spn.NewSum()
		// Add both M_i and G_i as children of B_i.
		B[i].AddChildW(M[i], rand.Float64())
		B[i].AddChildW(G[i], rand.Float64())
		B[i].AutoNormalize(true)
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
	R.AutoNormalize(true)
	for i := 0; i < K; i++ {
		// Add each S_i node to the root node.
		R.AddChildW(S[i], rand.Float64())
	}

	return R
}
