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

var width, height int

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

	lambda := spn.NewScopedCountingMultinomial(p.x+p.y*height, X, f)
	return lambda
}

// Poon is the Poon-Domingos generative SPN learning algorithm.
func Poon(k, b, w, h int, data []map[int]int) spn.SPN {
	width, height = w, h
	return genStructure(k, b, point{0, 0}, point{w, h}, data)
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
	sys.Printf("genStrucure(k=%d, b=%d, p0=%v, p1=%v, data)\n"+
		"m := %d, n := %d, q := %v\n", k, b, p0, p1, m, n, q.x, q.y)
	if n > 1 {
		for i := 1; i < n; i++ {
			q.x = int(math.Min(float64(p1.x), float64(q.x+k)))
			r := point{q.x, p0.y}
			pi := spn.NewProduct()
			a1, a2 := area(p0, q, k), area(r, p1, k)
			sys.Printf("i := %d -> n := %d\nq := %v\nr := %v\na1 := %d, a2 := %d\n", i, n, q, r, a1, a2)
			var c1 spn.SPN
			if a1 == 1 {
				sys.Println("a1 = 1. Create single leaf.")
				c1 = createLeaf(p0, q, b, data)
				sys.Printf("Created leaf from region (p0=%v, q=%v).\n", p0, q)
			} else if a1 == 2 {
				c1 = spn.NewProduct()
				s, t, u := point{p0.x + k, p0.y + k}, point{p0.x, p0.y + k}, point{p0.x + k, p0.y}
				cc1 := createLeaf(p0, s, b, data)
				sys.Printf("Created leaf cc1 from region (%v, %v).\n", p0, s)
				var cc2 spn.SPN
				if m > 1 {
					sys.Println("Special case for n > 1, a1 = 2.")
					cc2 = createLeaf(t, q, b, data)
					sys.Printf("Created leaf cc2 from region (%v, %v).\n", t, q)
				} else {
					sys.Println("Special case for n = 1, a1 = 2.")
					cc2 = createLeaf(u, q, b, data)
					sys.Printf("Created leaf cc2 from region (%v, %v).\n", u, q)
				}
				c1.AddChild(cc1)
				c1.AddChild(cc2)
			} else {
				sys.Printf("a1 > 2. Recurse: genStructure(k=%d, b=%d, p0=%v, q=%v, data)\n", k, b, p0, q)
				c1 = genStructure(k, b, p0, q, data)
				sys.Printf("End recursive call. genStrucure(k=%d, b=%d, p0=%v, p1=%v, data)\n"+
					"m := %d, n := %d, q := %v\n", k, b, p0, p1, m, n, q)
			}
			pi.AddChild(c1)
			var c2 spn.SPN
			if a2 == 1 {
				sys.Println("a2 = 1. Create single leaf.")
				c2 = createLeaf(r, p1, b, data)
				sys.Printf("Created leaf from region (r=%v, p1=%v).\n", r, p1)
			} else if a2 == 2 {
				sys.Println("a2 = 2. Create two leaves.")
				c2 = spn.NewProduct()
				s, t, u := point{r.x + k, r.y + k}, point{r.x, r.y + k}, point{r.x + k, r.y}
				cc1 := createLeaf(r, s, b, data)
				sys.Printf("Created leaf cc1 from region (%v, %v).\n", r, s)
				var cc2 spn.SPN
				if m > 1 {
					sys.Println("Special case for n > 1, a2 = 2.")
					cc2 = createLeaf(t, p1, b, data)
					sys.Printf("Created leaf cc2 from region (%v, %v).\n", t, p1)
				} else {
					sys.Println("Special case for n = 1, a2 = 2.")
					cc2 = createLeaf(u, p1, b, data)
					sys.Printf("Created leaf cc2 from region (%v, %v).\n", u, p1)
				}
				c2.AddChild(cc1)
				c2.AddChild(cc2)
			} else {
				sys.Printf("a2 > 2. Recurse: genStructure(k=%d, b=%d, r=%v, p1=%v, data)\n", k, b, r, p1)
				c2 = genStructure(k, b, r, p1, data)
				sys.Printf("End recursive call. genStrucure(k=%d, b=%d, p0=%v, p1=%v, data)\n"+
					"m := %d, n := %d, q := %v\n", k, b, p0, p1, m, n, q)
			}
			pi.AddChild(c2)
			sys.Printf("Creating new product node as child with weight w=1/n=%f.\n", 1.0/float64(n))
			S.AddChildW(pi, 1.0/float64(n))
		}
	}
	q.x, q.y = p1.x, p0.y
	// y-axis
	if m > 1 {
		for j := 1; j < m; j++ {
			q.y = int(math.Min(float64(p1.y), float64(q.y+k)))
			r := point{p0.x, q.y}
			pi := spn.NewProduct()
			a1, a2 := area(p0, q, k), area(r, p1, k)
			sys.Printf("j := %d -> m := %d\nq := %v\nr := %v\na1 := %d, a2 := %d\n", j, m, q, r, a1, a2)
			var c1 spn.SPN
			if a1 == 1 {
				sys.Println("a1 = 1. Create single leaf.")
				c1 = createLeaf(p0, q, b, data)
				sys.Printf("Created leaf from region (p0=%v, q=%v).\n", p0, q)
			} else if a1 == 2 {
				sys.Println("a1 = 2. Create two leaves.")
				c1 = spn.NewProduct()
				s, t, u := point{p0.x + k, p0.y + k}, point{p0.x, p0.y + k}, point{p0.x + k, p0.y}
				cc1 := createLeaf(p0, s, b, data)
				sys.Printf("Created leaf cc1 from region (%v, %v).\n", p0, s)
				var cc2 spn.SPN
				if n > 1 {
					sys.Println("Special case for m > 1, a1 = 2.")
					cc2 = createLeaf(u, q, b, data)
					sys.Printf("Created leaf cc2 from region (%v, %v).\n", u, q)
				} else {
					sys.Println("Special case for m = 1, a1 = 2.")
					cc2 = createLeaf(t, q, b, data)
					sys.Printf("Created leaf cc2 from region (%v, %v).\n", t, q)
				}
				c1.AddChild(cc1)
				c1.AddChild(cc2)
			} else {
				sys.Printf("a1 > 2. Recurse: genStructure(k=%d, b=%d, p0=%v, q=%v, data)\n", k, b, p0, q)
				c1 = genStructure(k, b, p0, q, data)
				sys.Printf("End recursive call. genStrucure(k=%d, b=%d, p0=%v, p1=%v, data)\n"+
					"m := %d, n := %d, q := %v\n", k, b, p0, p1, m, n, q)
			}
			pi.AddChild(c1)
			var c2 spn.SPN
			if a2 == 1 {
				sys.Println("a2 = 1. Create single leaf.")
				sys.Printf("Created leaf from region (r=%v, p1=%v).\n", r, p1)
				c2 = createLeaf(r, p1, b, data)
			} else if a2 == 2 {
				sys.Println("a2 = 2. Create two leaves.")
				c2 = spn.NewProduct()
				s, t, u := point{r.x + k, r.y + k}, point{r.x, r.y + k}, point{r.x + k, r.y}
				cc1 := createLeaf(r, s, b, data)
				sys.Printf("Created leaf cc1 from region (%v, %v).\n", r, s)
				var cc2 spn.SPN
				if n > 1 {
					sys.Println("Special case for m > 1, a2 = 2.")
					cc2 = createLeaf(u, p1, b, data)
					sys.Printf("Created leaf cc2 from region (%v, %v).\n", u, p1)
				} else {
					sys.Println("Special case for m = 1, a2 = 2.")
					cc2 = createLeaf(t, p1, b, data)
					sys.Printf("Created leaf cc2 from region (%v, %v).\n", t, p1)
				}
				c2.AddChild(cc1)
				c2.AddChild(cc2)
			} else {
				sys.Printf("a2 > 2. Recurse: genStructure(k=%d, b=%d, r=%v, p1=%v, data)\n", k, b, r, p1)
				c2 = genStructure(k, b, r, p1, data)
				sys.Printf("End recursive call. genStrucure(k=%d, b=%d, p0=%v, p1=%v, data)\n"+
					"m := %d, n := %d, q := %v\n", k, b, p0, p1, m, n, q)
			}
			pi.AddChild(c2)
			sys.Printf("Creating new product node as child with weight w=1/m=%f.\n", 1.0/float64(m))
			S.AddChildW(pi, 1.0/float64(m))
		}
	}

	return S
}
