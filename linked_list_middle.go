// Middle of the Linked List.
// using 2 pointers as "Slow & Fast"("tortoise and hare")

package main

func middleNode(head *LinkedListNode) *LinkedListNode {
    // initialize two pointers - slow and fast
    slow, fast := head, head
    
    // fast pointer движется в 2 раза быстрее чем slow pointer. когда slow pointer пройдет путь один раз, fast pointer успеет пройти этот путь 2 раза. fast и slow pointer встретятся в конце пути.
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    
    // when "fast" completes first traverse, "slow" will be at the middle
    return slow
}







