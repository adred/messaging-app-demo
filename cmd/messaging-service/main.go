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

	// Use the HTTP port from the configuration.
	log.Printf("Server starting on port %s", app.Config.HTTPPort)
	if err := http.ListenAndServe(":"+app.Config.HTTPPort, app.Router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
