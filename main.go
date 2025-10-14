package main

import (
	"log"
	"arek-muhammadiyah-be/config"
	"arek-muhammadiyah-be/database"
	"arek-muhammadiyah-be/middleware"
	"arek-muhammadiyah-be/route"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database
	database.ConnectDB()

	// Create Fiber app
	app := config.CreateApp()

	// Setup middleware
	middleware.Setup(app)

	// Setup routes
	route.Setup(app)

	// Start server
	log.Fatal(app.Listen("0.0.0.0:8080"))
}