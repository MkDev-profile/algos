package main

import (
	"fmt"
	"sync"
)

// ThreadSafeQueue is a generic thread-safe FIFO queue
type ThreadSafeQueue[T any] struct {
	items []T
	mu    sync.RWMutex
}

// NewThreadSafeQueue creates a new thread-safe queue
func NewThreadSafeQueue[T any]() *ThreadSafeQueue[T] {
	return &ThreadSafeQueue[T]{
		items: make([]T, 0),
	}
}

// Enqueue adds an item to the end of the queue
// complexity = O(1)
// p.s. O(n) if resize slice.
func (q *ThreadSafeQueue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)
}

// Dequeue removes and returns the item from the front of the queue
// Returns false if the queue is empty
// complexity = O(1)
func (q *ThreadSafeQueue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek returns the front item without removing it
// Returns false if the queue is empty
// complexity = O(1)
func (q *ThreadSafeQueue[T]) Peek() (T, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.items) == 0 {
		var zero T
		return zero, false
	}

	return q.items[0], true
}

// IsEmpty returns true if the queue is empty
// complexity = O(1)
func (q *ThreadSafeQueue[T]) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.items) == 0
}

// Size returns the number of items in the queue
// complexity = O(1)
func (q *ThreadSafeQueue[T]) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.items)
}

func Main_queue() {
	var wg sync.WaitGroup
	queue := NewThreadSafeQueue[int]()
	
	// Multiple producers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				item := id*10 + j
				queue.Enqueue(item)
				fmt.Printf("Producer %d enqueued: %d\n", id, item)
			}
		}(i)
	}
	
	// Multiple consumers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				if item, ok := queue.Dequeue(); ok {
					fmt.Printf("Consumer %d dequeued: %d\n", id, item)
				}
			}
		}(i)
	}
	
	wg.Wait()
	fmt.Printf("Final queue size: %d\n", queue.Size())
}

















