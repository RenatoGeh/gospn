package io

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	common "github.com/RenatoGeh/gospn/common"
	spn "github.com/RenatoGeh/gospn/spn"
	utils "github.com/RenatoGeh/gospn/utils"
)

// PGMFToData (PGM Folder to Data file). Each class is in a subfolder of dirname. dname is the
// output file. Arg dirname must be an absolute path. Arg dname must be the filename only.
func PGMFToData(dirname, dname string) (int, int, int) {
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
	cmarks := 1
	// Reserved dirname compiled for output. Also remove non-dirs.
	for i := 0; i < nsdirs; i++ {
		// m marks the spot.
		// Since for every removed item the slice shrinks by one, we keep track of the indices by
		// taking into account the subslices "translated" at the right moment.
		if subdirs[i] == "compiled" {
			m := i
			if len(mrkrm) > 0 {
				m = i - cmarks
				cmarks++
			}
			mrkrm = append(mrkrm, m)
		} else if fi, _ := os.Stat(utils.StringConcat(tpath, subdirs[i])); !fi.IsDir() {
			m := i
			if len(mrkrm) > 0 {
				m = i - cmarks
				cmarks++
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

	// Deal with magical number.
	instreams[0].Scan()

	// Read width, height and max value.
	w, h, max := -1, -1, -1
	instreams[0].Scan()
	fmt.Sscanf(instreams[0].Text(), "%d %d", &w, &h)
	instreams[0].Scan()
	fmt.Sscanf(instreams[0].Text(), "%d", &max)

	nin := len(instreams)
	// Move stream pointer to the right position.
	for i := 1; i < nin; i++ {
		instreams[i].Scan()
		instreams[i].Scan()
		instreams[i].Scan()
	}

	// Declare variables to data file.
	tt := w * h
	for i := 0; i < tt; i++ {
		fmt.Fprintf(out, "var %d %d\n", i, max+1)
	}
	fmt.Fprintf(out, "var %d %d\n", tt, nsdirs)

	for i := 0; i < nin; i++ {
		stream := instreams[i]
		for stream.Scan() {
			line := stream.Text()
			tokens := strings.Split(line, " ")
			ntokens := len(tokens)
			for j := 0; j < ntokens; j++ {
				tkn, err := strconv.Atoi(tokens[j])
				if err == nil {
					fmt.Fprintf(out, "%d ", tkn)
				}
			}
		}

		fmt.Fprintf(out, "%d\n", labels[i])
	}

	return w, h, max
}

// BufferedPGMFToData parses large quantities of files concurrently into a data file dname.
func BufferedPGMFToData(dirname, dname string) (int, int, int) {
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
	cmarks := 1
	// Reserved dirname compiled for output. Also remove non-dirs.
	for i := 0; i < nsdirs; i++ {
		// m marks the spot.
		// Since for every removed item the slice shrinks by one, we keep track of the indices by
		// taking into account the subslices "translated" at the right moment.
		if subdirs[i] == "compiled" {
			m := i
			if len(mrkrm) > 0 {
				m = i - cmarks
				cmarks++
			}
			mrkrm = append(mrkrm, m)
		} else if fi, _ := os.Stat(utils.StringConcat(tpath, subdirs[i])); !fi.IsDir() {
			m := i
			if len(mrkrm) > 0 {
				m = i - cmarks
				cmarks++
			}
			mrkrm = append(mrkrm, m)
		}
	}

	// Remove marked elements.
	for i := 0; i < len(mrkrm); i++ {
		j := mrkrm[i]
		subdirs, nsdirs = append(subdirs[:j], subdirs[j+1:]...), nsdirs-1
	}
	mrkrm = nil

	// Memorize filenames.
	var files []string
	var labels []int
	topdir := dirname + "/"
	for i := 0; i < nsdirs; i++ {
		sdir, err := os.Open(topdir + subdirs[i])
		if err != nil {
			fmt.Printf("Could not open subdirectory %s.\n", subdirs[i])
			sdir.Close()
			panic(err)
		}
		names, err := sdir.Readdirnames(-1)
		if err != nil {
			fmt.Printf("Coult not read files from subdirectory %s.\n", subdirs[i])
			sdir.Close()
			panic(err)
		}
		cmn := topdir + subdirs[i] + "/"
		for _, name := range names {
			files = append(files, cmn+name)
			labels = append(labels, i)
		}
		sdir.Close()
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

	// Parse the first file separately to find variables and image dimensions/max value.
	f, err := os.Open(files[0])
	if err != nil {
		fmt.Printf("Could not read file %s.\n", files[0])
		f.Close()
		panic(err)
	}
	stream := bufio.NewScanner(f)
	files = files[1:]

	// Deal with magical number.
	stream.Scan()
	w, h, max := -1, -1, -1
	stream.Scan()
	fmt.Sscanf(stream.Text(), "%d %d", &w, &h)
	stream.Scan()
	fmt.Sscanf(stream.Text(), "%d", &max)

	tt := w * h
	for i := 0; i < tt; i++ {
		fmt.Fprintf(out, "var %d %d\n", i, max+1)
	}
	fmt.Fprintf(out, "var %d %d\n", tt, nsdirs)

	for stream.Scan() {
		line := stream.Text()
		tokens := strings.Split(line, " ")
		ntokens := len(tokens)
		for i := 0; i < ntokens; i++ {
			tkn, err := strconv.Atoi(tokens[i])
			if err == nil {
				fmt.Fprintf(out, "%d ", tkn)
			}
		}
	}
	fmt.Fprintf(out, "%d\n", labels[0])
	labels = labels[1:]
	f.Close()

	var wg sync.WaitGroup
	nprocs, nrun := runtime.NumCPU(), 0
	cond, mutex := sync.NewCond(&sync.Mutex{}), &sync.Mutex{}
	nfiles := len(files)

	//fmt.Printf("nfiles: %d\n", nfiles)
	for i := 0; i < nfiles; i++ {
		cond.L.Lock()
		for nrun >= nprocs {
			cond.Wait()
		}
		nrun++
		cond.L.Unlock()
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mutex.Lock()
			//fmt.Printf("File: %s, id: %d\n", files[id], id)
			lf, err := os.Open(files[id])
			if err != nil {
				fmt.Printf("Could not read file %s.\n", files[id])
				panic(err)
			}
			mutex.Unlock()
			defer lf.Close()
			lstream := bufio.NewScanner(lf)
			// Magical number.
			lstream.Scan()
			// Width, height.
			lstream.Scan()
			// Max val.
			lstream.Scan()

			var vals []int
			for lstream.Scan() {
				line := lstream.Text()
				tokens := strings.Split(line, " ")
				ntokens := len(tokens)
				for j := 0; j < ntokens; j++ {
					tkn, err := strconv.Atoi(tokens[j])
					if err == nil {
						vals = append(vals, tkn)
					}
				}
			}
			vals = append(vals, labels[id])

			mutex.Lock()
			for j := 0; j < len(vals); j++ {
				fmt.Fprintf(out, "%d ", vals[j])
			}
			fmt.Fprintf(out, "%d\n", vals[len(vals)-1])
			mutex.Unlock()

			lf.Close()
			cond.L.Lock()
			nrun--
			cond.L.Unlock()
			cond.Signal()
		}(i)
	}

	wg.Wait()

	return w, h, max
}

// PGMFToEvidence (PGM file to evidence).
func PGMFToEvidence(dirname, dname string) (int, int, int) {
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
	cmarks := 1
	// Reserved dirname compiled for output. Also remove non-dirs.
	for i := 0; i < nsdirs; i++ {
		// m marks the spot.
		// Since for every removed item the slice shrinks by one, we keep track of the indices by
		// taking into account the subslices "translated" at the right moment.
		if fi, _ := os.Stat(utils.StringConcat(tpath, subdirs[i])); !fi.IsDir() ||
			subdirs[i] == "compiled" {
			m := i
			if len(mrkrm) > 0 {
				m = i - cmarks
				cmarks++
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

	// Deal with magical number.
	instreams[0].Scan()

	// Read width, height and max value.
	w, h, max := -1, -1, -1
	instreams[0].Scan()
	fmt.Sscanf(instreams[0].Text(), "%d %d", &w, &h)
	instreams[0].Scan()
	fmt.Sscanf(instreams[0].Text(), "%d", &max)

	nin := len(instreams)
	// Move stream pointer to the right position.
	for i := 1; i < nin; i++ {
		instreams[i].Scan()
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
		fmt.Fprintf(out, "var %d %d\n", i, max+1)
	}

	for i := 0; i < nin; i++ {
		stream := instreams[i]

		for stream.Scan() {
			line := stream.Text()
			tokens := strings.Split(line, " ")
			ntokens := len(tokens)
			for j := 0; j < ntokens; j++ {
				tkn, _ := strconv.Atoi(tokens[j])
				fmt.Fprintf(out, "%d ", tkn)
			}
		}

		fmt.Fprintf(out, "\n")
	}

	return w, h, max
}

// VarSetToPGM takes a state and draws according to the SPN that generated the instantiation.
func VarSetToPGM(filename string, state spn.VarSet, w, h, max int) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Could not create file [%s].\n", filename)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "P3\n%d %d\n%d\n", w, h, max)

	n := len(state)
	pixels := make([]int, n)
	for varid, val := range state {
		pixels[varid] = val
	}

	for i := 0; i < n; i++ {
		if i%71 == 0 {
			fmt.Fprintf(file, "\n")
		}
		p := pixels[i]
		common.DrawColorRGB(file, p, p, p)
		fmt.Fprintf(file, " ")
	}
}

// ImgCmplToPGM creates a new file distinguishing the original part of the image from the
// completion done by the SPN and indicated by typ.
func ImgCmplToPGM(filename string, orig, cmpl spn.VarSet, typ CmplType, w, h, max int) {
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

	fmt.Fprintf(file, "P3\n%d %d\n%d\n", w, h, max)

	n, j := w*h, 0
	for i := 0; i < n; i++ {
		if mid(i) {
			common.DrawColor(file, common.Red)
			goto cleanup
		} else if v, eo := orig[j]; eo {
			common.DrawColorRGB(file, v, v, v)
		} else {
			u, _ := cmpl[j]
			common.DrawColorRGB(file, 0, u, 0)
		}
		j++
	cleanup:
		fmt.Fprintf(file, " ")
		if i != 0 && i%w == 0 {
			fmt.Fprintf(file, "\n")
		}
	}
}
