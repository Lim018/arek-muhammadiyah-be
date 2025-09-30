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
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// Get all users
func (s *UserService) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := s.userRepo.GetAll(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}

// Get user by ID
func (s *UserService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.userRepo.GetByID(id)
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
}

// Create user
func (s *UserService) CreateUser(c *fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to hash password",
		})
	}

	user := &model.User{
		ID:        req.ID,
		Name:      req.Name,
		Password:  hashedPassword,
		RoleID:    req.RoleID,
		VillageID: req.VillageID,
		NIK:       req.NIK,
		Address:   req.Address,
	}

	if err := s.userRepo.Create(user); err != nil {
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
}

// Update user
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
		RoleID:    helper.GetUintPointer(req.RoleID, existing.RoleID),
		VillageID: helper.GetUintPointer(req.VillageID, existing.VillageID),
		NIK:       helper.GetStringPointer(req.NIK, existing.NIK),
		Address:   helper.GetStringPointer(req.Address, existing.Address),
		Photo:     helper.GetStringPointer(req.Photo, existing.Photo),
	}

	if err := s.userRepo.Update(id, updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "User updated successfully",
		Data:    updateData,
	})
}

// Delete user
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

// Bulk create via CSV
func (s *UserService) BulkCreate(c *fiber.Ctx) error {
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

	users, err := helper.ParseUsersFromCSV(src)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Failed to parse CSV: " + err.Error(),
		})
	}

	var createdUsers []model.User
	var failedUsers []string

	for _, userReq := range users {
		hashedPassword, err := utils.HashPassword(userReq.Password)
		if err != nil {
			failedUsers = append(failedUsers, userReq.ID)
			continue
		}

		u := &model.User{
			ID:        userReq.ID,
			Name:      userReq.Name,
			Password:  hashedPassword,
			RoleID:    userReq.RoleID,
			VillageID: userReq.VillageID,
			NIK:       userReq.NIK,
			Address:   userReq.Address,
		}

		if err := s.userRepo.Create(u); err != nil {
			failedUsers = append(failedUsers, userReq.ID)
			continue
		}

		createdUsers = append(createdUsers, *u)
	}

	if len(failedUsers) > 0 {
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
}
