package network

import (
	"errors"
	"math/rand"

	"github.com/ethersphere/bee/pkg/swarm"
)

const (
	oversaturation  = 20
	depthSaturation = 4
	depth           = 14
)

var ErrNode = errors.New("node fail")

type Node struct {
	overlay   swarm.Address
	bins      [][]*Node
	fail      bool
	trace     *Trace
	handeFunc HandleFunc
}

func NewNode(t *Trace, h HandleFunc) *Node {
	return &Node{
		overlay:   RandAddress(),
		bins:      make([][]*Node, swarm.MaxBins),
		trace:     t,
		handeFunc: h,
	}
}

func RandAddress() swarm.Address {
	b := make([]byte, 32)
	rand.Read(b)
	return swarm.NewAddress(b)
}

func (n *Node) SetFailureRate(failPercantage float64) {
	n.fail = rand.Float64() < failPercantage
}

func (n *Node) Add(peers []*Node) {
	for _, peer := range peers {
		po := swarm.Proximity(n.overlay.Bytes(), peer.overlay.Bytes())
		if len(n.bins[po]) < oversaturation && !n.overlay.Equal(peer.overlay) {
			n.bins[po] = append(n.bins[po], peer)
		}
	}
}

func (n *Node) Depth() int {
	for i, bin := range n.bins {
		if len(bin) < depthSaturation {
			return i
		}
	}

	return int(swarm.MaxPO)
}

func (n *Node) SetHandleFunc(h HandleFunc) {
	n.handeFunc = h
}

func (n *Node) Deepest() int {
	for i := len(n.bins) - 1; i >= 0; i-- {
		if len(n.bins[i]) >= 1 {
			return i
		}
	}

	return 0
}

func (n *Node) Prefix() string {
	return n.Addr().String()[:prefix]
}

func (n *Node) Push(ch swarm.Address) (*Node, error) {

	if n.fail {
		return nil, ErrNode
	}

	return n.handeFunc(ch, n)
}

func (n *Node) ClosestNode(addr swarm.Address, skipNodes ...swarm.Address) (int, *Node) {

	var closest *Node
	var bin int

	for b := range n.bins {
		for _, node := range n.bins[b] {

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

			if closest == nil || closer(addr, node.overlay, closest.overlay) {
				closest = node
				bin = b
			}
		}
	}

	if closer(addr, n.overlay, closest.overlay) {
		return 256, n
	}

	return bin, closest
}

func (n *Node) Addr() swarm.Address {
	return n.overlay
}

func closer(a, x, y swarm.Address) bool {
	cmp, _ := swarm.DistanceCmp(a.Bytes(), x.Bytes(), y.Bytes())
	return cmp == 1
}
