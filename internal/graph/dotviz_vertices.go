package graph

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
)

type initiatorVertex struct {
	node graph.Node
}

func (v initiatorVertex) ID() int64 {
	return v.node.ID()
}

func (v initiatorVertex) DOTID() string {
	return "START"
}

func (v initiatorVertex) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "shape", Value: "diamond"},
		{Key: "style", Value: "filled"},
		{Key: "color", Value: "blue"},
		{Key: "fontsize", Value: "60"},
		{Key: "fontcolor", Value: "white"},
		{Key: "margin", Value: "0.75"},
	}
}

var _ graph.Node = initiatorVertex{}
var _ dot.Node = initiatorVertex{}
var _ encoding.Attributer = initiatorVertex{}
