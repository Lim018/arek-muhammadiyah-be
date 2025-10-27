package database

import (
	"log"
	"time"
	"arek-muhammadiyah-be/app/model"
	"golang.org/x/crypto/bcrypt"
)

func RunSeeders() {
	if DB == nil {
		log.Fatal("Database not connected")
	}

	log.Println("🌱 Running seeders...")

	// === Roles ===
	log.Println("🔐 Seeding roles...")
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
		log.Printf("  ✓ Role: %s\n", r.Name)
	}

	// === Categories ===
	log.Println("📁 Seeding categories...")
	pendidikanDesc := "Bidang pendidikan dan pengembangan SDM"
	kesehatanDesc := "Bidang kesehatan dan pelayanan medis"
	ekonomiDesc := "Bidang ekonomi dan pemberdayaan usaha"
	keagamaanDesc := "Bidang keagamaan dan spiritual"

	categories := []model.Category{
		{Name: "Pendidikan", Description: &pendidikanDesc, Color: "#3B82F6", Icon: "book"},
		{Name: "Kesehatan", Description: &kesehatanDesc, Color: "#10B981", Icon: "heart"},
		{Name: "Ekonomi", Description: &ekonomiDesc, Color: "#F59E0B", Icon: "chart"},
		{Name: "Keagamaan", Description: &keagamaanDesc, Color: "#8B5CF6", Icon: "mosque"},
	}
	for _, c := range categories {
		DB.FirstOrCreate(&c, model.Category{Name: c.Name})
		log.Printf("  ✓ Category: %s\n", c.Name)
	}

	// === Users ===
	log.Println("👥 Seeding users...")
	log.Println("   Note: Village IDs must match your wilayah.json file")

	// Parse birth dates
	birthDate1 := parseDate("1990-05-15")
	birthDate2 := parseDate("1985-08-20")
	birthDate3 := parseDate("1992-03-10")
	birthDate4 := parseDate("1988-11-25")
	birthDate5 := parseDate("1995-07-30")

	// IMPORTANT: Sesuaikan village_id dengan data di wilayah.json Anda
	// Format: "3576011001" (10 digit - ID kelurahan)
	users := []model.User{
		{
			Name:      "Budi Santoso",
			Password:  "password123",
			BirthDate: birthDate1,
			Telp:      ptr("081234567890"),
			Gender:    ptr("male"),
			Job:       ptr("Guru"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3576011001"), // Kelurahan Blooto, Kec. Prajuritkulon, Kota Mojokerto
			NIK:       ptr("3576011234567890"),
			Address:   ptr("Jl. Merdeka No. 15, Blooto"),
			IsMobile:  true,
		},
		{
			Name:      "Siti Nurhaliza",
			Password:  "password123",
			BirthDate: birthDate2,
			Telp:      ptr("081234567891"),
			Gender:    ptr("female"),
			Job:       ptr("Dokter"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3576011002"), // Kelurahan Kauman
			NIK:       ptr("3576012234567891"),
			Address:   ptr("Jl. Kauman No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Ahmad Fauzi",
			Password:  "password123",
			BirthDate: birthDate3,
			Telp:      ptr("081234567892"),
			Gender:    ptr("male"),
			Job:       ptr("Pengusaha"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3576021001"), // Kelurahan Balongsari, Kec. Magersari
			NIK:       ptr("3576013234567892"),
			Address:   ptr("Jl. Balongsari No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Dewi Lestari",
			Password:  "password123",
			BirthDate: birthDate4,
			Telp:      ptr("081234567893"),
			Gender:    ptr("female"),
			Job:       ptr("Perawat"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3576021002"), // Kelurahan Gedongan
			NIK:       ptr("3576014234567893"),
			Address:   ptr("Jl. Gedongan No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Joko Widodo",
			Password:  "password123",
			BirthDate: birthDate5,
			Telp:      ptr("081234567894"),
			Gender:    ptr("male"),
			Job:       ptr("Wiraswasta"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3576031001"), // Kelurahan Jagalan, Kec. Kranggan
			NIK:       ptr("3576015234567894"),
			Address:   ptr("Jl. Jagalan No. 10"),
			IsMobile:  false,
		},
	}

	for _, u := range users {
		// Hash password dulu sebelum simpan
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
		if err != nil {
			log.Printf("  ✗ Gagal hash password untuk user %s: %v", u.Name, err)
			continue
		}
		u.Password = string(hashed)

		var existing model.User
		result := DB.Where("nik = ?", u.NIK).First(&existing)
		if result.Error == nil {
			log.Printf("  ⊘ User sudah ada: %s\n", u.Name)
			continue
		}

		if err := DB.Create(&u).Error; err != nil {
			log.Printf("  ✗ Gagal create user %s: %v\n", u.Name, err)
			continue
		}
		log.Printf("  ✓ User: %s (Age: %d, Village ID: %s)\n", u.Name, u.GetAge(), *u.VillageID)
	}

	log.Println("✅ All seeders completed successfully")
	log.Println("📍 Make sure your wilayah.json contains the village IDs used above")
}

// Helper functions
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