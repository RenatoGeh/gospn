package io

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	spn "github.com/RenatoGeh/gospn/src/spn"
	utils "github.com/RenatoGeh/gospn/src/utils"
)

// Creates a file filename and draws an SPN spn in Graphviz dot.
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
	queue := utils.QueueBFSPair{}
	queue.Enqueue(&utils.BFSPair{s, "", -1.0})
	for !queue.Empty() {
		currpair := queue.Dequeue()
		curr, pname, pw := currpair.Spn, currpair.Pname, currpair.Weight
		ch := s.Ch()
		nch := len(ch)

		name := "N"
		currt := curr.Type()

		// In case it is a sum node. Else product node.
		if currt == "sum" {
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
			fmt.Fprintf(file, "%s -- %s [label=\"%.3f\"];\n", pname, name, pw)
		}

		w := curr.Weights()
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
				tw := -1.0
				if w != nil {
					tw = w[i]
				}
				queue.Enqueue(&utils.BFSPair{c, name, tw})
			}
		}
	}

	fmt.Fprintf(file, "}")
}

// PBMFToData: PBM Folder to Data file. Each class is in a subfolder of dirname. dname is the
// output file. Arg dirname must be an absolute path. Arg dname must be the filename only.
func PBMFToData(dirname, dname string) {
	sdir, err := os.Open(dirname)

	if err != nil {
		fmt.Printf("Error. Could not open superdirectory [%s].\n", dirname)
		panic(err)
	}
	defer sdir.Close()

	subdirs, err := sdir.Readdirnames(-1)

	if err != nil {
		fmt.Printf("Error. Could not extract subdirectories from directory [%s].\n", dirname)
		panic(err)
	}

	nsdirs := len(subdirs)
	tpath := utils.StringConcat(dirname, "/")
	var mrkrm []int
	// Reserved dirname compiled for output. Also remove non-dirs.
	for i := 0; i < nsdirs; i++ {
		// m marks the spot.
		// Since for every removed item the slice shrinks by one, we keep track of the indices by
		// taking into account the subslices "translated" at the right moment.
		if subdirs[i] == "compiled" {
			var m int = i
			if len(mrkrm) > 0 {
				m = i - 1
			}
			mrkrm = append(mrkrm, m)
		} else if fi, _ := os.Stat(utils.StringConcat(tpath, subdirs[i])); !fi.IsDir() {
			var m int = i
			if len(mrkrm) > 0 {
				m = i - 1
			}
			mrkrm = append(mrkrm, m)
		}
	}

	// Remove marked elements.
	for i := 0; i < len(mrkrm); i++ {
		j := mrkrm[i]
		subdirs, nsdirs = append(subdirs[:j], subdirs[j+1:]...), nsdirs-1
	}

	// Memorize all subfiles.
	var instreams []*bufio.Scanner = nil
	var labels []int = nil
	for i := 0; i < nsdirs; i++ {
		sd, err := os.Open(utils.StringConcat(tpath, subdirs[i]))

		if err != nil {
			fmt.Printf("Error. Failed to open subdirectory [%s].\n", subdirs[i])
			panic(err)
		}
		defer sd.Close()

		sf, err := sd.Readdirnames(-1)

		if err != nil {
			fmt.Printf("Error. Failed to read files under [%s].\n", subdirs[i])
			panic(err)
		}

		spath := utils.StringConcat(utils.StringConcat(tpath, subdirs[i]), "/")
		nsf := len(sf)
		for j := 0; j < nsf; j++ {
			f, err := os.Open(utils.StringConcat(spath, sf[j]))

			if err != nil {
				fmt.Printf("Error. Failed to open file [%s%s].\n", spath, sf[j])
				panic(err)
			}
			defer f.Close()

			fmt.Printf("Name: %s\n", utils.StringConcat(spath, sf[j]))

			instreams = append(instreams, bufio.NewScanner(f))
			labels = append(labels, i)
		}
	}

	// Create compiled folder.
	cmpname, err := filepath.Abs(dirname)

	if err != nil {
		fmt.Printf("Error retrieving path [%s].\n", dirname)
		panic(err)
	}

	cmpname = utils.StringConcat(cmpname, "/compiled")
	if _, err := os.Stat(cmpname); os.IsNotExist(err) {
		os.Mkdir(cmpname, 0777)
	}

	cmpname = utils.StringConcat(cmpname, "/")
	out, err := os.Create(utils.StringConcat(cmpname, dname))

	if err != nil {
		fmt.Printf("Error creating output file [%s/%s].\n", cmpname, dname)
		panic(err)
	}
	defer out.Close()

	// Deal with P1.
	instreams[0].Scan()

	// Read width and height.
	w, h := -1, -1
	instreams[0].Scan()
	fmt.Sscanf(instreams[0].Text(), "%d %d", &w, &h)
	fmt.Printf("w: %d, h: %d\n", w, h)

	nin := len(instreams)
	// Move stream pointer to the right position.
	for i := 1; i < nin; i++ {
		instreams[i].Scan()
		instreams[i].Scan()
	}

	// Declare variables to data file.
	tt := w * h
	for i := 0; i < tt; i++ {
		fmt.Fprintf(out, "var %d 2\n", i)
	}

	for i := 0; i < nin; i++ {
		stream := instreams[i]

		for stream.Scan() {
			line := stream.Text()
			nline := len(line)
			for j := 0; j < nline; j++ {
				fmt.Fprintf(out, "%c ", line[j])
			}
		}

		fmt.Fprintf(out, "%d\n", labels[i])
	}
}

// PBM to Data file. If class is true, it's a classifying problem and will label as class.
func PBMToData(dirname, dname string, class int) {
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

		if class > 0 {
			fmt.Fprintf(out, "%d\n", class)
		} else {
			fmt.Fprintf(out, "\n")
		}
	}
}
