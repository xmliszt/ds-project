package secret

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/xmliszt/e-safe/pkg/file"
	"github.com/xmliszt/e-safe/util"
)

type Secret struct {
	Alias string `json:"Alias"` // Alias of secret (identifier)
	Value string `json:"Value"` // Value of secret
	Role  int    `json:"Role"`  // Int to identify role (clearance level)
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
		if len(allData) <= 0 {
			return nil, fmt.Errorf("key for secret not available")
		}
		if _, ok := allData[id]; !ok {
			return nil, fmt.Errorf("key for secret not available")
		}
		for key, val := range allData {
			if key == id {
				if val == nil {
					return nil, fmt.Errorf("key for secret not available")
				}
				secret, err := encodeSecret(val)
				if err != nil {
					return nil, err
				}
				return secret, nil
			}
		}
		unknownSecretIDError := errors.New("key for secret not available")
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
		for key, secretVal := range allData {
			if secretVal == nil {
				continue
			}
			keyInt, err := strconv.Atoi(key)
			if err != nil {
				return nil, err
			}
			// tail to head case
			var secret *Secret
			if from > to {
				if (keyInt >= from && keyInt <= int(^uint32(0))) || (keyInt >= 0 && keyInt <= to) {
					if keyInt >= from && keyInt <= to {
						secret, err = encodeSecret(secretVal)
						if err != nil {
							return nil, err
						}
						specificData[key] = secret
					}
				}
			} else {
				if keyInt >= from && keyInt <= to {
					secret, err = encodeSecret(secretVal)
					if err != nil {
						return nil, err
					}
					specificData[key] = secret
				}
			}

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

// RemoveSecret is different than deletion, our "deletion" is defined
// as simply 'update' the value of secret to be nil
// RemoveSecret delete the entire entry from the data
func RemoveSecret(pid int, key string) error {
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

	delete(allData, key)
	fileError = file.OverwriteDataFile(pid, allData)
	if fileError != nil {
		return fileError
	}
	return nil
}
