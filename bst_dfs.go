package main

import "fmt"

// DFS with different traversal orders

// pre-order traversal (current node result-им "в первую очередь")
func dfsTreePreOrder(root *TreeNode) {
	if root == nil {
		return
	}

	// Process current node (pre-order)
	fmt.Printf("%d ", root.Value)

	// Recursively visit left and right children
	dfsTreePreOrder(root.Left)
	dfsTreePreOrder(root.Right)
}

// in-order (current node result-им "в середине")
func dfsTreeInOrder(root *TreeNode) {
	if root == nil {
		return
	}

	dfsTreeInOrder(root.Left)
	fmt.Printf("%d ", root.Value) // Process node (in-order)
	dfsTreeInOrder(root.Right)
}

// post-order (current node result-им "в последнюю очередь")
func dfsTreePostOrder(root *TreeNode) {
	if root == nil {
		return
	}

	dfsTreePostOrder(root.Left)
	dfsTreePostOrder(root.Right)
	fmt.Printf("%d ", root.Value) // Process node (post-order)
}

func Main_bst_dfs() {
	fmt.Println("=== Tree DFS Example ===")

	// Create a sample binary tree
	//       1
	//      / \
	//     2   3
	//    / \
	//   4   5

	root := &TreeNode{Value: 1}
	root.Left = &TreeNode{Value: 2}
	root.Right = &TreeNode{Value: 3}
	root.Left.Left = &TreeNode{Value: 4}
	root.Left.Right = &TreeNode{Value: 5}

	fmt.Print("Pre-order DFS: ")
	dfsTreePreOrder(root) // Pre-order DFS: 1 2 4 5 3
	fmt.Println()

	fmt.Print("In-order DFS: ")
	dfsTreeInOrder(root) // In-order DFS: 4 2 5 1 3
	fmt.Println()

	fmt.Print("Post-order DFS: ")
	dfsTreePostOrder(root) // Post-order DFS: 4 5 2 3 1
	fmt.Println()
}
