package cluster

import (
	"container/heap"
	"math"
	"sort"
	//common "github.com/RenatoGeh/gospn/common"
	//utils "github.com/RenatoGeh/gospn/utils"
	"github.com/RenatoGeh/gospn/utils/cluster/metrics"
	"github.com/mpraski/clusters"
)

func optics(D [][]float64, mp int, eps, xi float64, F metrics.MetricF) []int {
	ci, err := clusters.OPTICS(mp, eps, xi, 0, metrics.EuclideanF)
	if err != nil {
		panic(err)
	}
	if err = ci.Learn(D); err != nil {
		panic(err)
	}
	return ci.Guesses()
}

func OPTICS2(D [][]int, mp int, eps, xi float64, F metrics.MetricF) []map[int][]int {
	E := copyMatrixF(D)
	G := optics(E, mp, eps, xi, F)
	m := -2
	for _, g := range G {
		if g > m {
			m = g
		}
	}
	return toCluster(m, D, G)
}

var distance = metrics.Euclidean
var undef = math.Inf(1)

type object struct {
	data  []int
	vst   bool
	index int
	rdist float64
	cdist float64
}

func (o object) get(index int) int {
	return o.data[index]
}

// Queue of objects.
type qObj struct {
	data []*object
}

func (q *qObj) enqueue(e *object) {
	q.data = append(q.data, e)
}
func (q *qObj) dequeue() *object {
	n := len(q.data)
	e := q.data[0]
	q.data = q.data[1:n]
	return e
}
func (q *qObj) peek() *object {
	return q.data[0]
}
func (q *qObj) get(i int) *object {
	return q.data[i]
}
func (q *qObj) size() int   { return len(q.data) }
func (q *qObj) empty() bool { return len(q.data) == 0 }

// End of queue of objects.

func getNeighbors(set []*object, o *object, eps float64) []*object {
	var nset []*object
	n := len(set)
	p1 := o.data
	for i := 0; i < n; i++ {
		if o.index == i {
			continue
		} else if distance(p1, set[i].data) <= eps {
			nset = append(nset, set[i])
		}
	}
	return nset
}

func max(n float64, m float64) float64 {
	if n > m {
		return n
	}
	return m
}

type objDist struct {
	dist  float64
	index int
}
type objDists []*objDist

func (od objDists) Len() int           { return len(od) }
func (od objDists) Swap(i, j int)      { od[i], od[j] = od[j], od[i] }
func (od objDists) Less(i, j int) bool { return od[i].dist < od[j].dist }

func getCoreDist(nb []*object, o *object, eps float64, mp int) float64 {
	n := len(nb)
	if n < mp {
		return undef
	}
	dists, p1 := make([]*objDist, n), o.data
	for i := 0; i < n; i++ {
		dists = append(dists, &objDist{dist: distance(p1, nb[i].data), index: i})
	}
	sort.Sort(objDists(dists))
	retval := dists[mp-1].dist
	dists = nil
	return retval
}

// Priority queue.
type item struct {
	obj   *object
	p     float64
	index int
}
type pQueue []*item

func (pq pQueue) Len() int { return len(pq) }
func (pq pQueue) Less(i, j int) bool {
	return pq[i].p < pq[j].p
}
func (pq pQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *pQueue) Push(x interface{}) {
	n := len(*pq)
	i := x.(*item)
	i.index = n
	*pq = append(*pq, i)
}
func (pq *pQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	i := old[n-1]
	i.index = -1 // for safety
	*pq = old[0 : n-1]
	return i
}
func (pq *pQueue) update(value *object, priority float64) {
	n, q := len(*pq), *pq
	for i := 0; i < n; i++ {
		if q[i].obj == value {
			q[i].obj = value
			q[i].p = priority
			heap.Fix(pq, q[i].index)
			return
		}
	}
}

func seedsUpdate(nb []*object, o *object, pq *pQueue) {
	cdist, n, p1 := o.cdist, len(nb), o.data
	for i := 0; i < n; i++ {
		if !nb[i].vst {
			rdist := max(cdist, distance(p1, nb[i].data))
			if o.rdist == undef {
				o.rdist = rdist
				pq.Push(&item{obj: o, p: rdist})
			} else if rdist < o.rdist {
				o.rdist = rdist
				pq.update(o, rdist)
			}
		}
	}
}

func expand(set []*object, o *object, eps float64, mp int, order *qObj, pq *pQueue) {
	nb := getNeighbors(set, o, eps)
	o.vst = true
	o.rdist = undef
	o.cdist = getCoreDist(nb, o, eps, mp)
	order.enqueue(o)
	if o.cdist != undef {
		seedsUpdate(nb, o, pq)
		for len(*pq) != 0 {
			co := pq.Pop().(*item).obj
			nb := getNeighbors(set, co, eps)
			co.vst = true
			co.cdist = getCoreDist(nb, co, eps, mp)
			order.enqueue(co)
			if co.cdist != undef {
				seedsUpdate(nb, co, pq)
			}
		}
	}
}

func extract(order *qObj, eps float64, mp int) []map[int][]int {
	idtrk := 0
	cids := make([][]*object, 1)

	n := order.size()
	const noise = 0
	for i := 0; i < n; i++ {
		o := order.get(i)
		if o.rdist > eps {
			if o.cdist <= eps {
				idtrk++
				stub := make([]*object, 1)
				stub[0] = o
				cids = append(cids, stub)
			} else {
				cids[noise] = append(cids[noise], o)
			}
		} else {
			cids[idtrk] = append(cids[idtrk], o)
		}
	}

	m := len(cids)
	clusters := make([]map[int][]int, m)
	for i := 0; i < m; i++ {
		clusters[i] = make(map[int][]int)
		p := len(cids[i])
		for j := 0; j < p; j++ {
			o := cids[i][j]
			clusters[i][o.index] = o.data
		}
	}

	return clusters
}

// OPTICS - Ordering points to identify the clustering structure (OPTICS).
// OPTICS is similar to DBSCAN with the exception that instead of an epsilon to bound the distance
// between points, OPTICS replaces that epsilon with a new epsilon that upper bounds the maximum
// possible epsilon a DBSCAN would take.
// Parameters:
//  - data is data matrix;
//  - eps is a maximum distance between density core points upper bound;
//  - mp is minimum number of points to be considered core point.
func OPTICS(data [][]int, eps float64, mp int) []map[int][]int {
	n := len(data)
	order := qObj{}
	set := make([]*object, n)
	var pq pQueue
	for i := 0; i < n; i++ {
		set[i] = &object{data: data[i], vst: false, index: i, rdist: undef, cdist: undef}
	}
	heap.Init(&pq)

	for i := 0; i < n; i++ {
		obj := set[i]
		if !obj.vst {
			expand(set, obj, eps, mp, &order, &pq)
		}
	}

	return extract(&order, eps, mp)
}
