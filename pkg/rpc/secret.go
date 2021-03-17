package rpc

import (
	"strconv"
)

type Secret struct {
	// Key   int
	Value string // Value of secret
	Role  int    // Int to identify role (clearance level)
}

// GetSecret returns the secret that is referenced by the specific id
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

// GetSecrets returns secrets available within the range given
func (n *Node) GetSecrets(from int, to int) (interface{}, error) {
	allData, fileError := n.ReadDataFile()
	if fileError != nil {
		return nil, fileError
	} else {
		specifcData := make(map[string]Secret)
		for dictKeyInt := from; dictKeyInt <= to; dictKeyInt++ {
			stringKey := strconv.Itoa(dictKeyInt)
			specifcData[stringKey] = allData[stringKey]
		}
		return specifcData, nil
	}
}

// PutSecret place a specific secret within the respective data file
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
		overwriteError := n.OverwriteDataFile(allData)
		if overwriteError != nil {
			// fmt.Println(overwriteError)
			return overwriteError
		}
		return nil

	}

}
