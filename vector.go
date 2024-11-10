package graph_search

import "math"

// Vector represents a mathematical vector in n-dimensional space.
// It contains a slice of float64 values representing the vector's components.
type Vector struct {
	ID         int
	Components []float64
}

func NewVector(id int, components []float64) Vector {
	return Vector{ID: id, Components: components}
}

// Add performs vector addition of two Vector objects.
// It returns a new Vector that is the sum of the current vector and the other vector.
//
// Parameters:
//   - other: The Vector to be added to the current Vector.
//
// Returns:
//   - A new Vector representing the sum of the two vectors.
//
// Panics:
//   - If the dimensions of the two vectors are not equal.
func (v Vector) Add(other Vector) Vector {
	// Check if the vectors have the same dimension
	if len(v.Components) != len(other.Components) {
		panic("Vectors must be of the same dimension to add")
	}

	// Create a new Vector to store the result
	result := Vector{
		Components: make([]float64, len(v.Components)),
	}

	// Iterate through the components and add them
	for i := range v.Components {
		result.Components[i] = v.Components[i] + other.Components[i]
	}

	// Return the resulting vector
	return result
}

// Subtract performs vector subtraction of two Vector objects.
// It returns a new Vector that is the difference between the current vector and the other vector.
//
// Parameters:
//   - other: The Vector to be subtracted from the current Vector.
//
// Returns:
//   - A new Vector representing the difference between the two vectors.
func (v Vector) Subtract(other Vector) Vector {
	return v.Add(other.Scale(-1))
}

// Scale multiplies the vector by a scalar value.
// It returns a new Vector with each component multiplied by the given scalar.
//
// Parameters:
//   - scalar: The float64 value to multiply each component by.
//
// Returns:
//   - A new Vector representing the scaled vector.
func (v Vector) Scale(scalar float64) Vector {
	result := Vector{
		Components: make([]float64, len(v.Components)),
	}
	for i, component := range v.Components {
		result.Components[i] = component * scalar
	}
	return result
}

// Dot calculates the dot product of two Vector objects.
// It returns a float64 representing the scalar result of the dot product.
//
// Parameters:
//   - other: The Vector to calculate the dot product with.
//
// Returns:
//   - A float64 value representing the dot product of the two vectors.
//
// Panics:
//   - If the dimensions of the two vectors are not equal.
func (v Vector) Dot(other Vector) float64 {
	result := 0.0
	for i, component := range v.Components {
		result += component * other.Components[i]
	}
	return result
}

// Magnitude calculates the length (magnitude) of the vector.
// It returns a float64 representing the Euclidean norm of the vector.
//
// Returns:
//   - A float64 value representing the magnitude of the vector.
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.Dot(v))
}

// Normalize returns a new Vector with the same direction as the original vector but with a magnitude of 1.
// If the original vector has a magnitude of 0, it returns a zero vector to avoid division by zero.
//
// Returns:
//   - A new Vector representing the normalized vector.
func (v Vector) Normalize() Vector {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector{Components: make([]float64, len(v.Components))}
	}
	return v.Scale(1 / mag)
}

// Unit is an alias for Normalize. It returns a new Vector with the same direction as the original vector but with a magnitude of 1.
//
// Returns:
//   - A new Vector representing the unit vector.
func (v Vector) Unit() Vector {
	return v.Normalize()
}

// Project calculates the vector projection of the current vector onto another vector.
// It returns a new Vector that represents the projection.
//
// Parameters:
//   - other: The Vector to project onto.
//
// Returns:
//   - A new Vector representing the projection of the current vector onto the other vector.
func (v Vector) Project(other Vector) Vector {
	u := other.Unit()
	return u.Scale(v.Dot(u))
}

// Between checks if the projection of vector x onto v falls between the origin and v.
// It returns true if the dot product of v and x is positive and less than the dot product of v with itself.
//
// Parameters:
//   - x: The Vector to check.
//
// Returns:
//   - A boolean indicating whether x projects between the origin and v.
func (v Vector) Between(x Vector) bool {
	return v.Dot(x) > 0 && v.Dot(x) < v.Dot(v)
}

// Copy returns a new Vector that is an exact copy of the original Vector.
// This method creates a deep copy, ensuring that modifications to the new Vector
// do not affect the original.
//
// Returns:
//   - A new Vector with the same components as the original Vector.
func (v Vector) Copy() Vector {
	newComponents := make([]float64, len(v.Components))
	copy(newComponents, v.Components)
	return Vector{Components: newComponents}
}

// Distance calculates the Euclidean distance between two vectors.
// This method computes the square root of the sum of the squared differences
// between corresponding components of the two vectors.
//
// Parameters:
//   - other: The Vector to calculate the distance to.
//
// Returns:
//   - float64: The Euclidean distance between the two vectors.
//
// Panics:
//   - If the vectors have different dimensions.
func (v Vector) Distance(other Vector) float64 {
	sumSquares := 0.0
	for i := 0; i < len(v.Components); i++ {
		diff := v.Components[i] - other.Components[i]
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares)
}

// DistanceSquared calculates the squared Euclidean distance between two vectors.
// This method computes the sum of the squared differences between corresponding
// components of the two vectors.
//
// Parameters:
//   - other: The Vector to calculate the squared distance to.
//
// Returns:
//   - float64: The squared Euclidean distance between the two vectors.
func (v Vector) DistanceSquared(other Vector) float64 {
	sumSquares := 0.0
	for i := 0; i < len(v.Components); i++ {
		diff := v.Components[i] - other.Components[i]
		sumSquares += diff * diff
	}
	return sumSquares
}

// IsZero checks if the vector is a zero vector (all components are zero).
//
// Returns:
//   - bool: true if the vector is a zero vector, false otherwise.
func (v Vector) IsZero() bool {
	for _, comp := range v.Components {
		if comp != 0 {
			return false
		}
	}
	return true
}

// Equals checks if two vectors are equal (have the same components).
//
// Parameters:
//   - other: The Vector to compare with.
//
// Returns:
//   - bool: true if the vectors are equal, false otherwise.
func (v Vector) Equals(other Vector) bool {
	if len(v.Components) != len(other.Components) {
		return false
	}
	for i, comp := range v.Components {
		if comp != other.Components[i] {
			return false
		}
	}
	return true
}
