package node

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/util"
)

// getPhysicalNodeID gets the physical node ID that the virtual node name belongs to
func getPhysicalNodeID(virtualNodeName string) (int, error) {
	parts := strings.Split(virtualNodeName, "-")
	nodeID, err := strconv.Atoi(parts[0])
	if err != nil {
		return -1, err
	}
	return nodeID, nil
}

// getVirtualLocationIndex returns the index where the virtual location
// can be found in the VirtualNodeLocation array
// if not found, return -1
func (n *Node) getVirtualLocationIndex(location int) int {
	for idx, loc := range n.VirtualNodeLocation {
		if loc == location {
			return idx
		}
	}
	return -1
}

// getReplayNodes gets the next n virtual nodes (n: replication factor)
// locations. The relay nodes must satisfy the following criteria:
// 1: the physical node cannot repeat, must be unqiue
// 2: the physical node cannot be the start node
func (n *Node) getRelayVirtualNodes(startLocation int) ([]int, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	replicationFactor := config.ConfigNode.ReplicationFactor

	// Check if physical node number is at least equal to replication factor
	log.Println(n.HeartBeatTable)
	aliveNodeCount := 0
	for _, alive := range n.HeartBeatTable {
		if alive {
			aliveNodeCount++
		}
	}
	if aliveNodeCount <= replicationFactor {
		return nil, fmt.Errorf("number of nodes alive [%d] cannot be smaller or equal to the replication factor [%d]. Please consider either creating more nodes or modifying the replication factor", aliveNodeCount, replicationFactor)
	}
	relayVirtualNodes := make([]int, replicationFactor)
	pickedPhysicalNodes := make([]int, replicationFactor)
	selector := 0
	idx := n.getVirtualLocationIndex(startLocation) + 1
	pickedPhysicalNodes = append(pickedPhysicalNodes, n.Pid) // dont choose myself as next relay
	for {
		if idx == len(n.VirtualNodeLocation) {
			idx = 0
		}
		loc := n.VirtualNodeLocation[idx]
		if selector >= replicationFactor {
			break
		}
		physicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[loc])
		if err != nil {
			return nil, err
		}
		if !util.IntInSlice(pickedPhysicalNodes, physicalNodeID) {
			relayVirtualNodes[selector] = loc
			pickedPhysicalNodes[selector] = physicalNodeID
			selector++
		}
		idx++
	}
	return relayVirtualNodes, nil
}
