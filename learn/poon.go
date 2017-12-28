package learn

import (
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/test"
	"github.com/RenatoGeh/gospn/utils"
	"math"
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

func createGMix(x1, y1, x2, y2, m, r int, D spn.Dataset) *region {
	S := spn.NewSum()
	n := r * r
	Mu, Sigma, V := make([]float64, n), make([]float64, n), make([]int, n)
	var l int
	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			p := x + y*w
			v := make([]int, max)
			for _, q := range D {
				v[q[p]]++
			}
			Mu[l], Sigma[l] = utils.MuSigma(v)
			// This is tricky. If the standard deviation is zero, then the probability is undefined,
			// which can cause problems. We alleviate this problem by setting it to 1 in such cases.
			// Since we're dealing with a discrete problem, it is fine to just set it to 1. But this
			// shouldn't be done for the continous or general case.
			if Sigma[l] == 0 {
				Sigma[l] = 1
			}
			V[l] = p
			l++
		}
	}
	w := 1.0 / float64(m)
	for i := 0; i < m; i++ {
		p := spn.NewProduct()
		for j := 0; j < n; j++ {
			z := spn.NewGaussianParams(V[j], Mu[j], Sigma[j])
			p.AddChild(z)
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
					R = createGMix(x1, y1, x2, y2, m, r, D)
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

// PoonUnwrap takes a VarSet product of MAP extraction from a Poon structure SPN and converts it
// back into a regular VarSet. The Poon structure allows for a lower resolution r to be applied to
// the image. This way, we compress the processed image, meaning the new atomic unit is a r x r
// "pixel". Each pixel is a gaussian mixture that represents the atomic unit, meaning a MAP
// extraction will yield a compressed VarSet. Applying the PoonUnwrap transformation decompresses
// the VarSet by simply applying the same value at the atomic unit P to each pixel in P.
func PoonUnwrap(set spn.VarSet, r, w int) spn.VarSet {
	dset := make(spn.VarSet)
	for k, v := range set {
		for i := 0; i < r; i++ {
			for j := 0; j < r; j++ {
				dset[k+i+j*w] = v
			}
		}
	}
	return dset
}
