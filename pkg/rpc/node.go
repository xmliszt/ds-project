package rpc

import (
	"fmt"
	"math/rand"
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
// Go Routine to handle the messages received
func (n *Node) HandleMessageReceived() {
	for {
		select {
		case msg, ok := <-n.RecvChannel:
			if ok {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
				fmt.Printf("Node [%d] receive: %v\n", n.Pid, msg)
				if msg.Payload["Hello"] == "world" {
					n.SendSignal(0, Data{
						From: n.Pid,
						To: 0,
						Payload: map[string]interface{}{
							"data": fmt.Sprintf("Hi there! Greeting from Node [%d]", n.Pid),
						},
					})
				}
			} else {
				continue
			}
		default:
			continue
		}
	}
	// for msg := range n.RecvChannel {
	// 	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	// 	fmt.Printf("Node [%d] receive: %v\n", n.Pid, msg)
	// 	if msg.Payload["Hello"] == "world" {
	// 		n.SendSignal(0, Data{
	// 			From: n.Pid,
	// 			To: 0,
	// 			Payload: map[string]interface{}{
	// 				"data": fmt.Sprintf("Hi there! Greeting from Node [%d]", n.Pid),
	// 			},
	// 		})
	// 	}
	// }
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
	fmt.Printf("Node [%d] has terminated!\n", n.Pid)
}