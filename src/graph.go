package main

type Graph struct {
	nodes int
	adj   [][]int
}

func NewGraph(nodes int) *Graph {
	adj := make([][]int, nodes)
	for i := range adj {
		adj[i] = make([]int, 0)
	}
	return &Graph{
		nodes: nodes,
		adj:   adj,
	}
}

func (g *Graph) AddEdge(u, v int) {
	g.adj[u] = append(g.adj[u], v)
}