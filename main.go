package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// GRAPH
type Graph struct {
	nodes int
	adj      [][]int
}

func NewGraph(nodes int) *Graph {
	adj := make([][]int, nodes)
	for i := range adj {
		adj[i] = make([]int, 0)
	}
	return &Graph{
		nodes: nodes,
		adj:      adj,
	}
}

func (g *Graph) AddEdge(u, v int) {
	g.adj[u] = append(g.adj[u], v)
}


// ALGORITHMS
func (g *Graph) TarjanSCC() [][]int {
	index := 0
	lowLink := make([]int, g.nodes)
	ids := make([]int, g.nodes)
	onStack := make([]bool, g.nodes)
	stack := make([]int, 0)
	result := make([][]int, 0)

	for i := 0; i < g.nodes; i++ {
		lowLink[i] = -1
		ids[i] = -1
		onStack[i] = false
	}

	for i := 0; i < g.nodes; i++ {
		if ids[i] == -1 {
			g.DFSSCC(i, &index, &lowLink, &ids, &onStack, &stack, &result)
		}
	}

	return result
}

func (g *Graph) DFSSCC(at int, index *int, lowLink, ids *[]int, onStack *[]bool, stack *[]int, result *[][]int) {
	*index++
	(*lowLink)[at] = *index
	(*ids)[at] = *index
	*stack = append(*stack, at)
	(*onStack)[at] = true

	for _, to := range g.adj[at] {
		if (*ids)[to] == -1 {
			g.DFSSCC(to, index, lowLink, ids, onStack, stack, result)
			(*lowLink)[at] = min((*lowLink)[at], (*lowLink)[to])
		} else if (*onStack)[to] {
			(*lowLink)[at] = min((*lowLink)[at], (*ids)[to])
		}
	}

	if (*lowLink)[at] == (*ids)[at] {
		component := make([]int, 0)
		for {
			node := (*stack)[len(*stack)-1]
			*stack = (*stack)[:len(*stack)-1]
			(*onStack)[node] = false
			component = append(component, node)
			if node == at {
				break
			}
		}
		*result = append(*result, component)
	}
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
	discovery := make([]int, g.nodes)
	lowLink := make([]int, g.nodes)
	bridges := make([][]int, 0)

	for i := 0; i < g.nodes; i++ {
		discovery[i] = -1
		lowLink[i] = -1
	}

	for i := 0; i < g.nodes; i++ {
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
	nodeSet := make(map[byte]bool)
	for _, line := range lines {
		edge := strings.Split(line, " ")
		if len(edge) != 2 {
			continue
		}
		u := edge[0][0]
		v := edge[1][0]
		nodeSet[u] = true
		nodeSet[v] = true
	}

	graph := NewGraph(len(nodeSet))
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
		for _, node := range component {
			fmt.Printf("%c ", node+'A')
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
