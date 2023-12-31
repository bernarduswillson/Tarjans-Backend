package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// enable CORS globally
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

	http.Handle("/", corsMiddleware(http.HandlerFunc(handleRoot)))

	// start the web server on port 8080
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {
		// read uploaded file
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Printf("Failed to read file: %s\n", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		log.Printf("Received file: %s\n", fileHeader.Filename)

		// read file content
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Failed to read file: %s\n", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
			// return errors.New("Internal Server Error")
		}

		// read graph from input
		lines := strings.Split(string(data), "\n")
		nodeSet := make(map[byte]bool)
		resultGraph := make([]string, 0)
		for _, line := range lines {
			edge := strings.Split(line, " ")
			if len(edge) != 2 {
				continue
			}
			u := edge[0][0]
			v := edge[1][0]
			nodeSet[u] = true
			nodeSet[v] = true
			resultGraph = append(resultGraph, fmt.Sprintf("%c%c", u, v))
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

		startTime := time.Now()

		// find SCC
		scc := graph.FindSCC()
		responseSCC := make([]string, 0)
		for _, component := range scc {
			for _, node := range component {
				comp := fmt.Sprintf("%c", node+'A')
				// iterate the graph edges to find adjacent nodes
				for _, edge := range graph.adj[node] {
					// check if the adjacent node is also in the component
					for _, compNode := range component {
						if edge == compNode {
							comp += fmt.Sprintf("%c", edge+'A')
						}
					}
				}
				responseSCC = append(responseSCC, comp)
			}
		}

		// find bridges
		bridges := graph.FindBridges()
		responseBridge := make([]string, 0)
		for _, bridge := range bridges {
			comp := ""
			u, v := bridge[0], bridge[1]
			comp += fmt.Sprintf("%c%c", u+'A', v+'A')
			responseBridge = append(responseBridge, comp)
		}

		// calculate elapsed time
		elapsedTime := time.Since(startTime).Nanoseconds()
		fmt.Printf("Elapsed time: %d nanoseconds\n", elapsedTime)

		response := map[string][]string{
			"graph":  resultGraph,
			"scc":    responseSCC,
			"bridge": responseBridge,
			"time":   {fmt.Sprintf("%d", elapsedTime)},
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal response: %s\n", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
