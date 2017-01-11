package common

// Stack is a stack.
type Stack struct {
	data []interface{}
}

// Push puts element e on top of the pointer stack s.
func (s *Stack) Push(e interface{}) {
	s.data = append(s.data, e)
}

// Pop removes and returns the last element of pointer stack s.
func (s *Stack) Pop() interface{} {
	n := len(s.data) - 1
	e := s.data[n]
	s.data[n] = nil
	s.data = s.data[:n]
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
