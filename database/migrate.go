package database

import (
	"arek-muhammadiyah-be/app/model"
	"log"
)

func RunMigrations() {
	if DB == nil {
		log.Fatal("Database not connected")
	}

	log.Println("🚀 Running database migrations...")

	models := []interface{}{
		&model.Role{},
		&model.Village{},
		&model.SubVillage{},  // Tambah SubVillage
		&model.Category{},
		&model.User{},
		&model.Document{},
		&model.Ticket{},
		&model.Article{},
	}

	for _, m := range models {
		if err := DB.AutoMigrate(m); err != nil {
			// Cuma tampilkan warning, tapi lanjut model berikutnya
			log.Printf("⚠️  Skipped migration for %T: %v\n", m, err)
			continue
		}
		log.Printf("✅ Migrated %T\n", m)
	}

	log.Println("✅ Migration completed successfully")
}

func RunMigrationsWithDrop() {
	if DB == nil {
		log.Fatal("Database not connected")
	}

	log.Println("⚠️  WARNING: Dropping all tables and recreating...")

	models := []interface{}{
		&model.Article{},
		&model.Ticket{},
		&model.Document{},
		&model.User{},
		&model.Category{},
		&model.SubVillage{},
		&model.Village{},
		&model.Role{},
	}

	for _, m := range models {
		if err := DB.Migrator().DropTable(m); err != nil {
			log.Printf("⚠️  Failed to drop table %T: %v\n", m, err)
		} else {
			log.Printf("🗑️  Dropped table %T\n", m)
		}
	}

	log.Println("🔄 Recreating tables...")
	RunMigrations()
}