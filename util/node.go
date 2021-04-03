package util

import (
	"strconv"
	"strings"
)

// GetPhysicalNodeID gets the physical node ID that the virtual node name belongs to
func GetPhysicalNodeID(virtualNodeName string) (int, error) {
	parts := strings.Split(virtualNodeName, "-")
	nodeID, err := strconv.Atoi(parts[0])
	if err != nil {
		return -1, err
	}
	return nodeID, nil
}
