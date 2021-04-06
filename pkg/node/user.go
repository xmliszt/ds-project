package node

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/api"
	"github.com/xmliszt/e-safe/pkg/message"
	"github.com/xmliszt/e-safe/pkg/user"
	"github.com/xmliszt/e-safe/util"
)

// User log in - return JWT token for authentication
func (n *Node) logIn(ctx echo.Context) error {
	u := new(user.User)
	if err := ctx.Bind(u); err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
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
					return ctx.JSON(http.StatusBadRequest, &api.Response{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					})
				}

				return ctx.JSON(http.StatusOK, &api.Response{
					Success: true,
					Error:   "",
					Data:    t,
				})
			}
		}
	}

	return ctx.JSON(http.StatusUnauthorized, &api.Response{
		Success: false,
		Error:   "Unauthorised",
	})
}

// Create a user - Sign up
func (n *Node) register(ctx echo.Context) error {
	newUser := new(user.User)
	if err := ctx.Bind(newUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	// Handle user registration
	config, err := config.GetConfig()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	var key []byte = make([]byte, 32)
	var pwd []byte
	keyStr := config.ConfigServer.Secret
	copy(key, []byte(keyStr))
	pwdStr := newUser.Password
	pwd = []byte(pwdStr)
	// Encrypt Password
	cipher, err := util.Encrypt(key, pwd)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	newUser.Password = string(cipher)
	// Contact Locksmith for lock - Centralized Server Locking
	request := &message.Request{
		From:    n.Pid,
		To:      0,
		Code:    message.ACQUIRE_USER_LOCK,
		Payload: nil,
	}
	var reply message.Reply
	err = message.SendMessage(n.RpcMap[0], "LockSmith.AcquireUserLock", request, &reply)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	userIDHash, err := util.GetHash(newUser.Username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &api.Response{
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		})
	}
	userIDInt := int(userIDHash)
	userID := strconv.Itoa(userIDInt)
	err = user.CreateUser(newUser, userID)
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
