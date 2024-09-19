package main

import (
	"digiauth/issuer"
	"digiauth/receiver"
	"log"
	"net/http"
)

func main() {
	// issuerRoute := issuer.RegisterRoutes()
	// log.Println("Starting  issuer server on :8080")
	// http.ListenAndServe(":8080", issuerRoute)
	// recieverRoute := reciever.RegisterRoutes()
	// log.Println("Starting reciever server on :6060")
	// http.ListenAndServe(":6060", recieverRoute)

	go func() {
		issuerRoute := issuer.RegisterRoutes()
		log.Println("Starting issuer server on :8080")
		if err := http.ListenAndServe(":8080", issuerRoute); err != nil {
			log.Fatalf("Issuer server failed: %v", err)
		}
	}()

	// Start receiver server
	go func() {
		receiverRoute := receiver.RegisterRoutes()
		log.Println("Starting receiver server on :6060")
		if err := http.ListenAndServe(":6060", receiverRoute); err != nil {
			log.Fatalf("Receiver server failed: %v", err)
		}
	}()
	// Block forever to keep the main Goroutine alive
	select {}
}
