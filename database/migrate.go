package database

import (
	"arek-muhammadiyah-be/app/model"
	"log"
)

func RunMigrations() {
	if DB == nil {
		log.Fatal("Database not connected")
	}

	// Auto migrate tables
	err := DB.AutoMigrate(
		&model.Role{},
		&model.Village{},
		&model.Category{},
		&model.User{},
		&model.Document{},
		&model.Ticket{},
		&model.Article{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration completed successfully")
}
