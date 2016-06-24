// Package spn contains the structure of an SPN.
package spn

// Product represents an SPN product node.
type Product struct {
	// Children nodes.
	ch []Node
	// Parent node.
	pa Node
	// Scope of this node.
	sc []int
}

// NewProduct returns an empty Product node with given parent.
func NewProduct(pa Node) *Product {
	p := &Product{}
	p.pa = pa
	p.sc = nil
	return p
}

// AddChild adds a new child to this product node.
func (p *Product) AddChild(c Node) {
	p.ch = append(p.ch, c)
	p.sc = nil
}

// Ch returns the set of children nodes.
func (p *Product) Ch() []Node { return p.ch }

// Pa returns the parent node.
func (p *Product) Pa() Node { return p.pa }

// Type returns the type of this node: 'product'.
func (p *Product) Type() string { return "product" }

// Sc returns the scope of this node.
func (p *Product) Sc() []int {
	if p.sc == nil {
		// We assume an SPN is, by the Gens-Domingos definition, decomposable (and thus consistent).
		// Therefore, all children nodes must have disjoint scopes pairwise.
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
func (p *Product) Value(valuation VarSet) float64 {
	var v float64 = 1
	n := len(p.ch)

	for i := 0; i < n; i++ {
		v *= (p.ch[i]).Value(valuation)
	}

	return v
}

// Max returns the MAP state given a valuation.
func (p *Product) Max(valuation VarSet) float64 {
	var v float64 = 1
	n := len(p.ch)

	for i := 0; i < n; i++ {
		v *= (p.ch[i]).Max(valuation)
	}

	return v
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (p *Product) ArgMax(valuation VarSet) (VarSet, float64) {
	argmax := make(VarSet)
	n := len(p.ch)
	pmap := 1.0

	for i := 0; i < n; i++ {
		chargs, chmap := p.ch[i].ArgMax(valuation)
		pmap *= chmap
		for k, v := range chargs {
			argmax[k] = v
		}
	}

	return argmax, pmap
}
