package graph

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/graph"

	amjson "github.com/sevein/amflow/internal/graph/encoding"
)

// Vertex is the fundamental unit to describe operations in the workflow.
type Vertex interface {
	graph.Node

	AMID() string
	Source() amjson.Vertex
	WithHighlight(bool)
}

// highlight knows if it needs to be highlighted.
type highlight struct {
	highlight bool
}

func (h highlight) IsHighlighted() bool {
	return h.highlight
}

func (h *highlight) WithHighlight(v bool) {
	h.highlight = v
}

// VertexChainLink is a Vertex.
type VertexChainLink struct {
	v   graph.Node
	src *amjson.Chain
	highlight
}

func (v VertexChainLink) ID() int64 {
	return v.v.ID()
}

func (v VertexChainLink) AMID() string {
	return v.src.ID()
}

func (v VertexChainLink) String() string {
	return fmt.Sprintf("VertexChainLink[%s] - %s", v.AMID(), v.src.Description["en"])
}

func (v VertexChainLink) Source() amjson.Vertex {
	return v.src
}

// VertexLink is a Vertex.
type VertexLink struct {
	v   graph.Node
	src *amjson.Link
	highlight
}

func (v VertexLink) ID() int64 {
	return v.v.ID()
}

func (v VertexLink) AMID() string {
	return v.src.ID()
}

func (v VertexLink) String() string {
	return fmt.Sprintf("VertexLink[%s] - %s", v.AMID(), v.src.Description["en"])
}

func (v VertexLink) Source() amjson.Vertex {
	return v.src
}

// VertexWatcheDir is a Vertex.
type VertexWatcheDir struct {
	v   graph.Node
	src *amjson.WatchedDirectory
	highlight
}

func (v VertexWatcheDir) ID() int64 {
	return v.v.ID()
}

func (v VertexWatcheDir) AMID() string {
	return v.src.ID()
}

func (v VertexWatcheDir) String() string {
	return fmt.Sprintf("VertexWatchedDir[%s]", v.AMID())
}

func (v VertexWatcheDir) Source() amjson.Vertex {
	return v.src
}

func (v VertexWatcheDir) isInitiator() bool {
	// v.src.Path == "/system/createAIC/"
	// v.src.Path == "/system/reingestAIP"
	return strings.HasPrefix(v.src.Path, "/activeTransfers/")
}
