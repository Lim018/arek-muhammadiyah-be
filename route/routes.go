package route

import (
	"arek-muhammadiyah-be/app/service"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, wilayahService *service.WilayahService) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Arek Muhammadiyah API is running",
			"version": "3.0.0",
		})
	})

	SetupAuthRoutes(app, wilayahService)
	SetupUserRoutes(app, wilayahService)
	SetupWilayahRoutes(app, wilayahService)
	SetupArticleRoutes(app)
	SetupTicketRoutes(app)
	SetupDocumentRoutes(app)
	SetupCategoryRoutes(app)
	SetupDashboardRoutes(app)
}