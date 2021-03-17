package rpc

// Testing GetUser
// 	send0Channel := make(chan *rpc.Data)
// 	rpcMap := make(map[int]chan *rpc.Data)
// 	boolean := true

// 	var myNode = rpc.Node{&boolean, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// 	// var testUserDetails1 = rpc.User{"Sudipta", "iLoveDistSys", 100}
// 	allUser, fileError := myNode.GetUser("DevBahl")
// 	if fileError != nil {
// 		fmt.Println(fileError)
// 	} else {
// 		fmt.Println(allUser)
// 		os.Exit(1)

// 	}

// Testing Create User
// recv0Channel := make(chan *rpc.Data)
// send0Channel := make(chan *rpc.Data)
// rpcMap := make(map[int]chan *rpc.Data)
// boolean := true

// var myNode = rpc.Node{&boolean, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// var testUserDetails1 = rpc.User{"Sudipta", "iLoveDistSys", 100}
// fileError := myNode.CreateUser(testUserDetails1, "1002256")
// if fileError != nil {
// 	fmt.Println(fileError)
// } else {
// 	fmt.Println("Created User")
// 	os.Exit(1)

// }

//Testing GetUsers
// recv0Channel := make(chan *rpc.Data)
// send0Channel := make(chan *rpc.Data)
// rpcMap := make(map[int]chan *rpc.Data)
// boolean := true

// var myNode = rpc.Node{&boolean, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// // var testUserDetails1 = rpc.User{"Sudipta", "iLoveDistSys", 100}
// allUsers, fileError := myNode.GetUsers()
// if fileError != nil {
// 	fmt.Println(fileError)
// } else {
// 	fmt.Println(allUsers)
// 	os.Exit(1)

// }
