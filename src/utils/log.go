package utils

import (
	"math"
)

// LogSum is the log of the sum of probabilities given by
// 	sum_i p_i -> P + ln(sum_i e^(ln(p_i) - P)), where P = max_i ln(p_i)
// Returns a float64 with the resulting log operation. To convert back use utils.AntiLog.
func LogSum(p ...float64) float64 {
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
func LogProd(p ...float64) float64 {
	lg, n := 0.0, len(p)
	for i := 0; i < n; i++ {
		lg += math.Log(p[i])
	}
	return lg
}

// AntiLog is the antilog with base e of l. It is equivalent to e raised to the power of l. Thus
// the following identity applies
// 	antiln(ln(k)) = k
// Returns a float64 that corresponds to the antilog of l.
func AntiLog(l float64) float64 {
	return math.Exp(l)
}

// LogPairSum is the result of the operation:
// 	ln(p1+p2)=ln(p1)+ln(1+p2/p1)
func LogSumPair(p1, p2 float64) float64 {
	return Log(p1) + Log(1.0+(p2/p1))
}

// Typedef for math.Log.
func Log(p float64) float64 { return math.Log(p) }

// Typedef for math.Inf.
func Inf(sig int) float64 { return math.Inf(sig) }
