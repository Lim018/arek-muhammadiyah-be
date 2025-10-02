package route

import (
	"arek-muhammadiyah-be/app/service"
	"arek-muhammadiyah-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupVillageRoutes(app *fiber.App) {
	villageService := service.NewVillageService()
	villages := app.Group("/api/villages")

	// Public routes
	villages.Get("/", villageService.GetAll)
	villages.Get("/map", villageService.GetWithUserCount)

	// Protected routes
	villages.Use(middleware.Authorization())
	villages.Use(middleware.AdminOnly())
	
	villages.Post("/", villageService.Create)
	villages.Put("/:id", villageService.Update)
	villages.Delete("/:id", villageService.Delete)
}