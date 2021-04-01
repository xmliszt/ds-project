package rpc

import (
	"testing"
	"time"

	"github.com/xmliszt/e-safe/pkg/data"
)

// Expected node to perform correct action upon receiving a particular type of message
func TestHandleMessageReceived(t *testing.T) {
	// Test for reply heartbeat
	mockRecvChannel := make(chan *data.Data, 1)
	node := &Node{
		Pid:         999,
		RecvChannel: make(chan *data.Data, 1),
		RpcMap:      make(map[int]chan *data.Data),
		// HeartBeatTable: nodeHeartbeat,
	}
	node.RpcMap[0] = mockRecvChannel
	go func() {
		node.RecvChannel <- &data.Data{From: 0, Payload: map[string]interface{}{"type": "CHECK_HEARTBEAT"}}
	}()
	go node.HandleMessageReceived()
	msg := <-mockRecvChannel
	if msg.From != 999 || msg.To != 0 || msg.Payload["type"] != "REPLY_HEARTBEAT" {
		t.Error("Node failed to handle CHECK_HEARTBEAT message. Received: ", msg)
	}

	// Test for update heartbeat
	nodeHeartbeat := map[int]bool{
		1: false,
		2: false,
		3: false,
	}
	mockRecvChannel2 := make(chan *data.Data, 1)
	node2 := &Node{
		Pid:            998,
		RecvChannel:    make(chan *data.Data, 1),
		RpcMap:         make(map[int]chan *data.Data),
		HeartBeatTable: nodeHeartbeat,
	}
	node2.RpcMap[0] = mockRecvChannel2
	heartbeatTable := map[int]bool{
		1: true,
		2: true,
		3: true,
	}
	go func() {
		node2.RecvChannel <- &data.Data{
			From: 0,
			To:   998,
			Payload: map[string]interface{}{
				"type": "UPDATE_HEARTBEAT",
				"data": heartbeatTable,
			},
		}
	}()
	go node2.HandleMessageReceived()

	time.Sleep(time.Second * 2)
	if heartbeatTable[1] != node2.HeartBeatTable[1] || heartbeatTable[2] != node2.HeartBeatTable[2] || heartbeatTable[3] != node2.HeartBeatTable[3] {
		t.Error("Node failed to handle UPDATE_HEARTBEAT message. Received: ", node2.HeartBeatTable)
	}
}
