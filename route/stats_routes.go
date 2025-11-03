package route

import (
	"arek-muhammadiyah-be/app/service"
	"arek-muhammadiyah-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupStatsRoutes(app *fiber.App) {
	statsService := service.NewStatsService()
	stats := app.Group("/api/stats", middleware.Authorization())

	stats.Get("/", statsService.GetGlobalStats)
	stats.Get("/:cityId", statsService.GetCityStats)
}