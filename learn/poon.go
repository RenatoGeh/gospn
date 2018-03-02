package learn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/common"
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

	smoothWeight = 1.0

	infTk   = 0
	countTk = 1
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
	id         int
	inner      []spn.SPN
	step       int
	mapIndex   int
	maxProd    float64
	maxSum     float64
	bestDecomp []*decomp
}

type decomp struct {
	p, q uint64 // Encoded positions.
	r, s int    // Indices of the children sum nodes of this decomposition's product node.
}

func newRegion(id int, inner []spn.SPN, step int) *region {
	var d []*decomp
	if step == 1 {
		d = nil
	} else {
		d = make([]*decomp, len(inner))
	}
	return &region{id, inner, step, -1, 100, 100, d}
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

func createSum(x1, y1, x2, y2, r int) (uint64, *region) {
	return Encode(x1, y1, x2, y2), newRegion(sumId, []spn.SPN{spn.NewSum()}, r)
}

func createRegion(m, r int) *region {
	z := make([]spn.SPN, m)
	for i := 0; i < m; i++ {
		z[i] = spn.NewSum()
	}
	return newRegion(regionId, z, r)
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

func createUnitRegion(x, y, n int, D spn.Dataset) *region {
	p := x + y*w
	V := make([]int, len(D))
	for i := range D {
		V[i] = D[i][p]
	}

	if n == 1 {
		mu, sigma := utils.MuSigma(V)
		return newRegion(gmixtureId, []spn.SPN{spn.NewGaussianParams(p, mu, sigma)}, 1)
	}

	sort.Ints(V)
	//sys.Printf("p: %d\n", p)
	Q := partitionQuantiles(V, n)
	S := make([]spn.SPN, n)
	//s := spn.NewSum()
	for i, q := range Q {
		g := spn.NewGaussianParams(p, q[0], q[1])
		//s.AddChildW(g, smoothWeight/float64(n))
		S[i] = g
	}

	//return &region{gmixtureId, []spn.SPN{s}}
	return newRegion(gmixtureId, S, 1)
}

func createRegions(D spn.Dataset, m, g, r int) map[uint64]*region {
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
						k, R := createSum(x1, y1, x2, y2, r)
						L[k] = R
					} else {
						k := Encode(x1, y1, x2, y2)
						R := createRegion(m, r)
						L[k] = R
					}
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
								R = createUnitRegion(x1, y1, g, D)
							} else {
								R = createRegion(m, r)
								if x2-x1 <= 4 || y2-y1 <= 4 {
									R.step = 1
								}
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
	//m := len(S.inner) * len(T.inner)
	for i := range S.inner {
		for j := range T.inner {
			pi := spn.NewProduct()
			pi.AddChild(S.inner[i])
			pi.AddChild(T.inner[j])
			for n := range R.inner {
				s := R.inner[n].(*spn.Sum)
				s.AddChildW(pi, float64(sys.Random.Intn(10)+1))
			}
		}
	}
}

func conRegions(r int, L map[uint64]*region) spn.SPN {
	l := Encode(0, 0, w, h)
	Q := common.Queue{}
	V := make(map[uint64]bool)
	Q.Enqueue(l)
	V[l] = true

	for !Q.Empty() {
		k := Q.Dequeue().(uint64)
		x1, y1, x2, y2 := Decode(k)
		R := L[k]
		var d int
		if x2-x1 <= r && y2-y1 <= r {
			d = 1
		} else {
			d = r
		}
		//sys.Printf("k=(%d, %d, %d, %d), d=%d\n", x1, y1, x2, y2, d)
		//sys.Println("  x-axis")
		for x := x1 + d; x < x2; x += d {
			p, q := Encode(x, y1, x2, y2), Encode(x1, y1, x, y2)
			//sys.Printf("    p=(%d, %d, %d, %d), q=(%d, %d, %d, %d)\n", x, y1, x2, y2, x1, y1, x, y2)
			S, T := L[p], L[q]
			linkRegions(R, S, T)
			if !V[p] {
				Q.Enqueue(p)
				V[p] = true
			}
			if !V[q] {
				Q.Enqueue(q)
				V[q] = true
			}
		}
		//sys.Println("  y-axis")
		for y := y1 + d; y < y2; y += d {
			p, q := Encode(x1, y, x2, y2), Encode(x1, y1, x2, y)
			//sys.Printf("    p=(%d, %d, %d, %d), q=(%d, %d, %d, %d)\n", x1, y, x2, y2, x1, y1, x2, y)
			S, T := L[p], L[q]
			linkRegions(R, S, T)
			if !V[p] {
				Q.Enqueue(p)
				V[p] = true
			}
			if !V[q] {
				Q.Enqueue(q)
				V[q] = true
			}
		}
	}

	return L[l].inner[0]
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
			} else {
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

func (R *region) compMap(k uint64, m, g, r int, I spn.VarSet, L map[uint64]*region, storer *spn.Storer, existingDecomps map[string]bool, existingProds map[spn.SPN]*decomp, decompToProd map[string]spn.SPN) {
	tab, _ := storer.Table(infTk)
	counts, _ := storer.Table(countTk)
	R.mapIndex = -1
	R.maxProd = math.Inf(-1)
	R.maxSum = math.Inf(-1)

	var U []int
	for i, s := range R.inner {
		if len(s.Ch()) == 0 {
			U = append(U, i)
		}
	}
	var cUnusedIndex int
	if len(U) > 0 {
		cUnusedIndex = U[sys.Random.Intn(len(U))]
	}
	step := R.step
	var D []*decomp
	x1, y1, x2, y2 := Decode(k)
	for x := x1 + step; x < x2; x += step {
		p, q := Encode(x, y1, x2, y2), Encode(x1, y1, x, y2)
		S, T := L[p], L[q]
		s, t := S.inner[S.mapIndex], T.inner[T.mapIndex]
		var m float64
		sv, _ := tab.Single(s)
		tv, _ := tab.Single(t)
		if sv == utils.LogZero || tv == utils.LogZero {
			m = utils.LogZero
		} else {
			m = sv + tv
		}
		if len(D) == 0 || m > R.maxProd {
			D, R.maxProd = nil, m
		}
		if m == R.maxProd {
			d := &decomp{p, q, S.mapIndex, T.mapIndex}
			D = append(D, d)
		}
	}
	for y := y1 + step; y < y2; y += step {
		p, q := Encode(x1, y, x2, y2), Encode(x1, y1, x2, y)
		S, T := L[p], L[q]
		s, t := S.inner[S.mapIndex], T.inner[T.mapIndex]
		var m float64
		sv, _ := tab.Single(s)
		tv, _ := tab.Single(t)
		if sv == utils.LogZero || tv == utils.LogZero {
			m = utils.LogZero
		} else {
			m = sv + tv
		}
		if len(D) == 0 || m > R.maxProd {
			D, R.maxProd = nil, m
		}
		if m == R.maxProd {
			d := &decomp{p, q, S.mapIndex, T.mapIndex}
			D = append(D, d)
		}
	}
	//sys.Printf("(%d, %d, %d, %d), step=%d\n", x1, y1, x2, y2, R.step)
	//sys.Printf("%v\n", D)
	cDecomp := D[sys.Random.Intn(len(D))]
	//sys.Printf("Selected: %v\n", cDecomp)

	var mapSumOpts []int
	var bestDecomp []*decomp

	for i, s := range R.inner {
		if len(s.Ch()) == 0 {
			continue
		}
		spn.StoreInference(s, I, infTk, storer)
		ch := s.Ch()
		sv, _ := tab.Single(s)
		sc, _ := counts.Single(s)
		slc := math.Log(sc)
		var mS float64
		for _, c := range ch {
			cv, _ := tab.Single(c)
			v := utils.LogSumExpPair(sv+slc, cv)
			if len(bestDecomp) == 0 || v > mS {
				bestDecomp = nil
				mS = v
			}
			if mS == v {
				bestDecomp = append(bestDecomp, existingProds[c])
			}
		}
		if !existsDecomp(cDecomp, existingDecomps) {
			v := R.maxProd
			if sv != utils.LogZero {
				v = utils.LogSumExpPair(R.maxProd, sv+slc)
			}
			// Add prior if necessary here.
			if len(bestDecomp) == 0 || v > mS {
				bestDecomp, mS = nil, v
				bestDecomp = append(bestDecomp, cDecomp)
			}
		}
		tab.StoreSingle(s, mS-math.Log(sc+1))
		R.bestDecomp[i] = bestDecomp[sys.Random.Intn(len(bestDecomp))]
		if len(mapSumOpts) == 0 || sv > R.maxSum {
			mapSumOpts, R.maxSum = nil, sv
		}
		if sv == R.maxSum {
			mapSumOpts = append(mapSumOpts, i)
		}
	}
	if cUnusedIndex >= 0 {
		n := R.inner[cUnusedIndex]
		c, _ := counts.Single(n)
		v := R.maxProd - math.Log(c+1)
		tab.StoreSingle(n, v)
		R.bestDecomp[cUnusedIndex] = cDecomp
		if len(mapSumOpts) == 0 || v > R.maxSum {
			R.maxSum, mapSumOpts = v, nil
			mapSumOpts = append(mapSumOpts, cUnusedIndex)
		}
	}
	R.mapIndex = mapSumOpts[sys.Random.Intn(len(mapSumOpts))]
	bD := R.bestDecomp[R.mapIndex]
	if !existsDecomp(bD, existingDecomps) {
		pi := spn.NewProduct()
		S, T := L[bD.p], L[bD.q]
		s, t := S.inner[bD.r], T.inner[bD.s]
		pi.AddChild(s)
		pi.AddChild(t)
		//sys.Printf("(%d, %d, %d, %d)\n", x1, y1, x2, y2)
		sum := R.inner[R.mapIndex].(*spn.Sum)
		sum.AddChildW(pi, 1.0)
		storeDecomp(bD, existingDecomps)
		storeProd(pi, bD, decompToProd)
		existingProds[pi] = bD
		c, _ := counts.Single(pi)
		counts.StoreSingle(pi, c+1.0)
	} else {
		pi := extractProd(bD, decompToProd)
		c, _ := counts.Single(pi)
		counts.StoreSingle(pi, c+1.0)
	}
}

func equalsDecomp(d, e *decomp) bool {
	return d.p == e.p && d.q == e.q && d.r == e.r && d.s == e.s
}

func existsDecomp(D *decomp, E map[string]bool) bool {
	s := fmt.Sprintf("%d,%d,%d,%d", D.p, D.q, D.r, D.s)
	v, e := E[s]
	return v && e
}

func storeDecomp(D *decomp, E map[string]bool) {
	s := fmt.Sprintf("%d,%d,%d,%d", D.p, D.q, D.r, D.s)
	E[s] = true
}

func storeProd(p spn.SPN, d *decomp, E map[string]spn.SPN) {
	s := fmt.Sprintf("%d,%d,%d,%d", d.p, d.q, d.r, d.s)
	E[s] = p
}

func extractProd(d *decomp, E map[string]spn.SPN) *spn.Product {
	s := fmt.Sprintf("%d,%d,%d,%d", d.p, d.q, d.r, d.s)
	v := E[s]
	return v.(*spn.Product)
}

func compUnitRegions(I spn.VarSet, L map[uint64]*region, st *spn.Storer) {
	tab, _ := st.Table(infTk)
	for x1 := 0; x1 < w; x1++ {
		x2 := x1 + 1
		for y1 := 0; y1 < h; y1++ {
			y2 := y1 + 1
			k := Encode(x1, y1, x2, y2)
			R := L[k]
			R.mapIndex = -1
			var m float64
			for i, s := range R.inner {
				g := s.(*spn.Gaussian)
				l := g.Value(I)
				if R.mapIndex == -1 || l > m {
					R.mapIndex, m = i, l
				}
				tab.StoreSingle(g, l)
			}
		}
	}
}

func mapInference(m, g, r int, I spn.VarSet, L map[uint64]*region, st *spn.Storer, D map[string]bool, P map[spn.SPN]*decomp, Q map[string]spn.SPN) {
	compUnitRegions(I, L, st)
	cw, ch := w/r, h/r
	// Fine regions first.
	for ca := 0; ca < cw; ca++ {
		for cb := 0; cb < ch; cb++ {
			for x := 1; x <= r; x++ {
				for y := 1; y <= r; y++ {
					for x1 := ca * r; x1 <= (ca+1)*r-x; x1++ {
						x2 := x1 + x
						for y1 := cb * r; y1 <= (cb+1)*r-y; y1++ {
							if x == 1 && y == 1 {
								continue
							}
							y2 := y1 + y
							k := Encode(x1, y1, x2, y2)
							R := L[k]
							R.compMap(k, m, g, r, I, L, st, D, P, Q)
						}
					}
				}
			}
		}
	}

	for ca := 1; ca <= cw; ca++ {
		for cb := 1; cb <= ch; cb++ {
			if ca == 1 && cb == 1 {
				continue
			}
			for x1 := 0; x1 <= w-ca*r; x1 += r {
				x2 := x1 + ca*r
				for y1 := 0; y1 <= h-cb*r; y1 += r {
					y2 := y1 + cb*r
					k := Encode(x1, y1, x2, y2)
					R := L[k]
					R.compMap(k, m, g, r, I, L, st, D, P, Q)
				}
			}
		}
	}
}

func maxThroughData(D spn.Dataset, m, g, r int, L map[uint64]*region) spn.SPN {
	const batchSize = 1
	st := spn.NewStorer()
	st.NewTicket()
	st.NewTicket()
	n := len(D) / batchSize
	E, P, Q := make(map[string]bool), make(map[spn.SPN]*decomp), make(map[string]spn.SPN)
	for q := 0; q < 1; q++ {
		for i := 0; i < n; i++ {
			l, u := i*batchSize, int(math.Min(float64((i+1)*batchSize), float64(len(D))))
			sys.Printf("%d: %d, %d\n", i, l, u)
			for j := l; j < u; j++ {
				I := D[j]
				sys.Printf("Starting mapInference on instance %d\n", j)
				sys.StartTimer()
				mapInference(m, g, r, I, L, st, E, P, Q)
				sys.Printf("mapInference took %s\n", sys.StopTimer())
				sys.Printf("Finished instance %d\n", j)
			}
		}
	}
	k := Encode(0, 0, w, h)
	return L[k].inner[0]
}

func PoonStructure(D spn.Dataset, m, g, r int) spn.SPN {
	w, h, max = sys.Width, sys.Height, sys.Max
	L := createRegions(D, m, g, r)
	S := maxThroughData(D, m, g, r, L)
	//S := conRegions(r, L)
	return S
}

func PoonTest(D spn.Dataset, I spn.VarSet, m, g, r int) spn.SPN {
	S := PoonStructure(D, m, g, r)
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

func BindedPoonGD(m, g, r int, eta, eps float64) LearnFunc {
	return func(_ map[int]Variable, data spn.Dataset) spn.SPN {
		return PoonGD(data, m, g, r, eta, eps)
	}
}

func PoonGD(D spn.Dataset, m, g, r int, eta, eps float64) spn.SPN {
	S := PoonStructure(D, m, g, r)

	//spn.ComputeScope(S)
	//chs := make(map[spn.SPN]bool)
	//for _, c := range S.Ch() {
	//for _, k := range c.Ch() {
	//chs[k] = true
	//}
	//}
	//var i int
	//for c, _ := range chs {
	//sc := c.Sc()
	//sys.Printf("c[%d] has scope (size %d):\n  %v\n", i, len(sc), sc)
	//i++
	//}
	//sys.Printf("Number of children: %d\n", i)

	sys.Println("Counting nodes...")
	spn.NormalizeSPN(S)
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
	//h := spn.ComputeHeight(S)
	//sys.Printf("Height: %d\n", h)
	//sys.Printf("Complete? %v, Decomposable? %v\n", spn.Complete(S), spn.Decomposable(S))
	//sys.Println("Maximizing the likelihood through gradient descent...")
	//spn.PrintSPN(S, "test_before.spn")
	//spn.NormalizeSPN(S)
	GenerativeHardBGD(S, eta, eps, D, nil, true, 50)
	//spn.NormalizeSPN(S)
	spn.PrintSPN(S, "test_after.spn")
	return S
}
