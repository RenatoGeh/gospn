package dennis

import (
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/utils/cluster"
	"sort"
)

func buildRegionGraph(D spn.Dataset, sc map[int]learn.Variable, k int, t float64) {
	M := learn.DataToMatrix(D)
	C := cluster.KMeans(k, M)
	G := newGraph(sc)
	n := G.root
	for i := 0; i < k; i++ {
		expandRegionGraph(G, n, C[i], t)
	}
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
