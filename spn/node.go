package spn

import (
	"github.com/RenatoGeh/gospn/learn/parameters"
)

// Node represents a node in an SPN.
type Node struct {
	// Children nodes.
	ch []SPN
	// Scope of this node.
	sc []int
}

// An SPN is a node.
type SPN interface {
	// Value returns the value of this node given an instantiation.
	Value(val VarSet) float64
	// Compute returns the soft value of this node's type given the children's values.
	Compute(cv []float64) float64
	// Max returns the MAP value of this node given an evidence.
	Max(val VarSet) float64
	// ArgMax returns the MAP value and state given an evidence.
	ArgMax(val VarSet) (VarSet, float64)
	// Ch returns the set of children of this node.
	Ch() []SPN
	// Sc returns the scope of this node.
	Sc() []int
	// Type returns the type of this node.
	Type() string
	// SubType returns the subtype of this node.
	SubType() string
	// AddChild adds a child to this node.
	AddChild(c SPN)
	// Returns the height of the graph.
	Height() int
	// Parameters returns the parameters of this object.
	Parameters() *parameters.P

	rawSc() []int
	setRawSc([]int)
}

func (n *Node) rawSc() []int {
	return n.sc
}

func (n *Node) setRawSc(sc []int) {
	n.sc = sc
}

// VarSet is a variable set specifying variables and their respective instantiations.
type VarSet map[int]int

// Dataset is a dataset indexed by instances.
type Dataset []map[int]int

// Value returns the value of this node given an instantiation. (virtual)
func (n *Node) Value(val VarSet) float64 {
	return -1
}

// Compute returns the soft value of this node's type given the children's values.
func (n *Node) Compute(cv []float64) float64 {
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

// Sc returns the scope of this node.
func (n *Node) Sc() []int {
	return n.sc
}

// Type returns the type of this node.
func (n *Node) Type() string {
	return "node"
}

// SubType returns the subtype of this node.
func (n *Node) SubType() string { return "node" }

// AddChild adds a child to this node.
func (n *Node) AddChild(c SPN) {
	n.ch = append(n.ch, c)
}

// Returns the height of the graph.
func (n *Node) Height() int {
	nc := len(n.ch)
	if nc > 0 {
		v := 0
		for i := 0; i < nc; i++ {
			t := n.ch[i].Height()
			if t > v {
				v = t
			}
		}
		return (v + 1)
	} else {
		return 0
	}
}

// Parameters returns the parameters of this object. If no bound parameter is found, binds default
// parameter values and returns.
func (n *Node) Parameters() *parameters.P {
	p, e := parameters.Retrieve(n)
	if !e {
		p = parameters.Default()
		parameters.Bind(n, p)
	}
	return p
}
