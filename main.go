package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Graph struct {
	vertices int
	adj      [][]int
}

func NewGraph(vertices int) *Graph {
	adj := make([][]int, vertices)
	for i := range adj {
		adj[i] = make([]int, 0)
	}
	return &Graph{
		vertices: vertices,
		adj:      adj,
	}
}

func (g *Graph) AddEdge(u, v int) {
	g.adj[u] = append(g.adj[u], v)
}

func (g *Graph) TarjanSCC() [][]int {
	index := 0
	lowLink := make([]int, g.vertices)
	ids := make([]int, g.vertices)
	onStack := make([]bool, g.vertices)
	stack := make([]int, 0)
	result := make([][]int, 0)

	var tarjanDFS func(int)
	tarjanDFS = func(v int) {
		lowLink[v] = index
		ids[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for _, u := range g.adj[v] {
			if ids[u] == 0 {
				tarjanDFS(u)
			}
			if onStack[u] && ids[u] < lowLink[v] {
				lowLink[v] = ids[u]
			}
		}

		if lowLink[v] == ids[v] {
			scc := make([]int, 0)
			w := -1
			for w != v {
				w, stack = stack[len(stack)-1], stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
			}
			result = append(result, scc)
		}
	}

	for v := 0; v < g.vertices; v++ {
		if ids[v] == 0 {
			tarjanDFS(v)
		}
	}

	return result
}

func (g *Graph) DFSBridge(u, parent int, discovery, lowLink *[]int, bridges *[][]int) {
	(*discovery)[u]++
	(*lowLink)[u] = (*discovery)[u]

	for _, v := range g.adj[u] {
		if (*discovery)[v] == -1 {
			(*discovery)[v] = (*discovery)[u] + 1
			(*lowLink)[v] = (*discovery)[v]
			g.DFSBridge(v, u, discovery, lowLink, bridges)
			if (*lowLink)[v] > (*discovery)[u] {
				(*bridges) = append((*bridges), []int{u, v})
			}
			(*lowLink)[u] = min((*lowLink)[u], (*lowLink)[v])
		} else if v != parent {
			(*lowLink)[u] = min((*lowLink)[u], (*discovery)[v])
		}
	}
}

func (g *Graph) FindBridges() [][]int {
	discovery := make([]int, g.vertices)
	lowLink := make([]int, g.vertices)
	bridges := make([][]int, 0)

	for i := 0; i < g.vertices; i++ {
		discovery[i] = -1
		lowLink[i] = -1
	}

	for i := 0; i < g.vertices; i++ {
		if discovery[i] == -1 {
			g.DFSBridge(i, -1, &discovery, &lowLink, &bridges)
		}
	}

	return bridges
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Read input from file
	filePath := "graph.txt"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file: %s\n", err.Error())
		return
	}

	// Read graph from input
	lines := strings.Split(string(data), "\n")
	vertexSet := make(map[byte]bool)
	for _, line := range lines {
		edge := strings.Split(line, " ")
		if len(edge) != 2 {
			continue
		}
		u := edge[0][0]
		v := edge[1][0]
		vertexSet[u] = true
		vertexSet[v] = true
	}

	graph := NewGraph(len(vertexSet))
	for _, line := range lines {
		edge := strings.Split(line, " ")
		if len(edge) != 2 {
			continue
		}
		u := edge[0][0]
		v := edge[1][0]
		graph.AddEdge(int(u-'A'), int(v-'A'))
	}

	// Find SCC
	scc := graph.TarjanSCC()
	fmt.Println("Strongly Connected Components:")
	for _, component := range scc {
		for _, vertex := range component {
			fmt.Printf("%c ", vertex+'A')
		}
		fmt.Println()
	}

	// Find bridges
	bridges := graph.FindBridges()
	fmt.Println("Bridges:")
	for _, bridge := range bridges {
		u, v := bridge[0], bridge[1]
		fmt.Printf("%c-%c\n", u+'A', v+'A')
	}
}
