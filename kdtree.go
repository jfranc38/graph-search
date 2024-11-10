package graph_search

import (
	"math"
	"sort"
)

type KDTree struct {
	root *node
}

// node represents a node in the k-d tree.
type node struct {
	v Vector // The point stored in this node
	l *node  // Left child node
	r *node  // Right child node
}

// BuildKDTree constructs a KDTree from a slice of vectors.
// A KDTree is a data structure used for efficient multi-dimensional range searches and nearest neighbor queries.
// It is a binary tree where each level splits the data along a different dimension, cycling through all dimensions as the tree deepens.
// This structure allows for logarithmic time complexity for range queries in multiple dimensions.
//
// Parameters:
//   - vectors: A slice of Vector to build the tree from.
//
// Returns:
//   - A pointer to the constructed KDTree.
//
// Time Complexity: O(n log n), where n is the number of input vectors.
// Space Complexity: O(n), as the tree stores all input points.

func BuildKDTree(vectors []Vector) *KDTree {
	return &KDTree{
		root: build(vectors, 0),
	}
}

// build constructs a k-d tree recursively from a slice of vectors.
// It returns the root node of the constructed (sub)tree.
//
// Parameters:
//   - vectors: A slice of Vector to build the tree from.
//   - depth: The current depth in the tree, used to determine the splitting axis.
//
// Returns:
//   - A pointer to the root node of the constructed (sub)tree.
func build(vectors []Vector, depth int) *node {
	// Base case: if the input slice is empty, return nil (empty subtree)
	if len(vectors) == 0 {
		return nil
	}

	// Determine the number of dimensions (k) from the first vector
	k := len(vectors[0].Components)

	// Calculate the current axis to split on, cycling through dimensions
	// Example:
	// If we have 3D vectors (x, y, z) and the current depth is 5:
	// k = 3 (3 dimensions: x, y, z)
	// axis = 5 % 3 = 2 (corresponds to z-axis)
	//
	// This ensures that we cycle through dimensions as we go deeper in the tree:
	// depth 0: x-axis (0 % 3 = 0)
	// depth 1: y-axis (1 % 3 = 1)
	// depth 2: z-axis (2 % 3 = 2)
	// depth 3: x-axis (3 % 3 = 0)
	// and so on...
	axis := depth % k

	// Sort the vectors based on their component values in the current axis
	sort.Slice(vectors, func(i, j int) bool {
		return vectors[i].Components[axis] < vectors[j].Components[axis]
	})

	// Find the median point
	medianIndex := len(vectors) / 2
	medianPoint := vectors[medianIndex]

	// Construct and return the current node
	return &node{
		v: medianPoint,                             // Store the median point in this node
		l: build(vectors[:medianIndex], depth+1),   // Recursively build left subtree
		r: build(vectors[medianIndex+1:], depth+1), // Recursively build right subtree
	}
}

// Query performs a range search on the KDTree to find all points within a given radius of a center point.
// It calls the internal query function starting from the root of the tree.
//
// Parameters:
//   - center: The center point of the search range.
//   - radius: The radius of the search range.
//
// Returns:
//   - A slice of Vector objects representing all points within the specified range.
func (t *KDTree) RangeQuery(center Vector, radius float64) []Vector {
	return rangeQuery(t.root, center, radius, 0)
}

// squaredDistance calculates the squared Euclidean distance between two vectors.
// This is more efficient than calculating the actual distance as it avoids the square root operation.
//
// Parameters:
//   - u: The first
//   - v: The second
//
// Returns:
//   - The squared Euclidean distance between u and v.
func squaredDistance(u, v Vector) float64 {
	sum := 0.0
	for i := range u.Components {
		diff := u.Components[i] - v.Components[i]
		sum += diff * diff
	}
	return sum
}

// rangeQuery performs a range search on the k-d tree to find all points within a given radius of a center point.
// It recursively traverses the tree, pruning branches that cannot contain points within the specified range.
//
// Parameters:
//   - node: The current node in the k-d tree being examined.
//   - center: The center point of the search range.
//   - radius: The radius of the search range.
//   - depth: The current depth in the tree, used to determine the splitting axis.
//
// Returns:
//   - A slice of Vector objects representing all points within the specified range.
//
// Example:
//
//	Suppose we have a 2D k-d tree with points: (1,2), (3,4), (5,6), (7,8)
//	We want to find all points within radius 2 of the center point (4,5)
//
//	1. Start at the root node (3,4)
//	2. Check if (3,4) is within radius 2 of (4,5) - it is, so add it to the result
//	3. Check if we need to search left subtree (we do, as 4-2 <= 3)
//	4. Check if we need to search right subtree (we do, as 4+2 >= 3)
//	5. Recursively search both subtrees
//	6. In the end, return [(3,4), (5,6)] as the result
func rangeQuery(node *node, center Vector, radius float64, depth int) []Vector {
	// Base case: if the node is nil, return an empty slice
	if node == nil {
		return nil
	}

	// Determine the number of dimensions and current axis
	k := len(node.v.Components)
	axis := depth % k

	// Initialize a slice to store points within the range
	pointsInRange := []Vector{}

	// Check if the current node's point is within the search radius
	if squaredDistance(node.v, center) <= radius*radius {
		pointsInRange = append(pointsInRange, node.v)
	}

	// Determine whether to search the left and/or right subtrees
	// Left subtree: search if the hypersphere's left bound is less than or equal to the current node's splitting value
	if node.l != nil && center.Components[axis]-radius <= node.v.Components[axis] {
		pointsInRange = append(pointsInRange, rangeQuery(node.l, center, radius, depth+1)...)
	}
	// Right subtree: search if the hypersphere's right bound is greater than or equal to the current node's splitting value
	if node.r != nil && center.Components[axis]+radius >= node.v.Components[axis] {
		pointsInRange = append(pointsInRange, rangeQuery(node.r, center, radius, depth+1)...)
	}

	return pointsInRange
}

// FindNearest finds the nearest neighbor to a target point in the KDTree.
//
// This method takes a target vector and returns the closest vector in the tree
// along with its distance from the target.
//
// Parameters:
//   - target: The target vector for which we're finding the nearest neighbor.
//
// Returns:
//   - The nearest vector found in the tree.
//   - The Euclidean distance between the target and the nearest
func (t *KDTree) FindNearest(target Vector) (Vector, float64) {
	best, bestDist := nearest(t.root, target, 0, nil, math.MaxFloat64)
	return best.v, bestDist
}

// nearest finds the nearest neighbor to a target point in the k-d tree.
//
// This function recursively traverses the k-d tree to find the node that is closest to the target
// It uses a depth-first search strategy and prunes branches that cannot contain a closer point.
//
// Parameters:
//   - n: The current node in the k-d tree.
//   - target: The target vector for which we're finding the nearest neighbor.
//   - depth: The current depth in the tree, used to determine the splitting axis.
//   - best: The current best (closest) node found so far.
//   - bestDist: The squared distance to the current best node.
//
// Returns:
//   - A pointer to the nearest node found.
//   - The squared distance to the nearest node.
//
// Example:
//
//	Suppose we have a 2D k-d tree with points: (2,3), (5,4), (9,6), (4,7), (8,1), (7,2)
//	We want to find the nearest neighbor to the target point (6,5)
//
//	1. Start at the root node (5,4)
//	2. Compare distance: (6-5)^2 + (5-4)^2 = 2, update best and bestDist
//	3. Target's x (6) > node's x (5), so search right subtree first
//	4. Move to (9,6), compare distance: (6-9)^2 + (5-6)^2 = 10, don't update best
//	5. No more right children, backtrack and check if left subtree needs searching
//	6. It does, so move to (7,2), compare distance: (6-7)^2 + (5-2)^2 = 10, don't update best
//	7. Continue this process for remaining nodes
//	8. In the end, return (5,4) as the nearest neighbor with distance 2
func nearest(n *node, target Vector, depth int, best *node, bestDist float64) (*node, float64) {
	if n == nil {
		return best, bestDist
	}
	k := len(target.Components)
	axis := depth % k

	// Calculate the distance from the target to the current node
	dist := squaredDistance(n.v, target)
	if dist < bestDist {
		bestDist = dist
		best = n
	}

	// Determine which subtree to search first
	var next, other *node

	if target.Components[axis] < n.v.Components[axis] {
		next = n.l
		other = n.r
	} else {
		next = n.r
		other = n.l
	}

	// Recursively search the next subtree
	best, bestDist = nearest(next, target, depth+1, best, bestDist)

	// Check if we need to search the other subtree
	if math.Abs(n.v.Components[axis]-target.Components[axis]) < math.Sqrt(bestDist) {
		best, bestDist = nearest(other, target, depth+1, best, bestDist)
	}

	return best, bestDist
}
