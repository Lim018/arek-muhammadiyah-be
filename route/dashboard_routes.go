package route

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupDashboardRoutes(app *fiber.App) {
	dashboard := app.Group("/api/dashboard", middleware.Authorization(), middleware.AdminOnly())

	dashboard.Get("/stats", func(c *fiber.Ctx) error {
		userRepo := repository.NewUserRepository()
		articleRepo := repository.NewArticleRepository()
		ticketRepo := repository.NewTicketRepository()

		// Get totals
		_, totalUsers, _ := userRepo.GetAll(1, 0)
		_, totalArticles, _ := articleRepo.GetAll(1, 0, nil)
		_, totalTickets, _ := ticketRepo.GetAll(1, 0, nil)

		// Get ticket stats
		ticketStatusCounts, _ := ticketRepo.GetCountByStatus()

		// Calculate total tickets
		totalTicketsSum := int64(0)
		for _, count := range ticketStatusCounts {
			totalTicketsSum += count
		}

		// Get wilayah stats (grouped by city)
		wilayahStats, _ := userRepo.GetWilayahStats()

		// Get gender stats
		genderStats := make(map[string]int64)
		_, totalMale, _ := userRepo.GetByGender("male", 1, 0)
		_, totalFemale, _ := userRepo.GetByGender("female", 1, 0)
		genderStats["male"] = totalMale
		genderStats["female"] = totalFemale

		stats := model.DashboardStats{
			TotalUsers:    totalUsers,
			TotalArticles: totalArticles,
			TotalTickets:  totalTickets,
			TicketStats: model.TicketStats{
				Unread:     ticketStatusCounts[model.TicketStatusUnread],
				Read:       ticketStatusCounts[model.TicketStatusRead],
				InProgress: ticketStatusCounts[model.TicketStatusInProgress],
				Resolved:   ticketStatusCounts[model.TicketStatusResolved],
				Rejected:     ticketStatusCounts[model.TicketStatusRejected],
				Total:      totalTicketsSum,
			},
			WilayahStats: wilayahStats,
			GenderStats:  genderStats,
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Dashboard statistics retrieved successfully",
			Data:    stats,
		})
	})
}