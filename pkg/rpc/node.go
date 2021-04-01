package rpc

import (
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/data"
)

// Node contains all the variables that are necessary to manage a node
type Node struct {
	IsCoordinator       *bool                   `validate:"required"`
	Pid                 int                     `validate:"gte=0"`    // Node ID
	Ring                []int                   `validate:"required"` // Ring structure of nodes
	RecvChannel         chan *data.Data         `validate:"required"` // Receiving channel
	RpcMap              map[int]chan *data.Data `validate:"required"` // Map node ID to their receiving channels
	HeartBeatTable      map[int]bool
	VirtualNodeLocation []int
	VirtualNodeMap      map[int]string
	Router              *echo.Echo
}

// HandleMessageReceived is a Go routine that handles the messages received
func (n *Node) HandleMessageReceived() {

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
			log.Printf("Node %d is the coordinator now!\n", n.Pid)
			n.StartRouter()
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
		log.Println("Virtual node ", virtualNode, "has started")
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
	log.Printf("Node [%d] has started!\n", n.Pid)

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
	log.Printf("Node [%d] has started!\n", n.Pid)
	go n.HandleMessageReceived()
}

// TearDown terminates node, closes all channels
func (n *Node) TearDown() {
	log.Println(n.RecvChannel)
	close(n.RecvChannel)
	log.Printf("Node [%d] has terminated!\n", n.Pid)
}

// Starts the router
func (n *Node) StartRouter() {
	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	log.Printf("Node %d listening to client's requests...\n", n.Pid)
	go func() {
		err := n.Router.Start(fmt.Sprintf(":%d", config.Port))
		if err != nil {
			log.Printf("Node %d REST server closed!\n", n.Pid)
		}
	}()
}

// Shutdown the router
func (n *Node) StopRouter() {
	err := n.Router.Close()
	if err != nil {
		panic(err)
	}
}
