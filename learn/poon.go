package learn

import (
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
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

func Poon(k, b, w, h int, data []map[int]int) spn.SPN {
	i_vars = 0
	return genStructure(k, b, point{0, 0}, point{w - 1, h - 1}, data)
}

// This function generates a dense SPN from a dataset of images. This function is an implementation
// of the algorithm proposed in Poon and Domingos' SPN article, following the description found in
// Dennis and Ventura's article on structural learning of SPNs through clustering of variables.
// More information can be found at
// https://github.com/RenatoGeh/spn_algos/blob/master/poon/poon.pdf
func genStructure(k, b int, p0, p1 point, data []map[int]int) spn.SPN {
	n := int(math.Ceil(float64(p1.x-p0.x) / float64(k)))
	m := int(math.Ceil(float64(p1.y-p0.y) / float64(k)))
	S := spn.NewSum()
	q := point{p0.x, p1.y}
	// x-axis
	sys.Printf("genStrucure(k=%d, b=%d, p0={%d, %d}, p1={%d, %d}, data)\n"+
		"m := %d, n := %d, q := {%d, %d}\n", k, b, p0.x, p0.y, p1.x, p1.y, m, n, q.x, q.y)
	if n > 1 {
		for i := 0; i < n-1; i++ {
			q.x = int(math.Min(float64(p1.x), float64(q.x+k)))
			r := point{q.x, p0.y}
			pi := spn.NewProduct()
			a1, a2 := area(p0, q, k), area(r, p1, k)
			sys.Printf("i := %d -> n := %d\nq := {%d, %d}\nr := {%d, %d}\na1 := %d, a2 := %d\n",
				i, n, q.x, q.y, r.x, r.y, a1, a2)
			var c1 spn.SPN
			if a1 == 1 {
				sys.Println("a1 = 1. Create single leaf.")
				c1 = createLeaf(p0, q, b, data)
				sys.Printf("Created leaf from region (p0={%d, %d}, q={%d, %d}).\n", p0.x, p0.y, q.x, q.y)
			} else if a1 == 2 {
				sys.Println("a1 = 2. Create two leaves.")
				c1 = spn.NewProduct()
				cc1 := createLeaf(p0, point{p0.x + k, p0.y + k}, b, data)
				sys.Printf("Created leaf cc1 from region ({%d, %d}, {%d, %d}).\n",
					p0.x, p0.y, p0.x+k, p0.y+k)
				cc2 := createLeaf(point{p0.x + k, p0.y}, q, b, data)
				sys.Printf("Created leaf cc2 from region ({%d, %d}, {%d, %d}).\n",
					p0.x, p0.y+k, q.x, q.y)
				c1.AddChild(cc1)
				c1.AddChild(cc2)
			} else {
				sys.Printf("a1 > 2. Recurse: genStructure(k=%d, b=%d, p0={%d, %d}, q={%d, %d}, data)\n",
					k, b, p0.x, p0.y, q.x, q.y)
				c1 = genStructure(k, b, p0, q, data)
				sys.Printf("End recursive call. genStrucure(k=%d, b=%d, p0={%d, %d}, p1={%d, %d}, data)\n"+
					"m := %d, n := %d, q := {%d, %d}\n", k, b, p0.x, p0.y, p1.x, p1.y, m, n, q.x, q.y)
			}
			pi.AddChild(c1)
			var c2 spn.SPN
			if a2 == 1 {
				sys.Println("a2 = 1. Create single leaf.")
				c2 = createLeaf(r, p1, b, data)
				sys.Printf("Created leaf from region (r={%d, %d}, p1={%d, %d}).\n", r.x, r.y, p1.x, p1.y)
			} else if a2 == 2 {
				sys.Println("a2 = 2. Create two leaves.")
				c2 = spn.NewProduct()
				cc1 := createLeaf(r, point{r.x + k, r.y + k}, b, data)
				sys.Printf("Created leaf cc1 from region ({%d, %d}, {%d, %d}).\n",
					r.x, r.y, r.x+k, r.y+k)
				cc2 := createLeaf(point{r.x + k, r.y}, p1, b, data)
				sys.Printf("Created leaf cc2 from region ({%d, %d}, {%d, %d}).\n",
					r.x, r.y+k, p1.x, p1.y)
				c2.AddChild(cc1)
				c2.AddChild(cc2)
			} else {
				sys.Printf("a2 > 2. Recurse: genStructure(k=%d, b=%d, r={%d, %d}, p1={%d, %d}, data)\n",
					k, b, r.x, r.y, p1.x, p1.y)
				c2 = genStructure(k, b, r, p1, data)
				sys.Printf("End recursive call. genStrucure(k=%d, b=%d, p0={%d, %d}, p1={%d, %d}, data)\n"+
					"m := %d, n := %d, q := {%d, %d}\n", k, b, p0.x, p0.y, p1.x, p1.y, m, n, q.x, q.y)
			}
			pi.AddChild(c2)
			sys.Printf("Creating new product node as child with weight w=1/n=%f.\n", 1.0/float64(n))
			S.AddChildW(pi, 1.0/float64(n))
		}
	}
	q.x, q.y = p1.x, p0.y
	// y-axis
	if m > 1 {
		for j := 0; j < m-1; j++ {
			q.y = int(math.Min(float64(p1.y), float64(q.y+k)))
			r := point{p0.x, q.y}
			pi := spn.NewProduct()
			a1, a2 := area(p0, q, k), area(r, p1, k)
			sys.Printf("j := %d -> m := %d\nq := {%d, %d}\nr := {%d, %d}\na1 := %d, a2 := %d\n",
				j, m, q.x, q.y, r.x, r.y, a1, a2)
			var c1 spn.SPN
			if a1 == 1 {
				sys.Println("a1 = 1. Create single leaf.")
				c1 = createLeaf(p0, q, b, data)
				sys.Printf("Created leaf from region (p0={%d, %d}, q={%d, %d}).\n", p0.x, p0.y, q.x, q.y)
			} else if a1 == 2 {
				sys.Println("a1 = 2. Create two leaves.")
				c1 = spn.NewProduct()
				cc1 := createLeaf(p0, point{p0.x + k, p0.y + k}, b, data)
				sys.Printf("Created leaf cc1 from region ({%d, %d}, {%d, %d}).\n",
					p0.x, p0.y, p0.x+k, p0.y+k)
				cc2 := createLeaf(point{p0.x, p0.y + k}, q, b, data)
				sys.Printf("Created leaf cc2 from region ({%d, %d}, {%d, %d}).\n",
					p0.x, p0.y+k, q.x, q.y)
				c1.AddChild(cc1)
				c1.AddChild(cc2)
			} else {
				sys.Printf("a1 > 2. Recurse: genStructure(k=%d, b=%d, p0={%d, %d}, q={%d, %d}, data)\n",
					k, b, p0.x, p0.y, q.x, q.y)
				c1 = genStructure(k, b, p0, q, data)
				sys.Printf("End recursive call. genStrucure(k=%d, b=%d, p0={%d, %d}, p1={%d, %d}, data)\n"+
					"m := %d, n := %d, q := {%d, %d}\n", k, b, p0.x, p0.y, p1.x, p1.y, m, n, q.x, q.y)
			}
			pi.AddChild(c1)
			var c2 spn.SPN
			if a2 == 1 {
				sys.Println("a2 = 1. Create single leaf.")
				sys.Printf("Created leaf from region (r={%d, %d}, p1={%d, %d}).\n", r.x, r.y, p1.x, p1.y)
				c2 = createLeaf(r, p1, b, data)
			} else if a2 == 2 {
				sys.Println("a2 = 2. Create two leaves.")
				c2 = spn.NewProduct()
				cc1 := createLeaf(r, point{r.x + k, r.y + k}, b, data)
				sys.Printf("Created leaf cc1 from region ({%d, %d}, {%d, %d}).\n",
					r.x, r.y, r.x+k, r.y+k)
				cc2 := createLeaf(point{r.x, r.y + k}, p1, b, data)
				sys.Printf("Created leaf cc2 from region ({%d, %d}, {%d, %d}).\n",
					r.x, r.y+k, p1.x, p1.y)
				c2.AddChild(cc1)
				c2.AddChild(cc2)
			} else {
				sys.Printf("a2 > 2. Recurse: genStructure(k=%d, b=%d, r={%d, %d}, p1={%d, %d}, data)\n",
					k, b, r.x, r.y, p1.x, p1.y)
				c2 = genStructure(k, b, r, p1, data)
				sys.Printf("End recursive call. genStrucure(k=%d, b=%d, p0={%d, %d}, p1={%d, %d}, data)\n"+
					"m := %d, n := %d, q := {%d, %d}\n", k, b, p0.x, p0.y, p1.x, p1.y, m, n, q.x, q.y)
			}
			pi.AddChild(c2)
			sys.Printf("Creating new product node as child with weight w=1/m=%f.\n", 1.0/float64(m))
			S.AddChildW(pi, 1.0/float64(m))
		}
	}

	return S
}
