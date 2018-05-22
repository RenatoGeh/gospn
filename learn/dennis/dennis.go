package dennis

import (
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/learn/parameters"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/utils/cluster"
	"github.com/RenatoGeh/gospn/utils/cluster/metrics"
	"sort"
)

func transpose(C map[int][]int) map[int][]int {
	var K []int
	for k, _ := range C {
		K = append(K, k)
	}
	sort.Ints(K)
	V := make(map[int][]int)
	for _, k := range K {
		c := C[k]
		for i, p := range c {
			V[i] = append(V[i], p)
		}
	}
	return V
}

func dataToCluster(D spn.Dataset) map[int][]int {
	C := make(map[int][]int)
	for i, I := range D {
		C[i] = make([]int, len(I))
		for j, v := range I {
			C[i][j] = v
		}
	}
	return C
}

func buildRegionGraph(D spn.Dataset, sc map[int]learn.Variable, k int, t float64) *graph {
	G := newGraph(sc)
	n := G.root
	if k > 1 {
		M := learn.DataToMatrix(D)
		C := cluster.KMedoid(k, M)
		for i := 0; i < k; i++ {
			//sys.Printf("Expanding region graph on cluster %d...\n", i)
			P := transpose(C[i])
			expandRegionGraph(G, n, P, t)
		}
	} else {
		C := dataToCluster(D)
		P := transpose(C)
		expandRegionGraph(G, n, P, t)
	}
	return G
}

func expandRegionGraph(G *graph, n *region, C map[int][]int, t float64) {
	sn := n.sc
	//sys.Printf("Partitioning scope of size %d...\n", len(sn))
	s1, s2 := partitionScope(sn, C)
	S := G.allScopes()
	//sys.Println("Trying to find similar scopes to S1 and S2...")
	for _, s := range S {
		if s.subsetOf(sn) && !s.equal(sn) {
			if s.similarTo(s1, s2, t) {
				s1, s2 = s, sn.minus(s)
				//sys.Printf("  Found similar scopes of size: %d, %d\n", len(s1), len(s2))
				break
			}
		}
	}
	//sys.Println("Validating regions from S1 and S2...")
	n1, n2 := G.validateRegion(s1), G.validateRegion(s2)
	if !G.existsPartition(n, n1, n2) {
		//sys.Println("No existing partition found. Creating new one...")
		p := newPartition()
		n.add(p)
		p.add(n1)
		p.add(n2)
		G.registerPartition(p, n, n1, n2)
	}
	if !S.contains(s1) && len(s1) > 1 {
		//sys.Println("Expanding on S1...")
		expandRegionGraph(G, n1, C, t)
	}
	if !S.contains(s2) && len(s2) > 1 {
		//sys.Println("Expanding on S2...")
		expandRegionGraph(G, n2, C, t)
	}
}

// Already transposes C and excludes any variables that are not present in s.
func clusterToMatrix(C map[int][]int, S scope) ([][]float64, map[int]int) {
	// S is always a subset of Sc(C), since C is complete. So we must restrict C wrt S and return a
	// matrix that has scope S.
	M := make([][]float64, len(S))
	V := make(map[int]int)
	var i int
	for k, v := range C {
		if _, e := S[k]; e {
			m := len(v)
			M[i] = make([]float64, m)
			for j := 0; j < m; j++ {
				M[i][j] = float64(v[j])
			}
			//copy(M[i], v)
			V[i] = k
			i++
		}
	}
	return M, V
}

func partitionScope(sn scope, C map[int][]int) (scope, scope) {
	M, V := clusterToMatrix(C, sn)
	//sys.Printf("%d, %d\n", len(sn), len(M))
	//sys.Println("Running k-means on variables...")
	S := cluster.KMeansF(2, M, metrics.EuclideanF)
	//sys.Println("Finished k-means.")
	//sys.Printf("  %d, %d, %d, %d\n", len(S[0]), len(S[1]), len(S[0])+len(S[1]), len(sn))
	s1, s2 := make(scope), make(scope)
	// Cluster/Partition 1
	for k, _ := range S[0] {
		// k is the variable ID, since M is the transpose of the dataset in matrix form.
		// If we took the regular cluster form of the dataset, k would be the instance index.
		s1[V[k]] = sn[V[k]]
	}
	// Cluster/Partition 2
	for k, _ := range S[1] {
		//sys.Printf("__ %d, %d\n", k, len(S[1][k]))
		//if _, e := sn[k]; !e {
		//sys.Printf("  %d does not exist in sn\n", k)
		//}
		s2[V[k]] = sn[V[k]]
	}
	//sys.Printf("Split scope into two partitions of size: %d, %d...\n", len(s1), len(s2))
	return s1, s2
}

func buildSPN(g *graph, D spn.Dataset, m, l int) spn.SPN {
	// Take the post-order first, since we go top-down.
	R := g.postorder()
	for _, r := range R {
		var N []spn.SPN
		if r == g.root {
			N = []spn.SPN{spn.NewSum()}
			r.rep = N
		} else {
			N = r.translate(D, m, l)
		}
		P := r.ch
		// If this for block executes, then we know that N is a set of sum nodes.
		for _, p := range P {
			C := p.ch
			// Assume |C|=2, since we partition scope into two.
			u, v := len(C[0].rep), len(C[1].rep)
			O := p.translate(u * v)
			//w := 1.0 / float64(len(O))
			// Add sum nodes from parent region to each product node in partition P.
			for _, n := range N {
				s := n.(*spn.Sum)
				for _, o := range O {
					//s.AddChildW(o, w)
					//s.AddChildW(o, 1.0)
					s.AddChildW(o, float64(sys.RandIntn(10)+1))
				}
			}
			// Add every combination of each child of P as children of each product node created.
			// We assume |C|=2 again, since we only partition the scope in two.
			var i int
			S1, S2 := C[0].rep, C[1].rep
			for _, c1 := range S1 {
				for _, c2 := range S2 {
					O[i].AddChild(c1)
					O[i].AddChild(c2)
					i++
				}
			}
		}
	}
	return g.root.rep[0]
}

func Structure(D spn.Dataset, sc map[int]learn.Variable, k, m, g int, t float64) spn.SPN {
	//sys.Println("Building region graph...")
	G := buildRegionGraph(D, sc, k, t)
	//sys.Println("Building SPN from region graph...")
	S := buildSPN(G, D, m, g)
	spn.NormalizeSPN(S)
	return S
}

func LearnGD(D spn.Dataset, sc map[int]learn.Variable, k, m, g int, t float64, P *parameters.P, i int) spn.SPN {
	S := Structure(D, sc, k, m, g, t)
	var ns, np, nl int
	spn.BreadthFirst(S, func(s spn.SPN) bool {
		switch t := s.Type(); t {
		case "sum":
			ns++
		case "product":
			np++
		default:
			nl++
		}
		return true
	})
	sys.Printf("Sum: %d, Products: %d, Leaves: %d, Total: %d\n", ns, np, nl, ns+np+nl)
	parameters.Bind(S, P)
	//spn.PrintSPN(S, fmt.Sprintf("test_before_%d.spn", i))
	return learn.Generative(S, D)
	//return S
}
