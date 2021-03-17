package rpc

// import (
// 	"fmt"
// 	"reflect"
// 	"testing"
// )

// func TestReadDataFile(t *testing.T) {

// 	recv0Channel := make(chan map[string]interface{})
// 	send0Channel := make(chan map[string]interface{})
// 	rpcMap := make(map[int]chan map[string]interface{})
// 	var myNode = Node{true, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// 	result, error := myNode.ReadDataFile()
// 	if error != nil {
// 		fmt.Println(error)
// 		// return error
// 	}
// 	expected := map[string]Secret{"126": Secret{"hashed_Vict0r1aSecret3", 3}, "127": Secret{"hashed_Vict0r1aSecret3", 3},
// 		"128": Secret{"hashed_Vict0r1aSecret3", 3}, "129": Secret{"hashed_Vict0r1aSecret3", 3}, "130": Secret{"hashed_Vict0r1aSecret3", 3}}
// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("The ReadDataFile return \n%v \nwhich is not as same as what we expected \n%v", result, expected)
// 	}

// }

// func TestReadUsersFile(t *testing.T) {

// 	recv0Channel := make(chan map[string]interface{})
// 	send0Channel := make(chan map[string]interface{})
// 	rpcMap := make(map[int]chan map[string]interface{})
// 	var myNode = Node{true, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// 	result, error := myNode.ReadUsersFile()
// 	if error != nil {
// 		fmt.Println(error)
// 		// return error
// 	}
// 	expected := map[string]User{"1003651": User{"DevBahl", "Reader", 3}, "1003652": User{"Jose", "Reader", 2}, "1003653": User{"Aaron", "Reader", 4}, "1003654": User{"JiaYi", "Reader", 5}, "1003655": User{"Yuxuan", "Reader", 1}}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("The ReadUsersFile return %v which is not as same as what we expected %v", result, expected)
// 	}
// }

// func TestWriteDataFile(t *testing.T) {

// 	recv0Channel := make(chan map[string]interface{})
// 	send0Channel := make(chan map[string]interface{})
// 	rpcMap := make(map[int]chan map[string]interface{})
// 	var myNode = Node{true, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// 	var testSecretDetails1 = Secret{"mySectret", 100}
// 	var testSecretDetails2 = Secret{"yourSecret", 1}
// 	testDataInput := make(map[string]Secret, 5)
// 	testDataInput["1006969"] = testSecretDetails1
// 	testDataInput["1007070"] = testSecretDetails2

// 	error := myNode.WriteDataFile(testDataInput)
// 	if error != nil {
// 		fmt.Println(error)
// 		// return error
// 	}
// 	result, readError := myNode.ReadDataFile()
// 	if readError != nil {
// 		fmt.Println(error)
// 		// return error
// 	}
// 	expected := map[string]Secret{"1006969": Secret{"mySectret", 100}, "1007070": {"yourSecret", 1}, "126": Secret{"hashed_Vict0r1aSecret3", 3}, "127": Secret{"hashed_Vict0r1aSecret3", 3},
// 		"128": Secret{"hashed_Vict0r1aSecret3", 3}, "129": Secret{"hashed_Vict0r1aSecret3", 3}, "130": Secret{"hashed_Vict0r1aSecret3", 3}}
// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("The WriteDataFile return \n %v \nwhich is not as same as what we expected \n%v", result, expected)
// 	}
// }

// func TestWriteUsersFile(t *testing.T) {

// 	recv0Channel := make(chan map[string]interface{})
// 	send0Channel := make(chan map[string]interface{})
// 	rpcMap := make(map[int]chan map[string]interface{})
// 	var myNode = Node{true, 3, []int{2, 3, 5, 7, 11, 13}, recv0Channel, send0Channel, rpcMap}

// 	var testUserDetails1 = User{"Sudipta", "iLoveDistSys", 100}
// 	var testUserDetails2 = User{"Juan", "iLoveDistSys", 1}
// 	testUserInput := make(map[string]User, 5)
// 	testUserInput["1006969"] = testUserDetails1
// 	testUserInput["1007070"] = testUserDetails2

// 	error := myNode.WriteUsersFile(testUserInput)
// 	if error != nil {
// 		fmt.Println(error)
// 		// return error
// 	}
// 	result, readError := myNode.ReadUsersFile()
// 	if readError != nil {
// 		fmt.Println(error)
// 		// return error
// 	}
// 	expected := map[string]User{"1003651": User{"DevBahl", "Reader", 3}, "1003652": User{"Jose", "Reader", 2}, "1003653": User{"Aaron", "Reader", 4}, "1003654": User{"JiaYi", "Reader", 5}, "1003655": User{"Yuxuan", "Reader", 1}, "1006969": User{"Sudipta", "iLoveDistSys", 100}, "1007070": User{"Juan", "iLoveDistSys", 1}}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("The WriteUsersFile return %v which is not as same as what we expected %v", result, expected)
// 	}
// }
