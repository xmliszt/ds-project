package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// FileMethods contains all the methods associated with manipulating OS files
type FileMethods interface {
	ReadUsersFile() map[string]User
	WriteUsersFile()
	ReadDataFile() map[string]Secret
	WriteDataFile()
}

// ReadUsersFile returns all the information from the global users.json file
// Code adapted from: https://tutorialedge.net/golang/parsing-json-with-golang/
func (n *Node) ReadUsersFile() (map[string]User, error) {

	jsonFile, osErr := os.Open("../../users.json")

	if osErr != nil { // if we os.Open returns an error then handle it
		// fmt.Println(osErr)
		return nil, osErr
		// exit program
	}

	fmt.Println("Successfully Opened users.json")

	defer jsonFile.Close()

	byteValue, readAllError := ioutil.ReadAll(jsonFile)
	if readAllError != nil {
		// fmt.Println(readAllError)
		return nil, readAllError
	}

	var fileContents map[string]User

	// Unmarshal parses the byteValue array to a type defined by fileContents
	marshalError := json.Unmarshal([]byte(byteValue), &fileContents)
	if marshalError != nil { // if we os.Open returns an error then handle it
		// fmt.Println(marshalError)
		return nil, marshalError
		// exit program
	}
	return fileContents, nil
}

// ReadDataFile returns all the information from the data.json of the respective node's local file
func (n *Node) ReadDataFile() (map[string]Secret, error) {

	filePath := dataFilePathNode(n.Pid)
	jsonFile, osErr := os.Open(filePath)

	if osErr != nil {
		// fmt.Println(osErr)
		return nil, osErr
	}

	fmt.Println("Successfully Opened", filePath)

	defer jsonFile.Close()

	byteValue, readAllError := ioutil.ReadAll(jsonFile)
	if readAllError != nil {
		// fmt.Println(readAllError)
		return nil, readAllError
	}

	var fileContents map[string]Secret
	marshalError := json.Unmarshal([]byte(byteValue), &fileContents)
	if marshalError != nil { // if we os.Open returns an error then handle it
		// fmt.Println(marshalError)
		return nil, marshalError
		// exit program
	}

	return fileContents, nil
}

// WriteUsersFile takes updates the user file with the new users provided
func (n *Node) WriteUsersFile(addUsers map[string]User) error {

	originalFileContent, readError := n.ReadUsersFile()

	if readError != nil {
		// fmt.Println(readError)
		return readError
	}
	fmt.Println("Original:\n", originalFileContent)

	// Update the values from the file
	for key, value := range addUsers {
		originalFileContent[key] = value
	}

	fmt.Println("Edited:\n", originalFileContent)

	file, marshallError := json.MarshalIndent(originalFileContent, "", " ")
	if marshallError != nil {
		// fmt.Println(marshallError)
		return marshallError

	}

	var writeError = ioutil.WriteFile("../../users.json", file, 0644)
	if writeError != nil {
		// fmt.Println(writeError)
		return writeError
	}
	// fmt.Println(n.ReadUsersFile())
	return nil
}

// WriteDataFile taks in the variable with map type then update the user file
func (n *Node) WriteDataFile(addData map[string]Secret) error {

	filePath := dataFilePathNode(n.Pid)
	originalFileContent, readError := n.ReadDataFile()

	if readError != nil {
		// fmt.Println(readError)
		return readError
	}
	fmt.Println("Original:\n", originalFileContent)

	for key, value := range addData {
		originalFileContent[key] = value
	}

	fmt.Println("Edited:\n", originalFileContent)

	file, marshallError := json.MarshalIndent(originalFileContent, "", " ")
	if marshallError != nil {
		// fmt.Println(marshallError)
		return marshallError
	}

	var writeError = ioutil.WriteFile(filePath, file, 0644)
	if writeError != nil {
		// fmt.Println(writeError)
		return writeError
	}
	// fmt.Println(n.ReadDataFile())
	return nil
}

func dataFilePathNode(nodePID int) string {
	id := strconv.Itoa(nodePID)
	dataFilePath := "../../nodeStorage/node" + id + "/data.json"
	return dataFilePath
}

// OverwriteDatafromFile taks in the variable with map type then update the user file
func (n *Node) OverwriteDataFile(addData map[string]Secret) error {

	filePath := dataFilePathNode(n.Pid)

	file, marshallError := json.MarshalIndent(addData, "", " ")
	if marshallError != nil {
		// fmt.Println(marshallError)
		return marshallError
	}

	var writeError = ioutil.WriteFile(filePath, file, 0644)
	if writeError != nil {
		// fmt.Println(writeError)
		return writeError
	}
	// fmt.Println(n.ReadDataFile())
	return nil
}
