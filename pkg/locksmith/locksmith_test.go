package locksmith

import (
	"testing"
	"time"

	"github.com/xmliszt/e-safe/pkg/data"
	"github.com/xmliszt/e-safe/pkg/rpc"
	"gopkg.in/go-playground/validator.v9"
)

// Expected a locksmith server to be created with all required fields
func TestInitializeLocksmith(t *testing.T) {
	validate := validator.New()
	locksmith := InitializeLocksmith()
	err := validate.Struct(locksmith)
	if err != nil {
		t.Error(err)
	}
}

// Expected n nodes to be initialized with all required fields
func TestInitializeNodes(t *testing.T) {
	locksmith := &LockSmith{
		LockSmithNode: &rpc.Node{
			Ring:   make([]int, 0),
			RpcMap: make(map[int]chan *data.Data),
		},
		Nodes: make(map[int]*rpc.Node),
	}
	locksmith.InitializeNodes(3)
	if len(locksmith.Nodes) < 3 || len(locksmith.LockSmithNode.Ring) < 3 || len(locksmith.LockSmithNode.RpcMap) < 3 {
		t.Errorf("Expected 3 nodes to be created, but have incomplete creation: %d", len(locksmith.Nodes))
	}
	for _, node := range locksmith.Nodes {
		validate := validator.New()
		err := validate.Struct(node)
		if err != nil {
			t.Error(err)
		}
	}
}

// Expected HeartbeatTable to update to true when receive a heartbeat reply
func TestHandleMessageReceived(t *testing.T) {
	receivingChannel := make(chan *data.Data, 1)
	locksmith := &LockSmith{
		LockSmithNode: &rpc.Node{
			RecvChannel: receivingChannel,
		},
		HeartBeatTable: make(map[int]bool),
	}
	locksmith.HeartBeatTable[1] = false
	go locksmith.HandleMessageReceived()
	locksmith.LockSmithNode.RecvChannel <- &data.Data{From: 1, Payload: map[string]interface{}{"type": "REPLY_HEARTBEAT"}}
	time.Sleep(time.Second * 1)
	if !locksmith.HeartBeatTable[1] {
		t.Errorf("Expected HeartbeatTable for Node 1 to be true, but instead it is still false.")
	}
}

// Expected 3 nodes to spin up and heartbeat table all update to true
func TestStartAllNodes(t *testing.T) {
	locksmith := &LockSmith{
		LockSmithNode: &rpc.Node{
			Pid:    0,
			Ring:   make([]int, 0),
			RpcMap: make(map[int]chan *data.Data),
		},
	}
	locksmith.Nodes = make(map[int]*rpc.Node)
	locksmith.HeartBeatTable = make(map[int]bool)
	iscoordinator := false
	for i := 1; i <= 3; i++ {
		newNode := &rpc.Node{
			Pid:           i,
			IsCoordinator: &iscoordinator,
			RecvChannel:   make(chan *data.Data),
		}
		locksmith.Nodes[i] = newNode
		locksmith.LockSmithNode.RpcMap[i] = newNode.RecvChannel
		locksmith.LockSmithNode.Ring = append(locksmith.LockSmithNode.Ring, i)
	}
	locksmith.StartAllNodes()
	for pid, alive := range locksmith.HeartBeatTable {
		if !alive {
			t.Errorf("Expected Node [%d] to be alive, but yet it is not alive!", pid)
		}
	}
}

// Expected length of Nodes to increase after spawning new node
func TestSpawnNewNode(t *testing.T) {
	locksmith := &LockSmith{
		LockSmithNode: &rpc.Node{
			Ring:   make([]int, 0),
			RpcMap: make(map[int]chan *data.Data),
		},
		Nodes: map[int]*rpc.Node{
			1: {},
			2: {},
			3: {},
		},
		HeartBeatTable: map[int]bool{
			1: true,
			2: true,
			3: true,
			4: false,
		},
	}

	initLenLocksmithNodes := len(locksmith.Nodes)

	locksmith.SpawnNewNode(initLenLocksmithNodes + 1)
	LenLocksmithNodes := len(locksmith.Nodes)

	if LenLocksmithNodes == initLenLocksmithNodes {
		t.Errorf("Expected to have %d number of nodes, but only have %d!", LenLocksmithNodes, initLenLocksmithNodes)
	}
}

func TestElection(t *testing.T) {
	mockChan := make(chan *data.Data)
	iscoordinator := false
	locksmith := &LockSmith{
		LockSmithNode: &rpc.Node{
			Ring: make([]int, 0),
			RpcMap: map[int]chan *data.Data{
				3: mockChan,
			},
		},
		Nodes: map[int]*rpc.Node{
			1: {},
			2: {},
			3: {
				RecvChannel:   mockChan,
				IsCoordinator: &iscoordinator,
			},
			4: {},
		},
		HeartBeatTable: map[int]bool{
			1: true,
			2: true,
			3: true,
			4: false,
		},
	}
	locksmith.Nodes[3].Start()
	locksmith.Election()
	time.Sleep(time.Second * 2)
	if !*locksmith.Nodes[3].IsCoordinator {
		t.Error("Node 3 supposed to be coordinator but it is not!")
	}
}
