package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/RenatoGeh/gospn/app"
	"github.com/RenatoGeh/gospn/io"
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/utils"
	//profile "github.com/pkg/profile"
)

var dataset = "olivetti_3bit"

func convertData() {
	cmn, _ := filepath.Abs("data/" + dataset + "/")
	io.BufferedPGMFToData(cmn, "all.data")
}

func main() {
	var p float64
	var clusters int
	var rseed int64
	var iterations int
	var concurrents int
	var mode string

	flag.Float64Var(&p, "p", 0.7, "Train/test partition ratio to be used for cross-validation. ")
	flag.IntVar(&clusters, "clusters", -1, "Number of clusters to be used during training. If "+
		"clusters = -1, GoSPN shall use DBSCAN. Else, if clusters = -2, then use OPTICS "+
		"(experimental). Else, if clusters > 0, then use k-means clustering with the indicated "+
		"number of clusters.")
	flag.Int64Var(&rseed, "rseed", -1, "Seed to be used when choosing which instances to be used as "+
		"training set and which to be used as testing set. If omitted, rseed defaults to -1, which "+
		"means GoSPN chooses a random seed according to the current time.")
	flag.IntVar(&iterations, "iterations", 1, "How many iterations to be run when running a "+
		"classification job. This allows for better, more general and randomized results, as some "+
		"test/train partitions may become degenerated.")
	flag.IntVar(&concurrents, "concurrents", -1, "GoSPN makes use of Go's natie concurrency and is "+
		"able to run on multiple cores in parallel. Argument concurrents defines the number of "+
		"concurrent jobs GoSPN should run at most. If concurrents <= 0, then concurrents = nCPU, "+
		"where nCPU is the number of CPUs the running machine has available.")
	flag.StringVar(&dataset, "dataset", dataset, "The name of the directory containing the "+
		"dataset structure inside the data folder. Setting -mode=data will cause a new given "+
		"dataset data file to be created. Omitting -mode or setting -mode to something different "+
		"than data will run a job on the given dataset.")
	flag.IntVar(&sys.Width, "width", sys.Width, "The width of the images to be classified or "+
		"completed.")
	flag.IntVar(&sys.Height, "height", sys.Height, "The height of the images to be classified or "+
		"completed.")
	flag.IntVar(&sys.Max, "max", sys.Max, "The maximum pixel value the images can have.")
	flag.StringVar(&mode, "mode", "cmpl", "Whether to convert a directory structure into a data "+
		"file (data), run an image completion job (cmpl) or a classification job (class).")
	flag.Float64Var(&sys.Pval, "pval", sys.Pval, "The significance value for the independence test.")
	flag.Float64Var(&sys.Eps, "eps", sys.Eps, "The epsilon minimum distance value for DBSCAN.")
	flag.IntVar(&sys.Mp, "mp", sys.Mp, "The minimum points density for DBSCAN.")
	flag.BoolVar(&sys.Verbose, "v", sys.Verbose, "Verbose mode.")

	flag.Parse()

	if p == 0 || p < 0 || p == 1 {
		fmt.Println("Argument p must be a float64 in range (0, 1).")
		return
	}
	if iterations <= 0 {
		fmt.Println("Argument iterations must be an integer greater than 0.")
		return
	}

	//defer profile.Start().Stop()

	in, _ := filepath.Abs("data/" + dataset + "/compiled")
	//out, _ := filepath.Abs("results/" + dataset + "/models")

	if mode == "data" {
		fmt.Printf("Converting dataset %s...\n", dataset)
		convertData()
		return
	} else if mode == "cmpl" {
		fmt.Printf("Running image completion on dataset %s with %d threads...\n", dataset, concurrents)
		lf := learn.BindedGens(clusters, sys.Pval, sys.Eps, sys.Mp)
		app.ImgCompletion(lf, utils.StringConcat(in, "/all.data"), concurrents)
		return
	} else if mode == "class" {
		lf := learn.BindedGens(clusters, sys.Pval, sys.Eps, sys.Mp)
		app.ImgBatchClassify(lf, dataset, p, rseed, clusters, iterations)
	} else if mode == "test" {
		//_, data, _ := io.ParseDataNL("data/digits/compiled/all.data")
		//_, data, _ := io.ParseDataNL("data/test/compiled/all.data")
		//sys.Width, sys.Height = 4, 4
		sys.Max = 256
		sys.Width, sys.Height = 48, 56
		//sys.Max = 2
		//sys.Width, sys.Height = 4, 4
		sys.Verbose = true
		//sys.MemConservative = true
		app.ImgTest("data/olivetti_padded/compiled/all.data", 2, 4, 4, 0.1, 1)
		//app.ImgTest("data/fourbyfour/compiled/all.data", 2, 1, 1, 0.1)
		//lf := learn.BindedPoonGD(2, 4, 0.1, 1)
		//sc, data, _ := io.ParseDataNL(filename)
		//S := lf(sc, data)
		//app.ImgClassify(lf, "data/digits/compiled/all.data", 0.3, -1)
		//app.ImgCompletion(lf, "data/olivetti_padded/compiled/all.data", 1)
		//learn.PoonTest(data, 2, 2)
	} else {
		fmt.Printf("Mode %s not found. Possible mode options:\n  cmpl, class, data\n", mode)
	}
}
