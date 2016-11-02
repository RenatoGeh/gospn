package common

import (
	spn "github.com/RenatoGeh/gospn/spn"
)

// BFSPair (Breadth-First Search Pair) is a tuple (SPN, string). See io/output:DrawGraph for more
// information.
type BFSPair struct {
	Spn    spn.SPN
	Pname  string
	Weight float64
}
