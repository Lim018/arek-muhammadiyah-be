package service

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// Login handler
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	user, err := s.userRepo.GetByID(req.ID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Success: false,
			Message: "Invalid credentials",
		})
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
    return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
        Success: false,
        Message: "Invalid credentials",
    })
}

	token, err := utils.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to generate token",
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Login successful",
		Data: fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

// Register handler
func (s *AuthService) Register(c *fiber.Ctx) error {
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
