package secret

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/xmliszt/e-safe/pkg/file"
	"github.com/xmliszt/e-safe/util"
)

type Secret struct {
	Alias string `json:"alias"` // Alias of secret (identifier)
	Value string `json:"value"` // Value of secret
	Role  int    `json:"role"`  // Int to identify role (clearance level)
}

func encodeSecret(data interface{}) (*Secret, error) {
	secret := &Secret{}
	for key, val := range data.(map[string]interface{}) {
		if key == "Role" {
			val = int(val.(float64))
		}
		err := util.SetField(secret, key, val)
		if err != nil {
			return nil, err
		}
	}
	return secret, nil
}

func decodeSecret(secret *Secret) interface{} {
	var r interface{} = secret
	return r
}

// GetSecret returns the secret that is referenced by the specific id
func GetSecret(pid int, id string) (*Secret, error) {
	allData, fileError := file.ReadDataFile(pid)
	if fileError != nil {
		return nil, fileError
	} else {
		for key, val := range allData {
			if key == id {
				secret, err := encodeSecret(val)
				if err != nil {
					return nil, err
				}
				return secret, nil
			}
		}
		unknownSecretIDError := errors.New("use	rname not available")
		return nil, unknownSecretIDError
	}

}

// GetSecrets returns secrets available within the range given
func GetSecrets(pid int, from int, to int) (map[string]*Secret, error) {
	allData, fileError := file.ReadDataFile(pid)
	if fileError != nil {
		return nil, fileError
	} else {
		specificData := make(map[string]*Secret)
		for dictKeyInt := from; dictKeyInt <= to; dictKeyInt++ {
			stringKey := strconv.Itoa(dictKeyInt)
			secretVal := allData[stringKey]
			secret, err := encodeSecret(secretVal)
			if err != nil {
				return nil, err
			}
			specificData[stringKey] = secret
		}
		return specificData, nil
	}
}

// PutSecret handles creation of a secret
// if the secret key exists, throw error
func PutSecret(pid int, key string, secret *Secret) error {
	allData, fileError := file.ReadDataFile(pid)
	if fileError != nil {
		return fileError
	}

	if _, ok := allData[key]; ok {
		return fmt.Errorf("key [%s] exists, creation failed", key)
	}

	secretVal := decodeSecret(secret)
	secretToPut := map[string]interface{}{key: secretVal}
	fileError = file.WriteDataFile(pid, secretToPut)
	if fileError != nil {
		return fileError
	}
	return nil
}

// UpdateSecret updates an existing secret or delete a secret (update it to null)
func UpdateSecret(pid int, key string, secret *Secret) error {
	var secretVal interface{}
	if secret == nil {
		secretVal = nil
	} else {
		secretVal = decodeSecret(secret)
	}
	allData, fileError := file.ReadDataFile(pid)
	if fileError != nil {
		return fileError
	}

	// check if the key exists in the data file
	// if exist, overwrite
	// if does not exist, throw error
	if _, ok := allData[key]; !ok {
		return fmt.Errorf("key [%s] does not exist", key)
	}
	newSecret := map[string]interface{}{key: secretVal}
	fileError = file.WriteDataFile(pid, newSecret)
	if fileError != nil {
		return fileError
	}
	return nil
}
