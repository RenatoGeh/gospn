package common

import (
	"github.com/RenatoGeh/gospn/sys"
)

// Stack is a stack of interface{}s.
type Stack struct {
	data []interface{}
	cap  int
}

// Push puts element e on top of the pointer stack s.
func (s *Stack) Push(e interface{}) {
	s.data = append(s.data, e)
	s.cap++
}

// Pop removes and returns the last element of pointer stack s.
func (s *Stack) Pop() interface{} {
	n := len(s.data) - 1
	e := s.data[n]
	s.data[n] = nil
	s.data = s.data[:n]
	if s.cap > sys.MemLowBoundShrink && s.cap >= (len(s.data)<<1) {
		s.Shrink()
	}
	return e
}

// Peek returns the top of the stack.
func (s *Stack) Peek() interface{} {
	return s.data[len(s.data)-1]
}

// Get returns the i-th element of the stack. Strongly discouraged, since this is a stack.
func (s *Stack) Get(i int) interface{} {
	return s.data[i]
}

// Size returns the size of pointer stack s.
func (s *Stack) Size() int { return len(s.data) }

// Empty returns whether pointer stack s is empty or not.
func (s *Stack) Empty() bool { return len(s.data) == 0 }

// Give is equivalent to Push.
func (s *Stack) Give(e interface{}) { s.Push(e) }

// Take is equivalent to Pop.
func (s *Stack) Take() interface{} { return s.Pop() }

// Shrink shrinks the queue to fit.
func (s *Stack) Shrink() {
	sys.Free()
	s.cap = len(s.data)
}

// Reset empties this queue.
func (s *Stack) Reset() { s.data, s.cap = nil, 0 }
