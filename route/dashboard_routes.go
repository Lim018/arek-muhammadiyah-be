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
		villageRepo := repository.NewVillageRepository()

		// Get totals
		_, totalUsers, _ := userRepo.GetAll(1, 0)
		_, totalArticles, _ := articleRepo.GetAll(1, 0, nil)
		_, totalTickets, _ := ticketRepo.GetAll(1, 0, nil)
		villages, _, _ := villageRepo.GetAll(100, 0, true)

		// Get ticket stats
		ticketStatusCounts, _ := ticketRepo.GetCountByStatus()
		// cardStatusStats, _ := userRepo.GetCardStatusStats()

		// Calculate total tickets
		totalTicketsSum := int64(0)
		for _, count := range ticketStatusCounts {
			totalTicketsSum += count
		}

		stats := model.DashboardStats{
			TotalUsers:    totalUsers,
			TotalArticles: totalArticles,
			TotalTickets:  totalTickets,
			TotalVillages: int64(len(villages)),
			TicketStats: model.TicketStats{
				Unread:     ticketStatusCounts[model.TicketStatusUnread],
				Read:       ticketStatusCounts[model.TicketStatusRead],
				InProgress: ticketStatusCounts[model.TicketStatusInProgress],
				Resolved:   ticketStatusCounts[model.TicketStatusResolved],
				Closed:     ticketStatusCounts[model.TicketStatusClosed],
				Total:      totalTicketsSum,
			},
			// CardStatusStats: cardStatusStats,
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Dashboard statistics retrieved successfully",
			Data:    stats,
		})
	})
}