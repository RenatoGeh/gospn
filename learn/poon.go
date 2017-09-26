package learn

import (
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
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

type region struct {
	id    int
	inner []spn.SPN
}

func encode(x1, y1, x2, y2 int) uint64 {
	return ((uint64(y1)*uint64(w)+uint64(x1))*uint64(w)+uint64(x2))*uint64(h) + uint64(y2)
}

func decode(k uint64) (x1, y1, x2, y2 int) {
	_w, _h := uint64(w), uint64(h)
	y2 = int(k % _h)
	c := (k - uint64(y2)) / _h
	x2 = int(c % _w)
	c = (c - uint64(x2)) / _w
	x1 = int(c % _w)
	y1 = int((c - uint64(x1)) / _w)
	return
}

func createSum(x1, y1, x2, y2 int) (uint64, *region) {
	return encode(x1, y1, x2, y2), &region{sumId, []spn.SPN{spn.NewSum()}}
}

func createRegion(m int) *region {
	z := make([]spn.SPN, m)
	for i := 0; i < m; i++ {
		z[i] = spn.NewSum()
	}
	return &region{regionId, z}
}

func createGauss(x1, y1, x2, y2, m int, D spn.Dataset) *region {
	S := spn.NewSum()
	z := make([]*spn.Gaussian, m)
	vals := make([][]int, m)
	for i := range vals {
		vals[i] = make([]int, max)
	}

	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			p := x + y*w
			// Partition pixel p dataset into m value slices
			for i := range D {
				k := i % m
				vals[k][D[i][p]]++
			}
		}
	}

	for i := 0; i < m; i++ {
		z[i] = spn.NewGaussian(x1+y1*w, vals[i])
		S.AddChildW(z[i], 1.0/float64(m))
	}

	return &region{gmixtureId, []spn.SPN{S}}
}

func createRegions(D spn.Dataset, m, r int) map[uint64]*region {
	L := make(map[uint64]*region)
	n := w * h
	for i := 0; i < n; i++ {
		x1 := i % w
		y1 := i / w
		for x2 := w - 1; x2 >= x1; x2-- {
			for y2 := h - 1; y2 >= y1; y2-- {
				if x1 == 0 && y1 == 0 && x2 == w-1 && y2 == h-1 {
					j, s := createSum(x1, y1, x2, y2)
					L[j] = s
					continue
				}
				var R *region
				dx, dy := x2-x1, y2-y1
				if dx < r || dy < r {
					x := int(math.Max(float64(x1+r), float64(x2)))
					y := int(math.Max(float64(y1+r), float64(y2)))
					l := encode(x1, y1, x, y)
					R = L[l]
				} else if dx == r && dy == r {
					R = createGauss(x1, y1, x2, y2, m, D)
				} else {
					R = createRegion(m)
				}
				k := encode(x1, y1, x2, y2)
				L[k] = R
			}
		}
	}
	return L
}

func leftQuadrant(S *region, x1, y1, x2, y2, m int, L map[uint64]*region) {
	// S equiv R1
	// T equiv R2
	// R equiv R
	for x := 0; x < x1; x++ {
		T := L[encode(x, y1, x1, y2)]
		R := L[encode(x, y1, x2, y2)]
		t, r, s := T.inner, R.inner, S.inner
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				pi := spn.NewProduct()
				pi.AddChild(s[i])
				pi.AddChild(t[j])
				for l := range r {
					r[l].AddChild(pi)
				}
			}
		}
	}
}

func bottomQuadrant(S *region, x1, y1, x2, y2, m int, L map[uint64]*region) {
	for y := 0; y < y1; y++ {
		T := L[encode(x1, y, x2, y1)]
		R := L[encode(x1, y, x2, y2)]
		t, r, s := T.inner, R.inner, S.inner
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				pi := spn.NewProduct()
				pi.AddChild(s[i])
				pi.AddChild(t[i])
				for l := range r {
					r[l].AddChild(pi)
				}
			}
		}
	}
}

func PoonStructure(w, h int, D spn.Dataset, m, r int) spn.SPN {
	L := createRegions(D, m, r)
	s := encode(0, 0, w-1, h-1)
	for k, R := range L {
		if k == s {
			continue
		}
		x1, y1, x2, y2 := decode(k)
		leftQuadrant(R, x1, y1, x2, y2, m, L)
		bottomQuadrant(R, x1, y1, x2, y2, m, L)
	}
	S := L[s].inner[0]
	return S
}

func PoonGD(w, h int, D spn.Dataset, m, r int, eta, eps float64) spn.SPN {
	S := PoonStructure(w, h, D, m, r)
	S = GenerativeGD(S, eta, eps, D, nil, true)
	return S
}
