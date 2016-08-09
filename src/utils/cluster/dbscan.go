package cluster

import (
	common "github.com/RenatoGeh/gospn/src/common"
	utils "github.com/RenatoGeh/gospn/src/utils"
	metrics "github.com/RenatoGeh/gospn/src/utils/cluster/metrics"
)

// Density-based spatial clustering of applications with noise (DBSCAN).
// Parameters:
//  - data is data matrix;
//  - eps is epsilon maximum distance between density core points;
//  - mp is minimum number of points to be considered core point.
func DBSCAN(data [][]int, eps float64, mp int) []map[int][]int {
	n, m := len(data), len(data[0])
	// Metric function.
	mfunc := metrics.Euclidean

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

	queue := common.QueueInteger{}
	queue.Enqueue(0)
	for !queue.Empty() {
		p := queue.Dequeue()

		// Neighbourhood of p.
		nbh := common.QueueInteger{}

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
				q := nbh.Dequeue()
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

	// Convert Union-Find format to []map[int][]int format.
	k := 0
	var clusters []map[int][]int = nil
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
