package rpc

import (
	"fmt"
)

type Secret struct {
	Key   int
	Value string // Value of secret
	Role  int    // Int to identify role (clearance level)
}

func (n Node) getSecret(key int) {
	fmt.Println("Hellu Warld")
}

func (n Node) getSecrets(from int, to int) {}

func (n Node) putSecret(key int, secret Secret) {}

func (n Node) deleteSecret(key int) {}
