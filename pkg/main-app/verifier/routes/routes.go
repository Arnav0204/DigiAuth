package verifier

import (
	controllers "digiauth/pkg/main-app/verifier/controllers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register-did", controllers.RegisterDID).Methods("POST")
	r.HandleFunc("/send-invitation", controllers.CreateInvitation).Methods("POST")
	r.HandleFunc("/receive-invitation", controllers.ReceiveInvitation).Methods("POST")
	r.HandleFunc("/connections", controllers.GetConnections).Methods("POST")
	r.HandleFunc("/send-presentation-request", controllers.SendProofRequest).Methods("POST")
	r.HandleFunc("/schemasGet", controllers.GetSchemasDB).Methods("GET")
	// r.HandleFunc("/records",controllers.GetRecords).Methods("POST")
	return r
}
