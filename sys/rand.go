package sys

import "math/rand"

var (
	seed   int64
	Random *rand.Rand
)

func init() {
	seed = 101010101
	Random = rand.New(rand.NewSource(seed))
}

// RandSeed sets the pseudo-random generator's seed and resets the
func RefreshRandom(s int64) *rand.Rand {
	seed = s
	Random = rand.New(rand.NewSource(seed))
	return Random
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
	K := Random.Perm(n)
	D := make([][2]int, m)
	var j int
	for i := 0; i < m; i++ {
		D[j] = C[K[i]]
		j++
	}
	return D
}
