package cluster

import (
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/utils/cluster/metrics"
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

	// Initializes using the Forgy method.
	//fmt.Println("Initializing K-means clustering via the Forgy method...")
	chkrnd := make(map[int]bool)
	clusters := make([]map[int][]int, k)
	means := make(map[int][]int, k)
	chkdata := make(map[int]int)
	for i := 0; i < k; i++ {
		var r int
		for ok := true; ok; _, ok = chkrnd[r] {
			r = sys.Random.Intn(n)
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

	//fmt.Println("Starting K-means until convergence...")
	nochange := 0
	i := 0
	for nochange < n {
		min, which := len(data[i])+1, -1
		if v, ok := chkdata[i]; ok {
			which = v
			min = metrics.Hamming(means[which], data[i])
		}
		for j := 0; j < k; j++ {
			if j != which {
				t := metrics.Hamming(means[j], data[i])
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
	//fmt.Println("Converged. Returning clusters...")
	clusters = make([]map[int][]int, k)
	for i = 0; i < k; i++ {
		clusters[i] = make(map[int][]int)
	}
	i = 0
	for i < n {
		clusters[chkdata[i]][i] = make([]int, len(data[i]))
		copy(clusters[chkdata[i]][i], data[i])
		i++
	}
	return clusters
}
