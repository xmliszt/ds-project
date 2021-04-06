package node

import (
	"log"
	"strconv"

	"github.com/xmliszt/e-safe/config"
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
	secretToStore := payload["data"].(secret.Secret)
	relayNodes := payload["relayNodes"].([]int)
	rf := payload["rf"].(int)

	// Write to respective node storage file
	writeErr := secret.PutSecret(n.Pid, hashedValue, &secretToStore)
	if writeErr != nil {
		log.Fatal("Data file write failed for node %d", n.Pid)
	}
	// nextVNode = util.MapHashToVNode()

	// Send eventual replication message to neighbouring nodes
	iHashedValue, _ := strconv.Atoi(hashedValue)
	uHashedValue := uint32(iHashedValue)
	nextVNodePid := util.FindNextVNode(n.Ring, n.VirtualNodeMap, n.VirtualNodeLocation, uHashedValue)
	nextVNodeActualPid := util.NodePidFromVNode(nextVNodePid)
	// Check if the next node is alive
	if n.checkHeartbeat(nextVNodeActualPid) {
		n.sendEventualRepMsg(rf-1, hashedValue, secretToStore, relayNodes)
	} else {
		n.sendEventualRepMsg(rf-2, hashedValue, secretToStore, relayNodes)
	}
	// Reply here
	*reply = message.Reply{
		From:    n.Pid,
		To:      request.From,
		ReplyTo: request.Code,
		Payload: map[string]interface{}{
			"success": true,
		},
	}
	return nil
}

// Performed by rf=1
func (n *Node) PerformEventualReplication(request *message.Request, reply *message.Reply) error {
	payload := request.Payload.(map[string]interface{})
	hashedValue := payload["hashedValue"].(string)
	secretToStore := payload["data"].(secret.Secret)
	relayNodes := payload["relayNodes"].([]int)
	rf := payload["rf"].(int)
	err := secret.PutSecret(n.Pid, hashedValue, &secretToStore)
	if err != nil {
		return err
	}
	if rf > 0 {
		n.sendEventualRepMsg(rf-1, hashedValue, secretToStore, relayNodes)
	}
	return nil
}

// Coordintor send message to
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
	err = n.relaySecretDeletion(replicationFactor, keyToDelete, relayVirtualNodes)
	if err != nil {
		return err
	}

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
		err := n.relaySecretDeletion(replicationFactor, keyToDelete, relayNodes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) PerformStrictDown(request *message.Request, reply *message.Reply) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	replicationFactor := config.ConfigNode.ReplicationFactor
	payload := request.Payload.(map[string]interface{})
	keyToStore := strconv.Itoa(payload["key"].(int))
	// Issue How to load a secret from the payload
	valueToStore := payload["secret"].(secret.Secret)
	// myLocation := payload["location"].(int)
	relayNodes := payload["nodes"].([]int)

	// Store
	err = secret.PutSecret(n.Pid, keyToStore, &valueToStore)
	if err != nil {
		return err
	}

	// Relay Strict Consistency
	err = n.sendEventualRepMsg(replicationFactor, keyToStore, valueToStore, relayNodes)
	if err != nil {
		return err
	} else {
		// ack back to the coord
		*reply = message.Reply{
			From:    n.Pid,
			To:      request.From,
			ReplyTo: request.Code,
			Payload: map[string]interface{}{
				"Strict Down excuted successful": true,
			},
		}
	}
	//log.Printf("Node %d deleted secret [%s] successfully!\n", n.Pid, keyToDelete)
	return nil
}

// Take in list from Coordinator and trigger Store and Replicate processes.
// Store data to itself and send message to next data on the list to conduct strict consistency
func (n *Node) StoreAndReplicate(request *message.Request, reply *message.Reply) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	replicationFactor := config.ConfigNode.ReplicationFactor
	payload := request.Payload.(map[string]interface{})
	keyToStore := strconv.Itoa(payload["key"].(int))
	// Issue How to load a secret from the payload
	valueToStore := payload["secret"].(secret.Secret)
	myLocation := payload["location"].(int)

	// Store
	err = secret.PutSecret(n.Pid, keyToStore, &valueToStore)
	if err != nil {
		return err
	}

	// get the next three locations of replicas
	relayVirtualNodes, err := n.getRelayVirtualNodes(myLocation)
	if err != nil {
		return err
	}

	// Relay Strict Consistency
	err = n.sendStrictRepMsg(replicationFactor, keyToStore, valueToStore, relayVirtualNodes)
	if err != nil {
		return err
	} else {
		// ack back to the coord
		*reply = message.Reply{
			From:    n.Pid,
			To:      request.From,
			ReplyTo: request.Code,
			Payload: map[string]interface{}{
				"success": true,
			},
		}
	}
	//log.Printf("Node %d deleted secret [%s] successfully!\n", n.Pid, keyToDelete)
	return nil
}
