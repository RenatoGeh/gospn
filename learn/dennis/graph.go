package dennis

import (
	"fmt"
)

type node interface {
	Type() string
}

type graph struct {
	root *region
	R    []*region
	P    map[string]*partition
	gId  uint64
}

func newGraph(sc scope) *graph {
	r := newRegion(sc)
	return &graph{r, []*region{r}, make(map[string]*partition), 0}
}

func (g *graph) registerRegion(r *region) {
	r.id = g.gId
	g.gId++
	g.R = append(g.R, r)
}

func (g *graph) registerPartition(p *partition, r, s, t *region) {
	f := fmt.Sprintf("%d,%d,%d", r.id, s.id, t.id)
	g.P[f] = p
}

func (g *graph) existsPartition(r, s, t *region) bool {
	f := fmt.Sprintf("%d,%d,%d", r.id, s.id, t.id)
	_, e := g.P[f]
	return e
}

// validateRegion either returns an existent region node in g or creates a new one and registers.
func (g *graph) validateRegion(s scope) *region {
	for _, r := range g.R {
		if s.equal(r.sc) {
			return r
		}
	}
	r := newRegion(s)
	g.registerRegion(r)
	return r
}

// allScopes returns the set of all region scopes.
func (g *graph) allScopes() scopeSlice {
	S := make(scopeSlice, len(g.R))
	for i, r := range g.R {
		S[i] = r.sc
	}
	return S
}

func postorderRec(n node, V map[node]bool, O *[]*region) {
	if t := n.Type(); t == "region" {
		r := n.(*region)
		if len(r.ch) != 0 {
			for _, c := range r.ch {
				if !V[c] {
					V[c] = true
					postorderRec(c, V, O)
				}
			}
		}
		*O = append(*O, r)
	} else /* partition */ {
		p := n.(*partition)
		for _, c := range p.ch {
			if !V[c] {
				V[c] = true
				postorderRec(c, V, O)
			}
		}
	}
}

// postorder returns a slice of the regions in postorder.
func (g *graph) postorder() []*region {
	V := make(map[node]bool)
	var O []*region
	V[g.root] = true
	postorderRec(g.root, V, &O)
	return O
}
