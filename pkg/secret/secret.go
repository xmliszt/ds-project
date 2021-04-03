package secret

import (
	"errors"
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

func decodeSecret(secret Secret) interface{} {
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

// PutSecret place a specific secret within the respective data file
func PutSecret(pid int, key string, secret Secret) error {
	// Write to this node node first.
	writeErr := writeSecret(pid, key, secret)
	if writeErr != nil {
		return nil
	} else {
		return writeErr
	}

	// Abstract this to another function
	// secretVal := decodeSecret(secret)
	// newSecret := map[string]interface{}{key: secretVal}
	// fileError := file.WriteDataFile(pid, newSecret)
	// if fileError != nil {
	// 	return nil
	// } else {
	// 	return fileError
	// }
}

func writeSecret(pid int, key string, secret Secret) error {
	secretVal := decodeSecret(secret)
	newSecret := map[string]interface{}{key: secretVal}
	fileError := file.WriteDataFile(pid, newSecret)
	if fileError != nil {
		return nil
	} else {
		return fileError
	}
}

// DeleteSecret deletes a secret from a data file from a given pid machine
func DeleteSecret(pid int, key string) error {
	allData, fileError := file.ReadDataFile(pid)
	if fileError != nil {
		return fileError
	} else {
		delete(allData, key)
		overwriteError := file.OverwriteDataFile(pid, allData)
		if overwriteError != nil {
			return overwriteError
		}
		return nil
	}
}
