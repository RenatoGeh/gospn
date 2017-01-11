package language

import (
	"github.com/RenatoGeh/gospn/spn"
	//"math"
)

// SquareProduct is a product node that squares its only child.
type SquareProduct struct {
	spn.Product
}

// NewSquareProduct creates a new SquareProduct
func NewSquareProduct() *SquareProduct {
	return &SquareProduct{*spn.NewProduct()}
}

// Soft is a common base for all soft inference methods.
func (p *SquareProduct) Soft(val spn.VarSet, key string) float64 {
	if _lv, ok := p.Stored(key); ok && p.Stores() {
		return _lv
	}

	ch := p.Ch()

	v := ch[0].Soft(val, key)
	v *= v

	p.Store(key, v)
	return v
}

// Value returns the value of this SPN given a set of valuations.
func (p *SquareProduct) Value(val spn.VarSet) float64 {
	return p.Soft(val, "soft")
}

// Derive derives this node only.
func (p *SquareProduct) Derive(wkey, nkey, ikey string) {
	ch := p.Ch()[0]

	st := ch.Storer()
	v, _ := p.Stored(nkey)
	u, _ := ch.Stored(ikey)
	st[nkey] += 2 * (v * u)

	//ch.Derive(wkey, nkey, ikey)
}
