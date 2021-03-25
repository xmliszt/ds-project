package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/pkg/secret"
	"github.com/xmliszt/e-safe/pkg/user"
)

// Put a secret
func PutSecret(ctx echo.Context) error {
	secret := new(secret.Secret)
	if err := ctx.Bind(secret); err != nil {
		return ctx.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	// Handle 3 replications and store data
	return ctx.String(http.StatusOK, fmt.Sprintf("Putting secret: %+v...", secret))
}

// Get a secret
func GetSecret(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	role, _ := strconv.Atoi(claims["role"].(string))

	alias := ctx.QueryParam("alias")
	if len(alias) < 1 {
		return ctx.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Error:   "Unknown URL params. 'alias' is not defined!",
			Data:    nil,
		})
	}
	// Handle getting a secret
	// Test a sample secret
	secret, err := secret.GetSecret(1, "126")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	if role > secret.Role {
		return ctx.JSON(http.StatusUnauthorized, &Response{
			Success: false,
			Error:   "Your role is too low for this information!",
			Data: &user.User{
				Username: claims["username"].(string),
				Role:     role,
			},
		})
	}
	return ctx.JSON(http.StatusOK, &Response{
		Success: true,
		Data: []interface{}{
			secret,
		},
	})
}

// Delete a secret
func DeleteSecret(ctx echo.Context) error {
	alias := ctx.QueryParam("alias")
	if len(alias) < 1 {
		return ctx.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Error:   "Unknown URL params. 'alias' is not defined!",
			Data:    nil,
		})
	}
	// Handle delete a secret
	err := secret.DeleteSecret(1, "126")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	return ctx.JSON(http.StatusOK, &Response{
		Success: true,
	})
}

// Get all secrets under a role
func GetAllSecrets(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	role, _ := strconv.Atoi(claims["role"].(string))

	fmt.Println("User role is: ", role)

	// Handle get all secrets
	return ctx.JSON(http.StatusOK, &Response{
		Success: true,
		Error:   "",
		Data: []*secret.Secret{
			{
				Role:  2,
				Value: "Sample secret",
				Alias: "It is a sample secret",
			},
		},
	})
}
