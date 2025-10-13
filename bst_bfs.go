package main

// BFS performs breadth-first search and returns values in level order ("по порядку of level-ов")
func (bst *BST) BFS() []int {
    if bst.Root == nil {
        return []int{}
    }

    var result []int
    queue := []*TreeNode{bst.Root}

    for len(queue) > 0 {
        // Dequeue the front node
        currentNode := queue[0]
        queue = queue[1:]

        // Process(handle) the current node:
		result = append(result, currentNode.Value)
		// OR: for example: search (p.s. breath-first-search):
		// if currentNode.Value == target { return true }

        // Enqueue left child
        if currentNode.Left != nil {
            queue = append(queue, currentNode.Left)
        }

        // Enqueue right child
        if currentNode.Right != nil {
            queue = append(queue, currentNode.Right)
        }
    }

    return result
}

// BFSWithLevels returns nodes grouped by levels
func (bst *BST) BFSWithLevels() [][]int {
    if bst.Root == nil {
        return [][]int{}
    }

    var result [][]int
    queue := []*TreeNode{bst.Root}

    for len(queue) > 0 {
        levelSize := len(queue)
        var currentLevel []int

        for i := 0; i < levelSize; i++ {
            // Dequeue
            currentNode := queue[0]
            queue = queue[1:]

            currentLevel = append(currentLevel, currentNode.Value)

            // Enqueue children
            if currentNode.Left != nil {
                queue = append(queue, currentNode.Left)
            }
            if currentNode.Right != nil {
                queue = append(queue, currentNode.Right)
            }
        }

        result = append(result, currentLevel)
    }

    return result
}
















