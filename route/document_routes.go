package route

import (
	"arek-muhammadiyah-be/app/service"
	"arek-muhammadiyah-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupDocumentRoutes(app *fiber.App) {
	documentService := service.NewDocumentService()
	documents := app.Group("/api/documents", middleware.Authorization())

	documents.Get("/", middleware.AdminOnly(), documentService.GetAll)
	documents.Get("/:id", documentService.GetByID)
	documents.Get("/user/:userId", documentService.GetByUserID)
	documents.Post("/", documentService.Create)
	documents.Delete("/:id", documentService.Delete)
}