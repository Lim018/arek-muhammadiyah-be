package service

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
	"arek-muhammadiyah-be/helper/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	userRepo       *repository.UserRepository
	wilayahService *WilayahService
}

func NewAuthService(wilayahService *WilayahService) *AuthService {
	return &AuthService{
		userRepo:       repository.NewUserRepository(),
		wilayahService: wilayahService,
	}
}

// enrichUserWithWilayah - Tambahkan info wilayah ke user
func (s *AuthService) enrichUserWithWilayah(user *model.User) {
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

	// Enrich dengan data wilayah
	s.enrichUserWithWilayah(user)

	token, err := utils.GenerateToken(fmt.Sprintf("%d", user.ID), user.RoleID)
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

	// Hash password
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
		Message: "User registered successfully",
		Data:    user,
	})
}