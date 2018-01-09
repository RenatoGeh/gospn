package learn

import (
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/test"
	"gonum.org/v1/gonum/floats"
	"math"
	"sort"
)

var (
	w   = sys.Width
	h   = sys.Height
	max = sys.Max
)

const (
	regionId = iota
	gmixtureId
	sumId
)

func getRegionType(t int) string {
	switch t {
	case regionId:
		return "region"
	case gmixtureId:
		return "gaussian_mix"
	case sumId:
		return "sum"
	}
	return "none"
}

type region struct {
	id    int
	inner []spn.SPN
}

// Encode is encode.
func Encode(x1, y1, x2, y2 int) uint64 {
	_w, _h := uint64(w+1), uint64(h+1)
	return ((uint64(y1)*_w+uint64(x1))*_w+uint64(x2))*_h + uint64(y2)
}

// Decode is decode.
func Decode(k uint64) (x1, y1, x2, y2 int) {
	_w, _h := uint64(w+1), uint64(h+1)
	y2 = int(k % _h)
	c := (k - uint64(y2)) / _h
	x2 = int(c % _w)
	c = (c - uint64(x2)) / _w
	x1 = int(c % _w)
	y1 = int((c - uint64(x1)) / _w)
	return
}

func createSum(x1, y1, x2, y2 int) (uint64, *region) {
	return Encode(x1, y1, x2, y2), &region{sumId, []spn.SPN{spn.NewSum()}}
}

func createRegion(m int) *region {
	z := make([]spn.SPN, m)
	for i := 0; i < m; i++ {
		z[i] = spn.NewSum()
	}
	return &region{regionId, z}
}

func partitionQuantiles(X []int, m int) [][]float64 {
	k := len(X)
	var l int
	if k/m <= 1 {
		l = 1
	} else {
		l = int(floats.Round(float64(k)/float64(m), 0))
	}
	P := make([][]float64, m)
	for i := 0; i < m; i++ {
		q := (i + 1) * l
		if i == m-1 {
			q = k
		}
		var Q []int
		for j := i * l; j < q && j < k; j++ {
			Q = append(Q, X[j])
		}
		// Compute mean and standard deviation.
		var mu, sigma float64
		n := float64(len(Q))
		for _, x := range Q {
			mu += float64(x)
		}
		mu /= n
		for _, x := range Q {
			d := float64(x) - mu
			sigma += d * d
		}
		sigma = math.Sqrt(sigma / n)
		if sigma == 0 {
			sigma = 1
		}
		P[i] = []float64{mu, sigma, n}
	}
	return P
}

func createUnitRegion(x, y, m int, D spn.Dataset) *region {
	p := x + y*w
	V := make([]int, len(D))
	for i := range D {
		V[i] = D[i][p]
	}
	sort.Ints(V)

	//sys.Printf("p: %d\n", p)
	S := make([]spn.SPN, m)
	Q := partitionQuantiles(V, m)
	for i, q := range Q {
		g := spn.NewGaussianParams(p, q[0], q[1])
		S[i] = g
	}

	return &region{gmixtureId, S}
}

func createRegions(D spn.Dataset, m, r int) map[uint64]*region {
	L := make(map[uint64]*region)

	// Coarse regions (i.e. regions that have area > r*r).
	cw, ch := w/r, h/r
	for ca := 1; ca <= cw; ca++ {
		for cb := 1; cb <= ch; cb++ {
			if ca == 1 && cb == 1 {
				continue
			}
			for x1 := 0; x1 <= w-ca*r; x1 += r {
				x2 := x1 + ca*r
				for y1 := 0; y1 <= h-cb*r; y1 += r {
					y2 := y1 + cb*r
					if ca == cw && cb == ch {
						k, R := createSum(x1, y1, x2, y2)
						L[k] = R
						continue
					}
					k := Encode(x1, y1, x2, y2)
					R := createRegion(m)
					L[k] = R
				}
			}
		}
	}

	// Fine regions (i.e. regions that have area <= r*r).
	for ca := 0; ca < cw; ca++ {
		for cb := 0; cb < ch; cb++ {
			for x := 1; x <= r; x++ {
				for y := 1; y <= r; y++ {
					for x1 := ca * r; x1 <= (ca+1)*r-x; x1++ {
						x2 := x1 + x
						for y1 := cb * r; y1 <= (cb+1)*r-y; y1++ {
							y2 := y1 + y
							k := Encode(x1, y1, x2, y2)
							var R *region
							if x == 1 && y == 1 {
								R = createUnitRegion(x1, y1, m, D)
							} else {
								R = createRegion(m)
							}
							L[k] = R
						}
					}
				}
			}
		}
	}

	return L
}

func linkRegions(R, S, T *region) {
	m := len(S.inner) * len(T.inner)
	for i := range S.inner {
		for j := range T.inner {
			pi := spn.NewProduct()
			pi.AddChild(S.inner[i])
			pi.AddChild(T.inner[j])
			for n := range R.inner {
				s := R.inner[n].(*spn.Sum)
				s.AddChildW(pi, 1.0/float64(m))
			}
		}
	}
}

func connectRegions(r int, L map[uint64]*region) spn.SPN {
	var Z spn.SPN
	cw, ch := w/r, h/r
	for ca := 1; ca <= cw; ca++ {
		for cb := 1; cb <= ch; cb++ {
			// Connects coarse regions to fine regions.
			if ca == 1 && cb == 1 {
				for x1 := 0; x1 < w; x1 += r {
					x2 := x1 + r
					for y1 := 0; y1 < h; y1 += r {
						y2 := y1 + r
						//sys.Printf("%d, %d, %d, %d\n", x1, y1, x2, y2)
						k := Encode(x1, y1, x2, y2)
						R := L[k]
						for x := x1 + 1; x < x2; x++ {
							p, q := Encode(x1, y1, x, y2), Encode(x, y1, x2, y2)
							S, T := L[p], L[q]
							linkRegions(R, S, T)
						}
						for y := y1 + 1; y < y2; y++ {
							p, q := Encode(x1, y, x2, y2), Encode(x1, y1, x2, y)
							S, T := L[p], L[q]
							linkRegions(R, S, T)
						}
					}
				}
				continue
			}
			for x1 := 0; x1 <= w-ca*r; x1 += r {
				x2 := x1 + ca*r
				for y1 := 0; y1 <= h-cb*r; y1 += r {
					y2 := y1 + cb*r
					k := Encode(x1, y1, x2, y2)
					R := L[k]
					if ca == cw && cb == ch {
						Z = R.inner[0]
					}
					//sys.Printf("R pos: (%d, %d, %d, %d)=%d, R=%v\n", x1, y1, x2, y2, k, R)
					for x := x1 + r; x < x2; x += r {
						p, q := Encode(x1, y1, x, y2), Encode(x, y1, x2, y2)
						S, T := L[p], L[q]
						//sys.Printf("p=%d=(%d, %d, %d, %d), q=%d=(%d, %d, %d, %d), S=%v, T=%v\n", p, x1, y1, x, y2, q, x, y1, x2, y2, S, T)
						linkRegions(R, S, T)
					}
					for y := y1 + r; y < y2; y += r {
						p, q := Encode(x1, y, x2, y2), Encode(x1, y1, x2, y)
						S, T := L[p], L[q]
						linkRegions(R, S, T)
					}
				}
			}
		}
	}

	for ca := 0; ca < cw; ca++ {
		for cb := 0; cb < ch; cb++ {
			for x := 1; x <= r; x++ {
				for y := 1; y <= r; y++ {
					for x1 := ca * r; x1 <= (ca+1)*r-x; x1++ {
						x2 := x1 + x
						for y1 := cb * r; y1 <= (cb+1)*r-y; y1++ {
							y2 := y1 + y
							if x == 1 && y == 1 {
								continue
							}
							k := Encode(x1, y1, x2, y2)
							R := L[k]
							for px := x1 + 1; px < x2; px++ {
								p, q := Encode(x1, y1, px, y2), Encode(px, y1, x2, y2)
								S, T := L[p], L[q]
								linkRegions(R, S, T)
							}
							for py := y1 + 1; py < y2; py++ {
								p, q := Encode(x1, py, x2, y2), Encode(x1, y1, x2, py)
								S, T := L[p], L[q]
								linkRegions(R, S, T)
							}
						}
					}
				}
			}
		}
	}
	return Z
}

func PoonStructure(D spn.Dataset, m, r int) (spn.SPN, map[uint64]*region) {
	w, h, max = sys.Width, sys.Height, sys.Max
	L := createRegions(D, m, r)
	S := connectRegions(r, L)
	return S, L
}

func cmpMarginal(px int, R *region, T spn.StorerTable) float64 {
	var t, d float64
	m := math.Inf(1)
	for _, s := range R.inner {
		l, e := T.Single(s)
		if l == math.Inf(-1) {
			continue
		}
		if m == math.Inf(1) || l > m {
			m = l
		}
		sys.Printf("m: %f, l: %f, e: %v\n", m, l, e)
	}
	for _, s := range R.inner {
		l, _ := T.Single(s)
		if l == math.Inf(-1) {
			continue
		}
		p := math.Exp(l - m)
		mu, _ := s.(*spn.Gaussian).Params()
		d += mu * p
		t += p
	}
	d /= t
	sys.Printf("m: %f, d: %f, t: %f\n", m, d, t)
	return d
}

func imageParams(I spn.VarSet) (float64, float64) {
	var mu, sigma float64
	n := float64(len(I))
	for p := range I {
		mu += float64(p)
	}
	mu /= n
	for p := range I {
		d := mu - float64(p)
		sigma += d * d
	}
	sigma = math.Sqrt(sigma / n)
	return mu, sigma
}

func pixelValue(mu, sigma, p float64) int {
	return int(p*sigma + mu)
}

func PoonCmpl(S spn.SPN, I spn.VarSet, L map[uint64]*region) spn.VarSet {
	st := spn.NewStorer()
	_, itk := spn.StoreInference(S, I, -1, st)
	_, dtk := DeriveSPN(S, st, -1, itk, nil)
	dT, _ := st.Table(dtk)
	mu, sigma := imageParams(I)
	J := make(spn.VarSet)

	for k, v := range I {
		J[k] = v
	}

	for x := 0; x < w/2; x++ {
		for y := 0; y < h; y++ {
			p := x + y*w
			k := Encode(x, y, x+1, y+1)
			pr := cmpMarginal(p, L[k], dT)
			v := pixelValue(mu, sigma, pr)
			J[p] = v
		}
	}

	return J
}

func PoonTest(D spn.Dataset, I spn.VarSet, m, r int) (spn.SPN, spn.VarSet) {
	S, L := PoonStructure(D, m, r)
	J := PoonCmpl(S, I, L)
	//sys.Printf("Complete? %v, Decomposable? %v\n", spn.Complete(S), spn.Decomposable(S))
	sys.Printf("%v\n", spn.ComputeScope(S))
	//test.DoBFS(S, func(s spn.SPN) bool {
	//sys.Printf("S: %p -> %v, Sc: %v\n", s, s, s.Sc())
	//return true
	//}, nil)
	//spn.PrintSPN(S, "test.spn")
	var sums, prods, leaves int
	test.DoBFS(S, func(s spn.SPN) bool {
		t := s.Type()
		if t == "sum" {
			sums++
		} else if t == "product" {
			prods++
			sys.Printf("Sc: %v\n", s.Sc())
		} else {
			leaves++
		}
		return true
	}, nil)
	sys.Printf("Sums: %d, Prods: %d, Leaves: %d\nTotal:%d\n", sums, prods, leaves, sums+prods+leaves)
	return S, J
}

func BindedPoonGD(m, r int, eta, eps float64) LearnFunc {
	return func(_ map[int]Variable, data spn.Dataset) spn.SPN {
		S, _ := PoonGD(data, m, r, eta, eps)
		return S
	}
}

func PoonGD(D spn.Dataset, m, r int, eta, eps float64) (spn.SPN, map[uint64]*region) {
	S, L := PoonStructure(D, m, r)
	sys.Println("Counting nodes...")
	//spn.NormalizeSPN(S)
	var sums, prods, leaves int
	test.DoBFS(S, func(s spn.SPN) bool {
		t := s.Type()
		if t == "sum" {
			sums++
		} else if t == "product" {
			prods++
		} else {
			leaves++
		}
		return true
	}, nil)
	sys.Printf("Sums: %d, Prods: %d, Leaves: %d\nTotal:%d\n", sums, prods, leaves, sums+prods+leaves)
	h := spn.ComputeHeight(S)
	sys.Printf("Height: %d\n", h)
	sys.Printf("Complete? %v, Decomposable? %v\n", spn.Complete(S), spn.Decomposable(S))
	sys.Println("Maximizing the likelihood through gradient descent...")
	spn.PrintSPN(S, "test_before.spn")
	S = GenerativeGD(S, eta, eps, D, nil, true)
	spn.PrintSPN(S, "test_after.spn")
	return S, L
}
