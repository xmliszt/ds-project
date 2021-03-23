package rpc

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/data"
)

// Node contains all the variables that are necessary to manage a node
type Node struct {
	IsCoordinator       *bool                   `validate:"required"`
	Pid                 int                     `validate:"gte=0"`    // Node ID
	Ring                []int                   `validate:"required"` // Ring structure of nodes
	RecvChannel         chan *data.Data         `validate:"required"` // Receiving channel
	SendChannel         chan *data.Data         `validate:"required"` // Sending channel
	RpcMap              map[int]chan *data.Data `validate:"required"` // Map node ID to their receiving channels
	HeartBeatTable      map[int]bool            // Heartbeat table
	VirtualNodeLocation []int
	VirtualNodeMap      map[int]string
}

// HandleMessageReceived is a Go routine that handles the messages received
func (n *Node) HandleMessageReceived() {

	// Test a dead node
	if n.Pid == 5 {
		go func() {
			time.Sleep(time.Second * 30)
			defer close(n.RecvChannel)
		}()
	}

	for msg := range n.RecvChannel {
		switch msg.Payload["type"] {
		case "CHECK_HEARTBEAT":
			n.SendSignal(0, &data.Data{
				From: n.Pid,
				To:   0,
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
		case "BROADCAST_VIRTUAL_NODES":
			location := msg.Payload["locationData"]
			virtualNode := msg.Payload["virtualNodeData"]
			n.VirtualNodeLocation = location.([]int)
			n.VirtualNodeMap = virtualNode.(map[int]string)
		}
	}
}

// Create virtual nodes
func (n *Node) CreateVirtualNodes(Pid int) error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}

	for i := 1; i <= conf.VirtualNodesCount; i++ {
		virtualNode := strconv.Itoa(Pid) + "-" + strconv.Itoa(i)
		location := config.GetHash(virtualNode)
		fmt.Println("Virtual node ", virtualNode, "has started")
		n.SendSignal(0, &data.Data{
			From: Pid,
			To:   0,
			Payload: map[string]interface{}{
				"type":            "UPDATE_VIRTUAL_NODE",
				"virtualNodeData": virtualNode,
				"locationData":    location,
			},
		})
	}
	return nil
}

// Start starts up a node, running receiving channel
func (n *Node) Start() {
	fmt.Printf("Node [%d] has started!\n", n.Pid)

	// Create virtual node
	n.CreateVirtualNodes(n.Pid)
	go n.HandleMessageReceived()
}

// Start starts up a node, running receiving channel
func (n *Node) StartDeadNode() {
	fmt.Printf("Node [%d] has started!\n", n.Pid)
	go n.HandleMessageReceived()
}

// TearDown terminates node, closes all channels
func (n *Node) TearDown() {
	close(n.RecvChannel)
	close(n.SendChannel)
	fmt.Printf("Node [%d] has terminated!\n", n.Pid)
}
