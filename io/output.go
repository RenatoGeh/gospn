package io

import (
	"fmt"
	"os"
	"path/filepath"

	common "github.com/RenatoGeh/gospn/common"
	spn "github.com/RenatoGeh/gospn/spn"
	utils "github.com/RenatoGeh/gospn/utils"
)

// CmplType indicates which type of image completion are we referring to.
type CmplType string

const (
	// Top image completion.
	Top CmplType = "TOP"
	// Bottom image completion.
	Bottom CmplType = "BOTTOM"
	// Left image completion.
	Left CmplType = "LEFT"
	// Right image completion.
	Right CmplType = "RIGHT"
)

// Orientations contains all orientations.
var Orientations = []CmplType{Top, Bottom, Left, Right}

// DrawGraphTools creates a file filename and draws an SPN spn in graph-tools. The resulting file
// is a python source code that outputs a PNG image of the graph.
func DrawGraphTools(filename string, s spn.SPN) {
	file, err := os.Create(filename)

	if err != nil {
		fmt.Printf("Error. Could not create file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	outname := utils.StringConcat(filename[0:len(filename)-len(filepath.Ext(filename))], ".png")

	fmt.Fprintf(file, "from graph_tool.all import *\n\n")
	fmt.Fprintf(file, "g = Graph(directed=True)\n")
	fmt.Fprintf(file, "vcolors = g.new_vertex_property(\"string\")\n")
	fmt.Fprintf(file, "vnames = g.new_vertex_property(\"string\")\n")
	fmt.Fprintf(file, "enames = g.new_edge_property(\"string\")\n\n")
	fmt.Fprintf(file, "def add_node(name, type):\n\tv=g.add_vertex()\n\tvnames[v]=name\n\t"+
		"vcolors[v]=type\n\treturn v\n\n")
	fmt.Fprintf(file, "def add_edge(o, t, name):\n\te=g.add_edge(o, t)\n\tenames[e]=name\n\treturn e\n\n")
	fmt.Fprintf(file, "def add_edge_nameless(o, t):\n\te=g.add_edge(o, t)\n\treturn e\n\n\n")

	// If the SPN is itself an univariate distribution, create a graph with a single node.
	if s.Type() == "leaf" {
		fmt.Fprintf(file, "add_node(\"X\")\n\n")
		fmt.Fprintf(file, "g.vertex_properties[\"name\"]=vnames\n")
		fmt.Fprintf(file, "g.vertex_properties[\"color\"]=vcolors\n")
		fmt.Fprintf(file, "\ngraph_draw(g, vertex_text=g.vertex_properties[\"name\"], "+
			"edge_text=enames, vertex_fill_color=g.vertex_properties[\"color\"], output=\"%s\")\n",
			outname)
		return
	}

	// Else, BFS the SPN and write nodes to filename.
	nvars, nsums, nprods := 0, 0, 0
	queue := common.QueueBFSPair{}
	queue.Enqueue(&common.BFSPair{Spn: s, Pname: "", Weight: -1.0})
	for !queue.Empty() {
		currpair := queue.Dequeue()
		curr, pname, pw := currpair.Spn, currpair.Pname, currpair.Weight
		ch := curr.Ch()
		nch := len(ch)

		name := "N"
		currt := curr.Type()

		// In case it is a sum node. Else product node.
		if currt == "sum" {
			name = fmt.Sprintf("S%d", nsums)
			fmt.Fprintf(file, "%s = add_node(\"+\", \"#ff3300\")\n", name)
			nsums++
		} else if currt == "product" {
			name = fmt.Sprintf("P%d", nprods)
			fmt.Fprintf(file, "%s = add_node(\"*\", \"#669900\")\n", name)
			nprods++
		}

		// If pname is empty, then it is the root node. Else, link parent node to current node.
		if pname != "" {
			if pw >= 0 {
				fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", pname, name, pw)
			} else {
				fmt.Fprintf(file, "add_edge_nameless(%s, %s)\n", pname, name)
			}
		}

		w := curr.Weights()
		// For each children, run the BFS.
		for i := 0; i < nch; i++ {
			c := ch[i]

			// If leaf, then simply write to the graphviz dot file. Else, recurse the BFS.
			if c.Type() == "leaf" {
				cname := fmt.Sprintf("X%d", nvars)
				fmt.Fprintf(file, "%s = add_node(\"X_%d\", \"#0066ff\")\n", cname, c.Sc()[0])
				nvars++
				if currt == "sum" {
					fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", name, cname, w[i])
				} else {
					fmt.Fprintf(file, "add_edge_nameless(%s, %s)\n", name, cname)
				}
			} else {
				tw := -1.0
				if w != nil {
					tw = w[i]
				}
				queue.Enqueue(&common.BFSPair{Spn: c, Pname: name, Weight: tw})
			}
		}
	}

	fmt.Fprintf(file, "g.vertex_properties[\"name\"]=vnames\n")
	fmt.Fprintf(file, "g.vertex_properties[\"color\"]=vcolors\n")
	//fmt.Fprintf(file, "\ngraph_draw(g, vertex_text=g.vertex_properties[\"name\"], "+
	//"edge_text=enames, vertex_fill_color=g.vertex_properties[\"color\"], "+
	//"output_size=[16384, 16384], output=\"%s\", bg_color=[1, 1, 1, 1])\n", outname)
	fmt.Fprintf(file, "\ngraph_draw(g, "+
		"edge_text=enames, vertex_fill_color=g.vertex_properties[\"color\"], "+
		"output_size=[16384, 16384], output=\"%s\", bg_color=[1, 1, 1, 1])\n", outname)
}

// DrawGraph creates a file filename and draws an SPN spn in Graphviz dot.
func DrawGraph(filename string, s spn.SPN) {
	file, err := os.Create(filename)

	if err != nil {
		fmt.Printf("Error. Could not create file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	fmt.Fprintf(file, "graph {\n")

	// If the SPN is itself an univariate distribution, create a graph with a single node.
	if s.Type() == "leaf" {
		fmt.Fprintf(file, "X1 [label=<X<sub>1</sub>>,shape=circle];\n")
		fmt.Fprintf(file, "}")
		file.Close()
		return
	}

	// Else, BFS the SPN and write nodes to filename.
	nvars, nsums, nprods := 0, 0, 0
	queue := common.QueueBFSPair{}
	queue.Enqueue(&common.BFSPair{Spn: s, Pname: "", Weight: -1.0})
	for !queue.Empty() {
		currpair := queue.Dequeue()
		curr, pname, pw := currpair.Spn, currpair.Pname, currpair.Weight
		ch := curr.Ch()
		nch := len(ch)

		name := "N"
		currt := curr.Type()

		// In case it is a sum node. Else product node.
		if currt == "sum" {
			name = fmt.Sprintf("S%d", nsums)
			fmt.Fprintf(file, "%s [label=\"+\",shape=circle];\n", name)
			nsums++
		} else if currt == "product" {
			name = fmt.Sprintf("P%d", nprods)
			fmt.Fprintf(file, "%s [label=<&times;>,shape=circle];\n", name)
			nprods++
		}

		// If pname is empty, then it is the root node. Else, link parent node to current node.
		if pname != "" {
			if pw >= 0 {
				fmt.Fprintf(file, "%s -- %s [label=\"%.3f\"];\n", pname, name, pw)
			} else {
				fmt.Fprintf(file, "%s -- %s\n", pname, name)
			}
		}

		w := curr.Weights()
		// For each children, run the BFS.
		for i := 0; i < nch; i++ {
			c := ch[i]

			// If leaf, then simply write to the graphviz dot file. Else, recurse the BFS.
			if c.Type() == "leaf" {
				cname := fmt.Sprintf("X%d", nvars)
				fmt.Fprintf(file, "%s [label=<X<sub>%d</sub>>,shape=circle];\n", cname, c.Sc()[0])
				nvars++
				if currt == "sum" {
					fmt.Fprintf(file, "%s -- %s [label=\"%.3f\"]\n", name, cname, w[i])
				} else {
					fmt.Fprintf(file, "%s -- %s\n", name, cname)
				}
			} else {
				tw := -1.0
				if w != nil {
					tw = w[i]
				}
				queue.Enqueue(&common.BFSPair{Spn: c, Pname: name, Weight: tw})
			}
		}
	}

	fmt.Fprintf(file, "}")
}
