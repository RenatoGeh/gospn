package io

import (
	"bufio"
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

	// If the SPN is itself an univariate distribution, create a graph with a single node.
	if spn.Type() == "leaf" {
		fmt.Fprintf(file, "X1 [label=<X<sub>1</sub>>,shape=circle];\n")
		fmt.Fprintf(file, "}")
		file.Close()
		return
	}

	// Else, BFS the SPN and write nodes to filename.
	nvars, nsums, nprods := 0, 0, 0
	queue := utils.QueueBFSPair{}
	queue.Enqueue(&utils.BFSPair{spn, ""})
	for !queue.Empty() {
		currpair := queue.Dequeue()
		curr, pname := currpair.Spn, currpair.Pname
		ch := spn.Ch()
		nch := len(ch)

		name := "N"

		// In case it is a sum node. Else product node.
		if curr.Type() == "sum" {
			name = fmt.Sprintf("S%d", nsums)
			fmt.Fprintf(file, "%s [label="+",shape=circle];\n", name, nsums)
			nsums++
		} else {
			name = fmt.Sprintf("P%d", nprods)
			fmt.Fprintf(file, "%s [label=<&times;>,shape=circle];\n", name, nprods)
			nprods++
		}

		// If pname is empty, then it is the root node. Else, link parent node to current node.
		if pname != "" {
			fmt.Fprintf(file, "%s -- %s\n", pname, name)
		}

		// For each children, run the BFS.
		for i := 0; i < nch; i++ {
			c := ch[i]

			// If leaf, then simply write to the graphviz dot file. Else, recurse the BFS.
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
}

func PBMToData(dirname, dname string) {
	dir, err := os.Open(dirname)

	if err != nil {
		fmt.Printf("Error. Could not open directory [%s].\n", dirname)
		panic(err)
	}
	defer dir.Close()

	filenames, err := dir.Readdirnames(-1)

	if err != nil {
		fmt.Printf("Error. Could not extract filenames from directory [%s].\n", dirname)
		panic(err)
	}

	in := make([]*os.File, len(filenames))
	nin := len(in)
	instream := make([]*bufio.Scanner, nin)

	tdir := utils.StringConcat(dirname, "/")
	for i := 0; i < nin; i++ {
		inname := utils.StringConcat(tdir, filenames[i])
		in[i], err = os.Open(inname)

		if err != nil {
			fmt.Printf("Error. Could not open file [%s]\n", inname)
			panic(err)
		}
		defer in[i].Close()

		instream[i] = bufio.NewScanner(in[i])
	}

	out, err := os.Create(dname)

	if err != nil {
		fmt.Printf("Error. Could not create file [%s].\n", dname)
		panic(err)
	}
	defer out.Close()

	// Deal with P1.
	instream[0].Scan()

	// Read width and height.
	w, h := -1, -1
	instream[0].Scan()
	fmt.Sscanf(instream[0].Text(), "%d %d", &w, &h)

	// Move stream pointer to the right position.
	for i := 1; i < nin; i++ {
		instream[i].Scan()
		instream[i].Scan()
	}

	// Declare variables to data file.
	tt := w * h
	for i := 0; i < tt; i++ {
		fmt.Fprintf(out, "var %d 2\n", i)
	}

	for i := 0; i < nin; i++ {
		stream := instream[i]

		for stream.Scan() {
			line := stream.Text()
			nline := len(line)
			for j := 0; j < nline; j++ {
				fmt.Fprintf(out, "%c ", line[j])
			}
		}
		fmt.Fprintf(out, "\n")
	}
}
