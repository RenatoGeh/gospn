package spn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/common"
	"os"
	"reflect"
)

const (
	cS = iota
	cP
	cL
)

// Print writes a text representation of SPN S to file of name filename.
func PrintSPN(S SPN, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error. Could not create file [%s].\n", filename)
		panic(err)
	}
	defer f.Close()

	L := make(map[SPN]string)
	T := []int{0, 0, 0}

	TopSortTarjanFunc(S, func(Z SPN) bool {
		t := Z.Type()
		var c rune
		var p int
		switch t {
		case "sum":
			c, p = 'S', cS
		case "product":
			c, p = 'P', cP
		default:
			c, p = 'L', cL
		}
		L[Z] = fmt.Sprintf("%c%d", c, T[p])
		T[p]++
		return true
	})

	Q := common.Queue{}
	V := make(map[SPN]bool)

	Q.Enqueue(S)
	V[S] = true

	for !Q.Empty() {
		s := Q.Dequeue().(SPN)
		fmt.Fprintf(f, "%s [\n", L[s])
		ch := s.Ch()
		var W []float64
		t := L[s][0]
		if t == 'S' {
			W = s.(*Sum).Weights()
		}
		if t == 'L' {
			st := reflect.TypeOf(s).String()
			var mu, sigma float64
			if st == "*spn.Multinomial" {
				cc := s.(*Multinomial)
				mu, sigma = cc.MuSigma()
			} else /* Gaussian */ {
				cc := s.(*Gaussian)
				mu, sigma = cc.dist.Mu, cc.dist.Sigma
			}
			//fmt.Printf("%.5f %.5f\n", mu, sigma)
			fmt.Fprintf(f, "  %.5f %.5f\n", mu, sigma)
		} else {
			for i, c := range ch {
				if t == 'S' {
					w := W[i]
					fmt.Fprintf(f, "  %s %.5f\n", L[c], w)
				} else {
					fmt.Fprintf(f, "  %s\n", L[c])
				}
				if !V[c] {
					Q.Enqueue(c)
					V[c] = true
				}
			}
		}
		fmt.Fprintf(f, "]\n")
	}
}
