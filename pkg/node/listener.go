package node

import (
	"fmt"
	"log"
	"strconv"
	"syscall"

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
	n.KillSignal <- syscall.SIGTERM
	return nil
}

// UpdateVirtualNodes updates the node's virtual node location and map
func (n *Node) UpdateVirtualNodes(request *message.Request, reply *message.Reply) error {
	payload := request.Payload.(map[string]interface{})
	locations := payload["virtualNodeLocation"]
	virtualNode := payload["virtualNodeMap"]
	rpcMap := payload["rpcMap"]
	n.VirtualNodeLocation = locations.([]int)
	n.VirtualNodeMap = virtualNode.(map[int]string)
	n.RpcMap = rpcMap.(map[int]string)
	log.Printf("Node %d updated virtual nodes: %v | %+v\n", n.Pid, n.VirtualNodeLocation, n.VirtualNodeMap)
	log.Printf("Node %d updated rpc map: %+v\n", n.Pid, n.RpcMap)
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
	hashedValue := payload["key"].(string)
	secretToStore := payload["secret"].(secret.Secret)
	relayNodes := payload["nodes"].([]int)
	rf := payload["rf"].(int)

	// Write to respective node storage file
	writeErr := secret.UpdateSecret(n.Pid, hashedValue, &secretToStore)
	if writeErr != nil {
		err := secret.PutSecret(n.Pid, hashedValue, &secretToStore)
		if err != nil {
			return err
		}
	}

	nextVNodeLocation := relayNodes[rf-1]
	nextVNodeName := n.VirtualNodeMap[nextVNodeLocation]
	nextVNodeActualPid, err := getPhysicalNodeID(nextVNodeName)
	if err != nil {
		return err
	}

	// Check if the next node is alive
	if n.checkHeartbeat(nextVNodeActualPid) {
		err := n.sendEventualRepMsg(rf-1, hashedValue, secretToStore, relayNodes)
		if err != nil {
			return err
		}
	} else {
		err := n.sendEventualRepMsg(rf-2, hashedValue, secretToStore, relayNodes)
		if err != nil {
			return err
		}
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
	hashedValue := payload["key"].(string)
	secretToStore := payload["secret"].(secret.Secret)
	relayNodes := payload["nodes"].([]int)
	rf := payload["rf"].(int)
	err := secret.UpdateSecret(n.Pid, hashedValue, &secretToStore)
	if err != nil {
		err := secret.PutSecret(n.Pid, hashedValue, &secretToStore)
		if err != nil {
			return err
		}
	}
	if rf > 1 {
		err := n.sendEventualRepMsg(rf-1, hashedValue, secretToStore, relayNodes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) DeleteSecret(request *message.Request, reply *message.Reply) error {
	payload := request.Payload.(map[string]interface{})
	keyToDelete := strconv.Itoa(payload["key"].(int))

	deletionErr := secret.UpdateSecret(n.Pid, keyToDelete, nil)
	if deletionErr != nil {
		return deletionErr
	}
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

// GetSecrets gets the secrets from given range and send back
func (n *Node) GetSecrets(request *message.Request, reply *message.Reply) error {
	payload := request.Payload.(map[string]interface{})
	fetchRange := payload["range"].([]int)
	toDelete := payload["delete"].(bool)
	from := fetchRange[0]
	to := fetchRange[1]

	secrets, err := secret.GetSecrets(n.Pid, from, to)
	if err != nil {
		return err
	}
	*reply = message.Reply{
		From:    n.Pid,
		To:      request.From,
		ReplyTo: request.Code,
		Payload: secrets,
	}

	if toDelete {
		for key := range secrets {
			err := secret.RemoveSecret(n.Pid, key)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *Node) GetAllSecrets(request *message.Request, reply *message.Reply) error {
	requestPayload := request.Payload.(map[string]interface{})
	role := requestPayload["role"].(int)

	var listOSecrets []*secret.Secret
	listOSecrets, getSecretsError := secret.GetAllNodeSecrets(n.Pid, role)
	if getSecretsError != nil {
		return getSecretsError
	} else {
		*reply = message.Reply{
			From:    n.Pid,
			To:      request.From,
			ReplyTo: request.Code,
			Payload: map[string]interface{}{
				"data": listOSecrets,
			},
		}
	}

	return nil
}

func (n *Node) PerformStrictDown(request *message.Request, reply *message.Reply) error {
	config, err := config.GetConfig()
	fmt.Println("start the strict down")
	if err != nil {
		return err
	}
	replicationFactor := config.ConfigNode.ReplicationFactor
	payload := request.Payload.(map[string]interface{})
	keyToStore := payload["key"].(string)
	// Issue How to load a secret from the payload
	valueToStore := payload["secret"].(secret.Secret)
	// myLocation := payload["location"].(int)
	relayNodes := payload["nodes"].([]int)

	// Store
	err = secret.UpdateSecret(n.Pid, keyToStore, &valueToStore)
	if err != nil {
		err := secret.PutSecret(n.Pid, keyToStore, &valueToStore)
		if err != nil {
			return err
		}
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
				"success": true,
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
	// myLocation := payload["location"].(int)
	relayVirtualNodes := payload["nodes"].([]int)

	// Store
	err = secret.UpdateSecret(n.Pid, keyToStore, &valueToStore)
	if err != nil {
		err := secret.PutSecret(n.Pid, keyToStore, &valueToStore)
		if err != nil {
			return err
		}
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
	return nil
}

func (n *Node) GetData(request *message.Request, reply *message.Reply) error {
	payload := request.Payload.(map[string]interface{})
	keyToSearch := payload["key"].(int)

	retrievedSecret, err := secret.GetSecret(n.Pid, strconv.Itoa(keyToSearch))
	if err != nil {
		log.Println("Unable to retrieve secret from file")
		return err
	}

	*reply = message.Reply{
		From:    n.Pid,
		To:      request.From,
		ReplyTo: request.Code,
		Payload: map[string]interface{}{
			"secret": retrievedSecret,
		},
	}

	return nil
}
