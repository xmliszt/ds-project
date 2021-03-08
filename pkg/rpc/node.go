package rpc

import (
	"fmt"
	"time"
)

type Node struct {
	IsCoordinator bool
	Pid int											// Node ID
	Ring []int										// Ring structure of nodes
	RecvChannel chan Data		// Receiving channel
	SendChannel chan Data 		// Sending channel
	RpcMap map[int]chan Data		// Map node ID to their receiving channels
}

// green part
// HandleMessageReceived is a Go routine that handles the messages received
func (n *Node) HandleMessageReceived() {
	
	// Test a dead node
	if n.Pid == 3 {
		go func() {
			time.Sleep(time.Second * 50)
			defer close(n.RecvChannel)
		}()
	}

	for msg := range n.RecvChannel {
		switch msg.Payload["type"] {
		case "CHECK_HEARTBEAT":
			n.SendSignal(0, Data{
				From: n.Pid,
				To: 0,
				Payload: map[string]interface{}{
					"type": "REPLY_HEARTBEAT",
					"data": nil,
				},
			})
		}
	}
}

// Start starts up a node, running receiving channel
func (n *Node) Start() {
	fmt.Printf("Node [%d] has started!\n", n.Pid)
	go n.HandleMessageReceived()
}

// TearDown terminates node, closes all channels
func (n *Node) TearDown() {
	close(n.RecvChannel)
	close(n.SendChannel)
	fmt.Printf("Node [%d] has terminated!\n", n.Pid)
}