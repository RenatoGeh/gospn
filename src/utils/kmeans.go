package utils

import (
	"fmt"
	"math"
	"math/rand"
)

// Distancer guarantees that the object can be measured in relation to a mean.
type Distancer interface {
	Distance(mean float64) float64
}

// Little hack to allow interface Distancer to allow float64 and []float64.
type Float float64
type FSlice []float64

// Distance between value Float and a given float64 mean.
func (f Float) Distance(mean float64) float64 {
	dist := math.Abs(float64(f) - mean)
	return dist * dist
}

// Distance between the mean of a float slice and a mean.
func (v FSlice) Distance(mean float64) float64 {
	m, n := 0.0, len(v)
	for i := 0; i < n; i++ {
		m += v[i]
	}
	d := math.Abs(m/float64(n) - mean)
	return d * d
}

func KMeansV(k int, data [][]int) []map[int][]int {
	n := len(data)

	// Initializes using the Forgy method.
	fmt.Println("Initializing K-means clustering via the Forgy method...")
	chkrnd := make(map[int]bool)
	clusters := make([]map[int][]int, k)
	means := make([]float64, k)
	chkdata := make(map[int]int)
	for i := 0; i < k; i++ {
		var r int
		for ok := false; ok; _, ok = chkrnd[r] {
			r = rand.Intn(n)
		}
		m := data[r]
		clusters[i] = make(map[int][]int)
		chkrnd[r], chkdata[r] = true, i
		clusters[i][r] = m
		mean, s := 0.0, len(m)
		for j := 0; j < s; j++ {
			mean += float64(m[j])
		}
		means[i] = mean / float64(s)
	}

	diff, diffsum := make([]float64, k), 0.0

	fmt.Println("Starting K-means until convergence...")
	for diffsum != 0 {
		for i := 0; i < n; i++ {
			min, mean, s, which := math.Inf(1), 0.0, len(data[i]), -1
			for j := 0; j < s; j++ {
				mean += float64(data[i][j])
			}
			mean /= float64(s)
			for j := 0; j < k; j++ {
				t := Float(mean).Distance(means[j])
				if t < min {
					min, which = t, j
				}
			}
			v, ok := chkdata[i]
			// Instance i has no attached cluster.
			if !ok {
				chkdata[i] = which
				clusters[which][i] = data[i]
			} else if v != which {
				// If instance has an earlier attached cluster.
				delete(clusters[v], i)
				clusters[which][i] = data[i]
				chkdata[i] = which
			}
		}

		// Recompute means and diff.
		diffsum = 0
		for i := 0; i < k; i++ {
			mean, s := 0.0, 0
			for _, value := range clusters[i] {
				m := len(clusters[i])
				s += m
				for j := 0; j < m; j++ {
					mean += float64(value[j])
				}
			}
			md := mean / float64(s)
			diff[i] = math.Abs(means[i] - md)
			diffsum += diff[i]
			means[i] = md
		}
	}
	fmt.Println("Converged. Returning clusters...")

	return clusters
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
