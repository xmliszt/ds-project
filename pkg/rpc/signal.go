package rpc

import (
	"github.com/xmliszt/e-safe/pkg/data"
)

type Signal interface {
	SendSignal(pid int, data map[string]interface{})
}

// SendSignal used for sending request to target Node
func (n *Node) SendSignal(pid int, data *data.Data) {
	defer func() {
		recover()
	}()
	// fmt.Printf("[%d] -> [%d]: %v\n", n.Pid, pid, data)
	sendingChannel := n.RpcMap[pid]
	sendingChannel <- data
}
