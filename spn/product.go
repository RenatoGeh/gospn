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
	return &Product{Node: NewNode()}
}

// Type returns the type of this node.
func (p *Product) Type() string { return "product" }

// Sc returns the scope of this node.
func (p *Product) Sc() map[int]int {
	if p.sc == nil {
		n := len(p.ch)
		for i := 0; i < n; i++ {
			chsc := p.ch[i].Sc()
			for k := range chsc {
				p.sc[k] = k
			}
		}
	}
	return p.sc
}

// Bsoft is a common base for all soft inference methods.
func (p *Product) Bsoft(val VarSet, key string) float64 {
	prev := p.Stored(key)
	if prev > 0 {
		return prev
	}

	n := len(p.ch)
	ch := p.Ch()
	var v float64

	for i := 0; i < n; i++ {
		v += ch[i].Bsoft(val, key)
	}

	p.Store(key, v)
	return v
}

// Value returns the value of this SPN given a set of valuations.
func (p *Product) Value(val VarSet) float64 {
	return p.Bsoft(val, "soft")
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
				s += p.ch[j].Stored("soft")
			}
		}
		*da += math.Log1p(math.Exp(p.pnode + s - *da))
	}

	for i := 0; i < n; i++ {
		p.ch[i].Derive()
	}
}
