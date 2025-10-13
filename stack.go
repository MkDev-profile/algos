package main

import ( 
    "errors" 
    "fmt" 
) 
 
// Stack is a generic stack implementation 
type Stack[T any] struct { 
    items []T 
} 
 
// NewStack creates a new empty generic stack 
func NewStack[T any]() *Stack[T] { 
    return &Stack[T]{ 
        items: make([]T, 0), 
    } 
} 
 
// Insert to end
// complexity = O(1)
// p.s. O(n) if resize slice.
func (s *Stack[T]) Push(item T) { 
    s.items = append(s.items, item) 
} 
 
// Remove from end
// complexity = O(1)
func (s *Stack[T]) Pop() (T, error) { 
    if s.IsEmpty() { 
        var zero T 
        return zero, errors.New("stack is empty") 
    } 
     
    lastIndex := len(s.items) - 1 
    item := s.items[lastIndex] 
    s.items = s.items[:lastIndex] 
    return item, nil 
} 
 
// Get from end
// complexity = O(1)
func (s *Stack[T]) Peek() (T, error) { 
    if s.IsEmpty() { 
        var zero T 
        return zero, errors.New("stack is empty") 
    } 
    return s.items[len(s.items)-1], nil 
} 
 
// complexity = O(1)
func (s *Stack[T]) IsEmpty() bool { 
    return len(s.items) == 0 
} 
 
// complexity = O(1)
func (s *Stack[T]) Size() int { 
    return len(s.items) 
} 
 
func Main_stack() {  

	// Integer stack 
    intStack := NewStack[int]() 
    intStack.Push(1) 
    intStack.Push(2) 
    intStack.Push(3) 
     
    fmt.Println("Integer stack:") 
    for !intStack.IsEmpty() { 
        item, _ := intStack.Pop() 
        fmt.Printf("%d ", item) 
    } 
    fmt.Println() 
     
    // String stack 
    stringStack := NewStack[string]() 
    stringStack.Push("hello") 
    stringStack.Push("world") 
     
    fmt.Println("String stack:") 
    for !stringStack.IsEmpty() { 
        item, _ := stringStack.Pop() 
        fmt.Printf("%s ", item) 
    } 
    fmt.Println() 
} 

/*

Theory:

--- Stack implementation based on linkedlist VS on slice:

-- linkedlist:
- push complexity = O(1)
- extra memory for field "next" field
- slower cache loading: nodes are allocated randomly in memory, so need fetch different cache-blocks(by 64 bytes).

-- slice:
- push complexity = "O(1) amortized" (if average rare slice resizings over many pushes, then average(amortized) complexity is O(1)).
- faster cache loading: elements are stored contiguously in memory, so fetch same cache-block for several nodes.

*/








