package cluster

import (
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/utils/cluster/metrics"
	"math"
)

func kMeansInsert(which int, means map[int][]int, clusters []map[int][]int, v []int) {
	l := len(means[which])
	for k := 0; k < l; k++ {
		//fmt.Printf("which %d k %d v[k] %d l %d len(clusters[which][k]) %d\n", which, k, v[k], l, len(clusters[which][k]))
		clusters[which][k][v[k]]++
		if clusters[which][k][v[k]] > clusters[which][k][means[which][k]] {
			means[which][k] = v[k]
		}
	}
}

func kMeansRemove(which int, means map[int][]int, clusters []map[int][]int, v []int) {
	l := len(means[which])
	for k := 0; k < l; k++ {
		clusters[which][k][v[k]]--
		if means[which][k] == v[k] {
			max, s := -1, len(clusters[which][k])
			for i := 0; i < s; i++ {
				if max < 0 || clusters[which][k][i] > max {
					max = i
				}
			}
			means[which][k] = max
		}
	}
}

func KMeans(k int, data [][]int) []map[int][]int {
	n := len(data)
	F := metrics.Hamming

	// Initializes using the Forgy method.
	//sys.Println("Initializing K-means clustering via the Forgy method...")
	chkrnd := make(map[int]bool)
	clusters := make([]map[int][]int, k)
	means := make(map[int][]int, k)
	chkdata := make(map[int]int)
	for i := 0; i < k; i++ {
		var r int
		for ok := true; ok; _, ok = chkrnd[r] {
			r = sys.RandIntn(n)
		}
		//fmt.Printf("%d vs %d\n", n, r)
		clusters[i] = make(map[int][]int)
		s := len(data[r])
		means[i] = make([]int, s)
		copy(means[i], data[r])
		for j := 0; j < s; j++ {
			max := -1
			for z := 0; z < n; z++ {
				if data[z][j] > max {
					max = data[z][j]
				}
			}
			clusters[i][j] = make([]int, max+1)
			//			for z := 0; z < max+1; z++ {
			//				clusters[i][j][z] = 0
			//			}
		}
		chkrnd[r], chkdata[r] = true, i
		kMeansInsert(i, means, clusters, data[r])
	}

	//sys.Println("Starting K-means until convergence...")
	nochange := 0
	i := 0
	for nochange < n {
		min, which := len(data[i])+1, -1
		if v, ok := chkdata[i]; ok {
			which = v
			min = F(means[which], data[i])
		}
		for j := 0; j < k; j++ {
			if j != which {
				t := F(means[j], data[i])
				if t < min {
					min, which = t, j
				}
			}
		}
		v, ok := chkdata[i]
		// Instance i has no attached cluster.
		if !ok {
			chkdata[i] = which
			//fmt.Println(data[i], " to cluster ", which)
			kMeansInsert(which, means, clusters, data[i])
			nochange = 0
		} else if v != which {
			// If instance has an earlier attached cluster.
			//fmt.Println(data[i], " from ", chkdata[i], " to cluster ", which)
			kMeansRemove(chkdata[i], means, clusters, data[i])
			chkdata[i] = which
			kMeansInsert(which, means, clusters, data[i])
			nochange = 0
		} else {
			nochange++
		}
		i++
		if i >= n {
			i = 0
		}
		//fmt.Println("0:", means[0], "  1:", means[1])
	}
	//sys.Println("Converged. Returning clusters...")
	clusters = make([]map[int][]int, k)
	for i = 0; i < k; i++ {
		clusters[i] = make(map[int][]int)
	}
	for i := 0; i < n; i++ {
		clusters[chkdata[i]][i] = make([]int, len(data[i]))
		copy(clusters[chkdata[i]][i], data[i])
	}
	return clusters
}

func KMeansF(k int, D [][]float64, F metrics.MetricF) []map[int][]float64 {
	if len(D) < k {
		panic("Length of dataset is smaller than number of clusters!")
	}

	n := len(D[0]) // Instance dimension size.
	m := len(D)
	M := make([][]float64, k) // Cluster means (centroids at each cluster).
	C := make([]int, m)       // Maps an instance to a cluster.
	for i := range M {
		M[i] = make([]float64, n)
	}
	for i := 0; i < m; i++ {
		C[i] = -1
	}
	R := make(map[int]bool)
	//sys.Println("Initialized through Forgy.")
	//sys.Printf("%d\n", len(D))
	// Forgy initialization.
	for i := 0; i < k; i++ {
		r := sys.RandIntn(m)
		for _, e := R[r]; e; _, e = R[r] {
			r = sys.RandIntn(m)
		}
		R[r] = true
		C[r] = i
		copy(M[i], D[r])
	}

	S := make([]int, k) // Number of instances in a cluster.
	//sys.Println("Running until convergence...")
	// Continue until there is no change (converges). If c = false, there was no change.
	for c := true; c; {
		c = false
		// Assignment step.
		for i, d := range D {
			kmin, min := C[i], math.Inf(1)
			if kmin != -1 {
				min = F(d, M[kmin])
			}
			for j := 0; j < k; j++ {
				u := F(d, M[j])
				if u < min {
					kmin, min = j, u
				}
			}
			if C[i] != kmin {
				c = true
				C[i] = kmin
			}
		}
		// Update step.
		for i := range M {
			for j := 0; j < n; j++ {
				M[i][j] = 0
			}
			S[i] = 0
		}
		for i := 0; i < m; i++ {
			p := C[i] // Which cluster instance i belongs to.
			S[p]++
			for j := 0; j < n; j++ {
				M[p][j] += D[i][j]
			}
		}
		for i := 0; i < k; i++ {
			for j := 0; j < n; j++ {
				if S[i] > 0 {
					M[i][j] /= float64(S[i])
				}
			}
		}
	}

	//sys.Println("Converged. Converting to output format...")
	// Convert to output format.
	O := make([]map[int][]float64, k)
	for i := range O {
		O[i] = make(map[int][]float64)
	}
	for i := range D {
		c := C[i]
		O[c][i] = make([]float64, n)
		copy(O[c][i], D[i])
	}
	return O
}
