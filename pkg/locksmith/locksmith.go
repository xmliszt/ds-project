package locksmith

import (
	"fmt"
	"sync"

	"github.com/xmliszt/e-safe/pkg/rpc"
)

type LockSmith struct {
	Node rpc.Node
	Nodes []*rpc.Node
	NodeNumber int
	HeartBeatTable map[int]bool
}

// Initialize Locksmith Server
func InitializeLocksmith(n int) *LockSmith {
	wg := &sync.WaitGroup{}
	receivingChannel := make(chan map[string]interface{}, 10)
	sendingChannel := make(chan map[string]interface{}, 10)
	locksmithServer := &LockSmith{
		Node: rpc.Node{
			Wg: wg,
			IsCoordinator: false,
			Pid: 0,
			RecvChannel: receivingChannel,
			SendChannel: sendingChannel,
		},
		NodeNumber: n,
	}

	locksmithServer.HeartBeatTable = make(map[int]bool)

	ring := make([]int, 0)
	rpcMap := make(map[int]chan map[string]interface{})

	rpcMap[0] = receivingChannel	// Add Locksmith receiving channel to RpcMap

	for i := 0; i < n; i ++ {
		nodeRecvChan := make(chan map[string]interface{}, 10)
		nodeSendChan := make(chan map[string]interface{}, 10)
		newNode := &rpc.Node{
			Wg: wg,
			IsCoordinator: false,
			Pid: i+1,
			RecvChannel: nodeRecvChan,
			SendChannel: nodeSendChan,
		}
		ring = append(ring, i+1)
		locksmithServer.Nodes = append(locksmithServer.Nodes, newNode)
		rpcMap[i+1] = nodeRecvChan
	}

	locksmithServer.Node.Ring = ring
	locksmithServer.Node.RpcMap = rpcMap
	
	for _, node := range locksmithServer.Nodes {
		node.Ring = ring
		node.RpcMap = rpcMap
	}

	return locksmithServer
}


// Call this function at the initialization
// Start up all created nodes
func (locksmith *LockSmith) StartAllNodes() {
	for _, node := range locksmith.Nodes {
		locksmith.Node.Wg.Add(1)
		node.Start()
		locksmith.HeartBeatTable[node.Pid] = true
	}
}

// Wait for all nodes to finish and exit
func (locksmith *LockSmith) MonitorNodes() {
	locksmith.Node.Wg.Wait()
	close(locksmith.Node.RecvChannel)
	close(locksmith.Node.SendChannel)
	fmt.Println("Teardown complete! ")
}