package network

import "sync"

type trace struct {
	nodes []*Node
	mux   sync.Mutex
}

func (t *trace) Add(n *Node) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.nodes = append(t.nodes, n)
}
