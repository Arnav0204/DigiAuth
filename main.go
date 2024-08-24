package main

import (
	"digiauth/issuer"
	"log"
	"net/http"
)

func main() {
	route := issuer.RegisterRoutes()
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", route)
}
