package rpc

// Signal interface has all the method signatures necessary for Signals
type Signal interface {
	SendSignal(pid int, data map[string]interface{})
}

// SendSignal used for sending request to target Node
// (n Node ) in front of the func means the func is a receiver for the node, only node can call it
// 1. send election signal to new coordinator
// 2. send check heartbeat signals to each node
// 3. pass updated heartbeat table to coordinator node
func (n Node) SendSignal(pid int, data map[string]interface{}) {
	sendingChannel := n.RPCMap[pid]
	data["from"] = n.Pid
	sendingChannel <- data
}
