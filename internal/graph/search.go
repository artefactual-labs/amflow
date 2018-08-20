package graph

import "gonum.org/v1/gonum/graph"

// List ancestors of a given vertice.
// It does depth-first search supported by a stack.
func ListAncestors(w *Workflow, amid string) int64s {
	var stack stack
	visited := make(int64s)
	to := w.VertexByAMID(amid)
	if to == nil {
		return nil
	}
	stack.Push(to)
	visited.Add(to.ID())
	for stack.Len() > 0 {
		cur := stack.Pop()
		curID := cur.ID()
		ancs := w.To(curID)
		for ancs.Next() {
			anc := ancs.Node()
			ancID := anc.ID()
			if visited.Has(ancID) {
				continue
			}
			visited.Add(ancID)
			stack.Push(anc)
		}
	}
	return visited
}

// stack is a LIFO we need for depth-first searches.
type stack []graph.Node

func (s *stack) Len() int { return len(*s) }

func (s *stack) Pop() graph.Node {
	v := *s
	v, n := v[:len(v)-1], v[len(v)-1]
	*s = v
	return n
}

func (s *stack) Push(n graph.Node) { *s = append(*s, n) }

// int64s stores int64 values without repetitions.
type int64s map[int64]struct{}

func (s int64s) Add(e int64) {
	s[e] = struct{}{}
}

func (s int64s) Has(e int64) bool {
	_, ok := s[e]
	return ok
}

func (s int64s) Remove(e int64) {
	delete(s, e)
}

func (s int64s) Count() int {
	return len(s)
}
