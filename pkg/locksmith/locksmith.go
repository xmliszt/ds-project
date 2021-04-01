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
		Coordinator:         config.ConfigNode.Number,
		Ring:                make([]int, 0),
		RpcMap:              make(map[int]string),
		VirtualNodeLocation: make([]int, 0),
		VirtualNodeMap:      make(map[int]string),
		HeartBeatTable:      make(map[int]bool),
	}

	// Populate each node's RPC listening addresses
	locksmith.RpcMap[0] = "localhost:5000" // Add Locksmith receiving RPC address
	for i := 1; i <= config.ConfigNode.Number; i++ {
		nodeAddr := fmt.Sprintf("localhost:%d", config.ConfigLocksmith.Port+i)
		locksmith.RpcMap[i] = nodeAddr
		locksmith.Ring = append(locksmith.Ring, i)
	}

	go locksmith.checkHeartbeat()                            // Start periodically checking Node's heartbeat
	go locksmith.assignCoordinator(config.ConfigNode.Number) // Assign the largest node to be the coordinator

	// Start RPC server
	log.Printf("Locksmith server listening on: %v\n", address)
	rpc.Register(locksmith)
	rpc.Accept(inbound)
}

// // HandleMessageReceived is a Go Routine to handle the messages received
// func (locksmith *LockSmith) HandleMessageReceived() {
// 	for msg := range locksmith.RecvChannel {
// 		switch msg.Payload["type"] {
// 		case "REPLY_HEARTBEAT":
// 			locksmith.HeartBeatTable[msg.From] = true
// 		case "UPDATE_VIRTUAL_NODE":
// 			location := int(msg.Payload["locationData"].(uint32))
// 			virtualNode := msg.Payload["virtualNodeData"]
// 			// log.Println("Locksmith has received from ", virtualNode.(string))
// 			// Update its own values
// 			locksmith.VirtualNodeLocation = append(locksmith.VirtualNodeLocation, location)
// 			locksmith.VirtualNodeMap[location] = virtualNode.(string)

// 			// Sort the location array
// 			sort.Ints(locksmith.VirtualNodeLocation)

// 			// log.Printf("---Map of virtual node's 'Location' : 'Virtual Node Id'---\n%v\n---Array of virtual node's location---\n%v\n", locksmith.LockSmithNode.VirtualNodeMap, locksmith.LockSmithNode.VirtualNodeLocation)
// 			// Broadcast to other nodes
// 			for _, pid := range locksmith.Ring {
// 				locksmith.SendSignal(pid, &data.Data{
// 					From: locksmith.Pid,
// 					To:   pid,
// 					Payload: map[string]interface{}{
// 						"type":            "BROADCAST_VIRTUAL_NODES",
// 						"virtualNodeData": locksmith.LockSmithNode.VirtualNodeMap,
// 						"locationData":    locksmith.LockSmithNode.VirtualNodeLocation,
// 					},
// 				})
// 			}
// 		}
// 	}
// }

// CheckHeartbeat periodically check if node is alive
func (locksmith *LockSmith) checkHeartbeat() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	for {
		for _, pid := range locksmith.Ring {
			time.Sleep(time.Second * time.Duration(config.HeartbeatInterval))
			heartbeatTableCopy := make(map[int]bool)
			for k, v := range locksmith.HeartBeatTable {
				heartbeatTableCopy[k] = v
			}
			nodeClient, err := rpc.Dial("tcp", locksmith.RpcMap[pid])
			if err != nil {
				// Node is down!
				locksmith.HeartBeatTable[pid] = false
				if pid == locksmith.Coordinator {
					locksmith.assignCoordinator(locksmith.getHighestAliveNodeID())
				}
			} else {
				locksmith.HeartBeatTable[pid] = true
				nodeClient.Close()
			}
			log.Println("Heartbeat Table: ", locksmith.HeartBeatTable)
			if !reflect.DeepEqual(locksmith.HeartBeatTable, heartbeatTableCopy) {
				// Broadcast updated Heartbeat table
				locksmith.broadcastHeartbeatTable()
			}
		}
	}
}

// broadcastHeartbeatTable sends heartbeat table to all nodes
func (locksmith *LockSmith) broadcastHeartbeatTable() {
	for _, pid := range locksmith.Ring {
		request := &message.Request{
			From:    locksmith.Pid,
			To:      pid,
			Code:    message.UPDATE_HEARTBEAT_TABLE,
			Payload: locksmith.HeartBeatTable,
		}
		message.SendMessage(locksmith.RpcMap[pid], "Node.UpdateHeartbeatTable", request, nil)
	}
}

// assignCoordinator assigns the given node as the new coordinator
func (locksmith *LockSmith) assignCoordinator(pid int) {
	locksmith.Coordinator = pid
	request := &message.Request{
		From:    locksmith.Pid,
		To:      pid,
		Code:    message.ASSIGN_COORDINATOR,
		Payload: nil,
	}
	message.SendMessage(locksmith.RpcMap[pid], "Node.AssignCoordinator", request, nil)
}

// getHighestAliveNodeID gets the highest node ID whose node is currently alive
func (locksmith *LockSmith) getHighestAliveNodeID() int {
	highestNodeID := 0
	for pid, alive := range locksmith.HeartBeatTable {
		if alive {
			if pid > highestNodeID {
				highestNodeID = pid
			}
		}
	}
	return highestNodeID
}

// // Spawn new nodes when a node is down
// func (locksmith *LockSmith) SpawnNewNode(n int) {
// 	router := api.GetRouter()
// 	nodeRecvChan := make(chan *data.Data, 1)
// 	isCoordinator := false
// 	newNode := &rpc.Node{
// 		IsCoordinator: &isCoordinator,
// 		Pid:           n,
// 		RecvChannel:   nodeRecvChan,
// 		Router:        router,
// 	}

// 	locksmith.Nodes[n] = newNode
// 	locksmith.LockSmithNode.RpcMap[n] = nodeRecvChan

// 	locksmith.Nodes[n].StartDeadNode()
// 	locksmith.LockSmithNode.HeartBeatTable[n] = true

// 	// Update ring
// 	found := util.IntInSlice(locksmith.LockSmithNode.Ring, n)
// 	if !found {
// 		locksmith.LockSmithNode.Ring = append(locksmith.LockSmithNode.Ring, n)
// 	}

// 	// Update node
// 	for _, node := range locksmith.Nodes {
// 		node.Ring = locksmith.LockSmithNode.Ring
// 		node.RpcMap = locksmith.LockSmithNode.RpcMap
// 	}

// }
