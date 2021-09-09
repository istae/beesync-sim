package network

import (
	"encoding/json"
	"math/rand"
)

const (
	oversaturation = 8
	prefix         = 16
)

type Network struct {
	nodes []*Node
	ids   map[string]int
	trace *trace
}

type trace struct {
	nodes []*Node
}

func (t *trace) Add(n *Node) {
	t.nodes = append(t.nodes, n)
}

func NewNetwork(count, connections int) *Network {

	t := &trace{}

	nodes := nodeBatch(count, t)
	cpy := make([]*Node, count)
	copy(cpy, nodes)

	ids := make(map[string]int)

	for i := 0; i < count; i++ {
		subset := rndSubset(cpy, connections)
		nodes[i].Add(subset...)
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

func nodeBatch(count int, t *trace) []*Node {
	var ret = make([]*Node, 0, count)
	for i := 0; i < count; i++ {
		ret = append(ret, NewNode(t))
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
