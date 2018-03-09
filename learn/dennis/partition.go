package dennis

import (
	"github.com/RenatoGeh/gospn/spn"
)

type partition struct {
	ch  []*region
	rep []*spn.Product
}

func newPartition() *partition {
	return &partition{nil, nil}
}

func (p *partition) add(r *region) {
	p.ch = append(p.ch, r)
}

func (p *partition) translate(n int) []*spn.Product {
	p.rep = make([]*spn.Product, n)
	for i := 0; i < n; i++ {
		p.rep[i] = spn.NewProduct()
	}
	return p.rep
}
