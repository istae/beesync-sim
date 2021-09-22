package network

import (
	"sync"
)

type Trace struct {
	edges []*edge
	mux   sync.Mutex
}

type edge struct {
	from *Node
	to   *Node
	err  bool
}

type traceNode struct {
	networkNode *Node
	children    *[]traceNode
}

func (t *Trace) Add(from, to *Node) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.edges = append(t.edges, &edge{from: from, to: to})
}

func (t *Trace) MarkErr(from, to *Node, err bool) {
	t.mux.Lock()
	defer t.mux.Unlock()
	for _, e := range t.edges {
		if e.from == from && e.to == to {
			e.err = err
		}
	}
}

func (t *Trace) Reset() {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.edges = nil
}

func (t *Trace) Count() int {
	t.mux.Lock()
	defer t.mux.Unlock()
	return len(t.edges)
}
