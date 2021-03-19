package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xmliszt/e-safe/pkg/user"
)

// Create a user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user user.User
	err := json.Unmarshal(reqBody, &user)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}
	fmt.Fprintf(w, "%+v", string(reqBody))
}
