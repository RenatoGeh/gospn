package utils

import (
	"math"
	"math/rand"
)

// Distancer guarantees that the object can be measured in relation to a mean.
type Distancer interface {
	Distance(mean float64) float64
}

// Little hack to allow interface Distancer to allow float64.
type Float float64

// Distance between value Float and a given float64 mean.
func (f Float) Distance(mean float64) float64 {
	dist := math.Abs(float64(f) - mean)
	return dist * dist
}

// KMeans clusters data into k clusters.
// Returns k slices each containing a map of elements belonging to their corresponding clusters,
// where this map has keys corresponding to indeces of instances in data and values as the actual
// instance values.
func KMeans(k int, data []int) []map[int]int {
	n := len(data)

	// Initializes using the Forgy method.
	chkrnd := make(map[int]bool)
	clusters := make([]map[int]int, k)
	means := make([]float64, k)
	for i := 0; i < k; i++ {
		var r int
		for ok := true; ok; _, ok = chkrnd[r] {
			r = rand.Intn(n)
		}
		m := data[r]
		clusters[i] = make(map[int]int)
		// Key is index r of data instance. Value is the actual value of r.
		clusters[i][r] = m
		means[i] = float64(m)
		chkrnd[r] = true
	}

	// Difference between iterations. If diff = {0,...,0}, then converged. Else repeat.
	diff := make([]float64, k)
	var diffsum float64 = 1

	// Key is a certain instance i. Value is in which cluster 0 <= j <= k instance i is.
	chkdata := make(map[int]int)

	for diffsum != 0 {
		// Update clusters comparing each instance with each cluster's mean.
		for i := 0; i < n; i++ {
			min := math.Inf(1)
			ind := -1
			for j := 0; j < k; j++ {
				dist := Float(data[i]).Distance(means[j])
				if dist < min {
					min = dist
					ind = j
				}
			}
			cl, ok := chkdata[i]
			if !ok {
				clusters[ind][i] = data[i]
				chkdata[i] = ind
			} else if cl != ind {
				// Deletes key i from the cluster it was assigned to.
				delete(clusters[cl], i)
				clusters[ind][i] = data[i]
				chkdata[i] = ind
			}
		}

		// Recompute means.
		for i := 0; i < k; i++ {
			// Iterate over map clusters[i]
			var m float64 = 0
			for _, value := range clusters[i] {
				m += float64(value)
			}
			// Compute the centroid of each cluster.
			m /= float64(len(clusters[i]))
			diff[i] = means[i] - m
			means[i] = m
		}

		// Check for convergence (i.e. no change in means).
		diffsum = 0
		for i := 0; i < k; i++ {
			diffsum += math.Abs(diff[i])
		}
	}

	return clusters
}
