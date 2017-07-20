package learn

import (
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"math"
)

// GenerativeGD performs a generative gradient descent parameter learning on SPN S. Argument eta is
// the learning rate; eps is the likelihood difference to consider convergence, the more will
// GenerativeGD try to fit data; data is the dataset; c is how we should perform the graph search.
// If a stack is used, perform a DFS. If a queue is used, BFS. If c is nil, we use a queue.
// Argument norm indicates whether GenerativeGD should normalize weights at each node.
func GenerativeGD(S spn.SPN, eta, eps float64, data []map[int]int, c common.Collection, norm bool) spn.SPN {
	if c == nil {
		c = &common.Queue{}
	}

	storage := NewStorer()
	wtk, stk, itk := storage.NewTicket(), storage.NewTicket(), storage.NewTicket()
	var ollh, llh float64
	sys.Println("Initiating Generative Gradient Descent...")
	for ok := true; ok; ok = (math.Abs(ollh-llh) > eps) {
		ollh = llh
		llh = 0.0
		for _, I := range data {
			// Store inference values under T[itk].
			sys.Println("Storing inference values...")
			StoreInference(S, I, itk, storage)
			// Store SPN derivatives under T[stk].
			sys.Println("Computing dS(X)/dS...")
			DeriveSPN(S, storage, stk, itk, c)
			// Store weights derivatives under T[wtk].
			sys.Println("Computing dS(X)/dW...")
			DeriveWeights(S, storage, wtk, stk, itk, c)
			// Apply gradient descent.
			sys.Println("Applying gradient descent...")
			applyGD(S, eta, wtk, storage, c, norm)
			// Reset derivative DP tables.
			storage.Reset(wtk)
			storage.Reset(stk)
			// Add current log-value to log-likelihood.
			lv, _ := storage.Single(itk, S)
			sys.Printf("Log-value ln(S(X)) = %.3f\n", lv)
			llh += lv
			// Reset inference DP table.
			storage.Reset(itk)
		}
		sys.Printf("Log-likelihood value at this iteration: llh = %.3f\n", llh)
		if sys.Verbose {
			dllh := math.Abs(ollh - llh)
			sys.Printf("Epsilon log-likelihood: eps = %.3f > %.3f ? %v \n", dllh, eps, dllh > eps)
		}
	}

	return S
}

func normalize(v []float64) {
	var norm float64
	for i := range v {
		norm += v[i]
	}
	for i := range v {
		v[i] /= norm
	}
}

// This is where the magic happens.
func applyGD(S spn.SPN, eta float64, wtk int, storage *Storer, c common.Collection, norm bool) {
	visited := make(map[spn.SPN]bool)
	wt, _ := storage.Table(wtk)
	c.Give(S)
	visited[S] = true
	for !c.Empty() {
		s := c.Take().(spn.SPN)
		ch := s.Ch()
		if s.Type() == "sum" {
			sum := s.(*spn.Sum)
			W := sum.Weights()
			dW, _ := wt.Value(s)
			for i := range W {
				delta := eta * math.Exp(dW[i])
				W[i] += delta
			}
			if norm {
				normalize(W)
			}
		}
		for _, cs := range ch {
			if cs.Type() != "leaf" && !visited[cs] {
				c.Give(cs)
				visited[cs] = true
			}
		}
	}
}
