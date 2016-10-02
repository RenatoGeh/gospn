package io

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	common "github.com/RenatoGeh/gospn/src/common"
	spn "github.com/RenatoGeh/gospn/src/spn"
	utils "github.com/RenatoGeh/gospn/src/utils"
)

// PBMFToData (PBM Folder to Data file). Each class is in a subfolder of dirname. dname is the
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
			m := i
			if len(mrkrm) > 0 {
				m = i - 1
			}
			mrkrm = append(mrkrm, m)
		} else if fi, _ := os.Stat(utils.StringConcat(tpath, subdirs[i])); !fi.IsDir() {
			m := i
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
	var instreams []*bufio.Scanner
	var labels []int
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
	fmt.Fprintf(out, "var %d %d\n", tt, nsdirs)

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

// PBMToData (PBM to Data file). If class is true, it's a classifying problem and will label as
// class.
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

// PBMFToEvidence (PBM file to evidence).
func PBMFToEvidence(dirname, dname string) {
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
		if fi, _ := os.Stat(utils.StringConcat(tpath, subdirs[i])); !fi.IsDir() ||
			subdirs[i] == "compiled" {
			m := i
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

	// Marks which class labels they are supposed to be classified as. Each int is the index of each
	// class label.
	var slabels []int

	// Memorize all subfiles.
	var instreams []*bufio.Scanner
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

			instreams = append(instreams, bufio.NewScanner(f))
		}
		slabels = append(slabels, nsf*i)
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

	nin := len(instreams)
	// Move stream pointer to the right position.
	for i := 1; i < nin; i++ {
		instreams[i].Scan()
		instreams[i].Scan()
	}

	// Declare labels.
	fmt.Fprintf(out, "labels %d ", len(slabels))
	for i, nslabels := 0, len(slabels); i < nslabels; i++ {
		if i == nslabels-1 {
			fmt.Fprintf(out, "%d\n", slabels[i])
		} else {
			fmt.Fprintf(out, "%d ", slabels[i])
		}
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

		fmt.Fprintf(out, "\n")
	}
}

// VarSetToPBM takes a state and draws according to the SPN that generated the instantiation.
func VarSetToPBM(filename string, state spn.VarSet, w, h int) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Could not create file [%s].\n", filename)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "P1\n%d %d\n", w, h)

	n := len(state)
	pixels := make([]int, n)
	for varid, val := range state {
		pixels[varid] = val
	}

	for i := 0; i < n; i++ {
		if i%71 == 0 {
			fmt.Fprintf(file, "\n")
		}
		fmt.Fprintf(file, "%d", pixels[i])
	}
}

// ImgCmplToPPM creates a new file distinguishing the original part of the image from the
// completion done by the SPN and indicated by typ.
func ImgCmplToPPM(filename string, orig, cmpl spn.VarSet, typ CmplType, w, h int) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Could not create file [%s].\n", filename)
		return
	}
	defer file.Close()

	var mid func(int) bool
	if typ == Top || typ == Bottom {
		h++
		mid = func(p int) bool {
			q := w * (h / 2)
			return p >= q && p < q+w
		}
	} else {
		w++
		mid = func(p int) bool {
			return p%w == w/2
		}
	}

	fmt.Fprintf(file, "P3\n%d %d\n255\n", w, h)

	n, j := w*h, 0
	for i := 0; i < n; i++ {
		if mid(i) {
			common.DrawColor(file, common.Red)
			goto cleanup
		} else if v, eo := orig[j]; eo {
			if v == 1 {
				common.DrawColor(file, common.Blue)
			} else {
				common.DrawColor(file, common.White)
			}
		} else {
			u, _ := cmpl[j]
			if u == 1 {
				common.DrawColor(file, common.Green)
			} else {
				common.DrawColor(file, common.White)
			}
		}
		j++
	cleanup:
		fmt.Fprintf(file, " ")
		if i != 0 && i%w == 0 {
			fmt.Fprintf(file, "\n")
		}
	}
}
