package spn

// Product represents a product node in an SPN.
type Product struct {
	Node
}

// NewProduct returns a new Product node.
func NewProduct() *Product {
	return &Product{}
}

// Type returns the type of this node.
func (p *Product) Type() string { return "product" }

// Sc returns the scope of this node.
func (p *Product) Sc() []int {
	if p.sc == nil {
		for i := range p.ch {
			p.sc = append(p.sc, p.ch[i].Sc()...)
		}
	}
	return p.sc
}

// Value returns the value of this SPN given a set of valuations.
func (p *Product) Value(val VarSet) float64 {
	n := len(p.ch)
	ch := p.Ch()
	var v float64

	for i := 0; i < n; i++ {
		v += ch[i].Value(val)
	}

	return v
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

// AddChild adds a child.
func (p *Product) AddChild(c SPN) {
	p.ch = append(p.ch, c)
	c.AddParent(p)
}
