package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	cwd, err := os.Getwd()
	fmt.Println(cwd)
	if err != nil {
		return nil, err
	}
	userFilePath := filepath.Join(cwd, "users.json")
	jsonFile, osErr := os.Open(userFilePath)

	if osErr != nil { // if we os.Open returns an error then handle it
		return nil, osErr
	}

	fmt.Println("Successfully Opened users.json")

	defer jsonFile.Close()

	byteValue, readAllError := ioutil.ReadAll(jsonFile)
	if readAllError != nil {
		return nil, readAllError
	}

	var fileContents map[string]User

	// Unmarshal parses the byteValue array to a type defined by fileContents
	marshalError := json.Unmarshal([]byte(byteValue), &fileContents)
	if marshalError != nil { // if we os.Open returns an error then handle it
		return nil, marshalError
	}
	return fileContents, nil
}

// ReadDataFile returns all the information from the data.json of the respective node's local file
func (n *Node) ReadDataFile() (map[string]Secret, error) {

	filePath, err := dataFilePathNode(n.Pid)
	if err != nil {
		return nil, err
	}
	jsonFile, osErr := os.Open(filePath)

	if osErr != nil {
		return nil, osErr
	}

	fmt.Println("Successfully Opened", filePath)

	defer jsonFile.Close()

	byteValue, readAllError := ioutil.ReadAll(jsonFile)
	if readAllError != nil {
		return nil, readAllError
	}

	var fileContents map[string]Secret
	marshalError := json.Unmarshal([]byte(byteValue), &fileContents)
	if marshalError != nil { // if we os.Open returns an error then handle it
		return nil, marshalError
	}

	return fileContents, nil
}

// WriteUsersFile takes updates the user file with the new users provided
func (n *Node) WriteUsersFile(addUsers map[string]User) error {

	originalFileContent, readError := n.ReadUsersFile()

	if readError != nil {
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
		return marshallError

	}

	cwd, err := os.Getwd()
	fmt.Println(cwd)
	if err != nil {
		return err
	}
	userFilePath := filepath.Join(cwd, "users.json")
	var writeError = ioutil.WriteFile(userFilePath, file, 0644)
	if writeError != nil {
		return writeError
	}
	return nil
}

// WriteDataFile taks in the variable with map type then update the user file
func (n *Node) WriteDataFile(addData map[string]Secret) error {

	filePath, err := dataFilePathNode(n.Pid)
	if err != nil {
		return err
	}
	originalFileContent, readError := n.ReadDataFile()

	if readError != nil {
		return readError
	}
	fmt.Println("Original:\n", originalFileContent)

	for key, value := range addData {
		originalFileContent[key] = value
	}

	fmt.Println("Edited:\n", originalFileContent)

	file, marshallError := json.MarshalIndent(originalFileContent, "", " ")
	if marshallError != nil {
		return marshallError
	}

	var writeError = ioutil.WriteFile(filePath, file, 0644)
	if writeError != nil {
		return writeError
	}
	return nil
}

func dataFilePathNode(nodePID int) (string, error) {
	id := strconv.Itoa(nodePID)
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dataFilePath := filepath.Join(cwd, "nodeStorage", "node"+id, "data.json")
	return dataFilePath, nil
}

// OverwriteDatafromFile taks in the variable with map type then update the user file
func (n *Node) OverwriteDataFile(addData map[string]Secret) error {

	filePath, err := dataFilePathNode(n.Pid)

	if err != nil {
		return err
	}

	file, marshallError := json.MarshalIndent(addData, "", " ")
	if marshallError != nil {
		return marshallError
	}

	var writeError = ioutil.WriteFile(filePath, file, 0644)
	if writeError != nil {
		return writeError
	}
	return nil
}
