package learn

import (
	"sort"

	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/utils"
	"github.com/RenatoGeh/gospn/utils/cluster"
	"github.com/RenatoGeh/gospn/utils/indep"
)

// BindedGens is a binded version of Gens.
func BindedGens(kclusters int, pval, eps float64, mp int) LearnFunc {
	return func(sc map[int]Variable, data spn.Dataset) spn.SPN {
		return Gens(sc, data, kclusters, pval, eps, mp)
	}
}

// Gens Learning Algorithm
// Based on the article
//	Learning the Structure of Sum Product Networks
//	Robert Gens and Pedro Domingos
//	International Conference on Machine Learning 30 (ICML 2013)
func Gens(sc map[int]Variable, data []map[int]int, kclusters int, pval, eps float64, mp int) spn.SPN {
	n := len(sc)

	sys.Printf("Sample size: %d, scope size: %d\n", len(data), n)

	// If the data's scope is unary, then we return a leaf (i.e. a univariate distribution).
	if n == 1 {
		sys.Println("Creating new leaf...")

		// m number of instantiations.
		m := len(data)
		// pr is the univariate probability distribution.
		var tv *Variable
		for _, v := range sc {
			tv = &v
		}
		counts := make([]int, tv.Categories)
		for i := 0; i < m; i++ {
			counts[data[i][tv.Varid]]++
		}

		leaf := spn.NewCountingMultinomial(tv.Varid, counts)
		//sys.Println("Leaf created.")
		return leaf
	}

	// Else we check for independent subsets of variables. We separate variables in k partitions,
	// where every partition is pairwise indepedent with each other.
	//sys.Println("Preparing to create new product node...")

	sys.Println("Creating VarDatas for Independency Test...")
	vdata, l := make([]*utils.VarData, n), 0
	for _, v := range sc {
		tn := len(data)
		// tdata is the transpose of data[k].
		tdata := make([]int, tn)
		for j := 0; j < tn; j++ {
			tdata[j] = data[j][v.Varid]
		}
		vdata[l] = utils.NewVarData(v.Varid, v.Categories, tdata)
		l++
	}

	sys.Println("Creating new Independency graph...")
	// Independency graph.
	igraph := indep.NewUFIndepGraph(vdata, pval)
	vdata = nil

	// If true, then we can partition the set of variables in data into independent subsets. This
	// means we can create a product node (since product nodes' children have disjoint scopes).
	if len(igraph.Kset) > 1 {
		sys.Println("Found independency. Separating independent sets.")

		//sys.Println("Found independency between variables. Creating new product node...")
		// prod is the new product node. m is the number of disjoint sets. kset is a shortcut.
		prod, m, kset := spn.NewProduct(), len(igraph.Kset), &igraph.Kset
		tn := len(data)
		for i := 0; i < m; i++ {
			//sort.Ints((*kset)[i])
			//}
			//nexti := 0
			//for done := 0; done < m; done++ {
			//i := 0
			//for (*kset)[i][0] != nexti {
			//i++
			//if i >= m {
			//i = 0
			//nexti++
			//}
			//}
			//nexti++
			// Data slices of the relevant vectors.
			tdata := make([]map[int]int, tn)
			// Number of variables in set of variables kset[i].
			s := len((*kset)[i])
			for j := 0; j < tn; j++ {
				tdata[j] = make(map[int]int)
				for l := 0; l < s; l++ {
					// Get the instanciations of variables in kset[i].
					//sys.Printf("[%d][%d] => %v vs %v | %v vs %v\n", j, k, (*kset)[i][k], len(data[j]), len(tdata[j]), k)
					k := (*kset)[i][l]
					tdata[j][k] = data[j][k]
				}
			}
			// Create new scope with new variables.
			nsc := make(map[int]Variable)
			for j := 0; j < s; j++ {
				t := (*kset)[i][j]
				nsc[t] = Variable{t, sc[t].Categories, ""}
			}
			//sys.Printf("LENGTH: %d\n", len(tdata))
			//sys.Println("Product node created. Recursing...")
			// Adds the recursive calls as children of this new product node.
			prod.AddChild(Gens(nsc, tdata, kclusters, pval, eps, mp))
		}
		return prod
	}
	igraph = nil

	// Else we perform k-clustering on the instances.
	sys.Println("No independency found. Preparing for clustering...")

	m := len(data)
	mdata := make([][]int, m)
	for i := 0; i < m; i++ {
		lc := len(data[i])
		mdata[i] = make([]int, lc)
		l := 0
		keys := make([]int, lc)
		for k := range data[i] {
			keys[l] = k
			l++
		}
		sort.Ints(keys)
		for j := 0; j < lc; j++ {
			mdata[i][j] = data[i][keys[j]]
		}
	}

	var clusters []map[int][]int
	if kclusters > 0 {
		sys.Printf("data: %d, mdata: %d\n", len(data), len(mdata))
		if len(mdata) < kclusters {
			//Fully factorized form.
			//All instances are approximately the same.
			prod := spn.NewProduct()
			m := len(data)
			for _, v := range sc {
				counts := make([]int, v.Categories)
				for i := 0; i < m; i++ {
					counts[data[i][v.Varid]]++
				}
				leaf := spn.NewCountingMultinomial(v.Varid, counts)
				prod.AddChild(leaf)
			}
			return prod
		}
		clusters = cluster.KMedoid(kclusters, mdata)
	} else if kclusters == -1 {
		clusters = cluster.DBSCAN(mdata, eps, mp)
	} else {
		clusters = cluster.OPTICS(mdata, eps, mp)
	}
	k := len(clusters)
	//sys.Printf("Clustering similar instances with %d clusters.\n", k)
	if k == 1 {
		// Fully factorized form.
		// All instances are approximately the same.
		prod := spn.NewProduct()
		m := len(data)
		for _, v := range sc {
			counts := make([]int, v.Categories)
			for i := 0; i < m; i++ {
				counts[data[i][v.Varid]]++
			}
			leaf := spn.NewCountingMultinomial(v.Varid, counts)
			counts = nil
			prod.AddChild(leaf)
		}
		return prod
	}
	mdata = nil

	sys.Println("Reformating clusters to appropriate format and creating sum node...")

	sum := spn.NewSum()
	for i := 0; i < k; i++ {
		ni := len(clusters[i])
		ndata := make([]map[int]int, ni)

		l := 0
		for k := range clusters[i] {
			ndata[l] = make(map[int]int)
			for index, value := range data[k] {
				ndata[l][index] = value
			}
			l++
		}

		nsc := make(map[int]Variable)
		for k, v := range sc {
			nsc[k] = v
		}

		//sys.Println("Created new sum node child. Recursing...")
		sum.AddChildW(Gens(nsc, ndata, kclusters, pval, eps, mp), float64(ni)/float64(len(data)))
	}

	clusters = nil
	return sum
}

// GensGauss uses Gaussians instead of Multinomials.
func GensGauss(sc map[int]Variable, data []map[int]int, kclusters int, pval, eps float64, mp int) spn.SPN {
	n := len(sc)

	sys.Printf("Sample size: %d, scope size: %d\n", len(data), n)

	// If the data's scope is unary, then we return a leaf (i.e. a univariate distribution).
	if n == 1 {
		sys.Println("Creating new leaf...")

		// m number of instantiations.
		m := len(data)
		// pr is the univariate probability distribution.
		var tv *Variable
		for _, v := range sc {
			tv = &v
		}
		counts := make([]int, tv.Categories)
		for i := 0; i < m; i++ {
			counts[data[i][tv.Varid]]++
		}

		leaf := spn.NewGaussian(tv.Varid, counts)
		//sys.Println("Leaf created.")
		return leaf
	}

	// Else we check for independent subsets of variables. We separate variables in k partitions,
	// where every partition is pairwise indepedent with each other.
	//sys.Println("Preparing to create new product node...")

	sys.Println("Creating VarDatas for Independency Test...")
	vdata, l := make([]*utils.VarData, n), 0
	for _, v := range sc {
		tn := len(data)
		// tdata is the transpose of data[k].
		tdata := make([]int, tn)
		for j := 0; j < tn; j++ {
			tdata[j] = data[j][v.Varid]
		}
		vdata[l] = utils.NewVarData(v.Varid, v.Categories, tdata)
		l++
	}

	sys.Println("Creating new Independency graph...")
	// Independency graph.
	igraph := indep.NewUFIndepGraph(vdata, pval)
	vdata = nil

	// If true, then we can partition the set of variables in data into independent subsets. This
	// means we can create a product node (since product nodes' children have disjoint scopes).
	if len(igraph.Kset) > 1 {
		sys.Println("Found independency. Separating independent sets.")

		//sys.Println("Found independency between variables. Creating new product node...")
		// prod is the new product node. m is the number of disjoint sets. kset is a shortcut.
		prod, m, kset := spn.NewProduct(), len(igraph.Kset), &igraph.Kset
		tn := len(data)
		for i := 0; i < m; i++ {
			//sort.Ints((*kset)[i])
			//}
			//nexti := 0
			//for done := 0; done < m; done++ {
			//i := 0
			//for (*kset)[i][0] != nexti {
			//i++
			//if i >= m {
			//i = 0
			//nexti++
			//}
			//}
			//nexti++
			// Data slices of the relevant vectors.
			tdata := make([]map[int]int, tn)
			// Number of variables in set of variables kset[i].
			s := len((*kset)[i])
			for j := 0; j < tn; j++ {
				tdata[j] = make(map[int]int)
				for l := 0; l < s; l++ {
					// Get the instanciations of variables in kset[i].
					//sys.Printf("[%d][%d] => %v vs %v | %v vs %v\n", j, k, (*kset)[i][k], len(data[j]), len(tdata[j]), k)
					k := (*kset)[i][l]
					tdata[j][k] = data[j][k]
				}
			}
			// Create new scope with new variables.
			nsc := make(map[int]Variable)
			for j := 0; j < s; j++ {
				t := (*kset)[i][j]
				nsc[t] = Variable{t, sc[t].Categories, ""}
			}
			//sys.Printf("LENGTH: %d\n", len(tdata))
			//sys.Println("Product node created. Recursing...")
			// Adds the recursive calls as children of this new product node.
			prod.AddChild(GensGauss(nsc, tdata, kclusters, pval, eps, mp))
		}
		return prod
	}
	igraph = nil

	// Else we perform k-clustering on the instances.
	sys.Println("No independency found. Preparing for clustering...")

	m := len(data)
	mdata := make([][]int, m)
	for i := 0; i < m; i++ {
		lc := len(data[i])
		mdata[i] = make([]int, lc)
		l := 0
		keys := make([]int, lc)
		for k := range data[i] {
			keys[l] = k
			l++
		}
		sort.Ints(keys)
		for j := 0; j < lc; j++ {
			mdata[i][j] = data[i][keys[j]]
		}
	}

	var clusters []map[int][]int
	if kclusters > 0 {
		sys.Printf("data: %d, mdata: %d\n", len(data), len(mdata))
		if len(mdata) < kclusters {
			//Fully factorized form.
			//All instances are approximately the same.
			prod := spn.NewProduct()
			m := len(data)
			for _, v := range sc {
				counts := make([]int, v.Categories)
				for i := 0; i < m; i++ {
					counts[data[i][v.Varid]]++
				}
				leaf := spn.NewGaussian(v.Varid, counts)
				prod.AddChild(leaf)
			}
			return prod
		}
		clusters = cluster.KMedoid(kclusters, mdata)
	} else if kclusters == -1 {
		clusters = cluster.DBSCAN(mdata, eps, mp)
	} else {
		clusters = cluster.OPTICS(mdata, eps, mp)
	}
	k := len(clusters)
	//sys.Printf("Clustering similar instances with %d clusters.\n", k)
	if k == 1 {
		// Fully factorized form.
		// All instances are approximately the same.
		prod := spn.NewProduct()
		m := len(data)
		for _, v := range sc {
			counts := make([]int, v.Categories)
			for i := 0; i < m; i++ {
				counts[data[i][v.Varid]]++
			}
			leaf := spn.NewGaussian(v.Varid, counts)
			counts = nil
			prod.AddChild(leaf)
		}
		return prod
	}
	mdata = nil

	sys.Println("Reformating clusters to appropriate format and creating sum node...")

	sum := spn.NewSum()
	for i := 0; i < k; i++ {
		ni := len(clusters[i])
		ndata := make([]map[int]int, ni)

		l := 0
		for k := range clusters[i] {
			ndata[l] = make(map[int]int)
			for index, value := range data[k] {
				ndata[l][index] = value
			}
			l++
		}

		nsc := make(map[int]Variable)
		for k, v := range sc {
			nsc[k] = v
		}

		//sys.Println("Created new sum node child. Recursing...")
		sum.AddChildW(GensGauss(nsc, ndata, kclusters, pval, eps, mp), float64(ni)/float64(len(data)))
	}

	clusters = nil
	return sum
}
