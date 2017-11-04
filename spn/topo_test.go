package spn

import (
	"testing"
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

func TestTopSortTarjan(t *testing.T) {
	G := initGraph()
	tSort := []int{9, 8, 5, 4, 3, 2, 7, 6, 1, 0}
	Q := TopSortTarjan(G[0])
	var i int
	for !Q.Empty() {
		u := Q.Dequeue().(*graphTest)
		if tSort[i] != u.i {
			t.Errorf("Expected %d, got %d.", tSort[i], u.i)
		}
		i++
	}
}

func TestTopSortTarjanRec(t *testing.T) {
	G := initGraph()
	tSort := []int{9, 8, 5, 7, 6, 1, 4, 2, 3, 0}
	Q := TopSortTarjanRec(G[0])
	var i int
	for !Q.Empty() {
		u := Q.Dequeue().(*graphTest)
		if tSort[i] != u.i {
			t.Errorf("Expected %d, got %d.", tSort[i], u.i)
		}
		i++
	}
}
