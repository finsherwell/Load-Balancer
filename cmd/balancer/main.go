package main

import (
	"log"

	"github.com/finsherwell/Load-Balancer/internal/server"
)

func main() {
	// Load configuration and start the server
	err := server.StartServer("config/config.json")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
