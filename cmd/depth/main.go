package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/ethersphere/bee/pkg/swarm"

	"sim/network"
)

func defaultPushHandleFunc(addr swarm.Address, base *network.Node) error {
	_, closest := base.ClosestNode(addr)
	if closest == base {
		return nil
	}

	return closest.Push(addr)
}

func main() {

	now := time.Now()

	rand.Seed(time.Now().UnixNano())

	f, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	t := &network.Trace{}

	net := network.NewNetwork(500, t, network.NodeOptions{
		NodeConnections: 200,
		FailPercantage:  0,
		PushHandle:      defaultPushHandleFunc,
	})

	depth := make(map[int]int)

	for i := 0; i < 1000; i++ {

		rndNode := net.RandNode(false)
		depth[rndNode.Depth()]++
	}

	fmt.Println(depth)

	fmt.Printf("in %v\n", time.Since(now))
}
