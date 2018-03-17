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
	//for ok := true; ok; ok = (math.Abs(ollh-llh) > eps) {
	for _l := 0; _l < 4; _l++ {
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

// GenerativeHardGD performs a generative gradient descent using hard inference.
func GenerativeHardGD(S spn.SPN, eta, eps float64, data spn.Dataset, c common.Collection, norm bool) spn.SPN {
	if c == nil {
		c = &common.Queue{}
	}

	storage := spn.NewStorer()
	dtk, itk := storage.NewTicket(), storage.NewTicket()
	var ollh, llh float64
	sys.Println("Initiating Generative Gradient Descent...")
	//for ok := true; ok; ok = (math.Abs(ollh-llh) > eps) {
	for _l := 0; _l < 5; _l++ {
		ollh = llh
		llh = 0.0
		n := len(data)
		var i int
		for _, I := range data {
			sys.Println("Storing inference values...")
			spn.StoreInference(S, I, itk, storage)
			lv, _ := storage.Single(itk, S)
			sys.Println("Computing hard derivatives...")
			DeriveHard(S, storage, dtk, I)
			sys.Println("Applying gradient descent...")
			applyHGD(S, eta, dtk, storage, norm)
			// Reset DP tables.
			storage.Reset(itk)
			storage.Reset(dtk)
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

// GenerativeBGD performs a generative batch gradient descent parameter learning on SPN S. Argument
// eta is the learning rate; eps is the likelihood difference to consider convergence, the more
// will GenerativeGD try to fit data; data is the dataset; c is how we should perform the graph
// search.  If a stack is used, perform a DFS. If a queue is used, BFS. If c is nil, we use a
// queue.  Argument norm indicates whether GenerativeGD should normalize weights at each node.
// bSize is the size of the batch.
//
// Batch means that all derivatives will be computed with the same structure and weights. Once we
// have completed a full iteration on the dataset, we then add all delta weights and apply them
// through gradient descent.
func GenerativeBGD(S spn.SPN, eta, eps float64, data spn.Dataset, c common.Collection, norm bool, bSize int) spn.SPN {
	if c == nil {
		c = &common.Queue{}
	}

	storage := spn.NewStorer()
	stk, itk, wtk := storage.NewTicket(), storage.NewTicket(), storage.NewTicket()
	var ollh, llh float64
	sys.Println("Initiating Generative Gradient Descent...")
	//for ok := true; ok; ok = (math.Abs(ollh-llh) > eps) {
	for _l := 0; _l < 1; _l++ {
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
			sys.Println("Computing dS(X)/dW...")
			DeriveWeightsBatch(S, storage, wtk, stk, itk, c)
			// Reset DP tables.
			storage.Reset(itk)
			storage.Reset(stk)
			// Add current log-value to log-likelihood.
			sys.Printf("Log-value ln(S(X)) = %.3f\n", lv)
			llh += lv
			i++
			if i%bSize == 0 {
				sys.Println("Applying gradient descent...")
				applyGD(S, eta, wtk, storage, c, norm)
				storage.Reset(wtk)
			}
			sys.Printf("Instance %d/%d.\n", i, n)
		}
		// Apply gradient descent.
		if i%bSize != 0 {
			sys.Println("Applying gradient descent...")
			applyGD(S, eta, wtk, storage, c, norm)
			storage.Reset(wtk)
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

// GenerativeHardBGD performs a batch generative gradient descent using hard inference.
func GenerativeHardBGD(S spn.SPN, eta, eps float64, data spn.Dataset, c common.Collection, norm bool, bSize int) spn.SPN {
	if c == nil {
		c = &common.Queue{}
	}

	storage := spn.NewStorer()
	dtk, itk := storage.NewTicket(), storage.NewTicket()
	var ollh, llh float64
	sys.Println("Initiating Generative Gradient Descent...")
	//for ok := true; ok; ok = (math.Abs(ollh-llh) > eps) {
	for _l := 0; _l < 5; _l++ {
		ollh = llh
		llh = 0.0
		n := len(data)
		var i int
		for _, I := range data {
			sys.Println("Storing inference values...")
			spn.StoreInference(S, I, itk, storage)
			lv, _ := storage.Single(itk, S)
			sys.Println("Computing hard derivatives...")
			DeriveHard(S, storage, dtk, I)
			storage.Reset(itk)
			i++
			if i%bSize == 0 {
				sys.Println("Applying gradient descent...")
				applyHGD(S, eta, dtk, storage, norm)
				storage.Reset(dtk)
			}
			// Add current log-value to log-likelihood.
			sys.Printf("Log-value ln(S(X)) = %.3f\n", lv)
			llh += lv
			sys.Printf("Instance %d/%d.\n", i, n)
		}
		if i%bSize != 0 {
			sys.Println("Applying gradient descent...")
			applyHGD(S, eta, dtk, storage, norm)
			storage.Reset(dtk)
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
				sys.Printf("delta=%.10f=%.2f*Exp(%.10f), W[%d]=%.10f\n", delta, eta, dW[i], i, W[i])
				//if delta > 1 {
				//delta = 1
				//} else if delta < -1 {
				//delta = -1
				//}
				W[i] += delta
				pW := sum.Weights()
				sys.Printf("W[%d]=%.10f, post-W[%d]=%.10f\n", i, W[i], i, pW[i])
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

func applyHGD(S spn.SPN, eta float64, tk int, st *spn.Storer, norm bool) {
	tab, _ := st.Table(tk)
	Q := common.Queue{}
	V := make(map[spn.SPN]bool)
	Q.Enqueue(S)
	V[S] = true
	for !Q.Empty() {
		s := Q.Dequeue().(spn.SPN)
		ch := s.Ch()
		if s.Type() == "sum" {
			v, e := tab.Value(s)
			if e {
				W := s.(*spn.Sum).Weights()
				for i, d := range v {
					w := W[i]
					delta := eta * math.Log(d+1) / w
					W[i] = w + delta
					//pW := s.(*spn.Sum).Weights()
					//sys.Printf("Ch: %d/%d, pre-W[%d]=%.10f, d: %.10f, eta*d/w: %.10f, post-W[%d]=%.10f\n", i, len(ch)-1, i, w, d, delta, i, pW[i])
					//for j := range ch {
					//if i != j {
					//sys.Printf("  sibling %d: %.10f\n", j, pW[j])
					//}
					//}
				}
				if norm {
					Normalize(W)
				}
			}
		}
		for _, c := range ch {
			if !V[c] && c.Type() != "leaf" {
				Q.Enqueue(c)
				V[c] = true
			}
		}
	}
}
