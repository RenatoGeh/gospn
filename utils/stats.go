package utils

import (
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
