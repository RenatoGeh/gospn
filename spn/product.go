package spn

import (
	//"fmt"
	//"math"
	"github.com/RenatoGeh/gospn/common"
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

// Soft is a common base for all soft inference methods.
func (p *Product) Soft(val VarSet, key string) float64 {
	if _lv, ok := p.Stored(key); ok && p.stores {
		return _lv
	}

	n := len(p.ch)
	ch := p.Ch()
	v := 1.0

	for i := 0; i < n; i++ {
		v *= ch[i].Soft(val, key)
	}

	p.Store(key, v)
	return v
}

// LSoft is Soft in logspace.
func (p *Product) LSoft(val VarSet, key string) float64 {
	if _lv, ok := p.Stored(key); ok {
		return _lv
	}

	n := len(p.ch)
	ch := p.Ch()
	var v float64

	for i := 0; i < n; i++ {
		v += ch[i].LSoft(val, key)
	}

	p.Store(key, v)
	return v
}

// Value returns the value of this SPN given a set of valuations.
func (p *Product) Value(val VarSet) float64 {
	return p.Soft(val, "soft")
}

// Max returns the MAP state given a valuation.
func (p *Product) Max(val VarSet) float64 {
	n := len(p.ch)
	v := 1.0

	for i := 0; i < n; i++ {
		v *= (p.ch[i]).Max(val)
	}

	return v
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (p *Product) ArgMax(val VarSet) (VarSet, float64) {
	argmax := make(VarSet)
	n := len(p.ch)
	v := 1.0

	for i := 0; i < n; i++ {
		chargs, chmap := p.ch[i].ArgMax(val)
		v *= chmap
		for k, val := range chargs {
			argmax[k] = val
		}
	}

	return argmax, v
}

// Derive derives this node only.
func (p *Product) Derive(wkey, nkey, ikey string, mode InfType) int {
	n := len(p.ch)

	if mode == SOFT {
		v, _ := p.Stored(nkey)
		for i := 0; i < n; i++ {
			s := 1.0
			for j := 0; j < n; j++ {
				if i != j {
					v, _ := p.ch[j].Stored(ikey)
					s *= v
				}
			}
			st := p.ch[i].Storer()
			st[nkey] += v * s
		}

		//for i := 0; i < n; i++ {
		//p.ch[i].Derive(wkey, nkey, ikey)
		//}
	}

	return -1
}

// RootDerive derives all nodes in a BFS fashion.
func (p *Product) RootDerive(wkey, nkey, ikey string, mode InfType) {
	q := common.Queue{}

	q.Enqueue(p)

	for !q.Empty() {
		t := q.Dequeue().(SPN)
		ch := t.Ch()

		r := t.Derive(wkey, nkey, ikey, mode)

		if ch != nil && r != 0 {
			if r < 0 {
				n := len(ch)
				for i := 0; i < n; i++ {
					q.Enqueue(ch[i])
				}
			} else {
				q.Enqueue(ch[r-1])
			}
		}
	}
}
