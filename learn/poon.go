package learn

import (
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/test"
	"github.com/RenatoGeh/gospn/utils"
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

func createUnitRegion(x1, y1, x2, y2, m, r int, G map[int]*spn.Gaussian, D spn.Dataset) *region {
	n := r * r
	S := make([]spn.SPN, m)
	P := make([]*spn.Product, m)
	Z := make([]*spn.Sum, n)
	for i := range S {
		S[i] = spn.NewSum()
		P[i] = spn.NewProduct()
	}
	for i := range S {
		for j := range P {
			S[i].(*spn.Sum).AddChildW(P[j], 1.0/float64(m))
		}
	}
	px := make([]int, n) // pixels
	for i, x := 0, x1; x < x2; x++ {
		for y := y1; y < y2; y, i = y+1, i+1 {
			p := x + y*w
			px[i] = p
		}
	}
	for i := range Z {
		// Create each gaussian g_i that is responsible for the i-th quantile out of m total.
		p := px[i]
		var X []int
		for _, d := range D {
			X = append(X, d[p])
		}
		sort.Ints(X)
		Q := partitionQuantiles(X, m)
		k := float64(len(X))
		Z[i] = spn.NewSum()
		for j := range P {
			P[j].AddChild(Z[i])
		}
		for _, q := range Q {
			//sys.Printf("p: %d, Mu: %f, Sigma: %f, Sample size: %d\n", p, q[0], q[1], int(q[2]))
			g := spn.NewGaussianParams(p, q[0], q[1])
			Z[i].AddChildW(g, q[2]/k)
		}
	}
	return &region{gmixtureId, S}
}

func createGMix(x1, y1, x2, y2, m, r int, G map[int]*spn.Gaussian, D spn.Dataset) *region {
	S := spn.NewSum()
	n := r * r
	Z := make([]*spn.Gaussian, n)
	var l int
	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			p := x + y*w
			if g, e := G[p]; e {
				Z[l] = g
			} else {
				v := make([]int, max)
				for _, q := range D {
					v[q[p]]++
				}
				mu, sigma := utils.MuSigma(v)
				// This is tricky. If the standard deviation is zero, then the probability is undefined,
				// which can cause problems. We alleviate this problem by setting it to 1 in such cases.
				// Since we're dealing with a discrete problem, it is fine to just set it to 1. But this
				// shouldn't be done for the continous or general case.
				if sigma == 0 {
					sigma = 1
				}
				Z[l] = spn.NewGaussianParams(p, mu, sigma)
				G[p] = Z[l]
				l++
			}
		}
	}
	w := 1.0 / float64(m)
	for i := 0; i < m; i++ {
		p := spn.NewProduct()
		for j := 0; j < n; j++ {
			p.AddChild(Z[j])
		}
		S.AddChildW(p, w)
	}

	return &region{gmixtureId, []spn.SPN{S}}
}

func createGauss(x1, y1, x2, y2, m int, A map[int]*spn.Gaussian, D spn.Dataset) *region {
	p := x1 + y1*w
	z := make([]spn.SPN, m)
	for i := 0; i < m; i++ {
		z[i] = A[p]
	}
	return &region{gmixtureId, z}
}

func createAtom(x, y, r, m int, D spn.Dataset) *spn.Gaussian {
	vals := make([]int, max)
	var mu, sigma, n float64
	for i := x; i < x+r; i++ {
		for j := y; j < x+r; j++ {
			p := x + y*w
			for i := range D {
				vals[D[i][p]]++
				n++
			}
		}
	}
	for i := range vals {
		mu += float64(i) * (float64(vals[i]) / float64(n))
	}
	for i := range vals {
		dx := float64(i) - mu
		sigma += float64(vals[i]) * dx * dx
	}
	sigma = math.Sqrt(sigma / float64(n))
	return spn.NewGaussianParams(x+y*w, mu, sigma)
}

func createAtoms(m, r int, D spn.Dataset) map[int]*spn.Gaussian {
	G := make(map[int]*spn.Gaussian)
	for x := 0; x < w; x += r {
		for y := 0; y < h; y += r {
			p := x + y*w
			G[p] = createAtom(x, y, r, m, D)
		}
	}
	return G
}

func createRegions(D spn.Dataset, m, r int) map[uint64]*region {
	L := make(map[uint64]*region)
	//atoms := createAtoms(m, r, D)
	atoms := make(map[int]*spn.Gaussian)
	n := w * h
	//var sq int
	for i := 0; i < n; i += r {
		if (i/w)%r != 0 {
			i += w * (r - 1)
		}
		x1 := i % w
		y1 := i / w
		for y2 := h; y2 > y1; y2 -= r {
			for x2 := w; x2 > x1; x2 -= r {
				//sq++
				if x1 == 0 && y1 == 0 && x2 == w && y2 == h {
					j, s := createSum(x1, y1, x2, y2)
					L[j] = s
					continue
				}
				var R *region
				dx, dy := x2-x1, y2-y1
				if dx < r || dy < r {
					//x := int(math.Max(float64(x1+r), float64(x2)))
					//y := int(math.Max(float64(y1+r), float64(y2)))
					//l := Encode(x1, y1, x, y)
					//R = L[l]
					continue
				} else if dx == r && dy == r {
					//R = createGMix(x1, y1, x2, y2, m, r, atoms, D)
					R = createUnitRegion(x1, y1, x2, y2, m, r, atoms, D)
					//R = createGauss(x1, y1, x2, y2, m, atoms, D)
				} else {
					R = createRegion(m)
				}
				k := Encode(x1, y1, x2, y2)
				L[k] = R
				//sys.Printf("(%d, %d, %d, %d)\n", x1, y1, x2, y2)
			}
		}
	}
	//sys.Printf("sq=%d\n", sq)
	return L
}

func leftQuadrant(S *region, x1, y1, x2, y2, m, rs int, L map[uint64]*region) {
	// S equiv R1
	// T equiv R2
	// R equiv R
	//sys.Printf("(%d, %d, %d, %d), S=%p\n", x1, y1, x2, y2, S)
	for x := 0; x < x1; x += rs {
		li, ri := Encode(x, y1, x1, y2), Encode(x, y1, x2, y2)
		T := L[li]
		R := L[ri]
		//sys.Printf("T=%p, S=%p, R=%p\n", T, S, R)
		//sys.Printf("T(%s), S(%s), R(%s)\n", getRegionType(T.id), getRegionType(S.id), getRegionType(R.id))
		t, r, s := T.inner, R.inner, S.inner
		for i := range s {
			for j := range t {
				pi := spn.NewProduct()
				pi.AddChild(s[i])
				pi.AddChild(t[j])
				//w := 1.0 / float64(len(r)*len(s))
				for l := range r {
					Z := r[l].(*spn.Sum)
					Z.AddChildW(pi, 1.0)
				}
			}
		}
	}
}

func topQuadrant(S *region, x1, y1, x2, y2, m, rs int, L map[uint64]*region) {
	//sys.Printf("(%d, %d, %d, %d), S=%p\n", x1, y1, x2, y2, S)
	for y := 0; y < y1; y += rs {
		T := L[Encode(x1, y, x2, y1)]
		R := L[Encode(x1, y, x2, y2)]
		//sys.Printf("T=%p, S=%p, R=%p\n", T, S, R)
		t, r, s := T.inner, R.inner, S.inner
		for i := range s {
			for j := range t {
				pi := spn.NewProduct()
				pi.AddChild(s[i])
				pi.AddChild(t[j])
				//w := 1.0 / float64(len(r)*len(s))
				for l := range r {
					Z := r[l].(*spn.Sum)
					Z.AddChildW(pi, 1.0)
				}
			}
		}
	}
}

func PoonStructure(D spn.Dataset, m, r int) spn.SPN {
	w, h, max = sys.Width, sys.Height, sys.Max
	sys.Println("Creating regions...")
	L := createRegions(D, m, r)
	s := Encode(0, 0, w, h)
	sys.Println("Joining regions...")
	for k, R := range L {
		if k == s {
			continue
		}
		x1, y1, x2, y2 := Decode(k)
		leftQuadrant(R, x1, y1, x2, y2, m, r, L)
		topQuadrant(R, x1, y1, x2, y2, m, r, L)
	}
	S := L[s].inner[0]
	return S
}

func PoonTest(D spn.Dataset, m, r int) spn.SPN {
	S := PoonStructure(D, m, r)
	sys.Println("Counting nodes...")
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
	return S
}

func BindedPoonGD(m, r int, eta, eps float64) LearnFunc {
	return func(_ map[int]Variable, data spn.Dataset) spn.SPN {
		return PoonGD(data, m, r, eta, eps)
	}
}

func PoonGD(D spn.Dataset, m, r int, eta, eps float64) spn.SPN {
	S := PoonStructure(D, m, r)
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
	//spn.PrintSPN(S, "test_before.spn")
	S = GenerativeBGD(S, eta, eps, D, nil, true, 100)
	//spn.PrintSPN(S, "test_after.spn")
	return S
}
