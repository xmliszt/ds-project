package message

import (
	"net/rpc"
)

const (
	UPDATE_HEARTBEAT_TABLE = 0
	GET_HEARTBEAT_TABLE    = 1
	ASSIGN_COORDINATOR     = 2
)

type Reply struct {
	From    int         // Node ID of which the reply is sent from
	To      int         // Node ID of which the reply is sent to
	ReplyTo int         // RPC code of the request that this reply is directed to
	Payload interface{} // The content of the reply
}

type Request struct {
	From    int         // Node ID of which the request is sent from
	To      int         // Node ID of which the request is sent to
	Code    int         // The requested RPC code
	Payload interface{} // The content of the request
}

// SendMessage delivers an RPC message to target address, with given method to call and parameters
func SendMessage(address string, method string, request *Request, reply *Reply) error {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return err
	}
	err = client.Call(method, request, reply)
	if err != nil {
		client.Close()
		return err
	}
	client.Close()
	return nil
}
