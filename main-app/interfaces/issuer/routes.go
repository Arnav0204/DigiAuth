package issuer

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register-certificate", RegisterSchema).Methods("POST")
	r.HandleFunc("/register-did", RegisterDID).Methods("POST")
	r.HandleFunc("/send-invitation", CreateInvitation).Methods("POST")
	r.HandleFunc("/receive-invitation", ReceiveInvitation).Methods("POST")
	r.HandleFunc("/connections", GetConnections).Methods("POST")
	r.HandleFunc("/issue-credential", IssueCredential).Methods("POST")
	return r
}
