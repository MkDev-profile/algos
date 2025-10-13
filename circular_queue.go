package main

import (
	"errors"
	"fmt"
)

type CircularQueue[T any] struct {
	items    []T
	front    int
	rear     int
	size     int
	capacity int
}

func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	if capacity <= 0 {
		panic("capacity must be positive")
	}

	return &CircularQueue[T]{
		items:    make([]T, capacity),
		front:    0,
		rear:     0,
		size:     0,
		capacity: capacity,
	}
}

func (q *CircularQueue[T]) Enqueue(item T) error {
	if q.IsFull() {
		return errors.New("queue is full")
	}

	q.items[q.rear] = item
	q.rear = (q.rear + 1) % q.capacity
	q.size++
	return nil
}

func (q *CircularQueue[T]) Dequeue() (T, error) {
	var zero T
	if q.IsEmpty() {
		return zero, errors.New("queue is empty")
	}

	item := q.items[q.front]
	q.front = (q.front + 1) % q.capacity
	q.size--
	return item, nil
}

func (q *CircularQueue[T]) Peek() (T, error) {
	var zero T
	if q.IsEmpty() {
		return zero, errors.New("queue is empty")
	}
	return q.items[q.front], nil
}

func (q *CircularQueue[T]) IsEmpty() bool {
	return q.size == 0
}

func (q *CircularQueue[T]) IsFull() bool {
	return q.size == q.capacity
}

func (q *CircularQueue[T]) Size() int {
	return q.size
}

func Main_circular_queue() {
	queue := NewCircularQueue[int](5)

	// Enqueue all elements
	for i := 0; i < 4; i++ {
		queue.Enqueue(i)
	}

	fmt.Printf("After enqueue: %#v\n\n", queue)

	// Dequeue several items
	for i := 0; i < 3; i++ {
		queue.Dequeue()
	}

    fmt.Printf("After dequeue: %#v\n\n", queue)

	// Enqueue more elements
	queue.Enqueue(4)

	// wrapping around
	queue.Enqueue(5) // insert to index = 0
	queue.Enqueue(6) // insert to index = 1

    fmt.Printf("After wrapping: %#v\n\n", queue)
}

/*

output:

$ go run .
After enqueue: &main.CircularQueue[int]{items:[]int{0, 1, 2, 3, 0}, front:0, rear:4, size:4, capacity:5}

After dequeue: &main.CircularQueue[int]{items:[]int{0, 1, 2, 3, 0}, front:3, rear:4, size:1, capacity:5}

After wrapping: &main.CircularQueue[int]{items:[]int{5, 6, 2, 3, 4}, front:3, rear:2, size:4, capacity:5}

*/


/*

CircularQueue use-cases:
- Buffer

*/













