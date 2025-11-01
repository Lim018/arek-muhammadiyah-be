package service

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type DocumentService struct {
	documentRepo *repository.DocumentRepository
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo: repository.NewDocumentRepository(),
	}
}

func (s *DocumentService) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	documents, total, err := s.documentRepo.GetAll(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Documents retrieved successfully",
		Data:       documents,
		Pagination: pagination,
	})
}

func (s *DocumentService) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	document, err := s.documentRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "Document not found",
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Document retrieved successfully",
		Data:    document,
	})
}

func (s *DocumentService) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("userId")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	documents, total, err := s.documentRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Documents retrieved successfully",
		Data:       documents,
		Pagination: pagination,
	})
}

func (s *DocumentService) Create(c *fiber.Ctx) error {
	var req model.CreateDocumentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	document := &model.Document{
		Title:       req.Title,
		Description: req.Description,
		FilePath:    req.FilePath,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		MimeType:    req.MimeType,
		TicketID:    *req.TicketID,
		ArticleID:   *req.ArticleID,
	}

	if err := s.documentRepo.Create(document); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Success: true,
		Message: "Document created successfully",
		Data:    document,
	})
}


func (s *DocumentService) Delete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	if _, err := s.documentRepo.GetByID(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "Document not found",
		})
	}

	if err := s.documentRepo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Document deleted successfully",
	})
}
