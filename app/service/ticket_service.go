package service

import (
	"crypto/rand"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"github.com/google/uuid"

	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TicketService struct {
	ticketRepo   *repository.TicketRepository
	documentRepo *repository.DocumentRepository
}

func NewTicketService() *TicketService {
	return &TicketService{
		ticketRepo:   repository.NewTicketRepository(),
		documentRepo: repository.NewDocumentRepository(),
	}
}

// ========== SEMUA METHOD KONSISTEN RETURN fiber.Handler ==========

func (s *TicketService) GetAllTickets() fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		statusStr := c.Query("status")
		
		var status *model.TicketStatus
		if statusStr != "" {
			s := model.TicketStatus(statusStr)
			status = &s
		}

		offset := (page - 1) * limit
		tickets, total, err := s.ticketRepo.GetAll(limit, offset, status)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		pagination := helper.CreatePagination(int64(page), int64(limit), total)

		return c.JSON(model.PaginatedResponse{
			Success:    true,
			Message:    "Tickets retrieved successfully",
			Data:       tickets,
			Pagination: pagination,
		})
	}
}

func (s *TicketService) GetUserTickets() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		offset := (page - 1) * limit
		tickets, total, err := s.ticketRepo.GetByUserID(userID, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		pagination := helper.CreatePagination(int64(page), int64(limit), total)

		return c.JSON(model.PaginatedResponse{
			Success:    true,
			Message:    "User tickets retrieved successfully",
			Data:       tickets,
			Pagination: pagination,
		})
	}
}

func (s *TicketService) GetTicketByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		ticket, err := s.ticketRepo.GetByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "Ticket not found",
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Ticket retrieved successfully",
			Data:    ticket,
		})
	}
}

func (s *TicketService) GetTicketStats() fiber.Handler {
	return func(c *fiber.Ctx) error {
		stats, err := s.ticketRepo.GetCountByStatus()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Ticket statistics retrieved successfully",
			Data:    stats,
		})
	}
}

func (s *TicketService) UpdateTicket() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		
		var req model.UpdateTicketRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		existing, err := s.ticketRepo.GetByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "Ticket not found",
			})
		}

		updateData := &model.Ticket{}
		if req.Status != nil {
			updateData.Status = *req.Status
			if *req.Status == model.TicketStatusResolved || *req.Status == model.TicketStatusRejected {
				now := time.Now()
				updateData.ResolvedAt = &now
			} else {
				updateData.ResolvedAt = existing.ResolvedAt
			}
		}
		if req.Resolution != nil {
			updateData.Resolution = req.Resolution
		}

		if err := s.ticketRepo.Update(uint(id), updateData); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		updatedTicket, _ := s.ticketRepo.GetByID(uint(id))

		return c.JSON(model.Response{
			Success: true,
			Message: "Ticket updated successfully",
			Data:    updatedTicket,
		})
	}
}

func (s *TicketService) DeleteTicket() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

		if _, err := s.ticketRepo.GetByID(uint(id)); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "Ticket not found",
			})
		}

		if err := s.ticketRepo.Delete(uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Ticket deleted successfully",
		})
	}
}

// CreateTicketWithFiles - Create ticket + upload files sekaligus
func (s *TicketService) CreateTicketWithFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		
		// Parse multipart form
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid form data",
			})
		}

		// Parse form fields
		var req model.CreateTicketRequest
		req.Title = c.FormValue("title")
		req.Description = c.FormValue("description")
		
		if categoryIDStr := c.FormValue("category_id"); categoryIDStr != "" {
			categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)
			categoryIDUint := uint(categoryID)
			req.CategoryID = &categoryIDUint
		}

		// Get uploaded files
		files := form.File["documents"]

		// Validate required fields
		if req.Title == "" || req.Description == "" {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Title and description are required",
			})
		}

		// Validate files jika ada
		if len(files) > 0 {
			for _, file := range files {
				if err := s.ValidateFile(file, 10*1024*1024, []string{
					"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp",
					"application/pdf",
					"application/msword",
					"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				}); err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(model.Response{
						Success: false,
						Message: err.Error(),
					})
				}
			}
		}

		// Start transaction
		tx := s.ticketRepo.Begin()

		// 1. Create ticket dulu
		ticket := &model.Ticket{
			UserID:      userID,
			CategoryID:  req.CategoryID,
			Title:       req.Title,
			Description: req.Description,
			Status:      model.TicketStatusUnread,
		}

		if err := tx.Create(ticket).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: "Failed to create ticket",
			})
		}

		// 2. Upload files jika ada
		if len(files) > 0 {
			for _, file := range files {
				if _, err := s.uploadFileToTicket(tx, ticket.ID, file); err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
						Success: false,
						Message: fmt.Sprintf("Failed to upload file: %v", err),
					})
				}
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: "Failed to save data",
			})
		}

		// Get complete ticket data
		finalTicket, _ := s.ticketRepo.GetByID(ticket.ID)

		return c.Status(fiber.StatusCreated).JSON(model.Response{
			Success: true,
			Message: "Ticket created successfully",
			Data:    finalTicket,
		})
	}
}

// UploadFilesToTicket - Upload files ke existing ticket
func (s *TicketService) UploadFilesToTicket() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ticketID, _ := strconv.ParseUint(c.Params("id"), 10, 32)

		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid form data",
			})
		}

		files := form.File["documents"]
		if len(files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "No files uploaded",
			})
		}

		// Validate files
		for _, file := range files {
			if err := s.ValidateFile(file, 10*1024*1024, []string{
				"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp",
				"application/pdf",
				"application/msword",
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			}); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(model.Response{
					Success: false,
					Message: err.Error(),
				})
			}
		}

		// Upload files
		tx := s.ticketRepo.Begin()
		for _, file := range files {
			if _, err := s.uploadFileToTicket(tx, uint(ticketID), file); err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
					Success: false,
					Message: err.Error(),
				})
			}
		}
		
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: "Failed to save data",
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Files uploaded successfully",
		})
	}
}

// GetTicketFiles - Get files dari ticket
func (s *TicketService) GetTicketFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ticketID, _ := strconv.ParseUint(c.Params("id"), 10, 32)

		files, _, err := s.documentRepo.GetByTicketID(uint(ticketID), 0, 0)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Files retrieved successfully",
			Data:    files,
		})
	}
}

// DeleteFile - Hapus file dari ticket
func (s *TicketService) DeleteFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fileIDParam := c.Params("fileId")

		// validasi format UUID
		fileID, err := uuid.Parse(fileIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid file ID format",
			})
		}

		// hapus file berdasarkan UUID
		if err := s.documentRepo.Delete(fileID.String()); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "File deleted successfully",
		})
	}
}

func (s *TicketService) ServeFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fileIDParam := c.Params("fileId")

		fileID, err := uuid.Parse(fileIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid file ID format",
			})
		}

		document, err := s.documentRepo.GetByID(fileID.String())
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "File not found",
			})
		}

		if _, err := os.Stat(document.FilePath); os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "File not found on server",
			})
		}

		return c.SendFile(document.FilePath)
	}
}

// ValidateFile - Validasi file
func (s *TicketService) ValidateFile(file *multipart.FileHeader, maxSize int64, allowedTypes []string) error {
	if file.Size > maxSize {
		return fmt.Errorf("file %s exceeds maximum size of %dMB", file.Filename, maxSize/(1024*1024))
	}

	contentType := file.Header.Get("Content-Type")
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return nil
		}
	}

	return fmt.Errorf("file type %s not allowed", contentType)
}


func (s *TicketService) uploadFileToTicket(tx *gorm.DB, ticketID uint, file *multipart.FileHeader) (*model.Document, error) {
	uploadDir := "./uploads/tickets"
	
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), s.generateRandomString(10), fileExt)
	filePath := filepath.Join(uploadDir, fileName)

	// Save file
	if err := s.saveFileToDisk(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file %s: %v", file.Filename, err)
	}

	// Create document
	description := fmt.Sprintf("Attachment for ticket #%d", ticketID)
	mimeType := file.Header.Get("Content-Type")
	document := &model.Document{
		TicketID:    &ticketID,
		FileName:    file.Filename,
		FilePath:    filePath,
		FileSize:    &file.Size,
		MimeType:    &mimeType,
		Title:       file.Filename,
		Description: &description,
	}

	if err := tx.Create(document).Error; err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create document record: %v", err)
	}

	return document, nil
}

func (s *TicketService) saveFileToDisk(file *multipart.FileHeader, destination string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	return err
}

func (s *TicketService) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
