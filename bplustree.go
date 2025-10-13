package main

// BPlusTreeNodeRecord represents a key-value pair in the database
type BPlusTreeNodeRecord struct {
	Key   int
	Value string
}

// BPlusTreeNode represents a B+Tree node
type BPlusTreeNode struct {
	isLeaf   bool
	keys     []int
	children []*BPlusTreeNode
	records  []*BPlusTreeNodeRecord
	next     *BPlusTreeNode // For leaf nodes, points to next leaf
	parent   *BPlusTreeNode
}

// BPlusTree represents the B+Tree structure
type BPlusTree struct {
	root     *BPlusTreeNode
	order    int // Maximum number of keys per node
	minKeys  int // Minimum number of keys per node (except root)
	leafHead *BPlusTreeNode // Head of the leaf linked list
}

// NewBPlusTree creates a new B+Tree with the given order
func NewBPlusTree(order int) *BPlusTree {
	if order < 3 {
		order = 3 // Minimum reasonable order
	}
	return &BPlusTree{
		root:    nil,
		order:   order,
		minKeys: (order - 1) / 2,
	}
}
