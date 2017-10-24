package spn

import (
	"github.com/RenatoGeh/gospn/sys"
)

// StorerTable is a DP table for Storer.
type StorerTable map[SPN]map[int]float64

// Storer allows for a mapping of node -> value, storing values for later use without having to
// recompute node values or derivatives (basically a Dynamic Programming table).
//
// Let T be an array of DP tables. T has n entries, with each entry T[i] being a distinct DP table
// and independent of other tables T[j], j != i. We start with T having zero entries. We call a
// ticket a key k such that T[k] is a new empty DP table. Each ticket is unique and represent
// distinct DP tables. A table T[k][S] has m entries, with each entry T[k][S][l] a float64.
//
// See NewTicket, Table and Value.
type Storer struct {
	// Tables.
	tables map[int]StorerTable
	// Number of tickets.
	tickets int
}

// NewStorer creates a new Storer pointer with an empty set of tables.
func NewStorer() *Storer {
	return &Storer{tables: make(map[int]StorerTable), tickets: 0}
}

// NewTicket creates a new Ticket k and creates an empty map T[k].
func (s *Storer) NewTicket() int {
	k := s.tickets
	s.tickets++
	s.tables[k] = make(StorerTable)
	return k
}

// Table returns the table at ticket position k (i.e. returns T[k]). It returns a double value,
// with the first being the table (if it exists), and the second a boolean indicating if such table
// exists.
func (s *Storer) Table(k int) (StorerTable, bool) {
	p, q := s.tables[k]
	return p, q
}

// Value returns the value of the SPN S in table T[k], returning two values: its value and whether
// the position T[k][S] exists.
func (s *Storer) Value(k int, S SPN) (map[int]float64, bool) {
	t, err := s.tables[k]
	if !err {
		return nil, err
	}
	return t.Value(S)
}

// Value returns the value of the SPN S in this StorerTable, returning two values: its value and
// whether the position exists.
func (t StorerTable) Value(S SPN) (map[int]float64, bool) {
	p, q := t[S]
	return p, q
}

// Entry returns the entry at Table T[k][S], given an SPN S and a position l, returning two values:
// the entry value and whether the position T[k][S][l] exists.
func (s *Storer) Entry(k int, S SPN, l int) (float64, bool) {
	t, err := s.tables[k]
	if !err {
		return 0, err
	}
	return t.Entry(S, l)
}

// Entry returns the entry at T[S][l], given an SPN S and a position l, returning two values:
// the entry value and whether the position T[S][l] exists.
func (s StorerTable) Entry(S SPN, l int) (float64, bool) {
	v, err := s[S]
	if !err {
		return 0, err
	}
	p, q := v[l]
	return p, q
}

// Single returns the first entry at Table T[k][S]. Equivalent to Entry(k, S, 0).
func (s *Storer) Single(k int, S SPN) (float64, bool) { return s.Entry(k, S, 0) }

// Single returns the first entry of this table. Equivalent to Entry(S, 0)
func (t StorerTable) Single(S SPN) (float64, bool) { return t.Entry(S, 0) }

// Store stores entry e in position T[k][S][l]. Returns whether the operation was successful.
func (s *Storer) Store(k int, S SPN, l int, e float64) bool {
	t, err := s.Table(k)
	if !err {
		return false
	}
	return t.Store(S, l, e)
}

// Store stores entry e in position [S][l]. Returns whether the operation was successful.
func (t StorerTable) Store(S SPN, l int, e float64) bool {
	_, err := t[S]
	if !err {
		t[S] = make(map[int]float64)
		t[S][l] = e
	} else {
		t[S][l] = e
	}
	return true
}

// StoreSingle sets the first entry of Table T[k][S] to v. Returns whether the operation was
// successful. Equivalent to Store(k, S, 0, v).
func (s *Storer) StoreSingle(k int, S SPN, v float64) bool { return s.Store(k, S, 0, v) }

// StoreSingle sets the first entry of this table to v. Returns whether the operation was
// successful. Equivalent to Store(S, 0, v)
func (t StorerTable) StoreSingle(S SPN, v float64) bool { return t.Store(S, 0, v) }

// Delete frees the memory at T[k]. A combination of Delete and NewTicket is NOT equivalent to
// using Reset. Prefer Reset over Delete + NewTicket.
func (s *Storer) Delete(k int) {
	s.tables[k] = nil
	sys.Free()
}

// Reset resets the values at T[k], deleting the map and creating a new one over it. The ticket
// remains the same. Reset returns the ticket k.
func (s *Storer) Reset(k int) int {
	s.tables[k] = nil
	sys.Free()
	s.tables[k] = make(StorerTable)
	return k
}
