package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/xmliszt/e-safe/pkg/locksmith"
	"github.com/xmliszt/e-safe/pkg/secret"
	"github.com/xmliszt/e-safe/pkg/user"
)

func main() {

	var mode int = 1
	flag.IntVar(&mode, "m", 1, "Select mode for different demo")
	flag.Parse()

	// Demo locksmith checks heartbeat
	if mode == 1 {
		err := locksmith.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if mode == 2 {
		nodeID := rand.Intn(3) + 1

		for {
			fmt.Printf("Node %d is pleased to serve you :)\n", nodeID)
			fmt.Print("Features to select:\n\t1. Get all users\n\t2. Create a new user\n\t3. Get a user\n\t4. Get a range of secrets\n\t5. Get a secret\n\t6. Put a secret\n\t7. Delete a secret\nEnter your option: ")
			var option int = 1
			fmt.Scanln(&option)

			// Read all users
			switch option {
			case 1:
				users, err := user.GetUsers()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(users)
			case 2:
				var uid string
				var username string
				var password string
				fmt.Print("Enter user ID: ")
				fmt.Scanln(&uid)
				fmt.Print("Enter username: ")
				fmt.Scanln(&username)
				fmt.Print("Enter password: ")
				fmt.Scanln(&password)
				newUser := user.User{
					Username: username,
					Password: password,
					Role:     rand.Intn(5) + 1,
				}
				err := user.CreateUser(newUser, uid)
				if err != nil {
					fmt.Println(err)
				}
			case 3:
				var username string
				fmt.Print("Enter username to search: ")
				fmt.Scanln(&username)
				user, err := user.GetUser(username)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(user)
			case 4:
				secrets, err := secret.GetSecrets(nodeID, 127, 130)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(secrets)
			case 5:
				var secretID string
				fmt.Print("Enter secret ID: ")
				fmt.Scanln(&secretID)
				secret, err := secret.GetSecret(nodeID, secretID)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(secret)
			case 6:
				var password string
				var role int
				fmt.Print("Enter the secret to store: ")
				fmt.Scanln(&password)
				fmt.Print("Enter your role associated with the secret: ")
				fmt.Scanln(&role)
				secretData := secret.Secret{
					Value: password,
					Role:  role,
				}
				err := secret.PutSecret(nodeID, "131", secretData)
				if err != nil {
					fmt.Println(err)
				}
			case 7:
				var secretID string
				fmt.Print("Enter the secret ID to delete: ")
				fmt.Scanln(&secretID)
				err := secret.DeleteSecret(nodeID, secretID)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("Refer to the JSON file in Node %d for update!\n", nodeID)
			}
		}
	}
}
