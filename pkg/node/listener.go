package node

import (
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

	log.Printf("Node %d done fetching secrets from %d to %d: %v\n", n.Pid, from, to, secrets)

	if toDelete {
		for key := range secrets {
			err := secret.RemoveSecret(n.Pid, key)
			if err != nil {
				return err
			}
		}
		log.Printf("Node %d done deleting secrets from %d to %d: %v\n", n.Pid, from, to, secrets)
	}

	return nil
}
