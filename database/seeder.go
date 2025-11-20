package database

import (
	"log"
	"time"
	"arek-muhammadiyah-be/app/model"
	"golang.org/x/crypto/bcrypt"
)

func ptr(s string) *string { return &s }
func uintPtr(u uint) *uint { return &u }

func parseDate(dateStr string) *time.Time {
	layout := "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Printf("⚠️  Failed to parse date %s: %v", dateStr, err)
		return nil
	}
	return &t
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}


func RunSeeders() {
	if DB == nil {
		log.Fatal("Database not connected")
	}

	log.Println("🌱 Running seeders...")

	// === Roles ===
	log.Println("🔑 Seeding roles...")
	adminDesc := "Administrator dengan akses penuh ke sistem"
	operatorDesc := "Operator sensus yang menginput data anggota"
	memberDesc := "Anggota organisasi"

	roles := []model.Role{
		{Name: "admin", Description: &adminDesc},
		{Name: "operator", Description: &operatorDesc},
		{Name: "member", Description: &memberDesc},
	}
	
	roleMap := make(map[string]uint)
	for _, r := range roles {
		DB.FirstOrCreate(&r, model.Role{Name: r.Name}) 
		roleMap[r.Name] = r.ID
		log.Printf("  ✓ Role: %s (ID: %d)\n", r.Name, r.ID)
	}

	log.Println("📁 Seeding categories...")
	pendidikanDesc := "Bidang pendidikan dan pengembangan SDM"
	kesehatanDesc := "Bidang kesehatan dan pelayanan medis"
	ekonomiDesc := "Bidang ekonomi dan pemberdayaan usaha"
	keagamaanDesc := "Bidang keagamaan dan spiritual"

	categories := []model.Category{
		{Name: "Pendidikan", Description: &pendidikanDesc, Color: "#3B82F6", Icon: "book"},
		{Name: "Kesehatan", Description: &kesehatanDesc, Color: "#3B82F6", Icon: "heart"},
		{Name: "Ekonomi", Description: &ekonomiDesc, Color: "#3B82F6", Icon: "chart"},
		{Name: "Keagamaan", Description: &keagamaanDesc, Color: "#3B82F6", Icon: "mosque"},
	}
	for _, c := range categories {
		DB.FirstOrCreate(&c, model.Category{Name: c.Name})
		log.Printf("  ✓ Category: %s\n", c.Name)
	}


	log.Println("👤 Seeding Admin User...")

	adminPassword := "password123"
	adminHashedPassword, err := HashPassword(adminPassword)
	if err != nil {
		log.Fatalf("  ✗ Gagal hash password admin: %v", err)
	}

	adminUser := model.User{
		Name:      "Admin Utama",
		Password:  adminHashedPassword,
		BirthDate: parseDate("1990-01-01"),
		Telp:      ptr("081234567890"),
		Gender:    ptr("male"),
		Job:       ptr("System Administrator"),
		RoleID:    uintPtr(roleMap["admin"]),
		VillageID: ptr("3578011001"),
		NIK:       ptr("3578010101900001"),
		Address:   ptr("Kantor Pusat"),
		IsMobile:  false,
	}

	var existing model.User
	result := DB.Where("nik = ? OR telp = ?", adminUser.NIK, adminUser.Telp).First(&existing)
	
	if result.Error == nil {
		log.Printf("  ⊘ Admin user sudah ada (NIK: %s / Telp: %s)\n", *adminUser.NIK, *adminUser.Telp)
	} else {
		if err := DB.Create(&adminUser).Error; err != nil {
			log.Printf("  ✗ Gagal create admin user: %v\n", err)
		} else {
			log.Printf("  ✓ Admin User: %s (Telp: %s, Pass: %s)\n", 
				adminUser.Name, *adminUser.Telp, adminPassword)
		}
	}

	log.Println("✅ Seeders (Roles, Categories, Admin) completed successfully")
}