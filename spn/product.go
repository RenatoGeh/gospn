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

// Value returns the value of this SPN given a set of valuations.
func (p *Product) Value(val VarSet) float64 {
	n := len(p.ch)
	ch := p.Ch()
	var v float64

	for i := 0; i < n; i++ {
		v += ch[i].Value(val)
	}
	p.s = v

	return v
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
