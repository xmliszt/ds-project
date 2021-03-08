package locksmith

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/xmliszt/e-safe/pkg/rpc"
)

type LockSmith struct {
	Node rpc.Node
	Nodes []*rpc.Node
	NodeNumber int
	HeartBeatTable map[int]bool
}

// Initialize Locksmith Server
func InitializeLocksmith(n int) {
	receivingChannel := make(chan rpc.Data)
	sendingChannel := make(chan rpc.Data)
	locksmithServer := &LockSmith{
		Node: rpc.Node{
			IsCoordinator: false,
			Pid: 0,
			RecvChannel: receivingChannel,
			SendChannel: sendingChannel,
		},
		NodeNumber: n,
	}

	locksmithServer.HeartBeatTable = make(map[int]bool)

	ring := make([]int, 0)
	rpcMap := make(map[int]chan rpc.Data)

	rpcMap[0] = receivingChannel	// Add Locksmith receiving channel to RpcMap

	for i := 1; i <= n; i ++ {
		nodeRecvChan := make(chan rpc.Data)
		nodeSendChan := make(chan rpc.Data)
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

	for _, pid := range locksmithServer.Node.Ring {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(1000)))
		locksmithServer.Node.SendSignal(pid, rpc.Data{
			From: locksmithServer.Node.Pid,
			To: pid,
			Payload: map[string]interface{}{
				"Hello": "world",
				"Age": 120,
			},
		})
	}

	locksmithServer.HandleMessageReceived()	// Run this as the main go routine, so do not need to create separate go routine
}

// Go Routine to handle the messages received
func (n *LockSmith) HandleMessageReceived() {
	for msg := range n.Node.RecvChannel {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		fmt.Println("Locksmith receive: ", msg)
	}
}

// Terminate node, close all channels
func (n *LockSmith) TearDown() {
	close(n.Node.RecvChannel)
	close(n.Node.SendChannel)
	fmt.Printf("Locksmith Server [%d] has terminated!\n", n.Node.Pid)
}

// Call this function at the initialization
// Start up all created nodes
func (locksmith *LockSmith) StartAllNodes() {
	for _, node := range locksmith.Nodes {
		node.Start()
		locksmith.HeartBeatTable[node.Pid] = true
	}
}

// Start teardown process of all created nodes
func (locksmit *LockSmith) EndAllNodes() {
	for _, node := range locksmit.Nodes {
		node.TearDown()
	}
}