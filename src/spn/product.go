// Package spn contains the structure of an SPN.
package spn

// Product represents an SPN product node.
type Product struct {
	// Children nodes.
	ch []Node
	// Parent node.
	pa Node
}

// NewProduct returns an empty Product node with given parent.
func NewProduct(pa Node) *Product {
	p := &Product{}
	p.pa = pa
	return p
}

// AddChild adds a new child to this product node.
func (p *Product) AddChild(c Node) {
	p.ch = append(p.ch, c)
}

// Ch returns the set of children nodes.
func (p *Product) Ch() []Node { return p.ch }

// Pa returns the parent node.
func (p *Product) Pa() Node { return p.pa }

// Type returns the type of this node: 'product'.
func (p *Product) Type() string { return "product" }

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
