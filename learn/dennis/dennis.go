package dennis

import (
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/utils/cluster"
	"sort"
)

func buildRegionGraph(D spn.Dataset, sc map[int]learn.Variable, k int, t float64) *graph {
	M := learn.DataToMatrix(D)
	C := cluster.KMeans(k, M)
	G := newGraph(sc)
	n := G.root
	for i := 0; i < k; i++ {
		expandRegionGraph(G, n, C[i], t)
	}
	return G
}

func expandRegionGraph(G *graph, n *region, C map[int][]int, t float64) {
	sn := n.sc
	s1, s2 := partitionScope(n, C)
	S := G.allScopes()
	for _, s := range S {
		if s.subsetOf(sn) {
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
	if S.contains(s1) && len(s1) > 1 {
		expandRegionGraph(G, n1, C, t)
	}
	if S.contains(s2) && len(s2) > 1 {
		expandRegionGraph(G, n2, C, t)
	}
}

// Already transposes C and excludes any variables that are not present in s.
func clusterToMatrix(C map[int][]int, s scope) [][]int {
	// Sc(C) is always a subset of s, since C is complete. So we must restrict C wrt s and return a
	// matrix that has scope s.
	m := len(s)
	S := make([]int, m)
	var l int
	for k, _ := range s {
		S[l] = k
		l++
	}
	sort.Ints(S)
	M := make([][]int, len(C))
	for i := range C {
		M[i] = make([]int, m)
		for j, k := range S {
			M[i][j] = C[i][k]
		}
	}
	return M
}

func partitionScope(r *region, C map[int][]int) (scope, scope) {
	sn := r.sc
	M := clusterToMatrix(C, sn)
	S := cluster.KMeans(2, M)
	s1, s2 := make(scope), make(scope)
	// Cluster/Partition 1
	for k, _ := range S[0] {
		// k is the variable ID, since M is the transpose of the dataset in matrix form.
		// If we took the regular cluster form of the dataset, k would be the instance index.
		s1[k] = sn[k]
	}
	// Cluster/Partition 2
	for k, _ := range S[1] {
		s2[k] = sn[k]
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
				for _, o := range O {
					s := n.(*spn.Sum)
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
	G := buildRegionGraph(D, sc, k, t)
	S := buildSPN(G, D, m)
	return S
}

func LearnGD(D spn.Dataset, sc map[int]learn.Variable, k, m int, t, eta, eps float64, norm bool) spn.SPN {
	S := Structure(D, sc, k, m, t)
	return learn.GenerativeHardGD(S, eta, eps, D, nil, norm)
}
