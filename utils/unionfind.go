package utils

// UFNode is a Union-Find node on a Union-Find tree.
// Holds an integer as value. In this case we wish to store the variable ID.
type UFNode struct {
	// Rank is the length of the longest.Path from root to a leaf.
	rank int
	//.Pa is.Parent of UFNode. Since we use the.Path-compression heuristic, this is usually the
	// representative of the set (i.e. the root node).
	Pa *UFNode
	// Variable ID.
	varid int
	// Children.
	Ch []*UFNode
}

// MakeSet creates a unary set with varid as representative of the resulting set.
func MakeSet(varid int) *UFNode {
	set := &UFNode{0, nil, varid, nil}
	set.Pa = set
	return set
}

// Find returns the representative of x's set.
func Find(x *UFNode) *UFNode {
	if x != x.Pa {
		x.Pa = Find(x.Pa)
	}
	return x.Pa
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
		y.Pa = x
		x.Ch = append(x.Ch, y)
		return x, 1
	}
	x.Pa = y
	y.Ch = append(y.Ch, x)
	if x.rank == y.rank {
		y.rank++
	}
	return y, 2
}

// UFVarids returns a slice with all varids in union-find tree x.
func UFVarids(x *UFNode) []int {
	n := len(x.Ch)

	if n == 0 {
		return []int{x.varid}
	}

	ch := []int{x.varid}

	for i := 0; i < n; i++ {
		ch = append(ch, UFVarids(x.Ch[i])...)
	}

	return ch
}
