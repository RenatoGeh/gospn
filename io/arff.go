package io

import (
	"bufio"
	"fmt"
	"github.com/RenatoGeh/gospn/learn"
	"os"
	"strconv"
	"strings"
)

// ParseArff takes an ARFF dataset file and returns three structures.
//
// The first is a map that maps VARID -> learn.Variable, containing the internal information
// necessary for learning.
// The second is a slice of maps that correspond to the instances of the dataset. Each element in
// this slice is a map representing this instance. This map is a function VARID -> Value of the
// variable represented by VARID.
// The third is a map containing the names/labels of variables when they are of type class or
// string. It is a function VAR_CLASSID -> string, where the string is the actual label.
//
// As an example, consider the ARFF dataset below:
//
// 	% Example dataset sampling a modified rain/slippery road scenario as seen on Adnan Darwiche's
// 	% Modeling and Reasoning with Bayesian Networks (Section 4.3).
// 	% We modified variable Winter, changing it to Season and made it into a numeric (yet
// 	% categorical) variable just to showcase how we deal with numeric variables.
// 	@RELATION weather
// 	% GoSPN doesn't (yet) support continuous variables. It does accept discrete values sent as
// 	% numeric type. In this case we assume a variable season that is discrete and has 4 possible
// 	% values: 0, 1, 2, 3 with 0-3 being numeric representations for spring-winter.
// 	@ATTRIBUTE season NUMERIC
// 	% We can also use the numeric type as boolean.
// 	@ATTRIBUTE sprinkler numeric
// 	% Or just use class. In the case class is used, ParseArff returns the labels describing the
// 	% valuations in the instances.
// 	@ATTRIBUTE rain {true,false}
// 	% We can also use string. Just like class, labels are returned separately.
// 	@ATTRIBUTE wet_grass string
// 	@ATTRIBUTE slippery STRING
// 	@data
// 	0,0,true,true,false
// 	0,1,false,false,true
// 	1,0,false,false,false
// 	1,1,false,true,false
// 	1,0,true,false,true
// 	2,0,true,true,true
// 	2,0,false,false,true
// 	3,0,true,false,false
//  3,1,false,true,false
//
// For numeric variables, we take the highest value in the dataset and set this value as the
// categorical upper bound of the variable.
func ParseArff(filename string) (name string, sc map[int]learn.Variable, vals []map[int]int,
	labels map[int]map[string]int) {
	in, err := os.Open(filename)

	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", filename)
		panic(err)
	}
	defer in.Close()

	stream := bufio.NewScanner(in)
	labels = make(map[int]map[string]int)
	sc = make(map[int]learn.Variable)
	var typs []string
	var counts map[int]int
	data := false
	for i, lc := 0, 0; stream.Scan(); lc++ {
		line := stream.Text()

		// Line is a comment.
		if line[0] == '%' {
			continue
		}

		if !data {
			_l := strings.ToLower(line)
			if strings.HasPrefix(_l, "@relation") {
				// Dataset name.
				name = strings.Fields(line)[1]
			} else if strings.HasPrefix(_l, "@attribute") {
				// Attributes.
				_f := strings.Fields(line)
				n, typ := _f[1], strings.Join(_f[2:], "")

				_t := strings.ToLower(typ)
				var cat int
				if _t == "numeric" {
					// Special treatment for numerics.
					typs = append(typs, _t)
				} else if _t == "string" {
					// Special treatment for strings.
					labels[i] = make(map[string]int)
					typs = append(typs, _t)
				} else {
					// Special treatment for class.
					l := strings.FieldsFunc(typ, func(c rune) bool {
						return c == ' ' || c == ',' || c == '{' || c == '}'
					})
					labels[i] = make(map[string]int)
					for j := range l {
						labels[i][l[j]] = j
					}
					cat = len(l)
					typs = append(typs, "class")
				}
				sc[i] = learn.Variable{Varid: i, Categories: cat, Name: n}
				i++
			} else if strings.HasPrefix(_l, "@data") {
				data = true
				i = 0
				counts = make(map[int]int)
			}
		} else {
			v := strings.FieldsFunc(line, func(c rune) bool {
				return c == ' ' || c == ','
			})

			vals = append(vals, make(map[int]int))
			for j := range v {
				if typs[j] == "numeric" {
					_v, err := strconv.Atoi(v[j])
					if err != nil {
						fmt.Printf("Error parsing line %d of file [%s].\n", lc, filename)
						panic(err)
					}
					vals[i][j] = _v
					_tv := sc[j]
					if _v+1 > _tv.Categories {
						_tv.Categories = _v + 1
						sc[j] = _tv
					}
				} else if typs[j] == "string" {
					tk := v[j]
					if _, e := labels[j][tk]; !e {
						_tv := sc[j]
						_tv.Categories++
						sc[j] = _tv
						labels[j][tk] = counts[j]
						counts[j]++
					}
					vals[i][j] = labels[j][tk]
				} else /* class */ {
					tk := v[j]
					if _, e := labels[j][tk]; !e {
						labels[j][tk] = counts[j]
						counts[j]++
					}
					vals[i][j] = labels[j][tk]
				}
			}
			i++
		}
	}
	return
}
