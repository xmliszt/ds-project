package rpc

import "fmt"

type Secret struct {
	// Key   int
	Value string // Value of secret
	Role  int    // Int to identify role (clearance level)
}

func (n Node) GetSecret(id string) (interface{}, error) {
	allData, fileError := n.ReadDataFile()
	if fileError != nil {
		return nil, fileError
	} else {
		for key := range allData {
			if key == id {
				return allData[key], nil
			}
		}
	}
	return nil, nil
}

func (n *Node) GetSecrets(from int, to int) (interface{}, error) {
	allData, fileError := n.ReadDataFile()
	if fileError != nil {
		return nil, fileError
	} else {
		return allData, nil
	}
}

func (n *Node) PutSecret(key string, secret Secret) error {
	newSecret := SecretToMapSecret(secret, key)
	fileError := n.WriteDataFile(newSecret)
	if fileError != nil {
		return nil
	} else {
		return fileError
	}
}

func SecretToMapSecret(secret Secret, key string) map[string]Secret {
	newSecret := make(map[string]Secret)
	newSecret[key] = secret
	return newSecret
}

func (n *Node) DeleteSecret(key string) error {
	allData, fileError := n.ReadDataFile()
	if fileError != nil {
		return fileError
	} else {
		delete(allData, key)
		fmt.Println("HERE:\n\n", allData, "\n\n")
		// for key, value := range allData {
		// 	n.PutSecret(key, value)
		// }
		n.OverwriteDataFile(allData)
		return nil

	}

}
