package rpc

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
	// allUsers := make(map[string]User)
	// var var1 interface{}
	// var var2 error
	allUsers, fileError := n.ReadUsersFile()
	// fmt.Println(allUsers, fileError)
	if fileError != nil {
		return nil, fileError
	} else {
		for key := range allUsers {
			//fmt.Println(allUsers[key].Username)
			if allUsers[key].Username == username {
				// var1, var2 = allUsers[key], nil
				return allUsers[key], nil
			}
		}
	}
	// map[string]User allUsers := n.ReadUsersFile();
	// return User{"User", "Pass", 2}
	return nil, nil
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
