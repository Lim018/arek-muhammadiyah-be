package model

import "time"

type User struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Name         string     `json:"name" gorm:"not null"`
	Password     string     `json:"-" gorm:"not null"`
	BirthDate    *time.Time `json:"birth_date"`
	Telp         *string    `json:"telp"`
	Gender       *string    `json:"gender" gorm:"type:varchar(20)"` // male, female
	Job          *string    `json:"job"`
	RoleID       *uint      `json:"role_id"`
	SubVillageID *uint      `json:"sub_village_id"`
	NIK          *string    `json:"nik" gorm:"unique"`
	Address      *string    `json:"address"`
	IsMobile     bool       `json:"is_mobile" gorm:"default:false"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	Role       *Role       `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	SubVillage *SubVillage `json:"sub_village,omitempty" gorm:"foreignKey:SubVillageID"`
	Articles   []Article   `json:"articles,omitempty"`
	Tickets    []Ticket    `json:"tickets,omitempty"`
	Documents  []Document  `json:"documents,omitempty"`
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

type Village struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"` // Kabupaten
	Code        string    `json:"code" gorm:"unique;not null"`
	Description *string   `json:"description"`
	Color       string    `json:"color" gorm:"default:'#3B82F6'"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	SubVillages []SubVillage `json:"sub_villages,omitempty"`
}

type SubVillage struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	VillageID   uint      `json:"village_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"` // Kecamatan
	Code        string    `json:"code" gorm:"unique;not null"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Village *Village `json:"village,omitempty" gorm:"foreignKey:VillageID"`
	Users   []User   `json:"users,omitempty"`
}