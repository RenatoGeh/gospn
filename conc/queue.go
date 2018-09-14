package conc

import (
	"runtime"
	"sync"
)

// SingleQueue is a concurrency queue that supports only one function. Each process waits its turn
// to execute. A number n of processes may run at concurrently. The queue waits for all processes
// to finish.
type SingleQueue struct {
	wg sync.WaitGroup
	n  int // Total number of allowed running processes.
	k  int // Number of currently running processes.
	c  *sync.Cond
}

// NewSingleQueue creates a new concurrency queue of type SingleQueue. Parameter n is the number of
// processes allowed to be run concurrently. If n <= 0, then n=runtime.NumCPU().
func NewSingleQueue(n int) *SingleQueue {
	if n <= 0 {
		n = runtime.NumCPU()
	}
	return &SingleQueue{n: n, k: 0, c: sync.NewCond(&sync.Mutex{})}
}

func (q *SingleQueue) Allowed() int { return q.n }

// Run queries the SingleQueue to run function f with id i. It may or may not run immediately.
func (q *SingleQueue) Run(f func(int), i int) {
	q.c.L.Lock()
	for q.k >= q.n {
		q.c.Wait()
	}
	q.k++
	q.c.L.Unlock()
	q.wg.Add(1)
	go func(id int) {
		defer q.wg.Done()
		f(id)
		q.c.L.Lock()
		q.k--
		q.c.L.Unlock()
		q.c.Signal()
	}(i)
}

// Wait waits for all the processes to finish.
func (q *SingleQueue) Wait() {
	q.wg.Wait()
}
