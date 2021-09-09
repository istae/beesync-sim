package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/ethersphere/bee/pkg/swarm/test"

	"sim/network"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	net := network.NewNetwork(10000, 5000)

	net.RandNode().Push(test.RandomAddress())

	ioutil.WriteFile("vis/trace-data.js", []byte(fmt.Sprintf(`trace = '%s'`, net.MarshallTrace())), os.ModePerm)
	ioutil.WriteFile("vis/network-data.js", []byte(fmt.Sprintf(`network = '%s'`, net.MarshallNetwork())), os.ModePerm)
}
