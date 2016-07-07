package utils

import (
	//"fmt"
	"math"
)

const eps = 1e-12

// Lower incomplete gamma.
func lgamma(x, s float64, regularized bool) float64 {
	if x == 0 {
		return 0
	}
	if x < 0 || s <= 0 {
		return math.NaN()
	}

	if x > 1.1 && x > s {
		if regularized {
			return 1.0 - ugamma(x, s, regularized)
		} else {
			return math.Gamma(s) - ugamma(x, s, regularized)
		}
	}

	var ft float64
	r := s
	c := 1.0
	pws := 1.0

	if regularized {
		logg, _ := math.Lgamma(s)
		ft = s*math.Log(x) - x - logg
	} else {
		ft = s*math.Log(x) - x
	}
	ft = math.Exp(ft)
	for c/pws > eps {
		r += 1
		c *= x / r
		pws += c
	}
	return pws * ft / s
}

// Upper incomplete gamma.
func ugamma(x, s float64, regularized bool) float64 {
	if x <= 1.1 || x <= s {
		if regularized {
			return 1 - lgamma(x, s, regularized)
		} else {
			return math.Gamma(s) - lgamma(x, s, regularized)
		}
	}

	f := 1.0 + x - s
	C := f
	D := 0.0
	var a, b, chg float64

	for i := 1; i < 10000; i++ {
		a = float64(i) * (s - float64(i))
		b = float64(i<<1) + 1.0 + x - s
		D = b + a*D
		C = b + a/C
		D = 1.0 / D
		chg = C * D
		f *= chg
		if math.Abs(chg-1) < eps {
			break
		}
	}
	if regularized {
		logg, _ := math.Lgamma(s)
		return math.Exp(s*math.Log(x) - x - logg - math.Log(f))
	} else {
		return math.Exp(s*math.Log(x) - x - math.Log(f))
	}
}

type ifctn func(float64) float64

func simpson38(f ifctn, a, b float64, n int) float64 {
	h := (b - a) / float64(n)
	h1 := h / 3
	sum := f(a) + f(b)
	for j := 3*n - 1; j > 0; j-- {
		if j%3 == 0 {
			sum += 2 * f(a+h1*float64(j))
		} else {
			sum += 3 * f(a+h1*float64(j))
		}
	}
	return h * sum / 8
}

func gammaIncQ(a, x float64) float64 {
	aa1 := a - 1
	var f ifctn = func(t float64) float64 {
		return math.Pow(t, aa1) * math.Exp(-t)
	}
	y := aa1
	h := 1.5e-2
	for f(y)*(x-y) > 2e-8 && y < x {
		y += .4
	}
	if y > x {
		y = x
	}
	return 1 - simpson38(f, 0, y, int(y/h/math.Gamma(a)))
}

func chisquare(dof int, distance float64) float64 {
	return gammaIncQ(.5*float64(dof), .5*distance)
}

// Lower incomplete Gamma function.
//func igamma(a, x float64) float64 {
//var sum float64 = 0.0
//var t float64 = 1.0 / a
//var n float64 = 1.0

//for t != 0 {
//sum += t
//t *= x / (a + n)
//n++
//}

//return math.Pow(x, a) * math.Exp(-x) * sum
//}

// Incomplete gamma convergence limit.
//const convgamma = 200

// Incomplete Gamma function.
//func Igamma(k, x float64) float64 {
//if x < 0.0 {
//return 0.0
//}

//s := (1.0 / x) * math.Pow(x, k) * math.Exp(-x)
//sum, nom, den := 1.0, 1.0, 1.0

//for i := 0; i < convgamma; i++ {
//nom *= x
//k++
//den *= k
//sum += nom / den
//}

//return sum * s
//}

//func igammac(a, x float64) float64 {
//if x <= 0 || a <= 0 {
//return 1.0
//} else if x < 1 || x < a {
//return 1.0 - Igamma(a, x)
//}

//lgamma, _ := math.Lgamma(a)
//ax := a*math.Log(x) - x - lgamma
//if ax < -709.78271289338399 {
//return 0.0
//}

//ax = math.Exp(ax)
//var y float64 = 1 - a
//var z float64 = x + y - 1
//c := 0.0
//p2 := 1.0
//q2 := x
//p1 := x + 1
//q1 := z * x
//ans := p1 / q1

//const eps = 0.000000000000001
//const bignum = 4503599627370496.0
//const invbignum = 2.22044604925031308085 * 0.0000000000000001

//var t float64 = -1.0
//var r float64

//for t > eps {
//c++
//y++
//z += 2
//yc := y * c
//pk := p1*z - p2*yc
//qk := q1*z - q2*yc

//if qk <= 0.0 {
//r = pk / qk
//t = math.Abs((ans - r) / r)
//ans = r
//} else {
//t = 1.0
//}

//p2 = p1
//p1 = pk
//q2 = q1
//q1 = qk

//if math.Abs(pk) > bignum {
//p2 = p2 * invbignum
//p1 = p1 * invbignum
//q2 = q2 * invbignum
//q1 = q1 * invbignum
//}
//}

//return ans * ax
//}

//func Igamma(a, x float64) float64 {
//if x <= 0 || a <= 0 {
//return 0.0
//} else if x > 1.0 && x > a {
//return 1.0 - igammac(a, x)
//}

//lgamma, _ := math.Lgamma(a)
//ax := a*math.Log(x) - x - lgamma

//if ax < -709.78271289338399 {
//return 0.0
//}

//ax = math.Exp(ax)
//var r float64 = a
//var c float64 = 1.0
//var ans float64 = 1.0

//const eps = 0.000000000000001

//for c/ans > eps {
//r++
//c = c * x / r
//ans += c
//}

//return ans * ax / a
//}

// Function chisquare returns the p-value of Pr(X^2 > cv).
// Compare this value to the significance level assumed. If chisquare < sigval, then we cannot
// accept the null hypothesis and thus the two variables are dependent.
//
// Thanks to Jacob F. W. for a tutorial on chi-square distributions.
// Source: http://www.codeproject.com/Articles/432194/How-to-Calculate-the-Chi-Squared-P-Value
func Chisquare(df int, cv float64) float64 {
	//fmt.Println("Running chi-square...")
	if cv < 0 || df < 1 {
		return 0.0
	}

	k := float64(df) * 0.5
	x := cv * 0.5

	//if df == 1 {
	//return math.Exp(-x/2.0) / (math.Sqrt2 * math.SqrtPi * math.Sqrt(x))
	//return (math.Pow(x, (k/2.0)-1.0) * math.Exp(-x/2.0)) / (math.Pow(2, k/2.0) * math.Gamma(k/2.0))
	//return lgamma(k/2.0, x/2.0, false) / math.Gamma(k/2.0)

	//} else if df == 2 {
	if df == 2 {
		return math.Exp(-x)
	}

	//fmt.Println("Computing incomplete lower gamma function...")
	pval := lgamma(x, k, false)

	if math.IsNaN(pval) || math.IsInf(pval, 0) || pval <= 1e-8 {
		return 1e-14
	}

	//fmt.Println("Computing gamma function...")
	pval /= math.Gamma(k)

	//fmt.Println("Returning chi-square value...")
	return 1.0 - pval
}

func Chisqr(df int, cv float64) float64 {
	return lgamma(float64(df)/2.0, cv/2.0, false) / math.Gamma(float64(df)/2.0)
}

/*
ChiSquareTest returns whether variable x and y are statistically independent.
We use the Chi-Square test to find correlations between the two variables.
Argument data is a table with the counting of each variable category, where the first axis is
the counting of each category of variable x and the second axis of variable y. The last element
of each row and column is the total counting. E.g.:

		+------------------------+
		|      X_1 X_2 X_3 total |
		| Y_1  100 200 100  400  |
		| Y_2   50 300  25  375  |
		|total 150 500 125  775  |
		+------------------------+

Argument p is the number of categories (or levels) in x.

Argument q is the number of categories (or levels) in y.

Returns true if independent and false otherwise.
*/
func ChiSquareTest(p, q int, data [][]int) bool {

	// df is the degree of freedom.
	//fmt.Println("Computing degrees of freedom...")
	df := (p - 1) * (q - 1)

	// Expected frequencies
	E := make([][]float64, p)
	for i := 0; i < p; i++ {
		E[i] = make([]float64, q)
	}

	//fmt.Printf("data: %v\n", data)

	//fmt.Println("Computing expected frequencies...")
	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			E[i][j] = float64(data[p][j]*data[i][q]) / float64(data[p][q])
			//fmt.Printf("E[%d][%d]: %d*%d/%d=%f\n", i, j, data[p][j], data[i][q], data[p][q], E[i][j])
		}
	}

	// Test statistic.
	//fmt.Println("Computing test statistic...")
	var chi float64 = 0
	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			diff := float64(data[i][j]) - E[i][j]
			chi += (diff * diff) / E[i][j]
		}
	}

	// Significance value.
	const sigval = 0.05

	// Compare cmd with sigval. If cmp < sigval, then dependent. Otherwise independent.
	//fmt.Println("Computing integral of p-value on chi-square distribution...")
	cmp := Chisquare(df, chi)

	//fmt.Println("Returning if integral >= significance value")
	//fmt.Printf("df: %d, chi: %f, cmp: %f\n", df, chi, cmp)
	return cmp >= sigval
}
