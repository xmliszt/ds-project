package node

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/api"
	"github.com/xmliszt/e-safe/pkg/message"
	"github.com/xmliszt/e-safe/pkg/secret"
	"github.com/xmliszt/e-safe/util"
)

// Put a secret, if exists, update. If does not exist, create a new one
func (n *Node) putSecret(ctx echo.Context) error {
	recievingSecret := new(secret.Secret)
	if err := ctx.Bind(recievingSecret); err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	// Handle 3 replications and store data
	// Hash the alias secret, get a string
	hashedAlias, err := util.GetHash(recievingSecret.Alias)
	if err != nil {
		// log.Fatal("Error when hashing the alias")
		return ctx.JSON(http.StatusInternalServerError, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}

	// Get the relayVirtualNodes
	vNodeLoc := util.MapHashToVNodeLoc(n.VirtualNodeMap, n.VirtualNodeLocation, hashedAlias)

	virtualNodesList, err := n.getRelayVirtualNodes(vNodeLoc)

	if err != nil {
		// log.Fatal("Error when geting the list of virtual nodes for replication")
		return ctx.JSON(http.StatusInternalServerError, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}

	nextPhysicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[vNodeLoc])
	if err != nil {
		// log.Fatal("Error when geting the list of virtual nodes for replication")
		return ctx.JSON(http.StatusInternalServerError, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}

	nextPhysicalNodeRpc := n.RpcMap[nextPhysicalNodeID]

	// Construct request message
	ownerRequest := &message.Request{
		From: n.Pid,
		To:   nextPhysicalNodeID,
		Code: message.STORE_AND_REPLICATE,
		Payload: map[string]interface{}{
			"rf":     3,
			"key":    int(hashedAlias),
			"secret": recievingSecret,
			"nodes":  virtualNodesList,
		},
	}

	// Check if owner node is alive
	var reply message.Reply
	// if(n.checkHeartbeat(nextPhysicalNodeID)){
	err = message.SendMessage(nextPhysicalNodeRpc, "Node.StoreAndReplicate", ownerRequest, &reply)
	if err != nil {
		log.Printf("Error sending message to owner node: %s\n", err)
		log.Println("Sending strict node down to next node")
		vNodeNextToDeadOwner := virtualNodesList[0]
		// Construct request message
		request := &message.Request{
			From: n.Pid,
			To:   nextPhysicalNodeID,
			Code: message.STRICT_OWNER_DOWN,
			Payload: map[string]interface{}{
				"rf":     2,
				"key":    strconv.Itoa(int((hashedAlias))),
				"secret": recievingSecret,
				"nodes":  virtualNodesList,
			},
		}

		vNodeNameNextToOwner := n.VirtualNodeMap[vNodeNextToDeadOwner]
		nextnextPhysicalNodeID, err := getPhysicalNodeID(vNodeNameNextToOwner)
		if err != nil {
			// log.Fatal("Error when geting the list of virtual nodes for replication")
			return ctx.JSON(http.StatusInternalServerError, &api.Response{
				Success: false,
				Error:   err.Error(),
				Data:    nil,
			})
		}

		nextnextPhysicalNodeRpc := n.RpcMap[nextnextPhysicalNodeID]
		err = message.SendMessage(nextnextPhysicalNodeRpc, "Node.PerformStrictDown", request, &reply)
		if err != nil {
			log.Fatal(err)
		}

	}

	payload := reply.Payload.(map[string]interface{})
	if payload["success"].(bool) {
		return ctx.JSON(http.StatusOK, &api.Response{
			Success: true,
		})
	} else {
		return ctx.JSON(http.StatusInternalServerError, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
}

func (n *Node) deleteSecret(ctx echo.Context) error {
	alias := ctx.QueryParam("alias")
	if len(alias) < 1 {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   "Unknown URL params. 'alias' is not defined!",
			Data:    nil,
		})
	}

	// Handle delete a secret
	uSecretHash, err := util.GetHash(alias)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	secretHash := int(uSecretHash)
	var listODeadNodes []int

	for key, value := range n.RpcMap {
		// Check each physical node's heartbeat
		if n.checkHeartbeat(key) {

			// Send each physical node a message asking for all secrets within that role's scope
			request := &message.Request{
				From: n.Pid,
				To:   key,
				Code: message.DELETE_SECRET_ALL_INSTANCES,
				Payload: map[string]interface{}{
					"key": secretHash,
				},
			}

			var reply message.Reply
			err := message.SendMessage(value, "Node.DeleteSecret", request, &reply)
			if err != nil {
				// This tecnhnically should not happen since we are already checking if there is a heartbeat
				// But we add just in case
				listODeadNodes = append(listODeadNodes, key)
			}
		} else {
			// Check if the node that is checked is the locksmith, then don't add
			if key != 0 {
				listODeadNodes = append(listODeadNodes, key)
			}
		}

	}
	return ctx.JSON(http.StatusOK, &api.Response{
		Success: true,
		Data: map[string]interface{}{
			"deadNodes": listODeadNodes,
		},
	})
}

// Get all secrets under a role
func (n *Node) getAllSecrets(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	role, _ := strconv.Atoi(claims["role"].(string))

	var finalOutputList []*secret.Secret

	// Global variable listOSecrets := []*secret.Secret
	var listOSecrets []*secret.Secret
	var listODeadNodes []int
	// for all physical nodes in the ring
	for key, value := range n.RpcMap {
		// Check each physical node's heartbeat
		if n.checkHeartbeat(key) {

			// Send each physical node a message asking for all secrets within that role's scope
			request := &message.Request{
				From: n.Pid,
				To:   key,
				Code: message.GET_ALL_SECRETS,
				Payload: map[string]interface{}{
					"role": role,
				},
			}

			var reply message.Reply
			err := message.SendMessage(value, "Node.GetAllSecrets", request, &reply)
			if err != nil {
				// This tecnhnically should not happen since we are already checking if there is a heartbeat
				// But we add just in case
				listODeadNodes = append(listODeadNodes, key)
			} else {
				replyPayload := reply.Payload.(map[string]interface{})
				dataPayload := replyPayload["data"].([]*secret.Secret)

				// Append all the data in the  replies into the global variable listOSecrets
				listOSecrets = append(listOSecrets, dataPayload...)
			}
		} else {
			// Check if the node that is checked is the locksmith, then don't add
			if key != 0 {
				listODeadNodes = append(listODeadNodes, key)
			}
		}

	}

	// Check for dupllicates
	duplicateMap := make(map[string]secret.Secret)
	for _, secretValue := range listOSecrets {
		if _, ok := duplicateMap[secretValue.Alias]; ok {
			continue
		} else {
			duplicateMap[secretValue.Alias] = *secretValue
			finalOutputList = append(finalOutputList, secretValue)
		}
	}

	// Send listOSecrets to client
	return ctx.JSON(http.StatusOK, &api.Response{
		Success: true,
		Error:   "",
		Data: map[string]interface{}{
			"role":      role,
			"data":      finalOutputList,
			"deadNodes": listODeadNodes,
		},
	})
}

// relaySecretDeletion sends deletion signal to subsequent node replica that has the copy of the secret
func (n *Node) relaySecretDeletion(rf int, key string, relayNodes []int) error {
	config, err := config.GetConfig()
	if err != nil {
		log.Printf("Node %d is unable to relay secret deletion to next node: %s\n", n.Pid, err)
		return err
	}
	nextNodeLoc := relayNodes[config.ConfigNode.ReplicationFactor-rf]
	nextPhysicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[nextNodeLoc])
	if err != nil {
		log.Printf("Node %d is unable to relay secret deletion to next node: %s\n", n.Pid, err)
		return err
	}
	nextNodeAddr := n.RpcMap[nextPhysicalNodeID]
	request := &message.Request{
		From: n.Pid,
		To:   nextPhysicalNodeID,
		Code: message.RELAY_DELETE_SECRET,
		Payload: map[string]interface{}{
			"rf":    rf,
			"key":   key,
			"nodes": relayNodes,
		},
	}

	var reply message.Reply
	err = message.SendMessage(nextNodeAddr, "Node.RelayDeleteSecret", request, &reply)
	if err != nil {
		for rf > 1 {
			err := n.relaySecretDeletion(rf-1, key, relayNodes)
			if err != nil {
				if rf == 2 {
					return err
				}
				rf--
			} else {
				break
			}
		}
		return nil
	}
	return nil
}
