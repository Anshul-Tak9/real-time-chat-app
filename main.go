package main

import (
	"log"
	"real-time-chat-app/db"
	"real-time-chat-app/router"
)

func main() {
	// Initialize MongoDB
	if err := db.InitializeDB(); err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}

	// // Start the WebSocket hub
	// go chat.RunHub()
	// go chat.StartRedisConsumer()

	// Get the fully-configured router
	r := router.Initialize()

	// Run the server
	if err := r.Run(":4000"); err != nil {
		panic(err)
	}
}
