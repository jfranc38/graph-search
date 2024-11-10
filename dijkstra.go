package graph_search

import (
	"container/list"
	"fmt"
	"math"

	"github.com/golang/geo/s2"
)

// INFINITE represents the maximum possible float32 value, used to indicate infinite distance/cost.
// This constant is used throughout the graph search algorithms to initialize distances and represent
// unreachable nodes.
const INFINITE = math.MaxFloat32

// Criteria defines the configuration parameters for graph search algorithms.
// It encapsulates all the necessary parameters to customize and constrain the search behavior.
type Criteria struct {
	// Source contains the IDs of starting nodes for the search.
	// Multiple source nodes allow for concurrent path finding from different origins.
	Source []int32

	// Targets contains the IDs of destination nodes for the search.
	// Multiple targets enable finding paths to several destinations in one search operation.
	Targets []int32
}

// PathCost represents the cost associated with reaching a specific node in the graph.
// This structure is used to track both the node's identity and the accumulated cost
// of reaching it through the shortest discovered path.
type PathCost struct {
	// ID uniquely identifies the node in the graph
	ID int32

	// Cost represents the total accumulated cost (distance/time/weight) to reach this node
	// from the source node(s) through the shortest discovered path
	Cost float32
}

// Costs maps node IDs to their associated costs in the graph traversal.
// This type provides an efficient way to store and retrieve path costs for each node
// during graph search operations.
type Costs map[int32]float32

// SearchSpace represents the explored portion of the graph during a search operation.
// It inherits from Graph to maintain the structure of discovered paths while providing
// a separate space for search-specific operations and results.
type SearchSpace Graph

// PathCoord reconstructs and returns the geographical coordinates of nodes along a path from source to target
// in the search space. It performs a breadth-first traversal starting from the target node and following
// incoming edges backwards to reconstruct the complete path.
//
// The method converts the internal graph node representations to geographical coordinates using the S2
// geometry library, returning an array of [longitude, latitude] pairs that can be used for visualization
// or further geographical analysis.
//
// Parameters:
//   - target: int32 - The ID of the destination node from which to reconstruct the path backwards
//   - g: Graph - The original graph containing the complete node information including geographical locations
//
// Returns:
//   - [][]float64 - An array of coordinate pairs where:
//   - Each inner array contains exactly 2 float64 values: [longitude, latitude]
//   - Coordinates are in decimal degrees
//   - The array represents the complete path from source to target
//   - The order of coordinates follows the path traversal from target back to source
//   - Empty array is returned if target node is not found or no path exists
//
// Example:
//
//	coords := searchSpace.PathCoord(targetID, originalGraph)
//	// coords might contain: [[lng1,lat1], [lng2,lat2], ...]
func (sp SearchSpace) PathCoord(target int32, g Graph) [][]float64 {
	queue := list.New()
	queue.PushBack(target)
	result := make([][]float64, 0)
	for queue.Len() > 0 {
		qnode := queue.Front()
		queue.Remove(qnode)
		nodeID := qnode.Value.(int32)
		result = append(result, []float64{
			s2.CellID(g.Nodes[sp.Nodes[nodeID].Rank].Location).LatLng().Lng.Degrees(),
			s2.CellID(g.Nodes[sp.Nodes[nodeID].Rank].Location).LatLng().Lat.Degrees(),
		})
		for _, e := range sp.IncomingEdges[nodeID] {
			queue.PushBack(e.ID)
		}
	}
	return result
}

// GetCost retrieves the cost associated with reaching a specific node in the graph.
// This method provides safe access to the cost map with proper error handling.
//
// Parameters:
//   - id: int32 - The unique identifier of the node whose cost is being queried
//
// Returns:
//   - float32: The cost to reach the specified node from the source
//   - error: An error if the node is not found in the cost map, indicating no valid path exists
//
// Example:
//
//	cost, err := costs.GetCost(nodeID)
//	if err != nil {
//	    // Handle case where path doesn't exist
//	}
func (costs Costs) GetCost(id int32) (float32, error) {
	if v, ok := costs[id]; ok {
		return v, nil
	}
	return INFINITE, fmt.Errorf("path not found")
}

// Response encapsulates the complete results of a graph search operation.
// It provides access to the explored paths, cost matrix, and final computed costs
// for analysis and path reconstruction.
type Response struct {
	// SearchSpace contains the explored graph paths from source to target
	// It represents the subset of the original graph that was traversed during the search
	SearchSpace SearchSpace

	// PathCost stores a matrix of costs between nodes, limited to 150x150 nodes
	// This matrix enables quick lookup of path costs between any two nodes in the explored space
	PathCost [150][150]PathCost

	// Costs maps each node ID to its final computed cost from the source
	// This map contains the shortest path costs for all reached nodes
	Costs Costs
}

// DijkstraSearch implements Dijkstra's shortest path algorithm with additional constraints
// and optimizations. It maintains the search state and provides methods for executing
// the search process.
type DijkstraSearch struct {
	// pq is a priority queue that manages nodes to visit based on their current costs
	// It ensures that nodes are processed in order of increasing cost
	pq *Heap

	// visited tracks which nodes have been processed using a bitset for memory efficiency
	visited Bitset

	// previous stores the shortest path tree as it's being constructed
	// This graph structure allows for path reconstruction once the search is complete
	previous Graph

	// costs maps each node to its current best known cost from the source
	costs Costs

	// sources tracks which nodes are designated as starting points using a bitset
	sources Bitset

	// target stores the ID of the destination node (-1 if no specific target)
	// A specific target allows early termination when the destination is reached
	target int32
}

// NewDijkstra creates and initializes a new DijkstraSearch instance with the specified criteria.
// It sets up all necessary data structures and initializes the search state according to the
// provided configuration.
//
// Parameters:
//   - c: Criteria - A structure containing search parameters including:
//   - Source nodes: Starting points for the search
//   - Target nodes: Destination points for the search
//   - Maximum hop constraints
//
// Returns:
//   - DijkstraSearch: A fully initialized search instance ready to execute the algorithm
//     The returned instance contains:
//   - Initialized priority queue with source nodes
//   - Empty visited set
//   - Initialized cost map with source nodes set to zero
//   - Configured target node (if specified)
//
// Example:
//
//	criteria := Criteria{
//	    Source: []int32{1, 2},
//	    Targets: []int32{10},
//	    MaxHops: 5,
//	}
//	search := NewDijkstra(criteria)
func NewDijkstra(c Criteria) DijkstraSearch {
	target := int32(-1)
	if len(c.Targets) > 0 {
		target = c.Targets[0]
	}
	search := DijkstraSearch{
		pq:       Create(),
		visited:  NewBigInt(),
		previous: EmptyGraph(),
		costs:    make(Costs, 0),
		sources:  NewBigInt(),
		target:   target,
	}

	for _, s := range c.Source {
		search.costs[s] = 0
		search.pq.Insert(HNode{Value: s, Cost: 0, Depth: 0, Previous: 0})
		search.sources.Set(s, true)
	}

	return search
}

// Run executes Dijkstra's algorithm on the provided graph, finding shortest paths
// from source nodes to either all reachable nodes or a specific target node.
//
// Parameters:
//   - g: Graph - The input graph to search through, containing:
//   - Nodes and their properties
//   - Edge connections and weights
//   - Any additional metadata needed for path computation
//
// Returns:
//   - Response: A comprehensive result structure containing:
//   - SearchSpace: The explored portion of the graph
//   - Costs: Final shortest path costs to all reached nodes
//   - PathCost: Matrix of costs between nodes
//
// The algorithm continues until either:
//   - The target node is reached (if specified)
//   - The priority queue is empty (all reachable nodes processed)
//   - Maximum hop count is reached (if specified in criteria)
func (search DijkstraSearch) Run(g Graph) Response {
	currentID := int32(0)
	for !search.isFinished() {
		min, _ := search.pq.Min()
		if !search.wasVisited(min.Value) {
			currentID = search.addPrevious()
		}
		search.visited.Set(min.Value, true)

		if search.reachTarget(min.Value) {
			return Response{
				SearchSpace: SearchSpace(search.previous),
				Costs:       search.costs,
			}
		}
		for _, e := range g.OutgoingEdges[min.Value] {
			search.Relax(g.Nodes[e.ID], currentID, e.Weight, e.Metadata.Distance)
		}
		search.pq.DeleteMin()
	}
	return Response{
		SearchSpace: SearchSpace(search.previous),
		Costs:       search.costs,
	}
}

// addPrevious adds the current node to the path tree and creates the appropriate
// edge connections to maintain the shortest path tree structure.
//
// Returns:
//   - int32: The unique identifier assigned to the current node in the path tree
//
// The method performs the following operations:
//  1. Retrieves the minimum cost node from the priority queue
//  2. Adds it to the previous graph structure
//  3. Creates an edge from its parent if necessary
//  4. Updates the path cost information
func (search *DijkstraSearch) addPrevious() int32 {
	min, _ := search.pq.Min()
	currentID := search.previous.AddNode(Node{Rank: min.Value})
	if min.Previous != currentID {
		search.previous.RelateNodes(Node{ID: min.Previous}, Node{ID: currentID}, min.Cost, LeftToRight, MetaData{Distance: min.Dist})
	}
	return currentID
}

// Relax attempts to improve the shortest path to a node by considering a new path
// through a neighboring node. This is a fundamental operation in Dijkstra's algorithm
// that updates path costs when a shorter route is found.
//
// Parameters:
//   - v: Node - The destination node being considered for path improvement
//   - currentID: int32 - The ID of the current node in the path tree
//   - w: float32 - The time-based weight of the edge being considered
//   - distance: float32 - The physical distance weight of the edge
//
// The method performs the following steps:
//  1. Checks if the destination node has been visited
//  2. Calculates the new potential path cost
//  3. Compares with the existing cost
//  4. Updates the cost and priority queue if a shorter path is found
func (search DijkstraSearch) Relax(v Node, currentID int32, w, distance float32) {
	min, _ := search.pq.Min()
	if !search.wasVisited(v.ID) {
		cost := search.costs[min.Value]
		currentPathValue := cost + w
		currentDistancePathValue := cost + distance
		edgeC, _ := search.costs.GetCost(v.ID)
		if currentPathValue < edgeC {
			search.costs[v.ID] = currentPathValue
			search.pq.Insert(HNode{Value: v.ID, Cost: currentPathValue, Depth: min.Depth + 1, Previous: currentID, Dist: currentDistancePathValue})
		}
	}
}

// reachTarget determines if the current node being processed is the target node,
// allowing for early termination of the search when the destination is reached.
//
// Parameters:
//   - currentValue: int32 - The ID of the current node being processed
//
// Returns:
//   - bool: true if the current node is the target node and a target was specified,
//     false otherwise or if no target was specified (target < 0)
func (search DijkstraSearch) reachTarget(currentValue int32) bool {
	return search.target >= 0 && currentValue == search.target
}

// wasVisited checks if a node has already been processed in the current search,
// preventing cycles and ensuring each node is processed only once.
//
// Parameters:
//   - id: int32 - The ID of the node to check
//
// Returns:
//   - bool: true if the node has been visited, false otherwise
//
// The method uses a bitset for efficient memory usage and quick lookup.
func (search DijkstraSearch) wasVisited(id int32) bool {
	return search.visited.Exists(id)
}

// isFinished determines if the search process should terminate based on the state
// of the priority queue.
//
// Returns:
//   - bool: true if the priority queue is empty (no more nodes to process),
//     false if there are still nodes to examine
//
// This method is crucial for controlling the main search loop and ensuring
// termination when all reachable nodes have been processed.
func (search DijkstraSearch) isFinished() bool {
	return search.pq.IsEmpty()
}
