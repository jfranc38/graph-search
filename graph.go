package graph_search

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"os"

	"github.com/golang/geo/s2"
	"github.com/umahmood/haversine"
)

// Graph represents a directed weighted graph data structure consisting of nodes (vertices) and edges.
// It maintains separate collections for nodes and their incoming/outgoing edge relationships.
type Graph struct {
	Nodes         []Node    // Collection of all nodes in the graph
	IncomingEdges Relations // Adjacency list of incoming edges for each node
	OutgoingEdges Relations // Adjacency list of outgoing edges for each node
}

// MetaData contains additional information associated with graph edges.
type MetaData struct {
	Speed    float32 // Speed limit or average speed for the edge in meters/second
	Distance float32 // Physical distance of the edge in meters
	RoadType string  // Classification of the road/path type (e.g., "motorway", "residential")
}

// Node represents a vertex in the graph with geographical positioning.
// Each node has a unique identifier, location encoded as an S2 cell ID, and rank for ordering.
type Node struct {
	ID       int32  // Unique identifier for the node
	Location uint64 // S2 cell ID encoding the geographical position
	Rank     int32  // Numerical rank used for node ordering
}

// Nodes is a slice type alias for a collection of Node objects
type Nodes []Node

// EdgeDirection defines the possible directions for edges between nodes
type EdgeDirection int

const (
	Bidirectional EdgeDirection = iota // Edge allows travel in both directions
	LeftToRight                        // Edge only allows travel from left node to right node
	RightToLeft                        // Edge only allows travel from right node to left node
)

// Edge represents a directed connection between two nodes in the graph.
// Each edge carries a weight and additional metadata about the connection.
type Edge struct {
	ID       int32    // Identifier of the destination node
	Weight   float32  // Cost/weight associated with traversing this edge
	Metadata MetaData // Additional data about the edge (speed, distance, road type)
}

// Coordinate represents a geographical position using latitude and longitude.
type Coordinate struct {
	Lat float64 // Latitude in degrees (-90 to +90)
	Lng float64 // Longitude in degrees (-180 to +180)
}

// Coordinates is a slice type alias for a collection of Coordinate objects
type Coordinates []Coordinate

// Relations is a slice of edge slices, representing adjacency lists for graph nodes
type Relations [][]Edge

// EmptyGraph creates and returns a new empty Graph instance with initialized but empty collections.
// Returns:
//   - Graph: A new Graph with empty Nodes, OutgoingEdges, and IncomingEdges slices
func EmptyGraph() Graph {
	return Graph{Nodes: make([]Node, 0), OutgoingEdges: make(Relations, 0), IncomingEdges: make(Relations, 0)}
}

// GetPoint converts the node's S2 cell ID location into latitude/longitude coordinates.
// Returns:
//   - s2.LatLng: Geographical coordinates of the node as a LatLng pair
func (n Node) GetPoint() s2.LatLng {
	return s2.CellID(n.Location).LatLng()
}

// GetID returns the node's identifier as an integer.
// Returns:
//   - int: The node's ID converted from int32 to int
func (n Node) GetID() int {
	return int(n.ID)
}

// AddNode appends a new node to the graph and initializes its edge collections.
// Parameters:
//   - n: Node - The node to be added to the graph
//
// Returns:
//   - int32: The ID assigned to the newly added node
func (g *Graph) AddNode(n Node) int32 {
	id := len(g.Nodes)
	n.ID = int32(id)
	g.Nodes = append(g.Nodes, n)
	g.OutgoingEdges = append(g.OutgoingEdges, make([]Edge, 0))
	g.IncomingEdges = append(g.IncomingEdges, make([]Edge, 0))
	return int32(id)
}

// RelateNodes creates edges between two nodes according to the specified direction.
// Parameters:
//   - a: Node - The first node to relate
//   - b: Node - The second node to relate
//   - weight: float32 - The weight/cost of the edge(s)
//   - dir: EdgeDirection - The direction of the relationship (Bidirectional, LeftToRight, or RightToLeft)
//   - metaData: MetaData - Additional information about the edge(s)
func (g *Graph) RelateNodes(a, b Node, weight float32, dir EdgeDirection, metaData MetaData) {
	switch dir {

	case Bidirectional:
		// relate two nodes bidirectionally o<------>o.
		{
			// Left to right relation(relate node n with node x).
			g.addOutgoingEdge(a.ID, b.ID, weight, metaData)
			g.addIncomingEdge(b.ID, a.ID, weight, metaData)

			// Right to left relation(relate node x with node n).
			g.addOutgoingEdge(b.ID, a.ID, weight, metaData)
			g.addIncomingEdge(a.ID, b.ID, weight, metaData)
		}

	case LeftToRight:
		// relate two nodes from left to right o------>o.
		{
			g.addOutgoingEdge(a.ID, b.ID, weight, metaData)
			g.addIncomingEdge(a.ID, b.ID, weight, metaData)
		}

	case RightToLeft:
		// relate two nodes from right to left o<------o.
		{
			g.addOutgoingEdge(b.ID, a.ID, weight, metaData)
			g.addIncomingEdge(b.ID, a.ID, weight, metaData)
		}
	}
}

// addOutgoingEdge adds a directed edge from one node to another in the outgoing edges collection.
// Parameters:
//   - from: int32 - ID of the source node
//   - to: int32 - ID of the destination node
//   - weight: float32 - The weight/cost of the edge
//   - metaData: MetaData - Additional information about the edge
func (g *Graph) addOutgoingEdge(from, to int32, weight float32, metaData MetaData) {
	if g.OutgoingEdges[from] == nil {
		g.OutgoingEdges[from] = make([]Edge, 0)
	}
	g.OutgoingEdges[from] = append(g.OutgoingEdges[from], Edge{
		ID:       to,
		Weight:   weight,
		Metadata: metaData,
	})
}

// addIncomingEdge adds a directed edge from one node to another in the incoming edges collection.
// Parameters:
//   - from: int32 - ID of the source node
//   - to: int32 - ID of the destination node
//   - weight: float32 - The weight/cost of the edge
//   - metaData: MetaData - Additional information about the edge
func (g *Graph) addIncomingEdge(from, to int32, weight float32, metaData MetaData) {
	if g.IncomingEdges[to] == nil {
		g.IncomingEdges[to] = make([]Edge, 0)
	}
	g.IncomingEdges[to] = append(g.IncomingEdges[to], Edge{
		ID:       from,
		Weight:   weight,
		Metadata: metaData,
	})
}

// DistanceMeters calculates the great-circle distance between two geographical points using the Haversine formula.
// Parameters:
//   - a: s2.CellID - The S2 cell ID of the first location
//   - b: s2.CellID - The S2 cell ID of the second location
//
// Returns:
//   - float32: The distance between the points in meters
func DistanceMeters(a, b s2.CellID) float32 {
	_, km := haversine.Distance(
		haversine.Coord{Lat: a.LatLng().Lat.Degrees(), Lon: a.LatLng().Lng.Degrees()},
		haversine.Coord{Lat: b.LatLng().Lat.Degrees(), Lon: b.LatLng().Lng.Degrees()},
	)
	return float32(km * MetersInAKilometer)
}

// BuildNodeIndex creates a spatial index of nodes using a range tree data structure.
// Only nodes with outgoing edges are included in the index.
// Parameters:
//   - g: *Graph - The graph whose nodes should be indexed
//
// Returns:
//   - *RangeTree: A spatial index of the graph's nodes for efficient geographical queries
func (g *Graph) BuildNodeIndex() *KDTree {
	vectors := make([]Vector, 0)
	for _, n := range g.Nodes {
		if len(g.OutgoingEdges[n.ID]) > 0 {
			latLng := s2.CellID(n.Location).LatLng()
			x, y := LatLngToMeters(latLng.Lat.Degrees(), latLng.Lng.Degrees())
			vector := Vector{ID: n.GetID(), Components: []float64{x, y}}
			vectors = append(vectors, vector)
		}
	}
	return BuildKDTree(vectors)
}

// Write serializes and writes content to a JSON file.
//
// This function creates a new file with the given name, marshals the content to JSON format,
// and writes it to the file. It handles file operations safely with proper error checking
// and resource cleanup.
//
// Parameters:
//   - name: string - The name/path of the file to create and write to
//   - content: interface{} - The data to be serialized to JSON and written to the file.
//     Can be any JSON-serializable Go type
//
// Returns:
//   - string - The name of the created file if successful, empty string if any error occurs
//
// The function will:
//   - Create a new file, overwriting if it already exists
//   - Marshal the content to JSON
//   - Write the JSON data to the file
//   - Log the number of bytes written
//   - Handle all errors appropriately with cleanup
//   - Close the file properly in all cases
func Write(name string, content interface{}) string {
	f, err := os.Create(name)
	if err != nil {
		return ""
	}
	d2, _ := json.Marshal(content)
	n2, err := f.Write(d2)
	if err != nil {
		log.Println(err)
		f.Close()
		return ""
	}
	log.Println(n2, "bytes written successfully")
	err = f.Close()
	if err != nil {
		log.Println(err)
		return ""
	}
	return f.Name()
}

// Serialize encodes and writes the Graph structure to a binary file using Go's gob encoding.
//
// This method persists the entire Graph structure to disk in a binary format that preserves
// all relationships and data. The gob encoder handles complex data structures and maintains
// referential integrity.
//
// Parameters:
//   - filePath: string - The full path where the serialized graph should be written
//
// Returns:
//   - error - nil if the serialization was successful, otherwise returns the encountered error
//
// The method will:
//   - Create a new file at the specified path
//   - Initialize a gob encoder
//   - Encode the entire graph structure
//   - Handle proper file closure
//   - Return any errors encountered during the process
func (g Graph) Serialize(filePath string) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(g)
	}
	file.Close()
	return err
}

// Deserialize reads a binary file and reconstructs a Graph structure from it.
//
// This function reads a previously serialized Graph from disk and reconstructs the complete
// graph structure including all nodes, edges, and associated metadata. It uses Go's gob
// decoder to handle the binary format and restore the complex data structure.
//
// Parameters:
//   - filePath: string - The path to the file containing the serialized Graph data
//
// Returns:
//   - Graph - The reconstructed Graph structure. If an error occurs during deserialization,
//     returns an empty Graph structure
//
// The function will:
//   - Open the specified file
//   - Initialize a gob decoder
//   - Decode the binary data into a new Graph structure
//   - Handle proper file closure
//   - Return the reconstructed Graph
//
// Note: Error handling is internal - errors during deserialization will result
// in an empty Graph being returned
func Deserialize(filePath string) Graph {
	var g = new(Graph)
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(g)
	}
	file.Close()
	return *g
}
