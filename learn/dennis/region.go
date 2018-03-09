package dennis

import (
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/utils"
)

type region struct {
	ch  []*partition
	sc  scope
	rep []spn.SPN
	id  uint64
}

func newRegion(sc scope) *region {
	return &region{nil, sc, nil, 0}
}

func (r *region) add(p *partition) {
	r.ch = append(r.ch, p)
}

func (r *region) translate(D spn.Dataset, m int) []spn.SPN {
	var v int
	for k, _ := range r.sc {
		v = k
	}
	r.rep = make([]spn.SPN, m)
	if len(r.sc) == 1 {
		X := learn.ExtractInstance(v, D)
		Q := utils.PartitionQuantiles(X, m)
		for i, q := range Q {
			r.rep[i] = spn.NewGaussianParams(v, q[0], q[1])
		}
	} else {
		for i := 0; i < m; i++ {
			r.rep[i] = spn.NewSum()
		}
	}
	return r.rep
}
