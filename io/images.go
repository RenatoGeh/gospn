package io

import (
	"github.com/RenatoGeh/gospn/spn"
)

// SplitHalf assumes O is an image with dimensions (w, h). It then splits O in half according to
// the given CmplType. The return value of SplitHalf is then the two spn.VarSet partitions.
func SplitHalf(O spn.VarSet, t CmplType, w, h int) (spn.VarSet, spn.VarSet) {
	P, Q := make(spn.VarSet), make(spn.VarSet)
	n := w * h
	if t == Top || t == Bottom {
		for i := 0; i < n; i++ {
			if i < n/2 {
				P[i] = O[i]
			} else {
				Q[i] = O[i]
			}
		}
		if t == Bottom {
			return Q, P
		}
		return P, Q
	}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			p := x + y*w
			if x < w/2 {
				P[p] = O[p]
			} else {
				Q[p] = O[p]
			}
		}
	}
	if t == Right {
		return Q, P
	}
	return P, Q
}
