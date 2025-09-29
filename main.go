package main

import (
	"log"
	"github.com/Lim018/arek-muhammadiyah-be/config"
	"github.com/Lim018/arek-muhammadiyah-be/database"
	"github.com/Lim018/arek-muhammadiyah-be/middleware"
	"github.com/Lim018/arek-muhammadiyah-be/route"
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
	log.Fatal(app.Listen(":8080"))
}