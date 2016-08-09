package cluster

import (
//common "github.com/RenatoGeh/gospn/src/common"
//utils "github.com/RenatoGeh/gospn/src/utils"
//metrics "github.com/RenatoGeh/gospn/src/utils/cluster/metrics"
)

// Ordering points to identify the clustering structure (OPTICS).
// OPTICS is similar to DBSCAN with the exception that instead of an epsilon to bound the distance
// between points, OPTICS replaces that epsilon with a new epsilon that upper bounds the maximum
// possible epsilon a DBSCAN would take.
// Parameters:
//  - data is data matrix;
//  - eps is a maximum distance between density core points upper bound;
//  - mp is minium number of points to be considered core point.
func OPTICS(data [][]int, eps float64, mp int) []map[int][]int {
	//n, m := len(data), len(data[0])
	//mfunc := metrics.Euclidean

	return nil
}
