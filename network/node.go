package network

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/ethersphere/bee/pkg/swarm"
)

const (
	oversaturation  = 20
	depthSaturation = 4
	depth           = 12
)

var ErrNode = errors.New("node fail")

type Node struct {
	overlay   swarm.Address
	bins      [][]*Node
	fail      bool
	trace     *trace
	handeFunc HandleFunc
}

func NewNode(t *trace, h HandleFunc, fail bool) *Node {
	return &Node{
		overlay:   RandAddress(),
		bins:      make([][]*Node, swarm.MaxBins),
		trace:     t,
		handeFunc: h,
		fail:      fail,
	}
}

func RandAddress() swarm.Address {
	b := make([]byte, 32)
	rand.Read(b)
	return swarm.NewAddress(b)
}

func (n *Node) Add(peers []*Node) {
	for _, peer := range peers {
		if n.Depth() >= 12 {
			return
		}
		po := swarm.Proximity(n.overlay.Bytes(), peer.overlay.Bytes())
		if len(n.bins[po]) < oversaturation && !n.overlay.Equal(peer.overlay) {
			n.bins[po] = append(n.bins[po], peer)
		}
	}
}

func (n *Node) Depth() int {
	for i, bin := range n.bins {
		if len(bin) <= depthSaturation {
			return i
		}
	}

	return int(swarm.MaxPO)
}

func (n *Node) Push(addr swarm.Address) error {

	n.trace.Add(n)
	fmt.Println(n.overlay)

	if n.fail {
		return ErrNode
	}

	return n.handeFunc(addr, n)
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
