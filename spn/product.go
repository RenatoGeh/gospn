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
func NewProduct() *Product {
	p := &Product{}
	p.pa, p.sc = nil, nil
	return p
}

// AddChild adds a new child to this product node.
func (p *Product) AddChild(c Node) {
	p.ch = append(p.ch, c)
	c.SetParent(p)
	p.sc = nil
}

// SetParent sets the parent node.
func (p *Product) SetParent(pa Node) { p.pa = pa }

// Ch returns the set of children nodes.
func (p *Product) Ch() []Node { return p.ch }

// Pa returns the parent node.
func (p *Product) Pa() Node { return p.pa }

// Type returns the type of this node: 'product'.
func (p *Product) Type() string { return "product" }

// Weights returns nil, since product nodes have no weights.
func (p *Product) Weights() []float64 { return nil }

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
	n := len(p.ch)
	ch := p.Ch()
	v := 0.0

	for i := 0; i < n; i++ {
		v += ch[i].Value(valuation)
	}
	//for i := 0; i < n; i++ {
	//vch := (p.ch[i]).Value(valuation)
	//cv[i] = vch
	//fmt.Printf("ch[%d] of type \"%s\" pa \"%s\": %f\n", i, p.ch[i].Type(), "product", vch)
	//v *= vch
	//}

	//v := utils.LogProd(cv...)

	//fmt.Printf("Value of product node: antiln(%f)=%f\n", v, utils.AntiLog(v))
	return v
}

// Max returns the MAP state given a valuation.
func (p *Product) Max(valuation VarSet) float64 {
	n := len(p.ch)
	v := 0.0

	for i := 0; i < n; i++ {
		//v *= (p.ch[i]).Max(valuation)
		v += (p.ch[i]).Max(valuation)
	}

	return v
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (p *Product) ArgMax(valuation VarSet) (VarSet, float64) {
	argmax := make(VarSet)
	n := len(p.ch)
	v := 0.0

	// Only a product node may have leaves as children. We must iterate through all children and
	// recurse through them. Since a product node's children must have disjoint scopes, each child
	// has different variables to be considered. If child is a leaf, then return its corresponding
	// MAP (trivial base case). Else we have a node with scope size greater than one and we must
	// recurse once again.
	for i := 0; i < n; i++ {
		chargs, chmap := p.ch[i].ArgMax(valuation)
		v += chmap
		//pmap *= chmap
		for k, val := range chargs {
			argmax[k] = val
		}
	}

	//fmt.Printf("Product: %f %v\n", utils.AntiLog(v), argmax)
	return argmax, v
}
