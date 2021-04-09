package node

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/pkg/api"
)

type MonitorInfo struct {
	VirtualNodesLocation map[string]int
	VirtualNodesCapacity map[string]int
	HeartbeatTable       map[int]bool
}

// getMonitorInfo fetches information on virtual nodes and send to client for monitor GUI
func (n *Node) getMonitorInfo(ctx echo.Context) error {
	nodesCapacity := make(map[string]int)
	nodesLocation := make(map[string]int)
	for idx, loc := range n.VirtualNodeLocation {
		name := n.VirtualNodeMap[loc]
		nodesLocation[name] = loc
		if idx+1 < len(n.VirtualNodeLocation) {
			nodesCapacity[name] = n.VirtualNodeLocation[idx+1] - 1 - loc
		} else {
			nodesCapacity[name] = int(^uint32(0)) - loc + n.VirtualNodeLocation[0] - 1
		}
	}

	info := &MonitorInfo{
		VirtualNodesLocation: nodesLocation,
		HeartbeatTable:       n.HeartBeatTable,
		VirtualNodesCapacity: nodesCapacity,
	}
	return ctx.JSON(http.StatusOK, &api.Response{
		Success: true,
		Data:    info,
	})
}
