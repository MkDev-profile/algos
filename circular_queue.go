package main

import (
	"errors"
	"fmt"
)

type CircularQueue[T any] struct {
	items    []T
	startIdx    int
	endIdx     int
	len     int
	cap int
}

func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	if capacity <= 0 {
		panic("capacity must be positive")
	}

	return &CircularQueue[T]{
		items:    make([]T, capacity),
		startIdx:    0,
		endIdx:     -1,
		len:     0,
		cap: capacity,
	}
}

func (q *CircularQueue[T]) Enqueue(item T) error {
	if q.len == q.cap {
		return errors.New("queue is full")
	}

	endIdx := (q.endIdx + 1) % q.cap
	q.items[endIdx] = item
	q.endIdx = endIdx
	q.len++
	return nil
}

func (q *CircularQueue[T]) Dequeue() (T, error) {
	var zero T
	if q.len == 0 {
		return zero, errors.New("queue is empty")
	}

	item := q.items[q.startIdx]
	q.startIdx = (q.startIdx + 1) % q.cap
	q.len--
	return item, nil
}

func (q *CircularQueue[T]) Peek() (T, error) {
	var zero T
	if q.len == 0 {
		return zero, errors.New("queue is empty")
	}
	
	return q.items[q.startIdx], nil
}

func (q *CircularQueue[T]) IsEmpty() bool {
	return q.len == 0
}

func (q *CircularQueue[T]) IsFull() bool {
	return q.len == q.cap
}

func (q *CircularQueue[T]) Count() int {
	return q.len
}

func Main_circular_queue() {
	queue := NewCircularQueue[int](5)

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)
	queue.Enqueue(4)

	fmt.Printf("%#+v\n", queue)

	queue.Dequeue()
	queue.Dequeue()

	fmt.Printf("%#+v\n", queue)

	// wrap-around
	queue.Enqueue(5)
	queue.Enqueue(6)
	queue.Enqueue(7)

	fmt.Printf("%#+v\n", queue)
}

/*

output:

&main.CircularQueue[int]{items:[]int{1, 2, 3, 4, 0}, startIdx:0, endIdx:3, len:4, cap:5}
&main.CircularQueue[int]{items:[]int{1, 2, 3, 4, 0}, startIdx:2, endIdx:3, len:2, cap:5}
&main.CircularQueue[int]{items:[]int{6, 7, 3, 4, 5}, startIdx:2, endIdx:1, len:5, cap:5}

*/


/*

CircularQueue use-cases:
- Buffer

*/













