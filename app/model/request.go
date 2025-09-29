package model

type LoginRequest struct {
	ID       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserRequest struct {
	ID        string  `json:"id" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Password  string  `json:"password" validate:"required,min=6"`
	RoleID    *uint   `json:"role_id"`
	VillageID *uint   `json:"village_id"`
	NIK       *string `json:"nik"`
	Address   *string `json:"address"`
}

type UpdateUserRequest struct {
	Name      *string `json:"name"`
	RoleID    *uint   `json:"role_id"`
	VillageID *uint   `json:"village_id"`
	NIK       *string `json:"nik"`
	Address   *string `json:"address"`
	Photo     *string `json:"photo"`
}

type CreateArticleRequest struct {
	CategoryID    *uint   `json:"category_id"`
	Title         string  `json:"title" validate:"required"`
	Content       string  `json:"content" validate:"required"`
	FeaturedImage *string `json:"featured_image"`
	IsPublished   *bool   `json:"is_published"`
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