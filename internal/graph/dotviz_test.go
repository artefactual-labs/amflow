package graph

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"

	amjson "github.com/artefactual-labs/amflow/internal/graph/encoding"
)

func TestWorkflow_SVG(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	// Create a new graph. You can also use populateGraph(t).
	w := New(nil)
	n0 := w.addVertex(&amjson.Link{})
	n1 := w.addVertex(&amjson.Link{})
	n2 := w.addVertex(&amjson.Link{})
	n3 := w.addVertex(&amjson.Link{})
	w.graph.SetEdge(simple.Edge{F: n0, T: n1})
	w.graph.SetEdge(simple.Edge{F: n2, T: n3})

	blob, err := w.SVG()
	require.NoError(t, err)

	var v struct {
		XMLName xml.Name `xml:"svg"`
	}
	err = xml.Unmarshal(blob, &v)
	require.NoError(t, err)
	require.Equal(t, "svg", v.XMLName.Local)
}

// Test that graph.Copy copies the topology but reuses our nodes/edges.
// See https://github.com/gonum/gonum/pull/897 for more details.
func Test_copy(t *testing.T) {
	// Original graph.
	src := New(nil)
	n0 := src.addVertex(&amjson.Link{})
	n1 := src.addVertex(&amjson.Link{})
	e0 := src.graph.NewEdge(n0, n1)
	src.graph.SetEdge(e0)

	// New graph.
	dst := simple.NewDirectedGraph()

	graph.Copy(dst, src)

	// Confirm that we have access to the same edge in the new graph.
	// Interface values are comparable. Two interface values are equal if they
	// have identical dynamic types and equal dynamic values or if both have
	// value nil.
	e := dst.Edge(0, 1)
	if e0 != e {
		t.Fail()
	}
}
