package utils

import (
	"gonum.org/v1/gonum/floats"
	"math"
)

// Mean returns the mean of a slice.
func Mean(c []int) float64 {
	var n int
	var m float64
	for i := range c {
		n += c[i]
	}
	for i := range c {
		m += float64(c[i]) / float64(n) * float64(i)
	}
	return m
}

// StdDev returns the standard deviation of a slice.
func StdDev(c []int) float64 {
	var n int
	var m float64
	for i := range c {
		n += c[i]
	}
	for i := range c {
		m += float64(c[i]) / float64(n) * float64(i)
	}
	var s float64
	for i := range c {
		d := float64(i) - m
		s += (float64(c[i]) / float64(n)) * d * d
	}
	return math.Sqrt(s)
}

// MuSigma returns both the mean and standard deviation of a slice.
func MuSigma(c []int) (float64, float64) {
	var n int
	var m float64
	for i := range c {
		n += c[i]
	}
	for i := range c {
		m += float64(c[i]) / float64(n) * float64(i)
	}
	var s float64
	for i := range c {
		d := float64(i) - m
		s += (float64(c[i]) / float64(n)) * d * d
	}
	return m, math.Sqrt(s)
}

// PartitionQuantiles takes a slice of values X of a single variable from the dataset and the
// number m of quantiles to partition the data. Returns a slice of pair of values containing mean
// and standard deviation (in this order).
func PartitionQuantiles(X []int, m int) [][]float64 {
	k := len(X)
	var l int
	if k/m <= 1 {
		l = 1
	} else {
		l = int(floats.Round(float64(k)/float64(m), 0))
	}
	P := make([][]float64, m)
	for i := 0; i < m; i++ {
		q := (i + 1) * l
		if i == m-1 {
			q = k
		}
		var Q []int
		for j := i * l; j < q && j < k; j++ {
			Q = append(Q, X[j])
		}
		// Compute mean and standard deviation.
		var mu, sigma float64
		n := float64(len(Q))
		for _, x := range Q {
			mu += float64(x)
		}
		mu /= n
		for _, x := range Q {
			d := float64(x) - mu
			sigma += d * d
		}
		sigma = math.Sqrt(sigma / n)
		if sigma == 0 {
			sigma = 0.5
		}
		P[i] = []float64{mu, sigma, n}
	}
	return P
}
