package learn

import (
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/spn"
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
func DeriveSPN(S spn.SPN, storage *Storer, tk, itk int, c common.Collection) (spn.SPN, int) {
	if tk < 0 {
		tk = storage.NewTicket()
	}
	if c == nil {
		c = &common.Queue{}
	}

	table, _ := storage.Table(tk)
	inf, _ := storage.Table(itk)
	table.StoreSingle(S, 0.0)
	c.Give(S)

	for !c.Empty() {
		s := c.Take().(spn.SPN)
		ch := s.Ch()
		pv, _ := table.Single(s)
		if s.Type() == "sum" {
			sum := s.(*spn.Sum)
			W := sum.Weights()
			for i, cs := range ch {
				v, e := table.Single(cs)
				if !e {
					table.StoreSingle(cs, math.Log(W[i])+pv)
				} else {
					table.StoreSingle(cs, math.Log(math.Exp(v)+math.Exp(math.Log(W[i])+pv)))
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
				if !e {
					table.StoreSingle(cs, t)
				} else {
					table.StoreSingle(cs, math.Log(math.Exp(v)+math.Exp(t)))
				}
			}
		}

		for _, cs := range ch {
			if cs.Type() != "leaf" {
				c.Give(cs)
			}
		}
	}

	return S, tk
}

// DeriveWeights computes the derivative dS/dW, where W is the multiset of weights in SPN S.
// The derivative of S with respect to W is given by
// 	dS/dw_{n,j} <- S_j * dS/dS_n, if S_n is a sum node
// It is only relevant to compute dS/dw_{n,j} in sum nodes since weights do not appear in product
// nodes. Argument S is the SPN to find the derivative of. Argument storage is the DP storage
// object we store the derivatives values and extract inference values from. Integers tk and itk
// are the tickets for the stored derivatives and inference values respectively. Collection c is
// the data type to be used for the graph search. If c is a stack, then DeriveWeights performs a
// depth-first search. If c is a queue, then DeriveWeights's graph search is a breadth-first
// search. The default value for c is Queue. DeriveWeights returns the SPN S and a ticket if tk is
// a negative value.
func DeriveWeights(S spn.SPN, storage *Storer, tk, itk int, c common.Collection) (spn.SPN, int) {
	if tk < 0 {
		tk = storage.NewTicket()
	}
	if c == nil {
		c = &common.Queue{}
	}

	table, _ := storage.Table(tk)
	table.StoreSingle(S, 0.0)
	c.Give(S)

	for !c.Empty() {
		s := c.Take().(spn.SPN)
		ch := s.Ch()
		pv, _ := table.Single(s)
		if s.Type() == "sum" {
			sum := s.(*spn.Sum)
			W := sum.Weights()
			for i, cs := range ch {
				v, e := table.Single(cs)
				if !e {
					table.StoreSingle(cs, math.Log(W[i])+pv)
				} else {
					table.StoreSingle(cs, math.Log(math.Exp(v)+math.Exp(math.Log(W[i])+pv)))
				}
			}
		}

		for _, cs := range ch {
			if cs.Type() != "leaf" {
				c.Give(cs)
			}
		}
	}

	return S, tk
}
