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

// // GetUser obtains a specific user based on their username provided
// func (n Node) GetUser(username string) User {
// 	return User{"User", "Pass", 2}
// }

// // CreateUser creates a user based on the User structure provided
// func (n Node) CreateUser(user User) {
// 	fmt.Println("Hellu")
// }

// // SimpleMethod for testing
// func SimpleMethod() {
// 	fmt.Println("Simple Method")
// }

// // GetUsers gets all  users
// func (n Node) GetUsers() []User {
// 	// user1 = User{"user1", "password1", 1}
// 	// user2 = User{"user2", "password2", 1}
// 	var userList = []User{User{"user2", "password2", 1}, User{"user2", "password2", 1}}
// 	return userList
// }
