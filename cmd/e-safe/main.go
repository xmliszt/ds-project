package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/xmliszt/e-safe/pkg/locksmith"
	"github.com/xmliszt/e-safe/pkg/rpc"
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
		recv0Channel := make(chan *rpc.Data)
		send0Channel := make(chan *rpc.Data)
		rpcMap := make(map[int]chan *rpc.Data)
		boolean := true

		nodeID := rand.Intn(3) + 1
		var myNode = &rpc.Node{
			IsCoordinator:  &boolean,
			Pid:            nodeID,
			Ring:           []int{2, 3, 5, 7, 11, 13},
			RecvChannel:    recv0Channel,
			SendChannel:    send0Channel,
			RpcMap:         rpcMap,
			HeartBeatTable: make(map[int]bool),
		}

		for {
			fmt.Printf("Node %d is pleased to serve you :)\n", nodeID)
			fmt.Print("Features to select:\n\t1. Get all users\n\t2. Create a new user\n\t3. Get a user\n\t4. Get a range of secrets\n\t5. Get a secret\n\t6. Put a secret\n\t7. Delete a secret\nEnter your option: ")
			var option int = 1
			fmt.Scanln(&option)

			// Read all users
			switch option {
			case 1:
				users, err := myNode.ReadUsersFile()
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
				newUser := map[string]rpc.User{
					uid: {
						Username: username,
						Password: password,
						Role:     rand.Intn(5) + 1,
					}}
				err := myNode.WriteUsersFile(newUser)
				if err != nil {
					fmt.Println(err)
				}
			case 3:
				var username string
				fmt.Print("Enter username to search: ")
				fmt.Scanln(&username)
				user, err := myNode.GetUser(username)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(user)
			case 4:
				secret, err := myNode.GetSecrets(127, 130)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(secret)
			case 5:
				var secretID string
				fmt.Print("Enter secret ID: ")
				fmt.Scanln(&secretID)
				secret, err := myNode.GetSecret(secretID)
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
				secret := rpc.Secret{
					Value: password,
					Role:  role,
				}
				err := myNode.PutSecret("131", secret)
				if err != nil {
					fmt.Println(err)
				}
			case 7:
				var secretID string
				fmt.Print("Enter the secret ID to delete: ")
				fmt.Scanln(&secretID)
				err := myNode.DeleteSecret(secretID)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("Refer to the JSON file in Node %d for update!\n", nodeID)
			}
		}
	}
}
