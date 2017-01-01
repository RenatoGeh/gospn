package spn

// Node represents a node in an SPN.
type Node struct {
	// Parent nodes.
	pa []SPN
	// Children nodes.
	ch []SPN
	// Scope of this node.
	sc map[int]int
	// Stores inference values.
	s map[string]float64
	// Signals this node to be the root of the SPN.
	root bool
	// Whether to store in DP table or not.
	stores bool
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
	Sc() map[int]int
	// Type returns the type of this node.
	Type() string
	// AddChild adds a child to this node.
	AddChild(c SPN)
	// AddParent adds a parent to this node.
	AddParent(p SPN)
	// Stored returns the stored soft inference value from the given key.
	Stored(key string) float64
	// Store stores an SPN evaluation for DP reasons.
	Store(key string, val float64)
	// SetStore sets whether the SPN should start storing evaluations on the DP table.
	SetStore(s bool)
	// Derive recursively derives this node and its children based on the last inference value.
	// Stores derivative values under key.
	Derive(wkey, nkey, ikey string)
	// Rootify signalizes this node is a root.
	Rootify(nkey string)
	// GenUpdate generatively updates weights given an eta learning rate.
	GenUpdate(eta float64, wkey string)
	// Storer returns DP table.
	Storer() map[string]float64
	// Common base for all soft inference methods.
	Soft(val VarSet, key string) float64
	// Normalizes the SPN.
	Normalize()
	// DiscUpdate discriminatively updates weights given an eta learning rate.
	DiscUpdate(eta, correct, expected float64, wckey, wekey string)
	// ResetDP resets a key on the DP table. If key is nil, resets everything.
	ResetDP(key string)
	// RResetDP recursively ResetDPs all children.
	RResetDP(key string)
}

// VarSet is a variable set specifying variables and their respective instantiations.
type VarSet map[int]int

// NewNode creates a new node value.
func NewNode(scope ...int) Node {
	m := len(scope)
	lsc := make(map[int]int)
	for i := 0; i < m; i++ {
		lsc[scope[i]] = scope[i]
	}
	return Node{sc: lsc, s: make(map[string]float64)}
}

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
func (n *Node) Sc() map[int]int {
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

// Stored returns the stored soft inference value from the given key.
func (n *Node) Stored(key string) float64 {
	if n.stores {
		return n.s[key]
	}
	return -1
}

// Store stores an SPN evaluation for DP reasons.
func (n *Node) Store(key string, val float64) {
	if n.stores {
		return
	}

	if key == "" {
		key = "default"
	}
	n.s[key] = val
}

// SetStore sets whether the SPN should start storing evaluations on the DP table.
func (n *Node) SetStore(s bool) {
	n.stores = s
	m := len(n.ch)

	for i := 0; i < m; i++ {
		n.ch[i].SetStore(s)
	}
}

// Derive recursively derives this node and its children based on the last inference value.
func (n *Node) Derive(wkey, nkey, ikey string) {}

// Rootify signalizes this node is a root.
func (n *Node) Rootify(nkey string) {
	n.Store(nkey, 1)
	n.root = true
}

// GenUpdate generatively updates weights given an eta learning rate.
func (n *Node) GenUpdate(eta float64, wkey string) {
	m := len(n.ch)

	for i := 0; i < m; i++ {
		n.ch[i].GenUpdate(eta, wkey)
	}
}

// Storer returns DP table.
func (n *Node) Storer() map[string]float64 { return n.s }

// Soft is a common base for all soft inference methods.
func (n *Node) Soft(val VarSet, key string) float64 { return -1 }

// Normalize normalizes the SPN's weights.
func (n *Node) Normalize() {
	m := len(n.ch)

	for i := 0; i < m; i++ {
		n.ch[i].Normalize()
	}
}

// ResetDP resets a key on the DP table. If key is nil, resets everything.
func (n *Node) ResetDP(key string) {
	if key == "" {
		for k := range n.s {
			n.s[k] = -1
		}
	} else {
		n.s[key] = -1
	}
}

// RResetDP recursively ResetDPs all children.
func (n *Node) RResetDP(key string) {
	m := len(n.ch)

	n.ResetDP(key)
	for i := 0; i < m; i++ {
		n.ch[i].ResetDP(key)
	}
}

// DiscUpdate discriminatively updates weights given an eta learning rate.
func (n *Node) DiscUpdate(eta, correct, expected float64, wckey, wekey string) {
	m := len(n.ch)

	for i := 0; i < m; i++ {
		n.ch[i].DiscUpdate(eta, correct, expected, wckey, wekey)
	}
}
