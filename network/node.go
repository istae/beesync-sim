package network

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/ethersphere/bee/pkg/swarm/test"
)

var NodeError = errors.New("node fail")

type Node struct {
	overlay swarm.Address
	bins    [][]*Node
	fail    bool
	trace   *trace
}

func NewNode(t *trace) *Node {
	return &Node{
		overlay: test.RandomAddress(),
		bins:    make([][]*Node, swarm.MaxBins),
		trace:   t,
	}
}

func (n *Node) Add(peers ...*Node) {
	for _, peer := range peers {
		po := swarm.Proximity(n.overlay.Bytes(), peer.overlay.Bytes())
		if len(n.bins[po]) < oversaturation && !n.overlay.Equal(peer.overlay) {
			n.bins[po] = append(n.bins[po], peer)
		}
	}
}

func (n *Node) Push(addr swarm.Address) error {

	n.trace.Add(n)
	fmt.Println(n.overlay)

	if n.fail {
		return NodeError
	}

	return defaultPushHandleFunc(addr, n)
}

func defaultPushHandleFunc(addr swarm.Address, base *Node) error {
	closest := base.ClosestNode(addr)
	if closest == base {
		return nil
	}

	return closest.Push(addr)
}

func (n *Node) ClosestNode(addr swarm.Address, skipNodes ...swarm.Address) *Node {

	var closest *Node

	for bin := range n.bins {
		for _, node := range n.bins[bin] {

			skip := false
			for _, skipNode := range skipNodes {
				if skipNode.Equal(node.overlay) {
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			if closest == nil {
				closest = node
			} else if closer(addr, node.overlay, closest.overlay) {
				closest = node
			}
		}
	}

	if closer(addr, n.overlay, closest.overlay) {
		return n
	}

	return closest
}

func closer(a, x, y swarm.Address) bool {
	cmp, _ := swarm.DistanceCmp(a.Bytes(), x.Bytes(), y.Bytes())
	return cmp == 1
}

func rndSubset(nodes []*Node, count int) []*Node {
	if count >= len(nodes) {
		return nodes
	}
	for i := 0; i < len(nodes); i++ {
		j := rand.Intn(len(nodes))
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
	return nodes[:count]
}
