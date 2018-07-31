package score

import (
	"fmt"
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"os"
)

// S stores classification scores.
type S struct {
	hits        int
	misses      int
	total       int
	predictions []Pair
}

// Pair of predicted and expected classification labels.
type Pair struct {
	Predicted int
	Expected  int
}

// NewScore returns a new empty score table.
func NewScore() *S { return &S{} }

// Hits returns the number of correct classifications.
func (s *S) Hits() int { return s.hits }

// Misses returns the number of incorrect classifications.
func (s *S) Misses() int { return s.misses }

// Total returns the number of correct + incorrect classifications.
func (s *S) Total() int { return s.total }

// Register adds the predicted label and the expected label to the score table.
func (s *S) Register(predicted int, expected int) {
	if predicted == expected {
		s.hits++
	} else {
		s.misses++
	}
	s.total++
	s.predictions = append(s.predictions, Pair{predicted, expected})
}

// String returns the textual representation of this score table.
func (s *S) String() string {
	var str string
	str = fmt.Sprintf("Hits: %d\nMisses: %d\nTotal: %d\nAccuracy: %.5f\n", s.hits, s.misses, s.total,
		float64(s.hits)/float64(s.total))
	str += fmt.Sprintf("Wrong predictions:\n")
	for _, p := range s.predictions {
		if p.Predicted != p.Expected {
			str += fmt.Sprintf("  Expected %d, got %d.\n", p.Expected, p.Predicted)
		}
	}
	return str
}

// Evaluate takes a dataset, an array of expected labels ordered according to the dataset, an SPN
// and the label variable, and registers each predicted and expected values of the label variable
// in the dataset.
func (s *S) Evaluate(T spn.Dataset, L []int, N spn.SPN, classVar *learn.Variable) {
	st := spn.NewStorer()
	tk := st.NewTicket()
	v := classVar.Varid
	sys.Println("Evaluating scores...")
	n := len(T) / 10
	for i, I := range T {
		if i > 0 && i%n == 0 {
			sys.Printf("... %d%% ...\n", int(100.0*(float64(i)/float64(len(T)))))
		}
		l := I[v]
		delete(I, v)
		_, _, M := spn.StoreMAP(N, I, tk, st)
		s.Register(M[v], L[i])
		st.Reset(tk)
		I[v] = l
	}
}

// Merge absorbs all the information from the given score.
func (s *S) Merge(t *S) {
	s.hits += t.hits
	s.misses += t.misses
	s.total += t.total
	s.predictions = append(s.predictions, t.predictions...)
}

// Add returns the result of adding the two scores. This function leaves the original scores
// untouched, returning a new score structure.
func Add(s, t *S) *S {
	m := NewScore()
	m.hits = s.hits + t.hits
	m.misses = s.misses + t.misses
	m.total = s.total + t.total
	m.predictions = append(s.predictions, t.predictions...)
	return m
}

// Save saves this score table's textual representation to a file.
func (s *S) Save(filename string) {
	f, e := os.Create(filename)
	defer f.Close()
	if e != nil {
		panic(e)
	}
	fmt.Fprintln(f, s.String())
}

// Clear clears the score table.
func (s *S) Clear() {
	s.hits, s.misses, s.total = 0, 0, 0
	s.predictions = nil
}
