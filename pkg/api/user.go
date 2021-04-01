package api

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/pkg/user"
)

// User log in - return JWT token for authentication
func LogIn(ctx echo.Context) error {
	u := new(user.User)
	if err := ctx.Bind(u); err != nil {
		return ctx.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}

	users, userErr := user.GetUsers()
	if userErr != nil {
		return userErr
	}

	for _, user := range users {
		if subtle.ConstantTimeCompare([]byte(user.Username), []byte(u.Username)) == 1 {
			if subtle.ConstantTimeCompare([]byte(user.Password), []byte(u.Password)) == 1 {
				token := jwt.New(jwt.SigningMethodHS256)
				claims := token.Claims.(jwt.MapClaims)
				claims["username"] = user.Username
				claims["role"] = fmt.Sprintf("%d", user.Role)
				claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

				t, err := token.SignedString([]byte("secret"))
				if err != nil {
					return ctx.JSON(http.StatusBadRequest, &Response{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					})
				}

				return ctx.JSON(http.StatusOK, &Response{
					Success: true,
					Error:   "",
					Data:    t,
				})
			}
		}
	}

	return ctx.JSON(http.StatusUnauthorized, &Response{
		Success: false,
		Error:   "Unauthorised",
	})
}

// Create a user - Sign up
func Register(ctx echo.Context) error {
	user := new(user.User)
	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	// Handle user registration
	return ctx.String(http.StatusOK, fmt.Sprintf("Register a new user: Username: %s, Password: %s, Role: %d", user.Username, user.Password, user.Role))
}
