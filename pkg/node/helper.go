package node

import (
	"fmt"
	"math"
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

	ownerNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[startLocation])
	if err != nil {
		return nil, err
	}
	pickedPhysicalNodes = append(pickedPhysicalNodes, ownerNodeID) // dont choose myself as next relay
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

// getReplicationLocations gets the locations of the virtual nodes which
// are the owners of the secrets whose replicas are supposed to store in
// this new born node.
func (n *Node) getReplicationLocations(location int) ([][]int, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	res := make([][]int, 0)
	nodesVisited := make([]int, 0)
	virtualNodeCount := 0
	var nextPhysicalNodeID int

	ownLocationIdx := n.getVirtualLocationIndex(location)

	// If my next and prev are same node, then i do not need to get replicas
	// As all the replicas will be store at the first virtual node of the
	// physical node that is different from its previous node
	nextIdx := ownLocationIdx + 1
	if nextIdx == len(n.VirtualNodeLocation) {
		nextIdx = 0
	}
	prevIdx := ownLocationIdx - 1
	if prevIdx == -1 {
		prevIdx = len(n.VirtualNodeLocation) - 1
	}
	nextNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[n.VirtualNodeLocation[nextIdx]])
	if err != nil {
		return nil, err
	}
	prevNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[n.VirtualNodeLocation[prevIdx]])
	if err != nil {
		return nil, err
	}
	if nextNodeID == prevNodeID {
		return [][]int{}, nil
	}

	// Get the next physical node ID that is not myself
	idx := ownLocationIdx + 1
	for {
		if idx == len(n.VirtualNodeLocation) {
			idx = 0
		}
		loc := n.VirtualNodeLocation[idx]
		physicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[loc])
		if err != nil {
			return nil, err
		}
		if physicalNodeID == n.Pid {
			idx++
		} else {
			nextPhysicalNodeID = physicalNodeID
			break
		}
	}

	// Get all previous virtual nodes
	idx = ownLocationIdx - 1
	for {
		if len(nodesVisited) >= config.ConfigNode.ReplicationFactor {
			break
		}
		if idx == -1 {
			idx = len(n.VirtualNodeLocation) - 1
		}
		var virtualNodeLoc int
		var prevVirtualNodeLoc int

		virtualNodeLoc = n.VirtualNodeLocation[idx]
		if idx-1 < 0 {
			prevVirtualNodeLoc = n.VirtualNodeLocation[len(n.VirtualNodeLocation)-int(math.Abs(float64(idx-1)))]
		} else {
			prevVirtualNodeLoc = n.VirtualNodeLocation[idx-1]
		}

		physicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[virtualNodeLoc])
		if err != nil {
			return nil, err
		}
		if !util.IntInSlice(nodesVisited, physicalNodeID) && physicalNodeID != n.Pid && physicalNodeID != nextPhysicalNodeID {
			if virtualNodeCount < config.ConfigNode.VirtualNodesCount-1 {
				virtualNodeCount++
				res = append(res, []int{physicalNodeID, prevVirtualNodeLoc, virtualNodeLoc})
			} else {
				nodesVisited = append(nodesVisited, physicalNodeID)
				res = append(res, []int{physicalNodeID, prevVirtualNodeLoc, virtualNodeLoc})
				virtualNodeCount = 0
			}
		}
		idx--
	}

	return res, nil
}

// Find physical node int from virtual node string
func GetPhysicalNode(vn string) (int, error) {
	var sPhysicalNode string

	for _, char := range vn {
		if string(char) != "-" {
			sPhysicalNode = sPhysicalNode + string(char)
		} else {
			break
		}
	}
	PhysicalNode, err := strconv.Atoi(sPhysicalNode)

	if err != nil {
		return -1, err
	}

	return PhysicalNode, nil
}
