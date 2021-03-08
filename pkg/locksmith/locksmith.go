package locksmith

import (
	"fmt"
	"time"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/rpc"
)

type LockSmith struct {
	Node *rpc.Node `validate:"required"`
	Nodes []*rpc.Node `validate:"required"`
	HeartBeatTable map[int]bool `validate:"required"`
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
	go locksmithServer.CheckHeartbeat()	// Start periodically checking Node's heartbeat
	locksmithServer.HandleMessageReceived()	// Run this as the main go routine, so do not need to create separate go routine

	return nil
}

// InitializeLocksmith initializes the locksmith server object
func InitializeLocksmith() *LockSmith {
	receivingChannel := make(chan *rpc.Data, 1)
	sendingChannel := make(chan *rpc.Data, 1)
	isCoordinator := false
	locksmithServer := &LockSmith{
		Node: &rpc.Node{
			IsCoordinator: &isCoordinator,
			Pid: 0,
			RecvChannel: receivingChannel,
			SendChannel: sendingChannel,
			Ring: make([]int, 0),
			RpcMap: make(map[int]chan *rpc.Data),
		},
		Nodes: make([]*rpc.Node, 0),
		HeartBeatTable: make(map[int]bool),
	}
	locksmithServer.Node.RpcMap[0] = receivingChannel	// Add Locksmith receiving channel to RpcMap
	return locksmithServer
}

// InitializeNodes initializes the number n nodes that Locksmith is going to create
func (locksmith *LockSmith) InitializeNodes(n int) {
	for i := 1; i <= n; i ++ {
		nodeRecvChan := make(chan *rpc.Data, 1)
		nodeSendChan := make(chan *rpc.Data, 1)
		isCoordinator := false
		newNode := &rpc.Node{
			IsCoordinator: &isCoordinator,
			Pid: i,
			RecvChannel: nodeRecvChan,
			SendChannel: nodeSendChan,
		}
		locksmith.Node.Ring = append(locksmith.Node.Ring, i)
		locksmith.Nodes = append(locksmith.Nodes, newNode)
		locksmith.Node.RpcMap[i] = nodeRecvChan
	}
	
	for _, node := range locksmith.Nodes {
		node.Ring = locksmith.Node.Ring
		node.RpcMap = locksmith.Node.RpcMap
	}
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

// StartAllNodes starts up all created nodes
func (locksmith *LockSmith) StartAllNodes() {
	for _, node := range locksmith.Nodes {
		node.Start()
		locksmith.HeartBeatTable[node.Pid] = true
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
			if locksmith.HeartBeatTable[pid] {
				go func(pid int) {
					locksmith.HeartBeatTable[pid] = false
					locksmith.Node.SendSignal(pid, &rpc.Data{
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
	}
	
// TearDown terminates node, closes all channels
func (locksmith *LockSmith) TearDown() {
	close(locksmith.Node.RecvChannel)
	close(locksmith.Node.SendChannel)
	fmt.Printf("Locksmith Server [%d] has terminated!\n", locksmith.Node.Pid)
}

// EndAllNodes starts teardown process of all created nodes
func (locksmit *LockSmith) EndAllNodes() {
	for _, node := range locksmit.Nodes {
		node.TearDown()
	}
}