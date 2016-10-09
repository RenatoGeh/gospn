package indep

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_cdf.h>
*/
import "C"

// ChiSquare returns the cumulative distribution function at point chi, that is:
// 	Pr(X^2 <= chi)
// Where X^2 is the chi-square distribution X^2(df), with df being the degree of freedom.
func ChiSquare(chi float64, df int) float64 {
	return float64(C.gsl_cdf_chisq_P(C.double(chi), C.double(df)))
}
