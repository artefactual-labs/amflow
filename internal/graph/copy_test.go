package graph

import (
	"testing"

	amjson "github.com/artefactual-labs/amflow/internal/graph/encoding"
	"gonum.org/v1/gonum/graph/simple"
)

func Test_gcopy(t *testing.T) {
	// Original graph.
	src := New(nil)
	n0 := src.addVertex(&amjson.Link{})
	n1 := src.addVertex(&amjson.Link{})
	e0 := src.graph.NewEdge(n0, n1)
	src.graph.SetEdge(e0)

	// New graph.
	dst := simple.NewDirectedGraph()

	gcopy(dst, src)

	// Confirm that we have access to the same edge in the new graph.
	// Interface values are comparable. Two interface values are equal if they
	// have identical dynamic types and equal dynamic values or if both have
	// value nil.
	e := dst.Edge(0, 1)
	if e0 != e {
		t.Fail()
	}
}
