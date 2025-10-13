package main

// BTreeNode represents a node in the B-tree
type BTreeNode struct {
	leaf     bool          // Whether this node is a leaf
	keys     []int         // Array of keys stored in the node
	children []*BTreeNode  // Array of child pointers
}

// BTree represents the B-tree data structure
type BTree struct {
	root *BTreeNode // Root node of the tree
	t    int        // Minimum degree (defines the range for number of keys)
}

// NewBTree creates a new B-tree with the given minimum degree
func NewBTree(t int) *BTree {
	if t < 2 {
		panic("Minimum degree must be at least 2")
	}
	return &BTree{
		root: nil,
		t:    t,
	}
}

// NewBTreeNode creates a new B-tree node
func NewBTreeNode(leaf bool, t int) *BTreeNode {
	return &BTreeNode{
		leaf:     leaf,
		keys:     make([]int, 0, 2*t-1), // Maximum 2t-1 keys
		children: make([]*BTreeNode, 0, 2*t), // Maximum 2t children
	}
}
