package route

import (
	"arek-muhammadiyah-be/app/service"
	"arek-muhammadiyah-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, wilayahService *service.WilayahService) {
	userService := service.NewUserService(wilayahService)
	users := app.Group("/api/users", middleware.Authorization())

	users.Get("/", userService.GetAll)
	users.Get("/mobile", userService.GetMobileUsers)
	users.Get("/:id", userService.GetByID)
	users.Post("/", middleware.AdminOnly(), userService.CreateUser)
	users.Put("/:id", userService.Update)
	users.Delete("/:id", middleware.AdminOnly(), userService.Delete)
	users.Get("/gender/:gender", userService.GetByGender)
	
	// Filter by wilayah
	users.Get("/city/:cityId", userService.GetByCity)
	users.Get("/district/:districtId", userService.GetByDistrict)
}