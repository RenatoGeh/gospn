package spn

// Node represents a node in an SPN.
type Node struct {
	// Parent nodes.
	pa []SPN
	// Children nodes.
	ch []SPN
	// Scope of this node.
	sc []int
	// Store last soft inference values.
	s float64
	// Store partial derivatives wrt parent.
	pnode float64
}

// An SPN is a node.
type SPN interface {
	// Value returns the value of this node given an instantiation.
	Value(val VarSet) float64
	// Max returns the MAP value of this node given an evidence.
	Max(val VarSet) float64
	// ArgMax returns the MAP value and state given an evidence.
	ArgMax(val VarSet) (VarSet, float64)
	// Ch returns the set of children of this node.
	Ch() []SPN
	// Pa returns the set of parents of this node.
	Pa() []SPN
	// Sc returns the scope of this node.
	Sc() []int
	// Type returns the type of this node.
	Type() string
	// AddChild adds a child to this node.
	AddChild(c SPN)
	// AddParent adds a parent to this node.
	AddParent(p SPN)
	// Stored returns the last stored soft inference value.
	Stored() float64
	// Derivative returns the partial derivative wrt its parent.
	Derivative() float64
	// Derive recursively derives this node and its children based on the last inference value.
	Derive()
	// Rootify signalizes this node is a root. The only change this does is set pnode=1.
	Rootify()
	// GenUpdate generatively updates weights given an eta learning rate.
	GenUpdate(eta float64)
}

// VarSet is a variable set specifying variables and their respective instantiations.
type VarSet map[int]int

// Value returns the value of this node given an instantiation. (virtual)
func (n *Node) Value(val VarSet) float64 {
	return -1
}

// Max returns the MAP value of this node given an evidence. (virtual)
func (n *Node) Max(val VarSet) float64 {
	return -1
}

// ArgMax returns the MAP value and state given an evidence. (virtual)
func (n *Node) ArgMax(val VarSet) (VarSet, float64) {
	return nil, -1
}

// Ch returns the set of children of this node.
func (n *Node) Ch() []SPN {
	return n.ch
}

// Pa returns the set of parents of this node.
func (n *Node) Pa() []SPN {
	return n.pa
}

// Sc returns the scope of this node.
func (n *Node) Sc() []int {
	return n.sc
}

// Type returns the type of this node.
func (n *Node) Type() string {
	return "node"
}

// AddChild adds a child to this node.
func (n *Node) AddChild(c SPN) {
	n.ch = append(n.ch, c)
	c.AddParent(n)
}

// AddParent adds a parent to this node.
func (n *Node) AddParent(p SPN) {
	n.pa = append(n.pa, p)
}

// Stored returns the last stored soft inference value.
func (n *Node) Stored() float64 {
	return n.s
}

// Derivative returns the derivative of this node.
func (n *Node) Derivative() float64 {
	return n.pnode
}

// Derive recursively derives this node and its children based on the last inference value.
func (n *Node) Derive() {}

// Rootify signalizes this node is a root. The only change this does is set pnode=1.
func (n *Node) Rootify() { n.pnode = 1 }

// GenUpdate generatively updates weights given an eta learning rate.
func (n *Node) GenUpdate(eta float64) {
	m := len(n.ch)

	for i := 0; i < m; i++ {
		n.ch[i].GenUpdate(eta)
	}
}
