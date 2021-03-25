package rpc

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/data"
	"github.com/xmliszt/e-safe/pkg/file"
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
		case "STORE_DATA":
			receivedPayload := msg.Payload["data"]                       // This should be the Secret
			hashedValue := fmt.Sprintf("%v", msg.Payload["hashedValue"]) // This should be the hashed value omt for that secret
			mapPayload := map[string]interface{}{
				hashedValue: receivedPayload,
			}
			file.WriteDataFile(n.Pid, mapPayload)
			// How are we going from n.Pid to next_id
			// hashedValue -> current virtual node number (1-1)
			// using Ring, current virtual node number -> next virtual node number
			var next_pid int
			for index, x := range n.Ring {

				hashedValueINT, err := strconv.Atoi(hashedValue)
				if err != nil {
					fmt.Println(err)
				}
				if hashedValueINT < x {
					// current_virtual_node := n.RingMap[x]
					next_virtual_node := n.VirtualNodeMap[n.Ring[(index+1)]]
					string_list := strings.Split(next_virtual_node, "-")
					next_pid, err = strconv.Atoi(string_list[0])
					if err != nil {
						fmt.Println(err)
					}
					break

				}

			}
			n.SendSignal(next_pid, &data.Data{
				From: n.Pid,
				To:   next_pid,
				Payload: map[string]interface{}{
					"type": "STRICT_CONSISTENCY",
					"data": mapPayload,
				},
			})
		case "STRICT_CONSISTENCY":

			//SendSignal for ack to owner node
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
		location, e := config.GetHash(virtualNode)
		if e != nil {
			return e
		}
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
func (n *Node) Start() error {
	fmt.Printf("Node [%d] has started!\n", n.Pid)

	// Create virtual node
	err := n.CreateVirtualNodes(n.Pid)
	if err != nil {
		return err
	}
	go n.HandleMessageReceived()
	return nil
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
