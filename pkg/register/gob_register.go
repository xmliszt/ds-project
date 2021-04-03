package register

import "encoding/gob"

// Register the data type that can be used in RPC messaging
// gob needs this to handle encoding and decoding of various types
var (
	heartbeatTable  = map[int]bool{}
	rpcMap          = map[int]string{}
	virtualNodeData = map[string]interface{}{}
)

func Regsiter() {
	gob.Register(heartbeatTable)
	gob.Register(rpcMap)
	gob.Register(virtualNodeData)
}
