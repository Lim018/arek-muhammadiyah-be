package database

import (
	"log"
	"arek-muhammadiyah-be/app/model"
	"golang.org/x/crypto/bcrypt"
)


func RunSeeders() {
	if DB == nil {
		log.Fatal("Database not connected")
	}


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

	// === Categories ===
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
	}

	// === Villages ===
	villages := []model.Village{
		{Name: "Gubeng", Code: "SBY-GBG-001", Description: ptr("Kelurahan Gubeng, Kecamatan Gubeng"), Color: "#3B82F6", IsActive: true},
		{Name: "Airlangga", Code: "SBY-ARL-002", Description: ptr("Kelurahan Airlangga, Kecamatan Gubeng"), Color: "#8B5CF6", IsActive: true},
		{Name: "Wonokromo", Code: "SBY-WNK-003", Description: ptr("Kelurahan Wonokromo, Kecamatan Wonokromo"), Color: "#10B981", IsActive: true},
		{Name: "Sawahan", Code: "SBY-SWH-004", Description: ptr("Kelurahan Sawahan, Kecamatan Sawahan"), Color: "#F59E0B", IsActive: true},
		{Name: "Genteng", Code: "SBY-GTG-005", Description: ptr("Kelurahan Genteng, Kecamatan Genteng"), Color: "#EF4444", IsActive: true},
	}
	for _, v := range villages {
		DB.FirstOrCreate(&v, model.Village{Code: v.Code})
	}

	// === Users ===
		// === Users ===
		users := []model.User{
			{Name: "Budi Santoso", Password: "password123", Telp: ptr("081234567890"), RoleID: uintPtr(1), VillageID: uintPtr(1), NIK: ptr("3578011234567890"), Address: ptr("Jl. Gubeng Pojok No. 15, Surabaya"), CardStatus: "delivered"},
			{Name: "Siti Nurhaliza", Password: "password123", Telp: ptr("081234567891"), RoleID: uintPtr(1), VillageID: uintPtr(2), NIK: ptr("3578012234567891"), Address: ptr("Jl. Airlangga No. 45, Surabaya"), CardStatus: "delivered"},
			{Name: "Ahmad Fauzi", Password: "password123", Telp: ptr("081234567892"), RoleID: uintPtr(2), VillageID: uintPtr(3), NIK: ptr("3578013234567892"), Address: ptr("Jl. Wonokromo No. 78, Surabaya"), CardStatus: "approved"},
			{Name: "Dewi Lestari", Password: "password123", Telp: ptr("081234567893"), RoleID: uintPtr(2), VillageID: uintPtr(4), NIK: ptr("3578014234567893"), Address: ptr("Jl. Sawahan No. 23, Surabaya"), CardStatus: "printed"},
			{Name: "Joko Widodo", Password: "password123", Telp: ptr("081234567894"), RoleID: uintPtr(3), VillageID: uintPtr(1), NIK: ptr("3578015234567894"), Address: ptr("Jl. Gubeng Kertajaya No. 10, Surabaya"), CardStatus: "delivered"},
		}
	
		for _, u := range users {
			// hash password dulu sebelum simpan
			hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
			if err != nil {
				log.Printf("Gagal hash password untuk user %s: %v", u.Name, err)
				continue
			}
			u.Password = string(hashed)
	
			DB.FirstOrCreate(&u, model.User{NIK: u.NIK})
		}

	log.Println("All seeders completed successfully")
}

// Helper functions
func ptr(s string) *string { return &s }
func uintPtr(u uint) *uint { return &u }
