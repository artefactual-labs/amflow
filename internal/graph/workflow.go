package graph

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	amjson "github.com/sevein/amflow/internal/graph/encoding"
)

var WorkflowSchemaBox = packr.New("workflow", "./schema")

// Workflow is a sequence of operations in Archivematica.
//
// It is modeled as a simple directed graph.
type Workflow struct {
	// Collection of vertices grouped by their types, e.g. chains, links or
	// watched directories. Vertex is an interface type.
	vertices map[nodeType]idNodeMapping

	// Underlying directed graph.
	graph *simple.DirectedGraph
}

// New returns a Workflow.
func New(data *amjson.WorkflowData) *Workflow {
	w := &Workflow{
		vertices: map[nodeType]idNodeMapping{},
		graph:    simple.NewDirectedGraph(),
	}
	for nt := range nodeTypes {
		w.vertices[nt] = idNodeMapping{}
	}
	if data != nil {
		w.load(data)
	}
	return w
}

type idNodeMapping map[string]Vertex

type nodeType byte

const (
	chainNodeType nodeType = iota
	linkNodeType
	watchedDirNodeType
)

var nodeTypes = map[nodeType]string{
	chainNodeType:      "chain",
	linkNodeType:       "link",
	watchedDirNodeType: "watchedDirectory",
}

// AddVertex adds a new vertex to the workflow.
func (w *Workflow) addVertex(v amjson.Vertex) Vertex {
	var (
		t      nodeType
		vertex Vertex
	)
	switch v := v.(type) {
	case *amjson.Chain:
		t = chainNodeType
		vertex = &VertexChainLink{
			v:   w.graph.NewNode(),
			src: v,
		}
	case *amjson.Link:
		t = linkNodeType
		vertex = &VertexLink{
			v:   w.graph.NewNode(),
			src: v,
		}
	case *amjson.WatchedDirectory:
		t = watchedDirNodeType
		vertex = &VertexWatcheDir{
			v:   w.graph.NewNode(),
			src: v,
		}
	}
	w.graph.AddNode(vertex)
	w.vertices[t][v.ID()] = vertex
	return vertex
}

// vertex returns a workflow vertex given its ID.
func (w *Workflow) vertex(id string) Vertex {
	for nt := range nodeTypes {
		if v, ok := w.vertices[nt][id]; ok {
			return v
		}
	}
	return nil
}

// hasMultipleComponents determines if every vertex is reachable from every
// other vertex. Currently, Archivematica workflows are not expected to have
// more than one component (subgraph). This is a property observed in the
// existing workflow dataset but it may stop being that way in the future.
func (w Workflow) hasMultipleComponents() bool {
	cc := topo.ConnectedComponents(graph.Undirect{G: w.graph})
	return len(cc) > 1
}

var ignoredLinks = map[string]string{
	"61c316a6-0a50-4f65-8767-1f44b1eeb6dd": "Link - Email fail report",
	"7d728c39-395f-4892-8193-92f086c0546f": "Link - Email fail report",
	"333532b9-b7c2-4478-9415-28a3056d58df": "Link - Move to the rejected directory",
	"19c94543-14cb-4158-986b-1d2b55723cd8": "Link - Cleanup rejected SIP",

	// The vertices that follow the ones above.
	"e780473a-0c10-431f-bab6-5d7238b2b70b": "...",
	"377f8ebb-7989-4a68-9361-658079ff8138": "...",
	"b2ef06b9-bca4-49da-bc5c-866d7b3c4bb1": "...",
	"828528c2-2eb9-4514-b5ca-dfd1f7cb5b8c": "...",
	"3467d003-1603-49e3-b085-e58aa693afed": "...",

	// The vertex that follows the ignored chain "Reject transfer".
	"ae5cdd0d-2f81-4935-a380-d5c6f1337d93": "Link - Cleanup rejected transfer",
}

var ignoredChains = map[string]string{
	"1b04ec43-055c-43b7-9543-bd03c6a778ba": "Chain - Reject transfer",
}

// load workflow data.
//
//
func (w *Workflow) load(data *amjson.WorkflowData) {

	// Links.
	_lns := make(map[string]*VertexLink)
	for id, item := range data.Links {
		if _, ignored := ignoredLinks[id]; ignored {
			continue
		}
		_lns[id] = w.addVertex(item).(*VertexLink)
	}

	// Chain links.
	_chs := make(map[string]*VertexChainLink)
	for id, item := range data.Chains {
		if _, ignored := ignoredChains[id]; ignored {
			continue
		}
		vertexSrc := w.addVertex(item).(*VertexChainLink)
		_chs[id] = vertexSrc
		if vertexDst, ok := _lns[item.LinkID]; ok {
			w.graph.SetEdge(w.graph.NewEdge(vertexSrc, vertexDst))
		}
	}

	// Watched directories.
	_wds := make(map[string]*VertexWatcheDir)
	for _, item := range data.WatchedDirectories {
		vertexSrc := w.addVertex(item).(*VertexWatcheDir)
		_wds[item.Path] = vertexSrc
		if vertexDst, ok := _chs[item.ChainID]; ok {
			w.graph.SetEdge(w.graph.NewEdge(vertexSrc, vertexDst))
		}
	}

	// Build a map of variables defined in TaskConfigSetUnitVariable links
	// and their respective links. This is going to be useful later to connect
	// pull links.
	_vars := map[string][]*VertexLink{}
	for _, node := range _lns {
		if node.src.Config.Model == "TaskConfigSetUnitVariable" {
			if match, ok := _lns[node.src.Config.ChainID]; ok {
				_vars[node.src.Config.Variable] = append(_vars[node.src.Config.Variable], match)
			}
		}
	}

	// Another pass to connect links.
	for _, vertexSrc := range _lns {

		// Connect to other links based on the fallback defined.
		if vertexSrc.src.FallbackLinkID != "" {
			if vertexDst, ok := _lns[vertexSrc.src.FallbackLinkID]; ok {
				w.graph.SetEdge(newDefaultFallbackEdge(vertexSrc, vertexDst, vertexSrc.src.FallbackJobStatus))
			}
		}

		// Connect to other links based on the exit codes.
		for code, ec := range vertexSrc.src.ExitCodes {
			if ec.LinkID == "" {
				continue
			}
			if vertexDst, ok := _lns[ec.LinkID]; ok {
				var hasFallback bool
				if e := w.graph.Edge(vertexSrc.ID(), vertexDst.ID()); e != nil {
					hasFallback = true
				}
				w.graph.SetEdge(newExitCodeEdge(vertexSrc, vertexDst, code, ec.JobStatus, hasFallback))
			}
		}

		switch {
		case vertexSrc.src.Config.Model == "MicroServiceChainChoice" && len(vertexSrc.src.Config.Choices) > 0:
			{
				for _, id := range vertexSrc.src.Config.Choices {
					if vertexDst, ok := _chs[id]; ok {
						w.graph.SetEdge(newChainChoiceEdge(vertexSrc, vertexDst))
					}
				}
			}
		case vertexSrc.src.Config.Manager == "linkTaskManagerUnitVariableLinkPull":
			{
				if values, ok := _vars[vertexSrc.src.Config.Variable]; ok {
					for _, vertexDst := range values {
						w.graph.SetEdge(w.graph.NewEdge(vertexSrc, vertexDst))
					}
				}
				if vertexSrc.src.Config.ChainID != "" {
					if vertexDst, ok := _lns[vertexSrc.src.Config.ChainID]; ok {
						w.graph.SetEdge(w.graph.NewEdge(vertexSrc, vertexDst))
					}
				}
			}
		// This section below declares edges for associations that are a result of filesystem moving operations that
		// MCPServer identifies by watching directories. We've found this mechanism to be undesirable and it will
		// probably change soon.
		case vertexSrc.src.Config.Manager == "linkTaskManagerDirectories":
			{
				if strings.HasPrefix(vertexSrc.src.Config.Execute, "move") {
					args := vertexSrc.src.Config.Arguments
					for path, vertexDst := range _wds {
						substr1 := fmt.Sprintf("%%watchedDirectories%s", path)
						substr2 := fmt.Sprintf("%%watchDirectoryPath%%%s", path[1:])
						if strings.Contains(args, substr1) || strings.Contains(args, substr2) {
							w.graph.SetEdge(newVirtualMovingDirBridge(vertexSrc, vertexDst))
						}
					}
				} else if vertexSrc.src.Description["en"] == "Create SIP from transfer objects" || vertexSrc.src.Description["en"] == "Create SIPs from TRIM transfer containers" {
					w.graph.SetEdge(newVirtualMovingDirBridge(vertexSrc, w.vertex("/system/autoProcessSIP")))
				}
			}
		}
	}
}

// Implement graph.Graph.
func (w *Workflow) Node(id int64) graph.Node           { return w.graph.Node(id) }
func (w *Workflow) Nodes() graph.Nodes                 { return w.graph.Nodes() }
func (w *Workflow) From(id int64) graph.Nodes          { return w.graph.From(id) }
func (w *Workflow) HasEdgeBetween(xid, yid int64) bool { return w.graph.HasEdgeBetween(xid, yid) }
func (w *Workflow) Edge(uid, vid int64) graph.Edge     { return w.graph.Edge(uid, vid) }

// Implement graph.Directed.
func (w *Workflow) HasEdgeFromTo(uid, vid int64) bool { return w.graph.HasEdgeFromTo(uid, vid) }
func (w *Workflow) To(id int64) graph.Nodes           { return w.graph.To(id) }
