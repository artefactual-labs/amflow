package graph

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/simple"
)

type exitCodeEdge struct {
	simple.Edge
	code   int
	status string

	// Indicates if there should be two edges (as we don't have a multigraph).
	hasFallback bool
}

func newExitCodeEdge(f graph.Node, t graph.Node, code int, status string, hasFallback bool) *exitCodeEdge {
	return &exitCodeEdge{
		Edge:        simple.Edge{F: f, T: t},
		code:        code,
		status:      status,
		hasFallback: hasFallback,
	}
}

func (e exitCodeEdge) Attributes() []encoding.Attribute {
	var fallbackDetail string
	if e.hasFallback {
		fallbackDetail = " [!]"
	}
	attrs := []encoding.Attribute{
		{Key: "label", Value: esc(fmt.Sprintf("%s (code %d)%s", e.status, e.code, fallbackDetail))},
	}
	if e.status == "Failed" {
		attrs = append(attrs, encoding.Attribute{Key: "color", Value: "red"})
	}
	return attrs
}

var _ encoding.Attributer = exitCodeEdge{}

type defaultFallbackEdge struct {
	simple.Edge
	status string
}

func newDefaultFallbackEdge(f graph.Node, t graph.Node, status string) *defaultFallbackEdge {
	return &defaultFallbackEdge{
		Edge:   simple.Edge{F: f, T: t},
		status: status,
	}
}

func (e defaultFallbackEdge) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "color", Value: "gray"},
		{Key: "label", Value: esc(fmt.Sprintf("Fallback (%s)", e.status))},
	}
}

var _ encoding.Attributer = defaultFallbackEdge{}

type chainChoiceEdge struct {
	simple.Edge
}

func newChainChoiceEdge(f graph.Node, t graph.Node) *chainChoiceEdge {
	return &chainChoiceEdge{
		Edge: simple.Edge{F: f, T: t},
	}
}

func (e chainChoiceEdge) Attributes() []encoding.Attribute {
	chain := e.Edge.T.(*VertexChainLink)
	return []encoding.Attribute{
		{Key: "label", Value: esc(chain.src.Description["en"])},
	}
}

type linkChoiceEdge struct {
	simple.Edge
}

func newLinkChoiceEdge(f graph.Node, t graph.Node) *linkChoiceEdge {
	return &linkChoiceEdge{
		Edge: simple.Edge{F: f, T: t},
	}
}

func (e linkChoiceEdge) Attributes() []encoding.Attribute {
	chain := e.Edge.T.(*VertexLink)
	return []encoding.Attribute{
		{Key: "label", Value: esc(chain.src.Description["en"])},
	}
}

var _ encoding.Attributer = chainChoiceEdge{}

type virtualMovingDirEdge struct {
	simple.Edge
}

func newVirtualMovingDirBridge(f graph.Node, t graph.Node) *virtualMovingDirEdge {
	return &virtualMovingDirEdge{
		Edge: simple.Edge{F: f, T: t},
	}
}

func (e virtualMovingDirEdge) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "color", Value: "gray"},
		{Key: "style", Value: "dashed"},
	}
}

var _ encoding.Attributer = virtualMovingDirEdge{}
