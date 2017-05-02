package io

import (
	"fmt"
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/spn"
	//"math/rand"
	"os"
)

// DrawRegions colors each SPN S regions assuming S represents an image. This is equivalent to
// coloring the feature maps or clusters of S. This function assumes the SPN is valid, of course.
// Draws l PGM images with each having k colors mapping regions, where k is the number of regional
// (sum or product) nodes in S and l is the number of layers of S.  Assumes S is an SPT
// (Sum-Product Tree).
func DrawRegions(S spn.SPN, filename string, w, h int, regionType string) {
	grid := make([]*common.Color, h*w)
	for i := range grid {
		grid[i] = common.Black
	}

	l := 0
	q := common.Queue{}
	ch := S.Ch()
	for i := range ch {
		q.Enqueue(ch[i])
	}
	pa := S

	for !q.Empty() {
		s := q.Dequeue().(spn.SPN)

		ipa := s.Pa()[0]
		if pa == ipa {
			if pa.Type() == regionType {
				tone := common.RandColor()
				sc := s.Sc()
				for i := range sc {
					grid[sc[i]] = tone
				}
			}
		} else {
			if pa.Type() == regionType {
				name := fmt.Sprintf("%s_l_%03d.pgm", filename, l)
				out, err := os.Create(name)
				if err != nil {
					fmt.Printf("Could not create file [%s].\n", name)
					out.Close()
					panic(err)
				}

				fmt.Fprintf(out, "P3\n%d %d\n255\n", w, h)

				for i := range grid {
					fmt.Fprintf(out, "%s", grid[i].String())
					if (i+1)%w == 0 {
						fmt.Fprintf(out, "\n")
					} else {
						fmt.Fprintf(out, " ")
					}
					grid[i] = common.Black
				}
				out.Close()
				l++
			}
			pa = ipa
		}

		ch = s.Ch()
		for i := range ch {
			q.Enqueue(ch[i])
		}
	}
}
