from graph_tool.all import *

g = Graph(directed=True)
vcolors = g.new_vertex_property("string")
vnames = g.new_vertex_property("string")
enames = g.new_edge_property("string")

def add_node(name, type):
	v=g.add_vertex()
	vnames[v]=name
	vcolors[v]=type
	return v

def add_edge(o, t, name):
	e=g.add_edge(o, t)
	enames[e]=name
	return e

def add_edge_nameless(o, t):
	e=g.add_edge(o, t)
	return e


S0 = add_node("+", "#ff3300")
P0 = add_node("*", "#669900")
add_edge(S0, P0, "0.200")
X0 = add_node("X_4", "#0066ff")
add_edge_nameless(P0, X0)
P1 = add_node("*", "#669900")
add_edge(S0, P1, "0.800")
X1 = add_node("X_5", "#0066ff")
add_edge_nameless(P1, X1)
S1 = add_node("+", "#ff3300")
add_edge_nameless(P0, S1)
X2 = add_node("X_0", "#0066ff")
add_edge(S1, X2, "0.300")
X3 = add_node("X_1", "#0066ff")
add_edge(S1, X3, "0.700")
S2 = add_node("+", "#ff3300")
add_edge_nameless(P1, S2)
X4 = add_node("X_2", "#0066ff")
add_edge(S2, X4, "0.400")
X5 = add_node("X_3", "#0066ff")
add_edge(S2, X5, "0.600")
g.vertex_properties["name"]=vnames
g.vertex_properties["color"]=vcolors

graph_draw(g, vertex_text=g.vertex_properties["name"], edge_text=enames, vertex_fill_color=g.vertex_properties["color"], output="/home/renatogeh/go/src/github.com/RenatoGeh/gospn/results/example/simplespn/spn.png")
