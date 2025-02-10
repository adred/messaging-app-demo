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
	defer app.DB.Close()
	defer app.RabbitMQ.Close()

	log.Printf("Server starting on port %s", "3001")
	if err := http.ListenAndServe(":3001", app.Router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
