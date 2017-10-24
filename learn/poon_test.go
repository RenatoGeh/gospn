package learn

import (
	"fmt"
	"github.com/RenatoGeh/gospn/sys"
	"testing"
)

func cmp(x1, y1, x2, y2, p1, p2, q1, q2 uint64) bool {
	return x1 == p1 && y1 == p2 && x2 == q1 && y2 == q2
}

func subRects(w, h int) int {
	return w * (w + 1) * h * (h + 1) / 4
}

func TestPoonEncoding(t *testing.T) {
	fmt.Println("PoonCodingTest")
	w, h := sys.Width, sys.Height
	n := w * h
	var conf int
	var sq int

	u := make(map[uint64]bool)
	for i := 0; i < n; i++ {
		y, x := i/w, i%w
		for l := h; l > y; l-- {
			for j := w; j > x; j-- {
				v := Encode(x, y, j, l)
				//fmt.Printf("Encoding: (%d, %d, %d, %d) -> %d\n", x, y, j, l, v)
				_, e := u[v]
				p, q, r, s := Decode(v)
				//fmt.Printf("Decoding: %d -> (%d, %d, %d, %d)\n", v, p, q, r, s)
				if e || !cmp(uint64(x), uint64(y), uint64(j), uint64(l), uint64(p), uint64(q), uint64(r), uint64(s)) {
					fmt.Printf("Conflict! (x1, y1, x2, y2) = (%d, %d, %d, %d) = %d <=> "+
						"(%d, %d, %d, %d)\n", x, y, j, l, v, p, q, r, s)
					conf++
				}
				u[v] = true
				sq++
			}
		}
	}
	nsq := subRects(w, h)
	var str string
	if conf == 0 {
		str = fmt.Sprintf("No conflicts!")
	} else {
		str = fmt.Sprintf("Found %d conflicts!", conf)
	}
	fmt.Printf("Result: %d stored instances. Should be %d instances.\nEquals? %v\n%s\n", sq, nsq,
		nsq == sq, str)
}
