package locksmith

import (
	"fmt"
	"time"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/rpc"
)

type LockSmith struct {
	Node rpc.Node
	Nodes []*rpc.Node
	HeartBeatTable map[int]bool
}

// InitializeLocksmith initializes the Locksmith Server
func InitializeLocksmith() error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	n := config.Number
	receivingChannel := make(chan rpc.Data, 1)
	sendingChannel := make(chan rpc.Data, 1)
	locksmithServer := &LockSmith{
		Node: rpc.Node{
			IsCoordinator: false,
			Pid: 0,
			RecvChannel: receivingChannel,
			SendChannel: sendingChannel,
		},
	}

	locksmithServer.HeartBeatTable = make(map[int]bool)

	ring := make([]int, 0)
	rpcMap := make(map[int]chan rpc.Data)

	rpcMap[0] = receivingChannel	// Add Locksmith receiving channel to RpcMap

	for i := 1; i <= n; i ++ {
		nodeRecvChan := make(chan rpc.Data, 1)
		nodeSendChan := make(chan rpc.Data, 1)
		newNode := &rpc.Node{
			IsCoordinator: false,
			Pid: i,
			RecvChannel: nodeRecvChan,
			SendChannel: nodeSendChan,
		}
		ring = append(ring, i)
		locksmithServer.Nodes = append(locksmithServer.Nodes, newNode)
		rpcMap[i] = nodeRecvChan
	}

	locksmithServer.Node.Ring = ring
	locksmithServer.Node.RpcMap = rpcMap
	
	for _, node := range locksmithServer.Nodes {
		node.Ring = ring
		node.RpcMap = rpcMap
	}

	locksmithServer.StartAllNodes()	// Spin up all created nodes
	fmt.Println("Locksmith [0] has started")
	go locksmithServer.CheckHeartbeat()	// Start periodically checking Node's heartbeat
	locksmithServer.HandleMessageReceived()	// Run this as the main go routine, so do not need to create separate go routine

	return nil
}

// HandleMessageReceived is a Go Routine to handle the messages received
func (locksmith *LockSmith) HandleMessageReceived() {
	for msg := range locksmith.Node.RecvChannel {
		switch msg.Payload["type"] {
		case "REPLY_HEARTBEAT":
			locksmith.HeartBeatTable[msg.From] = true
		}
	}
}

// TearDown terminates node, closes all channels
func (locksmith *LockSmith) TearDown() {
	close(locksmith.Node.RecvChannel)
	close(locksmith.Node.SendChannel)
	fmt.Printf("Locksmith Server [%d] has terminated!\n", locksmith.Node.Pid)
}

// StartAllNodes starts up all created nodes
func (locksmith *LockSmith) StartAllNodes() {
	for _, node := range locksmith.Nodes {
		node.Start()
		locksmith.HeartBeatTable[node.Pid] = true
	}
}

// EndAllNodes starts teardown process of all created nodes
func (locksmit *LockSmith) EndAllNodes() {
	for _, node := range locksmit.Nodes {
		node.TearDown()
	}
}

// CheckHeartbeat periodically check if node is alive
func (locksmith *LockSmith) CheckHeartbeat() {
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Fatal: Heartbeat checking has crashed. Reason: ", err)
		return
	}
	for {
		for _, pid := range locksmith.Node.Ring {
			time.Sleep(time.Second * time.Duration(config.HeartbeatInterval))
			go func(pid int) {
				locksmith.HeartBeatTable[pid] = false
				locksmith.Node.SendSignal(pid, rpc.Data{
					From: locksmith.Node.Pid,
					To: pid,
					Payload: map[string]interface{}{
						"type": "CHECK_HEARTBEAT",
						"data": nil,
					},
				})
				time.Sleep(time.Second * 1)
				fmt.Println("Hearbeat Table: ", locksmith.HeartBeatTable)
				if !locksmith.HeartBeatTable[pid] {
					time.Sleep(time.Second * time.Duration(config.HeartBeatTimeout))
					if !locksmith.HeartBeatTable[pid] {
						fmt.Printf("Node [%d] is dead! Need to create a new node!\n", pid)
						time.Sleep(time.Second * time.Duration(config.NodeCreationTimeout))	// allow sufficient time for node to restart, then resume heartbeat checking
					}
				}
			}(pid)
		}
	}
}