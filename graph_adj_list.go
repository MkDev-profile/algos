package main

type SimpleGraph struct {
	// матрица смежности
	// key = вершина
	// value = те узлы, с которыми связана вершина (напрямую)
	adjList map[int][]int
}

func NewSimpleGraph(vertices int) *SimpleGraph {
	return &SimpleGraph{
		adjList: make(map[int][]int),
	}
}

// заset-ить ребро (сторону, связь) графа
func (g *SimpleGraph) AddEdge(u, v int) {
	g.adjList[u] = append(g.adjList[u], v)
	g.adjList[v] = append(g.adjList[v], u) // For undirected graph
}








