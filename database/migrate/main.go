package main

import (
	"log"

	"arek-muhammadiyah-be/database"
)

func main() {
	// Hubungkan ke database
	log.Println("🔌 Connecting to database...")
	database.ConnectDB()

	// Jalankan migrasi
	log.Println("🚀 Running migrations...")
	database.RunMigrations()

	// Jalankan seeder
	log.Println("🌱 Running seeders...")
	database.RunSeeders()

	log.Println("✅ Migration and seeding completed successfully!")
}

// go run ./database/migrate
// tinggal run atas ini ajaa buat migrate
