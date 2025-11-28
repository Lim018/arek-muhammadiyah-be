package route

import (
	"arek-muhammadiyah-be/app/service"
	"arek-muhammadiyah-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, wilayahService *service.WilayahService) {
	authService := service.NewAuthService(wilayahService)
	auth := app.Group("/api/auth")

	// Public routes
	auth.Post("/login", authService.Login)
	auth.Post("/register", authService.Register)
	
	// Protected routes
	auth.Post("/logout", middleware.Authorization(), authService.Logout)
	auth.Get("/navbar", middleware.Authorization(), authService.GetNavbar)
}