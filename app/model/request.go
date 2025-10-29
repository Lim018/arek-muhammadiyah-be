package model

import "time"

type LoginRequest struct {
	ID       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserRequest struct {
	ID        string     `json:"id" validate:"required"`
	Name      string     `json:"name" validate:"required"`
	Password  string     `json:"password" validate:"required,min=6"`
	BirthDate *time.Time `json:"birth_date"`
	Telp      *string    `json:"telp"`
	Gender    *string    `json:"gender"`
	Job       *string    `json:"job"`
	RoleID    *uint      `json:"role_id"`
	VillageID *string    `json:"village_id"`
	NIK       *string    `json:"nik"`
	Address   *string    `json:"address"`
	IsMobile  *bool      `json:"is_mobile"`
}

type UpdateUserRequest struct {
	Name      *string    `json:"name"`
	BirthDate *time.Time `json:"birth_date"`
	Telp      *string    `json:"telp"`
	Gender    *string    `json:"gender"`
	Job       *string    `json:"job"`
	RoleID    *uint      `json:"role_id"`
	VillageID *string    `json:"village_id"`
	NIK       *string    `json:"nik"`
	Address   *string    `json:"address"`
}

type CreateArticleRequest struct {
	CategoryID   *uint   `json:"category_id"`
	Title        string  `json:"title" validate:"required"`
	Content      string  `json:"content" validate:"required"`
	FeatureImage *string `json:"feature_image"`
	IsPublished  *bool   `json:"is_published"`
}

type CreateTicketRequest struct {
	CategoryID  *uint  `json:"category_id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateTicketRequest struct {
	Status     *TicketStatus `json:"status"`
	Resolution *string       `json:"resolution"`
}

type CreateDocumentRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	FilePath    string  `json:"file_path" validate:"required"`
	FileName    string  `json:"file_name" validate:"required"`
	FileSize    *int64  `json:"file_size"`
	MimeType    *string `json:"mime_type"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	IsActive    *bool   `json:"is_active"`
}