package node

import (
	"log"

	"github.com/xmliszt/e-safe/pkg/message"
)

// UpdateHeartbeatTable updates the Heartbeat Table that the node has
func (n *Node) UpdateHeartbeatTable(request *message.Request, reply *message.Reply) error {
	n.HeartBeatTable = request.Payload.(map[int]bool)
	log.Printf("Node %d's Heartbeat Table is updated: %v", n.Pid, n.HeartBeatTable)
	return nil
}

// AssignCoordinator assigns the current node to be the coordinator and starts the router
func (n *Node) AssignCoordinator(request *message.Request, reply *message.Reply) error {
	n.IsCoordinator = true
	n.startRouter()
	return nil
}
