package node

import (
	"log"
	"strconv"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/message"
	"github.com/xmliszt/e-safe/pkg/secret"
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

// DeleteSecret deletes a secret - this is for owner node
func (n *Node) DeleteSecret(request *message.Request, reply *message.Reply) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	replicationFactor := config.ConfigNode.ReplicationFactor
	payload := request.Payload.(map[string]interface{})
	keyToDelete := strconv.Itoa(payload["key"].(int))
	myLocation := payload["location"].(int)
	err = secret.UpdateSecret(n.Pid, keyToDelete, nil)
	if err != nil {
		return err
	}
	// get the next three locations of replicas
	relayVirtualNodes, err := n.getRelayVirtualNodes(myLocation)
	if err != nil {
		return err
	}
	// relay deletion
	go n.relaySecretDeletion(replicationFactor, keyToDelete, relayVirtualNodes)

	*reply = message.Reply{
		From:    n.Pid,
		To:      request.From,
		ReplyTo: request.Code,
		Payload: map[string]interface{}{
			"success": true,
		},
	}

	log.Printf("Node %d deleted secret [%s] successfully!\n", n.Pid, keyToDelete)
	return nil
}

// RelayDeleteSecret deletes a copy of the secret
func (n *Node) RelayDeleteSecret(request *message.Request, reply *message.Reply) error {
	payload := request.Payload.(map[string]interface{})
	replicationFactor := payload["rf"].(int)
	keyToDelete := payload["key"].(string)
	relayNodes := payload["nodes"].([]int)
	err := secret.UpdateSecret(n.Pid, keyToDelete, nil)
	if err != nil {
		return err
	}
	log.Printf("Node %d deleted secret [%s] (replica) successfully!\n", n.Pid, keyToDelete)
	replicationFactor--
	if replicationFactor > 0 {
		go n.relaySecretDeletion(replicationFactor, keyToDelete, relayNodes)
	}
	return nil
}
