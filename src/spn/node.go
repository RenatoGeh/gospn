// Package spn contains the structure of an SPN.
package spn

// A node represents a node in a DAG. There can only be three types of nodes: a univariate
// distribution node, a sum node and a product node.
type Node interface {
	// Node value given a valuation.
	Value(valuation VarSet) float64
	// Returns the MAP state given a valuation.
	Max(valuation VarSet) float64
	// Set of children.
	Ch() []Node
	// Parent node. If returns nil, then it is a root node.
	Pa() Node
	// Node type: 'leaf', 'sum' or 'product'.
	Type() string
}
