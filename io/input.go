package io

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/spn"
)

func GetDataPath(dataset string) string {
	in, _ := filepath.Abs("data/" + dataset + "/compiled")
	return in
}

// GetPath gets the absolute path relative to relpath.
func GetPath(relpath string) string {
	rp, err := filepath.Abs(filepath.Clean(relpath))

	if err != nil {
		fmt.Printf("Error retrieving path \"%s\".\n", relpath)
		panic(err)
	}

	return rp
}

// ParseData reads from a file named filename and returns the scope and data map of the parsed data
// file.
func ParseData(filename string) (map[int]learn.Variable, []map[int]int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	sc := make(map[int]learn.Variable)

	scanner := bufio.NewScanner(file)

	var line string

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
		sc[varid] = learn.Variable{Varid: varid, Categories: cats}
	}

	n := len(sc)
	var data [][]int

	regex := regexp.MustCompile("[\\,\\s]+")
	// We assume complete data.
	k := 0
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
		k++
	}

	data = append(data, make([]int, n))
	s := regex.Split(line, -1)
	for i := 0; i < n; i++ {
		data[k][i], err = strconv.Atoi(s[i])
		if err != nil {
			fmt.Printf("Invalid string \"%s\" found in data file [%s].\n", s[i], filename)
			panic(err)
		}
	}

	m, n := len(data), len(data[0])
	cvntmap := make([]map[int]int, m)
	for i := 0; i < m; i++ {
		cvntmap[i] = make(map[int]int)
		for j := 0; j < n; j++ {
			cvntmap[i][j] = data[i][j]
		}
	}

	return sc, cvntmap
}

// ParseDataNL reads from a file named filename and returns the scope and data map of the parsed
// data file. This version doesn't add labels as variables, but return them separately as a slice.
func ParseDataNL(filename string) (map[int]learn.Variable, []map[int]int, []int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	sc := make(map[int]learn.Variable)

	scanner := bufio.NewScanner(file)

	var line string

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
		sc[varid] = learn.Variable{Varid: varid, Categories: cats}
	}

	n := len(sc) - 1
	var data [][]int

	delete(sc, n)
	var lbls []int

	regex := regexp.MustCompile("[\\,\\s]+")
	// We assume complete data.
	k := 0
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
		lbl, _ := strconv.Atoi(s[n])
		lbls = append(lbls, lbl)
		line = scanner.Text()
		k++
	}

	data = append(data, make([]int, n))
	s := regex.Split(line, -1)
	for i := 0; i < n; i++ {
		data[k][i], err = strconv.Atoi(s[i])
		if err != nil {
			fmt.Printf("Invalid string \"%s\" found in data file [%s].\n", s[i], filename)
			panic(err)
		}
	}
	lbl, _ := strconv.Atoi(s[n])
	lbls = append(lbls, lbl)

	m, n := len(data), len(data[0])
	cvntmap := make([]map[int]int, m)
	for i := 0; i < m; i++ {
		cvntmap[i] = make(map[int]int)
		for j := 0; j < n; j++ {
			cvntmap[i][j] = data[i][j]
		}
	}

	return sc, cvntmap, lbls
}

// ParseEvidence takes an evidence file that contains the instantiations of a subset of variables
// as evidence to be computed during inference. It may contain multiple instantiations.
//
// Returns a slice of maps, with each key corresponding to a variable ID and each associated value
// as the valuation of such variable; and the scope.
func ParseEvidence(filename string) (map[int]learn.Variable, []map[int]int, []int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	sc := make(map[int]learn.Variable)

	scanner := bufio.NewScanner(file)

	var line string

	// Get labels.
	scanner.Scan()
	line = scanner.Text()
	nslabels := 0
	fmt.Sscanf(line, "labels %d", &nslabels)
	slabels := make([]int, nslabels)
	tokens := strings.Split(line, " ")
	for i := 0; i < nslabels; i++ {
		slabels[i], err = strconv.Atoi(tokens[i+2])
		if err != nil {
			fmt.Printf("Invalid string \"%s\" found in data file [%s].\n", tokens[i], filename)
			panic(err)
		}
	}

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
		sc[varid] = learn.Variable{Varid: varid, Categories: cats}
	}

	n := len(sc)
	var data [][]int

	regex := regexp.MustCompile("[\\,\\s]+")
	k := 0
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
		k++
	}

	data = append(data, make([]int, n))
	s := regex.Split(line, -1)
	for i := 0; i < n; i++ {
		data[k][i], err = strconv.Atoi(s[i])
		if err != nil {
			fmt.Printf("Invalid string \"%s\" found in data file [%s].\n", s[i], filename)
			panic(err)
		}
	}

	m, n := len(data), len(data[0])
	cvntmap := make([]map[int]int, m)
	for i := 0; i < m; i++ {
		cvntmap[i] = make(map[int]int)
		for j := 0; j < n; j++ {
			cvntmap[i][j] = data[i][j]
		}
	}

	return sc, cvntmap, slabels
}

var glrand *rand.Rand
var glrseed int64 = -1

// ParsePartitionedData reads a data file and, with p probability, chooses ((1-p)*100)% of the data
// file to be used as evidence file. For instance, p=0.7 will create a map[int]learn.Variable,
// which contains the data variables, and two []map[int]int. The first []map[int]int returned is
// the training data, which composes 70% of the data file. The second map will return the evidence
// table with the remaining 30% data file. This partitioning is defined by the pseudo-random seed
// rseed. If rseed < 0, then use the default pseudo-random seed. It also returns the labels of each
// test line.
//
// Note: since this function "breaks" the order of classification, it returns a separate label
// containing the actual classification of each instantiation.
func ParsePartitionedData(filename string, p float64, rseed int64) (map[int]learn.Variable,
	[]map[int]int, []map[int]int, []int) {
	vartable, fdata := ParseData(filename)
	var rint func(n int) int

	if rseed < 0 {
		rint = rand.Intn
	} else {
		if glrseed < 0 {
			glrand, glrseed = rand.New(rand.NewSource(rseed)), rseed
		}
		rint = glrand.Intn
	}

	n := len(fdata)
	m := int((1 - p) * float64(n))
	test := make([]map[int]int, m)
	dels := make([]int, m)
	lbls := make([]int, m)

	// Marks instantiations that should serve as training set.
	for i := 0; i < m; i++ {
		l := rint(n)
		for fdata[l] == nil {
			l = rint(n)
		}
		test[i] = fdata[l]
		fdata[l] = nil
		dels[i] = l
	}
	sort.Ints(dels)

	// Discards marked lines from fdata.
	for i := 0; i < m; i++ {
		j := dels[i] - i
		fdata = append(fdata[:j], fdata[j+1:]...)
	}

	// All test lines have their real labels deleted and stored in a separate lbls slice.
	k := len(vartable) - 1
	for i := 0; i < m; i++ {
		lbls[i] = test[i][k]
		delete(test[i], k)
	}

	return vartable, fdata, test, lbls
}

// ReadFromFile reads an SPN from an spn mdl file.
func ReadFromFile(filename string) spn.SPN {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error. Could not create file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	return nil
}
