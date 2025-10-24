package route

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Health check
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Organization System API is running",
		})
	})

	// Setup all routes
	SetupAuthRoutes(app)
	SetupUserRoutes(app)
	SetupArticleRoutes(app)
	SetupTicketRoutes(app)
	SetupVillageRoutes(app)
	SetupSubVillageRoutes(app) 
	SetupDocumentRoutes(app)
	SetupCategoryRoutes(app)
	SetupDashboardRoutes(app)
}