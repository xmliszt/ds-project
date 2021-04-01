package locksmith

import (
	"log"
	"sort"
	"time"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/api"
	"github.com/xmliszt/e-safe/pkg/data"
	"github.com/xmliszt/e-safe/pkg/rpc"
	"github.com/xmliszt/e-safe/util"
)

type LockSmith struct {
	LockSmithNode *rpc.Node         `validate:"required"`
	Nodes         map[int]*rpc.Node `validate:"required"`
	Coordinator   int
}

// Start is the main function that starts the entire program
func Start() error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	locksmithServer, err := InitializeLocksmith()
	if err != nil {
		return err
	}
	go locksmithServer.HandleMessageReceived() // Run this as the main go routine, so do not need to create separate go routine
	locksmithServer.InitializeNodes(config.Number)
	log.Println("Locksmith [0] has started")

	e := locksmithServer.StartAllNodes()
	if e != nil {
		return e
	}

	// Simulate node failure
	go func() {
		targetNode := locksmithServer.Nodes[locksmithServer.Coordinator]
		time.Sleep(time.Second * 20)
		targetNode.StopRouter()
		targetNode.TearDown()
	}()

	locksmithServer.CheckHeartbeat() // Start periodically checking Node's heartbeat

	return nil
}

// InitializeLocksmith initializes the locksmith server object
func InitializeLocksmith() (*LockSmith, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	receivingChannel := make(chan *data.Data, config.Number*config.VirtualNodesCount)
	isCoordinator := false
	locksmithServer := &LockSmith{
		LockSmithNode: &rpc.Node{
			IsCoordinator:       &isCoordinator,
			Pid:                 0,
			RecvChannel:         receivingChannel,
			Ring:                make([]int, 0),
			RpcMap:              make(map[int]chan *data.Data),
			VirtualNodeLocation: make([]int, 0),
			VirtualNodeMap:      make(map[int]string),
			HeartBeatTable:      make(map[int]bool),
		},
		Nodes: make(map[int]*rpc.Node),
	}
	locksmithServer.LockSmithNode.RpcMap[0] = receivingChannel // Add Locksmith receiving channel to RpcMap
	return locksmithServer, nil
}

// InitializeNodes initializes the number n nodes that Locksmith is going to create
func (locksmith *LockSmith) InitializeNodes(n int) {
	for i := 1; i <= n; i++ {
		router := api.GetRouter()
		nodeRecvChan := make(chan *data.Data, 1)
		isCoordinator := false
		newNode := &rpc.Node{
			IsCoordinator:  &isCoordinator,
			Pid:            i,
			RecvChannel:    nodeRecvChan,
			HeartBeatTable: make(map[int]bool),
			Router:         router,
		}
		locksmith.LockSmithNode.Ring = append(locksmith.LockSmithNode.Ring, i)
		locksmith.Nodes[i] = newNode
		locksmith.LockSmithNode.RpcMap[i] = nodeRecvChan
	}

	for _, node := range locksmith.Nodes {
		node.Ring = locksmith.LockSmithNode.Ring
		node.RpcMap = locksmith.LockSmithNode.RpcMap
	}
}

// HandleMessageReceived is a Go Routine to handle the messages received
func (locksmith *LockSmith) HandleMessageReceived() {
	for msg := range locksmith.LockSmithNode.RecvChannel {
		switch msg.Payload["type"] {
		case "REPLY_HEARTBEAT":
			locksmith.LockSmithNode.HeartBeatTable[msg.From] = true
		case "UPDATE_VIRTUAL_NODE":
			location := int(msg.Payload["locationData"].(uint32))
			virtualNode := msg.Payload["virtualNodeData"]
			// log.Println("Locksmith has received from ", virtualNode.(string))
			// Update its own values
			locksmith.LockSmithNode.VirtualNodeLocation = append(locksmith.LockSmithNode.VirtualNodeLocation, location)
			locksmith.LockSmithNode.VirtualNodeMap[location] = virtualNode.(string)

			// Sort the location array
			sort.Ints(locksmith.LockSmithNode.VirtualNodeLocation)

			// log.Printf("---Map of virtual node's 'Location' : 'Virtual Node Id'---\n%v\n---Array of virtual node's location---\n%v\n", locksmith.LockSmithNode.VirtualNodeMap, locksmith.LockSmithNode.VirtualNodeLocation)
			// Broadcast to other nodes
			for _, pid := range locksmith.LockSmithNode.Ring {
				locksmith.LockSmithNode.SendSignal(pid, &data.Data{
					From: locksmith.LockSmithNode.Pid,
					To:   pid,
					Payload: map[string]interface{}{
						"type":            "BROADCAST_VIRTUAL_NODES",
						"virtualNodeData": locksmith.LockSmithNode.VirtualNodeMap,
						"locationData":    locksmith.LockSmithNode.VirtualNodeLocation,
					},
				})
			}
		}
	}
}

// StartAllNodes starts up all created nodes
func (locksmith *LockSmith) StartAllNodes() error {
	for pid, node := range locksmith.Nodes {
		err := node.Start()
		if err != nil {
			return err
		}
		locksmith.LockSmithNode.HeartBeatTable[pid] = true
	}
	coordinator := util.FindMax(locksmith.LockSmithNode.Ring)
	locksmith.Coordinator = coordinator
	// Send message to node to turn coordinator field to true
	locksmith.LockSmithNode.SendSignal(coordinator, &data.Data{
		From: locksmith.LockSmithNode.Pid,
		To:   coordinator,
		Payload: map[string]interface{}{
			"type": "YOU_ARE_COORDINATOR",
			"data": nil,
		},
	})
	return nil
}

// CheckHeartbeat periodically check if node is alive
func (locksmith *LockSmith) CheckHeartbeat() {
	config, err := config.GetConfig()
	if err != nil {
		log.Println("Fatal: Heartbeat checking has crashed. Reason: ", err)
		return
	}
	for {
		for _, pid := range locksmith.LockSmithNode.Ring {
			time.Sleep(time.Second * time.Duration(config.HeartbeatInterval))
			if locksmith.LockSmithNode.HeartBeatTable[pid] {
				go func(pid int) {
					locksmith.LockSmithNode.HeartBeatTable[pid] = false
					locksmith.LockSmithNode.SendSignal(pid, &data.Data{
						From: locksmith.LockSmithNode.Pid,
						To:   pid,
						Payload: map[string]interface{}{
							"type": "CHECK_HEARTBEAT",
							"data": nil,
						},
					})
					time.Sleep(time.Second * 1)
					// log.Println("Heartbeat Table: ", locksmith.LockSmithNode.HeartBeatTable)
					if !locksmith.LockSmithNode.HeartBeatTable[pid] {
						time.Sleep(time.Second * time.Duration(config.HeartBeatTimeout))
						if !locksmith.LockSmithNode.HeartBeatTable[pid] {
							log.Printf("Node [%d] is dead! Need to create a new node!\n", pid)
							// log.Println("LOCKSMITH", locksmith.LockSmithNode.VirtualNodeMap)
							// log.Println(locksmith.Nodes[pid].VirtualNodeMap)
							// Election process
							if *locksmith.Nodes[pid].IsCoordinator {
								locksmith.Election()
							}

							// Teardown the particular node in the Nodes
							delete(locksmith.Nodes, pid)

							// Send heartbeat table to all nodes
							locksmith.BroadcastHeartbeatTable()

							// Check and Restart all dead nodes
							locksmith.DeadNodeChecker()

							// Send heartbeat table to all nodes
							locksmith.BroadcastHeartbeatTable()

							time.Sleep(time.Second * time.Duration(config.NodeCreationTimeout)) // allow sufficient time for node to restart, then resume heartbeat checking
						}
					}
				}(pid)
			}
		}
	}
}

// Send heartbeat table to all nodes
func (locksmith *LockSmith) BroadcastHeartbeatTable() {
	for _, pid := range locksmith.LockSmithNode.Ring {
		locksmith.LockSmithNode.SendSignal(pid, &data.Data{
			From: locksmith.LockSmithNode.Pid,
			To:   pid,
			Payload: map[string]interface{}{
				"type": "UPDATE_HEARTBEAT",
				"data": locksmith.LockSmithNode.HeartBeatTable,
			},
		})
	}
}

func (locksmith *LockSmith) DeadNodeChecker() {
	for k, v := range locksmith.LockSmithNode.HeartBeatTable {
		if !v {
			locksmith.SpawnNewNode(k)
			log.Printf("Node [%d] has been revived!\n", k)
		}
	}
}

// Elect the highest surviving Pid node to be coordinator
func (locksmith *LockSmith) Election() {
	var potentialCandidate []int

	for k, v := range locksmith.LockSmithNode.HeartBeatTable {
		if v {
			potentialCandidate = append(potentialCandidate, k)
		}
	}

	coordinator := util.FindMax(potentialCandidate)
	locksmith.Coordinator = coordinator

	// Send message to node to turn coordinator field to true
	locksmith.LockSmithNode.SendSignal(coordinator, &data.Data{
		From: locksmith.LockSmithNode.Pid,
		To:   coordinator,
		Payload: map[string]interface{}{
			"type": "YOU_ARE_COORDINATOR",
			"data": nil,
		},
	})

	log.Printf("Node [%d] is currently the newly elected coordinator!\n", locksmith.Coordinator)
}

// Spawn new nodes when a node is down
func (locksmith *LockSmith) SpawnNewNode(n int) {
	router := api.GetRouter()
	nodeRecvChan := make(chan *data.Data, 1)
	isCoordinator := false
	newNode := &rpc.Node{
		IsCoordinator: &isCoordinator,
		Pid:           n,
		RecvChannel:   nodeRecvChan,
		Router:        router,
	}

	locksmith.Nodes[n] = newNode
	locksmith.LockSmithNode.RpcMap[n] = nodeRecvChan

	locksmith.Nodes[n].StartDeadNode()
	locksmith.LockSmithNode.HeartBeatTable[n] = true

	// Update ring
	found := util.IntInSlice(locksmith.LockSmithNode.Ring, n)
	if !found {
		locksmith.LockSmithNode.Ring = append(locksmith.LockSmithNode.Ring, n)
	}

	// Update node
	for _, node := range locksmith.Nodes {
		node.Ring = locksmith.LockSmithNode.Ring
		node.RpcMap = locksmith.LockSmithNode.RpcMap
	}

}

// TearDown terminates node, closes all channels
func (locksmith *LockSmith) TearDown() {
	close(locksmith.LockSmithNode.RecvChannel)
	log.Printf("Locksmith Server [%d] has terminated!\n", locksmith.LockSmithNode.Pid)
}

// EndAllNodes starts teardown process of all created nodes
func (locksmith *LockSmith) EndAllNodes() {
	for _, node := range locksmith.Nodes {
		node.TearDown()
	}
}
