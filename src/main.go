package main

import (
	"fmt"
	"path/filepath"

	io "github.com/RenatoGeh/gospn/src/io"
	utils "github.com/RenatoGeh/gospn/src/utils"
)

func main() {
	dir, err := filepath.Abs("../data/crt")

	if err != nil {
		fmt.Printf("Could not retrieve relative path \"../data/\".\n")
		panic(err)
	}

	io.PBMToData(utils.StringConcat(dir, "/circles"),
		utils.StringConcat(dir, "/compiled/circles.data"))
	io.PBMToData(utils.StringConcat(dir, "/rectangles"),
		utils.StringConcat(dir, "/compiled/rectangles.data"))
	io.PBMToData(utils.StringConcat(dir, "/triangles"),
		utils.StringConcat(dir, "/compiled/triangles.data"))
}
