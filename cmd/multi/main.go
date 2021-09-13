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

type hop struct {
	bin []int
}

func main() {

	now := time.Now()

	rand.Seed(time.Now().UnixNano())

	f, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	t := &network.Trace{}

	hopBins := make([]float32, 10)
	hopCount := make([]float32, 10)

	var avgHops float32

	net := network.NewNetwork(100000, t, network.NodeOptions{
		NodeConnections: 50000,
		FailPercantage:  0,
		PushHandle: func(addr swarm.Address, node *network.Node) error {

			bin, closest := node.ClosestNode(addr)
			if closest == node {
				return nil
			}

			hopBins[t.Count()-1] += float32(bin)
			hopCount[t.Count()-1]++

			return closest.Push(addr)
		},
	})

	for i := 0; i < 1000; i++ {
		chunk := network.RandAddress()

		rndNode := net.RandNode()

		err := rndNode.Push(chunk)
		if err != nil {
			fmt.Println(err)
		}

		avgHops += float32(t.Count()) - 1
		t.Reset()
	}

	for i, h := range hopBins {
		if hopCount[i] > 0 {
			fmt.Printf("hop %d avg bin %.2f\n", i, h/hopCount[i])
		}
	}

	fmt.Printf("avg hops %.2f\n", avgHops/1000)

	fmt.Printf("in %v\n", time.Since(now))
}
