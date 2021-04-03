package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// FileMethods contains all the methods associated with manipulating OS files
type FileMethods interface {
	ReadUsersFile() map[string]interface{}
	WriteUsersFile()
	ReadDataFile() map[string]interface{}
	WriteDataFile()
}

// ReadUsersFile returns all the information from the global users.json file
// Code adapted from: https://tutorialedge.net/golang/parsing-json-with-golang/
func ReadUsersFile() (map[string]interface{}, error) {

	userFilePath, err := GetUserFilePath()
	if err != nil {
		return nil, err
	}

	// if file not exists, create it
	var jsonFile *os.File
	var osErr error
	if _, err := os.Stat(userFilePath); err != nil {
		jsonFile, osErr = os.Create(userFilePath)
		if osErr != nil {
			return nil, osErr
		}
	} else {
		jsonFile, osErr = os.Open(userFilePath)
		if osErr != nil {
			return nil, osErr
		}
	}

	defer jsonFile.Close()

	byteValue, readAllError := ioutil.ReadAll(jsonFile)
	if readAllError != nil {
		return nil, readAllError
	}

	var fileContents map[string]interface{}

	// Unmarshal parses the byteValue array to a type defined by fileContents
	marshalError := json.Unmarshal([]byte(byteValue), &fileContents)
	if marshalError != nil {
		return nil, marshalError
	}
	return fileContents, nil
}

// ReadDataFile returns all the information from the data.json of the respective node's local file
func ReadDataFile(pid int) (map[string]interface{}, error) {

	filePath, err := GetNodeFilePath(pid)
	if err != nil {
		return nil, err
	}

	// if file not exists, create it
	var jsonFile *os.File
	var osErr error
	if _, err := os.Stat(filePath); err != nil {
		jsonFile, osErr = os.Create(filePath)
		if osErr != nil {
			return nil, osErr
		}
	} else {
		jsonFile, osErr = os.Open(filePath)
		if osErr != nil {
			return nil, osErr
		}
	}

	if osErr != nil {
		return nil, osErr
	}

	defer jsonFile.Close()

	byteValue, readAllError := ioutil.ReadAll(jsonFile)
	if readAllError != nil {
		return nil, readAllError
	}

	var fileContents map[string]interface{}
	marshalError := json.Unmarshal([]byte(byteValue), &fileContents)
	if marshalError != nil { // if we os.Open returns an error then handle it
		return nil, marshalError
	}

	return fileContents, nil
}

// WriteUsersFile takes updates the user file with the new users provided
func WriteUsersFile(addUsers map[string]interface{}) error {

	originalFileContent, log := ReadUsersFile()

	if log != nil {
		return log
	}

	// Update the values from the file
	for key, value := range addUsers {
		originalFileContent[key] = value
	}

	file, marshallError := json.MarshalIndent(originalFileContent, "", " ")
	if marshallError != nil {
		return marshallError

	}

	userFilePath, err := GetUserFilePath()
	if err != nil {
		return err
	}
	var writeError = ioutil.WriteFile(userFilePath, file, 0644)
	if writeError != nil {
		return writeError
	}
	return nil
}

// WriteDataFile taks in the variable with map type then update the user file
func WriteDataFile(pid int, addData map[string]interface{}) error {

	filePath, err := GetNodeFilePath(pid)
	if err != nil {
		return err
	}
	originalFileContent, log := ReadDataFile(pid)

	if log != nil {
		return log
	}

	for key, value := range addData {
		originalFileContent[key] = value
	}

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

// OverwriteDatafromFile taks in the variable with map type then update the user file
func OverwriteDataFile(pid int, addData map[string]interface{}) error {

	filePath, err := GetNodeFilePath(pid)
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
