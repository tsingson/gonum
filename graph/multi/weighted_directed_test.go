// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multi_test

import (
	"testing"

	"github.com/tsingson/gonum/graph"
	"github.com/tsingson/gonum/graph/internal/set"
	"github.com/tsingson/gonum/graph/multi"
	"github.com/tsingson/gonum/graph/testgraph"
)

func weightedDirectedBuilder(nodes []graph.Node, edges []graph.WeightedLine, self, absent float64) (g graph.Graph, n []graph.Node, e []graph.Edge, s, a float64, ok bool) {
	seen := make(set.Nodes)
	dg := multi.NewWeightedDirectedGraph()
	dg.EdgeWeightFunc = func(l graph.WeightedLines) float64 {
		// TODO(kortschak): Remove nil guard if nil iterators
		// are forbidden, https://github.com/gonum/gonum/issues/614.
		if l == nil || l.Len() == 0 {
			return absent
		}
		var w float64
		for l.Next() {
			w += l.WeightedLine().Weight()
		}
		l.Reset()
		return w
	}
	for _, n := range nodes {
		seen.Add(n)
		dg.AddNode(n)
	}
	for _, edge := range edges {
		f := dg.Node(edge.From().ID())
		if f == nil {
			f = edge.From()
		}
		t := dg.Node(edge.To().ID())
		if t == nil {
			t = edge.To()
		}
		cl := multi.WeightedLine{F: f, T: t, UID: edge.ID(), W: edge.Weight()}
		seen.Add(cl.F)
		seen.Add(cl.T)
		e = append(e, cl)
		dg.SetWeightedLine(cl)
	}
	if len(seen) != 0 {
		n = make([]graph.Node, 0, len(seen))
	}
	for _, sn := range seen {
		n = append(n, sn)
	}
	return dg, n, e, self, absent, true
}

func TestWeightedDirected(t *testing.T) {
	t.Run("EdgeExistence", func(t *testing.T) {
		testgraph.EdgeExistence(t, weightedDirectedBuilder)
	})
	t.Run("NodeExistence", func(t *testing.T) {
		testgraph.NodeExistence(t, weightedDirectedBuilder)
	})
	t.Run("ReturnAdjacentNodes", func(t *testing.T) {
		testgraph.ReturnAdjacentNodes(t, weightedDirectedBuilder, true)
	})
	t.Run("ReturnAllLines", func(t *testing.T) {
		testgraph.ReturnAllLines(t, weightedDirectedBuilder, true)
	})
	t.Run("ReturnAllNodes", func(t *testing.T) {
		testgraph.ReturnAllNodes(t, weightedDirectedBuilder, true)
	})
	t.Run("ReturnAllWeightedLines", func(t *testing.T) {
		testgraph.ReturnAllWeightedLines(t, weightedDirectedBuilder, true)
	})
	t.Run("ReturnNodeSlice", func(t *testing.T) {
		testgraph.ReturnNodeSlice(t, weightedDirectedBuilder, true)
	})
	t.Run("Weight", func(t *testing.T) {
		testgraph.Weight(t, weightedDirectedBuilder)
	})
}

// Tests Issue #27
func TestWeightedEdgeOvercounting(t *testing.T) {
	g := generateDummyWeightedGraph()

	if neigh := graph.NodesOf(g.From(int64(2))); len(neigh) != 2 {
		t.Errorf("Node 2 has incorrect number of neighbors got neighbors %v (count %d), expected 2 neighbors {0,1}", neigh, len(neigh))
	}
}

func generateDummyWeightedGraph() *multi.WeightedDirectedGraph {
	nodes := [4]struct{ srcID, targetID int }{
		{2, 1},
		{1, 0},
		{2, 0},
		{0, 2},
	}

	g := multi.NewWeightedDirectedGraph()

	for i, n := range nodes {
		g.SetWeightedLine(multi.WeightedLine{F: multi.Node(n.srcID), T: multi.Node(n.targetID), W: 1, UID: int64(i)})
	}

	return g
}

// Test for issue #123 https://github.com/gonum/graph/issues/123
func TestIssue123WeightedDirectedGraph(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	g := multi.NewWeightedDirectedGraph()

	n0 := g.NewNode()
	g.AddNode(n0)

	n1 := g.NewNode()
	g.AddNode(n1)

	g.RemoveNode(n0.ID())

	n2 := g.NewNode()
	g.AddNode(n2)
}
