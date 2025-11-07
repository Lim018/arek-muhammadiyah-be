package service

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
	"arek-muhammadiyah-be/helper/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	userRepo       *repository.UserRepository
	wilayahService *WilayahService
}

func NewUserService(wilayahService *WilayahService) *UserService {
	return &UserService{
		userRepo:       repository.NewUserRepository(),
		wilayahService: wilayahService,
	}
}

// EnrichUserWithWilayah - Tambahkan info wilayah ke user
func (s *UserService) enrichUserWithWilayah(user *model.User) {
	if user.VillageID != nil && *user.VillageID != "" {
		cityName, districtName, villageName := s.wilayahService.GetWilayahInfo(*user.VillageID)
		user.VillageName = villageName
		user.DistrictName = districtName
		user.CityName = cityName
		
		// Extract IDs from village_id
		// Format: "3576011001" -> city: "3576", district: "357601"
		if len(*user.VillageID) >= 6 {
			user.CityID = (*user.VillageID)[:4]
			user.DistrictID = (*user.VillageID)[:6]
		}
	}
}

func (s *UserService) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	filters := repository.UserFilterParams{
		Search:     c.Query("search"),
		CityID:     c.Query("city_id"),
		DistrictID: c.Query("district_id"),
	}

	users, total, err := s.userRepo.GetAllWithFilters(limit, offset, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	for i := range users {
		s.enrichUserWithWilayah(&users[i])
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}

func (s *UserService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "User not found",
		})
	}

	// Enrich dengan data wilayah
	s.enrichUserWithWilayah(user)

	return c.JSON(model.Response{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

func (s *UserService) CreateUser(c *fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to hash password",
		})
	}

	user := &model.User{
		Name:      req.Name,
		Password:  hashedPassword,
		BirthDate: req.BirthDate,
		Telp:      req.Telp,
		Gender:    req.Gender,
		Job:       req.Job,
		RoleID:    req.RoleID,
		VillageID: req.VillageID, // Simpan village_id dari JSON
		NIK:       req.NIK,
		Address:   req.Address,
		IsMobile:  helper.GetBoolValue(req.IsMobile, false),
	}

	if err := s.userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	// Enrich dengan data wilayah sebelum return
	s.enrichUserWithWilayah(user)

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

func (s *UserService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	existing, err := s.userRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "User not found",
		})
	}

	updateData := &model.User{
		Name:      helper.GetStringValue(req.Name, existing.Name),
		BirthDate: req.BirthDate,
		Telp:      helper.GetStringPointer(req.Telp, existing.Telp),
		Gender:    helper.GetStringPointer(req.Gender, existing.Gender),
		Job:       helper.GetStringPointer(req.Job, existing.Job),
		RoleID:    helper.GetUintPointer(req.RoleID, existing.RoleID),
		VillageID: helper.GetStringPointer(req.VillageID, existing.VillageID),
		NIK:       helper.GetStringPointer(req.NIK, existing.NIK),
		Address:   helper.GetStringPointer(req.Address, existing.Address),
	}

	if err := s.userRepo.Update(id, updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	// Get updated user
	updatedUser, _ := s.userRepo.GetByID(id)
	s.enrichUserWithWilayah(updatedUser)

	return c.JSON(model.Response{
		Success: true,
		Message: "User updated successfully",
		Data:    updatedUser,
	})
}

func (s *UserService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := s.userRepo.GetByID(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "User not found",
		})
	}

	if err := s.userRepo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

func (s *UserService) GetMobileUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := s.userRepo.GetMobileUsers(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	// Enrich dengan data wilayah
	for i := range users {
		s.enrichUserWithWilayah(&users[i])
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Mobile users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}

func (s *UserService) GetByGender(c *fiber.Ctx) error {
	gender := c.Params("gender")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := s.userRepo.GetByGender(gender, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	// Enrich dengan data wilayah
	for i := range users {
		s.enrichUserWithWilayah(&users[i])
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}

// GetByCity - Get users by city_id
func (s *UserService) GetByCity(c *fiber.Ctx) error {
	cityID := c.Params("cityId")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := s.userRepo.GetByCity(cityID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	// Enrich dengan data wilayah
	for i := range users {
		s.enrichUserWithWilayah(&users[i])
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}

// GetByDistrict - Get users by district_id
func (s *UserService) GetByDistrict(c *fiber.Ctx) error {
	districtID := c.Params("districtId")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := s.userRepo.GetByDistrict(districtID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	// Enrich dengan data wilayah
	for i := range users {
		s.enrichUserWithWilayah(&users[i])
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}