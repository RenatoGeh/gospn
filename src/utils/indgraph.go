package utils

import (
	_ "fmt"
)

/*
IndepGraph represents an independence graph.

An independence graph is an undirected graph that maps the (in)dependencies of a set of variable.
Let X={X_1,...,X_n} be the set of variables. We define an independence graph as an undirected
graph G=(X, E) where there exists an edge between a pair of vertices u,v in X iff there exists a
dependency between variables u and v. That is, if two variables are dependent than there exists
an edge between them. Otherwise there is no such edge.

The resulting graph after such construction is a graph with clusters of connected graphs. Let
H_1 and H_2 be two complete subgraphs in G. Then there exists no edge between any one vertex in
H_1 and another in H_2. This constitues an independency relation between these subgraphs. Thus we
say that sets of variables in H_1 are independent of sets of variables in H_2. We now show why this
is correct. Consider the following example (it can be extended to the general case easily):

Let X, Y and Z be variables. We will denote the symbol ~ as a dependency relation. That is, X ~ Y
means that X is dependent of Y. Consider the case where X ~ Y. Then there exists an edge between
X and Y. If Z is independent of both, then Y is disconnected from X-Y. The converse holds, since if
there exists no edge between them they are independent. Now consider X ~ Y and Y ~ Z. Since X-Y,
Y-Z and therefore the graph is connected. The last case is when everyone is independent of
everyone, in which case there are no edges and all variables are disconnected. We can assume X, Y
and Z as sets of variables for the general case.

To construct the graph, we can check for dependencies on each distinct pair of variables (u,v) of
set X. If there exists a dependency, add an edge u-v. Else, skip. It is clear that the complexity
for constructing such graph is O(n^2), since we must check each possible pairwise combination.

Once we have a constructed independency graph we must now discriminate each complete subgraph in
the independency graph. We can do this by Union-Find.

	Initially each vertex has its own set.
	For each vertex v:
		For each edge v-u:
			If u is not in the same set of v then
				Union(u, v)
			EndIf
		EndFor
	EndFor

After passing through every vertex, we have k connected subgraphs. These k subgraphs are indepedent
of each other. Return these k-sets.
*/
type IndepGraph struct {
	// Adjacency list containing each vertex and to which other vertices it is connected to and from.
	adjlist map[int][]int
	// This k-set contains the connected subgraphs that are completely separated from each other.
	Kset [][]int
}

// Constructs a new IndepGraph given a DataGroup.
func NewIndepGraph(data []*VarData) *IndepGraph {
	igraph := IndepGraph{make(map[int][]int), nil}
	n := len(data)

	// IDs and Reverse IDs.
	ids := make([]int, n)
	rids := make(map[int]int)

	for i := 0; i < n; i++ {
		ids[i] = data[i].Varid
		rids[ids[i]] = i
		igraph.adjlist[ids[i]] = []int{}
	}

	//fmt.Println("Constructing independency graph...")
	// Construct the indepedency graph by adding an edge if there exists a dependency relation.
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			v1, v2 := ids[i], ids[j]

			// Initialize the count matrix mdata.
			//fmt.Println("Initializing count matrix...")
			p, q := data[i].Categories, data[j].Categories
			mdata := make([][]int, p+1)
			for k := 0; k < p+1; k++ {
				mdata[k] = make([]int, q+1)
			}

			// len(data[i].Data) == len(data[j].Data) by definition.
			m := len(data[i].Data)
			for k := 0; k < m; k++ {
				mdata[data[i].Data[k]][data[j].Data[k]]++
			}

			//fmt.Println("Counting totals and assigning to edges...")
			// Total on the x axis, y axis and x+y respectively.
			tx, ty, tt := make([]int, q), 0, 0
			//fmt.Println("Y-axis...")
			for x := 0; x < p; x++ {
				ty = 0
				for y := 0; y < q; y++ {
					ty += mdata[x][y]
					tx[y] += mdata[x][y]
				}
				mdata[x][q] = ty
			}
			// Compute total on the x axis.
			//fmt.Println("X-axis...")
			for y := 0; y < q; y++ {
				mdata[p][y] = tx[y]
				tt += tx[y]
			}
			// Total total.
			mdata[p][q] = tt

			// Checks if variables i, j are independent.
			//fmt.Println("Checking for pairwise independency...")
			indep := ChiSquareTest(p, q, mdata, n-1)

			//fmt.Printf("%t\n", indep)
			// If not independent, then add an undirected edge i-j.
			if !indep {
				//fmt.Println("Not independent. Creating edge...")
				igraph.adjlist[v1] = append(igraph.adjlist[v1], v2)
				igraph.adjlist[v2] = append(igraph.adjlist[v2], v1)
			} //else {
			//fmt.Println("Independent. No edges.")
			//}
		}
	}

	// Union-Find to discriminate each set of connected variables that are fully disconnected of
	// another set of connected set of variables
	//fmt.Println("Finding disconnected subgraphs...")

	// Set of Union-Find trees.
	sets := make([]*UFNode, n)

	// At first every vertex has its own set.
	for i := 0; i < n; i++ {
		sets[i] = MakeSet(ids[i])
	}

	//fmt.Println("Preparing to test each vertex of the independency graph for disconnectivity...")
	// If a vertex u has an edge with another vertex v, then union sets that contain u and v.
	for i := 0; i < n; i++ {
		v1 := ids[i]
		m := len(igraph.adjlist[v1])
		for j := 0; j < m; j++ {
			v2 := igraph.adjlist[v1][j]
			rv2 := rids[v2]

			if Find(sets[i]) == Find(sets[rv2]) {
				continue
			}

			Union(sets[i], sets[rv2])
		}
	}

	igraph.Kset = nil
	for i := 0; i < n; i++ {
		if sets[i] == sets[i].Pa {
			igraph.Kset = append(igraph.Kset, UFVarids(sets[i]))
		}
	}

	return &igraph
}
