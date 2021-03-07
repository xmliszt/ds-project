package main

import (
	"github.com/xmliszt/e-safe/pkg/rpc"
)

func main() {

	// recv0Channel := make(chan map[string]interface{})
	// send0Channel := make(chan map[string]interface{})
	// send1Channel := make(chan map[string]interface{})
	// send2Channel := make(chan map[string]interface{})
	// rpcMap := make(map[int]chan map[string]interface{})
	// rpcMap[1] = send1Channel
	// rpcMap[2] = send2Channel
	// var myNode = rpc.Node{true, 1, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}
	// fmt.Println("Hellow World!")
	// user := myNode.getUser2("User")
	// fmt.Println(myNode.GetUser("User"))

	// rpc.SimpleMethod()

	rpc.ReadJSON("users.json")
}
