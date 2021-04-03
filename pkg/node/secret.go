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
	secret := new(secret.Secret)
	if err := ctx.Bind(secret); err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	// Handle 3 replications and store data
	return ctx.String(http.StatusOK, fmt.Sprintf("Putting secret: %+v...", secret))
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
	// Test a sample secret
	secret, err := secret.GetSecret(1, "126")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	if role > secret.Role {
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
			secret,
		},
	})
}

// Delete a secret
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
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	secretHash := int(uSecretHash)

	var targetLocation int
	for _, loc := range n.VirtualNodeLocation {
		if loc > secretHash {
			targetLocation = loc
			break
		}
	}

	targetNodeID, err := getPhysicalNodeID(n.VirtualNodeMap[targetLocation])
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}

	targetNodeAddr := n.RpcMap[targetNodeID]

	request := &message.Request{
		From: n.Pid,
		To:   targetNodeID,
		Code: message.DELETE_SECRET,
		Payload: map[string]interface{}{
			"key":      secretHash,
			"location": targetLocation,
		},
	}

	var reply message.Reply
	err = message.SendMessage(targetNodeAddr, "Node.DeleteSecret", request, &reply)
	if err != nil {
		// When owner node is down, we still need to try delete all replicas
		relayVirtualNodes, err := n.getRelayVirtualNodes(targetLocation)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &api.Response{
				Success: false,
				Error:   err.Error(),
				Data:    nil,
			})
		}
		config, err := config.GetConfig()
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &api.Response{
				Success: false,
				Error:   err.Error(),
				Data:    nil,
			})
		}
		err = n.relaySecretDeletion(config.ConfigNode.ReplicationFactor, strconv.Itoa(secretHash), relayVirtualNodes)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &api.Response{
				Success: false,
				Error:   err.Error(),
				Data:    nil,
			})
		}
		return ctx.JSON(http.StatusOK, &api.Response{
			Success: true,
		})
	} else {
		replyPayload := reply.Payload.(map[string]interface{})
		success := replyPayload["success"].(bool)
		if success {
			return ctx.JSON(http.StatusOK, &api.Response{
				Success: true,
			})
		} else {
			err := replyPayload["error"].(error)
			return ctx.JSON(http.StatusBadRequest, &api.Response{
				Success: false,
				Error:   err.Error(),
				Data:    nil,
			})
		}
	}
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
		log.Printf("Node %d relay secret deletion error: %s\n", n.Pid, err)
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
		return err
	}
	return nil
}
