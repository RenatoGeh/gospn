package utils

import (
	spn "github.com/RenatoGeh/gospn/src/spn"
)

// A BFSPair (Breadth-First Search Pair) is a tuple (SPN, string). See io/output:DrawGraph for more
// information.
type BFSPair struct {
	Spn    spn.SPN
	Pname  string
	Weight float64
}
