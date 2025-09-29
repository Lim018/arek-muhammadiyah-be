package route

import (
	"github.com/Lim018/arek-muhammadiyah-be/app/model"
	"github.com/Lim018/arek-muhammadiyah-be/app/service"
	"github.com/Lim018/arek-muhammadiyah-be/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SetupArticleRoutes(app *fiber.App) {
	articles := app.Group("/api/articles")
	articleService := service.NewArticleService()

	// Public routes
	articles.Get("/", func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		publishedStr := c.Query("published")
		
		var published *bool
		if publishedStr != "" {
			p := publishedStr == "true"
			published = &p
		}

		articles, pagination, err := articleService.GetAllArticles(page, limit, published)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.PaginatedResponse{
			Success:    true,
			Message:    "Articles retrieved successfully",
			Data:       articles,
			Pagination: pagination,
		})
	})

	articles.Get("/slug/:slug", func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		article, err := articleService.GetArticleBySlug(slug)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "Article not found",
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Article retrieved successfully",
			Data:    article,
		})
	})

	articles.Get("/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		article, err := articleService.GetArticleByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "Article not found",
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Article retrieved successfully",
			Data:    article,
		})
	})

	// Protected routes
	articles.Use(middleware.JWTMiddleware())

	articles.Post("/", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		var req model.CreateArticleRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		article, err := articleService.CreateArticle(userID, &req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(model.Response{
			Success: true,
			Message: "Article created successfully",
			Data:    article,
		})
	})

	articles.Put("/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		var req model.CreateArticleRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		article, err := articleService.UpdateArticle(uint(id), &req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Article updated successfully",
			Data:    article,
		})
	})

	articles.Delete("/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		err := articleService.DeleteArticle(uint(id))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Article deleted successfully",
		})
	})
}