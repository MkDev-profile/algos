// Detect Linked List Cycle.
// using 2 pointers as "Slow & Fast"("tortoise and hare")

package main

import (
	"fmt"
	"log"
)

// HasCycle detects if a linked list has a cycle
func HasCycle(head *LinkedListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	slow, fast := head, head
	log.Printf("slow=%d fast=%d\n", slow.Value, fast.Value)

	for fast != nil && fast.Next != nil {
		slow = slow.Next // slow передвигается на 1 узел
		fast = fast.Next.Next // fast передвигается на 2 узла

		log.Printf("slow=%d fast=%d\n", slow.Value, fast.Value)

		// если в графе есть петля, то fast и slow встретятся (fast успеет пройти путь в 2 раза быстрее чем slow. они встретятся на том узле, на котором начинается петля). 
		// p.s. offset = путь от начала списка до начала петли
		// fast path = offset + (cycle*X) 
		// slow path = offset + (cycle*logX) 
		if slow == fast {
			return true
		}
	}

	return false
}

// Helper function to create a linked list with cycle
func createLinkedListWithCycle(values []int, cyclePos int) *LinkedListNode {
	if len(values) == 0 {
		return nil
	}

	nodes := make([]*LinkedListNode, len(values))
	for i, val := range values {
		nodes[i] = &LinkedListNode{Value: val}
		if i > 0 {
			nodes[i-1].Next = nodes[i]
		}
	}

	// Create cycle if cyclePos is valid
	if cyclePos >= 0 && cyclePos < len(values) {
		nodes[len(values)-1].Next = nodes[cyclePos]
	}

	return nodes[0]
}

func Main_cycledetection() {
	// Example 1: No cycle
	head1 := createLinkedListWithCycle([]int{1, 2, 3, 4, 5}, -1)
	fmt.Printf("List 1 has cycle: %t\n", HasCycle(head1))
}

/*

head1 := createLinkedListWithCycle([]int{1, 2, 3, 4, 5}, 0)
PS D:\github\algos> go run main.go
2025/09/17 21:47:01 slow=1 fast=1
2025/09/17 21:47:01 slow=2 fast=3
2025/09/17 21:47:01 slow=3 fast=5
2025/09/17 21:47:01 slow=4 fast=2
2025/09/17 21:47:01 slow=5 fast=4
2025/09/17 21:47:01 slow=1 fast=1
List 1 has cycle: true

head1 := createLinkedListWithCycle([]int{1, 2, 3, 4, 5}, -1)
PS D:\github\algos> go run main.go
2025/09/17 23:30:31 slow=1 fast=1
2025/09/17 23:30:31 slow=2 fast=3
2025/09/17 23:30:31 slow=3 fast=5
List 1 has cycle: false

*/
