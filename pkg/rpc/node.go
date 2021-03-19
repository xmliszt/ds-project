package rpc

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/api"
	"github.com/xmliszt/e-safe/pkg/data"
)

// Node contains all the variables that are necessary to manage a node
type Node struct {
	IsCoordinator  bool
	Pid            int                     `validate:"gte=0"`    // Node ID
	Ring           []int                   `validate:"required"` // Ring structure of nodes
	RecvChannel    chan *data.Data         `validate:"required"` // Receiving channel
	SendChannel    chan *data.Data         `validate:"required"` // Sending channel
	RpcMap         map[int]chan *data.Data `validate:"required"` // Map node ID to their receiving channels
	HeartBeatTable map[int]bool            `validate:"required"` // Heartbeat table
	Handler        NodeHandler
}

type NodeHandler interface {
	Start()
	TearDown()
	HandlerMessageReceived()
	HandleAPIRequests()
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
			n.IsCoordinator = true
			fmt.Printf("Node [%d] is assigned to be the new coordinator!\n", n.Pid)
			go n.HandleAPIRequests()
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

// API Requests Handlers
func (n *Node) HandleAPIRequests() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Node %d listening to client's requests...\n", n.Pid)
	router := api.GetRouter()

	err = http.ListenAndServe(":"+fmt.Sprintf("%d", config.ConfigServer.Port), &router)
	if err != nil {
		log.Fatal(err)
	}
}
