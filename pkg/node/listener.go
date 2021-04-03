package node

import (
	"log"

	"github.com/xmliszt/e-safe/pkg/file"
	"github.com/xmliszt/e-safe/pkg/message"
	"github.com/xmliszt/e-safe/pkg/secret"
	"github.com/xmliszt/e-safe/util"
)

// UpdateRpcMap updates node's RPC Map
func (n *Node) UpdateRpcMap(request *message.Request, reply *message.Reply) error {
	n.RpcMap = request.Payload.(map[int]string)
	log.Printf("Node %d's RPC Map are updated: %+v\n", n.Pid, n.RpcMap)
	return nil
}

// UpdateHeartbeatTable updates the Heartbeat Table that the node has
func (n *Node) UpdateHeartbeatTable(request *message.Request, reply *message.Reply) error {
	n.HeartBeatTable = request.Payload.(map[int]bool)
	log.Printf("Node %d's Heartbeat Table is updated: %v", n.Pid, n.HeartBeatTable)
	return nil
}

// AssignCoordinator assigns the current node to be the coordinator and starts the router
func (n *Node) AssignCoordinator(request *message.Request, reply *message.Reply) error {
	log.Printf("Node %d is the new coordinator!\n", n.Pid)
	n.IsCoordinator = true
	go n.startRouter()
	return nil
}

// RemoveCoordinator removes the coordinator flag from this node and stop its router
func (n *Node) RemoveCoordinator(request *message.Request, reply *message.Reply) error {
	log.Printf("Node %d is no longer the coordinator!\n", n.Pid)
	n.IsCoordinator = false
	n.stopRouter()
	return nil
}

// UpdateVirtualNodes updates the node's virtual node location and map
func (n *Node) UpdateVirtualNodes(request *message.Request, reply *message.Reply) error {
	payload := request.Payload.(map[string]interface{})
	locations := payload["virtualNodeLocation"]
	virtualNode := payload["virtualNodeMap"]
	n.VirtualNodeLocation = locations.([]int)
	n.VirtualNodeMap = virtualNode.(map[int]string)
	log.Printf("Node %d updated virtual nodes: %v | %+v\n", n.Pid, n.VirtualNodeLocation, n.VirtualNodeMap)
	return nil
}

// receives the message from Coordinator and do what is on the board
// REMEMBER to capitalize the function name
// Strict Consistency with R = 2. Send ACK directly to coordinator
func (n *Node) OwnerNodeDown(request *message.Request, reply *message.Reply) error {
	// func (n *Node) StrictDown(replicationList []string){
	log.Printf("Owner Node down. ")
	// Coordintor send message to
	return nil
}

func (n *Node) StrictReplication(request *message.Request, reply *message.Reply) error {
	log.Printf("Begin strict replication")

	// parse secret
	payload := request.Payload.(map[string]interface{})
	hashedValue := payload["hashedValue"].(string)
	secret := payload["data"].(secret.Secret)

	dataToWrite := map[string]interface{}{
		hashedValue: secret,
	}

	// Write to respective node storage file
	writeErr := file.WriteDataFile(n.Pid, dataToWrite)
	if writeErr != nil {
		log.Fatal("Data file write failed for node %d", n.Pid)
	}

	// nextVNode = util.MapHashToVNode()

	nextVNodePid := util.FindNextVNode(n.Ring, n.VirtualNodeMap, n.VirtualNodeLocation, hashedValue)
	vNodeActualPid := util.NodePidFromVNode(nextVNodePid)
	// Check if the next node is alive
	if n.checkHeartbeat(vNodeActualPid) {

	}
	// Coordintor send message to
	return nil
}
