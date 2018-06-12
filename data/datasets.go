package data

import (
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

func Caltech() (map[int]learn.Variable, []map[int]int) {
	p, e := getDataset(caltech)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

func Digits() (map[int]learn.Variable, []map[int]int) {
	p, e := getDataset(digits)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

func DigitsX() (map[int]learn.Variable, []map[int]int) {
	p, e := getDataset(digitsX)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

func Olivetti() (map[int]learn.Variable, []map[int]int) {
	p, e := getDataset(olivetti)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

func OlivettiPadded() (map[int]learn.Variable, []map[int]int) {
	p, e := getDataset(olivettiPadded)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

func OlivettiBig() (map[int]learn.Variable, []map[int]int) {
	p, e := getDataset(olivettiBig)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}

func OlivettiSmall() (map[int]learn.Variable, []map[int]int) {
	p, e := getDataset(olivettiSmall)
	if e != nil {
		return nil, nil
	}
	v, d := io.ParseData(p)
	return v, d
}
