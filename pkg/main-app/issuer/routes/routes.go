package issuer

import (
	controllers "digiauth/pkg/main-app/issuer/controllers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register-certificate", controllers.RegisterSchema).Methods("POST")
	r.HandleFunc("/register-did", controllers.RegisterDID).Methods("POST")
	r.HandleFunc("/send-invitation", controllers.CreateInvitation).Methods("POST")
	r.HandleFunc("/receive-invitation", controllers.ReceiveInvitation).Methods("POST")
	r.HandleFunc("/connections", controllers.GetConnections).Methods("POST")
	r.HandleFunc("/issue-credential", controllers.IssueCredential).Methods("POST")
	r.HandleFunc("/creadted-schemas", controllers.GetSchemas).Methods("GET")
	return r
}
