package learn

import (
	"github.com/RenatoGeh/gospn/spn"
	"math"
)

// Point.
type point struct {
	x int
	y int
}

var i_vars int

// Computes surface area of rectangle (p, q) under the assumption the atomic unit is the rectangle
// (k, k).
func area(p, q point, k int) int {
	return int(math.Ceil(float64(p.x-q.x) * float64(p.y-q.y) / float64(k)))
}

// Creates a univariate distribution over pixels inside (p, q), where b is the number of possible
// valuations of each pixel and data is the dataset.
func createLeaf(p, q point, b int, data []map[int]int) spn.SPN {
	w, h := q.x-p.x, q.y-p.y
	n := w * h
	X := make([]int, n)

	// Selects all variables inside rectangle area (p, q). X is equivalent to this leaf's scope.
	var l int
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			k := (i + p.x) + j*w
			X[l] = k
			l++
		}
	}

	// Computes the frequency of each pixel value.
	f := make([]int, b)
	for _, v := range data {
		for _, i := range X {
			f[v[i]]++
		}
	}

	lambda := spn.NewScopedCountingMultinomial(i_vars, X, f)
	i_vars++
	return lambda
}

func Structure(k, b int, p0, p1 point, data []map[int]int) spn.SPN {
	n := int(math.Ceil(float64(p1.x-p0.x) / float64(k)))
	m := int(math.Ceil(float64(p1.y-p0.y) / float64(k)))
	S := spn.NewSum()
	q := point{p0.x, p1.y}
	// x-axis
	if n > 1 {
		for i := 0; i < n; i++ {
			q.x = int(math.Min(float64(p1.x), float64(q.x+k)))
			r := point{q.x, p0.y}
			pi := spn.NewProduct()
			a1, a2 := area(p0, q, k), area(r, p1, k)
			var c1 spn.SPN
			if a1 == 1 {
				c1 = createLeaf(p0, q, b, data)
			} else if a1 == 2 {
				c1 = spn.NewProduct()
				cc1 := createLeaf(p0, point{p0.x + k, p0.y + k}, b, data)
				cc2 := createLeaf(point{p0.x + k, p0.y}, q, b, data)
				c1.AddChild(cc1)
				c1.AddChild(cc2)
			} else {
				c1 = Structure(k, b, p0, q, data)
			}
			pi.AddChild(c1)
			var c2 spn.SPN
			if a2 == 1 {
				c2 = createLeaf(r, p1, b, data)
			} else if a2 == 2 {
				c2 = spn.NewProduct()
				cc1 := createLeaf(r, point{r.x + k, r.y + k}, b, data)
				cc2 := createLeaf(point{r.x + k, r.y}, p1, b, data)
				c2.AddChild(cc1)
				c2.AddChild(cc2)
			} else {
				c2 = Structure(k, b, r, p1, data)
			}
			pi.AddChild(c2)
			S.AddChildW(pi, 1.0/float64(n))
		}
	}
	q.x, q.y = p1.x, p0.y
	// y-axis
	if m > 1 {
		for j := 0; j < m; j++ {
			q.y = int(math.Min(float64(p1.y), float64(q.y+k)))
			r := point{p0.x, q.y}
			pi := spn.NewProduct()
			a1, a2 := area(p0, q, k), area(r, p1, k)
			var c1 spn.SPN
			if a1 == 1 {
				c1 = createLeaf(p0, q, b, data)
			} else if a1 == 2 {
				c1 = spn.NewProduct()
				cc1 := createLeaf(p0, point{p0.x + k, p0.y + k}, b, data)
				cc2 := createLeaf(point{p0.x, p0.y + k}, q, b, data)
				c1.AddChild(cc1)
				c1.AddChild(cc2)
			} else {
				c1 = Structure(k, b, p0, q, data)
			}
			pi.AddChild(c1)
			var c2 spn.SPN
			if a2 == 1 {
				c2 = createLeaf(r, p1, b, data)
			} else if a2 == 2 {
				c2 = spn.NewProduct()
				cc1 := createLeaf(r, point{r.x + k, r.y + k}, b, data)
				cc2 := createLeaf(point{r.x, r.y + k}, p1, b, data)
				c2.AddChild(cc1)
				c2.AddChild(cc2)
			} else {
				c2 = Structure(k, b, r, p1, data)
			}
			pi.AddChild(c2)
			S.AddChildW(pi, 1.0/float64(m))
		}
	}

	return S
}
