package io

import (
	"fmt"
	"os"

	spn "github.com/RenatoGeh/gospn/src/spn"
	utils "github.com/RenatoGeh/gospn/src/utils"
)

// Creates a file filename and draws an SPN spn in Graphviz dot.
func DrawGraph(filename string, spn spn.SPN) {
	file, err := os.Create(filename)

	if err != nil {
		fmt.Printf("Error. Could not create file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	fmt.Fprintf(file, "graph {\n")

	if spn.Type() == "leaf" {
		fmt.Fprintf(file, "X1 [label=<X<sub>1</sub>>,shape=circle];\n")
		fmt.Fprintf(file, "}")
		file.Close()
		return
	}

	// BFS the SPN and writes nodes to filename.
	nvars, nsums, nprods := 0, 0, 0
	queue := utils.QueueBFSPair{}
	queue.Enqueue(&utils.BFSPair{spn, ""})
	for !queue.Empty() {
		currpair := queue.Dequeue()
		curr, pname := currpair.Spn, currpair.Pname
		ch := spn.Ch()
		nch := len(ch)

		name := "N"

		if curr.Type() == "sum" {
			name = fmt.Sprintf("S%d", nsums)
			fmt.Fprintf(file, "%s [label="+",shape=circle];\n", name, nsums)
			nsums++
		} else {
			name = fmt.Sprintf("P%d", nprods)
			fmt.Fprintf(file, "%s [label=<&times;>,shape=circle];\n", name, nprods)
			nprods++
		}

		if pname != "" {
			fmt.Fprintf(file, "%s -- %s\n", pname, name)
		}

		for i := 0; i < nch; i++ {
			c := ch[i]

			if c.Type() == "leaf" {
				cname := fmt.Sprintf("X%d", nvars)
				fmt.Fprintf(file, "%s [label=<X<sub>%d</sub>>,shape=circle];\n", cname, nvars)
				nvars++
				fmt.Fprintf(file, "%s -- %s\n", name, cname)
			} else {
				queue.Enqueue(&utils.BFSPair{c, name})
			}
		}
	}

	fmt.Fprintf(file, "}")

	file.Close()
}
