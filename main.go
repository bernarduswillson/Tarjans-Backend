package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make([]Node, 0),
		Edges: make([]Edge, 0),
	}
}

func (g *Graph) AddNode(id, label string) {
	g.Nodes = append(g.Nodes, Node{
		ID:    id,
		Label: label,
	})
}

func (g *Graph) AddEdge(from, to string) {
	g.Edges = append(g.Edges, Edge{
		From: from,
		To:   to,
	})
}

func main() {
	// Enable CORS globally
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Define the HTTP route and handler
	http.Handle("/", corsMiddleware(http.HandlerFunc(handleRoot)))

	// Start the web server on port 8080
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Check if the request method is POST
	if r.Method == "POST" {
		// Read the uploaded file
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Printf("Failed to read file: %s\n", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Print the file path
		log.Printf("Received file: %s\n", fileHeader.Filename)

		// Read file content
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Failed to read file: %s\n", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Read graph from input
		lines := strings.Split(string(data), "\n")
		nodeSet := make(map[string]bool)
		edges := make([]Edge, 0)

		for _, line := range lines {
			edge := strings.Split(line, " ")
			if len(edge) != 2 {
				continue
			}
			u := edge[0]
			v := edge[1]
			nodeSet[u] = true
			nodeSet[v] = true
			edges = append(edges, Edge{
				From: u,
				To:   v,
			})
		}

		nodes := make([]Node, 0)
		for node := range nodeSet {
			nodes = append(nodes, Node{
				ID:    node,
				Label: node,
			})
		}

		graph := Graph{
			Nodes: nodes,
			Edges: edges,
		}

		// Find SCC
		scc := TarjanSCC(&graph)
		sccResult := make([][]string, len(scc))
		for i, component := range scc {
			sccResult[i] = make([]string, len(component))
			for j, node := range component {
				sccResult[i][j] = node
			}
		}

		// Find bridges
		bridges := FindBridges(&graph)
		bridgeResult := make([]Edge, len(bridges))
		for i, bridge := range bridges {
			bridgeResult[i] = bridge
		}

		response := map[string]interface{}{
			"result": sccResult,
			"graph": map[string]interface{}{
				"nodes": nodes,
				"edges": bridgeResult,
			},
		}

		// Print the result to console log
		log.Println(response)

		// Send the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}


func TarjanSCC(graph *Graph) [][]string {
	index := 0
	lowLink := make(map[string]int)
	ids := make(map[string]int)
	onStack := make(map[string]bool)
	stack := make([]string, 0)
	result := make([][]string, 0)

	for _, node := range graph.Nodes {
		lowLink[node.ID] = -1
		ids[node.ID] = -1
		onStack[node.ID] = false
	}

	for _, node := range graph.Nodes {
		if ids[node.ID] == -1 {
			DFSSCC(graph, node.ID, &index, &lowLink, &ids, &onStack, &stack, &result)
		}
	}

	return result
}

func DFSSCC(graph *Graph, at string, index *int, lowLink, ids *map[string]int, onStack *map[string]bool, stack *[]string, result *[][]string) {
	*index++
	(*lowLink)[at] = *index
	(*ids)[at] = *index
	*stack = append(*stack, at)
	(*onStack)[at] = true

	for _, edge := range graph.Edges {
		if edge.From != at {
			continue
		}
		to := edge.To
		if (*ids)[to] == -1 {
			DFSSCC(graph, to, index, lowLink, ids, onStack, stack, result)
			(*lowLink)[at] = min((*lowLink)[at], (*lowLink)[to])
		} else if (*onStack)[to] {
			(*lowLink)[at] = min((*lowLink)[at], (*ids)[to])
		}
	}

	if (*lowLink)[at] == (*ids)[at] {
		component := make([]string, 0)
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

func FindBridges(graph *Graph) []Edge {
	discovery := make(map[string]int)
	lowLink := make(map[string]int)
	bridges := make([]Edge, 0)

	for _, node := range graph.Nodes {
		discovery[node.ID] = -1
		lowLink[node.ID] = -1
	}

	for _, node := range graph.Nodes {
		if discovery[node.ID] == -1 {
			DFSBridge(graph, node.ID, "", &discovery, &lowLink, &bridges)
		}
	}

	return bridges
}

func DFSBridge(graph *Graph, u, parent string, discovery, lowLink *map[string]int, bridges *[]Edge) {
	(*discovery)[u]++
	(*lowLink)[u] = (*discovery)[u]

	for _, edge := range graph.Edges {
		if edge.From != u {
			continue
		}
		v := edge.To
		if (*discovery)[v] == -1 {
			(*discovery)[v] = (*discovery)[u] + 1
			(*lowLink)[v] = (*discovery)[v]
			DFSBridge(graph, v, u, discovery, lowLink, bridges)
			if (*lowLink)[v] > (*discovery)[u] {
				(*bridges) = append((*bridges), Edge{
					From: u,
					To:   v,
				})
			}
			(*lowLink)[u] = min((*lowLink)[u], (*lowLink)[v])
		} else if v != parent {
			(*lowLink)[u] = min((*lowLink)[u], (*discovery)[v])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
