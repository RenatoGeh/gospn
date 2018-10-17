package cluster

import (
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/utils"
	"github.com/RenatoGeh/gospn/utils/cluster/metrics"
)

func dbscanInternal(data [][]float64, eps float64, mp int) []*utils.UFNode {
	n := len(data)
	// Metric function.
	mfunc := metrics.EuclideanF

	// Distance matrix.
	dmatrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		dmatrix[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			dmatrix[i][j] = mfunc(data[i], data[j])
		}
	}

	// Regions.
	rgs := make([]*utils.UFNode, n)
	for i := 0; i < n; i++ {
		rgs[i] = utils.MakeSet(i)
	}

	// Visited points: 0 unvisited, 1 otherwise.
	vst, vindex := make([]int, n), 0

	queue := common.Queue{}
	queue.Enqueue(0)
	for !queue.Empty() {
		p := queue.Dequeue().(int)

		// Neighbourhood of p.
		nbh := common.Queue{}

		for i := 0; i < n; i++ {
			// Clause 1 (i != p):
			//  Pairs must be distinct.
			// Clause 2 (dmatrix[p][i] <= eps):
			//  Distance must be <= the epsilon parameter of max distance.
			// Clause 3 (utils.Find(rgs[i]) != utils.Find(rgs[p])):
			//  Pair is not already in the same cluster.
			if (i != p) && (dmatrix[p][i] <= eps) && (utils.Find(rgs[i]) != utils.Find(rgs[p])) {
				nbh.Enqueue(i)
			}
		}

		// Found dense neighbourhood.
		if nbh.Size() >= mp {
			for !nbh.Empty() {
				q := nbh.Dequeue().(int)
				utils.Union(rgs[p], rgs[q])
				vst[p], vst[q] = 1, 1
				queue.Enqueue(q)
			}
		}

		// Cluster has been formed. Select next non-clustered region.
		if queue.Empty() {
			for i := vindex; i < n; i++ {
				if vst[i] == 0 {
					queue.Enqueue(i)
					vindex = i + 1
				}
			}
		}
	}

	return rgs
}

// DBSCAN Density-based spatial clustering of applications with noise.
// Parameters:
//  - data is data matrix;
//  - eps is epsilon maximum distance between density core points;
//  - mp is minimum number of points to be considered core point.
func DBSCAN(data [][]int, eps float64, mp int) []map[int][]int {
	D := copyMatrixF(data)
	rgs := dbscanInternal(D, eps, mp)
	n, m := len(D), len(D[0])
	// Convert Union-Find format to []map[int][]int format.
	k := 0
	var clusters []map[int][]int
	for i := 0; i < n; i++ {
		if rgs[i].Pa == rgs[i] {
			clusters = append(clusters, make(map[int][]int))
			chs := utils.UFVarids(rgs[i])
			nchs := len(chs)
			for j := 0; j < nchs; j++ {
				l := chs[j]
				clusters[k][l] = make([]int, m)
				copy(clusters[k][l], data[l])
			}
			k++
		}
	}

	return clusters
}

func DBSCANData(data []map[int]int, eps float64, mp int) [][]map[int]int {
	M, Sc := learn.DataToMatrixF(data)
	rgs := dbscanInternal(M, eps, mp)
	G := make([]int, len(M))
	var k int
	for i, u := range rgs {
		if u.Pa == u {
			ch := utils.UFVarids(u)
			for _, c := range ch {
				G[c] = i
			}
			k++
		}
	}
	return guessToData(k, G, M, Sc)
}
