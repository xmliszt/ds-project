package rpc

import "fmt"

// Node contains all the variables that are necessary to manage a node
type Node struct {
	IsCoordinator bool
	Pid           int                                 // Node ID
	Ring          []int                               // Ring structure of nodes
	RecvChannel   chan map[string]interface{}         // Receiving channel
	SendChannel   chan map[string]interface{}         // Sending channel
	RPCMap        map[int]chan map[string]interface{} // Map node ID to their receiving channels
}

// HandleMessageReceived run as Go Routine to handle the messages received
func (n *Node) HandleMessageReceived() {
	for msg := range n.RecvChannel {
		fmt.Println("I receive: ", msg)
	}
}

// GetUser obtains a specific user based on their username provided
func (n Node) GetUser(username string) User {
	return User{"User", "Pass", 2}
}

// CreateUser creates a user based on the User structure provided
func (n Node) CreateUser(user User) {
	fmt.Println("Hellu")
}

// SimpleMethod for testing
func SimpleMethod() {
	fmt.Println("Simple Method")
}

// GetUsers gets all  users
func (n Node) GetUsers() []User {
	// user1 = User{"user1", "password1", 1}
	// user2 = User{"user2", "password2", 1}
	var userList = []User{User{"user2", "password2", 1}, User{"user2", "password2", 1}}
	return userList
}
