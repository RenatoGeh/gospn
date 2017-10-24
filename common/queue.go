package common

import (
	"github.com/RenatoGeh/gospn/sys"
)

// Queue is a queue of interface{}s.
type Queue struct {
	data []interface{}
	cap  int
}

// Enqueue inserts element e at the end of the queue.
func (q *Queue) Enqueue(e interface{}) {
	q.data = append(q.data, e)
	q.cap++
}

// Dequeue removes and returns the first element of the queue.
func (q *Queue) Dequeue() interface{} {
	n := len(q.data)
	e := q.data[0]
	q.data[0] = nil
	q.data = q.data[1:n]
	if q.cap > sys.MemLowBoundShrink && q.cap >= (len(q.data)<<1) {
		q.Shrink()
	}
	return e
}

// DequeueBack removes and returns the last element of the queue.
func (q *Queue) DequeueBack() interface{} {
	n := len(q.data) - 1
	e := q.data[n]
	q.data[n] = nil
	q.data = q.data[1:n]
	if q.cap > sys.MemLowBoundShrink && q.cap >= (len(q.data)<<1) {
		q.Shrink()
	}
	return e
}

// Peek returns the first element of the queue.
func (q *Queue) Peek() interface{} {
	return q.data[0]
}

// Get returns the i-th element of the queue.
func (q *Queue) Get(i int) interface{} {
	return q.data[i]
}

// Size returns the size of the queue.
func (q *Queue) Size() int { return len(q.data) }

// Empty returns whether the queue is empty or not.
func (q *Queue) Empty() bool { return len(q.data) == 0 }

// Give is equivalent to Enqueue.
func (q *Queue) Give(e interface{}) { q.Enqueue(e) }

// Take is equivalent to Dequeue.
func (q *Queue) Take() interface{} { return q.Dequeue() }

// Shrink shrinks the queue to fit.
func (q *Queue) Shrink() {
	sys.Free()
	q.cap = len(q.data)
}
