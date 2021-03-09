package rpc

import (
	"testing"
)

// Expected node to perform correct action upon receiving a particular type of message
func TestHandleMessageReceived(t *testing.T) {
	mockRecvChannel := make(chan *Data, 1)
	node := &Node{
		Pid: 999,
		RecvChannel: make(chan *Data, 1),
		RpcMap: make(map[int]chan *Data),
	}
	node.RpcMap[0] = mockRecvChannel
	go func(){
		node.RecvChannel <- &Data{From: 0, Payload: map[string]interface{}{"type": "CHECK_HEARTBEAT"}}
	}()
	go node.HandleMessageReceived()
	msg := <- mockRecvChannel 
	if msg.From != 999 || msg.To != 0 || msg.Payload["type"] != "REPLY_HEARTBEAT" {
		t.Error("Node failed to handle CHECK_HEARTBEAT message. Received: ", msg)
	}
}
