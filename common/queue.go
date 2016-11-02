package common

import (
	spn "github.com/RenatoGeh/gospn/spn"
)

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
