package graph

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
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
