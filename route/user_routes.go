package route

import (
	"github.com/Lim018/arek-muhammadiyah-be/app/model"
	"github.com/Lim018/arek-muhammadiyah-be/app/service"
	"github.com/Lim018/arek-muhammadiyah-be/helper"
	"github.com/Lim018/arek-muhammadiyah-be/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	users := app.Group("/api/users")
	userService := service.NewUserService()

	// Protected routes
	users.Use(middleware.JWTMiddleware())

	users.Get("/", func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		users, pagination, err := userService.GetAllUsers(page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.PaginatedResponse{
			Success:    true,
			Message:    "Users retrieved successfully",
			Data:       users,
			Pagination: pagination,
		})
	})

	users.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		user, err := userService.GetUserByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "User not found",
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "User retrieved successfully",
			Data:    user,
		})
	})

	users.Put("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req model.UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		user, err := userService.UpdateUser(id, &req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "User updated successfully",
			Data:    user,
		})
	})

	users.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id := c.Params("id")
		err := userService.DeleteUser(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "User deleted successfully",
		})
	})

	// Bulk create users from CSV
	users.Post("/bulk", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		file, err := c.FormFile("csv")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "CSV file required",
			})
		}

		src, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: "Failed to open CSV file",
			})
		}
		defer src.Close()

		userRequests, err := helper.ParseUsersFromCSV(src)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Failed to parse CSV: " + err.Error(),
			})
		}

		createdUsers, err := userService.BulkCreateUsers(userRequests)
		if err != nil {
			return c.Status(fiber.StatusPartialContent).JSON(model.Response{
				Success: false,
				Message: "Some users failed to create",
				Data:    createdUsers,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(model.Response{
			Success: true,
			Message: "Users created successfully",
			Data:    createdUsers,
		})
	})
}