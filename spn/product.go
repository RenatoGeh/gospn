package spn

import (
	"math"
)

// Product represents a product node in an SPN.
type Product struct {
	Node
}

// NewProduct returns a new Product node.
func NewProduct() *Product {
	return &Product{}
}

// Type returns the type of this node.
func (p *Product) Type() string { return "product" }

// Sc returns the scope of this node.
func (p *Product) Sc() []int {
	if p.sc == nil {
		n := len(p.ch)
		for i := 0; i < n; i++ {
			chsc := p.ch[i].Sc()
			k := len(chsc)
			for j := 0; j < k; j++ {
				p.sc = append(p.sc, chsc[j])
			}
		}
	}
	return p.sc
}

// Bsoft is a common base for all soft inference methods.
func (p *Product) Bsoft(val VarSet, where *float64) float64 {
	if *where > 0 {
		return *where
	}

	n := len(p.ch)
	ch := p.Ch()
	var v float64

	for i := 0; i < n; i++ {
		v += ch[i].Bsoft(val, where)
	}

	*where = v
	return v
}

// Value returns the value of this SPN given a set of valuations.
func (p *Product) Value(val VarSet) float64 {
	return p.Bsoft(val, &p.s)
}

// Max returns the MAP state given a valuation.
func (p *Product) Max(val VarSet) float64 {
	n := len(p.ch)
	var v float64

	for i := 0; i < n; i++ {
		v += (p.ch[i]).Max(val)
	}

	return v
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (p *Product) ArgMax(val VarSet) (VarSet, float64) {
	argmax := make(VarSet)
	n := len(p.ch)
	var v float64

	for i := 0; i < n; i++ {
		chargs, chmap := p.ch[i].ArgMax(val)
		v += chmap
		for k, val := range chargs {
			argmax[k] = val
		}
	}

	return argmax, v
}

// Derive recursively derives this node and its children based on the last inference value.
func (p *Product) Derive() {
	n := len(p.ch)

	for i := 0; i < n; i++ {
		da := p.ch[i].DrvtAddr()
		s := 0.0
		for j := 0; j < n; j++ {
			if i != j {
				s += p.ch[j].Stored()
			}
		}
		*da += math.Log1p(math.Exp(p.pnode + s - *da))
	}

	for i := 0; i < n; i++ {
		p.ch[i].Derive()
	}
}

// CondValue returns the value of this SPN queried on Y and conditioned on X.
// Let S be this SPN. If S is the root node, then CondValue(Y, X) = S(Y|X). Else we store the value
// of S(Y, X) in Y so that we don't need to recompute Union(Y, X) at every iteration.
func (p *Product) CondValue(Y VarSet, X VarSet) float64 {
	if p.root {
		for k, v := range X {
			Y[k] = v
		}
	}
	p.st = p.Bsoft(Y, &p.st)
	p.sb = p.Bsoft(X, &p.sb)
	p.scnd = p.st - p.sb

	// Store values for each sub-SPN.
	n := len(p.ch)
	for i := 0; i < n; i++ {
		p.ch[i].CondValue(Y, X)
	}

	return p.scnd
}
