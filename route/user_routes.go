package route

import (
	"arek-muhammadiyah-be/app/service"
	"arek-muhammadiyah-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userService := service.NewUserService()
	users := app.Group("/api/users", middleware.Authorization())

	users.Get("/", userService.GetAll)
	users.Get("/mobile", userService.GetMobileUsers)
	users.Get("/:id", userService.GetByID)
	users.Post("/", middleware.AdminOnly(), userService.CreateUser)
	users.Put("/:id", middleware.AdminOnly(), userService.Update)
	users.Delete("/:id", middleware.AdminOnly(), userService.Delete)
	users.Post("/bulk", middleware.AdminOnly(), userService.BulkCreate)
	users.Get("/village/:villageId", userService.GetByVillage)
	users.Get("/sub-village/:subVillageId", userService.GetBySubVillage)
	users.Get("/gender/:gender", userService.GetByGender)
}