package main

import (
	"digiauth/issuer"
	"digiauth/reciever"
	"log"
	"net/http"
)

func main() {
	issuerRoute := issuer.RegisterRoutes()
	log.Println("Starting  issuer server on :8080")
	http.ListenAndServe(":8080", issuerRoute)
	recieverRoute := reciever.RegisterRoutes()
	log.Println("Starting reciever server on :6060")
	http.ListenAndServe(":6060", recieverRoute)
}
