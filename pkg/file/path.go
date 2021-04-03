package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateStoragePath setup necessary directory for storage
// includes user and secret storage
func CreateStoragePath() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	nodeStoragePath := filepath.Join(cwd, "nodeStorage")
	if _, err := os.Stat(nodeStoragePath); os.IsNotExist(err) {
		err := os.Mkdir(nodeStoragePath, 0777)
		if err != nil {
			return err
		}
	}

	userPath := filepath.Join(cwd, "nodeStorage", "user")
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		err := os.Mkdir(userPath, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateNodeStoragePath creates the particular directory for a node's storage
// it will also create nodeStorage/ parent directory if it is not found
func CreateNodeStoragePath(nodeID int) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	nodeStoragePath := filepath.Join(cwd, "nodeStorage")
	if _, err := os.Stat(nodeStoragePath); os.IsNotExist(err) {
		err := os.Mkdir(nodeStoragePath, 0777)
		if err != nil {
			return err
		}
	}

	storagePath := filepath.Join(cwd, "nodeStorage", fmt.Sprintf("node%d", nodeID))
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		err := os.Mkdir(storagePath, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetNodeFilePath returns the file path of secret that the node stores
func GetNodeFilePath(nodeID int) (string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	storagePath := filepath.Join(cwd, "nodeStorage", fmt.Sprintf("node%d", nodeID))
	return filepath.Join(storagePath, "data.json"), nil
}

// GetUserFilePath returns the file path of user information
func GetUserFilePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	userPath := filepath.Join(cwd, "nodeStorage", "user")
	return filepath.Join(userPath, "users.json"), nil
}
