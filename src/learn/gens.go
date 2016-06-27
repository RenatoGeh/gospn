package learn

import (
	spn "github.com/RenatoGeh/gospn/src/spn"
	utils "github.com/RenatoGeh/gospn/src/utils"
)

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
// 		/src/utils/kmeans.go
//
// As for testing the independency between two variables we use the Chi-Square independence test,
// present in file:
//
// 		/src/utils/indtest.go
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
func Gens(sc map[int]Variable, data [][]int) spn.SPN {
	n := len(sc)

	// If the data's scope is unary, then we return a leaf (i.e. a univariate distribution).
	if n == 1 {
		// m number of instantiations.
		m := len(data)
		// pr is the univariate probability distribution.
		var tv *Variable
		for _, v := range sc {
			tv = &v
		}
		pr := make([]float64, tv.Categories)
		for i := 0; i < m; i++ {
			pr[data[i][0]]++
		}
		for i := 0; i < m; i++ {
			pr[i] /= float64(m)
		}

		leaf := spn.NewUnivDist(tv.Varid, pr)
		return leaf
	}

	// Else we check for independent subsets of variables. We separate variables in k partitions,
	// where every partition is pairwise indepedent with another.

	// vdata is the transpose of data.
	vdata, l := make([]*utils.VarData, n), 0
	for _, v := range sc {
		tn := len(data)
		// tdata is the transpose of data[l].
		tdata := make([]int, tn)
		for j := 0; j < tn; j++ {
			tdata[j] = data[j][l]
		}
		vdata[l] = utils.NewVarData(v.Varid, v.Categories, tdata)
		l++
	}

	// Independency graph.
	igraph := utils.NewIndepGraph(vdata)

	// If true, then we can partition the set of variables in data into independent subsets. This
	// means we can create a product node (since product nodes' children have disjoint scopes).
	if len(igraph.Kset) > 1 {
		// prod is the new product node. m is the number of disjoint sets. kset is a shortcut.
		prod, m, kset := spn.NewProduct(), len(igraph.Kset), &igraph.Kset
		tn := len(data)
		for i := 0; i < m; i++ {
			// Data slices of the relevant vectors.
			tdata := make([][]int, tn)
			// Number of variables in set of variables kset[i].
			s := len((*kset)[i])
			for j := 0; j < tn; j++ {
				tdata[i] = make([]int, s)
				for k := 0; k < s; k++ {
					// Get the instanciations of variables in kset[i].
					tdata[j][k] = data[j][(*kset)[i][k]]
				}
			}
			// Create new scope with new variables.
			nsc := make(map[int]Variable)
			for j := 0; j < s; j++ {
				t := (*kset)[i][j]
				nsc[t] = Variable{t, sc[t].Categories}
			}
			var ndata [][]int
			copy(ndata, tdata)
			// Adds the recursive calls as children of this new product node.
			prod.AddChild(Gens(nsc, ndata))
		}
		return prod
	}

	return nil
}
