package node

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/api"
)

// Node contains all the variables that are necessary to manage a node
type Node struct {
	IsCoordinator       bool           `validate:"required"`
	Pid                 int            `validate:"gte=0"`    // Node ID
	Ring                []int          `validate:"required"` // Ring structure of nodes
	RpcMap              map[int]string `validate:"required"` // Map node ID to their receiving address
	HeartBeatTable      map[int]bool
	VirtualNodeLocation []int
	VirtualNodeMap      map[int]string
	Router              *echo.Echo
}

// Start is the main function that starts the entire program
func Start(nodeID int) {

	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", config.ConfigLocksmith.Port+nodeID))
	if err != nil {
		log.Fatal(err)
	}
	inbound, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	routerBuilder := &api.RouterBuilder{}

	node := &Node{
		IsCoordinator:       false,
		Pid:                 nodeID,
		Ring:                make([]int, 0),
		RpcMap:              make(map[int]string),
		VirtualNodeLocation: make([]int, 0),
		VirtualNodeMap:      make(map[int]string),
		HeartBeatTable:      make(map[int]bool),
		Router:              routerBuilder.New(),
	}

	// Start RPC server
	log.Printf("Node %d listening on: %v\n", node.Pid, address)
	rpc.Register(node)
	rpc.Accept(inbound)
}

// // HandleMessageReceived is a Go routine that handles the messages received
// func (n *Node) HandleMessageReceived() {

// 	for msg := range n.RecvChannel {
// 		switch msg.Payload["type"] {
// 		case "YOU_ARE_COORDINATOR":
// 			isCoordinator := true
// 			n.IsCoordinator = &isCoordinator
// 			log.Printf("Node %d is the coordinator now!\n", n.Pid)
// 			n.StartRouter()
// 		case "BROADCAST_VIRTUAL_NODES":
// 			location := msg.Payload["locationData"]
// 			virtualNode := msg.Payload["virtualNodeData"]
// 			n.VirtualNodeLocation = location.([]int)
// 			n.VirtualNodeMap = virtualNode.(map[int]string)
// 		}
// 	}
// }

// // Create virtual nodes
// func (n *Node) CreateVirtualNodes(Pid int) error {
// 	conf, err := config.GetConfig()
// 	if err != nil {
// 		return err
// 	}

// 	for i := 1; i <= conf.VirtualNodesCount; i++ {
// 		virtualNode := strconv.Itoa(Pid) + "-" + strconv.Itoa(i)
// 		location, e := config.GetHash(virtualNode)
// 		if e != nil {
// 			return e
// 		}
// 		log.Println("Virtual node ", virtualNode, "has started")
// 		n.SendSignal(0, &data.Data{
// 			From: Pid,
// 			To:   0,
// 			Payload: map[string]interface{}{
// 				"type":            "UPDATE_VIRTUAL_NODE",
// 				"virtualNodeData": virtualNode,
// 				"locationData":    location,
// 			},
// 		})
// 	}
// 	return nil
// }

// // Start starts up a node, running receiving channel
// func (n *Node) Start() error {
// 	log.Printf("Node [%d] has started!\n", n.Pid)

// 	// Create virtual node
// 	err := n.CreateVirtualNodes(n.Pid)
// 	if err != nil {
// 		return err
// 	}
// 	go n.HandleMessageReceived()
// 	return nil
// }

// Starts the router
func (n *Node) startRouter() {
	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	log.Printf("Node %d listening to client's requests...\n", n.Pid)
	go func() {
		err := n.Router.Start(fmt.Sprintf(":%d", config.ConfigServer.Port))
		if err != nil {
			log.Printf("Node %d REST server closed!\n", n.Pid)
		}
	}()
}

// Shutdown the router
func (n *Node) stopRouter() {
	err := n.Router.Close()
	if err != nil {
		panic(err)
	}
}
