package model

import "time"

type User struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name" gorm:"not null"`
	Password   string     `json:"-" gorm:"not null"`
	Telp       *string    `json:"telp"`
	RoleID     *uint      `json:"role_id"`
	VillageID  *uint      `json:"village_id"`
	NIK        *string    `json:"nik" gorm:"unique"`
	Address    *string    `json:"address"`
	CardStatus string     `json:"card_status" gorm:"default:'pending'"`
	IsMobile   bool       `json:"is_mobile" gorm:"default:false"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	
	// Relations
	Role      *Role      `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Village   *Village   `json:"village,omitempty" gorm:"foreignKey:VillageID"`
	Articles  []Article  `json:"articles,omitempty"`
	Tickets   []Ticket   `json:"tickets,omitempty"`
	Documents []Document `json:"documents,omitempty"`
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
	Name        string    `json:"name" gorm:"not null"`
	Code        string    `json:"code" gorm:"unique;not null"`
	Description *string   `json:"description"`
	Color       string    `json:"color" gorm:"default:'#3B82F6'"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Users []User `json:"users,omitempty"`
}