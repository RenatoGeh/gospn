package language

import (
	"github.com/RenatoGeh/gospn/spn"
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
	if _lv, ok := p.Stored(key); ok {
		return _lv
	}

	ch := p.Ch()

	v := 2 * ch[0].Soft(val, key)

	p.Store(key, v)
	return v
}

// Value returns the value of this SPN given a set of valuations.
func (p *SquareProduct) Value(val spn.VarSet) float64 {
	return p.Soft(val, "soft")
}
