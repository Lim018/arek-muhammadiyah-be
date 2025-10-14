package database

import (
	"arek-muhammadiyah-be/app/model"
	"log"
)

func RunSeeders() {
	if DB == nil {
		log.Fatal("Database not connected")
	}

	// Seed roles
	adminDesc := "Administrator dengan akses penuh ke sistem"
	operatorDesc := "Operator sensus yang menginput data anggota"
	memberDesc := "Anggota organisasi"
	
	roles := []model.Role{
		{Name: "admin", Description: &adminDesc},
		{Name: "operator", Description: &operatorDesc},
		{Name: "member", Description: &memberDesc},
	}
	for _, r := range roles {
		DB.FirstOrCreate(&r, model.Role{Name: r.Name})
	}

	
	pendidikanDesc := "Bidang pendidikan dan pengembangan SDM"
	kesehatanDesc := "Bidang kesehatan dan pelayanan medis"
	ekonomiDesc := "Bidang ekonomi dan pemberdayaan usaha"
	keagamaanDesc := "Bidang keagamaan dan spiritual"
	
	categories := []model.Category{
		{Name: "Pendidikan", Description: &pendidikanDesc, Color: "#3B82F6"},
		{Name: "Kesehatan", Description: &kesehatanDesc, Color: "#10B981"},
		{Name: "Ekonomi", Description: &ekonomiDesc, Color: "#F59E0B"},
		{Name: "Keagamaan", Description: &keagamaanDesc, Color: "#8B5CF6"},
	}
	for _, c := range categories {
		DB.FirstOrCreate(&c, model.Category{Name: c.Name})
	}

	log.Println("Seeding completed successfully")
}
