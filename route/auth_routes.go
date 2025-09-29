package route

import (
	"github.com/Lim018/arek-muhammadiyah-be/app/model"
	"github.com/Lim018/arek-muhammadiyah-be/app/service"
	"github.com/Lim018/arek-muhammadiyah-be/helper"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/api/auth")
	userService := service.NewUserService()

	auth.Post("/login", func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		user, err := userService.Login(&req)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		token, err := helper.GenerateToken(user.ID, user.RoleID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: "Failed to generate token",
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Login successful",
			Data: fiber.Map{
				"user":  user,
				"token": token,
			},
		})
	})

	auth.Post("/register", func(c *fiber.Ctx) error {
		var req model.CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		user, err := userService.CreateUser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(model.Response{
			Success: true,
			Message: "User created successfully",
			Data:    user,
		})
	})
}