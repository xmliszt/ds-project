package rpc

import "fmt"

// Node contains all the variables that are necessary to manage a node
type Node struct {
	IsCoordinator bool
	Pid           int                                 // Node ID
	Ring          []int                               // Ring structure of nodes
	RecvChannel   chan map[string]interface{}         // Receiving channel
	SendChannel   chan map[string]interface{}         // Sending channel
	RPCMap        map[int]chan map[string]interface{} // Map node ID to their receiving channels
}

// HandleMessageReceived run as Go Routine to handle the messages received
func (n *Node) HandleMessageReceived() {
	for msg := range n.RecvChannel {
		fmt.Println("I receive: ", msg)
	}
}

// Start is the initializing function for the node
// func (n *Node) Start() error {
// 	// do whatever thing required for each node to start
// 	go n.HandleMessageReceived()
// }
