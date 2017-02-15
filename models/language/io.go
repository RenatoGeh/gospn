package language

import (
	"bufio"
	"fmt"
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/io"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/utils"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// T-rex: the king of regexers.
//const rex = "(\\w+|[\\.]+|[\\,\\!\\@\\#\\$\\%\\^\\&\\*\\(\\)\\;\\\\\\/\\|\\<\\>\\\"\\'\\:\\`" +
//"\\=\\-\\?\\+\\{\\}\\[\\]])"
const rex = "(\\w+)"

// Compile takes a plain text filename tfile and compiles it into a vocabulary file vfile. We
// treat punctuation as words and letters with accent marks as different characters (é != e). A
// vocabulary file contains K lines of word mapping where, for each line a number (which signals
// the id of a word) is followed by the word in question. Next we have a series of numbers that
// represent the id of each word in the order they appear in tfile.
func Compile(tfile, vfile string) {
	text, err := os.Open(io.GetPath(tfile))
	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", tfile)
		panic(err)
	}
	defer text.Close()

	vocab := make(map[string]int)
	vc := 0
	match := regexp.MustCompile(rex)

	var block []string
	cblock := 0
	nwords := 0

	// Read contents and store them into vocab and block.
	in := bufio.NewScanner(text)
	for in.Scan() {
		//fmt.Printf("Text: \"%s\"\nMatches:\n", in.Text())
		v := match.FindAllString(in.Text(), -1)
		nv := len(v)

		//for i := 0; i < nv; i++ {
		//fmt.Printf(" <%s>", v[i])
		//}
		//fmt.Printf("\n")

		if nv == 0 {
			continue
		}
		block = append(block, "")
		for i := 0; i < nv; i++ {
			str := strings.ToLower(v[i])
			_, ok := vocab[str]
			if !ok {
				vocab[str] = vc
				vc++
			}
			//fmt.Printf("%s -> %d\n", str, vocab[str])
			block[cblock] = utils.StringConcat(block[cblock], strconv.Itoa(vocab[str]))
			if i < nv-1 {
				block[cblock] = utils.StringConcat(block[cblock], " ")
			}
			nwords++
		}
		cblock++
	}

	if err := in.Err(); err != nil {
		fmt.Printf("Error parsing file [%s].\n", tfile)
		panic(err)
	}

	// Write contents into vfile.
	vocf, err := os.Create(io.GetPath(vfile))

	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", vfile)
		panic(err)
	}
	defer vocf.Close()

	// Number of vocabulary entries.
	fmt.Fprintf(vocf, "%d\n", len(vocab))
	for k, v := range vocab {
		// Write each entry as a pair (id, word).
		fmt.Fprintf(vocf, "%d %s\n", v, k)
	}
	// Number of words in block.
	fmt.Fprintf(vocf, "%d\n", nwords)
	for i := 0; i < cblock; i++ {
		// Write all lines as a list of ids.
		fmt.Fprintln(vocf, block[i])
	}
}

// Vocabulary is the in-memory representation of a .voc file.
type Vocabulary struct {
	// Entry slice : id -> word.
	entries []string
	// Number of previous words as evidence.
	n int
	// Translated block of text.
	block []int
	// Stream position indicator inside block.
	ptr int
	// Number of possible combinations.
	m int
}

// NewVocabulary constructs a new Vocabulary pointer.
func NewVocabulary(entries []string, block []int) *Vocabulary {
	return &Vocabulary{entries: entries, block: block, m: -1}
}

// Entries returns the entry map.
func (v *Vocabulary) Entries() []string { return v.entries }

// Translate returns the word corresponding to the given id.
func (v *Vocabulary) Translate(id int) string { return v.entries[id] }

// Size is the number of entries in this vocabulary.
func (v *Vocabulary) Size() int { return len(v.entries) }

// Combinations returns the number of possible combinations for Next.
func (v *Vocabulary) Combinations() int {
	if v.m >= 0 {
		return v.m
	}
	v.m = len(v.block) - v.n
	return v.m
}

// Set sets the number of previous words used as evidence and resets the ptr.
func (v *Vocabulary) Set(N int) {
	v.n = N
	v.ptr = N
}

// Next returns the next spn.VarSet of N+1 words to be used for training.
func (v *Vocabulary) Next() spn.VarSet {
	vs := make(spn.VarSet)
	vs[0] = v.block[v.ptr]
	for i := 1; i <= v.n; i++ {
		vs[v.n-i+1] = v.block[v.ptr-i]
	}
	v.ptr++
	return vs
}

// Rand returns a random word and its id amongst entries from this vocabulary.
func (v *Vocabulary) Rand() (string, int) {
	i := rand.Intn(len(v.entries))
	return v.entries[i], i
}

// Parse reads a vocabulary file vfile and returns an in-memory representation of it (Vocabulary).
func Parse(vfile string) *Vocabulary {
	voc, err := os.Open(io.GetPath(vfile))

	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", vfile)
		panic(err)
	}
	defer voc.Close()

	var n int
	fmt.Fscanf(voc, "%d", &n)

	entries := make([]string, n)
	for i := 0; i < n; i++ {
		var j int
		var str string
		fmt.Fscanf(voc, "%d %s ", &j, &str)
		entries[j] = str
	}

	var m int
	fmt.Fscanf(voc, "%d", &m)

	l, block := 0, make([]int, m)
	for i := 0; i < m; i++ {
		var k int
		fmt.Fscanf(voc, "%d", &k)
		block[l] = k
		l++
	}

	return NewVocabulary(entries, block)
}

// Write writes an SPN according to LMSPN to a .mdl file.
func Write(filename string, S spn.SPN, K, D, N int) {
	out, err := os.Create(io.GetPath(filename))

	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", filename)
		panic(err)
	}
	defer out.Close()

	fmt.Fprintf(out, "%d %d %d\n", K, D, N)

	// Root node and S_i product nodes.
	fmt.Fprintf(out, "# Weights going from root node to S layer.\n")
	q := common.Queue{}
	ch := S.Ch()
	root := S.(*spn.Sum)
	w := root.Weights()
	for i := 0; i < K; i++ {
		fmt.Fprintf(out, "%.15f ", w[i])
		q.Enqueue(ch[i])
	}
	fmt.Fprintf(out, "\n")

	// Discarting S_i nodes and retrieving B_i sum nodes.
	fmt.Fprintf(out, "# Weights going from B layer to G and M layers.\n")
	for i := 0; i < K; i++ {
		s := q.Dequeue().(spn.SPN)
		b := s.Ch()[0].(*spn.Sum)
		w, ch = b.Weights(), b.Ch()
		m := len(w)
		for j := 0; j < m; j++ {
			fmt.Fprintf(out, "%.15f ", w[j])
			q.Enqueue(ch[j])
		}
		fmt.Fprintf(out, "\n")
	}

	// Discarting G_i product nodes and retrieving M_i sum nodes.
	fmt.Fprintf(out, "# Weights going from M layer to H layer.\n")
	for i := 0; i < K; i++ {
		c1, c2 := q.Dequeue().(spn.SPN), q.Dequeue().(spn.SPN)
		var m *spn.Sum
		if c1.Type() == "sum" {
			m = c1.(*spn.Sum)
		} else if c2.Type() == "sum" {
			m = c2.(*spn.Sum)
		} else {
			fmt.Printf("This should never happen. HALP\n")
		}
		w, ch = m.Weights(), m.Ch()
		l := len(w)
		for j := 0; j < l; j++ {
			fmt.Fprintf(out, "%.15f ", w[j])
			q.Enqueue(ch[j])
		}
		fmt.Fprintf(out, "\n")
	}
	fmt.Fprintf(out, "\n# Weights going from H layer to feature vectors.\n")

	// H_i layer.
	for i := 0; i < N; i++ {
		for j := 0; j < D; j++ {
			h := q.Dequeue().(*SumVector)
			w := h.w
			m := h.n
			for l := 0; l < m; l++ {
				fmt.Fprintf(out, "%.15f ", w[l])
			}
			fmt.Fprintf(out, "\n")
		}
		fmt.Fprintf(out, "\n")
	}
}

// DrawGraphTools creates a file filename and draws an SPN spn in graph-tools. The resulting file
// is a python source code that outputs a PNG image of the graph.
func DrawGraphTools(filename string, s spn.SPN) {
	file, err := os.Create(filename)

	if err != nil {
		fmt.Printf("Error. Could not create file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	outname := utils.StringConcat(filename[0:len(filename)-len(filepath.Ext(filename))], ".png")

	fmt.Fprintf(file, "from graph_tool.all import *\n\n")
	fmt.Fprintf(file, "g = Graph(directed=True)\n")
	fmt.Fprintf(file, "vcolors = g.new_vertex_property(\"string\")\n")
	fmt.Fprintf(file, "vnames = g.new_vertex_property(\"string\")\n")
	fmt.Fprintf(file, "enames = g.new_edge_property(\"string\")\n\n")
	fmt.Fprintf(file, "def add_node(name, type):\n\tv=g.add_vertex()\n\tvnames[v]=name\n\t"+
		"vcolors[v]=type\n\treturn v\n\n")
	fmt.Fprintf(file, "def add_edge(o, t, name):\n\te=g.add_edge(o, t)\n\tenames[e]=name\n\treturn e\n\n")
	fmt.Fprintf(file, "def add_edge_nameless(o, t):\n\te=g.add_edge(o, t)\n\treturn e\n\n\n")

	read := make(map[spn.SPN]string)

	// Else, BFS the SPN and write nodes to filename.
	nvars, nsums, nprods, nsv := 0, 0, 0, 0
	queue := common.Queue{}
	queue.Enqueue(&io.BFSPair{Spn: s, Pname: "", Weight: -1.0})
	for !queue.Empty() {
		currpair := queue.Dequeue().(*io.BFSPair)
		curr, pname, pw := currpair.Spn, currpair.Pname, currpair.Weight
		ch := curr.Ch()
		nch := len(ch)

		name := "N"
		currt := curr.Type()

		// In case it is a sum node. Else product node.
		if currt == "sum" {
			if ename, exists := read[curr]; exists {
				if pw >= 0 {
					fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", pname, ename, pw)
				} else {
					fmt.Fprintf(file, "add_edge_nameless(%s, %s)\n", pname, ename)
				}
				continue
			}
			name = fmt.Sprintf("S%d", nsums)
			read[curr] = name
			fmt.Fprintf(file, "%s = add_node(\"+\", \"#ff3300\")\n", name)
			nsums++
		} else if currt == "product" {
			if ename, exists := read[curr]; exists {
				if pw >= 0 {
					fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", pname, ename, pw)
				} else {
					fmt.Fprintf(file, "add_edge_nameless(%s, %s)\n", pname, ename)
				}
				continue
			}
			name = fmt.Sprintf("P%d", nprods)
			read[curr] = name
			fmt.Fprintf(file, "%s = add_node(\"*\", \"#669900\")\n", name)
			nprods++
		} else if currt == "sum_vector" {
			if ename, exists := read[curr]; exists {
				if pw >= 0 {
					fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", pname, ename, pw)
				} else {
					fmt.Fprintf(file, "add_edge_nameless(%s, %s)\n", pname, ename)
				}
				continue
			}
			name = fmt.Sprintf("SV%d", nsv)
			read[curr] = name
			fmt.Fprintf(file, "%s = add_node(\"+\", \"#f48942\")\n", name)
			nsv++
		}

		// If pname is empty, then it is the root node. Else, link parent node to current node.
		if pname != "" {
			if pw >= 0 {
				fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", pname, name, pw)
			} else {
				fmt.Fprintf(file, "add_edge_nameless(%s, %s)\n", pname, name)
			}
		}

		var w []float64
		if curr.Type() == "sum" {
			w = (curr.(*spn.Sum).Weights())
		} else if curr.Type() == "sum_vector" {
			w = (curr.(*SumVector).w)
		}
		// For each chil, run the BFS.
		for i := 0; i < nch; i++ {
			c := ch[i]

			// If leaf, then simply write to the graphviz dot file. Else, recurse the BFS.
			if c.Type() == "leaf" {
				if ename, exists := read[c]; exists {
					if currt == "sum" {
						fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", name, ename, w[i])
					} else if currt == "sum_vector" {
						for j := 0; j < len(w); j++ {
							fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", name, ename, w[j])
						}
					} else {
						fmt.Fprintf(file, "add_edge_nameless(%s, %s)\n", name, ename)
					}
					continue
				}
				cname := fmt.Sprintf("X%d", nvars)
				read[c] = cname
				fmt.Fprintf(file, "%s = add_node(\"X_%d\", \"#0066ff\")\n", cname, c.Sc()[0])
				nvars++
				if currt == "sum" {
					fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", name, cname, w[i])
				} else if currt == "sum_vector" {
					for j := 0; j < len(w); j++ {
						fmt.Fprintf(file, "add_edge(%s, %s, \"%.3f\")\n", name, cname, w[j])
					}
				} else {
					fmt.Fprintf(file, "add_edge_nameless(%s, %s)\n", name, cname)
				}
			} else {
				tw := -1.0
				if w != nil {
					tw = w[i]
				}
				queue.Enqueue(&io.BFSPair{Spn: c, Pname: name, Weight: tw})
			}
		}
	}

	fmt.Fprintf(file, "g.vertex_properties[\"name\"]=vnames\n")
	fmt.Fprintf(file, "g.vertex_properties[\"color\"]=vcolors\n")
	//fmt.Fprintf(file, "\ngraph_draw(g, vertex_text=g.vertex_properties[\"name\"], "+
	//"edge_text=enames, vertex_fill_color=g.vertex_properties[\"color\"], "+
	//"output_size=[16384, 16384], output=\"%s\", bg_color=[1, 1, 1, 1])\n", outname)
	fmt.Fprintf(file, "\ngraph_draw(g, "+
		"edge_text=enames, vertex_fill_color=g.vertex_properties[\"color\"], "+
		"output_size=[16384, 16384], output=\"%s\", bg_color=[1, 1, 1, 1])\n", outname)
}

// Read reads an SPN from a .mdl file according to LMSPN.
func Read(filename string) (int, int, int, spn.SPN) {
	in, err := os.Open(io.GetPath(filename))

	if err != nil {
		fmt.Printf("Could not open file [%s].\n", filename)
		panic(err)
	}
	defer in.Close()

	var K, D, N int
	fmt.Fscanf(in, "%d %d %d", &K, &D, &N)

	// Root node and S_i product nodes.
	R := spn.NewSum()
	R.AutoNormalize(true)

	S := make([]*ProductIndicator, K)
	for i := 0; i < K; i++ {
		S[i] = NewProductIndicator(i)
	}

	for i := 0; i < K; i++ {
		var w float64
		fmt.Fscanf(in, "%f", &w)
		R.AddChildW(S[i], w)
	}

	// B layer.
	B := make([]*spn.Sum, K)
	for i := 0; i < K; i++ {
		B[i] = spn.NewSum()
		B[i].AutoNormalize(true)
	}

	for i := 0; i < K; i++ {
		S[i].AddChild(B[i])
	}

	// Adding M_i and G_i to B_i node.
	M, G := make([]*spn.Sum, K), make([]*SquareProduct, K)
	for i := 0; i < K; i++ {
		G[i] = NewSquareProduct()
		M[i] = spn.NewSum()
		var w1, w2 float64
		fmt.Fscanf(in, "%f %f", &w1, &w2)
		B[i].AddChildW(M[i], w1)
		B[i].AddChildW(G[i], w2)
		M[i].AutoNormalize(true)
	}

	// Add only weights for M_i sum node, since we depend on the H layer weights.
	T := N * D
	for i := 0; i < K; i++ {
		for j := 0; j < T; j++ {
			var w float64
			fmt.Fscanf(in, "%f", &w)
			M[i].AddWeight(w)
		}
	}

	// Get weights for H layer nodes.
	wmatrix := make([][]float64, D)
	for i := 0; i < D; i++ {
		wmatrix[i] = make([]float64, K)
		for j := 0; j < K; j++ {
			fmt.Fscanf(in, "%f", &wmatrix[i][j])
		}
	}

	V := make([]*Vector, N)
	for i := 0; i < N; i++ {
		V[i] = NewVector(i + 1)
	}

	// Create H nodes and add Vectors.
	H := make([][]*SumVector, N)
	for i := 0; i < N; i++ {
		H[i] = make([]*SumVector, D)
		for j := 0; j < D; j++ {
			H[i][j] = NewSumVector(wmatrix[j])
			H[i][j].AddChild(V[i])
		}
	}

	// Finally add the H nodes as children of the M layer.
	for i := 0; i < K; i++ {
		for j := 0; j < N; j++ {
			for l := 0; l < D; l++ {
				M[i].AddChild(H[j][l])
			}
		}
	}

	return K, D, N, R
}
