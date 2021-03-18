package rpc

import (
	"fmt"
)

type Signal interface {
	SendSignal(pid int, data map[string]interface{})
}

// SendSignal used for sending request to target Node
func (n Node) SendSignal(pid int, data *Data) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	fmt.Printf("[%d] -> [%d]: %v\n", n.Pid, pid, data)
	sendingChannel := n.RpcMap[pid]
	sendingChannel <- data
}
