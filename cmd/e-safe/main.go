package main

import (
	"fmt"
	"os"

	"github.com/xmliszt/e-safe/pkg/rpc"
)

func main() {
	// err := locksmith.Start()
	recv0Channel := make(chan *rpc.Data)
	send0Channel := make(chan *rpc.Data)
	rpcMap := make(map[int]chan *rpc.Data)
	boolean := true

	var myNode = rpc.Node{&boolean, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

	// var testUserDetails1 = rpc.User{"Sudipta", "iLoveDistSys", 100}
	// var testSecretDetails1 = rpc.Secret{"mySectret", 100}
	fileError := myNode.DeleteSecret("10004567")
	if fileError != nil {
		fmt.Println(fileError)
	} else {
		// fmt.Println(singleSecret)
		fmt.Println("Deleted Secret")
		os.Exit(1)

	}
}
