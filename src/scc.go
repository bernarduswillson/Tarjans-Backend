package main

func (g *Graph) FindSCC() [][]int {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
