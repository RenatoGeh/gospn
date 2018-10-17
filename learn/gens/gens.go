package gens

import (
	"fmt"
	"sync"

	"github.com/RenatoGeh/gospn/conc"
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/utils"
	"github.com/RenatoGeh/gospn/utils/cluster"
	"github.com/RenatoGeh/gospn/utils/indep"
)

// Binded is a binded version of Gens.
func Binded(kclusters int, pval, eps float64, mp int) learn.LearnFunc {
	return func(sc map[int]*learn.Variable, data spn.Dataset) spn.SPN {
		return Learn(sc, data, kclusters, pval, eps, mp)
	}
}

func LearnConcurrent(sc map[int]*learn.Variable, data []map[int]int, kclusters int, pval, eps float64, mp int, procs int) spn.SPN {
	n := len(sc)
	if n == 1 {
		var tv *learn.Variable
		for _, v := range sc {
			tv = v
		}
		return newMultinom(tv, data)
	}
	vdata := learn.DataToVarData(data, sc)
	igraph := indep.NewUFIndepGraph(vdata, pval)
	vdata = nil
	if len(igraph.Kset) > 1 {
		return indepStep(procs, kclusters, 0, pval, eps, mp, data, sc, igraph)
	}
	igraph = nil
	return clusterStep(procs, kclusters, 0, pval, eps, mp, data, sc)
}

// Learn runs the Gens Learning Algorithm
// Based on the article
//	Learning the Structure of Sum Product Networks
//	Robert Gens and Pedro Domingos
//	International Conference on Machine Learning 30 (ICML 2013)
func Learn(sc map[int]*learn.Variable, data []map[int]int, kclusters int, pval, eps float64, mp int) spn.SPN {
	n := len(sc)
	// If the data's scope is unary, then we return a leaf (i.e. a univariate distribution).
	if n == 1 {
		var tv *learn.Variable
		for _, v := range sc {
			tv = v
		}
		return newMultinom(tv, data)
	}

	// Else we check for independent subsets of variables. We separate variables in k partitions,
	// where every partition is pairwise indepedent with each other.
	vdata := learn.DataToVarData(data, sc)
	// Independency graph.
	igraph := indep.NewUFIndepGraph(vdata, pval)
	vdata = nil
	// If true, then we can partition the set of variables in data into independent subsets. This
	// means we can create a product node (since product nodes' children have disjoint scopes).
	if len(igraph.Kset) > 1 {
		return indepStep(1, kclusters, 0, pval, eps, mp, data, sc, igraph)
	}
	igraph = nil
	// Else we perform k-clustering on the instances.
	return clusterStep(1, kclusters, 0, pval, eps, mp, data, sc)
}

func newMultinom(v *learn.Variable, data []map[int]int) spn.SPN {
	m := len(data)
	counts := make([]int, v.Categories)
	for i := 0; i < m; i++ {
		counts[data[i][v.Varid]]++
	}
	return spn.NewCountingMultinomial(v.Varid, counts)
}

func newGaussMix(varid, g int, data []map[int]int) spn.SPN {
	X := learn.ExtractInstance(varid, data)
	Q := utils.PartitionQuantiles(X, g)
	s := spn.NewSum()
	for _, q := range Q {
		s.AddChildW(spn.NewGaussianParams(varid, q[0], q[1]), 1.0/float64(len(Q)))
	}
	return s
}

func newFullyFactorized(g int, D []map[int]int, Sc map[int]*learn.Variable) spn.SPN {
	prod := spn.NewProduct()
	if g <= 0 {
		m := len(D)
		for _, v := range Sc {
			counts := make([]int, v.Categories)
			for i := 0; i < m; i++ {
				counts[D[i][v.Varid]]++
			}
			leaf := spn.NewCountingMultinomial(v.Varid, counts)
			prod.AddChild(leaf)
		}
	} else {
		for _, v := range Sc {
			X := learn.ExtractInstance(v.Varid, D)
			Q := utils.PartitionQuantiles(X, g)
			z := spn.NewSum()
			for _, q := range Q {
				z.AddChildW(spn.NewGaussianParams(v.Varid, q[0], q[1]), 1.0/float64(len(Q)))
			}
			prod.AddChild(z)
		}
	}
	return prod
}

func indepStep(np, kc, g int, pval, eps float64, mp int, D []map[int]int, Sc map[int]*learn.Variable, igraph *indep.Graph) spn.SPN {
	Q := conc.NewSingleQueue(np)
	mu := &sync.Mutex{}
	prod, m, kset := spn.NewProduct(), len(igraph.Kset), &igraph.Kset
	tn := len(D)
	fmt.Println("Fork start | indep")
	step := func(id int) {
		tdata := make([]map[int]int, tn)
		s := len((*kset)[id])
		for j := 0; j < tn; j++ {
			tdata[j] = make(map[int]int)
			for l := 0; l < s; l++ {
				k := (*kset)[id][l]
				tdata[j][k] = D[j][k]
			}
		}
		nsc := make(map[int]*learn.Variable)
		for j := 0; j < s; j++ {
			t := (*kset)[id][j]
			nsc[t] = &learn.Variable{Varid: t, Categories: Sc[t].Categories, Name: ""}
		}
		var nc spn.SPN
		if g > 0 {
			nc = LearnGauss(nsc, tdata, kc, pval, eps, mp, g)
		} else {
			nc = Learn(nsc, tdata, kc, pval, eps, mp)
		}
		mu.Lock()
		prod.AddChild(nc)
		mu.Unlock()
	}
	for i := 0; i < m; i++ {
		if np != 1 {
			Q.Run(step, i)
		} else {
			step(i)
		}
	}
	if np != 1 {
		Q.Wait()
	}
	fmt.Println("Fork end | indep")
	return prod
}

func clusterStep(np, k, g int, pval, eps float64, mp int, D []map[int]int, Sc map[int]*learn.Variable) spn.SPN {
	Q := conc.NewSingleQueue(np)
	mu := &sync.Mutex{}
	var clusters [][]map[int]int
	if k > 0 {
		if len(D) < k {
			return newFullyFactorized(g, D, Sc)
		}
		clusters = cluster.KMeansDataI(k, D)
	} else {
		clusters = cluster.DBSCANData(D, eps, mp)
	}
	if c := len(clusters); c == 1 {
		return newFullyFactorized(g, D, Sc)
	}
	sum := spn.NewSum()
	fmt.Println("Fork start | clusters")
	step := func(id int) {
		nsc := learn.ReflectScope(Sc)
		var nc spn.SPN
		if g > 0 {
			nc = LearnGauss(nsc, clusters[id], k, pval, eps, mp, g)
		} else {
			nc = Learn(nsc, clusters[id], k, pval, eps, mp)
		}
		mu.Lock()
		sum.AddChildW(nc, float64(len(clusters[id]))/float64(len(D)))
		mu.Unlock()
	}
	for i := range clusters {
		if np != 1 {
			Q.Run(step, i)
		} else {
			step(i)
		}
	}
	if np != 1 {
		Q.Wait()
	}
	fmt.Println("Fork start | clusters")
	return sum
}

// LearnGaussConcurrent learns with gaussians concurrently.
func LearnGaussConcurrent(sc map[int]*learn.Variable, data []map[int]int, kclusters int, pval, eps float64, mp, g, procs int) spn.SPN {
	n := len(sc)
	if n == 1 {
		var v *learn.Variable
		for _, u := range sc {
			v = u
		}
		return newGaussMix(v.Varid, g, data)
	}
	vdata := learn.DataToVarData(data, sc)
	igraph := indep.NewUFIndepGraph(vdata, pval)
	vdata = nil
	if len(igraph.Kset) > 1 {
		return indepStep(procs, kclusters, g, pval, eps, mp, data, sc, igraph)
	}
	igraph = nil
	return clusterStep(procs, kclusters, g, pval, eps, mp, data, sc)
}

// LearnGauss uses Gaussians instead of Multinomials.
func LearnGauss(sc map[int]*learn.Variable, data []map[int]int, kclusters int, pval, eps float64, mp, g int) spn.SPN {
	n := len(sc)
	if n == 1 {
		var v *learn.Variable
		for _, u := range sc {
			v = u
		}
		return newGaussMix(v.Varid, g, data)
	}
	vdata := learn.DataToVarData(data, sc)
	igraph := indep.NewUFIndepGraph(vdata, pval)
	vdata = nil
	if len(igraph.Kset) > 1 {
		return indepStep(1, kclusters, g, pval, eps, mp, data, sc, igraph)
	}
	igraph = nil
	return clusterStep(1, kclusters, g, pval, eps, mp, data, sc)
}
