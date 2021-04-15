package message

import (
	"net/rpc"
)

const (
	SIGNAL_START                = 0
	UPDATE_HEARTBEAT_TABLE      = 1
	GET_HEARTBEAT_TABLE         = 2
	ASSIGN_COORDINATOR          = 3
	REMOVE_COORDINATOR          = 4
	CREATE_VIRTUAL_NODE         = 5
	UPDATE_VIRTUAL_NODES        = 6
	UPDATE_RPC_MAP              = 7
	STORE_AND_REPLICATE         = 8
	STRICT_STORE                = 9  // Take list of v.nodes and send it forward R-1 nodes
	EVENTUAL_STORE              = 10 // Sent by the neighbouring node to the next R-1 nodes
	STRICT_OWNER_DOWN           = 14 // When Owner is down, Sent to Strict Con Node\
	RELAY_DELETE_SECRET         = 19
	ACQUIRE_USER_LOCK           = 20
	RELEASE_USER_LOCK           = 21
	FETCH_ORIGINAL_SECRETS      = 22
	FETCH_REPLICA_SECRETS       = 23
	GIVE_ME_DATA                = 24
	DELETE_SECRET_ALL_INSTANCES = 25
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
