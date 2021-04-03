package node

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/pkg/api"
	"github.com/xmliszt/e-safe/pkg/secret"
	"github.com/xmliszt/e-safe/pkg/user"
)

// Put a secret
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
	err := secret.DeleteSecret(1, "131")
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
