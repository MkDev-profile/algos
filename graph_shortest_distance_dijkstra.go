// Dijkstra's algorithm

package main

import (
	"fmt"
	"math"
)

// Edge represents a connection between two nodes with a weight
type Edge struct {
	To     int // Target node index
	Weight int // Weight/cost of the edge
}

// WeightedGraph represents a weighted graph using adjacency lists
type WeightedGraph struct {
	Nodes int      // Number of nodes in the graph
	Edges [][]Edge // Adjacency list: edges[i] contains all edges from node i
}

// NewWeightedGraph creates a new graph with n nodes
func NewWeightedGraph(n int) *WeightedGraph {
	return &WeightedGraph{
		Nodes: n,
		Edges: make([][]Edge, n),
	}
}

// AddEdge adds a directed edge from node u to node v with weight w
func (g *WeightedGraph) AddEdge(u, v, w int) {
	g.Edges[u] = append(g.Edges[u], Edge{To: v, Weight: w})
}

// AddBidirectionalEdge adds an undirected edge between node u and node v with weight w
func (g *WeightedGraph) AddBidirectionalEdge(u, v, w int) {
	g.AddEdge(u, v, w)
	g.AddEdge(v, u, w)
}

// PriorityQueueItem represents an item in the priority queue
type PriorityQueueItem struct {
	Node     int // Node index
	Distance int // Current distance from source
	Index    int // Index in the heap (for internal use)
}

// PriorityQueue implements a min-heap based priority queue
type PriorityQueue struct {
	items []*PriorityQueueItem // Slice storing the heap items
}

// NewPriorityQueue creates a new empty priority queue
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		items: make([]*PriorityQueueItem, 0),
	}
}

// Len returns the number of items in the priority queue
func (pq *PriorityQueue) Len() int {
	return len(pq.items)
}

// Less compares two items for the min-heap property (smaller distance has higher priority)
func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.items[i].Distance < pq.items[j].Distance
}

// Swap exchanges two items in the heap and updates their indices
func (pq *PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].Index = i
	pq.items[j].Index = j
}

// Push adds a new item to the priority queue
func (pq *PriorityQueue) Push(item *PriorityQueueItem) {
	// Set the item's index to the current end position
	item.Index = len(pq.items)
	// Append the item to the slice
	pq.items = append(pq.items, item)
	// Fix the heap property by moving the item up if necessary
	pq.up(len(pq.items) - 1)
}

// Pop removes and returns the item with the smallest distance
func (pq *PriorityQueue) Pop() *PriorityQueueItem {
	if len(pq.items) == 0 {
		return nil
	}

	// The root item has the smallest distance
	root := pq.items[0]
	lastIndex := len(pq.items) - 1

	// Move the last item to the root
	pq.items[0] = pq.items[lastIndex]
	pq.items[0].Index = 0

	// Remove the last item
	pq.items = pq.items[:lastIndex]

	// If there are items left, fix the heap property from the root down
	if len(pq.items) > 0 {
		pq.down(0)
	}

	return root
}

// up moves an item up in the heap until the heap property is restored
func (pq *PriorityQueue) up(index int) {
	for {
		parent := (index - 1) / 2

		// If we're at the root or the parent is smaller/equal, we're done
		if index == 0 || pq.Less(parent, index) {
			break
		}

		// Swap with parent and continue
		pq.Swap(parent, index)
		index = parent
	}
}

// down moves an item down in the heap until the heap property is restored
func (pq *PriorityQueue) down(index int) {
	for {
		left := 2*index + 1
		right := 2*index + 2
		smallest := index

		// Find the smallest among current node and its children
		if left < len(pq.items) && pq.Less(left, smallest) {
			smallest = left
		}
		if right < len(pq.items) && pq.Less(right, smallest) {
			smallest = right
		}

		// If current node is the smallest, we're done
		if smallest == index {
			break
		}

		// Swap with the smallest child and continue
		pq.Swap(index, smallest)
		index = smallest
	}
}

// DijkstraResult contains the results of Dijkstra's algorithm
type DijkstraResult struct {
	Distances []int  // Shortest distances from source to each node
	Previous  []int  // Previous node in the shortest path (-1 if unreachable)
	Visited   []bool // Which nodes were visited during the algorithm
}

// Dijkstra performs Dijkstra's algorithm to find the shortest paths from source to all nodes
func (g *WeightedGraph) Dijkstra(source int) *DijkstraResult {
	// Initialize result structures
	distances := make([]int, g.Nodes)
	previous := make([]int, g.Nodes)
	visited := make([]bool, g.Nodes)

	// Initialize all distances to "infinity" and previous nodes to -1 (unreachable)
	for i := 0; i < g.Nodes; i++ {
		distances[i] = math.MaxInt32
		previous[i] = -1
		visited[i] = false
	}

	// Distance from source to itself is 0
	distances[source] = 0

	// Create priority queue and add the source node
	pq := NewPriorityQueue()

	// Only add the source node initially, we'll add others as needed
	pq.Push(&PriorityQueueItem{
		Node:     source,
		Distance: 0,
	})

	// Main algorithm loop: process nodes in order of increasing distance
	for pq.Len() > 0 {
		// Extract the node with the smallest distance
		currentItem := pq.Pop()
		currentNode := currentItem.Node

		fmt.Printf("  currentItem: (node=%d dist=%d) (visited=%t)\n", currentItem.Node, currentItem.Distance, visited[currentNode])

		// Skip if we've already processed this node with a better distance
		if visited[currentNode] {
			continue
		}

		// Mark the current node as visited
		visited[currentNode] = true

		// Explore all neighbors of the current node
		for _, edge := range g.Edges[currentNode] {
			neighbor := edge.To

			fmt.Printf("neighbor: %d (visited=%t)\n", neighbor, visited[neighbor])

			// Skip if we've already found the shortest path to this neighbor
			if visited[neighbor] {
				continue
			}

			// Calculate the new distance to the neighbor through the current node
			newDistance := distances[currentNode] + edge.Weight

			// If we found a shorter path to the neighbor, update it
			if newDistance < distances[neighbor] {
				distances[neighbor] = newDistance
				previous[neighbor] = currentNode

				// Push the updated node to the priority queue
				pq.Push(&PriorityQueueItem{
					Node:     neighbor,
					Distance: newDistance,
				})
			}

			fmt.Print("pq: ")
			for _,val := range pq.items {
				fmt.Printf("(node=%d dist=%d)", (*val).Node, (*val).Distance)
			}
			fmt.Println()

		}
	}

	return &DijkstraResult{
		Distances: distances,
		Previous:  previous,
		Visited:   visited,
	}
}

// GetShortestPath reconstructs the shortest path from source to target
// Returns the path as a slice of node indices and the total distance
// If no path exists, returns nil and -1
func (result *DijkstraResult) GetShortestPath(target int) ([]int, int) {
	// If the target is unreachable, return nil
	if result.Previous[target] == -1 && target != 0 {
		return nil, -1
	}

	// Reconstruct the path by following the previous pointers backwards
	path := make([]int, 0)
	current := target

	// Work backwards from target to source
	for current != -1 {
		path = append([]int{current}, path...)
		current = result.Previous[current]
	}

	return path, result.Distances[target]
}

// Example usage and demonstration
func Main_graph_shortest_distance_dijkstra() {
	// Create a graph with 6 nodes (0-5)
	graph := NewWeightedGraph(7)

	// Add edges to create a sample graph
	graph.AddBidirectionalEdge(0, 1, 1) // Node 0 -> Node 1 with weight 1
	graph.AddBidirectionalEdge(0, 2, 6)
	graph.AddBidirectionalEdge(0, 3, 5)
	graph.AddBidirectionalEdge(1, 4, 1)
	graph.AddBidirectionalEdge(1, 2, 8)
	graph.AddBidirectionalEdge(2, 3, 1)
	graph.AddBidirectionalEdge(2, 6, 1)
	graph.AddBidirectionalEdge(3, 5, 1)
	graph.AddBidirectionalEdge(4, 5, 1)
	graph.AddBidirectionalEdge(5, 6, 5)

	// Run Dijkstra's algorithm from node 0
	source := 0
	result := graph.Dijkstra(source)

	// Print distances from source to all nodes
	fmt.Printf("Shortest distances from node %d:\n", source)
	for i, dist := range result.Distances {
		if dist == math.MaxInt32 {
			fmt.Printf("  Node %d: unreachable\n", i)
		} else {
			fmt.Printf("  Node %d: %d\n", i, dist)
		}
	}

	fmt.Println("\nDetailed paths:")
	// Find and print shortest paths to all reachable nodes
	for target := 0; target < graph.Nodes; target++ {
		path, distance := result.GetShortestPath(target)

		if path != nil {
			fmt.Printf("  Path from %d to %d (distance: %d): ", source, target, distance)
			for i, node := range path {
				if i > 0 {
					fmt.Print(" -> ")
				}
				fmt.Printf("%d", node)
			}
			fmt.Println()
		} else if target != source {
			fmt.Printf("  No path from %d to %d\n", source, target)
		}
	}
}

/*

PS D:\github\algos> go run .
  currentItem: (node=0 dist=0) (visited=false)
neighbor: 1 (visited=false)
pq: (node=1 dist=1)
neighbor: 2 (visited=false)
pq: (node=1 dist=1)(node=2 dist=6)
neighbor: 3 (visited=false)
pq: (node=1 dist=1)(node=2 dist=6)(node=3 dist=5)
  currentItem: (node=1 dist=1) (visited=false)
neighbor: 0 (visited=true)
neighbor: 4 (visited=false)
pq: (node=4 dist=2)(node=2 dist=6)(node=3 dist=5)
neighbor: 2 (visited=false)
pq: (node=4 dist=2)(node=2 dist=6)(node=3 dist=5)
  currentItem: (node=4 dist=2) (visited=false)
neighbor: 1 (visited=true)
neighbor: 5 (visited=false)
pq: (node=5 dist=3)(node=2 dist=6)(node=3 dist=5)
  currentItem: (node=5 dist=3) (visited=false)
neighbor: 3 (visited=false)
pq: (node=3 dist=4)(node=2 dist=6)(node=3 dist=5)
neighbor: 4 (visited=true)
neighbor: 6 (visited=false)
pq: (node=3 dist=4)(node=2 dist=6)(node=3 dist=5)(node=6 dist=8)
  currentItem: (node=3 dist=4) (visited=false)
neighbor: 0 (visited=true)
neighbor: 2 (visited=false)
pq: (node=2 dist=5)(node=3 dist=5)(node=6 dist=8)(node=2 dist=6)
neighbor: 5 (visited=true)
  currentItem: (node=2 dist=5) (visited=false)
neighbor: 0 (visited=true)
neighbor: 1 (visited=true)
neighbor: 3 (visited=true)
neighbor: 6 (visited=false)
pq: (node=3 dist=5)(node=6 dist=6)(node=6 dist=8)(node=2 dist=6)
  currentItem: (node=3 dist=5) (visited=true)
  currentItem: (node=2 dist=6) (visited=true)
  currentItem: (node=6 dist=6) (visited=false)
neighbor: 2 (visited=true)
neighbor: 5 (visited=true)
  currentItem: (node=6 dist=8) (visited=true)
Shortest distances from node 0:
  Node 0: 0
  Node 1: 1
  Node 2: 5
  Node 3: 4
  Node 4: 2
  Node 5: 3
  Node 6: 6

Detailed paths:
  Path from 0 to 0 (distance: 0): 0
  Path from 0 to 1 (distance: 1): 0 -> 1
  Path from 0 to 2 (distance: 5): 0 -> 1 -> 4 -> 5 -> 3 -> 2
  Path from 0 to 3 (distance: 4): 0 -> 1 -> 4 -> 5 -> 3
  Path from 0 to 4 (distance: 2): 0 -> 1 -> 4
  Path from 0 to 5 (distance: 3): 0 -> 1 -> 4 -> 5
  Path from 0 to 6 (distance: 6): 0 -> 1 -> 4 -> 5 -> 3 -> 2 -> 6

*/



























