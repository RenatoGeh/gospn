package learn

import (
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"github.com/RenatoGeh/gospn/test"
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
	z := make([]spn.SPN, m)
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
	}

	return &region{gmixtureId, z}
}

func createRegions(D spn.Dataset, m, r int) map[uint64]*region {
	L := make(map[uint64]*region)
	n := w * h
	var sq int
	for i := 0; i < n; i++ {
		x1 := i % w
		y1 := i / w
		for x2 := w - 1; x2 > x1; x2-- {
			for y2 := h - 1; y2 > y1; y2-- {
				if x1 == 0 && y1 == 0 && x2 == w-1 && y2 == h-1 {
					j, s := createSum(x1, y1, x2, y2)
					L[j] = s
					continue
				}
				var R *region
				dx, dy := x2-x1, y2-y1
				//if dx < r || dy < r {
				//l := encode(x1, y1, x1+r, y1+r)
				//R = L[l]
				//sys.Printf("R := %p\n", R)
				//}
				if dx <= r && dy <= r {
					R = createGauss(x1, y1, x2, y2, m, D)
				} else {
					R = createRegion(m)
				}
				k := encode(x1, y1, x2, y2)
				L[k] = R
				sq++
			}
		}
	}
	return L
}

func leftQuadrant(S *region, x1, y1, x2, y2, m, r int, L map[uint64]*region) {
	// S equiv R1
	// T equiv R2
	// R equiv R
	for x := 0; x < x1; x++ {
		if x1-x < r {
			continue
		}
		li, ri := encode(x, y1, x1, y2), encode(x, y1, x2, y2)
		T := L[li]
		R := L[ri]
		t, r, s := T.inner, R.inner, S.inner
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				pi := spn.NewProduct()
				pi.AddChild(s[i])
				pi.AddChild(t[j])
				w := 1.0 / float64(len(r))
				for l := range r {
					Z := r[l].(*spn.Sum)
					Z.AddChildW(pi, w)
				}
			}
		}
	}
}

func bottomQuadrant(S *region, x1, y1, x2, y2, m, r int, L map[uint64]*region) {
	for y := 0; y < y1; y++ {
		if y1-y < r {
			continue
		}
		T := L[encode(x1, y, x2, y1)]
		R := L[encode(x1, y, x2, y2)]
		t, r, s := T.inner, R.inner, S.inner
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				pi := spn.NewProduct()
				pi.AddChild(s[i])
				pi.AddChild(t[i])
				w := 1.0 / float64(len(r))
				for l := range r {
					Z := r[l].(*spn.Sum)
					Z.AddChildW(pi, w)
				}
			}
		}
	}
}

func PoonStructure(D spn.Dataset, m, r int) spn.SPN {
	w, h, max = sys.Width, sys.Height, sys.Max
	sys.Println("Creating regions...")
	L := createRegions(D, m, r)
	s := encode(0, 0, w-1, h-1)
	sys.Println("Joining regions...")
	for k, R := range L {
		if k == s {
			continue
		}
		x1, y1, x2, y2 := decode(k)
		leftQuadrant(R, x1, y1, x2, y2, m, r, L)
		bottomQuadrant(R, x1, y1, x2, y2, m, r, L)
	}
	S := L[s].inner[0]
	return S
}

func PoonGD(D spn.Dataset, m, r int, eta, eps float64) spn.SPN {
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
	sys.Println("Maximizing the likelihood through gradient descent...")
	S = GenerativeGD(S, eta, eps, D, nil, true)
	return S
}
