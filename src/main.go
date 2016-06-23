package main

import (
	"fmt"
	"github.com/RenatoGeh/gospn/src/learn"
)

func main() {
	fmt.Printf("igamma = %f\nchisquare=%f\n", learn.Igamma(2, 1), learn.Chisquare(3, 15.2))
}
