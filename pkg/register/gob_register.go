package register

import (
	"encoding/gob"

	"github.com/xmliszt/e-safe/pkg/secret"
)

// Register the data type that can be used in RPC messaging
// gob needs this to handle encoding and decoding of various types
var (
	heartbeatTable  = map[int]bool{}
	rpcMap          = map[int]string{}
	virtualNodeData = map[string]interface{}{}
	secretType      = secret.Secret{}
	secretsMap      = map[string]*secret.Secret{}
	secretsList     = []*secret.Secret{}
)

func Regsiter() {
	gob.Register(heartbeatTable)
	gob.Register(rpcMap)
	gob.Register(virtualNodeData)
	gob.Register(secretType)
	gob.Register(secretsMap)
	gob.Register(secretsList)
}
