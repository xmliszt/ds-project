package rpc

import (
	"fmt"
	"sync"
)

type Node struct {
	Wg *sync.WaitGroup
	IsCoordinator bool
	Pid int											// Node ID
	Ring []int										// Ring structure of nodes
	RecvChannel chan map[string]interface{}			// Receiving channel
	SendChannel chan map[string]interface{} 		// Sending channel
	RpcMap map[int]chan map[string]interface{}		// Map node ID to their receiving channels
}
// green part
// Go Routine to handle the messages received
func (n *Node) HandleMessageReceived() {
	for msg := range n.RecvChannel {
		fmt.Println("I receive: ", msg)
	}
	defer n.TearDown()
}

// Start up a node, running receiving channel
func (n *Node) Start() {
	fmt.Printf("Node [%d] has started!\n", n.Pid)
	go n.HandleMessageReceived()
}

// Terminate node, close all channels
func (n *Node) TearDown() {
	close(n.RecvChannel)
	close(n.SendChannel)
	n.Wg.Done()
	fmt.Printf("Node [%d] has terminated!\n", n.Pid)
}