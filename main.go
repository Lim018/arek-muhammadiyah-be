package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"arek-muhammadiyah-be/config"
	"arek-muhammadiyah-be/database"
	"arek-muhammadiyah-be/middleware"
	"arek-muhammadiyah-be/route"
)

func main() {
	// Define flags
	migrate := flag.Bool("migrate", false, "Run database migrations")
	seed := flag.Bool("seed", false, "Run database seeders")
	fresh := flag.Bool("fresh", false, "Drop all tables and run fresh migrations (WARNING: deletes all data)")
	help := flag.Bool("help", false, "Show help message")
	
	flag.Parse()

	// Load environment variables
	config.LoadEnv()

	// Initialize database
	database.ConnectDB()

	// Show help
	if *help {
		showHelp()
		return
	}

	// Handle migration flags
	if *fresh {
		log.Println("⚠️  Running FRESH migrations (dropping all tables)...")
		database.RunMigrationsWithDrop()
		
		if *seed {
			log.Println("🌱 Running seeders after fresh migration...")
			database.RunSeeders()
		}
		
		log.Println("✅ Fresh migration completed!")
		return
	}

	if *migrate {
		log.Println("🚀 Running migrations...")
		database.RunMigrations()
		
		if *seed {
			log.Println("🌱 Running seeders...")
			database.RunSeeders()
		}
		
		log.Println("✅ Migration completed!")
		return
	}

	if *seed {
		log.Println("🌱 Running seeders only...")
		database.RunSeeders()
		log.Println("✅ Seeding completed!")
		return
	}

	// Normal server startup (no flags)
	log.Println("🚀 Starting server...")
	startServer()
}

func startServer() {
	// Create Fiber app
	app := config.CreateApp()

	// Setup middleware
	middleware.Setup(app)

	// Setup routes
	route.Setup(app)

	// Start server
	log.Println("✅ Server is running on http://localhost:8080")
	log.Fatal(app.Listen("0.0.0.0:8080"))
}

func showHelp() {
	fmt.Println(`
╔════════════════════════════════════════════════════════════════╗
║          Arek Muhammadiyah Backend - Command Line              ║
╚════════════════════════════════════════════════════════════════╝

USAGE:
	go run main.go [flags]

FLAGS:
	-migrate        Run database migrations
	-seed           Run database seeders
	-fresh          Drop all tables and run fresh migrations (⚠️  DELETES ALL DATA)
	-help           Show this help message

EXAMPLES:
	# Start server normally
	go run main.go

	# Run migrations only
	go run main.go -migrate

	# Run migrations and seed data
	go run main.go -migrate -seed

	# Run seeders only (tables must exist)
	go run main.go -seed

	# Fresh install (drop all tables, migrate, and seed)
	go run main.go -fresh -seed

	# Show help
	go run main.go -help

NOTES:
	- Use -fresh flag with CAUTION as it will DELETE ALL DATA
	- Always backup your database before running -fresh
	- Run migrations before seeders for fresh installations

════════════════════════════════════════════════════════════════`)
	os.Exit(0)
}