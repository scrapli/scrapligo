package util

import (
	"sync"
)

// Queue is a simple queue structure to store and queue/requeue/dequeue bytes.
type Queue struct {
	queue [][]byte
	depth int
	lock  *sync.Mutex
}

// NewQueue returns a prepared Queue object.
func NewQueue() *Queue {
	return &Queue{
		lock: &sync.Mutex{},
	}
}

// Requeue prepends some bytes to the front of the queue.
func (q *Queue) Requeue(b []byte) {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := [][]byte{b}
	q.queue = append(n, q.queue...)

	q.depth++
}

// Enqueue queues some bytes at the end of the queue.
func (q *Queue) Enqueue(b []byte) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queue = append(q.queue, b)
	q.depth++
}

// Dequeue returns the first bytes in the queue.
func (q *Queue) Dequeue() []byte {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.queue) == 0 {
		return nil
	}

	b := q.queue[0]

	q.queue = q.queue[1:]
	q.depth--

	return b
}

// GetDepth returns the depth of the queue.
func (q *Queue) GetDepth() int {
	return q.depth
}
