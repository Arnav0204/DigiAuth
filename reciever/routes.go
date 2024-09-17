package reciever

import "github.com/gorilla/mux"

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	// r.HandleFunc("/create-credential-defination", CreateCredentialDefination).Methods("POST")
	// r.HandleFunc("/register-certificate", RegisterSchema).Methods("POST")
	// r.HandleFunc("/register-did", RegisterDID).Methods("POST")
	return r
}
