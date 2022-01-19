package graph

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"go.uber.org/multierr"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	amjson "github.com/artefactual-labs/amflow/internal/graph/encoding"
)

var WorkflowSchemaBox = packr.New("workflow", "./schema")

// Workflow is a sequence of operations in Archivematica.
//
// It is modeled as a simple directed graph.
type Workflow struct {
	// Underlying directed graph.
	graph *simple.DirectedGraph

	// Internal mappings for convenience.
	vxByID   map[int64]Vertex
	vxByAMID map[string]Vertex
}

// New returns a Workflow.
func New(data *amjson.WorkflowData) *Workflow {
	w := &Workflow{
		graph:    simple.NewDirectedGraph(),
		vxByID:   map[int64]Vertex{},
		vxByAMID: map[string]Vertex{},
	}
	if data != nil {
		w.load(data)
	}
	return w
}

// AddVertex adds a new vertex to the workflow.
func (w *Workflow) addVertex(v amjson.Vertex) Vertex {
	var vertex Vertex
	switch v := v.(type) {
	case *amjson.Chain:
		vertex = &VertexChainLink{
			v:   w.graph.NewNode(),
			src: v,
		}
	case *amjson.Link:
		vertex = &VertexLink{
			v:   w.graph.NewNode(),
			src: v,
		}
	case *amjson.WatchedDirectory:
		vertex = &VertexWatcheDir{
			v:   w.graph.NewNode(),
			src: v,
		}
	}
	w.graph.AddNode(vertex)
	w.vxByID[vertex.ID()] = vertex
	w.vxByAMID[vertex.AMID()] = vertex
	return vertex
}

// Vertex returns a workflow vertex given its AMD.
func (w Workflow) VertexByAMID(amid string) Vertex {
	v, ok := w.vxByAMID[amid]
	if !ok {
		return nil
	}
	return v
}

// Vertex returns a workflow vertex given its ID.
func (w Workflow) VertexByID(id int64) Vertex {
	v, ok := w.vxByID[id]
	if !ok {
		return nil
	}
	return v
}

// hasMultipleComponents determines if every vertex is reachable from every
// other vertex. Currently, Archivematica workflows are not expected to have
// more than one component (subgraph). This is a property observed in the
// existing workflow dataset but it may stop being that way in the future.
func (w Workflow) hasMultipleComponents() bool {
	cc := topo.ConnectedComponents(graph.Undirect{G: w.graph})
	return len(cc) > 1
}

// Check looks for inconsistencies in the graph. It returns a multiError
// created with the go.uber.org/multierr module.
func (w Workflow) Check() error {
	var err error
	for _, item := range w.watchedDirs() {
		id := item.ID()
		f, t := w.From(id), w.To(id)
		if f.Len() == 0 {
			err = multierr.Append(err, fmt.Errorf("[%s] watched directory does not make references", item.AMID()))
		}
		if t.Len() == 0 && !item.isInitiator() {
			err = multierr.Append(err, fmt.Errorf("[%s] watched directory is not referenced", item.AMID()))
		}
	}
	for _, item := range w.chains() {
		f, t := w.From(item.ID()), w.To(item.ID())
		if f.Len() == 0 {
			err = multierr.Append(err, fmt.Errorf("[%s] chain does not make references", item.AMID()))
		}
		if t.Len() == 0 {
			err = multierr.Append(err, fmt.Errorf("[%s] chain is not referenced", item.AMID()))
		}
	}
	for _, item := range w.links() {
		// Number of ascendants and descendants.
		fl := w.To(item.ID()).Len()
		tl := w.From(item.ID()).Len()

		// Orphan link!
		if fl == 0 && !item.src.Start {
			err = multierr.Append(err, fmt.Errorf("[%s] link is not referenced", item.AMID()))
		}

		// Do we have a transition to a watchedDir?
		var refsWatchedDir bool
		iter := w.From(item.ID())
		for iter.Next() {
			_, ok := iter.Node().(*VertexWatcheDir)
			if ok {
				refsWatchedDir = true
				break
			}
		}

		// Look for terminal links that look like false positives.
		if item.src.End {
			switch {
			case tl == 0:
				// Terminal link that actually terminates, nothing wrong here.
			case tl >= 1:

				// Okay to terminal links that only points to a WD, assuming
				// it's end of package but hard to tell from here.
				if tl == 1 && refsWatchedDir {
					continue
				}

				// We have a few terminal links that either:
				// - Refer to a watched directory *and* to another link.
				// - Does not refer to a watched directory *but* refers to a link.
				//
				// In both cases, they're really only terminal if the don't
				// continue after execution. E.g. "Move transfer to backlog" is
				// ending when the execution succeeds but it continues otherwise.
				//
				// MCPServer will not mark the package as completed in case of
				// false positives, but it will low a warning.
				err = multierr.Append(err, fmt.Errorf("[%s] link is terminal but has alternative paths [children=%d] [refsWD=%t]", item.AMID(), tl, refsWatchedDir))
			}
		}

		// Look for unidentified terminal links (false negatives).
		if !item.src.End && tl == 0 {
			err = multierr.Append(err, fmt.Errorf("[%s] link should be terminal", item.AMID()))
		}

		// TODO: This could guess a false negative but it's most likely the end
		// of a chain transitioning to a new chain within the same package. Not
		// something we necessarily need to convert into a terminal link.
		// if !item.src.End && tl == 1 && refsWatchedDir { err = multierr.Append(err, fmt.Errorf("[%s] link should be terminal", item.AMID())) }
	}
	return err
}

func (w Workflow) watchedDirs() []*VertexWatcheDir {
	ret := []*VertexWatcheDir{}
	for _, v := range w.vxByID {
		vwd, ok := v.(*VertexWatcheDir)
		if ok {
			ret = append(ret, vwd)
		}
	}
	return ret
}

func (w Workflow) chains() []*VertexChainLink {
	ret := []*VertexChainLink{}
	for _, v := range w.vxByID {
		vwd, ok := v.(*VertexChainLink)
		if ok {
			ret = append(ret, vwd)
		}
	}
	return ret
}

func (w Workflow) links() []*VertexLink {
	ret := []*VertexLink{}
	for _, v := range w.vxByID {
		vwd, ok := v.(*VertexLink)
		if ok {
			ret = append(ret, vwd)
		}
	}
	return ret
}

// load workflow data. It includes vertices and edges. The latter are mostly
// explicit in the workflow data, excepting move filesystem operations.
func (w *Workflow) load(data *amjson.WorkflowData) {
	// Links.
	_lns := make(map[string]*VertexLink)
	for id, item := range data.Links {
		_lns[id] = w.addVertex(item).(*VertexLink)
	}

	// Chain links.
	_chs := make(map[string]*VertexChainLink)
	for id, item := range data.Chains {
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
		case vertexSrc.src.Config.Model == "MicroServiceChainChoice" && len(vertexSrc.src.Config.ChainChoices) > 0:
			{
				for _, id := range vertexSrc.src.Config.ChainChoices {
					if vertexDst, ok := _chs[id]; ok {
						w.graph.SetEdge(newChainChoiceEdge(vertexSrc, vertexDst))
					}
				}
			}
		case vertexSrc.src.Config.Manager == "linkTaskManagerChoice":
			if len(vertexSrc.src.Config.LinkChoices) > 0 {
				for _, linkChoice := range vertexSrc.src.Config.LinkChoices {
					if vertexDst, ok := _lns[linkChoice.LinkID]; ok {
						w.graph.SetEdge(newLinkChoiceEdge(vertexSrc, vertexDst))
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
		// This section below declares edges for associations that are a result
		// of filesystem moving operations that MCPServer identifies by watching
		// directories. We've found this mechanism to be undesirable and it will
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
					if transitionDir := w.VertexByAMID("/system/autoProcessSIP"); transitionDir != nil {
						w.graph.SetEdge(newVirtualMovingDirBridge(vertexSrc, w.VertexByAMID("/system/autoProcessSIP")))
					}
				}
			}
		}
	}
}

// Implement graph.Graph.
func (w Workflow) Node(id int64) graph.Node           { return w.graph.Node(id) }
func (w Workflow) Nodes() graph.Nodes                 { return w.graph.Nodes() }
func (w Workflow) From(id int64) graph.Nodes          { return w.graph.From(id) }
func (w Workflow) HasEdgeBetween(xid, yid int64) bool { return w.graph.HasEdgeBetween(xid, yid) }
func (w Workflow) Edge(uid, vid int64) graph.Edge     { return w.graph.Edge(uid, vid) }

var _ graph.Graph = Workflow{}

// Implement graph.Directed.
func (w Workflow) HasEdgeFromTo(uid, vid int64) bool { return w.graph.HasEdgeFromTo(uid, vid) }
func (w Workflow) To(id int64) graph.Nodes           { return w.graph.To(id) }

var _ graph.Directed = Workflow{}
