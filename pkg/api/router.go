package api

import "github.com/gorilla/mux"

func GetRouter() mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// General
	router.HandleFunc("/", Home).Methods("GET")

	// User
	router.HandleFunc("/user", CreateUser).Methods("POST")

	// Secret
	router.HandleFunc("/secret", GetSecret).Methods("GET")
	router.HandleFunc("/secret", PutSecret).Methods("POST")

	return *router
}
