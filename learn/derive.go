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
		c = &common.Stack{}
	}

	table, _ := storage.Table(tk)
	inf, _ := storage.Table(itk)
	table[S] = 0.0
	c.Give(S)

	for !c.Empty() {
		s := c.Take().(spn.SPN)
		ch := s.Ch()
		pv := table[s]
		if s.Type() == "sum" {
			sum := s.(*spn.Sum)
			W := sum.Weights()
			for i, cs := range ch {
				v, e := table[cs]
				if !e {
					table[cs] = math.Log(W[i]) + pv
				} else {
					table[cs] = math.Log(math.Exp(v) + W[i]*math.Exp(pv))
				}
			}
		} else /* there can never be a case where s is a leaf, therefore s is a product */ {
			for i, cs := range ch {
				v, e := table[cs]
				t := pv
				for j := range ch {
					if j != i {
						t += inf[ch[j]]
					}
				}
				if !e {
					table[cs] = t
				} else {
					table[cs] = math.Log(math.Exp(v) + math.Exp(t))
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
