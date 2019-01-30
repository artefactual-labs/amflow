package graph

import "gonum.org/v1/gonum/graph"

// gcopy mimics graph.Copy but it uses the same edges instead of creating new
// objects so we don't lose their extra properties.
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
			// The original implementation does the following:
			//   v := to.Node()
			//   dst.SetEdge(dst.NewEdge(u, v))
			// This is what we're doing instead to preserve the origina value:
			dst.SetEdge(src.Edge(u.ID(), to.Node().ID()))
		}
	}
}
