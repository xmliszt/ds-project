package rpc

import "fmt"

type Node struct {
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
}