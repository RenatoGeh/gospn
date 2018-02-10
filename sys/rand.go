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
