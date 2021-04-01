package locksmith

import "github.com/xmliszt/e-safe/pkg/message"

// GetHeartbeatTable returns the heartbeat table as a RPC reply
func (locksmith *LockSmith) GetHeartbeatTable(request *message.Request, reply *message.Reply) error {
	*reply = message.Reply{
		From:    locksmith.Pid,
		To:      request.From,
		ReplyTo: request.Code,
		Payload: locksmith.HeartBeatTable,
	}
	return nil
}
