package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	io "github.com/RenatoGeh/gospn/io"
	learn "github.com/RenatoGeh/gospn/learn"
	spn "github.com/RenatoGeh/gospn/spn"
	utils "github.com/RenatoGeh/gospn/utils"
	//profile "github.com/pkg/profile"
)

const dataset = "olivetti_3bit"

const (
	width  int = 46
	height int = 56
	max    int = 8
)

func halfImg(s spn.SPN, set spn.VarSet, typ io.CmplType, w, h int) (spn.VarSet, spn.VarSet) {
	cmpl, half := make(spn.VarSet), make(spn.VarSet)
	var criteria func(int) bool

	switch typ {
	case io.Top:
		criteria = func(p int) bool {
			return p < w*(h/2)
		}
	case io.Bottom:
		criteria = func(p int) bool {
			return p >= w*(h/2)
		}
	case io.Left:
		criteria = func(p int) bool {
			return p%w < w/2
		}
	case io.Right:
		criteria = func(p int) bool {
			return p%w >= w/2
		}
	}

	for k, v := range set {
		if !criteria(k) {
			half[k] = v
		}
	}

	cmpl, _ = s.ArgMax(half)

	for k := range half {
		delete(cmpl, k)
	}

	return cmpl, half
}

func classify(filename string, p float64, rseed int64, kclusters int) (spn.SPN, int, int) {
	vars, train, test, lbls := io.ParsePartitionedData(filename, p, rseed)
	s := learn.Gens(vars, train, kclusters)

	lines, n := len(test), len(vars)
	nclass := vars[n-1].Categories

	//fmt.Println("Drawing the MPE state of each class instance:")
	//evclass := make(spn.VarSet)
	//for i := 0; i < nclass; i++ {
	//evclass[n-1] = i
	//mpe, _ := s.ArgMax(evclass)
	//filename := fmt.Sprintf("mpe_%d.pbm", i)
	//delete(mpe, n-1)
	//io.VarSetToPBM(filename, mpe, width, height)
	//fmt.Printf("Class %d drawn to %s.\n", i, filename)
	//}

	corrects := 0
	for i := 0; i < lines; i++ {
		imax, max, prs := -1, -1.0, make([]float64, nclass)
		pz := s.Value(test[i])
		fmt.Printf("Testing instance %d. Should be classified as %d.\n", i, lbls[i])
		for j := 0; j < nclass; j++ {
			test[i][n-1] = j
			px := s.Value(test[i])
			prs[j] = utils.AntiLog(px - pz)
			fmt.Printf("  Pr(X=%d|E) = antilog(%.10f) = %.10f\n", j, px-pz, prs[j])
			if prs[j] > max {
				max, imax = prs[j], j
			}
		}
		fmt.Printf("Instance %d should be classified as %d. SPN classified as %d.\n", i, lbls[i], imax)
		if imax == lbls[i] {
			corrects++
		} else {
			fmt.Printf("--------> INCORRECT! <--------\n")
		}
		delete(test[i], n-1)
	}

	fmt.Printf("\n========= Iteration Results ========\n")
	fmt.Printf("  Correct classifications: %d/%d\n", corrects, lines)
	fmt.Printf("  Percentage of correct hits: %.2f%%\n", 100.0*(float64(corrects)/float64(lines)))
	fmt.Println("======================================")

	reps := make([]map[int]int, nclass)
	for i := 0; i < lines; i++ {
		if reps[lbls[i]] == nil {
			reps[lbls[i]] = test[i]
		}
	}
	for i := 0; i < nclass; i++ {
		for _, v := range io.Orientations {
			fmt.Printf("Drawing %s completion for digit %d.\n", v, i)
			cmpl, half := halfImg(s, reps[i], v, width, height)
			io.ImgCmplToPPM(fmt.Sprintf("cmpl_%d-%s.ppm", i, v), half, cmpl, v, width, height)
		}
	}

	return s, corrects, lines
}

func randVarSet(s spn.SPN, sc map[int]learn.Variable, n int) spn.VarSet {
	nsc := len(sc)
	vs := make(spn.VarSet)

	for i := 0; i < n; i++ {
		r := rand.Intn(nsc)
		id := sc[r]
		v := int(rand.NormFloat64()*(float64(id.Categories)/6) + float64(id.Categories/2))
		if v >= id.Categories {
			v = id.Categories - 1
		} else if v < 0 {
			v = 0
		}
		vs[id.Varid] = v
	}

	mpe, _ := s.ArgMax(vs)
	vs = nil
	return mpe
}

func imageCompletion(filename string, kclusters int, concurrents int) {
	fmt.Printf("Parsing data from [%s]...\n", filename)
	sc, data, lbls := io.ParseDataNL(filename)
	ndata := len(data)

	// Concurrency control.
	var wg sync.WaitGroup
	var nprocs int
	if concurrents <= 0 {
		nprocs = runtime.NumCPU()
	} else {
		nprocs = concurrents
	}
	nrun := 0
	cond := sync.NewCond(&sync.Mutex{})
	cpmutex := &sync.Mutex{}

	for i := 0; i < ndata; i++ {
		cond.L.Lock()
		for nrun >= nprocs {
			cond.Wait()
		}
		nrun++
		cond.L.Unlock()
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var train []map[int]int
			var ldata []map[int]int
			lsc := make(map[int]learn.Variable)

			cpmutex.Lock()
			for k, v := range sc {
				lsc[k] = v
			}
			for j := 0; j < ndata; j++ {
				ldata = append(ldata, make(map[int]int))
				for k, v := range data[j] {
					ldata[j][k] = v
				}
			}
			cpmutex.Unlock()

			chosen := ldata[id]
			for j := 0; j < ndata; j++ {
				if id != j && lbls[j] != lbls[id] {
					train = append(train, ldata[j])
				}
			}

			fmt.Printf("P-%d: Training SPN with %d clusters against instance %d...\n", id, kclusters, id)
			s := learn.Gens(lsc, train, kclusters)

			for _, v := range io.Orientations {
				fmt.Printf("P-%d: Drawing %s image completion for instance %d.\n", id, v, id)
				cmpl, half := halfImg(s, chosen, v, width, height)
				io.ImgCmplToPGM(fmt.Sprintf("cmpl_%d-%s.pgm", id, v), half, cmpl, v, width, height, max-1)
				cmpl, half = nil, nil
			}
			fmt.Printf("P-%d: Drawing MPE image for instance %d.\n", id, id)
			io.VarSetToPGM(fmt.Sprintf("mpe_cmpl_%d.pgm", id), randVarSet(s, lsc, 100),
				width, height, max-1)

			//out, _ := filepath.Abs("results/" + dataset + "/models")
			//io.DrawGraphTools(utils.StringConcat(out, "/all.py"), s)

			// Force garbage collection.
			s = nil
			train = nil
			lsc = nil
			ldata = nil

			cond.L.Lock()
			nrun--
			cond.L.Unlock()
			cond.Signal()
		}(i)
	}
	wg.Wait()
}

func convertData() {
	cmn, _ := filepath.Abs("data/" + dataset + "/")
	io.PGMFToData(cmn, "all.data")
}

func main() {
	p := 0.7
	kclusters := -1
	var rseed int64 = -1
	iterations := 1
	concurrents := -1
	var err error

	//defer profile.Start().Stop()

	if len(os.Args) > 5 {
		concurrents, err = strconv.Atoi(os.Args[5])
		if err != nil {
			fmt.Printf("Argument invalid. Argument concurrents must be an integer.\n")
			return
		}
	}
	if len(os.Args) > 4 {
		iterations, err = strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Printf("Argument invalid. Argument iterations must be an integer greater than zero.\n")
			return
		}
	}
	if len(os.Args) > 3 {
		kclusters, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf("Argument invalid. Argument kcluster must be an integer.\n")
			return
		}
	}
	if len(os.Args) > 2 {
		rseed, err = strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Printf("Argument invalid. Argument rseed must be a 64-bit integer.\n")
			return
		}
	}
	if len(os.Args) > 1 {
		p, err = strconv.ParseFloat(os.Args[1], 64)
		if err != nil || p < 0 || p >= 1 {
			if p == -1 {
				fmt.Printf("Converting dataset %s...", dataset)
				convertData()
				return
			}
			fmt.Printf("Argument invalid. Argument p must be a 64-bit float in the interval (0, 1).")
			return
		}
	}

	in, _ := filepath.Abs("data/" + dataset + "/compiled")
	out, _ := filepath.Abs("results/" + dataset + "/models")

	if p == 0 {
		fmt.Printf("Running image completion on dataset %s...\n", dataset)
		imageCompletion(utils.StringConcat(in, "/all.data"), kclusters, concurrents)
		return
	}

	fmt.Printf("Running cross-validation test with p = %.2f%%, random seed = %d and kclusters = %d "+
		"on the dataset = %s.\n", 100.0*p, rseed, kclusters, dataset)
	fmt.Printf("Iterations to run: %d\n\n", iterations)

	corrects, total := 0, 0
	for i := 0; i < iterations; i++ {
		fmt.Printf("+-----------------------------------------------+\n")
		fmt.Printf("|================ Iteration %d ==================|\n", i+1)
		fmt.Printf("+-----------------------------------------------+\n")
		s, c, t := classify(utils.StringConcat(in, "/all.data"), p, rseed, kclusters)
		corrects, total = corrects+c, total+t
		fmt.Printf("+-----------------------------------------------+\n")
		fmt.Printf("|============= End of Iteration %d =============|\n", i+1)
		fmt.Printf("+-----------------------------------------------+\n")
		io.DrawGraphTools(utils.StringConcat(out, "/all.py"), s)
	}
	fmt.Printf("---------------------------------\n")
	fmt.Printf(">>>>>>>>> Final Results <<<<<<<<<\n")
	fmt.Printf("  Correct classifications: %d/%d\n", corrects, total)
	fmt.Printf("  Percentage of correct hits: %.2f%%\n", 100.0*(float64(corrects)/float64(total)))
	fmt.Printf("---------------------------------\n")
}
