package route

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
	"arek-muhammadiyah-be/middleware"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

func SetupCategoryRoutes(app *fiber.App) {
	categories := app.Group("/api/categories")
	categoryRepo := repository.NewCategoryRepository()

	// Public routes
	categories.Get("/", func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		activeOnly := c.Query("active", "false") == "true"
		offset := (page - 1) * limit

		categories, total, err := categoryRepo.GetAll(limit, offset, activeOnly)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		pagination := helper.CreatePagination(int64(page), int64(limit), total)
		return c.JSON(model.PaginatedResponse{
			Success:    true,
			Message:    "Categories retrieved successfully",
			Data:       categories,
			Pagination: pagination,
		})
	})

	categories.Get("/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		category, err := categoryRepo.GetByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "Category not found",
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Category retrieved successfully",
			Data:    category,
		})
	})

	// Protected routes
	categories.Use(middleware.Authorization())
	categories.Use(middleware.AdminOnly())

	categories.Post("/", func(c *fiber.Ctx) error {
		var req model.CreateCategoryRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		category := &model.Category{
			Name:        req.Name,
			Description: req.Description,
			Color:       helper.GetStringValue(req.Color, "#10B981"),
			IsActive:    helper.GetBoolValue(req.IsActive, true),
		}

		if err := categoryRepo.Create(category); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(model.Response{
			Success: true,
			Message: "Category created successfully",
			Data:    category,
		})
	})

	categories.Put("/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		var req model.CreateCategoryRequest
		
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		existing, err := categoryRepo.GetByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "Category not found",
			})
		}

		updateData := &model.Category{
			Name:        helper.GetStringValue(&req.Name, existing.Name),
			Description: helper.GetStringPointer(req.Description, existing.Description),
			Color:       helper.GetStringValue(req.Color, existing.Color),
			IsActive:    helper.GetBoolValue(req.IsActive, existing.IsActive),
		}

		if err := categoryRepo.Update(uint(id), updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Category updated successfully",
			Data:    updateData,
		})
	})

	categories.Delete("/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		
		if _, err := categoryRepo.GetByID(uint(id)); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "Category not found",
			})
		}

		if err := categoryRepo.Delete(uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Category deleted successfully",
		})
	})
}