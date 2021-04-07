package user

import (
	"errors"

	"github.com/xmliszt/e-safe/pkg/file"
	"github.com/xmliszt/e-safe/util"
)

// User contains all the variables that are found in User database
type User struct {
	Username string `json:"Username"` // Username as string
	Password string `json:"Password"` // Password as string
	Role     int    `json:"Role"`     // Int to identify role (clearance level)
}

// UserMethods contains all method references for interacting with the User database
type UserMethods interface {
	GetUser(username string) User
	CreateUser(user User)
	GetUsers() []User
}

func encodeUser(data interface{}) (*User, error) {
	user := &User{}
	for key, val := range data.(map[string]interface{}) {
		if key == "Role" {
			val = int(val.(float64))
		}
		err := util.SetField(user, key, val)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func decodeUser(user *User) interface{} {
	var r interface{} = user
	return r
}

// GetUser obtains a specific user based on their username provided
func GetUser(username string) (*User, error) {
	allUsers, fileError := file.ReadUsersFile()
	if fileError != nil {
		return nil, fileError
	} else {
		for _, val := range allUsers {
			user, err := encodeUser(val)
			if err != nil {
				return nil, err
			}
			if user.Username == username {
				return user, nil
			}

		}
		unknownUsernameError := errors.New("username not available")
		return nil, unknownUsernameError
	}
}

// CreateUser creates a user based on the User structure provided
func CreateUser(user *User, userID string) error {
	newUserVal := decodeUser(user)
	newUser := map[string]interface{}{userID: newUserVal}
	fileError := file.WriteUsersFile(newUser)
	if fileError != nil {
		return nil
	} else {
		return fileError
	}
}

// GetUsers gets all  users
func GetUsers() (map[string]*User, error) {
	allUsersVal, fileError := file.ReadUsersFile()
	allUsers := make(map[string]*User)
	for key, val := range allUsersVal {
		user, err := encodeUser(val)
		if err != nil {
			return nil, err
		}
		allUsers[key] = user
	}
	if fileError != nil {
		return nil, fileError
	} else {
		return allUsers, nil
	}
}
