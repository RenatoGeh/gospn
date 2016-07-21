package main

import (
	"fmt"
	"math"
	"math/rand"
	"path/filepath"
	"sort"

	common "github.com/RenatoGeh/gospn/src/common"
	io "github.com/RenatoGeh/gospn/src/io"
	learn "github.com/RenatoGeh/gospn/src/learn"
	spn "github.com/RenatoGeh/gospn/src/spn"
	utils "github.com/RenatoGeh/gospn/src/utils"
)

const dataset = "digits"

func learn_test() spn.SPN {
	comps, err := filepath.Abs("../data/" + dataset + "/compiled")

	if err != nil {
		fmt.Printf("Error on finding data/" + dataset + "/compiled.\n")
		panic(err)
	}

	res, err := filepath.Abs("../results/" + dataset + "/models/all")

	if err != nil {
		fmt.Printf("Error on finding results/" + dataset + "/models.\n")
		panic(err)
	}

	fmt.Printf("Input path:\n%s\nOutput path:\n%s\nLearning...\n", comps, res)
	s := learn.Gens(io.ParseData(utils.StringConcat(comps, "/all.data")))
	fmt.Printf("Drawing graph...\n")
	io.DrawGraphTools(utils.StringConcat(res, "/all.py"), s)

	return s
}

func indep_test() {
	fmt.Printf("Chi-square: %f\n", 1-utils.Chisquare(1, 6.73))

	data := [][]int{
		{200, 150, 50, 400},
		{250, 300, 50, 600},
		{450, 450, 100, 1000}}
	fmt.Printf("Indep? %t\n", utils.ChiSquareTest(2, 3, data, 1))
}

func parse_test() {
	sc, data := io.ParseData(io.GetPath("../data/digits/compiled/all.data"))

	for k, v := range sc {
		fmt.Printf("[k=%d] varid: %d, categories: %d\n", k, v.Varid, v.Categories)
	}

	n, m := len(data), len(data[0])
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			fmt.Printf("%d ", data[i][j])
		}
		fmt.Printf("\n")
	}
}

func convert_imgs() {
	cmn, _ := filepath.Abs("../data/" + dataset + "/")
	io.PBMFToData(cmn, "all.data")
}

func cvntev_imgs() {
	cmn, _ := filepath.Abs("../data/" + dataset + "_test/")
	io.PBMFToEvidence(cmn, "all.data")
}

func drawgraph_test() {
	l1, l2, l3, l4 := spn.NewEmptyUnivDist(0, 2), spn.NewEmptyUnivDist(1, 2), spn.NewEmptyUnivDist(2, 2), spn.NewEmptyUnivDist(3, 2)
	s1, s2 := spn.NewSum(), spn.NewSum()
	s1.AddChildW(l1, 0.3)
	s1.AddChildW(l2, 0.7)
	s2.AddChildW(l3, 0.4)
	s2.AddChildW(l4, 0.6)
	p1, p2 := spn.NewProduct(), spn.NewProduct()
	l5, l6 := spn.NewEmptyUnivDist(4, 2), spn.NewEmptyUnivDist(5, 2)
	p1.AddChild(s1)
	p1.AddChild(l5)
	p2.AddChild(s2)
	p2.AddChild(l6)
	s := spn.NewSum()
	s.AddChildW(p1, 0.2)
	s.AddChildW(p2, 0.8)

	path, _ := filepath.Abs("../results/example/simplespn")
	io.DrawGraphTools(utils.StringConcat(path, "/spn.py"), s)

	fmt.Println("Testing probabilities...")

	vset := make(spn.VarSet)
	vset[2], vset[1], vset[4] = 1, 0, 1
	val := s.Value(vset)
	fmt.Printf("Pr(X_1=0, X_2=1, X_4=1)=antiln(%f)=%f.\n", val, utils.AntiLog(val))
	delete(vset, 2)
	delete(vset, 1)
	delete(vset, 4)
	vset[4], vset[3], vset[1], vset[0], vset[2] = 0, 0, 1, 1, 0
	val = s.Value(vset)
	fmt.Printf("Pr(X_1=1, X_2=1, X_3=0, X_4=0, X_5=0)=antiln(%f)=%f.\n", val, utils.AntiLog(val))
	for i := 0; i < 5; i++ {
		delete(vset, i)
	}
	for i := 0; i < 6; i++ {
		vset[i] = 1
	}
	val = s.Value(vset)
	fmt.Printf("Pr(for all 0<=i<6, X_i=1)=antiln(%f)=%f.\n", val, utils.AntiLog(val))
}

func queue_test() {
	queue := common.QueueBFSPair{}
	queue.Enqueue(&common.BFSPair{nil, "1", 1})
	queue.Enqueue(&common.BFSPair{nil, "2", 2})
	queue.Enqueue(&common.BFSPair{nil, "3", 3})

	for !queue.Empty() {
		e := queue.Dequeue()
		fmt.Printf("\"%s\" - %f\n", e.Pname, e.Weight)
		fmt.Printf("Size: %d\n", queue.Size())
	}
	fmt.Printf("Size: %d\n", queue.Size())

	queue.Enqueue(&common.BFSPair{nil, "4", 4})
	fmt.Printf("Size: %d\n", queue.Size())
	queue.Enqueue(&common.BFSPair{nil, "5", 5})
	fmt.Printf("Size: %d\n", queue.Size())
	t := queue.Dequeue()
	fmt.Printf("\"%s\" - %f\n", t.Pname, t.Weight)
	fmt.Printf("Size: %d\n", queue.Size())
	queue.Enqueue(&common.BFSPair{nil, "6", 6})
	t = queue.Dequeue()
	fmt.Printf("\"%s\" - %f\n", t.Pname, t.Weight)
	fmt.Printf("Size: %d\n", queue.Size())

	for !queue.Empty() {
		e := queue.Dequeue()
		fmt.Printf("\"%s\" - %f\n", e.Pname, e.Weight)
		fmt.Printf("Size: %d\n", queue.Size())
	}
}

func classify_test() {
	s := learn_test()
	sc, ev, tlabels := io.ParseEvidence(io.GetPath("../data/" + dataset + "_test/compiled/all.data"))

	nsc := len(sc)
	nv := len(tlabels)

	totals := make([]int, nv)
	corrects := make([]int, nv)

	tlabels = append(tlabels, int(^uint(0)>>1))
	c, l := 0, 0
	for _, ve := range ev {
		fmt.Printf("Test %d...\n", c)
		fmt.Printf("X is supposed to be %d.\n", l)
		prs := make([]float64, nv)
		pz := s.Value(ve)
		for i := 0; i < nv; i++ {
			vset := make(spn.VarSet)
			for k, v := range ve {
				vset[k] = v
			}
			vset[nsc] = i
			px := s.Value(vset)
			pr := px - pz
			prs[i] = utils.AntiLog(pr)
			fmt.Printf("Pr(X=%d|E)=%f/%f=%.50f\n", i, px, pz, prs[i])
		}

		max, imax := 0.0, 0
		for i := 0; i < nv; i++ {
			if max < prs[i] {
				max, imax = prs[i], i
			}
		}
		fmt.Printf("Classified as class %d when it's supposed to be %d.\n", imax, l)

		totals[l]++
		if l == imax {
			corrects[l]++
		}

		c++
		if tlabels[l+1] == c {
			l++
		}
	}

	fmt.Printf("\n=========== Overall Results ============\n")
	for i := 0; i < nv; i++ {
		perc := 100.0 * (float64(corrects[i]) / float64(totals[i]))
		fmt.Printf("Class %d:\n  Total instances: %d\n  Correct instances: %d\n  Correctness "+
			"percentage: %.3f%%\n\n", i, totals[i], corrects[i], perc)
	}
	fmt.Printf("========================================")

	//argmax, max := s.ArgMax(ev[0])
	//arg, ok := argmax[600]
	//fmt.Printf("argmax_X Pr(X|E) = [%t, %d] %f\n", ok, arg, utils.AntiLog(max))
}

func log_test() {
	const n = 50
	pr, w := make([]float64, n), make([]float64, n)
	for i := 0; i < n; i++ {
		pr[i] = rand.Float64()
		w[i] = rand.Float64()
	}
	sumv, sum, prod := make([]float64, n), 0.0, 1.0
	for i := 0; i < n; i++ {
		sumv[i] = w[i] * pr[i]
		sum += w[i] * pr[i]
		prod *= pr[i]
	}
	ls, s := utils.AntiLog(utils.LogSum(sumv...)), sum
	lp, p := utils.AntiLog(utils.LogProd(pr...)), prod
	fmt.Printf("SUM:  (Lval=%.50f) == (Rval=%.50f) ? %t\n  DIFF: %.50f\n",
		ls, s, ls == s, math.Abs(ls-s))
	fmt.Printf("PROD: (Lval=%.50f) == (Rval=%.50f) ? %t\n  DIFF: %.50f\n",
		lp, p, lp == p, math.Abs(lp-p))
}

func discgraph_test() {
	graph := make(map[int][]int)

	// 20 nodes in this graph.
	const N = 20

	// Supposed to be 6 disconnected subgraphs in this graph.
	// Subgraph 1
	graph[0], graph[1], graph[2] = []int{1}, []int{0, 2}, []int{1}
	// Subgraph 2
	graph[3], graph[4], graph[5], graph[6] = []int{4}, []int{3, 4, 5}, []int{4, 6}, []int{4, 5}
	// Subgraph 3
	graph[7] = []int{}
	// Subgraph 4
	graph[8], graph[9], graph[10], graph[11] = []int{9, 14}, []int{8, 10}, []int{9, 11}, []int{10, 12}
	graph[12], graph[13], graph[14] = []int{11, 13, 14}, []int{12, 14}, []int{12, 13}
	// Subgraph 5
	graph[15] = []int{}
	// Subgraph 6
	graph[16], graph[17], graph[18] = []int{17, 18, 19}, []int{16, 18, 19}, []int{16, 17, 19}
	graph[19] = []int{16, 17, 18}

	sets := make([]*utils.UFNode, N)

	for i := 0; i < N; i++ {
		sets[i] = utils.MakeSet(i)
	}

	for i := 0; i < N; i++ {
		m := len(graph[i])
		for j := 0; j < m; j++ {
			t := graph[i][j]
			if utils.Find(sets[i]) == utils.Find(sets[t]) {
				continue
			}
			utils.Union(sets[i], sets[t])
		}
	}

	var subgraphs [][]int = nil
	// Find roots.
	for i := 0; i < N; i++ {
		if sets[i] == sets[i].Pa {
			subgraphs = append(subgraphs, utils.UFVarids(sets[i]))
		}
	}

	k := len(subgraphs)
	fmt.Printf("There are %d disconnected subgraphs.\n", k)
	for i := 0; i < k; i++ {
		fmt.Printf("Subgraph %d has %d elements:\n%v\n", i+1, len(subgraphs[i]), subgraphs[i])
	}

	M := 30
	sc := make(map[int]learn.Variable)
	data := make([]map[int]int, M)

	fmt.Printf("\nVariables:\n")
	for i := 0; i < N; i++ {
		sc[i] = learn.Variable{i, N * M}
		fmt.Printf("Var %d %d\n", i, N*M)
	}

	fmt.Printf("\nData:\n")
	for i := 0; i < M; i++ {
		data[i] = make(map[int]int)
		for j := 0; j < N; j++ {
			data[i][j] = j + i*N
			fmt.Printf("%3d ", data[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	kset := &subgraphs
	for i := 0; i < len(subgraphs); i++ {
		tn := len(data)
		tdata := make([]map[int]int, tn)
		s := len((*kset)[i])
		for j := 0; j < tn; j++ {
			tdata[j] = make(map[int]int)
			for l := 0; l < s; l++ {
				k := (*kset)[i][l]
				tdata[j][k] = data[j][k]
				fmt.Printf("V:%d=%3d ", k, tdata[j][k])
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
		nsc := make(map[int]learn.Variable)
		for j := 0; j < s; j++ {
			t := (*kset)[i][j]
			nsc[t] = learn.Variable{t, sc[t].Categories}
			fmt.Printf("Variable: %d, %d\n", t, sc[t].Categories)
		}
		fmt.Printf("\n")
	}
}

func kmeans_test() {
	data := [][]int{{0, 1, 2}, {2, 3, 4}, {4, 5, 6}, {6, 7, 8}, {8, 9, 10}, {1, 2, 3}, {3, 4, 5},
		{5, 6, 7}, {7, 8, 9}, {7, 5, 3}, {0, 0, 1}, {9, 9, 10}, {0, 5, 10}}
	k := 3
	clusters := utils.KMeansV(k, data)

	for i := 0; i < k; i++ {
		fmt.Printf("Cluster %d:\n", i)
		for k, v := range clusters[i] {
			fmt.Printf("[%d]=%d ", k, v)
		}
		fmt.Printf("\n")
	}

	mdata := make([]map[int]int, len(data))
	fmt.Printf("mdata:\n")
	for i := 0; i < len(data); i++ {
		mdata[i] = make(map[int]int)
		for j := 0; j < len(data[i]); j++ {
			mdata[i][j] = data[i][j]
			fmt.Printf("%d ", mdata[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	for i := 0; i < k; i++ {
		ni := len(clusters[i])
		ndata := make([]map[int]int, ni)

		l := 0
		for k, _ := range clusters[i] {
			ndata[l] = make(map[int]int)
			fmt.Printf("%d:\n", k)
			for index, value := range mdata[k] {
				ndata[l][index] = value
				fmt.Printf("[%d]=%d ", index, value)
			}
			fmt.Printf("\n")
			l++
		}

		fmt.Printf("Clusters %d:\n", i)
		for j := 0; j < ni; j++ {
			keys := make([]int, len(ndata[j]))
			t := 0
			for _, v := range ndata[j] {
				keys[t] = v
				t++
			}
			sort.Ints(keys)
			for tt := 0; tt < len(ndata[j]); tt++ {
				fmt.Printf("%d ", keys[tt])
			}

			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}
}

func vardata_test() {
	n, m := 40, 30

	sc := make(map[int]learn.Variable)
	data := make([]map[int]int, m)

	for i := 0; i < n; i++ {
		sc[i] = learn.Variable{i, 11}
	}

	for i := 0; i < m; i++ {
		data[i] = make(map[int]int)
		for j := 0; j < n; j++ {
			data[i][j] = (j*i)%3 + (j+i)%4*(2+j%2)
		}
	}

	fmt.Printf("Data:\n")
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf("%3d ", data[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	vdata, l := make([]*utils.VarData, n), 0
	indices := make([]int, n)
	for _, v := range sc {
		tn := len(data)
		// tdata is the transpose of data[k].
		tdata := make([]int, tn)
		for j := 0; j < tn; j++ {
			tdata[j] = data[j][v.Varid]
		}
		vdata[l] = utils.NewVarData(v.Varid, v.Categories, tdata)
		indices[v.Varid] = l
		l++
	}

	for i := 0; i < n; i++ {
		j := indices[i]
		fmt.Printf("vdata[%d]:\n  %d\n  %d\n  %v\n\n", i, vdata[j].Varid, vdata[j].Categories, vdata[j].Data)
	}

	igraph := utils.NewIndepGraph(vdata)

	kset := &igraph.Kset

	for i := 0; i < len(igraph.Kset); i++ {
		sort.Ints((*kset)[i])
	}

	for i := 0; i < len(igraph.Kset); i++ {
		tn := len(data)
		tdata := make([]map[int]int, tn)
		s := len((*kset)[i])
		for j := 0; j < tn; j++ {
			tdata[j] = make(map[int]int)
			for l := 0; l < s; l++ {
				k := (*kset)[i][l]
				tdata[j][k] = data[j][k]
				fmt.Printf("V:%d=%3d ", k, tdata[j][k])
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
		nsc := make(map[int]learn.Variable)
		for j := 0; j < s; j++ {
			t := (*kset)[i][j]
			nsc[t] = learn.Variable{t, sc[t].Categories}
			fmt.Printf("Variable: %d, %d\n", t, sc[t].Categories)
		}
		fmt.Printf("\n")
	}
}

func maptoslice_test() {
	N, M := 20, 10

	sc := make(map[int]learn.Variable)
	data := make([]map[int]int, M)

	for i := 0; i < N; i++ {
		sc[i] = learn.Variable{i, N * M}
	}

	for i := 0; i < M; i++ {
		data[i] = make(map[int]int)
		for j := 0; j < N; j++ {
			data[i][j] = j + i*N
		}
	}

	fmt.Printf("Data:\n")
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf("%3d ", data[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	m := len(data)
	mdata := make([][]int, m)
	for i := 0; i < m; i++ {
		lc := len(data[i])
		mdata[i] = make([]int, lc)
		l := 0
		keys := make([]int, lc)
		for k, _ := range data[i] {
			keys[l] = k
			l++
		}
		sort.Ints(keys)
		for j := 0; j < lc; j++ {
			mdata[i][j] = data[i][keys[j]]
		}
	}

	fmt.Printf("Mdata:\n")
	for i := 0; i < len(mdata); i++ {
		for j := 0; j < len(mdata[i]); j++ {
			fmt.Printf("%3d ", mdata[i][j])
		}
		fmt.Printf("\n")
	}
}

func main() {
	//indep_test()
	//learn_test()
	//convert_imgs()
	//cvntev_imgs()
	//parse_test()
	//drawgraph_test()
	//queue_test()
	classify_test()
	//log_test()
	//discgraph_test()
	//kmeans_test()
	//vardata_test()
	//maptoslice_test()
}
