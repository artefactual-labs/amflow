package graph

import "gonum.org/v1/gonum/graph"

// gcopy mimics graph.Copy but it uses the same edges instead of creating new
// objects, so the properties are carried over.
func gcopy(dst graph.Builder, src graph.Graph) {
	nodes := src.Nodes()
	for nodes.Next() {
		dst.AddNode(nodes.Node())
	}
	nodes.Reset()
	for nodes.Next() {
		u := nodes.Node()
		to := src.From(u.ID())
		for to.Next() {
			dst.SetEdge(src.Edge(u.ID(), to.Node().ID()))
		}
	}
}
