package rpc

import (
	"fmt"
	"time"
)

type Node struct {
	IsCoordinator *bool `validate:"required"`
	Pid int `validate:"gte=0"` 											// Node ID
	Ring []int `validate:"required"`								// Ring structure of nodes
	RecvChannel chan *Data	`validate:"required"`			// Receiving channel
	SendChannel chan *Data `validate:"required"`			// Sending channel
	RpcMap map[int]chan *Data `validate:"required"`	// Map node ID to their receiving channels
	HeartBeatTable map[int]bool // Heartbeat table
}

// green part
// HandleMessageReceived is a Go routine that handles the messages received
func (n *Node) HandleMessageReceived() {
	
	// Test a dead node
	if n.Pid == 5 {
		go func() {
			time.Sleep(time.Second * 12)
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
		case "UPDATE_HEARTBEAT":
			heartbeatTable := msg.Payload["data"]
			n.HeartBeatTable = heartbeatTable.(map[int]bool)
		case "YOU_ARE_COORDINATOR":
			isCoordinator := true
			n.IsCoordinator = &isCoordinator
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