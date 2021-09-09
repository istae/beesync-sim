package network

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/ethersphere/bee/pkg/swarm"
)

const (
	prefix = 6
)

type Network struct {
	nodes []*Node
	ids   map[string]int
	trace *trace
}

type HandleFunc func(swarm.Address, *Node) error

type NodeOptions struct {
	PushHandle      HandleFunc
	NodeConnections int
	FailPercantage  float32
}

func NewNetwork(count int, o NodeOptions) *Network {

	t := &trace{}

	nodes := nodeBatch(count, t, o)
	ids := make(map[string]int)

	for i := 0; i < count; i++ {
		end := rand.Intn(len(nodes)-o.NodeConnections) + o.NodeConnections
		nodes[i].Add(nodes[end-o.NodeConnections : end])
		ids[nodes[i].overlay.ByteString()] = i
	}

	return &Network{
		nodes: nodes,
		ids:   ids,
		trace: t,
	}
}

func (n *Network) RandNode() *Node {
	return n.nodes[rand.Intn(len(n.nodes))]
}

func nodeBatch(count int, t *trace, o NodeOptions) []*Node {
	var ret = make([]*Node, 0, count)
	for i := 0; i < count; i++ {
		ret = append(ret, NewNode(t, o.PushHandle, false))
	}
	return ret
}

func (net *Network) MarshallNetwork() []byte {

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

func (net *Network) MarshallTrace() []byte {

	var ret []jTrace

	for i := 1; i < len(net.trace.nodes); i++ {
		from := net.trace.nodes[i-1]
		to := net.trace.nodes[i]

		ret = append(ret, jTrace{
			From:   net.ids[from.overlay.ByteString()],
			To:     net.ids[to.overlay.ByteString()],
			Arrows: "to",
			Label:  fmt.Sprintf("%d", i),
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
