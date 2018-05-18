package sys

import (
	"math/rand"
	"sync"
)

var (
	seed   int64
	Random *rand.Rand
	mu     sync.Mutex
)

func init() {
	seed = 101010101
	Random = rand.New(rand.NewSource(seed))
}

// RandSeed sets the pseudo-random generator's seed and resets the
func RefreshRandom(s int64) *rand.Rand {
	mu.Lock()
	seed = s
	Random = rand.New(rand.NewSource(seed))
	mu.Unlock()
	return Random
}

func RandIntn(n int) int {
	mu.Lock()
	r := Random.Intn(n)
	mu.Unlock()
	return r
}

func RandFloat64() float64 {
	mu.Lock()
	r := Random.Float64()
	mu.Unlock()
	return r
}

func RandNormFloat64() float64 {
	mu.Lock()
	r := Random.NormFloat64()
	mu.Unlock()
	return r
}

// Seed returns the pseudo-random seed.
func Seed() int64 { return seed }

func RandomComb(u int, v int, m int) [][2]int {
	n := u * v
	C := make([][2]int, n)
	var t int
	for i := 0; i < u; i++ {
		for j := 0; j < v; j++ {
			C[t][0], C[t][1] = i, j
			t++
		}
	}
	mu.Lock()
	K := Random.Perm(n)
	mu.Unlock()
	D := make([][2]int, m)
	var j int
	for i := 0; i < m; i++ {
		D[j] = C[K[i]]
		j++
	}
	return D
}
