package cluster

import (
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/utils/cluster/metrics"
)

func kMedoidInsert(which int, means []int, clusters []map[int][]int, v []int, i int) {
	clusters[which][i] = make([]int, len(v))
	copy(clusters[which][i], v)
	smean, s := 0, 0
	for _, value := range clusters[which] {
		smean += metrics.Hamming(clusters[which][means[which]], value)
		s += metrics.Hamming(v, value)
	}
	if s < smean {
		means[which] = i
	}
}

func kMedoidRemove(which int, means []int, clusters []map[int][]int, v []int, i int) {
	delete(clusters[which], i)
	if means[which] == i {
		best := 2000000000
		for j, value := range clusters[which] {
			s := 0
			for _, v := range clusters[which] {
				s += metrics.Hamming(v, value)
			}
			if s < best {
				best = s
				means[which] = j
			}
		}
	}
}

func KMedoid(k int, data [][]int) []map[int][]int {
	n := len(data)

	// Initializes using the Forgy method.
	//fmt.Println("Initializing K-means clustering via the Forgy method...")
	chkrnd := make(map[int]bool)
	clusters := make([]map[int][]int, 1)
	means := make([]int, 1)
	chkdata := make(map[int]int)
	for i := 0; i < k; i++ {
		var r int
		ok := true
		for ok && len(chkrnd) < n {
			for ok = true; ok; _, ok = chkrnd[r] {
				r = sys.RandIntn(n)
			}
			chkrnd[r] = true
			for ii := 0; ii < i && !ok; ii++ {
				lr := len(data[r])
				j := 0
				for j < lr {
					if data[r][j] != data[means[ii]][j] {
						break
					}
					j++
				}
				if j >= lr {
					ok = true
				}
			}
		}
		if ok {
			break
		}
		//fmt.Printf("medoid %d %d\n", i, r)
		if i > 0 {
			clusters = append(clusters, make(map[int][]int))
			means = append(means, r)
		} else {
			clusters[0] = make(map[int][]int)
			means[0] = r
		}
		chkdata[r] = i
		kMedoidInsert(i, means, clusters, data[r], r)
	}
	//fmt.Println("k", k, "n", n)

	//fmt.Println("Starting K-means until convergence...")
	nochange := 0
	i := 0
	for nochange < n {
		min, which := len(data[i])+1, -1
		if v, ok := chkdata[i]; ok {
			which = v
			min = metrics.Hamming(clusters[which][means[which]], data[i])
		}
		for j := 0; j < k; j++ {
			if j != which {
				t := metrics.Hamming(clusters[j][means[j]], data[i])
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
			kMedoidInsert(which, means, clusters, data[i], i)
			nochange = 0
		} else if v != which {
			// If instance has an earlier attached cluster.
			//fmt.Println(data[i], " from ", chkdata[i], " to cluster ", which)
			kMedoidRemove(chkdata[i], means, clusters, data[i], i)
			chkdata[i] = which
			kMedoidInsert(which, means, clusters, data[i], i)
			nochange = 0
		} else {
			nochange++
		}
		i++
		if i >= n {
			i = 0
		}
		//fmt.Println("0:",clusters[0][means[0]], "  1:", clusters[1][means[1]])
	}
	//fmt.Println("Converged. Returning clusters...")
	clusters = make([]map[int][]int, k)
	for i = 0; i < k; i++ {
		clusters[i] = make(map[int][]int)
	}
	i = 0
	for i < n {
		//fmt.Println("i",i,"chkdata",chkdata[i],"len",len(data[i]))
		clusters[chkdata[i]][i] = make([]int, len(data[i]))
		copy(clusters[chkdata[i]][i], data[i])
		i++
	}
	return clusters
}
