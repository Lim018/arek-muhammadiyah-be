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
	log.Println("🔑 Seeding roles...")
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
	log.Println("   Note: Village IDs based on Jawa Timur wilayah.json")

	users := []model.User{
		// === KOTA SURABAYA (3578) - 15 users ===
		{
			Name:      "Dr. Ahmad Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1985-01-15"),
			Telp:      ptr("081234000001"),
			Gender:    ptr("male"),
			Job:       ptr("Dokter Spesialis"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3578011001"), // Kelurahan Simokerto
			NIK:       ptr("3578011234567001"),
			Address:   ptr("Jl. Simokerto Raya No. 15"),
			IsMobile:  true,
		},
		{
			Name:      "Siti Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1990-03-20"),
			Telp:      ptr("081234000002"),
			Gender:    ptr("female"),
			Job:       ptr("Guru SD"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3578011002"),
			NIK:       ptr("3578011234567002"),
			Address:   ptr("Jl. Kapasan No. 22"),
			IsMobile:  true,
		},
		{
			Name:      "Budi Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1988-05-10"),
			Telp:      ptr("081234000003"),
			Gender:    ptr("male"),
			Job:       ptr("Pengusaha"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3578021001"), // Kelurahan Tegalsari
			NIK:       ptr("3578021234567003"),
			Address:   ptr("Jl. Tegalsari No. 45"),
			IsMobile:  false,
		},
		{
			Name:      "Dewi Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1992-07-25"),
			Telp:      ptr("081234000004"),
			Gender:    ptr("female"),
			Job:       ptr("Perawat"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3578031001"), // Kelurahan Genteng
			NIK:       ptr("3578031234567004"),
			Address:   ptr("Jl. Genteng Besar No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Andi Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1987-09-30"),
			Telp:      ptr("081234000005"),
			Gender:    ptr("male"),
			Job:       ptr("PNS"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3578041001"), // Kelurahan Bubutan
			NIK:       ptr("3578041234567005"),
			Address:   ptr("Jl. Bubutan No. 88"),
			IsMobile:  true,
		},
		{
			Name:      "Rina Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1995-02-14"),
			Telp:      ptr("081234000006"),
			Gender:    ptr("female"),
			Job:       ptr("Dosen"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3578051001"), // Kelurahan Gubeng
			NIK:       ptr("3578051234567006"),
			Address:   ptr("Jl. Gubeng Kertajaya No. 67"),
			IsMobile:  true,
		},
		{
			Name:      "Hadi Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1989-11-08"),
			Telp:      ptr("081234000007"),
			Gender:    ptr("male"),
			Job:       ptr("Arsitek"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3578061001"), // Kelurahan Wonokromo
			NIK:       ptr("3578061234567007"),
			Address:   ptr("Jl. Wonokromo No. 34"),
			IsMobile:  false,
		},
		{
			Name:      "Maya Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1993-04-18"),
			Telp:      ptr("081234000008"),
			Gender:    ptr("female"),
			Job:       ptr("Apoteker"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3578071001"), // Kelurahan Rungkut
			NIK:       ptr("3578071234567008"),
			Address:   ptr("Jl. Rungkut Madya No. 99"),
			IsMobile:  true,
		},
		{
			Name:      "Fajar Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1991-06-22"),
			Telp:      ptr("081234000009"),
			Gender:    ptr("male"),
			Job:       ptr("IT Consultant"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3578081001"), // Kelurahan Tenggilis
			NIK:       ptr("3578081234567009"),
			Address:   ptr("Jl. Tenggilis Mejoyo No. 56"),
			IsMobile:  true,
		},
		{
			Name:      "Indah Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1994-08-05"),
			Telp:      ptr("081234000010"),
			Gender:    ptr("female"),
			Job:       ptr("Akuntan"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3578091001"), // Kelurahan Mulyorejo
			NIK:       ptr("3578091234567010"),
			Address:   ptr("Jl. Mulyorejo Utara No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Rizki Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1986-12-12"),
			Telp:      ptr("081234000011"),
			Gender:    ptr("male"),
			Job:       ptr("Manager"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3578101001"), // Kelurahan Sukolilo
			NIK:       ptr("3578101234567011"),
			Address:   ptr("Jl. Keputih No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Lina Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1990-10-30"),
			Telp:      ptr("081234000012"),
			Gender:    ptr("female"),
			Job:       ptr("Psikolog"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3578111001"), // Kelurahan Sememi
			NIK:       ptr("3578111234567012"),
			Address:   ptr("Jl. Sememi No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Joko Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1988-03-17"),
			Telp:      ptr("081234000013"),
			Gender:    ptr("male"),
			Job:       ptr("Pilot"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3578121001"), // Kelurahan Kenjeran
			NIK:       ptr("3578121234567013"),
			Address:   ptr("Jl. Kenjeran No. 112"),
			IsMobile:  true,
		},
		{
			Name:      "Putri Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1996-05-09"),
			Telp:      ptr("081234000014"),
			Gender:    ptr("female"),
			Job:       ptr("Desainer"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3578131001"), // Kelurahan Pacar Keling
			NIK:       ptr("3578131234567014"),
			Address:   ptr("Jl. Pacar Keling No. 67"),
			IsMobile:  true,
		},
		{
			Name:      "Eko Surabaya",
			Password:  "password123",
			BirthDate: parseDate("1984-07-21"),
			Telp:      ptr("081234000015"),
			Gender:    ptr("male"),
			Job:       ptr("Insinyur"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3578141001"), // Kelurahan Sambikerep
			NIK:       ptr("3578141234567015"),
			Address:   ptr("Jl. Sambikerep No. 34"),
			IsMobile:  false,
		},

		// === KABUPATEN SIDOARJO (3515) - 8 users ===
		{
			Name:      "Ahmad Sidoarjo",
			Password:  "password123",
			BirthDate: parseDate("1987-02-10"),
			Telp:      ptr("081235000001"),
			Gender:    ptr("male"),
			Job:       ptr("Guru SMA"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3515011001"), // Kelurahan Sidoarjo
			NIK:       ptr("3515011234567001"),
			Address:   ptr("Jl. Pahlawan No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Siti Sidoarjo",
			Password:  "password123",
			BirthDate: parseDate("1991-04-15"),
			Telp:      ptr("081235000002"),
			Gender:    ptr("female"),
			Job:       ptr("Bidan"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3515011002"),
			NIK:       ptr("3515011234567002"),
			Address:   ptr("Jl. Diponegoro No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Budi Sidoarjo",
			Password:  "password123",
			BirthDate: parseDate("1989-06-20"),
			Telp:      ptr("081235000003"),
			Gender:    ptr("male"),
			Job:       ptr("Petani"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3515021001"), // Kecamatan Buduran
			NIK:       ptr("3515021234567003"),
			Address:   ptr("Jl. Raya Buduran No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Dewi Sidoarjo",
			Password:  "password123",
			BirthDate: parseDate("1993-08-25"),
			Telp:      ptr("081235000004"),
			Gender:    ptr("female"),
			Job:       ptr("Perawat"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3515031001"), // Kecamatan Candi
			NIK:       ptr("3515031234567004"),
			Address:   ptr("Jl. Candi No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Andi Sidoarjo",
			Password:  "password123",
			BirthDate: parseDate("1986-10-30"),
			Telp:      ptr("081235000005"),
			Gender:    ptr("male"),
			Job:       ptr("Wiraswasta"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3515041001"), // Kecamatan Porong
			NIK:       ptr("3515041234567005"),
			Address:   ptr("Jl. Porong No. 56"),
			IsMobile:  true,
		},
		{
			Name:      "Rina Sidoarjo",
			Password:  "password123",
			BirthDate: parseDate("1995-12-05"),
			Telp:      ptr("081235000006"),
			Gender:    ptr("female"),
			Job:       ptr("Mahasiswa"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3515051001"), // Kecamatan Krembung
			NIK:       ptr("3515051234567006"),
			Address:   ptr("Jl. Krembung No. 89"),
			IsMobile:  false,
		},
		{
			Name:      "Hadi Sidoarjo",
			Password:  "password123",
			BirthDate: parseDate("1988-01-18"),
			Telp:      ptr("081235000007"),
			Gender:    ptr("male"),
			Job:       ptr("Teknisi"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3515061001"), // Kecamatan Tanggulangin
			NIK:       ptr("3515061234567007"),
			Address:   ptr("Jl. Tanggulangin No. 34"),
			IsMobile:  true,
		},
		{
			Name:      "Maya Sidoarjo",
			Password:  "password123",
			BirthDate: parseDate("1992-03-22"),
			Telp:      ptr("081235000008"),
			Gender:    ptr("female"),
			Job:       ptr("Kasir"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3515071001"), // Kecamatan Jabon
			NIK:       ptr("3515071234567008"),
			Address:   ptr("Jl. Jabon No. 67"),
			IsMobile:  false,
		},

		// === KABUPATEN GRESIK (3525) - 6 users ===
		{
			Name:      "Ahmad Gresik",
			Password:  "password123",
			BirthDate: parseDate("1985-05-12"),
			Telp:      ptr("081236000001"),
			Gender:    ptr("male"),
			Job:       ptr("Nelayan"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3525011001"), // Kelurahan Gresik Kota
			NIK:       ptr("3525011234567001"),
			Address:   ptr("Jl. Veteran No. 15"),
			IsMobile:  true,
		},
		{
			Name:      "Siti Gresik",
			Password:  "password123",
			BirthDate: parseDate("1990-07-20"),
			Telp:      ptr("081236000002"),
			Gender:    ptr("female"),
			Job:       ptr("Pedagang"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3525011002"),
			NIK:       ptr("3525011234567002"),
			Address:   ptr("Jl. Ahmad Yani No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Budi Gresik",
			Password:  "password123",
			BirthDate: parseDate("1988-09-25"),
			Telp:      ptr("081236000003"),
			Gender:    ptr("male"),
			Job:       ptr("Buruh Pabrik"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3525021001"), // Kecamatan Kebomas
			NIK:       ptr("3525021234567003"),
			Address:   ptr("Jl. Kebomas No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Dewi Gresik",
			Password:  "password123",
			BirthDate: parseDate("1992-11-30"),
			Telp:      ptr("081236000004"),
			Gender:    ptr("female"),
			Job:       ptr("Guru TK"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3525031001"), // Kecamatan Manyar
			NIK:       ptr("3525031234567004"),
			Address:   ptr("Jl. Manyar No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Andi Gresik",
			Password:  "password123",
			BirthDate: parseDate("1987-02-08"),
			Telp:      ptr("081236000005"),
			Gender:    ptr("male"),
			Job:       ptr("Supir"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3525041001"), // Kecamatan Bungah
			NIK:       ptr("3525041234567005"),
			Address:   ptr("Jl. Bungah No. 56"),
			IsMobile:  true,
		},
		{
			Name:      "Rina Gresik",
			Password:  "password123",
			BirthDate: parseDate("1994-04-15"),
			Telp:      ptr("081236000006"),
			Gender:    ptr("female"),
			Job:       ptr("Karyawan Swasta"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3525051001"), // Kecamatan Ujungpangkah
			NIK:       ptr("3525051234567006"),
			Address:   ptr("Jl. Ujungpangkah No. 89"),
			IsMobile:  false,
		},

		// === KABUPATEN MOJOKERTO (3516) - 5 users ===
		{
			Name:      "Ahmad Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1986-03-18"),
			Telp:      ptr("081237000001"),
			Gender:    ptr("male"),
			Job:       ptr("Guru"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3516011001"), // Kecamatan Puri
			NIK:       ptr("3516011234567001"),
			Address:   ptr("Jl. Puri No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Siti Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1991-05-22"),
			Telp:      ptr("081237000002"),
			Gender:    ptr("female"),
			Job:       ptr("Perawat"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3516011002"),
			NIK:       ptr("3516011234567002"),
			Address:   ptr("Jl. Jetis No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Budi Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1989-07-27"),
			Telp:      ptr("081237000003"),
			Gender:    ptr("male"),
			Job:       ptr("Petani"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3516021001"), // Kecamatan Trawas
			NIK:       ptr("3516021234567003"),
			Address:   ptr("Jl. Trawas No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Dewi Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1993-09-30"),
			Telp:      ptr("081237000004"),
			Gender:    ptr("female"),
			Job:       ptr("Bidan"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3516031001"), // Kecamatan Pacet
			NIK:       ptr("3516031234567004"),
			Address:   ptr("Jl. Pacet No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Andi Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1987-11-12"),
			Telp:      ptr("081237000005"),
			Gender:    ptr("male"),
			Job:       ptr("Wiraswasta"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3516041001"), // Kecamatan Ngoro
			NIK:       ptr("3516041234567005"),
			Address:   ptr("Jl. Ngoro No. 56"),
			IsMobile:  true,
		},

		// === KOTA MOJOKERTO (3576) - 5 users ===
		{
			Name:      "Fajar Kota Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1990-01-10"),
			Telp:      ptr("081238000001"),
			Gender:    ptr("male"),
			Job:       ptr("PNS"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3576011001"), // Kelurahan Blooto
			NIK:       ptr("3576011234567001"),
			Address:   ptr("Jl. Merdeka No. 15, Blooto"),
			IsMobile:  true,
		},
		{
			Name:      "Indah Kota Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1992-03-15"),
			Telp:      ptr("081238000002"),
			Gender:    ptr("female"),
			Job:       ptr("Dokter"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3576011002"), // Kelurahan Kauman
			NIK:       ptr("3576012234567002"),
			Address:   ptr("Jl. Kauman No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Rizki Kota Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1988-05-20"),
			Telp:      ptr("081238000003"),
			Gender:    ptr("male"),
			Job:       ptr("Pengusaha"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3576021001"), // Kelurahan Balongsari
			NIK:       ptr("3576013234567003"),
			Address:   ptr("Jl. Balongsari No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Lina Kota Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1994-07-25"),
			Telp:      ptr("081238000004"),
			Gender:    ptr("female"),
			Job:       ptr("Apoteker"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3576021002"), // Kelurahan Gedongan
			NIK:       ptr("3576014234567004"),
			Address:   ptr("Jl. Gedongan No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Joko Kota Mojokerto",
			Password:  "password123",
			BirthDate: parseDate("1986-09-30"),
			Telp:      ptr("081238000005"),
			Gender:    ptr("male"),
			Job:       ptr("Manager"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3576031001"), // Kelurahan Jagalan
			NIK:       ptr("3576015234567005"),
			Address:   ptr("Jl. Jagalan No. 10"),
			IsMobile:  true,
		},

		// === KABUPATEN MALANG (3507) - 7 users ===
		{
			Name:      "Ahmad Malang",
			Password:  "password123",
			BirthDate: parseDate("1985-04-12"),
			Telp:      ptr("081239000001"),
			Gender:    ptr("male"),
			Job:       ptr("Dosen"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3507011001"), // Kecamatan Donomulyo
			NIK:       ptr("3507011234567001"),
			Address:   ptr("Jl. Donomulyo No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Siti Malang",
			Password:  "password123",
			BirthDate: parseDate("1990-06-17"),
			Telp:      ptr("081239000002"),
			Gender:    ptr("female"),
			Job:       ptr("Guru"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3507011002"),
			NIK:       ptr("3507011234567002"),
			Address:   ptr("Jl. Kalipare No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Budi Malang",
			Password:  "password123",
			BirthDate: parseDate("1988-08-22"),
			Telp:      ptr("081239000003"),
			Gender:    ptr("male"),
			Job:       ptr("Petani"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3507021001"), // Kecamatan Dampit
			NIK:       ptr("3507021234567003"),
			Address:   ptr("Jl. Dampit No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Dewi Malang",
			Password:  "password123",
			BirthDate: parseDate("1992-10-27"),
			Telp:      ptr("081239000004"),
			Gender:    ptr("female"),
			Job:       ptr("Perawat"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3507031001"), // Kecamatan Tirtoyudo
			NIK:       ptr("3507031234567004"),
			Address:   ptr("Jl. Tirtoyudo No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Andi Malang",
			Password:  "password123",
			BirthDate: parseDate("1987-12-30"),
			Telp:      ptr("081239000005"),
			Gender:    ptr("male"),
			Job:       ptr("Wiraswasta"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3507041001"), // Kecamatan Ampelgading
			NIK:       ptr("3507041234567005"),
			Address:   ptr("Jl. Ampelgading No. 56"),
			IsMobile:  true,
		},
		{
			Name:      "Rina Malang",
			Password:  "password123",
			BirthDate: parseDate("1995-02-05"),
			Telp:      ptr("081239000006"),
			Gender:    ptr("female"),
			Job:       ptr("Mahasiswa"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3507051001"), // Kecamatan Poncokusumo
			NIK:       ptr("3507051234567006"),
			Address:   ptr("Jl. Poncokusumo No. 89"),
			IsMobile:  false,
		},
		{
			Name:      "Hadi Malang",
			Password:  "password123",
			BirthDate: parseDate("1989-04-10"),
			Telp:      ptr("081239000007"),
			Gender:    ptr("male"),
			Job:       ptr("Teknisi"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3507061001"), // Kecamatan Wajak
			NIK:       ptr("3507061234567007"),
			Address:   ptr("Jl. Wajak No. 34"),
			IsMobile:  true,
		},

		// === KOTA MALANG (3573) - 6 users ===
		{
			Name:      "Fajar Kota Malang",
			Password:  "password123",
			BirthDate: parseDate("1986-01-15"),
			Telp:      ptr("081240000001"),
			Gender:    ptr("male"),
			Job:       ptr("Pengacara"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3573011001"), // Kelurahan Lowokwaru
			NIK:       ptr("3573011234567001"),
			Address:   ptr("Jl. Soekarno Hatta No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Indah Kota Malang",
			Password:  "password123",
			BirthDate: parseDate("1991-03-20"),
			Telp:      ptr("081240000002"),
			Gender:    ptr("female"),
			Job:       ptr("Dokter Gigi"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3573011002"),
			NIK:       ptr("3573011234567002"),
			Address:   ptr("Jl. Veteran No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Rizki Kota Malang",
			Password:  "password123",
			BirthDate: parseDate("1989-05-25"),
			Telp:      ptr("081240000003"),
			Gender:    ptr("male"),
			Job:       ptr("Arsitek"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3573021001"), // Kelurahan Klojen
			NIK:       ptr("3573021234567003"),
			Address:   ptr("Jl. Ijen No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Lina Kota Malang",
			Password:  "password123",
			BirthDate: parseDate("1993-07-30"),
			Telp:      ptr("081240000004"),
			Gender:    ptr("female"),
			Job:       ptr("Psikolog"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3573031001"), // Kelurahan Blimbing
			NIK:       ptr("3573031234567004"),
			Address:   ptr("Jl. Blimbing No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Joko Kota Malang",
			Password:  "password123",
			BirthDate: parseDate("1987-09-12"),
			Telp:      ptr("081240000005"),
			Gender:    ptr("male"),
			Job:       ptr("Marketing Manager"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3573041001"), // Kelurahan Kedungkandang
			NIK:       ptr("3573041234567005"),
			Address:   ptr("Jl. Kedungkandang No. 56"),
			IsMobile:  true,
		},
		{
			Name:      "Putri Kota Malang",
			Password:  "password123",
			BirthDate: parseDate("1995-11-18"),
			Telp:      ptr("081240000006"),
			Gender:    ptr("female"),
			Job:       ptr("Content Creator"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3573051001"), // Kelurahan Sukun
			NIK:       ptr("3573051234567006"),
			Address:   ptr("Jl. Sukun No. 89"),
			IsMobile:  false,
		},

		// === KABUPATEN JEMBER (3509) - 5 users ===
		{
			Name:      "Ahmad Jember",
			Password:  "password123",
			BirthDate: parseDate("1985-02-08"),
			Telp:      ptr("081241000001"),
			Gender:    ptr("male"),
			Job:       ptr("Petani Tembakau"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3509011001"), // Kecamatan Kencong
			NIK:       ptr("3509011234567001"),
			Address:   ptr("Jl. Kencong No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Siti Jember",
			Password:  "password123",
			BirthDate: parseDate("1990-04-13"),
			Telp:      ptr("081241000002"),
			Gender:    ptr("female"),
			Job:       ptr("Guru SD"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3509011002"),
			NIK:       ptr("3509011234567002"),
			Address:   ptr("Jl. Gumukmas No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Budi Jember",
			Password:  "password123",
			BirthDate: parseDate("1988-06-18"),
			Telp:      ptr("081241000003"),
			Gender:    ptr("male"),
			Job:       ptr("Sopir Bus"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3509021001"), // Kecamatan Puger
			NIK:       ptr("3509021234567003"),
			Address:   ptr("Jl. Puger No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Dewi Jember",
			Password:  "password123",
			BirthDate: parseDate("1992-08-23"),
			Telp:      ptr("081241000004"),
			Gender:    ptr("female"),
			Job:       ptr("Bidan Desa"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3509031001"), // Kecamatan Wuluhan
			NIK:       ptr("3509031234567004"),
			Address:   ptr("Jl. Wuluhan No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Andi Jember",
			Password:  "password123",
			BirthDate: parseDate("1987-10-28"),
			Telp:      ptr("081241000005"),
			Gender:    ptr("male"),
			Job:       ptr("Pedagang"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3509041001"), // Kecamatan Ambulu
			NIK:       ptr("3509041234567005"),
			Address:   ptr("Jl. Ambulu No. 56"),
			IsMobile:  true,
		},

		// === KABUPATEN BANYUWANGI (3510) - 5 users ===
		{
			Name:      "Fajar Banyuwangi",
			Password:  "password123",
			BirthDate: parseDate("1986-03-05"),
			Telp:      ptr("081242000001"),
			Gender:    ptr("male"),
			Job:       ptr("Tour Guide"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3510011001"), // Kecamatan Banyuwangi
			NIK:       ptr("3510011234567001"),
			Address:   ptr("Jl. Ahmad Yani No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Indah Banyuwangi",
			Password:  "password123",
			BirthDate: parseDate("1991-05-10"),
			Telp:      ptr("081242000002"),
			Gender:    ptr("female"),
			Job:       ptr("Perawat"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3510011002"),
			NIK:       ptr("3510011234567002"),
			Address:   ptr("Jl. Diponegoro No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Rizki Banyuwangi",
			Password:  "password123",
			BirthDate: parseDate("1989-07-15"),
			Telp:      ptr("081242000003"),
			Gender:    ptr("male"),
			Job:       ptr("Nelayan"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3510021001"), // Kecamatan Rogojampi
			NIK:       ptr("3510021234567003"),
			Address:   ptr("Jl. Rogojampi No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Lina Banyuwangi",
			Password:  "password123",
			BirthDate: parseDate("1993-09-20"),
			Telp:      ptr("081242000004"),
			Gender:    ptr("female"),
			Job:       ptr("Guru SMP"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3510031001"), // Kecamatan Kabat
			NIK:       ptr("3510031234567004"),
			Address:   ptr("Jl. Kabat No. 23"),
			IsMobile:  true,
		},
		{
			Name:      "Joko Banyuwangi",
			Password:  "password123",
			BirthDate: parseDate("1987-11-25"),
			Telp:      ptr("081242000005"),
			Gender:    ptr("male"),
			Job:       ptr("Petani Kopi"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3510041001"), // Kecamatan Kalibaru
			NIK:       ptr("3510041234567005"),
			Address:   ptr("Jl. Kalibaru No. 56"),
			IsMobile:  true,
		},

		// === KABUPATEN KEDIRI (3506) - 4 users ===
		{
			Name:      "Ahmad Kediri",
			Password:  "password123",
			BirthDate: parseDate("1985-01-20"),
			Telp:      ptr("081243000001"),
			Gender:    ptr("male"),
			Job:       ptr("Guru SMA"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3506011001"), // Kecamatan Mojo
			NIK:       ptr("3506011234567001"),
			Address:   ptr("Jl. Mojo No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Siti Kediri",
			Password:  "password123",
			BirthDate: parseDate("1990-03-25"),
			Telp:      ptr("081243000002"),
			Gender:    ptr("female"),
			Job:       ptr("Perawat"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3506011002"),
			NIK:       ptr("3506011234567002"),
			Address:   ptr("Jl. Semen No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Budi Kediri",
			Password:  "password123",
			BirthDate: parseDate("1988-05-30"),
			Telp:      ptr("081243000003"),
			Gender:    ptr("male"),
			Job:       ptr("Buruh Pabrik Rokok"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3506021001"), // Kecamatan Gurah
			NIK:       ptr("3506021234567003"),
			Address:   ptr("Jl. Gurah No. 78"),
			IsMobile:  false,
		},
		{
			Name:      "Dewi Kediri",
			Password:  "password123",
			BirthDate: parseDate("1992-07-12"),
			Telp:      ptr("081243000004"),
			Gender:    ptr("female"),
			Job:       ptr("Bidan"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3506031001"), // Kecamatan Pare
			NIK:       ptr("3506031234567004"),
			Address:   ptr("Jl. Pare No. 23"),
			IsMobile:  true,
		},

		// === KOTA KEDIRI (3571) - 3 users ===
		{
			Name:      "Fajar Kota Kediri",
			Password:  "password123",
			BirthDate: parseDate("1986-02-14"),
			Telp:      ptr("081244000001"),
			Gender:    ptr("male"),
			Job:       ptr("PNS"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3571011001"), // Kelurahan Pesantren
			NIK:       ptr("3571011234567001"),
			Address:   ptr("Jl. Dhoho No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Indah Kota Kediri",
			Password:  "password123",
			BirthDate: parseDate("1991-04-19"),
			Telp:      ptr("081244000002"),
			Gender:    ptr("female"),
			Job:       ptr("Dokter"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3571011002"),
			NIK:       ptr("3571011234567002"),
			Address:   ptr("Jl. Brawijaya No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Rizki Kota Kediri",
			Password:  "password123",
			BirthDate: parseDate("1989-06-24"),
			Telp:      ptr("081244000003"),
			Gender:    ptr("male"),
			Job:       ptr("Pengusaha"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3571021001"), // Kelurahan Mojoroto
			NIK:       ptr("3571021234567003"),
			Address:   ptr("Jl. Mojoroto No. 78"),
			IsMobile:  false,
		},

		// === KABUPATEN PROBOLINGGO (3513) - 3 users ===
		{
			Name:      "Ahmad Probolinggo",
			Password:  "password123",
			BirthDate: parseDate("1985-05-08"),
			Telp:      ptr("081245000001"),
			Gender:    ptr("male"),
			Job:       ptr("Petani"),
			RoleID:    uintPtr(1),
			VillageID: ptr("3513011001"), // Kecamatan Sukapura
			NIK:       ptr("3513011234567001"),
			Address:   ptr("Jl. Sukapura No. 12"),
			IsMobile:  true,
		},
		{
			Name:      "Siti Probolinggo",
			Password:  "password123",
			BirthDate: parseDate("1990-07-13"),
			Telp:      ptr("081245000002"),
			Gender:    ptr("female"),
			Job:       ptr("Guru"),
			RoleID:    uintPtr(2),
			VillageID: ptr("3513011002"),
			NIK:       ptr("3513011234567002"),
			Address:   ptr("Jl. Leces No. 45"),
			IsMobile:  true,
		},
		{
			Name:      "Budi Probolinggo",
			Password:  "password123",
			BirthDate: parseDate("1988-09-18"),
			Telp:      ptr("081245000003"),
			Gender:    ptr("male"),
			Job:       ptr("Nelayan"),
			RoleID:    uintPtr(3),
			VillageID: ptr("3513021001"), // Kecamatan Paiton
			NIK:       ptr("3513021234567003"),
			Address:   ptr("Jl. Paiton No. 78"),
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
		log.Printf("  ✓ User: %s (Age: %d, City: %s, Village ID: %s)\n", 
			u.Name, u.GetAge(), (*u.VillageID)[:4], *u.VillageID)
	}

	log.Println("✅ All seeders completed successfully")
	log.Printf("📊 Total users seeded: %d\n", len(users))
	log.Println("🗺️ Coverage:")
	log.Println("   - Kota Surabaya (3578): 15 users")
	log.Println("   - Kabupaten Sidoarjo (3515): 8 users")
	log.Println("   - Kabupaten Gresik (3525): 6 users")
	log.Println("   - Kabupaten Mojokerto (3516): 5 users")
	log.Println("   - Kota Mojokerto (3576): 5 users (merged with 3516 on map)")
	log.Println("   - Kabupaten Malang (3507): 7 users")
	log.Println("   - Kota Malang (3573): 6 users")
	log.Println("   - Kabupaten Jember (3509): 5 users")
	log.Println("   - Kabupaten Banyuwangi (3510): 5 users")
	log.Println("   - Kabupaten Kediri (3506): 4 users")
	log.Println("   - Kota Kediri (3571): 3 users")
	log.Println("   - Kabupaten Probolinggo (3513): 3 users")
	log.Println("📝 Make sure your wilayah.json contains the village IDs used above")
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