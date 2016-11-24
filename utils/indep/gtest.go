package indep

import (
	//"fmt"
	"math"
)

// G-Test of Independence

// GTest is the G-Test log-likelihood independence test.
func GTest(p, q int, data [][]int, n int) bool {

	E := make([][]float64, p)
	for i := 0; i < p; i++ {
		E[i] = make([]float64, q)
	}

	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			E[i][j] = float64(data[p][j]*data[i][q]) / float64(data[p][q])
		}
	}

	df := (p - 1) * (q - 1)
	sum := 0.0

	for i := 0; i < p; i++ {
		//cx := float64(data[i][q])
		for j := 0; j < q; j++ {
			if E[i][j] == 0 {
				continue
			}
			o := float64(data[i][j])
			//r := (o * float64(n)) / (float64(data[p][j]) * cx)
			//if r == 0 {
			//continue
			//}
			//sum += o * math.Log(r)
			if o == 0 {
				continue
			}
			sum += o * math.Log(o/E[i][j])
		}
	}

	sum *= 2
	sigval := 0.0001 / 2
	cmp := ChiSquare(sum, df)
	//fmt.Printf("G: df: %d, g: %f, cmp: %.50f, sigval: %.50f\n", df, sum, cmp, sigval)

	return cmp >= sigval
}
