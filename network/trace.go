package network

import "sync"

type Trace struct {
	nodes []*Node
	mux   sync.Mutex
}

func (t *Trace) Add(n *Node) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.nodes = append(t.nodes, n)
}

func (t *Trace) Reset() {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.nodes = nil
}
