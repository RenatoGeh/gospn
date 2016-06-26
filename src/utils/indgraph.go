package utils

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
}
