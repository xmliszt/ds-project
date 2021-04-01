package locksmith

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"reflect"
	"time"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/message"
)

type LockSmith struct {
	Pid                 int            // Node ID
	Ring                []int          // Ring structure of nodes
	Coordinator         int            // Indicate the coordinator node number
	RpcMap              map[int]string // Map node ID to their receiving address
	HeartBeatTable      map[int]bool
	VirtualNodeLocation []int
	VirtualNodeMap      map[int]string
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
	}

	// Populate each node's RPC listening addresses
	locksmith.RpcMap[0] = "localhost:5000" // Add Locksmith receiving RPC address

	go locksmith.checkHeartbeat()           // Start periodically checking Node's heartbeat
	go locksmith.monitorCoordinatorStatus() // Start periodically monitor and update coordinator

	// Start RPC server
	log.Printf("Locksmith server listening on: %v\n", address)
	rpc.Register(locksmith)
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
				// Node is down!
				locksmith.HeartBeatTable[pid] = false
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

// broadcastHeartbeatTable sends heartbeat table to all nodes
func (locksmith *LockSmith) broadcastHeartbeatTable(excludeNodeID interface{}) {
	if excludeNodeID != nil {
		excludeNodeID = excludeNodeID.(int)
	}
	for pid, address := range locksmith.RpcMap {
		if pid == excludeNodeID || pid == 0 {
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
			log.Fatal(err)
		}
	}
}

// monitorCoordinatorStatus monitors heartbeat table and always assign the highest alive node as coordinator
func (locksmith *LockSmith) monitorCoordinatorStatus() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	for {
		time.Sleep(time.Second * time.Duration(config.CoordinatorMonitorInterval))
		newCoordinatorID := locksmith.Coordinator
		for pid, alive := range locksmith.HeartBeatTable {
			if alive && pid > locksmith.Coordinator {
				newCoordinatorID = pid
			}
		}
		if newCoordinatorID != locksmith.Coordinator {
			locksmith.assignCoordinator(newCoordinatorID)
		}
	}
}

// assignCoordinator assigns the given node as the new coordinator
func (locksmith *LockSmith) assignCoordinator(pid int) {
	var reply message.Reply
	var request *message.Request

	// Remove the old coordinator
	if locksmith.Coordinator > 0 {
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
