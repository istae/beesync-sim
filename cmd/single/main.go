package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/ethersphere/bee/pkg/swarm"

	"sim/network"
)

func defaultPushHandle(addr swarm.Address, base *network.Node) error {
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
		NodeConnections: 499,
		FailPercantage:  0,
		PushHandle:      defaultPushHandle,
	})

	chunk := network.RandAddress()

	rndNode := net.RandNode(false)

	err := rndNode.Push(chunk)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("depth %d\n", rndNode.Depth())
	fmt.Printf("deepest %d\n", rndNode.Deepest())

	fmt.Printf("in %v\n", time.Since(now))

	ioutil.WriteFile("vis/trace-data.js", []byte(fmt.Sprintf(`trace = '%s'`, net.MarshallTraceEdges())), os.ModePerm)
	ioutil.WriteFile("vis/network-data.js", []byte(fmt.Sprintf(`network = '%s'`, net.MarshallTraceNode())), os.ModePerm)
}
