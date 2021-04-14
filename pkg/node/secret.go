package node

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/api"
	"github.com/xmliszt/e-safe/pkg/message"
	"github.com/xmliszt/e-safe/pkg/secret"
	"github.com/xmliszt/e-safe/pkg/user"
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
	fmt.Println("This is the alias for the secret", recievingSecret.Alias)
	hashedAlias, err := util.GetHash(recievingSecret.Alias)
	fmt.Println("This is the hashedAlias", hashedAlias)
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
	fmt.Println("This is the VNodeLoc", vNodeLoc)

	virtualNodesList, err := n.getRelayVirtualNodes(vNodeLoc)
	fmt.Println("This is the VirtualNodeList", virtualNodesList)

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
	fmt.Println("this is the nextPhyscialNodeID", nextPhysicalNodeID)

	nextPhysicalNodeRpc := n.RpcMap[nextPhysicalNodeID]
	fmt.Println("this is the nextPhysicalNodeRpc", nextPhysicalNodeRpc)

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
		log.Println(err)
		log.Println("Error sending message to owner node. It is dead")
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
			log.Fatal("Node next to owner is dead. Problemo, should never happen")
		}

	}

	payload := reply.Payload.(map[string]interface{})
	if payload["success"].(bool) {
		return ctx.JSON(http.StatusInternalServerError, &api.Response{
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

// Get a secret - deprecated
func (n *Node) getSecret(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	role, _ := strconv.Atoi(claims["role"].(string))

	alias := ctx.QueryParam("alias")
	if len(alias) < 1 {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   "Unknown URL params. 'alias' is not defined!",
			Data:    nil,
		})
	}
	// Handle getting a secret

	// Need to decide where is the secret
	hashedAlias, err := util.GetHash(alias)
	if err != nil {
		log.Println("Hashing error")
		return err
	}

	ownerVNodeLoc := util.MapHashToVNodeLoc(n.VirtualNodeMap, n.VirtualNodeLocation, hashedAlias)
	ownerPhysicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[ownerVNodeLoc])
	if err != nil {
		// log.Fatal("Error when geting the list of virtual nodes for replication")
		return ctx.JSON(http.StatusInternalServerError, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}

	ownerNodeAddress := n.RpcMap[ownerPhysicalNodeID]
	ownerRequest := &message.Request{
		From: n.Pid,
		To:   ownerPhysicalNodeID,
		Code: message.GIVE_ME_DATA,
		Payload: map[string]interface{}{
			"key": int(hashedAlias),
		},
	}

	var reply message.Reply
	// if(n.checkHeartbeat(nextPhysicalNodeID)){
	err = message.SendMessage(ownerNodeAddress, "Node.GetData", ownerRequest, &reply)
	if err != nil {
		virtualNodesList, err := n.getRelayVirtualNodes(ownerVNodeLoc)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &api.Response{
				Success: false,
				Error:   err.Error(),
				Data:    nil,
			})
		}

		nextPhysicalNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[virtualNodesList[0]])
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &api.Response{
				Success: false,
				Error:   err.Error(),
				Data:    nil,
			})
		}
		ownerNodeAddress := n.RpcMap[nextPhysicalNodeID]

		nextNodeRequest := &message.Request{
			From: n.Pid,
			To:   nextPhysicalNodeID,
			Code: message.GIVE_ME_DATA,
			Payload: map[string]interface{}{
				"key": int(hashedAlias),
			},
		}

		err = message.SendMessage(ownerNodeAddress, "Node.GetData", nextNodeRequest, &reply)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &api.Response{
				Success: false,
				Error:   err.Error(),
				Data:    nil,
			})
		}

	}

	payload := reply.Payload.(map[string]interface{})
	retrievedSecret := payload["secret"].(secret.Secret)

	if role > retrievedSecret.Role {
		return ctx.JSON(http.StatusUnauthorized, &api.Response{
			Success: false,
			Error:   "Your role is too low for this information!",
			Data: &user.User{
				Username: claims["username"].(string),
				Role:     role,
			},
		})
	}

	return ctx.JSON(http.StatusOK, &api.Response{
		Success: true,
		Data: []interface{}{
			retrievedSecret,
		},
	})

	// Ask the respective node to get the secret
	// Reply based on the reply from the respective node
	// Test a sample secret
	// secret, err := secret.GetSecret(1, "126")
	// if err != nil {
	// 	return ctx.JSON(http.StatusBadRequest, &api.Response{
	// 		Success: false,
	// 		Error:   err.Error(),
	// 		Data:    nil,
	// 	})
	// }
	// if role > secret.Role {
	// 	return ctx.JSON(http.StatusUnauthorized, &api.Response{
	// 		Success: false,
	// 		Error:   "Your role is too low for this information!",
	// 		Data: &user.User{
	// 			Username: claims["username"].(string),
	// 			Role:     role,
	// 		},
	// 	})
	// }
	// return ctx.JSON(http.StatusOK, &api.Response{
	// 	Success: true,
	// 	Data: []interface{}{
	// 		secret,
	// 	},
	// })
}

// Delete a secret
// func (n *Node) deleteSecret(ctx echo.Context) error {
// 	alias := ctx.QueryParam("alias")
// 	if len(alias) < 1 {
// 		return ctx.JSON(http.StatusBadRequest, &api.Response{
// 			Success: false,
// 			Error:   "Unknown URL params. 'alias' is not defined!",
// 			Data:    nil,
// 		})
// 	}

// 	// Handle delete a secret
// 	uSecretHash, err := util.GetHash(alias)
// 	if err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, &api.Response{
// 			Success: false,
// 			Error:   err.Error(),
// 			Data:    nil,
// 		})
// 	}
// 	secretHash := int(uSecretHash)

// 	var targetLocation int
// 	for _, loc := range n.VirtualNodeLocation {
// 		if loc > secretHash {
// 			targetLocation = loc
// 			break
// 		}
// 	}

// 	targetNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[targetLocation])
// 	if err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, &api.Response{
// 			Success: false,
// 			Error:   err.Error(),
// 			Data:    nil,
// 		})
// 	}

// 	targetNodeAddr := n.RpcMap[targetNodeID]

// 	request := &message.Request{
// 		From: n.Pid,
// 		To:   targetNodeID,
// 		Code: message.DELETE_SECRET,
// 		Payload: map[string]interface{}{
// 			"key":      secretHash,
// 			"location": targetLocation,
// 		},
// 	}

// 	var reply message.Reply
// 	err = message.SendMessage(targetNodeAddr, "Node.DeleteSecret", request, &reply)
// 	if err != nil {
// 		// When owner node is down, we still need to try delete all replicas
// 		relayVirtualNodes, err := n.getRelayVirtualNodes(targetLocation)
// 		if err != nil {
// 			return ctx.JSON(http.StatusInternalServerError, &api.Response{
// 				Success: false,
// 				Error:   err.Error(),
// 				Data:    nil,
// 			})
// 		}
// 		config, err := config.GetConfig()
// 		if err != nil {
// 			return ctx.JSON(http.StatusInternalServerError, &api.Response{
// 				Success: false,
// 				Error:   err.Error(),
// 				Data:    nil,
// 			})
// 		}
// 		err = n.relaySecretDeletion(config.ConfigNode.ReplicationFactor, strconv.Itoa(secretHash), relayVirtualNodes)
// 		if err != nil {
// 			return ctx.JSON(http.StatusInternalServerError, &api.Response{
// 				Success: false,
// 				Error:   err.Error(),
// 				Data:    nil,
// 			})
// 		}
// 		return ctx.JSON(http.StatusOK, &api.Response{
// 			Success: true,
// 		})
// 	} else {
// 		replyPayload := reply.Payload.(map[string]interface{})
// 		success := replyPayload["success"].(bool)
// 		if success {
// 			return ctx.JSON(http.StatusOK, &api.Response{
// 				Success: true,
// 			})
// 		} else {
// 			err := replyPayload["error"].(error)
// 			return ctx.JSON(http.StatusInternalServerError, &api.Response{
// 				Success: false,
// 				Error:   err.Error(),
// 				Data:    nil,
// 			})
// 		}
// 	}
// }

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
				return ctx.JSON(http.StatusInternalServerError, &api.Response{
					Success: false,
					Error:   "One of the node is down ",
				})
			}
		}

	}
	return ctx.JSON(http.StatusOK, &api.Response{
		Success: true,
	})
	// return nil
}

// Get all secrets under a role
func (n *Node) getAllSecrets(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	role, _ := strconv.Atoi(claims["role"].(string))

	fmt.Println("User role is: ", role)

	// Handle get all secrets
	return ctx.JSON(http.StatusOK, &api.Response{
		Success: true,
		Error:   "",
		Data: map[string]interface{}{
			"role": role,
			"data": []*secret.Secret{
				{
					Role:  2,
					Value: "Sample secret",
					Alias: "It is a sample secret",
				},
			},
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
