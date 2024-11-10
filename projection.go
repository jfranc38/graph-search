package graph_search

import "math"

// LatLngToMeters converts latitude and longitude to X and Y coordinates in meters.
//
// Parameters:
//   - lat: float64 - The latitude in degrees.
//   - lng: float64 - The longitude in degrees.
//
// Returns:
//   - x: float64 - The X coordinate in meters.
//   - y: float64 - The Y coordinate in meters.
func LatLngToMeters(lat, lng float64) (x, y float64) {
	R := 6378137.0
	φ := lat * (math.Pi / 180.0)
	λ := lng * (math.Pi / 180.0)
	x = R * λ
	y = R * math.Log(math.Tan((math.Pi/4)+(φ/2)))
	return x, y
}

// MetersToLatLng converts X and Y coordinates in meters to latitude and longitude.
//
// Parameters:
//   - x: float64 - The X coordinate in meters.
//   - y: float64 - The Y coordinate in meters.
//
// Returns:
//   - lat: float64 - The latitude in degrees.
//   - lng: float64 - The longitude in degrees.
func MetersToLatLng(x, y float64) (lat, lng float64) {
	R := 6378137.0
	λ := x / R
	φ := 2*math.Atan(math.Exp(y/R)) - (math.Pi / 2)
	lat = φ * (180.0 / math.Pi)
	lng = λ * (180.0 / math.Pi)
	return lat, lng
}
