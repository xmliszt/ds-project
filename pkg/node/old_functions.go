// package node

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/xmliszt/e-safe/config"
// 	"github.com/xmliszt/e-safe/pkg/data"
// 	"github.com/xmliszt/e-safe/pkg/file"
// 	"github.com/xmliszt/e-safe/pkg/secret"
// )

// // Node contains all the variables that are necessary to manage a node
// // type Node struct {
// // 	IsCoordinator       *bool                   `validate:"required"`
// // 	Pid                 int                     `validate:"gte=0"`    // Node ID
// // 	Ring                []int                   `validate:"required"` // Ring structure of nodes
// // 	RecvChannel         chan *data.Data         `validate:"required"` // Receiving channel
// // 	SendChannel         chan *data.Data         `validate:"required"` // Sending channel
// // 	RpcMap              map[int]chan *data.Data `validate:"required"` // Map node ID to their receiving channels
// // 	HeartBeatTable      map[int]bool            // Heartbeat table
// // 	VirtualNodeLocation []int
// // 	VirtualNodeMap      map[int]string
// // 	Signal              chan error
// // }

// // HandleMessageReceived is a Go routine that handles the messages received
// func (n *Node) HandleMessageReceived() {

// 	// Test a dead node
// 	if n.Pid == 5 {
// 		go func() {
// 			time.Sleep(time.Second * 30)
// 			defer close(n.RecvChannel)
// 		}()
// 	}

// 	for msg := range n.RecvChannel {
// 		switch msg.Payload["type"] {
// 		case "CHECK_HEARTBEAT":
// 			n.SendSignal(0, &data.Data{
// 				From: n.Pid,
// 				To:   0,
// 				Payload: map[string]interface{}{
// 					"type": "REPLY_HEARTBEAT",
// 					"data": nil,
// 				},
// 			})
// 		case "UPDATE_HEARTBEAT":
// 			heartbeatTable := msg.Payload["data"]
// 			n.HeartBeatTable = heartbeatTable.(map[int]bool)
// 		case "YOU_ARE_COORDINATOR":
// 			isCoordinator := true
// 			n.IsCoordinator = &isCoordinator
// 		// Secret sent to node that it is hased to
// 		case "STORE_AND_REPLICATE":
// 			// Store the data to the main node. Check if the neighbouring node is alive. If not then send to the next node.
// 			n.StartReplication(msg)
// 			// receivedPayload := msg.Payload["data"]                       // This should be the Secret
// 			// hashedValue := fmt.Sprintf("%v", msg.Payload["hashedValue"]) // This should be the hashed value omt for that secret
// 			// mapPayload := map[string]interface{}{
// 			// 	hashedValue: receivedPayload,
// 			// }
// 			// file.WriteDataFile(n.Pid, mapPayload)
// 			// // How are we going from n.Pid to next_id
// 			// // hashedValue -> current virtual node number (1-1)
// 			// // using Ring, current virtual node number -> next virtual node number
// 			// var next_pid int
// 			// for index, x := range n.Ring {

// 			// 	hashedValueINT, err := strconv.Atoi(hashedValue)
// 			// 	if err != nil {
// 			// 		fmt.Println(err)
// 			// 	}
// 			// 	if hashedValueINT < x {
// 			// 		// current_virtual_node := n.RingMap[x]
// 			// 		next_virtual_node := n.VirtualNodeMap[n.Ring[(index+1)]]
// 			// 		string_list := strings.Split(next_virtual_node, "-")
// 			// 		next_pid, err = strconv.Atoi(string_list[0])
// 			// 		if err != nil {
// 			// 			fmt.Println(err)
// 			// 		}
// 			// 		break

// 			// 	}

// 			// }
// 			// n.SendSignal(next_pid, &data.Data{
// 			// 	From: n.Pid,
// 			// 	To:   next_pid,
// 			// 	Payload: map[string]interface{}{
// 			// 		"type": "STRICT_CONSISTENCY",
// 			// 		"data": mapPayload,
// 			// 	},
// 			// })
// 		// Sent by the owner node to the neighbouring node
// 		case "STRICT_CONSISTENCY":
// 			// Replicate secret into the node and then call for eventual consistency in the next R-1 nodes
// 			n.StrictReplication(msg)
// 		// Sent by the neighbouring node to the next R-1 nodes
// 		case "EVENTUAL_STORE":
// 			n.EventualReplication(msg)
// 		// Sent by the neighbouring node to the owner node
// 		case "ACK_OWNER_NODE":
// 			// TODO: need to send signal to coordinator
// 			n.SendSignal(0, &data.Data{
// 				From: n.Pid,
// 				To:   0,
// 				Payload: map[string]interface{}{
// 					"type": "ASK_COORDINATOR",
// 					"data": nil,
// 				},
// 			})
// 			// Need to ask locksmith who is the coordinator
// 			// Sent by the owner node to the coordinator
// 		case "REPLY_COORDINATOR":
// 			n.ackCoordinator(msg)
// 		case "ACK_COORDINATOR":
// 			n.Signal <- nil
// 			// Reply to coordintor that write was successful

// 		// Replication Factor Value -- (decrement)
// 		// Secret to be stored
// 		// Check if Replication Factor Value == 0. If yes, stop process - data successfully replicated.
// 		// If not then send another EVEN_CONSISTENCY message to the next node

// 		//SendSignal for ack to owner node
// 		case "BROADCAST_VIRTUAL_NODES":
// 			location := msg.Payload["locationData"]
// 			virtualNode := msg.Payload["virtualNodeData"]
// 			n.VirtualNodeLocation = location.([]int)
// 			n.VirtualNodeMap = virtualNode.(map[int]string)
// 		}
// 	}
// }

// // Create virtual nodes
// func (n *Node) CreateVirtualNodes(Pid int) error {
// 	conf, err := config.GetConfig()
// 	if err != nil {
// 		return err
// 	}

// 	for i := 1; i <= conf.VirtualNodesCount; i++ {
// 		virtualNode := strconv.Itoa(Pid) + "-" + strconv.Itoa(i)
// 		location, e := config.GetHash(virtualNode)
// 		if e != nil {
// 			return e
// 		}
// 		fmt.Println("Virtual node ", virtualNode, "has started")
// 		n.SendSignal(0, &data.Data{
// 			From: Pid,
// 			To:   0,
// 			Payload: map[string]interface{}{
// 				"type":            "UPDATE_VIRTUAL_NODE",
// 				"virtualNodeData": virtualNode,
// 				"locationData":    location,
// 			},
// 		})
// 	}
// 	return nil
// }

// // Start starts up a node, running receiving channel
// func (n *Node) Start() error {
// 	fmt.Printf("Node [%d] has started!\n", n.Pid)

// 	// Create virtual node
// 	err := n.CreateVirtualNodes(n.Pid)
// 	if err != nil {
// 		return err
// 	}
// 	go n.HandleMessageReceived()
// 	return nil
// }

// // Start starts up a node, running receiving channel
// func (n *Node) StartDeadNode() {
// 	fmt.Printf("Node [%d] has started!\n", n.Pid)
// 	go n.HandleMessageReceived()
// }

// // TearDown terminates node, closes all channels
// func (n *Node) TearDown() {
// 	close(n.RecvChannel)
// 	close(n.SendChannel)
// 	fmt.Printf("Node [%d] has terminated!\n", n.Pid)
// }

// // This is called by the coordinator node
// func (n *Node) PutSecret(alias string, value string, role int) error {
// 	// Figure out which node to go to
// 	aliasHash, e := config.GetHash(alias)
// 	if e != nil {
// 		fmt.Errorf("Hashing error")
// 	}
// 	stringifiedHash := fmt.Sprintf("%c", aliasHash)
// 	nodePid, _ := n.mapHashToPid(aliasHash)
// 	// Check if node is alive
// 	nodeAlive := n.checkHeartbeat(nodePid)

// 	// Prep secret struct
// 	incomingSecret := secret.Secret{
// 		Value: value,
// 		Role:  role,
// 	}

// 	// Send replication signal
// 	if nodeAlive {
// 		n.SendSignal(nodePid, &data.Data{
// 			From: n.Pid,
// 			To:   nodePid,
// 			Payload: map[string]interface{}{
// 				"type":              "STORE_AND_REPLICATE",
// 				"replicationFactor": 3, // TODO: need to write this in the config.go
// 				"hashedValue":       stringifiedHash,
// 				"data":              incomingSecret,
// 			}})

// 		// Wait for acknowledge from the owner node
// 		err := <-n.Signal
// 		return err
// 	} else {
// 		// return fmt.Errorf("Owner Node not alive")
// 		fmt.Println("Owner Node not alive\nSending to subsequent node")
// 		nextPid, _ := n.findNextPid(aliasHash)
// 		n.SendSignal(nextPid, &data.Data{
// 			From: n.Pid,
// 			To:   nextPid,
// 			Payload: map[string]interface{}{
// 				"type":              "STRICT_REPLICATION", // Add to handleMessageReceived
// 				"replicationFactor": 2,                    // TODO: need to write this in the config.go
// 				"hashedValue":       aliasHash,
// 				"data":              incomingSecret,
// 			}})
// 		err := <-n.Signal
// 		return err
// 	}

// }

// // This is called by the owner node
// // It stores the secret & begins the replication process
// func (n *Node) StartReplication(incomingData *data.Data) error {
// 	// may got problem
// 	key := incomingData.Payload["hashedValue"].(string)
// 	secret := incomingData.Payload["data"]
// 	secretToAdd := map[string]interface{}{
// 		incomingData.Payload["hashedValue"].(string): incomingData.Payload["data"],
// 	}
// 	file.WriteDataFile(n.Pid, secretToAdd)

// 	// next_node := n.Pid + 1
// 	ukey, err := strconv.Atoi(key)
// 	if err != nil {
// 		return err
// 	}
// 	nextNode, nextVirtualNode := n.findNextPid(uint32(ukey))
// 	// AndVirtualNodeId, virtual_node_name := n.mapHashToPid(uint32(next_node))
// 	// Check if next node is alive
// 	nodeAlive := n.checkHeartbeat(nextNode)
// 	// Send replication signal
// 	if nodeAlive {

// 		n.SendSignal(nextNode, &data.Data{
// 			From: n.Pid,
// 			To:   nextNode,
// 			Payload: map[string]interface{}{
// 				"type":              "STRICT_CONSISTENCY",
// 				"replicationFactor": 3,   // TODO: Take from config.go and minus 1
// 				"hashedValue":       key, // This is the string of the hash
// 				"data":              secret,
// 			}})

// 		// Wait for acknowledge from the neighbouring node
// 		return nil
// 	} else {
// 		// return fmt.Errorf("Owner Node not alive")
// 		fmt.Println("Subsequent Node not alive\nSending to the next node in sequence")
// 		// TODO: Find the node next to the next node
// 		nextNextNode, _ := n.findNextWithPid(nextVirtualNode)

// 		n.SendSignal(nextNextNode, &data.Data{
// 			From: n.Pid,
// 			To:   nextNextNode,
// 			Payload: map[string]interface{}{
// 				"type":              "STRICT_CONSISTENCY", // Add to handleMessageReceived
// 				"replicationFactor": 2,                    // TODO: need to write this in the config.go
// 				"hashedValue":       key,
// 				"data":              secret,
// 			},
// 		})
// 		return nil
// 	}
// 	// return nil

// }

// // this is called by the subsequent node
// func (n *Node) strictReplication(incomingData *data.Data) error {
// 	stringifiedHash := incomingData.Payload["hashedValue"].(string)
// 	secret := incomingData.Payload["data"]
// 	secretToAdd := map[string]interface{}{
// 		incomingData.Payload["hashedValue"].(string): incomingData.Payload["data"],
// 	}
// 	err := file.WriteDataFile(n.Pid, secretToAdd)
// 	if err != nil {
// 		return err
// 	} else {
// 		uIntHashValue, uIntConvertError := strconv.ParseUint(stringifiedHash, 32, 32)
// 		if uIntConvertError != nil {
// 			return uIntConvertError
// 		}
// 		ownerPid, _ := n.mapHashToPid(uint32(uIntHashValue))
// 		_, nextVirtualNode := n.findNextPid(uint32(uIntHashValue))

// 		// ack to owner node
// 		n.SendSignal(ownerPid, &data.Data{
// 			From: n.Pid,
// 			To:   ownerPid,
// 			Payload: map[string]interface{}{
// 				"type": "ACK_OWNER_NODE",
// 				"data": nil,
// 			},
// 		})
// 		//  whether the replication factor is 2 or 3
// 		replication_fac := incomingData.Payload["replicationFactor"].(int)
// 		if replication_fac == 2 {
// 			nextPid, _ := n.findNextWithPid(nextVirtualNode)
// 			//  check whetehr next node alive
// 			if n.checkHeartbeat(nextPid) {
// 				n.SendSignal(nextPid, &data.Data{
// 					From: n.Pid,
// 					To:   nextPid,
// 					Payload: map[string]interface{}{
// 						"type":              "EVENTUAL_STORE",
// 						"replicationFactor": 1,
// 						"hashedValue":       stringifiedHash,
// 						"data":              secret,
// 					},
// 				})
// 			}
// 		} else if replication_fac == 3 {
// 			nextPid, nextNextVirtualNode := n.findNextWithPid(nextVirtualNode)
// 			if n.checkHeartbeat(nextPid) {
// 				n.SendSignal(nextPid, &data.Data{
// 					From: n.Pid,
// 					To:   nextPid,
// 					Payload: map[string]interface{}{
// 						"type":              "EVENTUAL_STORE",
// 						"replicationFactor": 2,
// 						"hashedValue":       stringifiedHash,
// 						"data":              secret,
// 					},
// 				})
// 			} else {
// 				nextNextPid, _ := n.findNextWithPid(nextNextVirtualNode)
// 				n.SendSignal(nextNextPid, &data.Data{
// 					From: n.Pid,
// 					To:   nextNextPid,
// 					Payload: map[string]interface{}{
// 						"type":              "EVENTUAL_STORE",
// 						"replicationFactor": 1,
// 						"hashedValue":       stringifiedHash,
// 						"data":              secret,
// 					},
// 				})
// 			}
// 		}

// 		return nil
// 	}

// 	// Ask the subsequent node to replicate
// }

// //  EventualReplication is called by the nodes that will not ACK back to the owner
// func (n *Node) EventualReplication(incomingData *data.Data) error {
// 	// messageType := data.Payload["type"]
// 	secret := incomingData.Payload["data"]
// 	stringifiedHash := incomingData.Payload["hashedValue"].(string)
// 	replicationFactor := incomingData.Payload["replicationFactor"].(int)
// 	secretToAdd := map[string]interface{}{
// 		incomingData.Payload["hashedValue"].(string): incomingData.Payload["data"],
// 	}
// 	err := file.WriteDataFile(n.Pid, secretToAdd)
// 	if replicationFactor == 1 {
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		uIntHashValue, uIntConvertError := strconv.ParseUint(stringifiedHash, 32, 32)
// 		if uIntConvertError != nil {
// 			return uIntConvertError
// 		}
// 		// ownerPid, _ := n.mapHashToPid(uint32(uIntHashValue))
// 		_, nextVirtualNode := n.findNextPid(uint32(uIntHashValue))

// 		nextPid, _ := n.findNextWithPid(nextVirtualNode)
// 		// nextPid, nextNextVirtualNode := n.findNextWithPid(nextVirtualNode)
// 		var currentRepFacor = incomingData.Payload["replicationFactor"]
// 		if n.checkHeartbeat(nextPid) && currentRepFacor == 2 {
// 			n.SendSignal(nextPid, &data.Data{
// 				From: n.Pid,
// 				To:   nextPid,
// 				Payload: map[string]interface{}{
// 					"type":              "EVENTUAL_STORE",
// 					"replicationFactor": 1, // TODO: Need to implement R-2 from config.go
// 					"hashedValue":       stringifiedHash,
// 					"data":              secret,
// 				},
// 			})
// 		}
// 		// else {
// 		// 	nextNextPid, _ := n.findNextWithPid(nextNextVirtualNode)
// 		// 	n.SendSignal(nextNextPid, &data.Data{
// 		// 		From: n.Pid,
// 		// 		To:   nextNextPid,
// 		// 		Payload: map[string]interface{}{
// 		// 			"type":              "EVENTUAL_STORE",
// 		// 			"replicationFactor": 1,
// 		// 			"hashedValue":       stringifiedHash,
// 		// 			"data":              secret,
// 		// 		},
// 		// 	})
// 		// }
// 		return nil
// 	}
// 	return nil
// }

// // Takes in virtual node pid to return next pid & virtual node name
// func (n *Node) findNextWithPid(virtualNodePid string) (int, string) {
// 	var currentLocation uint32
// 	var nextLocation uint32
// 	var nextVirtualNodeName string
// 	var next_pid int
// 	var err error
// 	for k, v := range n.VirtualNodeMap {
// 		if v == virtualNodePid {
// 			currentLocation = uint32(k)
// 		}
// 	}

// 	for idx, location := range n.VirtualNodeLocation {
// 		if location == int(currentLocation) {
// 			nextLocation = uint32(n.VirtualNodeLocation[idx+1])
// 		}
// 	}
// 	nextVirtualNodeName = n.VirtualNodeMap[int(nextLocation)]
// 	string_list := strings.Split(nextVirtualNodeName, "-")
// 	next_pid, err = strconv.Atoi(string_list[0])
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return next_pid, nextVirtualNodeName

// }

// // Takes in hasedValue to find the next node's pid and next node's virtual node name
// // func (n *Node) findNextPid(hashedValue uint32) (int, string) {
// // 	var nextVirtualNode string
// // 	var nextPid int
// // 	var err error
// // 	for idx, location := range n.VirtualNodeLocation {
// // 		if int(hashedValue) < location {
// // 			// current_virtual_node := n.RingMap[x]
// // 			nextVirtualNode = n.VirtualNodeMap[n.Ring[(idx+1)]]
// // 			string_list := strings.Split(nextVirtualNode, "-")
// // 			nextPid, err = strconv.Atoi(string_list[0])
// // 			if err != nil {
// // 				fmt.Println(err)
// // 			}
// // 			break

// // 		}
// // 	}
// // 	return nextPid, nextVirtualNode
// // }

// func (n *Node) checkHeartbeat(pid int) bool {
// 	return n.HeartBeatTable[pid]

// }

// func (n *Node) mapHashToPid(hashedValue uint32) (int, string) {
// 	// n := Node
// 	var pid int
// 	var err error
// 	var virtual_node_name string
// 	for location := range n.VirtualNodeLocation {

// 		if int(hashedValue) < location {

// 			virtual_node_name = n.VirtualNodeMap[location]
// 			string_list := strings.Split(virtual_node_name, "-")
// 			pid, err = strconv.Atoi(string_list[0])
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			break
// 		}
// 		continue
// 	}
// 	return pid, virtual_node_name
// }

// func (n *Node) ackCoordinator(incomingData *data.Data) error {
// 	// ack to owner node
// 	// owner id is the coordinator id
// 	coordinatorID := incomingData.Payload["coordinatorID"].(int)
// 	n.SendSignal(coordinatorID, &data.Data{
// 		From: n.Pid,
// 		To:   coordinatorID,
// 		Payload: map[string]interface{}{
// 			"type": "ACK_COORDINATOR",
// 			"data": nil,
// 		},
// 	})
// 	return nil
// }

// // Replication function
// // Function to check if the subsequent nodes are alive
// // SendSignal to each of the subsequent nodes to replicate in them
// // Strict consistency on the first
// // And eventual on the rest
