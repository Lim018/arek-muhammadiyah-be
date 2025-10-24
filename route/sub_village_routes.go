package route

import (
	"arek-muhammadiyah-be/app/service"
	"arek-muhammadiyah-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupSubVillageRoutes(app *fiber.App) {
	subVillageService := service.NewSubVillageService()
	subVillages := app.Group("/api/sub-villages")

	// Public routes
	subVillages.Get("/", subVillageService.GetAll)
	subVillages.Get("/village/:villageId", subVillageService.GetByVillageID)
	subVillages.Get("/map", subVillageService.GetWithUserCount)
	subVillages.Get("/stats", subVillageService.GetWithStats)

	// Protected routes
	subVillages.Use(middleware.Authorization())
	subVillages.Use(middleware.AdminOnly())
	
	subVillages.Post("/", subVillageService.Create)
	subVillages.Put("/:id", subVillageService.Update)
	subVillages.Delete("/:id", subVillageService.Delete)
}