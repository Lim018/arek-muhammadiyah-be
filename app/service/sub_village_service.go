package service

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type SubVillageService struct {
	subVillageRepo *repository.SubVillageRepository
}

func NewSubVillageService() *SubVillageService {
	return &SubVillageService{
		subVillageRepo: repository.NewSubVillageRepository(),
	}
}

func (s *SubVillageService) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	activeOnly := c.Query("active", "false") == "true"
	offset := (page - 1) * limit

	subVillages, total, err := s.subVillageRepo.GetAll(limit, offset, activeOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Sub villages retrieved successfully",
		Data:       subVillages,
		Pagination: pagination,
	})
}

func (s *SubVillageService) GetByVillageID(c *fiber.Ctx) error {
	villageID, _ := strconv.ParseUint(c.Params("villageId"), 10, 32)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	subVillages, total, err := s.subVillageRepo.GetByVillageID(uint(villageID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Sub villages retrieved successfully",
		Data:       subVillages,
		Pagination: pagination,
	})
}

func (s *SubVillageService) GetWithUserCount(c *fiber.Ctx) error {
	subVillages, err := s.subVillageRepo.GetWithUserCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Sub villages with user count retrieved successfully",
		Data:    subVillages,
	})
}

func (s *SubVillageService) GetWithStats(c *fiber.Ctx) error {
	subVillages, err := s.subVillageRepo.GetWithCompleteStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Sub villages with statistics retrieved successfully",
		Data:    subVillages,
	})
}

func (s *SubVillageService) Create(c *fiber.Ctx) error {
	var req model.CreateSubVillageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	subVillage := &model.SubVillage{
		VillageID:   req.VillageID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    helper.GetBoolValue(req.IsActive, true),
	}

	if err := s.subVillageRepo.Create(subVillage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Success: true,
		Message: "Sub village created successfully",
		Data:    subVillage,
	})
}

func (s *SubVillageService) Update(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var req model.CreateSubVillageRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	existing, err := s.subVillageRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "Sub village not found",
		})
	}

	updateData := &model.SubVillage{
		VillageID:   req.VillageID,
		Name:        helper.GetStringValue(&req.Name, existing.Name),
		Code:        helper.GetStringValue(&req.Code, existing.Code),
		Description: helper.GetStringPointer(req.Description, existing.Description),
		IsActive:    helper.GetBoolValue(req.IsActive, existing.IsActive),
	}

	if err := s.subVillageRepo.Update(uint(id), updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Sub village updated successfully",
		Data:    updateData,
	})
}

func (s *SubVillageService) Delete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	
	if _, err := s.subVillageRepo.GetByID(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "Sub village not found",
		})
	}

	if err := s.subVillageRepo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Sub village deleted successfully",
	})
}