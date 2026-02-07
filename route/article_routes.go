package route

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/service"
	"arek-muhammadiyah-be/middleware"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	articles.Use(middleware.Authorization())

	articles.Get("/category/:categoryID", func(c *fiber.Ctx) error {
    categoryIDParam := c.Params("categoryID")
    categoryID, err := strconv.ParseUint(categoryIDParam, 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response{
            Success: false,
            Message: "Invalid category ID",
        })
    }

    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    articles, pagination, err := articleService.GetArticlesByCategory(uint(categoryID), page, limit)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
            Success: false,
            Message: err.Error(),
        })
    }

    return c.JSON(model.PaginatedResponse{
        Success:    true,
        Message:    "Articles filtered by category retrieved successfully",
        Data:       articles,
        Pagination: pagination,
    })
})


	articles.Post("/", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)

		// Check Content-Type untuk menentukan cara parsing
		contentType := c.Get("Content-Type")
		var req model.CreateArticleRequest

		if len(contentType) >= 19 && contentType[:19] == "multipart/form-data" {
			// Parse multipart form data
			form, err := c.MultipartForm()
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(model.Response{
					Success: false,
					Message: "Failed to parse form data: " + err.Error(),
				})
			}

			// Extract title
			if titles, ok := form.Value["title"]; ok && len(titles) > 0 {
				req.Title = titles[0]
			}

			// Extract content
			if contents, ok := form.Value["content"]; ok && len(contents) > 0 {
				req.Content = contents[0]
			}

			// Extract category_id
			if categoryIDs, ok := form.Value["category_id"]; ok && len(categoryIDs) > 0 && categoryIDs[0] != "" {
				catID, err := strconv.ParseUint(categoryIDs[0], 10, 32)
				if err == nil {
					catUint := uint(catID)
					req.CategoryID = &catUint
				}
			}

			// Extract is_published
			if isPublished, ok := form.Value["is_published"]; ok && len(isPublished) > 0 {
				published := isPublished[0] == "true"
				req.IsPublished = &published
			}

			// Handle feature_image upload
			if files, ok := form.File["feature_image"]; ok && len(files) > 0 {
				file := files[0]

				// Generate unique filename
				ext := filepath.Ext(file.Filename)
				newFileName := fmt.Sprintf("feature_%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
				uploadPath := filepath.Join("uploads", "articles", newFileName)

				// Ensure directory exists
				if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
						Success: false,
						Message: "Failed to create upload directory",
					})
				}

				// Save file
				if err := c.SaveFile(file, uploadPath); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
						Success: false,
						Message: "Failed to save feature image",
					})
				}

				req.FeatureImage = &uploadPath
			}

			// Handle document uploads (lampiran)
			if files, ok := form.File["documents"]; ok {
				for _, file := range files {
					ext := filepath.Ext(file.Filename)
					newFileName := fmt.Sprintf("doc_%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
					uploadPath := filepath.Join("uploads", "documents", newFileName)

					if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
						continue
					}

					if err := c.SaveFile(file, uploadPath); err != nil {
						continue
					}

					fileSize := file.Size
					mimeType := file.Header.Get("Content-Type")
					req.Documents = append(req.Documents, model.CreateDocumentRequest{
						Title:    file.Filename,
						FilePath: uploadPath,
						FileName: file.Filename,
						FileSize: &fileSize,
						MimeType: &mimeType,
					})
				}
			}
		} else {
			// Fallback ke JSON parsing
			if err := c.BodyParser(&req); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(model.Response{
					Success: false,
					Message: "Invalid request body",
				})
			}
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