package data

import (
	"fmt"
	"github.com/RenatoGeh/gospn/io"
	"github.com/RenatoGeh/gospn/learn"
	"os"
)

const (
	// Dataset relative path.
	DatasetPath = ".cache/data/"

	// Dataset names.
	caltech        = "caltech"
	digits         = "digits"
	digitsX        = "digits_x"
	olivetti       = "olivetti"
	olivettiPadded = "olivetti_padded_u"
	olivettiBig    = "olivetti_big"
	olivettiSmall  = "olivetti_small"

	// Dataset extensions.
	DatasetExtension = ".data"

	// Dataset download upstream.
	upstreamPrepend = "https://raw.githubusercontent.com/RenatoGeh/datasets/master/"
	// Dataset download upstream append.
	upstreamAppend = "/compiled/all.data"
)

func upstreamURL(d string) string {
	return upstreamPrepend + d + upstreamAppend
}

func fullPath(d string) string {
	return DatasetPath + d + DatasetExtension
}

func exists(d string) bool {
	_, e := os.Stat(io.GetPath(d))
	return e == nil || !os.IsNotExist(e)
}

func getDataset(d string) (string, error) {
	u, p := upstreamURL(d), fullPath(d)
	if !exists(DatasetPath) {
		os.MkdirAll(DatasetPath, os.ModePerm)
	}
	e := io.DownloadFromURL(u, p, false)
	return p, e
}

// Caltech downloads a partition of the Caltech-101 dataset containing only certain categories.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and dataset indexed by variables' ID.
func Caltech() (map[int]*learn.Variable, []map[int]int) {
	p, e := getDataset(caltech)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

// Digits downloads the digits dataset containing handwritten digits from 0 to 9.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and dataset indexed by variables' ID.
func Digits() (map[int]*learn.Variable, []map[int]int) {
	p, e := getDataset(digits)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

// DigitsX downloads the digits-x dataset, an extended version of digits with more variance.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and dataset indexed by variables' ID.
func DigitsX() (map[int]*learn.Variable, []map[int]int) {
	p, e := getDataset(digitsX)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

// Olivetti downloads a downscaled Olivetti Faces dataset from Bell Labs.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and dataset indexed by variables' ID.
func Olivetti() (map[int]*learn.Variable, []map[int]int) {
	p, e := getDataset(olivetti)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

// OlivettiPadded downloads a downscaled Olivetti Faces dataset with left and right sides padded by
// uniformly distributed pixels such that both width and height are divisible by four.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and dataset indexed by variables' ID.
func OlivettiPadded() (map[int]*learn.Variable, []map[int]int) {
	p, e := getDataset(olivettiPadded)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

// OlivettiBig downloads the original Olivetti Faces dataset.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and dataset indexed by variables' ID.
func OlivettiBig() (map[int]*learn.Variable, []map[int]int) {
	p, e := getDataset(olivettiBig)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

// OlivettiSmall downloads a smaller version of Olivetti.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and dataset indexed by variables' ID.
func OlivettiSmall() (map[int]*learn.Variable, []map[int]int) {
	p, e := getDataset(olivettiSmall)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

func pullMNIST(n, b int) (map[int]*learn.Variable, []map[int]int, []map[int]int) {
	if !exists(DatasetPath) {
		os.MkdirAll(DatasetPath, os.ModePerm)
	}
	ta, ra := fmt.Sprintf("test_%d_%d.data", b, n), fmt.Sprintf("train_%d_%d.data", b, n)
	ub := upstreamPrepend + "mnist/compiled/"
	ut, ur := ub+ta, ub+ra
	pb := DatasetPath + "mnist_"
	pt, pr := pb+ta, pb+ra
	e := io.DownloadFromURL(ut, pt, false)
	if e != nil {
		return nil, nil, nil
	}
	e = io.DownloadFromURL(ur, pr, false)
	v, dt := io.ParseData(pt)
	_, dr := io.ParseData(pr)
	return v, dt, dr
}

// MNIST1000 downloads a subset of 1000 MNIST samples.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and a pair of test and training dataset indexed by variables' ID.
func MNIST1000() (map[int]*learn.Variable, []map[int]int, []map[int]int) {
	return pullMNIST(1000, 256)
}

// MNIST2000 downloads a subset of 2000 MNIST samples.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and a pair of test and training dataset indexed by variables' ID.
func MNIST2000() (map[int]*learn.Variable, []map[int]int, []map[int]int) {
	return pullMNIST(2000, 256)
}

// MNIST3Bits1000 downloads a subset of 1000 MNIST samples with 3-bit resolution.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and a pair of test and training dataset indexed by variables' ID.
func MNIST3Bits1000() (map[int]*learn.Variable, []map[int]int, []map[int]int) {
	return pullMNIST(1000, 8)
}

// MNIST3Bits2000 downloads a subset of 2000 MNIST samples with 3-bit resolution.
// For more information: https://github.com/RenatoGeh/datasets.
// Returns scope (variables) and a pair of test and training dataset indexed by variables' ID.
func MNIST3Bits2000() (map[int]*learn.Variable, []map[int]int, []map[int]int) {
	return pullMNIST(2000, 8)
}
