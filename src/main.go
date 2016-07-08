package main

import (
	"fmt"
	"path/filepath"

	io "github.com/RenatoGeh/gospn/src/io"
	learn "github.com/RenatoGeh/gospn/src/learn"
	spn "github.com/RenatoGeh/gospn/src/spn"
	utils "github.com/RenatoGeh/gospn/src/utils"
)

func learn_test() {
	comps, err := filepath.Abs("../data/crtsf/compiled")

	if err != nil {
		fmt.Printf("Error on finding data/crtsf/compiled.\n")
		panic(err)
	}

	res, err := filepath.Abs("../results/crtsf/models/all")

	if err != nil {
		fmt.Printf("Error on finding results/crt/models.\n")
		panic(err)
	}

	fmt.Printf("Input path:\n%s\nOutput path:\n%s\nLearning...\n", comps, res)
	s := learn.Gens(io.ParseData(utils.StringConcat(comps, "/all.data")))
	fmt.Printf("Drawing graph...\n")
	io.DrawGraph(utils.StringConcat(res, "/all.dot"), s)
}

func indep_test() {
	fmt.Printf("Chi-square: %f\n", 1-utils.Chisquare(1, 6.73))

	data := [][]int{
		{200, 150, 50, 400},
		{250, 300, 50, 600},
		{450, 450, 100, 1000}}
	fmt.Printf("Indep? %t\n", utils.ChiSquareTest(2, 3, data))
}

func parse_test() {
	sc, data := io.ParseData(io.GetPath("../data/crtsf/compiled/all.data"))

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
	cmn, _ := filepath.Abs("../data/crtsf/")
	io.PBMFToData(cmn, "all.data")
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

	path, _ := filepath.Abs("../results/crtsf/models/all")
	io.DrawGraph(utils.StringConcat(path, "/all.dot"), s)
}

func queue_test() {
	queue := utils.QueueBFSPair{}
	queue.Enqueue(&utils.BFSPair{nil, "1", 1})
	queue.Enqueue(&utils.BFSPair{nil, "2", 2})
	queue.Enqueue(&utils.BFSPair{nil, "3", 3})

	for !queue.Empty() {
		e := queue.Dequeue()
		fmt.Printf("\"%s\" - %f\n", e.Pname, e.Weight)
	}
	fmt.Printf("Size: %d\n", queue.Size())
}

func main() {
	//indep_test()
	//learn_test()
	//convert_imgs()
	//parse_test()
	drawgraph_test()
	//queue_test()
}
