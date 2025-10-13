package main

// Insert adds a new value to the BST
func (bst *BST) Insert(value int) {
    if bst.Root == nil {
        bst.Root = &TreeNode{Value: value}
        return
    }
    bst.Root.insert(value)
}

func (node *TreeNode) insert(value int) {
    if value < node.Value {
        if node.Left == nil {
            node.Left = &TreeNode{Value: value}
        } else {
            node.Left.insert(value)
        }
    } else if value > node.Value {
        if node.Right == nil {
            node.Right = &TreeNode{Value: value}
        } else {
            node.Right.insert(value)
        }
    }
    // If value == node.Value, don't insert duplicates
}











