package dennis

import (
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/utils/cluster"
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

func buildRegionGraph(D spn.Dataset, sc map[int]learn.Variable, k int, t float64) *graph {
	M := learn.DataToMatrix(D)
	C := cluster.KMeans(k, M)
	G := newGraph(sc)
	n := G.root
	for i := 0; i < k; i++ {
		sys.Printf("Expanding region graph on cluster %d...\n", i)
		P := transpose(C[i])
		expandRegionGraph(G, n, P, t)
	}
	return G
}

func expandRegionGraph(G *graph, n *region, C map[int][]int, t float64) {
	sn := n.sc
	s1, s2 := partitionScope(sn, C)
	S := G.allScopes()
	for _, s := range S {
		if s.subsetOf(sn) && !s.equal(sn) {
			if s.similarTo(s1, s2, t) {
				s1, s2 = s, sn.minus(s)
				break
			}
		}
	}
	n1, n2 := G.validateRegion(s1), G.validateRegion(s2)
	if !G.existsPartition(n, n1, n2) {
		p := newPartition()
		n.add(p)
		p.add(n1)
		p.add(n2)
		G.registerPartition(p, n, n1, n2)
	}
	if !S.contains(s1) && len(s1) > 1 {
		expandRegionGraph(G, n1, C, t)
	}
	if !S.contains(s2) && len(s2) > 1 {
		expandRegionGraph(G, n2, C, t)
	}
}

// Already transposes C and excludes any variables that are not present in s.
func clusterToMatrix(C map[int][]int, S scope) ([][]int, map[int]int) {
	// S is always a subset of Sc(C), since C is complete. So we must restrict C wrt S and return a
	// matrix that has scope S.
	M := make([][]int, len(S))
	V := make(map[int]int)
	var i int
	for k, v := range C {
		if _, e := S[k]; e {
			M[i] = make([]int, len(v))
			copy(M[i], v)
			V[i] = k
			i++
		}
	}
	return M, V
}

func partitionScope(sn scope, C map[int][]int) (scope, scope) {
	M, V := clusterToMatrix(C, sn)
	//sys.Printf("%d, %d\n", len(sn), len(M))
	S := cluster.KMeans(2, M)
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
	return s1, s2
}

func buildSPN(g *graph, D spn.Dataset, m int) spn.SPN {
	// Take the post-order first, since we go top-down.
	R := g.postorder()
	for _, r := range R {
		var N []spn.SPN
		if r == g.root {
			N = []spn.SPN{spn.NewSum()}
			r.rep = N
		} else {
			N = r.translate(D, m)
		}
		P := r.ch
		// If this for block executes, then we know that N is a set of sum nodes.
		for _, p := range P {
			C := p.ch
			// Assume |C|=2, since we partition scope into two.
			O := p.translate(m * m)
			w := 1.0 / float64(len(O))
			// Add sum nodes from parent region to each product node in partition P.
			for _, n := range N {
				s := n.(*spn.Sum)
				for _, o := range O {
					s.AddChildW(o, w)
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

func Structure(D spn.Dataset, sc map[int]learn.Variable, k, m int, t float64) spn.SPN {
	sys.Println("Building region graph...")
	G := buildRegionGraph(D, sc, k, t)
	sys.Println("Building SPN from region graph...")
	S := buildSPN(G, D, m)
	return S
}

func LearnGD(D spn.Dataset, sc map[int]learn.Variable, k, m int, t, eta, eps float64, norm bool) spn.SPN {
	S := Structure(D, sc, k, m, t)
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
	return learn.GenerativeHardBGD(S, eta, eps, D, nil, norm, 50)
}
