package model

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
}

type UserWithStats struct {
	User
	TotalArticles int `json:"total_articles"`
	TotalTickets  int `json:"total_tickets"`
	TotalDocuments int `json:"total_documents"`
}

type TicketStats struct {
	Unread     int64 `json:"unread"`
	Read       int64 `json:"read"`
	InProgress int64 `json:"in_progress"`
	Resolved   int64 `json:"resolved"`
	Closed     int64 `json:"closed"`
	Total      int64 `json:"total"`
}

type DashboardStats struct {
	TotalUsers     int64        `json:"total_users"`
	TotalArticles  int64        `json:"total_articles"`
	TotalTickets   int64        `json:"total_tickets"`
	TotalVillages  int64        `json:"total_villages"`
	TicketStats    TicketStats  `json:"ticket_stats"`
	CardStatusStats map[string]int64 `json:"card_status_stats"`
}

type VillageWithUserCount struct {
	Village
	TotalUsers int `json:"total_users"`
}