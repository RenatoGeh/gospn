package score

import (
	"fmt"
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/spn"
	"os"
)

// S stores classification scores.
type S struct {
	Hits        int
	Misses      int
	Total       int
	Predictions []Pair
}

// Pair of predicted and expected classification labels.
type Pair struct {
	Predicted int
	Expected  int
}

// NewScore returns a new empty score table.
func NewScore() *S { return &S{} }

// Register adds the predicted label and the expected label to the score table.
func (s *S) Register(predicted int, expected int) {
	if predicted == expected {
		s.Hits++
	} else {
		s.Misses++
	}
	s.Total++
	s.Predictions = append(s.Predictions, Pair{predicted, expected})
}

// String returns the textual representation of this score table.
func (s *S) String() string {
	var str string
	str = fmt.Sprintf("Hits: %d\nMisses: %d\nTotal: %d\nAccuracy: %.5f\n", s.Hits, s.Misses, s.Total,
		float64(s.Hits)/float64(s.Total))
	str += fmt.Sprintf("Wrong predictions:\n")
	for _, p := range s.Predictions {
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
	fmt.Println("Evaluating scores...")
	n := len(T) / 10
	for i, I := range T {
		if i > 0 && i%n == 0 {
			fmt.Printf("... %d%% ...\n", int(100.0*(float64(i)/float64(len(T)))))
		}
		delete(I, v)
		_, _, M := spn.StoreMAP(N, I, tk, st)
		s.Register(M[v], L[i])
		st.Reset(tk)
	}
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
