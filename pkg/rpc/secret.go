package rpc

import (
	"log"

	"github.com/xmliszt/e-safe/config"
)

// Coordinator node to handle deleting a secret
func (n *Node) DeleteSecret(alias string) error {
	// Hash alias into location token
	uloc, err := config.GetHash(alias)
	if err != nil {
		return err
	}

	loc := int(uloc)
	var targetVirtualNodeName string
	for _, virtualLoc := range n.VirtualNodeLocation {
		if virtualLoc >= loc {
			targetVirtualNodeName = n.VirtualNodeMap[virtualLoc]
		}
	}
	log.Println(targetVirtualNodeName)
	return nil
}
