package main

import (
	"github.com/RenatoGeh/gospn/io"
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/sys"
)

func main() {
	sc, data := io.ParseData("data/olivetti_3bit/compiled/all.data")
	sys.Verbose = false
	S := learn.Poon(1, sc[0].Categories, 46, 56, data)
	io.DrawGraphTools("poon.py", S)
}
