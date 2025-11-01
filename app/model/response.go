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
	Age            int `json:"age"`
	TotalArticles  int `json:"total_articles"`
	TotalTickets   int `json:"total_tickets"`
	TotalDocuments int `json:"total_documents"`
}

type TicketStats struct {
	Unread     int64 `json:"unread"`
	Read       int64 `json:"read"`
	InProgress int64 `json:"in_progress"`
	Resolved   int64 `json:"resolved"`
	Rejected   int64 `json:"rejected"`
	Total      int64 `json:"total"`
}

type DashboardStats struct {
	TotalUsers    int64                  `json:"total_users"`
	TotalArticles int64                  `json:"total_articles"`
	TotalTickets  int64                  `json:"total_tickets"`
	TicketStats   TicketStats            `json:"ticket_stats"`
	WilayahStats  map[string]interface{} `json:"wilayah_stats"`
	GenderStats   map[string]int64       `json:"gender_stats"`
}

// Wilayah statistics
type CityStats struct {
	CityID       string `json:"city_id"`
	CityName     string `json:"city_name"`
	TotalMembers int64  `json:"total_members"`
	TotalMale    int64  `json:"total_male"`
	TotalFemale  int64  `json:"total_female"`
	TotalMobile  int64  `json:"total_mobile"`
}

type DistrictStats struct {
	DistrictID   string `json:"district_id"`
	DistrictName string `json:"district_name"`
	CityID       string `json:"city_id"`
	CityName     string `json:"city_name"`
	TotalMembers int64  `json:"total_members"`
	TotalMale    int64  `json:"total_male"`
	TotalFemale  int64  `json:"total_female"`
}

type VillageStats struct {
	VillageID    string `json:"village_id"`
	VillageName  string `json:"village_name"`
	DistrictID   string `json:"district_id"`
	DistrictName string `json:"district_name"`
	CityID       string `json:"city_id"`
	CityName     string `json:"city_name"`
	TotalMembers int64  `json:"total_members"`
}

// Wilayah info response (enriched from JSON)
type WilayahInfo struct {
	VillageID    string `json:"village_id"`
	VillageName  string `json:"village_name"`
	DistrictID   string `json:"district_id"`
	DistrictName string `json:"district_name"`
	CityID       string `json:"city_id"`
	CityName     string `json:"city_name"`
}
