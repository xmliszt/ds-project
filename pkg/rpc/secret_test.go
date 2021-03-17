package rpc

// Testing GetSecret
// recv0Channel := make(chan *rpc.Data)
// send0Channel := make(chan *rpc.Data)
// rpcMap := make(map[int]chan *rpc.Data)
// boolean := true

// var myNode = rpc.Node{&boolean, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// // var testUserDetails1 = rpc.User{"Sudipta", "iLoveDistSys", 100}
// singleSecret, fileError := myNode.GetSecret("1006969")
// if fileError != nil {
// 	fmt.Println(fileError)
// } else {
// 	fmt.Println(singleSecret)
// 	os.Exit(1)

// }

// Testing GetSecrets
// recv0Channel := make(chan *rpc.Data)
// send0Channel := make(chan *rpc.Data)
// rpcMap := make(map[int]chan *rpc.Data)
// boolean := true

// var myNode = rpc.Node{&boolean, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// // var testUserDetails1 = rpc.User{"Sudipta", "iLoveDistSys", 100}
// groupSecret, fileError := myNode.GetSecrets(127, 129)
// if fileError != nil {
// 	fmt.Println(fileError)
// } else {
// 	fmt.Println(groupSecret)
// 	os.Exit(1)

// }

// Testing PutSecret
// recv0Channel := make(chan *rpc.Data)
// send0Channel := make(chan *rpc.Data)
// rpcMap := make(map[int]chan *rpc.Data)
// boolean := true

// var myNode = rpc.Node{&boolean, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// // var testUserDetails1 = rpc.User{"Sudipta", "iLoveDistSys", 100}
// var testSecretDetails1 = rpc.Secret{"mySectret", 100}
// fileError := myNode.PutSecret("10004568", testSecretDetails1)
// if fileError != nil {
// 	fmt.Println(fileError)
// } else {
// 	// fmt.Println(singleSecret)
// 	fmt.Println("Putted Secret")
// 	os.Exit(1)

// }

// Testing Delete Secret
// recv0Channel := make(chan *rpc.Data)
// send0Channel := make(chan *rpc.Data)
// rpcMap := make(map[int]chan *rpc.Data)
// boolean := true

// var myNode = rpc.Node{&boolean, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// // var testUserDetails1 = rpc.User{"Sudipta", "iLoveDistSys", 100}
// // var testSecretDetails1 = rpc.Secret{"mySectret", 100}
// fileError := myNode.DeleteSecret("10004567")
// if fileError != nil {
// 	fmt.Println(fileError)
// } else {
// 	// fmt.Println(singleSecret)
// 	fmt.Println("Deleted Secret")
// 	os.Exit(1)

// }
