package utils

// UFNode is a Union-Find node on a Union-Find tree.
// Holds an integer as value. In this case we wish to store the variable ID.
type UFNode struct {
	// Rank is the length of the longest path from root to a leaf.
	rank int
	// Pa is parent of UFNode. Since we use the path-compression heuristic, this is usually the
	// representative of the set (i.e. the root node).
	pa *UFNode
	// Variable ID.
	varid int
}

// MakeSet creates a unary set with varid as representative of the resulting set.
func MakeSet(varid int) *UFNode {
	set := &UFNode{0, nil, varid}
	set.pa = set
	return set
}

// Find returns the representative of x's set.
func Find(x *UFNode) *UFNode {
	if x != x.pa {
		x.pa = Find(x.pa)
	}
	return x.pa
}

// Union takes the sets S_1 and S_2, where x is in S_1 and y is in S_2 and unifies S_1 with S_2,
// returning the union's representative.
// Returns the union set and 1 if the first argument is the new representative or 2 if the second.
func Union(x, y *UFNode) (*UFNode, int) {
	x, y = Find(x), Find(y)
	if x == y {
		return nil, -1
	}

	if x.rank > y.rank {
		y.pa = x
		return x, 1
	} else {
		x.pa = y
		if x.rank == y.rank {
			y.rank++
		}
		return y, 2
	}
}
