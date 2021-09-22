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

type HandleFunc func(swarm.Address, *Node) (*Node, error)

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

	if count%10 != 0 {
		panic("count must be a multiple of 10")
	}

	const parallel = 10
	window := count / parallel

	N := len(nodes) - o.NodeConnections
	if o.NodeConnections == len(nodes) {
		N = 1
	}

	for i := 0; i < parallel; i++ {

		wg.Add(1)

		start := window * i
		go func(index int, end int) {

			for i := index; i < end; i++ {
				end := rand.Intn(N)
				nodes[i].Add(nodes[end : end+o.NodeConnections])

				mux.Lock()
				ids[nodes[i].overlay.ByteString()] = i
				mux.Unlock()
			}

			wg.Done()
		}(start, start+window)
	}

	wg.Wait()

	return &Network{
		nodes: nodes,
		ids:   ids,
		trace: t,
	}
}

func (n *Network) RandNode(fail bool) *Node {
	for {
		node := n.nodes[rand.Intn(len(n.nodes))]
		if fail == node.fail {
			return node
		}
	}
}

func (n *Network) SetFailureRate(failure float64) {
	for _, n := range n.nodes {
		n.SetFailureRate(failure)
	}
}

func (net *Network) SetHandleFunc(h HandleFunc) {
	for _, n := range net.nodes {
		n.handeFunc = h
	}
}

func nodeBatch(count int, t *Trace, o NodeOptions) []*Node {
	var ret = make([]*Node, 0, count)
	for i := 0; i < count; i++ {
		ret = append(ret, NewNode(t, o.PushHandle))
	}
	return ret
}

func (net *Network) MarshallNetwork() []byte {

	type jnode struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	}

	var j []jnode

	for _, node := range net.trace.edges {

		f := func(overlay swarm.Address) {

			id := net.ids[overlay.ByteString()]

			skip := false
			for _, jn := range j {

				if jn.ID == id {
					skip = true
					break
				}
			}

			if skip {
				return
			}

			j = append(j, jnode{
				ID:    id,
				Label: overlay.String()[:prefix],
			})
		}

		f(node.from.overlay)
		f(node.to.overlay)
	}

	b, _ := json.Marshal(j)
	return b
}

type jTrace struct {
	From   int    `json:"from"`
	To     int    `json:"to"`
	Arrows string `json:"arrows,omitempty"`
	Label  string `json:"label,omitempty"`
	Color  string `json:"color,omitempty"`
}

func (net *Network) MarshallTraceEdges(chunk swarm.Address) []byte {

	var ret []jTrace

	for i := 0; i < len(net.trace.edges); i++ {
		from := net.trace.edges[i].from
		to := net.trace.edges[i].to

		var color = ""
		if net.trace.edges[i].err {
			color = "red"
		}

		po := swarm.Proximity(to.Addr().Bytes(), chunk.Bytes())

		ret = append(ret, jTrace{
			From:   net.ids[from.overlay.ByteString()],
			To:     net.ids[to.overlay.ByteString()],
			Arrows: "to",
			Label:  fmt.Sprintf("%d (po %d)", i, po),
			Color:  color,
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
