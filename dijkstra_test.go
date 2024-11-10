package graph_search

import (
	"testing"
)

func TestConditionalDijkstra_ShortestPath(t *testing.T) {
	nodeA, nodeB, nodeC, nodeD, nodeE, nodeF := Node{ID: 0}, Node{ID: 1}, Node{ID: 2}, Node{ID: 3},
		Node{ID: 4}, Node{ID: 5}
	g := Graph{Nodes: make([]Node, 0, 6)}

	for _, n := range []Node{nodeA, nodeB, nodeC, nodeD, nodeE, nodeF} {
		g.AddNode(n)
	}

	g.RelateNodes(nodeA, nodeB, 1, Bidirectional, MetaData{})
	g.RelateNodes(nodeA, nodeE, 2, Bidirectional, MetaData{})
	g.RelateNodes(nodeE, nodeF, 2, Bidirectional, MetaData{})
	g.RelateNodes(nodeF, nodeD, 2, Bidirectional, MetaData{})
	g.RelateNodes(nodeB, nodeC, 1, Bidirectional, MetaData{})
	g.RelateNodes(nodeC, nodeD, 1, Bidirectional, MetaData{})

	//   b --------1-------c
	//  / 1                 1 \
	// a --2-- e --2-- f --2-- d
	response := NewDijkstra(Criteria{
		Source:  []int32{0}, //a
		Targets: []int32{5}, //f

	}).Run(g)

	expectedDistance := float32(4.0)
	c, _ := response.Costs.GetCost(5)
	if expectedDistance != c {
		t.Fatalf("got %f, expected %f", c, expectedDistance)
	}

}
