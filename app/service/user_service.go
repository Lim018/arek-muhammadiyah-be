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
		Name:         req.Name,
		Password:     hashedPassword,
		BirthDate:    req.BirthDate,
		Telp:         req.Telp,
		Gender:       req.Gender,
		Job:          req.Job,
		RoleID:       req.RoleID,
		SubVillageID: req.SubVillageID,
		NIK:          req.NIK,
		Address:      req.Address,
		IsMobile:     helper.GetBoolValue(req.IsMobile, false),
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
		Name:         helper.GetStringValue(req.Name, existing.Name),
		BirthDate:    req.BirthDate,
		Telp:         helper.GetStringPointer(req.Telp, existing.Telp),
		Gender:       helper.GetStringPointer(req.Gender, existing.Gender),
		Job:          helper.GetStringPointer(req.Job, existing.Job),
		RoleID:       helper.GetUintPointer(req.RoleID, existing.RoleID),
		SubVillageID: helper.GetUintPointer(req.SubVillageID, existing.SubVillageID),
		NIK:          helper.GetStringPointer(req.NIK, existing.NIK),
		Address:      helper.GetStringPointer(req.Address, existing.Address),
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
			Name:         userReq.Name,
			Password:     hashedPassword,
			BirthDate:    userReq.BirthDate,
			Telp:         userReq.Telp,
			Gender:       userReq.Gender,
			Job:          userReq.Job,
			RoleID:       userReq.RoleID,
			SubVillageID: userReq.SubVillageID,
			NIK:          userReq.NIK,
			Address:      userReq.Address,
			IsMobile:     helper.GetBoolValue(userReq.IsMobile, false),
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

func (s *UserService) GetBySubVillage(c *fiber.Ctx) error {
	subVillageID, _ := strconv.ParseUint(c.Params("subVillageId"), 10, 32)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := s.userRepo.GetBySubVillage(uint(subVillageID), limit, offset)
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

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return c.JSON(model.PaginatedResponse{
		Success:    true,
		Message:    "Users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}