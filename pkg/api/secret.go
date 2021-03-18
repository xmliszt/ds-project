package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xmliszt/e-safe/pkg/secret"
)

// Put a secret
func PutSecret(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var secret secret.Secret
	err := json.Unmarshal(reqBody, &secret)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}
	fmt.Fprintf(w, "%+v", string(reqBody))
}

// Get a secret
func GetSecret(w http.ResponseWriter, r *http.Request) {

}
