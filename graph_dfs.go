package main

import "fmt"

// DFS recursive implementation for graph
func (g *SimpleGraph) DFS(start int) {
	visited := make(map[int]bool)
	g.dfsRecursive(start, visited)
}

func (g *SimpleGraph) dfsRecursive(vertex int, visited map[int]bool) {
	// Mark current vertex as visited

	fmt.Printf("%d ", vertex)
	visited[vertex] = true

	// Visit all adjacent vertices
	for _, neighbor := range g.adjList[vertex] {
		if !visited[neighbor] {
			g.dfsRecursive(neighbor, visited)
		}
	}
}

// DFS iterative implementation using stack
func (g *SimpleGraph) DFSIterative(start int) {
	visited := make(map[int]bool) // будет хранить все узлы графа
	stack := []int{start}

	for len(stack) > 0 {
		// Pop from stack
		vertex := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// handle and mark as visited
		fmt.Printf("%d ", vertex) // some processing (обработка) узла
		visited[vertex] = true

		// Push all unvisited neighbors to stack
		// то есть более-стартовые узлы не попадут в стек повторно (например если cycle(петля) в графе)
		for _, neighbor := range g.adjList[vertex] {
			if !visited[neighbor] {
				stack = append(stack, neighbor)
			}
		}
	}
}

func Main_graph_dfs() {
	fmt.Println("\n=== Graph DFS Example ===")

	// Create a sample graph
	// 0 -- 1
	// |    |
	// 2 -- 3
	graph := NewSimpleGraph(4)
	graph.AddEdge(0, 1)
	graph.AddEdge(0, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(2, 3)

	fmt.Print("Recursive DFS starting from vertex 0: ")
	graph.DFS(0)
	fmt.Println()

	fmt.Print("Iterative DFS starting from vertex 0: ")
	graph.DFSIterative(0)
	fmt.Println()
}
