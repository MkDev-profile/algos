// heap for generic types which implement cmp.ordered interface.

// MinHeap is data structure that implements usecase "PriorityQueue(очередь приоритетов)".

/*

Complexity:
Insert: O(log n)
ExtractMin: O(log n)
Peek: O(1)
Heapify: O(n)

Space: O(n)

*/


package main

import (
    "cmp"
    "errors"
    "fmt"
)

// MinHeap is a generic min-heap that works with any ordered type
type MinHeap[T cmp.Ordered] struct {
    items []T
    size  int
}

// NewMinHeap creates a new generic MinHeap
func NewMinHeap[T cmp.Ordered]() *MinHeap[T] {
    return &MinHeap[T]{
        items: make([]T, 0),
        size:  0,
    }
}

// Insert adds a new value to the heap
func (h *MinHeap[T]) Insert(value T) {
    h.items = append(h.items, value)
    h.size++
    h.bubbleUp(h.size - 1)
}

// ExtractMin removes and returns the minimum value
func (h *MinHeap[T]) ExtractMin() (T, error) {
    var zero T
    if h.IsEmpty() {
        return zero, errors.New("heap is empty")
    }

    min := h.items[0]
    lastIndex := h.size - 1
    h.items[0] = h.items[lastIndex]
    h.items = h.items[:lastIndex]
    h.size--
    
    if h.size > 0 {
        h.bubbleDown(0)
    }
    
    return min, nil
}

// Peek returns the minimum value without removing it
func (h *MinHeap[T]) Peek() (T, error) {
    var zero T
    if h.IsEmpty() {
        return zero, errors.New("heap is empty")
    }
    return h.items[0], nil
}

// bubbleUp and bubbleDown implementations are similar but use generic comparisons
func (h *MinHeap[T]) bubbleUp(index int) {
    for index > 0 {
        parent := (index - 1) / 2
        if h.items[index] >= h.items[parent] {
            break
        }
        h.items[index], h.items[parent] = h.items[parent], h.items[index]
        index = parent
    }
}

func (h *MinHeap[T]) bubbleDown(index int) {
    for {
        left := 2*index + 1
        right := 2*index + 2
        smallest := index

        if left < h.size && h.items[left] < h.items[smallest] {
            smallest = left
        }

        if right < h.size && h.items[right] < h.items[smallest] {
            smallest = right
        }

        if smallest == index {
            break
        }

        h.items[index], h.items[smallest] = h.items[smallest], h.items[index]
        index = smallest
    }
}

// Other methods (Size, IsEmpty, etc.) remain the same as non-generic version
func (h *MinHeap[T]) Size() int {
    return h.size
}

func (h *MinHeap[T]) IsEmpty() bool {
    return h.size == 0
}

func Main_min_heap() {
    heap := NewMinHeap[int]()
    
    elements := []int{5, 3, 8, 1, 9, 2, 7}

    for _, elem := range elements {
        heap.Insert(elem)
        min, _ := heap.Peek()
        fmt.Printf("Inserted %d, current min: %d\n", elem, min)
    }
    
    for !heap.IsEmpty() {
        min, _ := heap.ExtractMin()
        fmt.Printf("%d ", min)
    }
}









