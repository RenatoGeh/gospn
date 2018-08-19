package io

import (
	"fmt"
	"os"

	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/spn"
)

// VarSetToPPM takes a state and draws according to the SPN that generated the instantiation.
func VarSetToPPM(filename string, state spn.VarSet, w, h, max int) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Could not create file [%s].\n", filename)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "P6\n%d %d\n%d\n", w, h, max)

	n := len(state)
	pixels := make([]int, n)
	//fmt.Printf("len(pixels)=%d\n", n)
	for varid, val := range state {
		//fmt.Printf("[%d] = %d\n", varid, val)
		pixels[varid] = val
	}

	for i := 0; i < n; i++ {
		if (i+1)%w == 0 {
			fmt.Fprintf(file, "\n")
		}
		p := pixels[i]
		common.DrawColorRGB(file, p, p, p)
		fmt.Fprintf(file, " ")
	}
}
