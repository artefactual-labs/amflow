package graph

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"

	amjson "github.com/sevein/amflow/internal/graph/encoding"
)

func populate(t *testing.T) *Workflow {
	bytes, _ := ioutil.ReadFile("./schema/example.json")
	data, err := amjson.LoadWorkflowData(bytes)
	if err != nil && t != nil {
		t.Fatal(err)
	}
	return New(data)
}

// Making sure that this is fast enough. It seems that the example JSON doc
// can be decoded/loaded into the graph in less than 10ms, that's pretty good.
func BenchmarkPopulate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		populate(nil)
	}
}

func TestComponents(t *testing.T) {
	w := New(nil)
	n0 := w.addVertex(&amjson.Link{})
	n1 := w.addVertex(&amjson.Link{})
	n2 := w.addVertex(&amjson.Link{})
	n3 := w.addVertex(&amjson.Link{})
	w.graph.SetEdge(simple.Edge{F: n0, T: n1})
	w.graph.SetEdge(simple.Edge{F: n2, T: n3})

	assert.True(t, w.hasMultipleComponents())
}

func TestEncode(t *testing.T) {
	// TODO: workflows need to be encoded back to JSON when needed.
}

func TestWalk(t *testing.T) {
	w := populate(nil)
	init := w.vertex("/activeTransfers/standardTransfer")

	t.Log(w.graph.Nodes().Len())

	visitfn := func(u, v graph.Node) {
		t.Log(u, v)
	}

	t.Log(" ~~~ Depth first ~~~")
	dfs := &traverse.DepthFirst{Visit: visitfn}
	dfsUntil := func(v graph.Node) bool {
		return false
	}
	dfs.Walk(w, init, dfsUntil)

	t.Log(" ~~~ Breadth first ~~~")
	bfs := &traverse.BreadthFirst{Visit: visitfn}
	bfsUntil := func(v graph.Node, depth int) bool {
		return false
	}
	bfs.Walk(w, init, bfsUntil)
}
