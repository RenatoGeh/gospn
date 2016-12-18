package sspn

// Node represents a node in an SPN.
type Node struct {
	// Parent nodes.
	pa []*Node
	// Children nodes.
	ch []*Node
	// Scope of this node.
	sc []int
}

// An SPN is a node.
type SPN Node

// VarSet is a variable set specifying variables and their respective instantiations.
type VarSet map[int]int

// Struct methods are statically linked, and thus are faster than using interfaces for methods that
// are ideally virtual. Interface method evaluation is done in runtime. We use this dirty trick of
// using useless ideally-virtual methods for this reason. Since each node is "self-contained" in
// the sense that it never contains second-level methods, this practice is safe.

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
func (n *Node) Ch() []*Node {
	return n.ch
}

// Pa returns the set of parents of this node.
func (n *Node) Pa() []*Node {
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
func (n *Node) AddChild(c *Node) {
	n.ch = append(n.ch, c)
	c.pa = append(n.pa, n)
}

// Weights returns weights if sum product. Returns nil otherwise.
func (n *Node) Weights() []float64 {
	return nil
}
