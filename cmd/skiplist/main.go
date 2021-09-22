package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/ethersphere/bee/pkg/swarm"

	"sim/network"
)

// func defaultPushHandle(addr swarm.Address, base *network.Node) error {
// 	_, closest := base.ClosestNode(addr)
// 	if closest == base {
// 		return nil
// 	}

// 	return closest.Push(addr)
// }

func main() {

	rand.Seed(time.Now().UnixNano())

	f, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	trace := &network.Trace{}

	list := make(map[string]*peerSkipList)

	skiplist := func(addr swarm.Address) *peerSkipList {
		l := list[addr.ByteString()]
		if l == nil {
			l = &peerSkipList{
				skip: make(map[string]map[string]time.Time),
			}
			list[addr.ByteString()] = l
		}
		return l
	}

	h := func(retry int) func(swarm.Address, *network.Node) (*network.Node, error) {
		return func(ch swarm.Address, base *network.Node) (*network.Node, error) {

			l := skiplist(base.Addr())

			var skipPeers []swarm.Address

			for i := 0; i < retry; i++ {

				_, closest := base.ClosestNode(ch, append(skipPeers, l.ChunkSkipPeers(ch)...)...)
				if closest == base {
					return base, nil
				}

				skipPeers = append(skipPeers, closest.Addr())

				trace.Add(base, closest)

				n, err := closest.Push(ch)
				if err != nil {
					l.Add(ch, closest.Addr(), time.Minute*30)
					trace.MarkErr(base, closest, true)
				} else {
					return n, nil
				}
			}

			return nil, errors.New("no push")

		}
	}

	net := network.NewNetwork(50000, trace, network.NodeOptions{
		NodeConnections: 50000,
		PushHandle:      h(3),
	})

	net.SetFailureRate(0.3)

	count := 1000

	rndNode := net.RandNode(false)
	rndNode.SetHandleFunc(h(3))

	for i := 0; i < count; i++ {

		chunk := network.RandAddress()

		storer, err := rndNode.Push(chunk)
		if err == nil {

			swarm.Proximity(storer.Addr().Bytes(), chunk.Bytes())

			fmt.Printf("from %s to %s bin %d\n", rndNode.Prefix(), storer.Prefix(), swarm.Proximity(storer.Addr().Bytes(), chunk.Bytes()))
		}

		ioutil.WriteFile("vis/trace-data.js", []byte(fmt.Sprintf(`trace = '%s'`, net.MarshallTraceEdges(chunk))), os.ModePerm)
		ioutil.WriteFile("vis/network-data.js", []byte(fmt.Sprintf(`network = '%s'`, net.MarshallNetwork())), os.ModePerm)

		trace.Reset()

		var b []byte = make([]byte, 1)
		os.Stdin.Read(b)
	}
}

type peerSkipList struct {
	sync.Mutex

	// key is chunk address, value is map of peer address to expiration
	skip map[string]map[string]time.Time
}

func (l *peerSkipList) Add(chunk, peer swarm.Address, expire time.Duration) {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.skip[chunk.ByteString()]; !ok {
		l.skip[chunk.ByteString()] = make(map[string]time.Time)
	}
	l.skip[chunk.ByteString()][peer.ByteString()] = time.Now().Add(expire)
}

func (l *peerSkipList) ChunkSkipPeers(ch swarm.Address) (peers []swarm.Address) {
	l.Lock()
	defer l.Unlock()

	if p, ok := l.skip[ch.ByteString()]; ok {
		for peer, exp := range p {
			if time.Now().Before(exp) {
				peers = append(peers, swarm.NewAddress([]byte(peer)))
			}
		}
	}
	return peers
}

func (l *peerSkipList) PruneChunk(chunk swarm.Address) {
	l.Lock()
	defer l.Unlock()
	delete(l.skip, chunk.ByteString())
}

func (l *peerSkipList) PruneExpired() {
	l.Lock()
	defer l.Unlock()

	now := time.Now()

	for k, v := range l.skip {
		kc := len(v)
		for kk, vv := range v {
			if vv.Before(now) {
				delete(v, kk)
				kc--
			}
		}
		if kc == 0 {
			// prune the chunk too
			delete(l.skip, k)
		}
	}
}
