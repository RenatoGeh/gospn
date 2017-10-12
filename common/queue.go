package common

import (
	"github.com/RenatoGeh/gospn/spn"
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

/*************************************************************************************************/

// QueueSPN is a queue of spn.SPNs.
type QueueSPN struct {
	data []spn.SPN
}

// Enqueue inserts element e at the end of the queue.
func (q *QueueSPN) Enqueue(e spn.SPN) {
	q.data = append(q.data, e)
}

// Dequeue removes and returns the first element of the queue.
func (q *QueueSPN) Dequeue() spn.SPN {
	n := len(q.data)
	e := q.data[0]
	q.data[0] = nil
	q.data = q.data[1:n]
	return e
}

// Peek returns the first element of the queue.
func (q *QueueSPN) Peek() spn.SPN {
	return q.data[0]
}

// Get returns the i-th element of the queue.
func (q *QueueSPN) Get(i int) spn.SPN {
	return q.data[i]
}

// Size returns the size of the queue.
func (q *QueueSPN) Size() int { return len(q.data) }

// Empty returns whether the queue is empty or not.
func (q *QueueSPN) Empty() bool { return len(q.data) == 0 }

/*************************************************************************************************/

// QueueBFSPair is a queue of BFSPairs.
type QueueBFSPair struct {
	data []*BFSPair
}

// Enqueue inserts element e at the end of the queue.
func (q *QueueBFSPair) Enqueue(e *BFSPair) {
	q.data = append(q.data, e)
}

// Dequeue removes and returns the first element of the queue.
func (q *QueueBFSPair) Dequeue() *BFSPair {
	n := len(q.data)
	e := q.data[0]
	q.data[0] = nil
	q.data = q.data[1:n]
	return e
}

// Peek returns the first element of the queue.
func (q *QueueBFSPair) Peek() *BFSPair {
	return q.data[0]
}

// Get returns the i-th element of the queue.
func (q *QueueBFSPair) Get(i int) *BFSPair {
	return q.data[i]
}

// Size returns the size of the queue.
func (q *QueueBFSPair) Size() int { return len(q.data) }

// Empty returns whether the queue is empty or not.
func (q *QueueBFSPair) Empty() bool { return len(q.data) == 0 }

/*************************************************************************************************/

// QueueInteger is a queue of integers.
type QueueInteger struct {
	data []int
}

// Enqueue inserts element e at the end of the queue.
func (q *QueueInteger) Enqueue(e int) {
	q.data = append(q.data, e)
}

// Dequeue removes and returns the first element of the queue.
func (q *QueueInteger) Dequeue() int {
	n := len(q.data)
	e := q.data[0]
	q.data = q.data[1:n]
	return e
}

// Peek returns the first element of the queue.
func (q *QueueInteger) Peek() int {
	return q.data[0]
}

// Get returns the i-th element of the queue.
func (q *QueueInteger) Get(i int) int {
	return q.data[i]
}

// Size returns the size of the queue.
func (q *QueueInteger) Size() int { return len(q.data) }

// Empty returns whether the queue is empty or not.
func (q *QueueInteger) Empty() bool { return len(q.data) == 0 }
