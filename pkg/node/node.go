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

	if len(node.VirtualNodeLocation) > config.VirtualNodesCount {
		err = node.updateData(nodeID) // Update data
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
func (n *Node) updateData(nodeID int) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	// do for all virtual ndoes
	for i := 1; i <= config.VirtualNodesCount; i++ {
		virtualNode := strconv.Itoa(nodeID) + "-" + strconv.Itoa(i)
		ulocation, e := util.GetHash(virtualNode)
		location := int(ulocation)
		if e != nil {
			return e
		}

		for idx, loc := range n.VirtualNodeLocation {
			if loc == location {
				fmt.Println("loc", loc)
				fmt.Println("idx", idx)
				var nextVirtualNodeLocation int
				if idx+1 == len(n.VirtualNodeLocation) {
					nextVirtualNodeLocation = n.VirtualNodeLocation[0]
				} else {
					nextVirtualNodeLocation = n.VirtualNodeLocation[idx+1]
				}
				var prevVirtualNodeLocation int
				if idx-1 < 0 { // loop from the tail of the slice if index of location is less than replication factor
					prevVirtualNodeLocation = n.VirtualNodeLocation[len(n.VirtualNodeLocation)-1]
				} else {
					prevVirtualNodeLocation = n.VirtualNodeLocation[idx-1]
				}
				nextVirtualNode := n.VirtualNodeMap[nextVirtualNodeLocation]
				nextPhysicalNode, err := util.GetPhysicalNode(nextVirtualNode)
				if err != nil {
					return err
				}
				// grab all data from the node
				// TODO: change this to contact using RPC
				dataFromNextNode, err := secret.GetSecrets(nextPhysicalNode, prevVirtualNodeLocation, location)
				fmt.Println("range is", prevVirtualNodeLocation, location, "physicalnode", nextPhysicalNode, "datafromnextnode", dataFromNextNode)
				if err != nil {
					return err
				}

				// put secret to itself
				for k, v := range dataFromNextNode {
					err := secret.PutSecret(n.Pid, k, v)
					if err != nil {
						return err
					}
				}

				// TODO: Get replica from previous nodes using RPC
				replicationLocation, err := n.getReplicationLocation(location)
				if err != nil {
					return err
				}
				fmt.Println("replication location", replicationLocation)
				for _, slice := range replicationLocation {
					fmt.Println("slice0", slice[0], "slice1", slice[1], "slice2", slice[2])
					dataFromPrevNode, err := secret.GetSecrets(slice[0], slice[1], slice[2])
					fmt.Println("data from prev node", dataFromPrevNode)
					if err != nil {
						return err
					}
					// put secret to itself
					for k, v := range dataFromPrevNode {
						err := secret.PutSecret(nodeID, k, v)
						if err != nil {
							return err
						}
					}
				}
				// TODO: Delete from the next physical node using RPC
				break
			}
		}
	}
	return nil
}

func (n *Node) getReplicationLocation(location int) ([][]int, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	res := make([][]int, 0)
	nodesVisited := make([]int, 0)
	loc := util.GetIndex(n.VirtualNodeLocation, location)
	var virtualNode string

	for len(nodesVisited) <= config.ReplicationFactor {
		if loc < 0 { // loop from the tail of the slice if index of location is less than replication factor
			moduloLoc := (loc * -1) % len(n.VirtualNodeLocation)
			virtualNode = n.VirtualNodeMap[n.VirtualNodeLocation[len(n.VirtualNodeLocation)-moduloLoc]]
		} else {
			virtualNode = n.VirtualNodeMap[n.VirtualNodeLocation[loc]]
		}
		fmt.Println("virtualnode", virtualNode)
		physicalNode, err := util.GetPhysicalNode(virtualNode)
		if err != nil {
			return nil, err
		}
		fmt.Println("physical node", physicalNode)
		if !util.IntInSlice(nodesVisited, physicalNode) {
			nodesVisited = append(nodesVisited, physicalNode)
			uvirtualNodeLoc, err := util.GetHash(virtualNode)
			if err != nil {
				return nil, err
			}
			virtualNodeLoc := int(uvirtualNodeLoc)
			virtualNodeIdx := util.GetIndex(n.VirtualNodeLocation, virtualNodeLoc)
			prevVirtualNodeIdx := virtualNodeIdx - 1
			var prevVirtualNodeLoc int
			if prevVirtualNodeIdx == -1 { // loop from the tail of the slice if index of location is less than replication factor
				prevVirtualNodeLoc = n.VirtualNodeLocation[len(n.VirtualNodeLocation)-1]
			} else {
				prevVirtualNodeLoc = n.VirtualNodeLocation[prevVirtualNodeIdx]
			}
			slice := []int{physicalNode, prevVirtualNodeLoc, virtualNodeLoc}
			res = append(res, slice)
			fmt.Println("res", res)
			loc -= 1
		}
	}

	return res, nil
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
