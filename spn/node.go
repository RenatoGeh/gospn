package spn

import (
//"fmt"
//"github.com/RenatoGeh/gospn/common"
)

// DiscStorer stores relevant information for DiscUpdate.
type DiscStorer struct {
	// The SPN's root.
	s SPN
	// The correct VarSet.
	c VarSet
	// The expected VarSet.
	e VarSet
	// Whether to store.
	store bool
	// Correct key.
	ckey string
	// Expected key.
	ekey string
}

// NewDiscStorer creates a new DiscStorer.
func NewDiscStorer(s SPN, c, e VarSet, ckey, ekey string) *DiscStorer {
	return &DiscStorer{s, c, e, true, ckey, ekey}
}

// Store sets DiscStorer to store previously computed values.
func (ds *DiscStorer) Store(store bool) { ds.store = store }

// Correct returns the value of the stored SPN given a correct valuation.
func (ds *DiscStorer) Correct() float64 {
	if v, ok := ds.s.Stored(ds.ckey); ds.store && ok {
		return v
	}
	//ds.s.RResetDP(ds.ckey)
	l := ds.s.Stores()
	if l {
		ds.s.SetStore(false)
	}
	val := ds.s.Soft(ds.c, ds.ckey)
	if l {
		ds.s.SetStore(l)
	}
	return val
}

// Expected returns the value of the stored SPN given an expected valuation.
func (ds *DiscStorer) Expected() float64 {
	if v, ok := ds.s.Stored(ds.ekey); ds.store && ok {
		return v
	}
	//ds.s.RResetDP(ds.ekey)
	l := ds.s.Stores()
	if l {
		ds.s.SetStore(false)
	}
	val := ds.s.Soft(ds.e, ds.ekey)
	if l {
		ds.s.SetStore(l)
	}
	return val
}

// CorrectSet returns the correct VarSet.
func (ds *DiscStorer) CorrectSet() VarSet { return ds.c }

// ExpectedSet returns the expected VarSet.
func (ds *DiscStorer) ExpectedSet() VarSet { return ds.e }

// InfType is Inference Type (soft or hard).
type InfType int

const (
	// SOFT is soft inference
	SOFT InfType = iota
	// HARD is hard inference
	HARD
)

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
	Stored(key string) (float64, bool)
	// Store stores an SPN evaluation for DP reasons.
	Store(key string, val float64)
	// SetStore sets whether the SPN should start storing evaluations on the DP table.
	SetStore(s bool)
	// Stores returns whether this node stores.
	Stores() bool
	// Derive derives this node only.
	Derive(wkey, nkey, ikey string, mode InfType) int
	// RootDerive derives all nodes in a BFS fashion.
	RootDerive(wkey, nkey, ikey string, mode InfType)
	// Rootify signalizes this node is a root.
	Rootify(nkey string)
	// GenUpdate generatively updates weights given an eta learning rate.
	GenUpdate(eta float64, wkey string)
	// Storer returns DP table.
	Storer() map[string]float64
	// Common base for all soft inference methods.
	Soft(val VarSet, key string) float64
	// LSoft is Soft in logspace.
	LSoft(val VarSet, key string) float64
	// Normalizes the SPN.
	Normalize()
	// DiscUpdate discriminatively updates weights given an eta learning rate.
	DiscUpdate(eta float64, ds *DiscStorer, wckey, wekey string, mode InfType)
	// ResetDP resets a key on the DP table. If key is nil, resets everything.
	ResetDP(key string)
	// RResetDP recursively ResetDPs all children.
	RResetDP(key string)
	// L2 regularization weight penalty.
	L2() float64
	// SetL2 changes the L2 regularization weight penalty throughout all SPN.
	SetL2(float64)
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
func (n *Node) Stored(key string) (float64, bool) {
	if val, ok := n.s[key]; ok && n.stores {
		return val, true
	}
	return 0, false
}

// Store stores an SPN evaluation for DP reasons.
func (n *Node) Store(key string, val float64) {
	if !n.stores {
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
// Return 0 to stop BFS. Return -1 to BFS through all children. Return i>0 to BFS through child i.
func (n *Node) Derive(wkey, nkey, ikey string, mode InfType) int { return -1 }

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

// LSoft is Soft in logspace.
func (n *Node) LSoft(val VarSet, key string) float64 { return -1 }

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
		n.s = make(map[string]float64)
	} else {
		delete(n.s, key)
	}
}

// RResetDP recursively ResetDPs all children.
func (n *Node) RResetDP(key string) {
	m := len(n.ch)

	n.ResetDP(key)
	for i := 0; i < m; i++ {
		n.ch[i].RResetDP(key)
	}
}

// DiscUpdate discriminatively updates weights given an eta learning rate.
func (n *Node) DiscUpdate(eta float64, ds *DiscStorer, wckey, wekey string, mode InfType) {
	m := len(n.ch)

	for i := 0; i < m; i++ {
		n.ch[i].DiscUpdate(eta, ds, wckey, wekey, mode)
	}
}

// RootDerive derives all nodes in a BFS fashion.
func (n *Node) RootDerive(wkey, nkey, ikey string, mode InfType) {}

// Stores returns whether this node stores.
func (n *Node) Stores() bool { return n.stores }

// L2 regularization weight penalty.
func (n *Node) L2() float64 { return 0 }

// SetL2 changes the L2 regularization weight penalty throughout all SPN.
func (n *Node) SetL2(l float64) {
	m := len(n.ch)
	for i := 0; i < m; i++ {
		n.ch[i].SetL2(l)
	}
}
