package locksmith

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"reflect"
	"strconv"
	"time"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/message"
	"github.com/xmliszt/e-safe/util"
)

type LockSmith struct {
	Pid                 int            // Node ID
	Ring                []int          // Ring structure of nodes
	Coordinator         int            // Indicate the coordinator node number
	RpcMap              map[int]string // Map node ID to their receiving address
	HeartBeatTable      map[int]bool
	VirtualNodeLocation []int
	VirtualNodeMap      map[int]string
	RequestQueue        []int // Request queue for granting lock for accessing user.json, an array of nodeID
}

// Start is the main function that starts the entire program
func Start() {

	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", config.ConfigLocksmith.Port))
	if err != nil {
		log.Fatal(err)
	}
	inbound, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	locksmith := &LockSmith{
		Pid:                 0,
		Coordinator:         0,
		RpcMap:              make(map[int]string),
		VirtualNodeLocation: make([]int, 0),
		VirtualNodeMap:      make(map[int]string),
		HeartBeatTable:      make(map[int]bool),
		RequestQueue:        make([]int, 0),
	}

	// Populate each node's RPC listening addresses
	locksmith.RpcMap[0] = "localhost:5000" // Add Locksmith receiving RPC address

	go locksmith.checkHeartbeat()           // Start periodically checking Node's heartbeat
	go locksmith.monitorCoordinatorStatus() // Start periodically monitor and update coordinator

	// Start RPC server
	log.Printf("Locksmith server listening on: %v\n", address)
	err = rpc.Register(locksmith)
	if err != nil {
		log.Fatal(err)
	}
	rpc.Accept(inbound)
}

// checkHeartbeat periodically check if node is alive
func (locksmith *LockSmith) checkHeartbeat() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	for {
		time.Sleep(time.Second * time.Duration(config.HeartbeatInterval))
		for pid := range locksmith.RpcMap {
			if pid == 0 {
				continue
			}
			heartbeatTableCopy := make(map[int]bool)
			for k, v := range locksmith.HeartBeatTable {
				heartbeatTableCopy[k] = v
			}
			nodeClient, err := rpc.Dial("tcp", locksmith.RpcMap[pid])
			if err != nil {
				// Node is down! If heartbeat table shows previously alive
				if locksmith.HeartBeatTable[pid] {
					locksmith.HeartBeatTable[pid] = false
					// if happened to be the coordinator
					if locksmith.Coordinator == pid {
						locksmith.Coordinator = 0 // reset coordinator
					}
					// Remove virtual nodes
					err := locksmith.removeVirtualNodes(pid)
					if err != nil {
						log.Fatal(err)
					}
				}
			} else {
				locksmith.HeartBeatTable[pid] = true
				nodeClient.Close()
			}
			if !reflect.DeepEqual(locksmith.HeartBeatTable, heartbeatTableCopy) {
				// Broadcast updated Heartbeat table
				locksmith.broadcastHeartbeatTable(nil)
			}
		}
		log.Println("Heartbeat Table: ", locksmith.HeartBeatTable)
	}
}

// removeVirtualNodes removes the dead node's virtual node locations and map
func (locksmith *LockSmith) removeVirtualNodes(nodeID int) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	virtualLocations := make([]int, 0)

	for i := 1; i <= config.VirtualNodesCount; i++ {
		virtualNode := strconv.Itoa(nodeID) + "-" + strconv.Itoa(i)
		ulocation, e := util.GetHash(virtualNode)
		location := int(ulocation)
		if e != nil {
			return e
		}

		virtualLocations = append(virtualLocations, location)
	}

	// Remove from map
	for _, location := range virtualLocations {
		delete(locksmith.VirtualNodeMap, location)
	}

	// Remove from location
	newLocations := make([]int, 0)
	for _, location := range locksmith.VirtualNodeLocation {
		if !util.IntInSlice(virtualLocations, location) {
			newLocations = append(newLocations, location)
		}
	}
	locksmith.VirtualNodeLocation = newLocations

	err = locksmith.broadcastVirtualNodes()
	if err != nil {
		return err
	}
	return nil
}

// broadcastHeartbeatTable sends heartbeat table to all nodes
func (locksmith *LockSmith) broadcastHeartbeatTable(excludeNodeID interface{}) {
	if excludeNodeID != nil {
		excludeNodeID = excludeNodeID.(int)
	}
	for pid, address := range locksmith.RpcMap {
		if pid == excludeNodeID || pid == 0 || !locksmith.HeartBeatTable[pid] {
			continue
		}
		request := &message.Request{
			From:    locksmith.Pid,
			To:      pid,
			Code:    message.UPDATE_HEARTBEAT_TABLE,
			Payload: locksmith.HeartBeatTable,
		}
		var reply message.Reply
		err := message.SendMessage(address, "Node.UpdateHeartbeatTable", request, &reply)
		if err != nil {
			log.Printf("Locksmith failed to send Heartbeat Table to Node %d: %s\n", pid, err)
		}
	}
}

// broadcastVirtualNodes sends the modified virtual nodes to every alive nodes
// it is only done when a node is dead and virtual nodes are modified
func (locksmith *LockSmith) broadcastVirtualNodes() error {
	// Relay updated virtual nodes to others
	for pid, address := range locksmith.RpcMap {
		if pid == locksmith.Pid || !locksmith.HeartBeatTable[pid] {
			continue
		}
		request := &message.Request{
			From: locksmith.Pid,
			To:   pid,
			Code: message.UPDATE_VIRTUAL_NODES,
			Payload: map[string]interface{}{
				"virtualNodeMap":      locksmith.VirtualNodeMap,
				"virtualNodeLocation": locksmith.VirtualNodeLocation,
			},
		}
		var reply message.Reply
		err := message.SendMessage(address, "Node.UpdateVirtualNodes", request, &reply)
		if err != nil {
			return err
		}
	}
	return nil
}

// monitorCoordinatorStatus monitors heartbeat table and always assign the highest alive node as coordinator
func (locksmith *LockSmith) monitorCoordinatorStatus() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	for {
		time.Sleep(time.Second * time.Duration(config.CoordinatorMonitorInterval))
		newCoordinatorID := 0
		for pid, alive := range locksmith.HeartBeatTable {
			if alive && pid > newCoordinatorID {
				newCoordinatorID = pid
			}
		}
		if newCoordinatorID != locksmith.Coordinator && newCoordinatorID > 0 {
			locksmith.assignCoordinator(newCoordinatorID)
		}
	}
}

// assignCoordinator assigns the given node as the new coordinator
func (locksmith *LockSmith) assignCoordinator(pid int) {
	var reply message.Reply
	var request *message.Request

	// Remove the old coordinator
	if locksmith.Coordinator > 0 && locksmith.HeartBeatTable[locksmith.Coordinator] {
		request = &message.Request{
			From:    locksmith.Pid,
			To:      pid,
			Code:    message.REMOVE_COORDINATOR,
			Payload: nil,
		}
		err := message.SendMessage(locksmith.RpcMap[locksmith.Coordinator], "Node.RemoveCoordinator", request, &reply)
		if err != nil {
			panic(err)
		}
	}

	// Assign new coordinator
	locksmith.Coordinator = pid
	request = &message.Request{
		From:    locksmith.Pid,
		To:      pid,
		Code:    message.ASSIGN_COORDINATOR,
		Payload: nil,
	}
	err := message.SendMessage(locksmith.RpcMap[pid], "Node.AssignCoordinator", request, &reply)
	if err != nil {
		panic(err)
	}
}
