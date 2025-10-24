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

	// === Villages (Kabupaten) ===
	log.Println("🏘️  Seeding villages (kabupaten)...")
	villages := []model.Village{
		{Name: "Surabaya", Code: "KAB-SBY-001", Description: ptr("Kota Surabaya"), Color: "#3B82F6", IsActive: true},
		{Name: "Sidoarjo", Code: "KAB-SDA-002", Description: ptr("Kabupaten Sidoarjo"), Color: "#8B5CF6", IsActive: true},
		{Name: "Gresik", Code: "KAB-GRK-003", Description: ptr("Kabupaten Gresik"), Color: "#10B981", IsActive: true},
	}
	for _, v := range villages {
		DB.FirstOrCreate(&v, model.Village{Code: v.Code})
		log.Printf("  ✓ Village: %s (%s)\n", v.Name, v.Code)
	}

	// === SubVillages (Kecamatan) ===
	log.Println("🏘️  Seeding sub villages (kecamatan)...")
	
	// Ambil village IDs
	var surabaya, sidoarjo, gresik model.Village
	DB.Where("code = ?", "KAB-SBY-001").First(&surabaya)
	DB.Where("code = ?", "KAB-SDA-002").First(&sidoarjo)
	DB.Where("code = ?", "KAB-GRK-003").First(&gresik)

	subVillages := []model.SubVillage{
		// Surabaya
		{VillageID: surabaya.ID, Name: "Gubeng", Code: "KEC-GBG-001", Description: ptr("Kecamatan Gubeng"), IsActive: true},
		{VillageID: surabaya.ID, Name: "Wonokromo", Code: "KEC-WNK-002", Description: ptr("Kecamatan Wonokromo"), IsActive: true},
		{VillageID: surabaya.ID, Name: "Sawahan", Code: "KEC-SWH-003", Description: ptr("Kecamatan Sawahan"), IsActive: true},
		{VillageID: surabaya.ID, Name: "Genteng", Code: "KEC-GTG-004", Description: ptr("Kecamatan Genteng"), IsActive: true},
		{VillageID: surabaya.ID, Name: "Tegalsari", Code: "KEC-TGS-005", Description: ptr("Kecamatan Tegalsari"), IsActive: true},
		
		// Sidoarjo
		{VillageID: sidoarjo.ID, Name: "Sidoarjo", Code: "KEC-SDA-001", Description: ptr("Kecamatan Sidoarjo"), IsActive: true},
		{VillageID: sidoarjo.ID, Name: "Waru", Code: "KEC-WRU-002", Description: ptr("Kecamatan Waru"), IsActive: true},
		{VillageID: sidoarjo.ID, Name: "Gedangan", Code: "KEC-GDG-003", Description: ptr("Kecamatan Gedangan"), IsActive: true},
		
		// Gresik
		{VillageID: gresik.ID, Name: "Gresik", Code: "KEC-GRK-001", Description: ptr("Kecamatan Gresik"), IsActive: true},
		{VillageID: gresik.ID, Name: "Kebomas", Code: "KEC-KBM-002", Description: ptr("Kecamatan Kebomas"), IsActive: true},
	}
	for _, sv := range subVillages {
		DB.FirstOrCreate(&sv, model.SubVillage{Code: sv.Code})
		log.Printf("  ✓ SubVillage: %s (%s)\n", sv.Name, sv.Code)
	}

	// === Users ===
	log.Println("👥 Seeding users...")
	
	// Ambil sub_village IDs
	var gubeng, wonokromo, sawahan, genteng, tegalsari model.SubVillage
	DB.Where("code = ?", "KEC-GBG-001").First(&gubeng)
	DB.Where("code = ?", "KEC-WNK-002").First(&wonokromo)
	DB.Where("code = ?", "KEC-SWH-003").First(&sawahan)
	DB.Where("code = ?", "KEC-GTG-004").First(&genteng)
	DB.Where("code = ?", "KEC-TGS-005").First(&tegalsari)

	// Parse birth dates
	birthDate1 := parseDate("1990-05-15")
	birthDate2 := parseDate("1985-08-20")
	birthDate3 := parseDate("1992-03-10")
	birthDate4 := parseDate("1988-11-25")
	birthDate5 := parseDate("1995-07-30")

	users := []model.User{
		{
			Name:         "Budi Santoso",
			Password:     "password123",
			BirthDate:    birthDate1,
			Telp:         ptr("081234567890"),
			Gender:       ptr("male"),
			Job:          ptr("Guru"),
			RoleID:       uintPtr(1),
			SubVillageID: &gubeng.ID,
			NIK:          ptr("3578011234567890"),
			Address:      ptr("Jl. Gubeng Pojok No. 15, Surabaya"),
			IsMobile:     true,
		},
		{
			Name:         "Siti Nurhaliza",
			Password:     "password123",
			BirthDate:    birthDate2,
			Telp:         ptr("081234567891"),
			Gender:       ptr("female"),
			Job:          ptr("Dokter"),
			RoleID:       uintPtr(1),
			SubVillageID: &wonokromo.ID,
			NIK:          ptr("3578012234567891"),
			Address:      ptr("Jl. Wonokromo No. 45, Surabaya"),
			IsMobile:     true,
		},
		{
			Name:         "Ahmad Fauzi",
			Password:     "password123",
			BirthDate:    birthDate3,
			Telp:         ptr("081234567892"),
			Gender:       ptr("male"),
			Job:          ptr("Pengusaha"),
			RoleID:       uintPtr(2),
			SubVillageID: &sawahan.ID,
			NIK:          ptr("3578013234567892"),
			Address:      ptr("Jl. Sawahan No. 78, Surabaya"),
			IsMobile:     false,
		},
		{
			Name:         "Dewi Lestari",
			Password:     "password123",
			BirthDate:    birthDate4,
			Telp:         ptr("081234567893"),
			Gender:       ptr("female"),
			Job:          ptr("Perawat"),
			RoleID:       uintPtr(2),
			SubVillageID: &genteng.ID,
			NIK:          ptr("3578014234567893"),
			Address:      ptr("Jl. Genteng No. 23, Surabaya"),
			IsMobile:     true,
		},
		{
			Name:         "Joko Widodo",
			Password:     "password123",
			BirthDate:    birthDate5,
			Telp:         ptr("081234567894"),
			Gender:       ptr("male"),
			Job:          ptr("Wiraswasta"),
			RoleID:       uintPtr(3),
			SubVillageID: &tegalsari.ID,
			NIK:          ptr("3578015234567894"),
			Address:      ptr("Jl. Tegalsari No. 10, Surabaya"),
			IsMobile:     false,
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
		log.Printf("  ✓ User: %s (Age: %d)\n", u.Name, u.GetAge())
	}

	log.Println("✅ All seeders completed successfully")
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