// thread-safe generic circular singly linked list.

package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// LNode represents a node in the singly linked list
type LNode[T comparable] struct {
    value T
    next  *LNode[T]
}

// SinglyLinkedList represents a thread-safe singly linked list
type SinglyLinkedList[T comparable] struct {
    head  *LNode[T]
    tail  *LNode[T]
    size  int32
    mutex sync.RWMutex
}

// New creates a new empty singly linked list
func New[T comparable]() *SinglyLinkedList[T] {
    return &SinglyLinkedList[T]{
        head: nil,
        tail: nil,
        size: 0,
    }
}

// Remove removes the first occurrence of a value from the list
func (list *SinglyLinkedList[T]) Remove(value T) bool {
    list.mutex.Lock()
    defer list.mutex.Unlock()

    if list.head == nil {
        return false
    }

    if list.head.value == value {
        list.head = list.head.next
        if list.head == nil {
            list.tail = nil
        }
        atomic.AddInt32(&list.size, -1)
        return true
    }

    current := list.head
    for current.next != nil {
        if current.next.value == value {
            current.next = current.next.next
            
            if current.next == nil {
                list.tail = current
            }
            
            atomic.AddInt32(&list.size, -1)
            return true
        }
        current = current.next
    }

    return false
}

// Search checks if a value exists in the list
func (list *SinglyLinkedList[T]) Search(value T) (*LNode[T], bool) {
    list.mutex.RLock()
    defer list.mutex.RUnlock()

    current := list.head
    for current != nil {
        if current.value == value {
            return current, true
        }
        current = current.next
    }
    return nil, false
}

// Size returns the number of elements in the list
func (list *SinglyLinkedList[T]) Size() int {
    return int(atomic.LoadInt32(&list.size))
}

// IsEmpty checks if the list is empty
func (list *SinglyLinkedList[T]) IsEmpty() bool {
    return list.Size() == 0
}

// ForEach applies a function to each element in the list
func (list *SinglyLinkedList[T]) ForEach(fn func(T)) {
    list.mutex.RLock()
    defer list.mutex.RUnlock()

    current := list.head
    for current != nil {
        fn(current.value)
        current = current.next
    }
}

func (list *SinglyLinkedList[T]) InsertAtPosition(position int, value T) bool {
    list.mutex.Lock()
    defer list.mutex.Unlock()

    // Validate position
    if position < 0 || position > int(atomic.LoadInt32(&list.size)) {
        return false
    }

    newNode := &LNode[T]{value: value, next: nil}

    // Insert at head (position 0)
    if position == 0 {
        newNode.next = list.head
        list.head = newNode
        if list.tail == nil {
            list.tail = newNode
        }
        atomic.AddInt32(&list.size, 1)
        return true
    }

    // Insert at tail (position == size)
    if position == int(atomic.LoadInt32(&list.size)) {
        list.tail.next = newNode
        list.tail = newNode
        atomic.AddInt32(&list.size, 1)
        return true
    }

    // Insert in the middle
    current := list.head
    for i := 0; i < position-1; i++ {
        current = current.next
    }
    
    newNode.next = current.next
    current.next = newNode
    atomic.AddInt32(&list.size, 1)
    return true
}

// Reverse reverses the linked list in-place - O(n)
func (list *SinglyLinkedList[T]) Reverse() {
    list.mutex.Lock()
    defer list.mutex.Unlock()
    
    var prev *LNode[T]
    current := list.head
    list.tail = list.head
    
    for current != nil {
        next := current.next
        current.next = prev
        prev = current
        current = next
    }
    
    list.head = prev
}

// String returns a string representation of the list
func (list *SinglyLinkedList[T]) String() string {
    list.mutex.RLock()
    defer list.mutex.RUnlock()

    if list.IsEmpty() {
        return "[]"
    }

    var sb strings.Builder
    sb.WriteString("[")

    current := list.head
    fmt.Fprintf(&sb, "%v", current.value)

    for current.next != nil {
        current = current.next
        fmt.Fprintf(&sb, " -> %v", current.value)
    }

    sb.WriteString("]")
    return sb.String()
}









