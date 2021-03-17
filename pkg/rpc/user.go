package rpc

import "errors"

// User contains all the variables that are found in User database
type User struct {
	Username string // Username as string
	Password string // Password as string
	Role     int    // Int to identify role (clearance level)
}

// UserMethods contains all method references for interacting with the User database
type UserMethods interface {
	GetUser(username string) User
	CreateUser(user User)
	GetUsers() []User
}

// GetUser obtains a specific user based on their username provided
func (n *Node) GetUser(username string) (interface{}, error) {
	allUsers, fileError := n.ReadUsersFile()
	if fileError != nil {
		return nil, fileError
	} else {
		for key := range allUsers {
			if allUsers[key].Username == username {
				return allUsers[key], nil
			}

		}
		unknownUsernameError := errors.New("use	rname not available")
		return nil, unknownUsernameError
	}
}

// CreateUser creates a user based on the User structure provided
func (n *Node) CreateUser(user User, userID string) error {
	newUser := UserToMapUser(user, userID)
	fileError := n.WriteUsersFile(newUser)
	if fileError != nil {
		return nil
	} else {
		return fileError
	}
}

func UserToMapUser(user User, userID string) map[string]User {
	newUser := make(map[string]User)
	newUser[userID] = user
	return newUser
}

// GetUsers gets all  users
func (n *Node) GetUsers() (map[string]User, error) {
	allUsers, fileError := n.ReadUsersFile()
	if fileError != nil {
		return nil, fileError
	} else {
		return allUsers, nil
	}
}
