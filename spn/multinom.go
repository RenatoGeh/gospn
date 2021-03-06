package spn

import (
	"bytes"
	"fmt"
	"github.com/RenatoGeh/gospn/sys"
	"math"
)

// Mode of a univariate distribution.
type Mode struct {
	// Value of variable when it is the highest.
	index int
	// Highest value of variable.
	val float64
}

func computeMode(pr []float64) (int, float64) {
	var i int
	var m float64
	for j, p := range pr {
		if p > m {
			i, m = j, p
		}
	}
	return i, m
}

// Multinomial represents a multinomial distribution.
type Multinomial struct {
	Node
	// Variable ID
	varid int
	// Discrete probability distribution
	pr []float64
	// Mode of pr. We pre-compute this to save time.
	mode Mode
}

// NewMultinomial constructs a new Multinomial.
func NewMultinomial(varid int, dist []float64) *Multinomial {
	n := len(dist)
	m := -1.0
	miv := make([]int, n)
	var mi int
	nmi := 0
	for i := 0; i < n; i++ {
		if dist[i] > m {
			m = dist[i]
			miv[0] = i
			mi = i
			nmi = 1
		} else {
			if dist[i] == m {
				miv[nmi] = i
				nmi++
			}
		}
	}
	if nmi > 1 {
		mi = miv[sys.RandIntn(nmi)]
	}
	return &Multinomial{Node{sc: []int{varid}}, varid, dist, Mode{mi, m}}
}

// NewCountingMultinomial constructs a new Multinomial from a count slice.
func NewCountingMultinomial(varid int, counts []int) *Multinomial {
	n := len(counts)

	pr := make([]float64, n)
	s := 0.0
	for i := 0; i < n; i++ {
		s += 1.0 + float64(counts[i])
		pr[i] = float64(1 + counts[i])
	}

	for i := 0; i < n; i++ {
		pr[i] /= float64(s)
	}

	m := -1.0
	var mi int
	miv := make([]int, n)
	nmi := 0
	for i := 0; i < n; i++ {
		if pr[i] > m {
			m = pr[i]
			miv[0] = i
			mi = i
			nmi = 1
		} else {
			if pr[i] == m {
				miv[nmi] = i
				nmi++
			}
		}
	}
	if nmi > 1 {
		mi = miv[sys.RandIntn(nmi)]
	}

	return &Multinomial{Node{sc: []int{varid}}, varid, pr, Mode{mi, m}}
}

// NewScopedCountingMultinomial does the same as NewCountingMultinomial except it allows multiple
// variable scope.
func NewScopedCountingMultinomial(varid int, esc []int, counts []int) *Multinomial {
	n := len(counts)

	pr := make([]float64, n)
	s := 0.0
	for i := 0; i < n; i++ {
		s += 1.0 + float64(counts[i])
		pr[i] = float64(1 + counts[i])
	}

	for i := 0; i < n; i++ {
		pr[i] /= float64(s)
	}

	var m float64
	var mi int

	for i := 0; i < n; i++ {
		if pr[i] > m {
			m = pr[i]
			mi = i
		}
	}

	return &Multinomial{Node{sc: esc}, varid, pr, Mode{mi, m}}
}

// NewEmptyMultinomial constructs a new empty Multinomial for learning.  Argument m is the
// cardinality of varid.
func NewEmptyMultinomial(varid, m int) *Multinomial {
	pr := make([]float64, m)

	for i := 0; i < m; i++ {
		pr[i] = 1.0 / float64(m)
	}

	return &Multinomial{Node{sc: []int{varid}}, varid, pr, Mode{0, pr[0]}}
}

// Type returns the type of this node.
func (m *Multinomial) Type() string { return "leaf" }

// Pr returns the discrete probability distribution.
func (m *Multinomial) Pr() []float64 { return m.pr }

// Value returns the probability of a certain valuation. That is Pr(X=val[varid]), where
// Pr is a probability function over a Multinomial distribution.
func (m *Multinomial) Value(val VarSet) float64 {
	v, ok := val[m.varid]
	if ok {
		return math.Log(m.pr[v])
	}
	return 0.0 // ln(1.0) = 0.0
}

// Max returns the MAP state given a valuation.
func (m *Multinomial) Max(val VarSet) float64 {
	v, ok := val[m.varid]
	if ok {
		return math.Log(m.pr[v])
	}
	return math.Log(m.mode.val)
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (m *Multinomial) ArgMax(val VarSet) (VarSet, float64) {
	retval := make(VarSet)
	v, ok := val[m.varid]

	if ok {
		retval[m.varid] = v
		return retval, math.Log(m.pr[v])
	}

	retval[m.varid] = m.mode.index
	return retval, math.Log(m.mode.val)
}

// Mean returns the mean of the distribution.
func (m *Multinomial) Mean() float64 {
	var mu float64
	for i, p := range m.pr {
		mu += p * float64(i)
	}
	return mu
}

// StdDev returns the standard deviation of the distribution.
func (m *Multinomial) StdDev() float64 {
	mu := m.Mean()
	var sd float64
	for i, _ := range m.pr {
		d := float64(i) - mu
		sd += d * d
	}
	sd = math.Sqrt(sd / float64(len(m.pr)))
	return sd
}

// MuSigma returns the mean and standard deviation of the distribution.
func (m *Multinomial) MuSigma() (float64, float64) {
	mu := m.Mean()
	var sd float64
	for i, _ := range m.pr {
		d := float64(i) - mu
		sd += d * d
	}
	sd = math.Sqrt(sd / float64(len(m.pr)))
	return mu, sd
}

// GobEncode serializes this multinomial node.
func (m *Multinomial) GobEncode() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%d %d", m.varid, len(m.pr))
	for _, p := range m.pr {
		fmt.Fprintf(&b, " %f", p)
	}
	return b.Bytes(), nil
}

// GobDecode unserializes this multinomial node.
func (m *Multinomial) GobDecode(data []byte) error {
	b := bytes.NewBuffer(data)
	var n int
	_, err := fmt.Fscanf(b, "%d %d", &m.varid, &n)
	if err != nil {
		return err
	}
	m.sc = []int{m.varid}
	m.pr = make([]float64, n)
	for i := 0; i < n; i++ {
		_, err = fmt.Fscanf(b, "%f", &m.pr[i])
		if err != nil {
			return err
		}
	}
	imode, mode := computeMode(m.pr)
	m.mode.index, m.mode.val = imode, mode
	return nil
}

// SubType returns this leaf's subtype.
func (m *Multinomial) SubType() string { return "multinomial" }
