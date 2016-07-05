package main

import (
	"fmt"
	"path/filepath"

	io "github.com/RenatoGeh/gospn/src/io"
	learn "github.com/RenatoGeh/gospn/src/learn"
	utils "github.com/RenatoGeh/gospn/src/utils"
)

func learn_test() {
	comps, err := filepath.Abs("../data/crt/compiled")

	if err != nil {
		fmt.Printf("Error on finding data/crt/compiled.\n")
		panic(err)
	}

	res, err := filepath.Abs("../results/crt/models/circles")

	if err != nil {
		fmt.Printf("Error on finding results/crt/models.\n")
		panic(err)
	}

	fmt.Printf("Input path:\n%s\nOutput path:\n%s\nLearning...\n", comps, res)
	s := learn.Gens(io.ParseData(utils.StringConcat(comps, "/circles.data")))
	fmt.Printf("Drawing graph...\n")
	io.DrawGraph(utils.StringConcat(res, "/circles.dot"), s)
}

func indep_test() {
	fmt.Printf("Chi-square: %f\n", utils.Chisquare(2, 15.2))

	data := [][]int{
		{200, 150, 50, 400},
		{250, 300, 50, 600},
		{450, 450, 100, 1000}}
	fmt.Printf("Indep? %t\n", utils.ChiSquareTest(2, 3, data))
}

func convert_imgs() {
	cmn, _ := filepath.Abs("../data/crtsf/")
	io.PBMFToData(cmn, "all.data")
}

func main() {
	//indep_test()
	//learn_test()
	convert_imgs()
}
