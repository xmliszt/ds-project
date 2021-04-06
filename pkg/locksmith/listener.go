package locksmith

import (
	"fmt"
	"sort"

	"github.com/xmliszt/e-safe/pkg/message"
)

// SignalStart receives the signal from sending node and return RPC map to the node
func (locksmith *LockSmith) SignalStart(request *message.Request, reply *message.Reply) error {
	newNodeID := request.From
	locksmith.RpcMap[newNodeID] = fmt.Sprintf("localhost:%d", 5000+newNodeID)
	*reply = message.Reply{
		From:    locksmith.Pid,
		To:      request.From,
		ReplyTo: request.Code,
		Payload: locksmith.RpcMap,
	}
	return nil
}

// GetHeartbeatTable returns the heartbeat table as a RPC reply
func (locksmith *LockSmith) CreateVirtualNodes(request *message.Request, reply *message.Reply) error {
	payload := request.Payload.(map[string]interface{})
	ilocations := payload["virtualNodeLocation"]
	ivirtualNodesMap := payload["virtualNodeMap"]

	locations := ilocations.([]int)
	virtualNodesMap := ivirtualNodesMap.(map[int]string)

	// Update its own values
	locksmith.VirtualNodeLocation = append(locksmith.VirtualNodeLocation, locations...)

	for loc, vNodeName := range virtualNodesMap {
		locksmith.VirtualNodeMap[loc] = vNodeName
	}

	// Sort the location array
	sort.Ints(locksmith.VirtualNodeLocation)

	*reply = message.Reply{
		From:    locksmith.Pid,
		To:      request.From,
		ReplyTo: request.Code,
		Payload: map[string]interface{}{
			"virtualNodeMap":      locksmith.VirtualNodeMap,
			"virtualNodeLocation": locksmith.VirtualNodeLocation,
		},
	}
	return nil
}

// AcquireUserLock puts the request into the queue and serves the request one by one in FIFO order
func (locksmith *LockSmith) AcquireUserLock(request *message.Request, reply *message.Reply) error {
	requestNodeID := request.From
	// Add request to queue
	locksmith.RequestQueue = append(locksmith.RequestQueue, requestNodeID)
	for {
		if locksmith.RequestQueue[0] == requestNodeID {
			break
		}
	}
	*reply = message.Reply{
		From:    locksmith.Pid,
		To:      requestNodeID,
		ReplyTo: request.Code,
		Payload: nil,
	}
	return nil
}

// ReleaseUserLock pop the top request from the queue
func (locksmith *LockSmith) ReleaseUserLock(request *message.Request, reply *message.Reply) error {
	locksmith.RequestQueue = locksmith.RequestQueue[1:]
	*reply = message.Reply{
		From:    locksmith.Pid,
		To:      request.From,
		ReplyTo: request.Code,
		Payload: nil,
	}
	return nil
}
