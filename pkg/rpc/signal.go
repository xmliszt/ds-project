package rpc

import (
	"fmt"
)

type Signal interface {
	SendSignal(pid int, data map[string]interface{}) 
}

// Sending request to target Node
// (n Node ) in front of the func means the func is a receiver for the node, only node can call it
// 1. send election signal to new coordinator
// 2. send check heartbeat signals to each node
// 3. pass updated heartbeat table to coordinator node
func (n Node) SendSignal(pid int, data Data) {
	defer func(){
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	fmt.Printf("[%d] -> [%d]: %v\n", n.Pid, pid, data)
	sendingChannel := n.RpcMap[pid]
	sendingChannel <- data
}
