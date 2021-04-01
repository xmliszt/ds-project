package rpc

import (
	"log"

	"github.com/xmliszt/e-safe/pkg/data"
)

type Signal interface {
	SendSignal(pid int, data map[string]interface{})
}

// SendSignal used for sending request to target Node
func (n *Node) SendSignal(pid int, data *data.Data) {
	defer func() {
		r := recover()
		if r != nil {
			log.Println("Signal recovery: ", r)
		}
	}()
	// fmt.Printf("[%d] -> [%d]: %v\n", n.Pid, pid, data)
	sendingChannel := n.RpcMap[pid]
	sendingChannel <- data
}
