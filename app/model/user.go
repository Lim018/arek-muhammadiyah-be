package model

import (
	"time"
)

type User struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"not null"`
	Password  string     `json:"-" gorm:"not null"`
	RoleID    *uint      `json:"role_id"`
	VillageID *uint      `json:"village_id"`
	NIK       *string    `json:"nik"`
	Address   *string    `json:"address"`
	Photo     *string    `json:"photo"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	
	// Relations
	Role     *Role     `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Village  *Village  `json:"village,omitempty" gorm:"foreignKey:VillageID"`
	Articles []Article `json:"articles,omitempty"`
	Tickets  []Ticket  `json:"tickets,omitempty"`
	Documents []Document `json:"documents,omitempty"`
}

type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Users     []User     `json:"users,omitempty"`
	RoleMenus []RoleMenu `json:"role_menus,omitempty"`
}

type Village struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	Name        string      `json:"name" gorm:"not null"`
	Code        *string     `json:"code" gorm:"unique"`
	Description *string     `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	
	// Relations
	Users       []User       `json:"users,omitempty"`
	SubVillages []SubVillage `json:"sub_villages,omitempty"`
}

type SubVillage struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	VillageID   uint      `json:"village_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Code        *string   `json:"code"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Village *Village `json:"village,omitempty" gorm:"foreignKey:VillageID"`
}