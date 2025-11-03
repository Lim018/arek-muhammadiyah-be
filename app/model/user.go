package model

import "time"

type User struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name" gorm:"not null"`
	Password   string     `json:"-" gorm:"not null"`
	BirthDate  *time.Time `json:"birth_date"`
	Telp       *string    `json:"telp"`
	Gender     *string    `json:"gender" gorm:"type:varchar(20)"` // male, female
	Job        *string    `json:"job"`
	RoleID     *uint      `json:"role_id"`
	
	// Wilayah - hanya simpan Village ID (Kelurahan)
	VillageID  *string    `json:"village_id" gorm:"type:varchar(20);index"` // ID dari JSON (contoh: "3576011001")
	
	NIK        *string    `json:"nik" gorm:"unique"`
	Address    *string    `json:"address"`
	IsMobile   bool       `json:"is_mobile" gorm:"default:false"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	// Relations
	Role     *Role     `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Articles []Article `json:"articles,omitempty"`
	Tickets  []Ticket  `json:"tickets,omitempty"`
	// Documents []Document `json:"documents,omitempty"`
	
	// Virtual fields - akan diisi dari JSON saat query
	VillageName  string `json:"village_name,omitempty" gorm:"-"`
	DistrictID   string `json:"district_id,omitempty" gorm:"-"`
	DistrictName string `json:"district_name,omitempty" gorm:"-"`
	CityID       string `json:"city_id,omitempty" gorm:"-"`
	CityName     string `json:"city_name,omitempty" gorm:"-"`
}

// Method to calculate age
func (u *User) GetAge() int {
	if u.BirthDate == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - u.BirthDate.Year()
	if now.YearDay() < u.BirthDate.YearDay() {
		age--
	}
	return age
}

type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Users []User `json:"users,omitempty"`
}

// Struct untuk parse JSON wilayah Indonesia
type City struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Districts []District `json:"districts"`
}

type District struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Villages []Village `json:"villages"`
}

type Village struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}