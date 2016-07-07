package main

import (
	"fmt"
	"path/filepath"

	io "github.com/RenatoGeh/gospn/src/io"
	learn "github.com/RenatoGeh/gospn/src/learn"
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

func main() {
	//indep_test()
	learn_test()
	//convert_imgs()
	//parse_test()
}
