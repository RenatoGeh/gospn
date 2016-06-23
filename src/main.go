package main

import (
	"fmt"
	"github.com/RenatoGeh/gospn/src/learn"
)

func main() {
	fmt.Printf("igamma = %f\nchisquare=%f\n", learn.Igamma(1, 2), learn.Chisquare(3, 23.1))
}
