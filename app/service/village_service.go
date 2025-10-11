package service

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type VillageService struct {
	villageRepo *repository.VillageRepository
}

func NewVillageService() *VillageService {
	return &VillageService{
		villageRepo: repository.NewVillageRepository(),
	}
}

func (s *VillageService) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	activeOnly := c.Query("active", "false") == "true"
	offset := (page - 1) * limit

	villages, total, err := s.villageRepo.GetAll(limit, offset, activeOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Villages retrieved successfully",
		Data:       villages,
		Pagination: pagination,
	})
}

func (s *VillageService) GetWithUserCount(c *fiber.Ctx) error {
	villages, err := s.villageRepo.GetWithUserCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Villages with user count retrieved successfully",
		Data:    villages,
	})
}

// NEW: Get villages with complete stats (for map)
func (s *VillageService) GetWithStats(c *fiber.Ctx) error {
	villages, err := s.villageRepo.GetWithCompleteStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Villages with statistics retrieved successfully",
		Data:    villages,
	})
}

func (s *VillageService) Create(c *fiber.Ctx) error {
	var req model.CreateVillageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	village := &model.Village{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       helper.GetStringValue(req.Color, "#3B82F6"),
		IsActive:    helper.GetBoolValue(req.IsActive, true),
	}

	if err := s.villageRepo.Create(village); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Success: true,
		Message: "Village created successfully",
		Data:    village,
	})
}

func (s *VillageService) Update(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var req model.CreateVillageRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	existing, err := s.villageRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "Village not found",
		})
	}

	updateData := &model.Village{
		Name:        helper.GetStringValue(&req.Name, existing.Name),
		Code:        helper.GetStringValue(&req.Code, existing.Code),
		Description: helper.GetStringPointer(req.Description, existing.Description),
		Color:       helper.GetStringValue(req.Color, existing.Color),
		IsActive:    helper.GetBoolValue(req.IsActive, existing.IsActive),
	}

	if err := s.villageRepo.Update(uint(id), updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Village updated successfully",
		Data:    updateData,
	})
}

func (s *VillageService) Delete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	
	if _, err := s.villageRepo.GetByID(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "Village not found",
		})
	}

	if err := s.villageRepo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Village deleted successfully",
	})
}