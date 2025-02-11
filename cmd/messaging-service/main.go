package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the application using Wire.
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer app.RabbitMQ.Close()

	log.Printf("Server starting on port %s", "3000")
	if err := http.ListenAndServe(":3000", app.Router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
