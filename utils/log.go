package utils

import (
	"math"
)

var (
	// LogZero = ln(0) = -inf
	LogZero float64
	EpsZero float64
)

func init() {
	LogZero = math.Inf(-1)
	EpsZero = 0
}

// LogSumLog is a function to compute the log of sum of logs.
func LogSumLog(v []float64, s []int) (float64, int) {
	max, imax, simax := math.Inf(-1), -1, 1
	n := len(v)
	for i := 0; i < n; i++ {
		if s[i] != 0 && v[i] > max {
			max, imax, simax = v[i], i, s[i]
		}
	}
	if imax == -1 {
		return 0.0, 0
	}
	p, r := max, 0.0
	for i := 0; i < n; i++ {
		if i != imax && s[i] != 0 {
			r += math.Exp(v[i]-p) * float64(s[i]) * float64(simax)
		}
	}
	if r < -1.0 {
		return p + math.Log(-1.0-r), -simax
	}
	return p + math.Log1p(r), simax
}

// LogSum is the log of the sum of probabilities given by
// 	sum_i p_i -> P + ln(sum_i e^(ln(p_i) - P)), where P = max_i ln(p_i)
// Returns a float64 with the resulting log operation. To convert back use utils.AntiLog.
func LogSum(p []float64) float64 {
	pi, n := math.Inf(-1), len(p)
	for i := 0; i < n; i++ {
		lp := math.Log(p[i])
		if lp > pi {
			pi = lp
		}
	}

	lg := 0.0
	for i := 0; i < n; i++ {
		lg += math.Exp(math.Log(p[i]) - pi)
	}
	lg = pi + math.Log(lg)

	return lg
}

// LogProd is the log of the product of probabilities given by
// 	prod_i p_i -> sum_i ln(p_i)
// Returns a float64 with the resulting log operation. To convert back use utils.AntiLog.
func LogProd(p []float64) float64 {
	var lg float64
	for _, v := range p {
		if v == LogZero {
			return v
		}
		lg += math.Log(v)
	}
	return lg
}

// AntiLog is the antilog with base e of l. It is equivalent to e raised to the power of l. Thus
// the following identities apply
// 	antiln(ln(k)) = k
// 	ln(antiln(k)) = k
// Returns a float64 that corresponds to the antilog of l.
func AntiLog(l float64) float64 {
	return math.Exp(l)
}

// LogSumPair is the result of the operation:
// 	ln(p1+p2)=ln(p1)+ln(1+p2/p1)
func LogSumPair(p1, p2 float64) float64 {
	return Log(p1) + Log(1.0+(p2/p1))
}

// Log is a typedef for math.Log.
func Log(p float64) float64 { return math.Log(p) }

// Inf is a typedef for math.Inf.
func Inf(sig int) float64 { return math.Inf(sig) }

// Trim removes elements from a slice of floats that have the same value as c.
func Trim(a []float64, c float64) []float64 {
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] == c {
			n := len(a) - 1
			a[n], a[i] = a[i], a[n]
			a = a[:n]
		}
	}
	return a
}

// LogSumExp takes a slice of floats a={a_1,...,a_n} and computes ln(exp(a_1)+...+exp(a_n)).
func LogSumExp(a []float64) float64 {
	a = Trim(a, LogZero)
	max := a[0]
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	if math.IsInf(max, 0) {
		return max
	}
	var l float64
	for _, v := range a {
		l += math.Exp(v - max)
	}
	return math.Log(l) + max
}

// LogSumExpPair takes two floats l and r and computes ln(l+r). Particular case of LogSumExp.
func LogSumExpPair(l, r float64) float64 {
	var max, min float64
	if l >= r {
		max, min = l, r
	} else {
		max, min = r, l
	}
	if math.IsInf(max, 0) {
		return max
	} else if math.IsInf(min, 0) {
		// When min=LogZero=-inf (i.e. exp(min)=0), ln(exp(max) + exp(min))=ln(exp(max))=max.
		return max
	}
	return math.Log(1.0+math.Exp(min-max)) + max
}
