package service

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
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

// Logout handler
func (s *AuthService) Logout(c *fiber.Ctx) error {
	// Karena menggunakan JWT stateless, logout hanya perlu menghapus token di client side
	// Server hanya mengirim response sukses
	return c.JSON(model.Response{
		Success: true,
		Message: "Logout successful",
	})
}

// Navbar handler - Get user info untuk navbar
func (s *AuthService) GetNavbar(c *fiber.Ctx) error {
	// Get user ID dari JWT token (sudah di-extract oleh middleware)
	// Middleware menggunakan "user_id" dengan underscore
	userIDInterface := c.Locals("user_id")
	if userIDInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Success: false,
			Message: "Unauthorized",
		})
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Success: false,
			Message: "Invalid user ID format",
		})
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "User not found",
		})
	}

	// Get role name
	roleName := "user"
	if user.Role != nil {
		roleName = user.Role.Name
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Navbar data retrieved successfully",
		Data: fiber.Map{
			"user": fiber.Map{
				"name": user.Name,
				"role": roleName,
			},
		},
	})
}