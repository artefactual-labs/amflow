package graph

import (
	"fmt"

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

var _ graph.Node = initiatorVertex{}

func (v initiatorVertex) DOTID() string {
	return fmt.Sprintf(`START`)
}

var _ dot.Node = initiatorVertex{}

func (v initiatorVertex) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		encoding.Attribute{"shape", "diamond"},
		encoding.Attribute{"style", "filled"},
		encoding.Attribute{"color", "blue"},
		encoding.Attribute{"fontsize", "60"},
		encoding.Attribute{"fontcolor", "white"},
		encoding.Attribute{"margin", "0.75"},
	}
}

var _ encoding.Attributer = initiatorVertex{}
