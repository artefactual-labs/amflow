package graph

import (
	"fmt"

	"gonum.org/v1/gonum/graph"

	amjson "github.com/sevein/amflow/internal/graph/encoding"
)

// Vertex is the fundamental unit to describe operations in the workflow.
type Vertex interface {
	graph.Node

	AMID() string
}

// VertexChainLink is a Vertex.
type VertexChainLink struct {
	v   graph.Node
	src *amjson.Chain
}

func (v *VertexChainLink) ID() int64 {
	return v.v.ID()
}

func (v *VertexChainLink) AMID() string {
	return v.src.ID()
}

func (v *VertexChainLink) String() string {
	return fmt.Sprintf("VertexChainLink[%s] - %s", v.AMID(), v.src.Description["en"])
}

// VertexLink is a Vertex.
type VertexLink struct {
	v   graph.Node
	src *amjson.Link
}

func (v *VertexLink) ID() int64 {
	return v.v.ID()
}

func (v *VertexLink) AMID() string {
	return v.src.ID()
}

func (v *VertexLink) String() string {
	return fmt.Sprintf("VertexLink[%s] - %s", v.AMID(), v.src.Description["en"])
}

// VertexWatcheDir is a Vertex.
type VertexWatcheDir struct {
	v   graph.Node
	src *amjson.WatchedDirectory
}

func (v *VertexWatcheDir) ID() int64 {
	return v.v.ID()
}

func (v *VertexWatcheDir) AMID() string {
	return v.src.ID()
}

func (v *VertexWatcheDir) String() string {
	return fmt.Sprintf("WatchedDir[%s]", v.AMID())
}
