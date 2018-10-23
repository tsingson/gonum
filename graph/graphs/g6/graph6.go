// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package g6 implements implicit graphs specified by graph6 format strings.
package g6 // import "gonum.org/v1/gonum/graph/graphs/g6"

import (
	"fmt"
	"math/big"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

// Graph is a graph6-represented graph.
//
// See https://users.cecs.anu.edu.au/~bdm/data/formats.txt for details
// and https://hog.grinvin.org/ for a source of interesting graphs in graph6
// format.
type Graph string

var _ graph.Graph = Graph("")

// Node returns the node with the given ID if it exists in the graph, and nil
// otherwise.
func (g Graph) Node(id int64) graph.Node {
	if id < 0 || numberOf(g) <= id {
		return nil
	}
	return simple.Node(id)
}

// Nodes returns all the nodes in the graph.
func (g Graph) Nodes() graph.Nodes {
	return iterator.NewImplicitNodes(0, int(numberOf(g)), func(id int) graph.Node { return simple.Node(id) })
}

// From returns all nodes that can be reached directly from the node with the
// given ID.
func (g Graph) From(id int64) graph.Nodes {
	if g.Node(id) == nil {
		return nil
	}
	return &g6Iterator{g: g, from: id, to: -1}
}

// HasEdgeBetween returns whether an edge exists between nodes with IDs xid
// and yid without considering direction.
func (g Graph) HasEdgeBetween(xid, yid int64) bool {
	if xid == yid {
		return false
	}
	if xid < 0 || numberOf(g) <= xid {
		return false
	}
	if yid < 0 || numberOf(g) <= yid {
		return false
	}
	return isSet(bitFor(xid, yid), g)
}

// Edge returns the edge from u to v, with IDs uid and vid, if such an edge
// exists and nil otherwise. The node v must be directly reachable from u as
// defined by the From method.
func (g Graph) Edge(uid, vid int64) graph.Edge {
	if !g.HasEdgeBetween(uid, vid) {
		return nil
	}
	return simple.Edge{simple.Node(uid), simple.Node(vid)}
}

// g6Iterator is a graph.Nodes for graph6 graph edges.
type g6Iterator struct {
	g    Graph
	from int64
	to   int64
}

var _ graph.Nodes = (*g6Iterator)(nil)

func (i *g6Iterator) Next() bool {
	n := numberOf(i.g)
	for i.to < n-1 {
		i.to++
		if i.to != i.from && isSet(bitFor(i.from, i.to), i.g) {
			return true
		}
	}
	return false
}

func (i *g6Iterator) Len() int {
	var cnt int
	n := numberOf(i.g)
	for to := i.to; to < n-1; {
		to++
		if to != i.from && isSet(bitFor(i.from, to), i.g) {
			cnt++
		}
	}
	return cnt
}

func (i *g6Iterator) Reset() { i.to = -1 }

func (i *g6Iterator) Node() graph.Node { return simple.Node(i.to) }

// numberOf returns the graph6-encoded number corresponding to g.
func numberOf(g Graph) int64 {
	if len(g) == 0 {
		return -1
	}
	if g[0] != 126 {
		return int64(g[0] - 63)
	}
	if g[1] != 126 {
		return int64(g[1]-63)<<12 | int64(g[2]-63)<<6 | int64(g[3]-63)
	}
	return int64(g[2]-63)<<30 | int64(g[3]-63)<<24 | int64(g[4]-63)<<18 | int64(g[5]-63)<<12 | int64(g[6]-63)<<6 | int64(g[7]-63)
}

// bitFor returns the index into the graph6 adjacency matrix for xid--yid.
func bitFor(xid, yid int64) int {
	if xid < yid {
		xid, yid = yid, xid
	}
	return int((xid*xid-xid)/2 + yid)
}

// isSet returns whether the given bit of the adjacency matrix is set.
func isSet(bit int, g Graph) bool {
	switch {
	case g[0] != 126:
		g = g[1:]
	case g[1] != 126:
		g = g[4:]
	default:
		g = g[8:]
	}
	if bit/6 >= len(g) {
		panic("g6: index out of range")
	}
	return (g[bit/6]-63)&(1<<uint(5-bit%6)) != 0
}

func (g Graph) GoString() string {
	bin, m6 := binary(g)
	format := fmt.Sprintf("%%d:%%0%db", m6)
	return fmt.Sprintf(format, numberOf(g), bin)
}

func binary(g Graph) (b *big.Int, l int) {
	switch {
	case g[0] != 126:
		g = g[1:]
	case g[1] != 126:
		g = g[4:]
	default:
		g = g[8:]
	}
	b = &big.Int{}
	var c big.Int
	for i := range g {
		c.SetUint64(uint64(g[len(g)-i-1] - 63))
		c.Lsh(&c, uint(6*i))
		b.Or(b, &c)
	}
	return b, len(g) * 6
}
