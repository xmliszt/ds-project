package node

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/message"
	"github.com/xmliszt/e-safe/pkg/secret"
	"github.com/xmliszt/e-safe/util"
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
	KillSignal          chan os.Signal // For signalling shutdown of router server
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

	node := &Node{
		IsCoordinator:       false,
		Pid:                 nodeID,
		RpcMap:              make(map[int]string),
		VirtualNodeLocation: make([]int, 0),
		VirtualNodeMap:      make(map[int]string),
		HeartBeatTable:      make(map[int]bool),
		KillSignal:          make(chan os.Signal, 1),
	}

	signal.Notify(node.KillSignal, syscall.SIGTERM)

	err = node.signalNodeStart() // Send start signal to Locksmith
	if err != nil {
		log.Fatal(err)
	}
	err = node.createVirtualNodes() // Create virtual nodes
	if err != nil {
		log.Fatal(err)
	}

	// If the number of nodes present (excluding Locksmith) is greater than replication factor
	// Then start data re-distribution
	if len(node.RpcMap)-1 > config.ConfigNode.ReplicationFactor+1 {
		log.Printf("Node %d starts data re-distribution!\n", node.Pid)
		err := node.updateData() // Update data
		if err != nil {
			log.Fatal(err)
		}
	}

	// Start RPC server
	log.Printf("Node %d listening on: %v\n", node.Pid, address)
	err = rpc.Register(node)
	if err != nil {
		log.Fatal(err)
	}
	rpc.Accept(inbound)
}

// updataData grabs data from the next clockwise node
// for the replicated data, it will grab from the previous nodes
func (n *Node) updateData() error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	// do for all virtual nodes
	for i := 1; i <= config.VirtualNodesCount; i++ {
		virtualNode := strconv.Itoa(n.Pid) + "-" + strconv.Itoa(i)
		ulocation, e := util.GetHash(virtualNode)
		location := int(ulocation)
		if e != nil {
			return e
		}

		var nextPhysicalNodeID int
		var prevVirtualNodeLocation int

		ownLocationIdx := n.getVirtualLocationIndex(location)
		var nextVirtualNodeName string

		// Get the next physical node ID that is not myself
		idx := ownLocationIdx + 1
		for {
			if idx == len(n.VirtualNodeLocation) {
				idx = 0
			}
			loc := n.VirtualNodeLocation[idx]
			physicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[loc])
			if err != nil {
				return err
			}
			if physicalNodeID == n.Pid {
				idx++
			} else {
				nextPhysicalNodeID = physicalNodeID
				if ownLocationIdx-1 < 0 {
					prevVirtualNodeLocation = n.VirtualNodeLocation[len(n.VirtualNodeLocation)-int(math.Abs(float64(ownLocationIdx-1)))]
				} else {
					prevVirtualNodeLocation = n.VirtualNodeLocation[ownLocationIdx-1]
				}
				nextVirtualNodeName = n.VirtualNodeMap[loc]
				break
			}
		}

		// grab original data from the next node
		originalSecretMigrationRequest := &message.Request{
			From: n.Pid,
			To:   nextPhysicalNodeID,
			Code: message.FETCH_ORIGINAL_SECRETS,
			Payload: map[string]interface{}{
				"range":  []int{prevVirtualNodeLocation, location},
				"delete": true, // If true, the target node will delete the data after sending
			},
		}
		var originalSecretMigrationReply message.Reply
		err = message.SendMessage(n.RpcMap[nextPhysicalNodeID], "Node.GetSecrets", originalSecretMigrationRequest, &originalSecretMigrationReply)
		if err != nil {
			return err
		}
		fetchedSecrets := originalSecretMigrationReply.Payload.(map[string]*secret.Secret)
		log.Printf("Virtual Node %s fetched original secrets from Virtual Node %s: %v\n", virtualNode, nextVirtualNodeName, fetchedSecrets)

		// put secret to itself
		for k, v := range fetchedSecrets {
			err := secret.PutSecret(n.Pid, k, v)
			if err != nil {
				return err
			}
		}

		// Get replica from previous nodes using RPC
		replicationLocation, err := n.getReplicationLocations(location)
		if err != nil {
			return err
		}
		for _, slice := range replicationLocation {
			nodeID, from, to := slice[0], slice[1], slice[2]

			replicaSecretMigrationRequest := &message.Request{
				From: n.Pid,
				To:   nodeID,
				Code: message.FETCH_REPLICA_SECRETS,
				Payload: map[string]interface{}{
					"range":  []int{from, to},
					"delete": false, // if false, the target node will retain the data after sending
				},
			}
			var replicaSecretMigrationReply message.Reply
			err = message.SendMessage(n.RpcMap[nodeID], "Node.GetSecrets", replicaSecretMigrationRequest, &replicaSecretMigrationReply)
			if err != nil {
				return err
			}
			fetchedReplicas := originalSecretMigrationReply.Payload.(map[string]*secret.Secret)
			log.Printf("Virtual Node %s fetched replica secrets from Virtual Node %s: %v\n", virtualNode, n.VirtualNodeMap[to], fetchedReplicas)

			// put secret to itself
			for k, v := range fetchedReplicas {
				err := secret.PutSecret(n.Pid, k, v)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// signalNodeStart sends a signal to Locksmith server that the node has started
// it is for Locksmith server to respond with the current RPC map-
func (n *Node) signalNodeStart() error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	request := &message.Request{
		From:    n.Pid,
		To:      0,
		Code:    message.SIGNAL_START,
		Payload: nil,
	}
	var reply message.Reply
	err = message.SendMessage(fmt.Sprintf("localhost:%d", config.ConfigLocksmith.Port), "LockSmith.SignalStart", request, &reply)
	if err != nil {
		return err
	}
	n.RpcMap = reply.Payload.(map[int]string)
	log.Printf("Node %d RPC map updated: %+v\n", n.Pid, n.RpcMap)
	// Relay updated RPC map to others
	for pid, address := range n.RpcMap {
		if pid == n.Pid || pid == 0 {
			continue
		}
		request = &message.Request{
			From:    n.Pid,
			To:      pid,
			Code:    message.UPDATE_RPC_MAP,
			Payload: n.RpcMap,
		}
		err = message.SendMessage(address, "Node.UpdateRpcMap", request, &reply)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

// Create virtual nodes
func (n *Node) createVirtualNodes() error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	virtualNodesData := make(map[int]string)
	virtualLocations := make([]int, 0)

	for i := 1; i <= config.VirtualNodesCount; i++ {
		virtualNode := strconv.Itoa(n.Pid) + "-" + strconv.Itoa(i)
		ulocation, e := util.GetHash(virtualNode)
		location := int(ulocation)
		if e != nil {
			return e
		}

		virtualNodesData[location] = virtualNode
		virtualLocations = append(virtualLocations, location)
	}
	request := &message.Request{
		From: n.Pid,
		To:   0,
		Code: message.CREATE_VIRTUAL_NODE,
		Payload: map[string]interface{}{
			"virtualNodeMap":      virtualNodesData,
			"virtualNodeLocation": virtualLocations,
		},
	}
	var reply message.Reply
	err = message.SendMessage(n.RpcMap[0], "LockSmith.CreateVirtualNodes", request, &reply)
	if err != nil {
		return err
	}
	payload := reply.Payload.(map[string]interface{})
	n.VirtualNodeMap = payload["virtualNodeMap"].(map[int]string)
	n.VirtualNodeLocation = payload["virtualNodeLocation"].([]int)
	log.Printf("Node %d has created virtual nodes: %+v | %+v\n", n.Pid, n.VirtualNodeLocation, n.VirtualNodeMap)

	// Relay updated virtual nodes to others
	for pid, address := range n.RpcMap {
		if pid == n.Pid || pid == 0 {
			continue
		}
		request = &message.Request{
			From: n.Pid,
			To:   pid,
			Code: message.UPDATE_VIRTUAL_NODES,
			Payload: map[string]interface{}{
				"virtualNodeMap":      n.VirtualNodeMap,
				"virtualNodeLocation": n.VirtualNodeLocation,
			},
		}
		err = message.SendMessage(address, "Node.UpdateVirtualNodes", request, &reply)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

// Starts the router -> Graceful shutdown
func (n *Node) startRouter() {
	n.Router = n.getRouter()
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err := n.Router.Start(fmt.Sprintf(":%d", config.ConfigServer.Port))
		if err != nil {
			log.Printf("Node %d REST server closed!\n", n.Pid)
		}
	}()
	<-n.KillSignal // Blocking, until kill signal received
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = n.Router.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
