package main

import (
	"fmt"
	"sort"
)

// Insert adds a key-value pair to the B+Tree
func (t *BPlusTree) Insert(key int, value string) {
	record := &BPlusTreeNodeRecord{Key: key, Value: value}
	
	if t.root == nil {
		// Create first leaf node
		leaf := &BPlusTreeNode{
			isLeaf:  true,
			keys:    []int{key},
			records: []*BPlusTreeNodeRecord{record},
		}
		t.root = leaf
		t.leafHead = leaf
		return
	}
	
	// Find the leaf node where this key should be inserted
	leaf := t.findLeaf(key)
	
	// Insert into leaf node
	t.insertIntoLeaf(leaf, key, record)
}

// FindLeaf finds the leaf node that should contain the given key
func (t *BPlusTree) findLeaf(key int) *BPlusTreeNode {
	current := t.root
	
	for !current.isLeaf {
		// Find the appropriate child to follow
		idx := sort.SearchInts(current.keys, key)
		if idx < len(current.keys) && key >= current.keys[idx] {
			idx++ // Go to the right child if key >= current key
		}
		if idx < len(current.children) {
			current = current.children[idx]
		} else if len(current.children) > 0 {
			current = current.children[len(current.children)-1]
		} else {
			break
		}
	}
	
	return current
}

// InsertIntoLeaf inserts a key-value pair into a leaf node
func (t *BPlusTree) insertIntoLeaf(leaf *BPlusTreeNode, key int, record *BPlusTreeNodeRecord) {
	// Find insertion position
	idx := sort.SearchInts(leaf.keys, key)
	
	// Insert key and record
	leaf.keys = append(leaf.keys, 0)
	leaf.records = append(leaf.records, nil)
	
	// Shift elements to make space
	copy(leaf.keys[idx+1:], leaf.keys[idx:])
	copy(leaf.records[idx+1:], leaf.records[idx:])
	
	// Insert new element
	leaf.keys[idx] = key
	leaf.records[idx] = record
	
	// Check if node needs splitting
	if len(leaf.keys) > t.order {
		t.splitLeafNode(leaf)
	}
}

// SplitLeafNode splits an overflowing leaf node
func (t *BPlusTree) splitLeafNode(leaf *BPlusTreeNode) {
	mid := len(leaf.keys) / 2
	
	// Create new leaf node
	newLeaf := &BPlusTreeNode{
		isLeaf:  true,
		keys:    make([]int, len(leaf.keys[mid:])),
		records: make([]*BPlusTreeNodeRecord, len(leaf.records[mid:])),
		next:    leaf.next,
		parent:  leaf.parent,
	}
	
	copy(newLeaf.keys, leaf.keys[mid:])
	copy(newLeaf.records, leaf.records[mid:])
	
	// Update original leaf
	leaf.keys = leaf.keys[:mid]
	leaf.records = leaf.records[:mid]
	leaf.next = newLeaf
	
	// Update leaf linked list
	if leaf == t.leafHead {
		t.leafHead = leaf
	}
	
	// Insert the new leaf into the parent
	t.insertIntoParent(leaf, newLeaf.keys[0], newLeaf)
}

// InsertIntoParent inserts a new child into the parent node after a split
func (t *BPlusTree) insertIntoParent(left *BPlusTreeNode, key int, right *BPlusTreeNode) {
	parent := left.parent
	
	if parent == nil {
		// Create new root
		t.root = &BPlusTreeNode{
			keys:     []int{key},
			children: []*BPlusTreeNode{left, right},
		}
		left.parent = t.root
		right.parent = t.root
		return
	}
	
	// Find insertion position in parent
	idx := sort.SearchInts(parent.keys, key)
	
	// Insert key and child
	parent.keys = append(parent.keys, 0)
	parent.children = append(parent.children, nil)
	
	// Shift elements
	copy(parent.keys[idx+1:], parent.keys[idx:])
	copy(parent.children[idx+1:], parent.children[idx:])
	
	// Insert new elements
	parent.keys[idx] = key
	parent.children[idx+1] = right
	right.parent = parent
	
	// Check if parent needs splitting
	if len(parent.keys) > t.order {
		t.splitInternalNode(parent)
	}
}

// SplitInternalNode splits an overflowing internal node
func (t *BPlusTree) splitInternalNode(node *BPlusTreeNode) {
	mid := len(node.keys) / 2
	promoteKey := node.keys[mid]
	
	// Create new internal node
	newNode := &BPlusTreeNode{
		keys:     make([]int, len(node.keys[mid+1:])),
		children: make([]*BPlusTreeNode, len(node.children[mid+1:])),
		parent:   node.parent,
	}
	
	copy(newNode.keys, node.keys[mid+1:])
	copy(newNode.children, node.children[mid+1:])
	
	// Update children's parent pointers
	for _, child := range newNode.children {
		child.parent = newNode
	}
	
	// Update original node
	node.keys = node.keys[:mid]
	node.children = node.children[:mid+1]
	
	// Insert new node into parent
	t.insertIntoParent(node, promoteKey, newNode)
}

// Search finds a record by key
func (t *BPlusTree) Search(key int) (*BPlusTreeNodeRecord, bool) {
	if t.root == nil {
		return nil, false
	}
	
	leaf := t.findLeaf(key)
	
	// Search within the leaf node
	idx := sort.SearchInts(leaf.keys, key)
	if idx < len(leaf.keys) && leaf.keys[idx] == key {
		return leaf.records[idx], true
	}
	
	return nil, false
}

// RangeSearch finds all records with keys in the range [start, end]
func (t *BPlusTree) RangeSearch(start, end int) []*BPlusTreeNodeRecord {
	if t.root == nil || start > end {
		return nil
	}
	
	results := []*BPlusTreeNodeRecord{}
	current := t.findLeaf(start)
	
	for current != nil {
		for i, key := range current.keys {
			if key >= start && key <= end {
				results = append(results, current.records[i])
			}
			if key > end {
				return results
			}
		}
		current = current.next
	}
	
	return results
}

// Delete removes a key from the B+Tree
func (t *BPlusTree) Delete(key int) bool {
	if t.root == nil {
		return false
	}
	
	leaf := t.findLeaf(key)
	idx := sort.SearchInts(leaf.keys, key)
	
	if idx >= len(leaf.keys) || leaf.keys[idx] != key {
		return false // Key not found
	}
	
	// Remove key and record
	leaf.keys = append(leaf.keys[:idx], leaf.keys[idx+1:]...)
	leaf.records = append(leaf.records[:idx], leaf.records[idx+1:]...)
	
	// Handle underflow
	if len(leaf.keys) < t.minKeys && leaf != t.root {
		t.handleLeafUnderflow(leaf)
	}
	
	// If root becomes empty after deletion
	if t.root != nil && len(t.root.keys) == 0 && !t.root.isLeaf {
		t.root = t.root.children[0]
		t.root.parent = nil
	}
	
	return true
}

// HandleLeafUnderflow handles underflow in leaf nodes
func (t *BPlusTree) handleLeafUnderflow(leaf *BPlusTreeNode) {
	// Try borrowing from left sibling
	leftSibling := t.getLeftSibling(leaf)
	if leftSibling != nil && len(leftSibling.keys) > t.minKeys {
		// Borrow from left sibling
		borrowedKey := leftSibling.keys[len(leftSibling.keys)-1]
		borrowedRecord := leftSibling.records[len(leftSibling.records)-1]
		
		leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
		leftSibling.records = leftSibling.records[:len(leftSibling.records)-1]
		
		// Insert borrowed element at beginning of leaf
		leaf.keys = append([]int{borrowedKey}, leaf.keys...)
		leaf.records = append([]*BPlusTreeNodeRecord{borrowedRecord}, leaf.records...)
		
		// Update parent's key
		t.updateParentKey(leaf, leaf.keys[0])
		return
	}
	
	// Try borrowing from right sibling
	rightSibling := t.getRightSibling(leaf)
	if rightSibling != nil && len(rightSibling.keys) > t.minKeys {
		// Borrow from right sibling
		borrowedKey := rightSibling.keys[0]
		borrowedRecord := rightSibling.records[0]
		
		rightSibling.keys = rightSibling.keys[1:]
		rightSibling.records = rightSibling.records[1:]
		
		// Insert borrowed element at end of leaf
		leaf.keys = append(leaf.keys, borrowedKey)
		leaf.records = append(leaf.records, borrowedRecord)
		
		// Update parent's key for right sibling
		t.updateParentKey(rightSibling, rightSibling.keys[0])
		return
	}
	
	// Merge with sibling
	if leftSibling != nil {
		t.mergeLeaves(leftSibling, leaf)
	} else if rightSibling != nil {
		t.mergeLeaves(leaf, rightSibling)
	}
}

// GetLeftSibling returns the left sibling of a node
func (t *BPlusTree) getLeftSibling(node *BPlusTreeNode) *BPlusTreeNode {
	if node.parent == nil {
		return nil
	}
	
	parent := node.parent
	idx := -1
	for i, child := range parent.children {
		if child == node {
			idx = i
			break
		}
	}
	
	if idx > 0 {
		return parent.children[idx-1]
	}
	return nil
}

// GetRightSibling returns the right sibling of a node
func (t *BPlusTree) getRightSibling(node *BPlusTreeNode) *BPlusTreeNode {
	if node.parent == nil {
		return nil
	}
	
	parent := node.parent
	idx := -1
	for i, child := range parent.children {
		if child == node {
			idx = i
			break
		}
	}
	
	if idx >= 0 && idx < len(parent.children)-1 {
		return parent.children[idx+1]
	}
	return nil
}

// UpdateParentKey updates the key in the parent that points to the node
func (t *BPlusTree) updateParentKey(node *BPlusTreeNode, newKey int) {
	if node.parent == nil {
		return
	}
	
	parent := node.parent
	for i, child := range parent.children {
		if child == node {
			if i > 0 {
				parent.keys[i-1] = newKey
			}
			break
		}
	}
}

// MergeLeaves merges two leaf nodes
func (t *BPlusTree) mergeLeaves(left, right *BPlusTreeNode) {
	// Merge keys and records
	left.keys = append(left.keys, right.keys...)
	left.records = append(left.records, right.records...)
	left.next = right.next
	
	// Remove right node from parent
	parent := right.parent
	if parent != nil {
		idx := -1
		for i, child := range parent.children {
			if child == right {
				idx = i
				break
			}
		}
		
		if idx >= 0 {
			// Remove key and child
			if idx > 0 {
				parent.keys = append(parent.keys[:idx-1], parent.keys[idx:]...)
			}
			parent.children = append(parent.children[:idx], parent.children[idx+1:]...)
			
			// Check if parent underflows
			if len(parent.keys) < t.minKeys && parent != t.root {
				t.handleInternalUnderflow(parent)
			}
		}
	}
}

// HandleInternalUnderflow handles underflow in internal nodes
func (t *BPlusTree) handleInternalUnderflow(node *BPlusTreeNode) {
	// Simplified implementation - in practice, this would be more complex
	// For this sample, we'll just leave it as a placeholder
}

// PrintTree prints the B+Tree structure
func (t *BPlusTree) PrintTree() {
	if t.root == nil {
		fmt.Println("Empty tree")
		return
	}
	
	queue := []*BPlusTreeNode{t.root}
	level := 0
	
	for len(queue) > 0 {
		fmt.Printf("Level %d: ", level)
		nextQueue := []*BPlusTreeNode{}
		
		for _, node := range queue {
			if node.isLeaf {
				fmt.Printf("[LEAF: %v] ", node.keys)
			} else {
				fmt.Printf("[INTERNAL: %v] ", node.keys)
				nextQueue = append(nextQueue, node.children...)
			}
		}
		
		fmt.Println()
		queue = nextQueue
		level++
	}
}

// PrintLeaves prints all leaf nodes in order
func (t *BPlusTree) PrintLeaves() {
	fmt.Print("Leaf nodes: ")
	current := t.leafHead
	for current != nil {
		fmt.Printf("%v -> ", current.keys)
		current = current.next
	}
	fmt.Println("nil")
}

// FakeDB represents a simple database using B+Tree
type FakeDB struct {
	tree *BPlusTree
}

// NewFakeDB creates a new fake database
func NewFakeDB() *FakeDB {
	return &FakeDB{
		tree: NewBPlusTree(4), // Order 4 B+Tree
	}
}

// Set inserts or updates a key-value pair
func (db *FakeDB) Set(key int, value string) {
	db.tree.Insert(key, value)
}

// Get retrieves a value by key
func (db *FakeDB) Get(key int) (string, bool) {
	record, found := db.tree.Search(key)
	if found {
		return record.Value, true
	}
	return "", false
}

// Delete removes a key-value pair
func (db *FakeDB) Delete(key int) bool {
	return db.tree.Delete(key)
}

// RangeQuery returns all key-value pairs in the range [start, end]
func (db *FakeDB) RangeQuery(start, end int) map[int]string {
	records := db.tree.RangeSearch(start, end)
	result := make(map[int]string)
	for _, record := range records {
		result[record.Key] = record.Value
	}
	return result
}

// Display shows the current state of the database
func (db *FakeDB) Display() {
	fmt.Println("=== Fake Database ===")
	db.tree.PrintTree()
	fmt.Println("=== Leaf Structure ===")
	db.tree.PrintLeaves()
}

func Main_bplustree_crud() {
	// Create a fake database
	db := NewFakeDB()
	
	// Insert some sample data
	fmt.Println("Inserting sample data...")
	db.Set(10, "Alice")
	db.Set(20, "Bob")
	db.Set(5, "Charlie")
	db.Set(15, "Diana")
	db.Set(25, "Eve")
	db.Set(30, "Frank")
	db.Set(35, "Grace")
	db.Set(40, "Henry")
	
	// Display the tree structure
	db.Display()
	
	// Search for specific keys
	fmt.Println("\n=== Search Operations ===")
	if value, found := db.Get(15); found {
		fmt.Printf("Key 15: %s\n", value)
	}
	if value, found := db.Get(99); found {
		fmt.Printf("Key 99: %s\n", value)
	} else {
		fmt.Println("Key 99: Not found")
	}
	
	// Range query
	fmt.Println("\n=== Range Query (10-30) ===")
	results := db.RangeQuery(10, 30)
	for key, value := range results {
		fmt.Printf("Key %d: %s\n", key, value)
	}
	
	// Delete operation
	fmt.Println("\n=== Delete Operation ===")
	if db.Delete(15) {
		fmt.Println("Key 15 deleted successfully")
	}
	if db.Delete(99) {
		fmt.Println("Key 99 deleted successfully")
	} else {
		fmt.Println("Key 99 not found for deletion")
	}
	
	// Display after deletion
	fmt.Println("\n=== After Deletion ===")
	db.Display()
	
	// Final range query
	fmt.Println("\n=== Final Range Query (5-35) ===")
	results = db.RangeQuery(5, 35)
	for key, value := range results {
		fmt.Printf("Key %d: %s\n", key, value)
	}
}
