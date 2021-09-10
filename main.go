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

func defaultPushHandleFunc(addr swarm.Address, base *network.Node) error {
	closest := base.ClosestNode(addr)
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

	net := network.NewNetwork(100000, network.NodeOptions{
		NodeConnections: 50000,
		FailPercantage:  0,
		PushHandle:      defaultPushHandleFunc,
	})

	chunk := network.RandAddress()
	fmt.Println(chunk)

	rndNode := net.RandNode()

	err := rndNode.Push(chunk)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("depth: %d\n", rndNode.Depth())
	fmt.Printf("in %v\n", time.Since(now))

	ioutil.WriteFile("vis/trace-data.js", []byte(fmt.Sprintf(`trace = '%s'`, net.MarshallTrace())), os.ModePerm)
	ioutil.WriteFile("vis/network-data.js", []byte(fmt.Sprintf(`network = '%s'`, net.MarshallNetwork())), os.ModePerm)
}
