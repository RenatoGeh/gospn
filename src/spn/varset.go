// Package spn contains the structure of an SPN.
package spn

// A VarSet is a variable set specifying variables and their respective instantiations.
// It's just a map with int keys and int values.
// Each key is a varid. Its associated value is the variable's value.
type VarSet map[int]int
