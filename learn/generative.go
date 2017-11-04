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
func GenerativeGD(S spn.SPN, eta, eps float64, data spn.Dataset, c common.Collection, norm bool) spn.SPN {
	if c == nil {
		c = &common.Queue{}
	}

	storage := spn.NewStorer()
	stk, itk := storage.NewTicket(), storage.NewTicket()
	//wtk := storage.NewTicket()
	var ollh, llh float64
	sys.Println("Initiating Generative Gradient Descent...")
	for ok := true; ok; ok = (math.Abs(ollh-llh) > eps) {
		ollh = llh
		llh = 0.0
		n := len(data)
		var i int
		for _, I := range data {
			// Store inference values under T[itk].
			sys.Println("Storing inference values...")
			spn.StoreInference(S, I, itk, storage)
			lv, _ := storage.Single(itk, S)
			// Store SPN derivatives under T[stk].
			sys.Println("Computing dS(X)/dS...")
			DeriveSPN(S, storage, stk, itk, c)
			// Store weights derivatives under T[wtk].
			//sys.Println("Computing dS(X)/dW...")
			//DeriveWeights(S, storage, wtk, stk, itk, c)
			// Apply gradient descent.
			sys.Println("Applying gradient descent...")
			DeriveApplyWeights(S, eta, storage, stk, itk, c, norm)
			//applyGD(S, eta, wtk, storage, c, norm)
			// Reset DP tables.
			storage.Reset(itk)
			storage.Reset(stk)
			//applyGD(S, eta, wtk, storage, c, norm)
			// Reset weight derivative DP table.
			//storage.Reset(wtk)
			// Add current log-value to log-likelihood.
			sys.Printf("Log-value ln(S(X)) = %.3f\n", lv)
			llh += lv
			i++
			sys.Printf("Instance %d/%d.\n", i, n)
		}
		sys.Printf("Log-likelihood value at this iteration: llh = %.3f\n", llh)
		if sys.Verbose {
			dllh := math.Abs(ollh - llh)
			sys.Printf("Epsilon log-likelihood: eps = %.3f > %.3f ? %v \n", dllh, eps, dllh > eps)
		}
	}
	sys.Println("Generative gradient descent done. Returning...")

	return S
}

// This is where the magic happens.
func applyGD(S spn.SPN, eta float64, wtk int, storage *spn.Storer, c common.Collection, norm bool) {
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
				Normalize(W)
			}
		}
		for _, cs := range ch {
			if cs.Type() != "leaf" && !visited[cs] {
				c.Give(cs)
				visited[cs] = true
			}
		}
	}
	visited = nil
	c = nil
	sys.Free()
}