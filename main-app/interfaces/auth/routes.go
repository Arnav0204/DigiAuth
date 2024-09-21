package auth

import "github.com/gorilla/mux"

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", RegisterUser).Methods("POST")
	r.HandleFunc("/login", LoginUser).Methods("POST")
	return r
}
