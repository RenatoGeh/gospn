package io

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	learn "github.com/RenatoGeh/gospn/src/learn"
)

// Reads from a file named filename and returns a matrix of
func ParseData(filename string) (map[int]learn.Variable, [][]int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	sc := make(map[int]learn.Variable)

	scanner := bufio.NewScanner(file)

	var line string = ""

	// Get variable definitions.
	for {
		if !scanner.Scan() {
			break
		}
		line = scanner.Text()
		if line[0] != 'v' {
			break
		}
		var varid, cats int
		fmt.Sscanf(line, "var %d %d", &varid, &cats)
		sc[varid] = learn.Variable{varid, cats}
	}

	n := len(sc)
	var data [][]int

	regex := regexp.MustCompile("[\\,\\s]+")
	// We assume complete data.
	for i := 0; scanner.Scan(); i++ {
		data = append(data, make([]int, n))
		s := regex.Split(line, -1)
		for j := 0; j < n; j++ {
			data[i][j], err = strconv.Atoi(s[j])
			if err != nil {
				fmt.Printf("Invalid string \"%s\" found in data file [%s].\n", s[j], filename)
				panic(err)
			}
		}
		line = scanner.Text()
	}

	return sc, data
}
