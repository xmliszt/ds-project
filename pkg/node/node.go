package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"strings"
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

	// Start RPC server
	log.Printf("Node %d listening on: %v\n", node.Pid, address)
	err = rpc.Register(node)
	if err != nil {
		log.Fatal(err)
	}
	rpc.Accept(inbound)
}

// signalNodeStart sends a signal to Locksmith server that the node has started
// it is for Locksmith server to respond with the current RPC map
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

// Strict Consistency with R = 2. Send ACK directly to coordinator
func (n *Node) strictDown(rf int, key string, secret secret.Secret, relayNodes []int) error {
	config, err := config.GetConfig()
	if err != nil {
		log.Printf("Node %d is unable to strict down to next node: %s\n", n.Pid, err)
		return err
	}
	nextNodeLoc := relayNodes[config.ConfigNode.ReplicationFactor-rf]
	nextPhysicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[nextNodeLoc])
	if err != nil {
		log.Printf("Node %d is unable to strict down to next node: %s\n", n.Pid, err)
		return err
	}
	nextNodeAddr := n.RpcMap[nextPhysicalNodeID]
	request := &message.Request{
		From: n.Pid,
		To:   nextPhysicalNodeID,
		Code: message.STRICT_OWNER_DOWN,
		Payload: map[string]interface{}{
			"rf":     rf,
			"key":    key,
			"secret": secret,
			"nodes":  relayNodes,
		},
	}

	var reply message.Reply
	err = message.SendMessage(nextNodeAddr, "Node.PerformStrictDown", request, &reply)
	return nil
}

func (n *Node) performEventualRep(rf int, key string, secret secret.Secret, relayNodes []int) error {

	return nil
}

// Strict node sends to Even node
func (n *Node) sendEventualRepMsg(rf int, key string, secret secret.Secret, relayNodes []int) error {
	config, err := config.GetConfig()
	if err != nil {
		log.Printf("Node %d is unable to relay secret deletion to next node: %s\n", n.Pid, err)
		return err
	}
	nextNodeLoc := relayNodes[config.ConfigNode.ReplicationFactor-rf]
	nextPhysicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[nextNodeLoc])
	if err != nil {
		log.Printf("Node %d is unable to relay secret deletion to next node: %s\n", n.Pid, err)
		return err
	}
	nextNodeAddr := n.RpcMap[nextPhysicalNodeID]
	request := &message.Request{
		From: n.Pid,
		To:   nextPhysicalNodeID,
		Code: message.EVENTUAL_STORE,
		Payload: map[string]interface{}{
			"rf":     rf,
			"key":    key,
			"secret": secret,
			"nodes":  relayNodes,
		},
	}

	var reply message.Reply
	err = message.SendMessage(nextNodeAddr, "Node.PerformEventualReplication", request, &reply)
	return nil
}

// Takes in Hash value. Will generate list of nodes to store the data
// need to check the logic later
// implemented in helper.go already
func (n *Node) generateReplicationList(hashValue uint32) []string {
	var VirtualNodeList []string
	var nodeIdList []int
	var positionForOwner int
	nextPid, nextVirtualNode := n.findNextPid(hashValue)
	VirtualNodeList = append(VirtualNodeList, nextVirtualNode)
	nodeIdList = append(nodeIdList, nextPid)
	hashValueForOwnerNode, _ := util.GetHash(nextVirtualNode)
	for idx, location := range n.VirtualNodeLocation {
		if location == int(hashValueForOwnerNode) {
			positionForOwner = idx
		}
	}
	if len(VirtualNodeList) > 3 {
		for indx, position := range n.VirtualNodeLocation {
			if indx > positionForOwner {
				nextId := n.findPidByVname(n.VirtualNodeMap[position])
				// need to check whether the id in the nodeIdList or not
				if contain(nodeIdList, nextId) {
					continue
				} else {
					nodeIdList = append(nodeIdList, nextId)
					VirtualNodeList = append(VirtualNodeList, n.VirtualNodeMap[position])
				}
			}
		}
	}
	return VirtualNodeList

}

func contain(nodeIdList []int, id int) bool {
	var result bool
	for id1 := range nodeIdList {
		if id1 == id {
			result = true
		} else {
			result = false
		}
	}
	return result
}

// No input. Will generate list of nodes to store the data.
func (n *Node) generateRecoveryList() {}

// TODO: move to util
// This function will takes in the hashed value then find the next node's pid and next node's virtual node name
func (n *Node) findNextPid(hashedValue uint32) (int, string) {
	var nextVirtualNode string
	var nextPid int
	var err error
	for idx, location := range n.VirtualNodeLocation {
		if int(hashedValue) < location {
			// current_virtual_node := n.RingMap[x]
			nextVirtualNode = n.VirtualNodeMap[n.Ring[(idx+1)]]
			string_list := strings.Split(nextVirtualNode, "-")
			nextPid, err = strconv.Atoi(string_list[0])
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}
	return nextPid, nextVirtualNode
}

// TODO: Already in Util, consider removing
func (n *Node) findPidByVname(vName string) int {
	string_list := strings.Split(vName, "-")
	Pid, _ := strconv.Atoi(string_list[0])
	return Pid
}

func (n *Node) checkHeartbeat(pid int) bool {
	return n.HeartBeatTable[pid]
}

// Start Strict Replication
func (n *Node) sendStrictRepMsg(rf int, key string, value secret.Secret, relayNodes []int) error {
	config, err := config.GetConfig()
	if err != nil {
		log.Printf("Node %d is unable to relay strict consistency to next node: %s\n", n.Pid, err)
		return err
	}
	nextNodeLoc := relayNodes[config.ConfigNode.ReplicationFactor-rf]
	nextPhysicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[nextNodeLoc])
	if err != nil {
		log.Printf("Node %d is unable to relay strict consistency to next node: %s\n", n.Pid, err)
		return err
	}
	nextNodeAddr := n.RpcMap[nextPhysicalNodeID]
	request := &message.Request{
		From: n.Pid,
		To:   nextPhysicalNodeID,
		Code: message.STRICT_STORE,
		Payload: map[string]interface{}{
			"rf":     rf,
			"key":    key,
			"secret": value,
			"nodes":  relayNodes,
		},
	}

	var reply message.Reply
	err = message.SendMessage(nextNodeAddr, "Node.StrictReplication", request, &reply)
	replyPayload := reply.Payload.(map[string]interface{})
	if replyPayload["success"].(bool) == true {
		// ISSUE: How do I respond to the caller with a positive ACK.
		return nil
	}
	if err != nil {
		log.Printf("Node %d strict consistency error: %s\n", n.Pid, err)
		return err
	}
	return nil
}
