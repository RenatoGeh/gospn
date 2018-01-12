// Package containing graph manipulation functions.
package spn

type graphTest struct {
	Node
	i int
}

func (g *graphTest) Add(n ...SPN) {
	g.ch = append(g.ch, n...)
}

func (g *graphTest) Type() string {
	if g.ch == nil {
		return "leaf"
	}
	return "node"
}
