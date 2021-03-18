package locksmith

import (
	"fmt"
	"time"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/data"
	"github.com/xmliszt/e-safe/pkg/rpc"
	"github.com/xmliszt/e-safe/util"
)

type LockSmith struct {
	LockSmithNode  *rpc.Node         `validate:"required"`
	Nodes          map[int]*rpc.Node `validate:"required"`
	HeartBeatTable map[int]bool      `validate:"required"`
	Coordinator    int
}

// Start is the main function that starts the entire program
func Start() error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	locksmithServer := InitializeLocksmith()
	locksmithServer.InitializeNodes(config.Number)
	locksmithServer.StartAllNodes()

	fmt.Println("Locksmith [0] has started")
	go locksmithServer.CheckHeartbeat()     // Start periodically checking Node's heartbeat
	locksmithServer.HandleMessageReceived() // Run this as the main go routine, so do not need to create separate go routine

	return nil
}

// InitializeLocksmith initializes the locksmith server object
func InitializeLocksmith() *LockSmith {
	receivingChannel := make(chan *data.Data, 1)
	sendingChannel := make(chan *data.Data, 1)
	isCoordinator := false
	locksmithServer := &LockSmith{
		LockSmithNode: &rpc.Node{
			IsCoordinator: &isCoordinator,
			Pid:           0,
			RecvChannel:   receivingChannel,
			SendChannel:   sendingChannel,
			Ring:          make([]int, 0),
			RpcMap:        make(map[int]chan *data.Data),
		},
		Nodes:          make(map[int]*rpc.Node),
		HeartBeatTable: make(map[int]bool),
	}
	locksmithServer.LockSmithNode.RpcMap[0] = receivingChannel // Add Locksmith receiving channel to RpcMap
	return locksmithServer
}

// InitializeNodes initializes the number n nodes that Locksmith is going to create
func (locksmith *LockSmith) InitializeNodes(n int) {
	for i := 1; i <= n; i++ {
		nodeRecvChan := make(chan *data.Data, 1)
		nodeSendChan := make(chan *data.Data, 1)
		isCoordinator := false
		newNode := &rpc.Node{
			IsCoordinator: &isCoordinator,
			Pid:           i,
			RecvChannel:   nodeRecvChan,
			SendChannel:   nodeSendChan,
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
			locksmith.HeartBeatTable[msg.From] = true
		}
	}
}

// StartAllNodes starts up all created nodes
func (locksmith *LockSmith) StartAllNodes() {
	for pid, node := range locksmith.Nodes {
		node.Start()
		locksmith.HeartBeatTable[pid] = true
	}
	coordinator := util.FindMax(locksmith.LockSmithNode.Ring)
	// Send message to node to turn coordinator field to true
	locksmith.LockSmithNode.SendSignal(coordinator, &data.Data{
		From: locksmith.LockSmithNode.Pid,
		To:   coordinator,
		Payload: map[string]interface{}{
			"type": "YOU_ARE_COORDINATOR",
			"data": nil,
		},
	})
}

// CheckHeartbeat periodically check if node is alive
func (locksmith *LockSmith) CheckHeartbeat() {
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Fatal: Heartbeat checking has crashed. Reason: ", err)
		return
	}
	for {
		for _, pid := range locksmith.LockSmithNode.Ring {
			time.Sleep(time.Second * time.Duration(config.HeartbeatInterval))
			if locksmith.HeartBeatTable[pid] {
				go func(pid int) {
					locksmith.HeartBeatTable[pid] = false
					locksmith.LockSmithNode.SendSignal(pid, &data.Data{
						From: locksmith.LockSmithNode.Pid,
						To:   pid,
						Payload: map[string]interface{}{
							"type": "CHECK_HEARTBEAT",
							"data": nil,
						},
					})
					time.Sleep(time.Second * 1)
					fmt.Println("Heartbeat Table: ", locksmith.HeartBeatTable)
					if !locksmith.HeartBeatTable[pid] {
						time.Sleep(time.Second * time.Duration(config.HeartBeatTimeout))
						if !locksmith.HeartBeatTable[pid] {
							fmt.Printf("Node [%d] is dead! Need to create a new node!\n", pid)

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
				"data": locksmith.HeartBeatTable,
			},
		})
		fmt.Printf("Node [%d] has updated its heartbeat table from locksmith\n", pid)
	}

}

func (locksmith *LockSmith) DeadNodeChecker() {
	for k, v := range locksmith.HeartBeatTable {
		if !v {
			locksmith.SpawnNewNode(k)
			fmt.Printf("Node [%d] has been revived!\n", k)
		}
	}
}

// Elect the highest surviving Pid node to be coordinator
func (locksmith *LockSmith) Election() {
	var potentialCandidate []int

	for k, v := range locksmith.HeartBeatTable {
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

	// isCoordinator := true
	// locksmith.Nodes[coordinator].IsCoordinator = &isCoordinator

	fmt.Printf("Node [%d] is currently the newly elected coordinator!\n", locksmith.Coordinator)
}

// Spawn new nodes when a node is down
func (locksmith *LockSmith) SpawnNewNode(n int) {
	nodeRecvChan := make(chan *data.Data, 1)
	nodeSendChan := make(chan *data.Data, 1)
	isCoordinator := false
	newNode := &rpc.Node{
		IsCoordinator: &isCoordinator,
		Pid:           n,
		RecvChannel:   nodeRecvChan,
		SendChannel:   nodeSendChan,
	}

	locksmith.Nodes[n] = newNode
	locksmith.LockSmithNode.RpcMap[n] = nodeRecvChan

	locksmith.Nodes[n].Start()
	locksmith.HeartBeatTable[n] = true

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
	close(locksmith.LockSmithNode.SendChannel)
	fmt.Printf("Locksmith Server [%d] has terminated!\n", locksmith.LockSmithNode.Pid)
}

// EndAllNodes starts teardown process of all created nodes
func (locksmith *LockSmith) EndAllNodes() {
	for _, node := range locksmith.Nodes {
		node.TearDown()
	}
}
