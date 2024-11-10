// Package graph_search provides functionality for creating and manipulating graph structures from OpenStreetMap (OSM)
// Protocolbuffer Binary Format (PBF) files. It includes tools for parsing OSM data, building directed weighted graphs,
// and performing spatial operations on geographical data.
package graph_search

import (
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/golang/geo/s2"
	"github.com/qedus/osmpbf"
)

// BuildGraph constructs a graph from an OSM PBF file, processing nodes and ways to create a connected road network.
// It filters ways based on road type tags and builds edges between connected nodes.
//
// Parameters:
//   - path: string - File path to the OSM PBF file to process
//
// Returns:
//   - Graph: A constructed graph containing nodes and edges representing the road network
//     The graph includes:
//   - Nodes with geographical coordinates stored as S2 cell IDs
//   - Edges with weights based on travel time/distance
//   - Metadata including speed limits, distances, and road types
func BuildGraph(path string) Graph {
	decoder, file := openAndDecodePBF(path)
	nodes := buildCoverageNodes(path)
	ways := make(map[int64][]int32)
	g := Graph{Nodes: make([]Node, 0, len(nodes))}

	for {
		obj, err := decoder.Decode()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		switch obj := obj.(type) {
		case *osmpbf.Node:
			buildNode(&g, obj, nodes)
		case *osmpbf.Way:
			if validWay(*obj) {
				buildWay(&g, obj, nodes, ways)
			}
		}
	}

	_ = file.Close()
	nodes = nil
	return g
}

// buildNode creates and adds a node to the graph based on OSM node data.
// The node is only added if its OSM ID exists in the provided nodes map.
//
// Parameters:
//   - g: *Graph - Pointer to the graph being constructed
//   - node: *osmpbf.Node - OSM node data containing location information
//   - nodes: map[int64]int32 - Map of valid OSM node IDs to internal graph IDs
//
// The function modifies the graph by adding nodes and updates the nodes map with internal IDs
func buildNode(g *Graph, node *osmpbf.Node, nodes map[int64]int32) {
	osmID := node.ID
	if _, ok := nodes[osmID]; ok {
		id := g.AddNode(Node{
			Location: coordinatesToCellID(node.Lat, node.Lon),
		})
		nodes[osmID] = id
	}
}

// buildWay creates edges in the graph based on OSM way data. It processes sequences of nodes
// that form a way, calculating distances and travel times between consecutive nodes.
//
// Parameters:
//   - g: *Graph - Pointer to the graph being constructed
//   - way: *osmpbf.Way - OSM way data containing node sequences and tags
//   - nodes: map[int64]int32 - Map of valid node IDs
//   - ways: map[int64][]int32 - Map to store processed way segments
//
// The function modifies the graph by:
//   - Adding edges between consecutive nodes in the way
//   - Setting edge weights based on distance and speed limits
//   - Including metadata about road type and travel characteristics
func buildWay(g *Graph, way *osmpbf.Way, nodes map[int64]int32, ways map[int64][]int32) {
	speed := 50 // Default speed in km/h
	for i := 0; i < len(way.NodeIDs)-1; i++ {
		idA, ok1 := nodes[way.NodeIDs[i]]
		idB, ok2 := nodes[way.NodeIDs[i+1]]

		if !ok1 || !ok2 {
			continue
		}

		nodeA := g.Nodes[idA]
		nodeB := g.Nodes[idB]
		distance := DistanceMeters(s2.CellID(nodeA.Location), s2.CellID(nodeB.Location))
		roadType := "n/a"
		if highwayTag, found := way.Tags[Highway]; found {
			roadType = strings.ToLower(highwayTag)
		}
		g.RelateNodes(nodeA, nodeB, distance, edgeDirectionFromWay(*way), MetaData{
			Speed:    float32(speed),
			Distance: distance,
			RoadType: roadType,
		})
		ways[way.ID] = append(ways[way.ID], nodeA.ID)
		if i == len(way.NodeIDs)-2 {
			ways[way.ID] = append(ways[way.ID], nodeB.ID)
		}
	}
}

// buildCoverageNodes creates a map of valid nodes from the input file.
// It processes the file to identify nodes that are part of valid road segments.
//
// Parameters:
//   - path: string - Path to the OSM PBF file to process
//
// Returns:
//   - map[int64]int32: A map where keys are OSM node IDs and values are internal graph node IDs
func buildCoverageNodes(path string) map[int64]int32 {
	nodes := determineValidNodesFromFile(path)
	log.Println("Valid nodes from file: ", len(nodes))

	return nodes
}

// calculateTimeAndDistance computes travel time and physical distance between two geographical points.
//
// Parameters:
//   - origin: s2.CellID - S2 cell ID of the starting point
//   - target: s2.CellID - S2 cell ID of the ending point
//   - velocityKMH: float64 - Travel speed in kilometers per hour
//
// Returns:
//   - float32: Travel time in minutes
//   - float32: Distance in meters
func calculateTimeAndDistance(origin, target s2.CellID, velocityKMH float64) (float32, float32) {
	distanceM := DistanceMeters(origin, target)
	distanceKM := float64(distanceM / MetersInAKilometer)

	timeMinutes := (distanceKM / velocityKMH) * MinutesInAnHour
	return float32(timeMinutes), distanceM
}

// determineValidNodesFromFile processes an OSM PBF file to identify nodes that are part of valid ways.
//
// Parameters:
//   - path: string - Path to the OSM PBF file
//
// Returns:
//   - map[int64]int32: Map of valid OSM node IDs to sequential internal IDs
//
// The function filters nodes based on their presence in valid ways (roads, paths, etc.)
func determineValidNodesFromFile(path string) map[int64]int32 {
	d, f := openAndDecodePBF(path)

	result := make(map[int64]int32)
	i := 0
	for {
		if o, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch o := o.(type) {
			case *osmpbf.Way:
				w := *o
				if validWay(w) {
					for _, n := range w.NodeIDs {
						if _, ok := result[n]; !ok {
							result[n] = int32(i)
							i++
						}
					}
				}
			}
		}
	}
	_ = f.Close()
	return result
}

// coordinatesToCellID converts latitude and longitude coordinates to an S2 cell ID.
//
// Parameters:
//   - lat: float64 - Latitude in degrees (-90 to +90)
//   - lng: float64 - Longitude in degrees (-180 to +180)
//
// Returns:
//   - uint64: S2 cell ID at the configured cell level (CellLevel)
func coordinatesToCellID(lat, lng float64) uint64 {
	return uint64(s2.CellFromPoint(s2.PointFromLatLng(
		s2.LatLngFromDegrees(lat, lng))).ID().Parent(CellLevel))
}

// validWay determines if an OSM way represents a valid road segment for inclusion in the graph.
//
// Parameters:
//   - w: osmpbf.Way - OSM way to validate
//
// Returns:
//   - bool: true if the way represents a valid road type, false otherwise
//
// Valid road types include: motorway, trunk, primary, secondary, tertiary, residential, and their variants
func validWay(w osmpbf.Way) bool {
	tags := map[string]struct{}{
		Motorway: {}, MotorwayLink: {}, Trunk: {},
		TrunkLink: {}, Primary: {}, PrimaryLink: {},
		Secondary: {}, SecondaryLink: {}, Tertiary: {},
		TertiaryLink: {}, Residential: {},
		Unclassified: {}, LivingStreet: {},
	}

	_, ok := tags[(w.Tags)[Highway]]
	return ok
}

// edgeDirectionFromWay determines the directionality of a road segment based on OSM tags.
//
// Parameters:
//   - w: osmpbf.Way - OSM way to analyze
//
// Returns:
//   - EdgeDirection: One of:
//   - LeftToRight: One-way from start to end
//   - Bidirectional: Two-way traffic allowed
//
// The direction is determined by oneway tags and special cases like roundabouts
func edgeDirectionFromWay(w osmpbf.Way) EdgeDirection {
	tags := w.Tags
	if oneWay, ok := tags[Oneway]; ok && oneWay == Yes {
		return LeftToRight
	}
	if junction, ok := tags[Junction]; ok && junction == Roundabout {
		return LeftToRight
	}
	return Bidirectional
}

// openAndDecodePBF opens an OSM PBF file and creates an optimized decoder for processing.
//
// Parameters:
//   - path: string - Path to the OSM PBF file
//
// Returns:
//   - *osmpbf.Decoder: Configured PBF decoder
//   - *os.File: Open file handle
//
// The function configures the decoder for optimal performance using maximum buffer size
// and parallel processing based on available CPU cores
func openAndDecodePBF(path string) (*osmpbf.Decoder, *os.File) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	d := osmpbf.NewDecoder(f)
	d.SetBufferSize(osmpbf.MaxBlobSize)
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	return d, f
}
