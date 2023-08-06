package main

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