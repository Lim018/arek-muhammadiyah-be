package model

import "time"

type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description *string   `json:"description"`
	Color       string    `json:"color" gorm:"default:'#10B981'"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Articles []Article `json:"articles,omitempty"`
	Tickets  []Ticket  `json:"tickets,omitempty"`
}

type Article struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        string    `json:"user_id" gorm:"not null"`
	CategoryID    *uint     `json:"category_id"`
	Title         string    `json:"title" gorm:"not null"`
	Slug          string    `json:"slug" gorm:"unique;not null"`
	Content       string    `json:"content" gorm:"not null"`
	FeatureImage  *string   `json:"feature_image"`
	IsPublished   bool      `json:"is_published" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// Relations
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type TicketStatus string

const (
	TicketStatusUnread     TicketStatus = "unread"
	TicketStatusRead       TicketStatus = "read"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusResolved   TicketStatus = "resolved"
	TicketStatusClosed     TicketStatus = "closed"
)

type Ticket struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	UserID      string       `json:"user_id" gorm:"not null"`
	CategoryID  *uint        `json:"category_id"`
	Title       string       `json:"title" gorm:"not null"`
	Description string       `json:"description" gorm:"not null"`
	Status      TicketStatus `json:"status" gorm:"default:'unread'"`
	Resolution  *string      `json:"resolution"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	ResolvedAt  *time.Time   `json:"resolved_at"`
	
	// Relations
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type Document struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	Description *string   `json:"description"`
	FilePath    string    `json:"file_path" gorm:"not null"`
	FileName    string    `json:"file_name" gorm:"not null"`
	FileSize    *int64    `json:"file_size"`
	MimeType    *string   `json:"mime_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}