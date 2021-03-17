package rpc

import (
	"fmt"
	"time"
)

// Node contains all the variables that are necessary to manage a node
type Node struct {
<<<<<<< HEAD
	IsCoordinator bool
	Pid           int                                 // Node ID
	Ring          []int                               // Ring structure of nodes
	RecvChannel   chan map[string]interface{}         // Receiving channel
	SendChannel   chan map[string]interface{}         // Sending channel
	RPCMap        map[int]chan map[string]interface{} // Map node ID to their receiving channels
}

// HandleMessageReceived run as Go Routine to handle the messages received
=======
	IsCoordinator *bool `validate:"required"`
	Pid int `validate:"gte=0"` 											// Node ID
	Ring []int `validate:"required"`								// Ring structure of nodes
	RecvChannel chan *Data	`validate:"required"`			// Receiving channel
	SendChannel chan *Data `validate:"required"`			// Sending channel
	RpcMap map[int]chan *Data `validate:"required"`	// Map node ID to their receiving channels
}

// green part
// HandleMessageReceived is a Go routine that handles the messages received
>>>>>>> dev
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
			n.SendSignal(0, &Data{
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

<<<<<<< HEAD
// Start is the initializing function for the node
// func (n *Node) Start() error {
// 	// do whatever thing required for each node to start
// 	go n.HandleMessageReceived()
// }
=======
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
>>>>>>> dev
