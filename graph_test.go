package graph_search

import (
	"fmt"
	"testing"

	"github.com/paulmach/go.geojson"
)

func TestGraphSearch(t *testing.T) {
	graph := BuildGraph("testdata/colombia-latest.osm.pbf")
	fmt.Println(len(graph.Nodes))

	rangeTree := graph.BuildNodeIndex()

	source := Coordinate{Lat: 6.1997796925416395, Lng: -75.57815231451204}
	target := Coordinate{Lat: 6.197606519075109, Lng: -75.55768012592779}

	sourceX, sourceY := LatLngToMeters(source.Lat, source.Lng)
	targetX, targetY := LatLngToMeters(target.Lat, target.Lng)

	projectedSource, _ := rangeTree.FindNearest(Vector{Components: []float64{sourceX, sourceY}})
	projectedTarget, _ := rangeTree.FindNearest(Vector{Components: []float64{targetX, targetY}})

	response := NewDijkstra(Criteria{
		Source:  []int32{int32(projectedSource.ID)},
		Targets: []int32{int32(projectedTarget.ID)},
	}).Run(graph)
	distance, _ := response.Costs.GetCost(int32(projectedTarget.ID))
	targetSearchSpace := response.SearchSpace.Nodes[len(response.SearchSpace.Nodes)-1].ID
	p := response.SearchSpace.PathCoord(targetSearchSpace, graph)

	fc := geojson.NewFeatureCollection()
	fc.AddFeature(geojson.NewLineStringFeature(p))
	Write("testdata/route.geojson", fc)
	fmt.Printf("Total distance: %.2f meters\n", distance)
}
