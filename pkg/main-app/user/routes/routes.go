package receiver

import (
	controllers "digiauth/pkg/main-app/user/controllers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register-did", controllers.RegisterDID).Methods("POST")
	r.HandleFunc("/send-invitation", controllers.CreateInvitation).Methods("POST")
	r.HandleFunc("/receive-invitation", controllers.ReceiveInvitation).Methods("POST")
	r.HandleFunc("/connections", controllers.GetConnections).Methods("POST")
	r.HandleFunc("/credentials", controllers.GetCredentials).Methods("GET")
	r.HandleFunc("/send-presentation", controllers.SendPresentation).Methods("POST")
	return r
}
