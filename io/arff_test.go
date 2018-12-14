package io

import (
	"fmt"
	"testing"
)

func TestParseArff(t *testing.T) {
	name, sc, vals, labels := ParseArff("test.arff")

	fmt.Printf("Relation name: %s.\n", name)
	fmt.Println("Scope:")
	for k, v := range sc {
		fmt.Printf("  Variable [%d] = { varid: %d, categories: %d, name: %s }\n", k, v.Varid, v.Categories, v.Name)
	}
	for i := range labels {
		fmt.Printf("Variable [%d] labels:\n", i)
		for k, v := range labels[i] {
			fmt.Printf("  Key: %d, Value: %s\n", v, k)
		}
	}
	fmt.Println("Data:")
	for i := range vals {
		fmt.Printf("(%d) = {", i)
		for j, v := range vals[i] {
			fmt.Printf(" [%d]=%d", j, v)
		}
		fmt.Println(" }")
	}
}
