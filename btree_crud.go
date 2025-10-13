package main

import (
	"fmt"
)

// Search searches for a key in the B-tree
// Returns the node containing the key and the index, or nil if not found
func (btree *BTree) Search(key int) (*BTreeNode, int) {
	if btree.root == nil {
		return nil, -1
	}
	return btree.searchRecursive(btree.root, key)
}

// searchRecursive helper function for searching
func (btree *BTree) searchRecursive(node *BTreeNode, key int) (*BTreeNode, int) {
	// Find the first key greater than or equal to the search key
	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}

	// If the key is found at this node
	if i < len(node.keys) && key == node.keys[i] {
		return node, i
	}

	// If this is a leaf node and key not found, return nil
	if node.leaf {
		return nil, -1
	}

	// Recurse to the appropriate child
	return btree.searchRecursive(node.children[i], key)
}

// Insert inserts a new key into the B-tree
func (btree *BTree) Insert(key int) {
	// If tree is empty, create root node
	if btree.root == nil {
		btree.root = NewBTreeNode(true, btree.t)
		btree.root.keys = append(btree.root.keys, key)
		return
	}

	// If root is full, split it first
	if len(btree.root.keys) == 2*btree.t-1 {
		// Create new root
		newRoot := NewBTreeNode(false, btree.t)
		newRoot.children = append(newRoot.children, btree.root)
		
		// Split the old root
		btree.splitChild(newRoot, 0)
		
		// Update root
		btree.root = newRoot
	}

	// Insert key into the tree
	btree.insertNonFull(btree.root, key)
}

// insertNonFull inserts a key into a non-full node
func (btree *BTree) insertNonFull(node *BTreeNode, key int) {
	// Start from the rightmost key
	i := len(node.keys) - 1

	if node.leaf {
		// For leaf node: insert key at correct position
		
		// Create space for new key
		node.keys = append(node.keys, 0)
		
		// Find location to insert and shift keys to the right
		for i >= 0 && key < node.keys[i] {
			node.keys[i+1] = node.keys[i]
			i--
		}
		
		// Insert the new key
		node.keys[i+1] = key
	} else {
		// For internal node: find child to descend into
		
		// Find child which will contain the new key
		for i >= 0 && key < node.keys[i] {
			i--
		}
		i++ // Adjust index to point to correct child
		
		// If the child is full, split it first
		if len(node.children[i].keys) == 2*btree.t-1 {
			btree.splitChild(node, i)
			
			// After splitting, the middle key goes up to current node
			// Check which child will now have the new key
			if key > node.keys[i] {
				i++
			}
		}
		
		// Recurse to the appropriate child
		btree.insertNonFull(node.children[i], key)
	}
}

// splitChild splits the child node of the given parent at specified index
func (btree *BTree) splitChild(parent *BTreeNode, childIndex int) {
	t := btree.t
	child := parent.children[childIndex]
	
	// Create new node which will store (t-1) keys of child
	newChild := NewBTreeNode(child.leaf, t)
	
	// Copy the last (t-1) keys from child to newChild
	newChild.keys = append(newChild.keys, child.keys[t:]...)
	
	// If child is not leaf, copy the last t children
	if !child.leaf {
		newChild.children = append(newChild.children, child.children[t:]...)
	}
	
	// Reduce the number of keys in child
	child.keys = child.keys[:t-1]
	
	// Reduce the number of children in child if not leaf
	if !child.leaf {
		child.children = child.children[:t]
	}
	
	// Create space in parent for new child pointer
	parent.children = append(parent.children, nil)
	copy(parent.children[childIndex+2:], parent.children[childIndex+1:])
	parent.children[childIndex+1] = newChild
	
	// Create space in parent for new key (middle key from child)
	parent.keys = append(parent.keys, 0)
	copy(parent.keys[childIndex+1:], parent.keys[childIndex:])
	
	// Move the middle key from child to parent
	parent.keys[childIndex] = child.keys[t-1]
}

// Delete deletes a key from the B-tree
func (btree *BTree) Delete(key int) {
	if btree.root == nil {
		return
	}

	btree.deleteRecursive(btree.root, key)

	// If root becomes empty after deletion
	if len(btree.root.keys) == 0 {
		if btree.root.leaf {
			btree.root = nil
		} else {
			btree.root = btree.root.children[0]
		}
	}
}

// deleteRecursive recursively deletes a key from the subtree rooted at node
func (btree *BTree) deleteRecursive(node *BTreeNode, key int) {
	t := btree.t
	idx := btree.findKey(node, key)

	// Case 1: Key is present in this node
	if idx < len(node.keys) && node.keys[idx] == key {
		if node.leaf {
			// Case 1a: Key is in leaf node - simply remove it
			btree.removeFromLeaf(node, idx)
		} else {
			// Case 1b: Key is in internal node
			btree.removeFromInternal(node, idx)
		}
	} else {
		// Case 2: Key is not in this node
		if node.leaf {
			// Key doesn't exist in tree
			return
		}

		// Flag indicating if key is in the last child
		flag := (idx == len(node.keys))

		// If the child has less than t keys, fill it
		if len(node.children[idx].keys) < t {
			btree.fill(node, idx)
		}

		// If last child was merged, recurse on (idx-1)th child
		if flag && idx > len(node.keys) {
			btree.deleteRecursive(node.children[idx-1], key)
		} else {
			btree.deleteRecursive(node.children[idx], key)
		}
	}
}

// findKey finds the first key >= given key
func (btree *BTree) findKey(node *BTreeNode, key int) int {
	idx := 0
	for idx < len(node.keys) && node.keys[idx] < key {
		idx++
	}
	return idx
}

// removeFromLeaf removes key at given index from leaf node
func (btree *BTree) removeFromLeaf(node *BTreeNode, idx int) {
	node.keys = append(node.keys[:idx], node.keys[idx+1:]...)
}

// removeFromInternal removes key at given index from internal node
func (btree *BTree) removeFromInternal(node *BTreeNode, idx int) {
	t := btree.t
	key := node.keys[idx]

	// If the preceding child has at least t keys
	if len(node.children[idx].keys) >= t {
		// Find predecessor (largest key in left subtree)
		pred := btree.getPredecessor(node, idx)
		node.keys[idx] = pred
		// Delete predecessor recursively
		btree.deleteRecursive(node.children[idx], pred)
	} else if len(node.children[idx+1].keys) >= t {
		// If the succeeding child has at least t keys
		// Find successor (smallest key in right subtree)
		succ := btree.getSuccessor(node, idx)
		node.keys[idx] = succ
		// Delete successor recursively
		btree.deleteRecursive(node.children[idx+1], succ)
	} else {
		// If both children have less than t keys, merge them
		btree.merge(node, idx)
		btree.deleteRecursive(node.children[idx], key)
	}
}

// getPredecessor gets the predecessor of keys[idx]
func (btree *BTree) getPredecessor(node *BTreeNode, idx int) int {
	// Move to the rightmost node until we reach a leaf
	cur := node.children[idx]
	for !cur.leaf {
		cur = cur.children[len(cur.children)-1]
	}
	// Return the last key of the leaf
	return cur.keys[len(cur.keys)-1]
}

// getSuccessor gets the successor of keys[idx]
func (btree *BTree) getSuccessor(node *BTreeNode, idx int) int {
	// Move to the leftmost node until we reach a leaf
	cur := node.children[idx+1]
	for !cur.leaf {
		cur = cur.children[0]
	}
	// Return the first key of the leaf
	return cur.keys[0]
}

// fill fills the child at idx which has less than t-1 keys
func (btree *BTree) fill(node *BTreeNode, idx int) {
	t := btree.t
	if idx != 0 && len(node.children[idx-1].keys) >= t {
		// Borrow from previous sibling
		btree.borrowFromPrev(node, idx)
	} else if idx != len(node.keys) && len(node.children[idx+1].keys) >= t {
		// Borrow from next sibling
		btree.borrowFromNext(node, idx)
	} else {
		// Merge with sibling
		if idx != len(node.keys) {
			btree.merge(node, idx)
		} else {
			btree.merge(node, idx-1)
		}
	}
}

// borrowFromPrev borrows a key from the previous child
func (btree *BTree) borrowFromPrev(node *BTreeNode, idx int) {
	child := node.children[idx]
	sibling := node.children[idx-1]

	// Move all keys in child one step ahead
	child.keys = append([]int{0}, child.keys...)
	// Move all child pointers one step ahead if not leaf
	if !child.leaf {
		child.children = append([]*BTreeNode{nil}, child.children...)
	}

	// Set child's first key equal to keys[idx-1] from current node
	child.keys[0] = node.keys[idx-1]

	// Move sibling's last key up to current node
	if !child.leaf {
		child.children[0] = sibling.children[len(sibling.children)-1]
		sibling.children = sibling.children[:len(sibling.children)-1]
	}

	node.keys[idx-1] = sibling.keys[len(sibling.keys)-1]
	sibling.keys = sibling.keys[:len(sibling.keys)-1]
}

// borrowFromNext borrows a key from the next child
func (btree *BTree) borrowFromNext(node *BTreeNode, idx int) {
	child := node.children[idx]
	sibling := node.children[idx+1]

	// keys[idx] is inserted as the last key in child
	child.keys = append(child.keys, node.keys[idx])

	// Sibling's first child is inserted as the last child of child
	if !child.leaf {
		child.children = append(child.children, sibling.children[0])
		sibling.children = sibling.children[1:]
	}

	// The first key from sibling is inserted into keys[idx]
	node.keys[idx] = sibling.keys[0]

	// Remove the first key from sibling
	sibling.keys = sibling.keys[1:]
}

// merge merges child at idx with child at idx+1
func (btree *BTree) merge(node *BTreeNode, idx int) {
	child := node.children[idx]
	sibling := node.children[idx+1]

	// Pull a key from the current node and insert it into (t-1)th position of child
	child.keys = append(child.keys, node.keys[idx])

	// Copy keys from sibling to child
	child.keys = append(child.keys, sibling.keys...)

	// Copy children from sibling to child if not leaf
	if !child.leaf {
		child.children = append(child.children, sibling.children...)
	}

	// Remove the key from current node
	node.keys = append(node.keys[:idx], node.keys[idx+1:]...)

	// Remove the sibling pointer
	node.children = append(node.children[:idx+1], node.children[idx+2:]...)
}

// Traverse performs an in-order traversal of the B-tree
func (btree *BTree) Traverse() []int {
	result := make([]int, 0)
	if btree.root != nil {
		btree.traverseRecursive(btree.root, &result)
	}
	return result
}

// traverseRecursive helper for in-order traversal
func (btree *BTree) traverseRecursive(node *BTreeNode, result *[]int) {
	for i := 0; i < len(node.keys); i++ {
		if !node.leaf {
			btree.traverseRecursive(node.children[i], result)
		}
		*result = append(*result, node.keys[i])
	}
	if !node.leaf {
		btree.traverseRecursive(node.children[len(node.keys)], result)
	}
}

// PrintTree prints the B-tree structure
func (btree *BTree) PrintTree() {
	if btree.root == nil {
		fmt.Println("Tree is empty")
		return
	}
	btree.printTreeRecursive(btree.root, 0)
}

// printTreeRecursive helper for printing tree structure
func (btree *BTree) printTreeRecursive(node *BTreeNode, level int) {
	fmt.Printf("Level %d: ", level)
	for i := 0; i < len(node.keys); i++ {
		fmt.Printf("%d ", node.keys[i])
	}
	fmt.Println()

	if !node.leaf {
		for i := 0; i < len(node.children); i++ {
			btree.printTreeRecursive(node.children[i], level+1)
		}
	}
}

// GetHeight returns the height of the B-tree
func (btree *BTree) GetHeight() int {
	if btree.root == nil {
		return 0
	}
	return btree.getHeightRecursive(btree.root)
}

// getHeightRecursive helper for calculating height
func (btree *BTree) getHeightRecursive(node *BTreeNode) int {
	if node.leaf {
		return 1
	}
	return 1 + btree.getHeightRecursive(node.children[0])
}

// Example usage and testing
func Main_btree_crud() {
	// Create a B-tree with minimum degree 3
	btree := NewBTree(3)

	// Insert keys
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17, 3, 1, 8, 25}
	fmt.Println("Inserting keys:", keys)
	for _, key := range keys {
		btree.Insert(key)
	}

	// Print tree structure
	fmt.Println("\nB-tree structure:")
	btree.PrintTree()

	// Traverse and print keys in sorted order
	fmt.Println("\nIn-order traversal:", btree.Traverse())

	// Search for keys
	searchKeys := []int{6, 15, 25}
	for _, key := range searchKeys {
		node, index := btree.Search(key)
		if node != nil {
			fmt.Printf("Key %d found at index %d\n", key, index)
		} else {
			fmt.Printf("Key %d not found\n", key)
		}
	}

	// Delete some keys
	deleteKeys := []int{6, 13, 7, 4, 20}
	fmt.Println("\nDeleting keys:", deleteKeys)
	for _, key := range deleteKeys {
		btree.Delete(key)
		fmt.Printf("After deleting %d: %v\n", key, btree.Traverse())
	}

	// Print final tree structure
	fmt.Println("\nFinal B-tree structure:")
	btree.PrintTree()

	// Print height of the tree
	fmt.Printf("Height of B-tree: %d\n", btree.GetHeight())
}
