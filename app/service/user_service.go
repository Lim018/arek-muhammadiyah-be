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
		Name:       req.Name,
		Password:   hashedPassword,
		Telp:       req.Telp,
		RoleID:     req.RoleID,
		VillageID:  req.VillageID,
		NIK:        req.NIK,
		Address:    req.Address,
		CardStatus: helper.GetStringValue(req.CardStatus, "pending"),
		IsMobile:   helper.GetBoolValue(req.IsMobile, false),
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
		Name:       helper.GetStringValue(req.Name, existing.Name),
		Telp:       helper.GetStringPointer(req.Telp, existing.Telp),
		RoleID:     helper.GetUintPointer(req.RoleID, existing.RoleID),
		VillageID:  helper.GetUintPointer(req.VillageID, existing.VillageID),
		NIK:        helper.GetStringPointer(req.NIK, existing.NIK),
		Address:    helper.GetStringPointer(req.Address, existing.Address),
		CardStatus: helper.GetStringValue(req.CardStatus, existing.CardStatus),
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
			Name:       userReq.Name,
			Password:   hashedPassword,
			Telp:       userReq.Telp,
			RoleID:     userReq.RoleID,
			VillageID:  userReq.VillageID,
			NIK:        userReq.NIK,
			Address:    userReq.Address,
			CardStatus: helper.GetStringValue(userReq.CardStatus, "pending"),
			IsMobile:   helper.GetBoolValue(userReq.IsMobile, false),
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
			Data: fiber.Map{
				"created": createdUsers,
				"failed":  failedUsers,
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Success: true,
		Message: "Users created successfully",
		Data:    createdUsers,
	})
}

func (s *UserService) GetByVillage(c *fiber.Ctx) error {
	villageID, _ := strconv.ParseUint(c.Params("villageId"), 10, 32)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := s.userRepo.GetByVillage(uint(villageID), limit, offset)
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

func (s *UserService) GetByCardStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := s.userRepo.GetByCardStatus(status, limit, offset)
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

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Mobile users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}