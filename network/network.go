package network

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"

	"github.com/ethersphere/bee/pkg/swarm"
)

const (
	prefix = 6
)

type Network struct {
	nodes []*Node
	ids   map[string]int
	trace *Trace
}

type HandleFunc func(swarm.Address, *Node) error

type NodeOptions struct {
	PushHandle      HandleFunc
	NodeConnections int
	FailPercantage  float32
}

func NewNetwork(count int, t *Trace, o NodeOptions) *Network {

	nodes := nodeBatch(count, t, o)
	ids := make(map[string]int)

	var mux sync.Mutex
	var wg sync.WaitGroup

	const parallel = 8
	window := count / parallel

	for i := 0; i < parallel; i++ {
		wg.Add(1)

		start := window * i
		go func(index int, end int) {

			defer wg.Done()

			for i := index; i < end; i++ {
				end := rand.Intn(len(nodes) - o.NodeConnections)
				nodes[i].Add(nodes[end : end+o.NodeConnections])

				mux.Lock()
				ids[nodes[i].overlay.ByteString()] = i
				mux.Unlock()
			}

		}(start, start+window)
	}

	wg.Wait()

	return &Network{
		nodes: nodes,
		ids:   ids,
		trace: t,
	}
}

func (n *Network) RandNode() *Node {
	return n.nodes[rand.Intn(len(n.nodes))]
}

func nodeBatch(count int, t *Trace, o NodeOptions) []*Node {
	var ret = make([]*Node, 0, count)
	for i := 0; i < count; i++ {
		ret = append(ret, NewNode(t, o.PushHandle, false))
	}
	return ret
}

func (net *Network) MarshallTraceNode() []byte {

	type jnode struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	}

	var j []jnode

	for _, node := range net.trace.nodes {
		j = append(j, jnode{
			ID:    net.ids[node.overlay.ByteString()],
			Label: node.overlay.String()[:prefix],
		})
	}

	b, _ := json.Marshal(j)
	return b
}

type jTrace struct {
	From   int    `json:"from"`
	To     int    `json:"to"`
	Arrows string `json:"arrows,omitempty"`
	Label  string `json:"label,omitempty"`
}

func (net *Network) MarshallTraceEdges() []byte {

	var ret []jTrace

	for i := 1; i < len(net.trace.nodes); i++ {
		from := net.trace.nodes[i-1]
		to := net.trace.nodes[i]

		ret = append(ret, jTrace{
			From:   net.ids[from.overlay.ByteString()],
			To:     net.ids[to.overlay.ByteString()],
			Arrows: "to",
			Label:  fmt.Sprintf("%d (bin %d)", i, swarm.Proximity(from.overlay.Bytes(), to.overlay.Bytes())),
		})
	}

	b, _ := json.Marshal(ret)
	return b
}

// func (n *network) MarshallConnections(nodes []*node) []byte {

// 	var ret []jTrace

// 	for i := 1; i < len(nodes); i++ {
// 		from := nodes[i-1]

// 		for bin := range from.bins {
// 			for _, to := range from.bins[bin] {
// 				ret = append(ret, jTrace{
// 					From: n.ids[from.overlay.ByteString()],
// 					To:   n.ids[to.overlay.ByteString()],
// 				})
// 			}
// 		}
// 	}

// 	b, _ := json.Marshal(ret)
// 	return b
// }
