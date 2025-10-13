package main

// node of BinarySearchTree
type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

// binary search tree
type BST struct {
	Root *TreeNode
}

// creates new empty BST
func NewBST() *BST {
	return &BST{}
}
