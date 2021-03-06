package learn

import (
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/utils"
	"math"
)

// DeriveSPN computes the derivative dS/dS_i, for every child of S: S_i. The base case dS/dS is
// trivial and equal to one. For each child S_i and parent node S_n, the derivative is given by:
//  dS/dS_i <- dS/dS_i + w_{n,i} * dS/dS_n, if S_n is sum node
//  dS/dS_j <- dS/dS_j + dS/dS_n * \prod_{k \in \Ch(n) \setminus \{j\}} S_k
// Where w_{n,i} is the weight of edge S_n -> S_i and Ch(n) is the set of children of n.
// In other words, the derivative of a sum node with respect to the SPN is the weighted sum of the
// derivatives of its parent nodes. For product nodes, the derivative is a sum where the elements
// of such sum are the products of each parent node multiplied by all its siblings.
// It is relevant to note that since GoSPN treats values in logspace, all the derivatives are too
// in logspace. Argument tk is the ticket to be used for storing the derivatives. Argument itk is
// the ticket for the stored values of S(X) (i.e. soft inference). A Collection is required for the
// graph search, though if Collection c is nil, then we use a Queue. If a Queue is used, then the
// graph search is a breadth-first, if a Stack is used, then it performs a depth-first search.
// If tk < 0, then a new ticket will be created and returned alongside the SPN S.
// Returns the SPN S and the ticket used.
func DeriveSPN(S spn.SPN, storage *spn.Storer, tk, itk int, c common.Collection) (spn.SPN, int) {
	if tk < 0 {
		tk = storage.NewTicket()
	}
	if c == nil {
		c = &common.Queue{}
	}

	table, _ := storage.Table(tk)
	inf, _ := storage.Table(itk)
	visited := make(map[spn.SPN]bool)
	table.StoreSingle(S, 0.0)
	c.Give(S)
	visited[S] = true

	for !c.Empty() {
		s := c.Take().(spn.SPN)
		ch := s.Ch()
		pv, _ := table.Single(s)
		//sys.Printf("pv=%f, parent: %s\n", pv, s.Type())
		if s.Type() == "sum" {
			sum := s.(*spn.Sum)
			W := sum.Weights()
			for i, cs := range ch {
				v, e := table.Single(cs)
				if !e {
					table.StoreSingle(cs, math.Log(W[i])+pv)
				} else {
					table.StoreSingle(cs, utils.LogSumExpPair(v, math.Log(W[i])+pv))
				}
			}
		} else /* there can never be a case where s is a leaf, therefore s is a product */ {
			for i, cs := range ch {
				v, e := table.Single(cs)
				t := pv
				for j := range ch {
					if j != i {
						_v, _ := inf.Single(ch[j])
						t += _v
					}
				}
				//sys.Printf("  i: %d, v: %f, e: %v, t: %f, ch: %p\n", i, v, e, t, cs)
				if !e {
					table.StoreSingle(cs, t)
				} else {
					//table.StoreSingle(cs, math.Log(math.Exp(v)+math.Exp(t)))
					table.StoreSingle(cs, utils.LogSumExpPair(v, t))
				}
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
	return S, tk
}

// DeriveWeights computes the derivative dS/dW, where W is the multiset of weights in SPN S.
// The derivative of S with respect to W is given by
// 	dS/dw_{n,j} <- S_j * dS/dS_n, if S_n is a sum node
// It is only relevant to compute dS/dw_{n,j} in sum nodes since weights do not appear in product
// nodes. Argument S is the SPN to find the derivative of. Argument storage is the DP storage
// object we store the derivatives values and extract inference values from. Integers tk, dtk and
// itk are the tickets for where to store dS/dW, where to locate dS/dS_i and stored inference
// values respectively. Collection c is the data type to be used for the graph search. If c is a
// stack, then DeriveWeights performs a depth-first search. If c is a queue, then DeriveWeights's
// graph search is a breadth-first search. The default value for c is Queue. DeriveWeights returns
// the SPN S and a ticket if tk is a negative value.
func DeriveWeights(S spn.SPN, storage *spn.Storer, tk, dtk, itk int, c common.Collection) (spn.SPN, int) {
	if tk < 0 {
		tk = storage.NewTicket()
	}
	if c == nil {
		c = &common.Queue{}
	}

	wt, _ := storage.Table(tk)
	st, _ := storage.Table(dtk)
	it, _ := storage.Table(itk)
	visited := make(map[spn.SPN]bool)
	c.Give(S)
	visited[S] = true

	for !c.Empty() {
		s := c.Take().(spn.SPN)
		ch := s.Ch()
		pv, _ := st.Single(s)
		if s.Type() == "sum" {
			for i, cs := range ch {
				v, _ := it.Single(cs)
				wt.Store(s, i, v+pv)
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
	return S, tk
}

// DeriveWeightsBatch computes the derivative dS/dW, where W is the multiset of weights in SPN S
// and adds it to the given Storer.
// The derivative of S with respect to W is given by
// 	dS/dw_{n,j} <- S_j * dS/dS_n, if S_n is a sum node
// It is only relevant to compute dS/dw_{n,j} in sum nodes since weights do not appear in product
// nodes. Argument S is the SPN to find the derivative of. Argument storage is the DP storage
// object we store the derivatives values and extract inference values from. Integers tk, dtk and
// itk are the tickets for where to store dS/dW, where to locate dS/dS_i and stored inference
// values respectively. Collection c is the data type to be used for the graph search. If c is a
// stack, then DeriveWeights performs a depth-first search. If c is a queue, then DeriveWeights's
// graph search is a breadth-first search. The default value for c is Queue. DeriveWeights returns
// the SPN S and a ticket if tk is a negative value.
func DeriveWeightsBatch(S spn.SPN, storage *spn.Storer, tk, dtk, itk int, c common.Collection) (spn.SPN, int) {
	if tk < 0 {
		tk = storage.NewTicket()
	}
	if c == nil {
		c = &common.Queue{}
	}

	wt, _ := storage.Table(tk)
	st, _ := storage.Table(dtk)
	it, _ := storage.Table(itk)
	visited := make(map[spn.SPN]bool)
	c.Give(S)
	visited[S] = true

	for !c.Empty() {
		s := c.Take().(spn.SPN)
		ch := s.Ch()
		pv, _ := st.Single(s)
		if s.Type() == "sum" {
			for i, cs := range ch {
				v, _ := it.Single(cs)
				dw, e := wt.Entry(s, i)
				ndw := v + pv
				if e {
					ndw = utils.LogSumExpPair(ndw, dw)
				}
				//sys.Printf("ndw: %.10f, v: %.10f, pv: %.10f, dw: %.10f\n", ndw, v, pv, dw)
				wt.Store(s, i, ndw)
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
	return S, tk
}

func Normalize(v []float64) {
	var min float64
	for _, u := range v {
		if u < min {
			min = u
		}
	}
	if min < 0 {
		for i := range v {
			v[i] -= min
		}
	}

	var norm float64
	for i := range v {
		norm += v[i]
	}
	//var c float64
	for i := range v {
		v[i] /= norm
		//c += v[i]
	}
	//max := math.Inf(-1)
	//for i := range v {
	//if v[i] > max {
	//max = v[i]
	//}
	//}
	//for i := range v {
	//v[i] /= max
	//}
}

// DeriveApplyWeights does not store the weight derivatives like DeriveWeights. Instead, it
// computes and applies the gradient on the go.
func DeriveApplyWeights(S spn.SPN, eta float64, storage *spn.Storer, dtk, itk int, c common.Collection, norm bool) spn.SPN {
	visited := make(map[spn.SPN]bool)
	if c == nil {
		c = &common.Queue{}
	}
	st, _ := storage.Table(dtk)
	it, _ := storage.Table(itk)
	c.Give(S)
	visited[S] = true
	for !c.Empty() {
		s := c.Take().(spn.SPN)
		ch := s.Ch()
		pv, _ := st.Single(s)
		if s.Type() == "sum" {
			sum := s.(*spn.Sum)
			W := sum.Weights()
			for i, cs := range ch {
				v, _ := it.Single(cs)
				W[i] += eta * math.Exp(v+pv)
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
	return S
}

// DeriveHard performs hard inference (MAP) derivation on the SPN. The hard derivative is the
// number of times MAP states pass a certain weight. The delta weight is then computed as
//  eta*c/w
// where eta is the learning rate, c is the number of times hard inference passed through weight w
// and w is the weight of the current edge.
func DeriveHard(S spn.SPN, st *spn.Storer, tk int, I spn.VarSet) int {
	if tk < 0 {
		tk = st.NewTicket()
	}
	tab, _ := st.Table(tk)

	T := spn.TraceMAP(S, I)
	Q := common.Queue{}
	V := make(map[spn.SPN]bool)
	Q.Enqueue(S)
	V[S] = true

	for !Q.Empty() {
		s := Q.Dequeue().(spn.SPN)
		ch := s.Ch()
		switch t := s.Type(); t {
		case "product":
			for _, c := range ch {
				if c.Type() != "leaf" && !V[c] {
					Q.Enqueue(c)
					V[c] = true
				}
			}
		case "sum":
			mi, e := T[s]
			if e {
				v, _ := tab.Entry(s, mi)
				tab.Store(s, mi, v+1)
				//fmt.Printf("%d ", int(v+1))
				if mc := ch[mi]; !V[mc] {
					Q.Enqueue(mc)
					V[mc] = true
				}
			}
		}
	}
	//fmt.Println("")

	return tk
}
