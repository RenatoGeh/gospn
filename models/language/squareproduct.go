package language

import (
	"github.com/RenatoGeh/gospn/spn"
	"math"
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

// Derive recursively derives this node and its children based on the last inference value.
func (p *SquareProduct) Derive(wkey, nkey, ikey string) {
	ch := p.Ch()[0]

	st := ch.Storer()
	v, _ := p.Stored(nkey)
	u, _ := ch.Stored(ikey)
	st[nkey] += math.Log1p(math.Exp(math.Log(2) + v + u - st[nkey]))

	//for i := 0; i < n; i++ {
	//s := 0.0
	//for j := 0; j < n; j++ {
	//if i != j {
	//v, _ := p.ch[j].Stored(ikey)
	//s += v
	//}
	//}
	//st := p.ch[i].Storer()
	//v, _ := p.Stored(nkey)
	//st[nkey] += math.Log1p(math.Exp(v + s - st[nkey]))
	//}

	ch.Derive(wkey, nkey, ikey)
}
