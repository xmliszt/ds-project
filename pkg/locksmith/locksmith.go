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
	NodeNumber int
	HeartBeatTable map[int]bool
}

// InitializeLocksmith initializes the Locksmith Server
func InitializeLocksmith() error {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}
	n := config.Number
	receivingChannel := make(chan rpc.Data, n+1)
	sendingChannel := make(chan rpc.Data, n+1)
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
		nodeRecvChan := make(chan rpc.Data, n+1)
		nodeSendChan := make(chan rpc.Data, n+1)
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
		fmt.Printf("Node [%d] recv channel: %v\n", node.Pid, node.RecvChannel)
	}

	locksmithServer.StartAllNodes()	// Spin up all created nodes

	go locksmithServer.CheckHeartbeat(config.HeartbeatInterval)	// Start periodically checking Node's heartbeat
	
	locksmithServer.HandleMessageReceived()	// Run this as the main go routine, so do not need to create separate go routine
	return nil
}

// HandleMessageReceived is a Go Routine to handle the messages received
func (n *LockSmith) HandleMessageReceived() {
	for {
		select {
		case msg, ok := <-n.Node.RecvChannel:
			if ok {
				switch msg.Payload["type"] {
				case "REPLY_HEARTBEAT":
					alive := msg.Payload["data"].(bool)
					if !alive {
						n.updateHeartbeatTable(msg.From, alive)
					}
				}
			} else {
				continue
			}
		default:
			continue
		}
	}
}

// TearDown terminates node, closes all channels
func (n *LockSmith) TearDown() {
	close(n.Node.RecvChannel)
	close(n.Node.SendChannel)
	fmt.Printf("Locksmith Server [%d] has terminated!\n", n.Node.Pid)
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
func (locksmith *LockSmith) CheckHeartbeat(interval int) {
	for _, pid := range locksmith.Node.Ring {
		time.Sleep(time.Second * time.Duration(interval))
		locksmith.Node.SendSignal(pid, rpc.Data{
			From: locksmith.Node.Pid,
			To: pid,
			Payload: map[string]interface{}{
				"type": "CHECK_HEARTBEAT",
				"data": nil,
			},
		})
	}
}

// updateHeartbeatTable updates the heartbeat table
func (locksmith *LockSmith) updateHeartbeatTable(pid int, val bool) {
	table := locksmith.HeartBeatTable
	fmt.Println(table, locksmith.HeartBeatTable)
	if oldVal, ok := table[pid]; ok {
		if oldVal != val {
			table[pid] = val
		}
	}
}