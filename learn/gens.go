package learn

import (
	"fmt"
	//"math"
	"sort"

	spn "github.com/RenatoGeh/gospn/spn"
	utils "github.com/RenatoGeh/gospn/utils"
	cluster "github.com/RenatoGeh/gospn/utils/cluster"
	indep "github.com/RenatoGeh/gospn/utils/indep"
)

// Gens Learning Algorithm
// We refer to this structural learning algorithm as the Gens Algorithm for structural learning.
// The full article describing this algorithm schema can be found at:
//
// 		http://spn.cs.washington.edu/pubs.shtml
//
// Under the name of
//
// 		Learning the Structure of Sum-Product Networks
// 		Robert Gens and Pedro Domingos; ICML 2013
//
// For clustering we use k-means clustering. Our implementation can be seen in file:
//
// 		/utils/cluster/kmeans.go
//
// As for testing the independency between two variables we use the Chi-Square independence test,
// present in file:
//
// 		/utils/cluster/indtest.go
//
// Function Gens takes as input a matrix of data instances, where the columns are variables and
// lines are the observed instantiations of each variable.
//
// 		+-----+------------------------------+
// 		|     | X_1   X_2   X_3   ...   X_n  |
// 		+-----+------------------------------+
// 		| I_1 | x_11  x_12  x_13  ...   x_1n |
// 		| I_2 | x_21  x_22  x_23  ...   x_2n |
// 		|  .  |  .     .     .     .     .   |
// 		|  .  |  .     .     .     .     .   |
// 		|  .  |  .     .     .     .     .   |
// 		| I_m | x_m1  x_m2  x_m3  ...   x_mn |
// 		+-----+------------------------------+
//
// Where X={X_1,...,X_n} is the set of variables and I={I_1,...,I_m} is the set of instances.
// Each x_ij is the i-th observed instantiation of X_j.
func Gens(sc map[int]Variable, data []map[int]int, kclusters int) spn.SPN {
	n := len(sc)

	fmt.Printf("Sample size: %d, scope size: %d\n", len(data), n)

	// If the data's scope is unary, then we return a leaf (i.e. a univariate distribution).
	if n == 1 {
		fmt.Println("Creating new leaf...")

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

		leaf := spn.NewCountingUnivDist(tv.Varid, counts)
		//fmt.Println("Leaf created.")
		return leaf
	}

	// Else we check for independent subsets of variables. We separate variables in k partitions,
	// where every partition is pairwise indepedent with each other.
	//fmt.Println("Preparing to create new product node...")

	fmt.Println("Creating VarDatas for Independency Test...")
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

	fmt.Println("Creating new Independency graph...")
	// Independency graph.
	igraph := indep.NewUFIndepGraph(vdata)
	vdata = nil

	// If true, then we can partition the set of variables in data into independent subsets. This
	// means we can create a product node (since product nodes' children have disjoint scopes).
	if len(igraph.Kset) > 1 {
		fmt.Println("Found independency. Separating independent sets.")

		//fmt.Println("Found independency between variables. Creating new product node...")
		// prod is the new product node. m is the number of disjoint sets. kset is a shortcut.
		prod, m, kset := spn.NewProduct(), len(igraph.Kset), &igraph.Kset
		tn := len(data)
		for i := 0; i < m; i++ {
			// Data slices of the relevant vectors.
			tdata := make([]map[int]int, tn)
			// Number of variables in set of variables kset[i].
			s := len((*kset)[i])
			for j := 0; j < tn; j++ {
				tdata[j] = make(map[int]int)
				for l := 0; l < s; l++ {
					// Get the instanciations of variables in kset[i].
					//fmt.Printf("[%d][%d] => %v vs %v | %v vs %v\n", j, k, (*kset)[i][k], len(data[j]), len(tdata[j]), k)
					k := (*kset)[i][l]
					tdata[j][k] = data[j][k]
				}
			}
			// Create new scope with new variables.
			nsc := make(map[int]Variable)
			for j := 0; j < s; j++ {
				t := (*kset)[i][j]
				nsc[t] = Variable{t, sc[t].Categories}
			}
			//fmt.Printf("LENGTH: %d\n", len(tdata))
			//fmt.Println("Product node created. Recursing...")
			// Adds the recursive calls as children of this new product node.
			prod.AddChild(Gens(nsc, tdata, kclusters))
		}
		return prod
	}
	igraph = nil

	// Else we perform k-clustering on the instances.
	fmt.Println("No independency found. Preparing for clustering...")

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
		fmt.Printf("data: %d, mdata: %d\n", len(data), len(mdata))
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
				leaf := spn.NewCountingUnivDist(v.Varid, counts)
				prod.AddChild(leaf)
			}
			return prod
		}
		clusters = cluster.KMeansV(kclusters, mdata)
	} else if kclusters == -1 {
		clusters = cluster.DBSCAN(mdata, 4, 4)
	} else {
		clusters = cluster.OPTICS(mdata, 10, 4)
	}
	k := len(clusters)
	//fmt.Printf("Clustering similar instances with %d clusters.\n", k)
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
			leaf := spn.NewCountingUnivDist(v.Varid, counts)
			prod.AddChild(leaf)
		}
		return prod
	}
	mdata = nil

	fmt.Println("Reformating clusters to appropriate format and creating sum node...")

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

		//fmt.Println("Created new sum node child. Recursing...")
		sum.AddChildW(Gens(nsc, ndata, kclusters), float64(ni)/float64(len(data)))
	}

	clusters = nil
	return sum
}
