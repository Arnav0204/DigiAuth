package main

import (
	"context"
	"digiauth/auth"
	"digiauth/database"
	"digiauth/issuer"
	"digiauth/receiver"
	"log"
	"net/http"
)

func main() {

	database.InitDB()
	defer func() {
		if err := database.DB.Close(context.Background()); err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}()

	go func() {
		authRoute := auth.RegisterRoutes()
		log.Println("Starting auth server on :1010")
		if err := http.ListenAndServe(":1010", authRoute); err != nil {
			log.Fatalf("Issuer server failed: %v", err)
		}
	}()

	// start issuer server
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
