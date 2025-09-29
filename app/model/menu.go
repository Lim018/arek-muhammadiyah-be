package model

import "time"

type Menu struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description *string   `json:"description"`
	Icon        *string   `json:"icon"`
	URL         *string   `json:"url"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	SubMenus  []SubMenu  `json:"sub_menus,omitempty"`
	RoleMenus []RoleMenu `json:"role_menus,omitempty"`
}

type SubMenu struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	MenuID      uint      `json:"menu_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description *string   `json:"description"`
	URL         *string   `json:"url"`
	Icon        *string   `json:"icon"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Menu *Menu `json:"menu,omitempty" gorm:"foreignKey:MenuID"`
}

type RoleMenu struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RoleID    uint      `json:"role_id" gorm:"not null"`
	MenuID    uint      `json:"menu_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relations
	Role *Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Menu *Menu `json:"menu,omitempty" gorm:"foreignKey:MenuID"`
}