package route

import (
	"arek-muhammadiyah-be/app/service"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, wilayahService *service.WilayahService) {
	authService := service.NewAuthService(wilayahService)
	auth := app.Group("/api/auth")

	auth.Post("/login", authService.Login)
	auth.Post("/register", authService.Register)
}