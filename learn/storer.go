package learn

import (
	"github.com/RenatoGeh/gospn/spn"
)

// Storer allows for a mapping of node -> value, storing values for later use without having to
// recompute node values or derivatives (basically a Dynamic Programming table).
//
// Let T be an array of DP tables. T has n entries, with each entry T[i] being a distinct DP table
// and independent of other tables T[j], j != i. We start with T having zero entries. We call a
// ticket a key k such that T[k] is a new empty DP table. Each ticket is unique and represent
// distinct DP tables.
//
// See NewTicket, Table and Value.
type Storer struct {
	// Tables.
	tables map[int]map[spn.SPN]float64
	// Number of tickets.
	tickets int
}

// NewStorer creates a new Storer pointer with an empty set of tables.
func NewStorer() *Storer {
	return &Storer{tables: make(map[int]map[spn.SPN]float64), tickets: 0}
}

// NewTicket creates a new Ticket k and creates an empty map T[k].
func (s *Storer) NewTicket() int {
	k := s.tickets
	s.tickets++
	s.tables[k] = make(map[spn.SPN]float64)
	return k
}

// Table returns the table at ticket position k (i.e. returns T[k]). It returns a double value,
// with the first being the table (if it exists), and the second a boolean indicating if such table
// exists.
func (s *Storer) Table(k int) (map[spn.SPN]float64, bool) {
	p, q := s.tables[k]
	return p, q
}

// Value returns the value of the SPN S in table T[k], returning two values: its value and whether
// the position T[k][S] exists.
func (s *Storer) Value(k int, S spn.SPN) (float64, bool) {
	t, err := s.tables[k]
	if !err {
		return 0, err
	}
	p, q := t[S]
	return p, q
}

// Delete frees the memory at T[k]. A combination of Delete and NewTicket is NOT equivalent to
// using Reset. Prefer Reset over Delete + NewTicket.
func (s *Storer) Delete(k int) {
	s.tables[k] = nil
}

// Reset resets the values at T[k], deleting the map and creating a new one over it. The ticket
// remains the same. Reset returns the ticket k.
func (s *Storer) Reset(k int) int {
	s.tables[k] = make(map[spn.SPN]float64)
	return k
}
