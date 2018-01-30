package spn

import (
	//"fmt"
	"github.com/RenatoGeh/gospn/common"
	"math/rand"
	"testing"
)

const (
	rSeed  = 101
	maxInt = 50
)

func initGraph() []*graphTest {
	G := make([]*graphTest, 10)
	for i := range G {
		G[i] = &graphTest{}
		G[i].i = i
	}
	G[0].Add(G[1], G[2], G[3])
	G[1].Add(G[5], G[6])
	G[2].Add(G[4])
	G[3].Add(G[4], G[5])
	G[5].Add(G[8])
	G[6].Add(G[7])
	G[7].Add(G[9])
	G[8].Add(G[9])
	return G
}

// Topological sorting may not exist in this graph.
func createRandomGraph(seed int64) []*graphTest {
	R := rand.New(rand.NewSource(seed))
	n := R.Intn(maxInt) + 1
	G := make([]*graphTest, n)
	for i := 0; i < n; i++ {
		G[i] = &graphTest{}
		G[i].i = i
	}

	for i := 0; i < n; i++ {
		// Not too dense.
		k := rand.Intn(n / 4)
		for j := 0; j < k; j++ {
			c := rand.Intn(n)
			for c == i {
				c = rand.Intn(n)
			}
			G[i].Add(G[c])
		}
	}
	return G
}

// There is always at least one topological sorting in a DAG.
func generateDAG(seed int64) []*graphTest {
	const (
		maxNodesPerLayer = 20
		maxLayers        = 30
	)
	R := rand.New(rand.NewSource(seed))
	h := R.Intn(maxLayers) + 1
	L := make([][]*graphTest, h)
	var G []*graphTest
	var c int
	for i := 0; i < h; i++ {
		n := R.Intn(maxNodesPerLayer) + 1
		L[i] = make([]*graphTest, n)
		for j := 0; j < n; j++ {
			v := &graphTest{}
			L[i][j] = v
			v.i = c
			c++
			G = append(G, v)
		}
	}
	for i := 0; i < h-1; i++ {
		k := len(L[i+1])
		n := len(L[i])
		for j := 0; j < n; j++ {
			l := R.Intn(k) + 1
			U := R.Perm(k)
			for t := 0; t < l; t++ {
				L[i][j].Add(L[i+1][U[t]])
			}
		}
	}
	return G
}

func TestTopSortTarjanSmall(t *testing.T) {
	G := initGraph()
	tSort := []int{9, 8, 5, 4, 3, 2, 7, 6, 1, 0}
	Q := common.Queue{}
	TopSortTarjan(G[0], &Q)
	var i int
	for !Q.Empty() {
		u := Q.Dequeue().(*graphTest)
		if tSort[i] != u.i {
			t.Errorf("Expected %d, got %d.", tSort[i], u.i)
		}
		i++
	}
}

func TestTopSortTarjanRecSmall(t *testing.T) {
	G := initGraph()
	tSort := []int{9, 8, 5, 7, 6, 1, 4, 2, 3, 0}
	Q := common.Queue{}
	TopSortTarjanRec(G[0], &Q)
	var i int
	for !Q.Empty() {
		u := Q.Dequeue().(*graphTest)
		if tSort[i] != u.i {
			t.Errorf("Expected %d, got %d.", tSort[i], u.i)
		}
		i++
	}
}

func reaches(u, v *graphTest) bool {
	V := make(map[*graphTest]bool)
	Q := common.Queue{}
	Q.Enqueue(u)
	V[u] = true
	for !Q.Empty() {
		s := Q.Dequeue().(*graphTest)
		if s == v {
			return true
		}
		ch := s.Ch()
		for _, c := range ch {
			k := c.(*graphTest)
			if !V[k] {
				Q.Enqueue(k)
				V[k] = true
			}
		}
	}
	return false
}

func reachable(u *graphTest) *common.Queue {
	V := make(map[*graphTest]bool)
	Q := common.Queue{}
	R := &common.Queue{}
	Q.Enqueue(u)
	V[u] = true
	for !Q.Empty() {
		s := Q.Dequeue().(*graphTest)
		ch := s.Ch()
		for _, c := range ch {
			k := c.(*graphTest)
			if !V[k] {
				Q.Enqueue(k)
				R.Enqueue(k)
				V[k] = true
			}
		}
	}
	return R
}

func unitTest(seed int64, sort func(SPN, common.Collection) common.Collection) bool {
	G := generateDAG(seed)
	//G := initGraph()
	Q := common.Queue{}
	sort(G[0], &Q)
	V := make(map[*graphTest]bool)

	//fmt.Println("Topo sort:")
	for !Q.Empty() {
		u := Q.Dequeue().(*graphTest)
		//fmt.Printf("  %d\n", u.i)
		V[u] = true
		R := reachable(u)
		for !R.Empty() {
			v := R.Dequeue().(*graphTest)
			if !V[v] {
				return false
			}
		}
	}
	//fmt.Println("")

	return true
}

func TestTopSortTarjan(t *testing.T) {
	if !unitTest(rSeed, TopSortTarjan) {
		t.Error("Expected true, got false.")
	}
}

func TestTopSortTarjanRec(t *testing.T) {
	if !unitTest(rSeed, TopSortTarjanRec) {
		t.Error("Expected true, got false.")
	}
}

func TestTopSortDFS(t *testing.T) {
	if !unitTest(rSeed, TopSortDFS) {
		t.Error("Expected true, got false.")
	}
}
