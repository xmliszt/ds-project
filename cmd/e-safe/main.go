package main

import (
<<<<<<< HEAD
	"github.com/xmliszt/e-safe/pkg/rpc"
)

func main() {

	recv0Channel := make(chan map[string]interface{})
	send0Channel := make(chan map[string]interface{})
	// send1Channel := make(chan map[string]interface{})
	// send2Channel := make(chan map[string]interface{})
	rpcMap := make(map[int]chan map[string]interface{})
	// rpcMap[1] = send1Channel
	// rpcMap[2] = send2Channel
	var myNode = rpc.Node{true, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}
	// fmt.Println("Hellow World!")
	// user := myNode.getUser2("User")
	// fmt.Println(myNode.GetUser("User"))

	// var testUserDetails User
	// var testUserDetails1 = rpc.User{"Sudipta", "iLoveDistSys", 100}
	// var testUserDetails2 = rpc.User{"Juan", "iLoveDistSys", 1}
	// // var testUserInput map[string]rpc.User
	// testUserInput := make(map[string]rpc.User, 5)
	// testUserInput["1006969"] = testUserDetails1
	// testUserInput["1007070"] = testUserDetails2

	var testSecretDetails1 = rpc.Secret{"mySectret", 100}
	var testSecretDetails2 = rpc.Secret{"yourSecret", 1}
	// var testUserInput map[string]rpc.User
	testDataInput := make(map[string]rpc.Secret, 5)
	testDataInput["1006969"] = testSecretDetails1
	testDataInput["1007070"] = testSecretDetails2

	// rpc.SimpleMethod()

	// myNode.WriteUsersFile(testUserInput)
	myNode.WriteDataFile(testDataInput)
	// myNode.ReadDataFile()
	// rpc.ReadJSONOri("user2.json")
}
=======
	"fmt"
	"os"

	"github.com/xmliszt/e-safe/pkg/locksmith"
)

func main() {
	err := locksmith.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
>>>>>>> dev
