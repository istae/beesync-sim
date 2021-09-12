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

	net := network.NewNetwork(100000, t, network.NodeOptions{
		NodeConnections: 50000,
		FailPercantage:  0,
		PushHandle:      defaultPushHandleFunc,
	})

	for i := 0; i < 10; i++ {
		chunk := network.RandAddress()
		fmt.Println(chunk)

		rndNode := net.RandNode()

		err := rndNode.Push(chunk)
		if err != nil {
			fmt.Println(err)
		}
		t.Reset()
	}
	fmt.Printf("in %v\n", time.Since(now))
}
