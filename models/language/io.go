package language

import (
	"bufio"
	"fmt"
	"github.com/RenatoGeh/gospn/io"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/utils"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// T-rex: the king of regexers.
const rex = "(\\w+|[\\.]+|[\\,\\!\\@\\#\\$\\%\\^\\&\\*\\(\\)\\;\\\\\\/\\|\\<\\>\\\"\\'\\:\\`" +
	"\\=\\-\\?\\+\\{\\}\\[\\]])"

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
		v := match.FindAllString(in.Text(), -1)
		nv := len(v)

		if nv == 0 {
			continue
		}

		block = append(block, "")
		for i := 0; i < nv; i++ {
			str := strings.ToLower(v[i])
			id, ok := vocab[str]
			if !ok {
				vocab[str] = vc
				vc++
			} else {
				block[cblock] = utils.StringConcat(block[cblock], strconv.Itoa(id))
				block[cblock] = utils.StringConcat(block[cblock], " ")
				nwords++
			}
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
		vs[i] = v.block[v.ptr-i]
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
	fmt.Fscanf(voc, "%d ", &n)

	entries := make([]string, n)
	for i := 0; i < n; i++ {
		var j int
		var str string
		fmt.Fscanf(voc, "%d %s ", &j, &str)
		entries[j] = str
	}

	var m int
	fmt.Fscanf(voc, "%d ", &m)

	l, block := 0, make([]int, m)
	for i := 0; i < m; i++ {
		var k int
		fmt.Fscanf(voc, "%d ", &k)
		block[l] = k
		l++
	}

	return NewVocabulary(entries, block)
}
