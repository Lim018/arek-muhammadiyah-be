package service

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper/utils"
	"fmt"
	"time"

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

func (s *AuthService) enrichUserWithWilayah(user *model.User) {
	if user.VillageID != nil && *user.VillageID != "" {
		cityName, districtName, villageName := s.wilayahService.GetWilayahInfo(*user.VillageID)
		user.VillageName = villageName
		user.DistrictName = districtName
		user.CityName = cityName
		
		if len(*user.VillageID) >= 6 {
			user.CityID = (*user.VillageID)[:4]
			user.DistrictID = (*user.VillageID)[:6]
		}
	}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	user, err := s.userRepo.GetByTelp(req.Telp)
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

	s.enrichUserWithWilayah(user)

	// Generate Access Token (1 jam)
	accessToken, err := utils.GenerateAccessToken(fmt.Sprintf("%d", user.ID), user.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to generate access token",
		})
	}

	// Generate Refresh Token (30 hari)
	refreshToken, err := utils.GenerateRefreshToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to generate refresh token",
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Login successful",
		Data: fiber.Map{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// RefreshToken handler - Generate new access token using refresh token
func (s *AuthService) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Refresh token is required",
		})
	}

	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Success: false,
			Message: "Invalid or expired refresh token",
		})
	}

	// Get user dari database untuk mendapatkan roleID terbaru
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Success: false,
			Message: "User not found",
		})
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(claims.UserID, user.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to generate access token",
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Token refreshed successfully",
		Data: fiber.Map{
			"access_token": accessToken,
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

func (s *AuthService) Register(c *fiber.Ctx) error {
	var req model.CreateUserRequest // Pastikan struct ini punya field Name, Telp, Password
	if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
					Success: false,
					Message: "Invalid request body",
			})
	}

	
	// 2. Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
					Success: false,
					Message: "Failed to hash password",
			})
	}

	// 3. Set Default Values untuk User HP
	// PENTING: Hardcode RoleID agar user tidak bisa jadi admin lewat API ini
	defaultRoleID := uint(3) // Sesuaikan dengan ID role "Warga" atau "User" di DB kamu
	isMobile := true

	user := &model.User{
			Name:      req.Name,
			Password:  hashedPassword,
			Telp:      req.Telp,
			RoleID:    &defaultRoleID, // Paksa jadi User biasa
			IsMobile:  isMobile,       // Paksa jadi Mobile User
			
	}

	// 4. Simpan ke Database
	if err := s.userRepo.Create(user); err != nil {
			// Handle error duplicate entry (misal telp sudah ada)
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
					Success: false,
					Message: "Gagal registrasi: " + err.Error(),
			})
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
			Success: true,
			Message: "Registrasi berhasil, silakan login",
			Data:    user,
	})
}

func (s *AuthService) ForgotPasswordResetDefault(c *fiber.Ctx) error {
	var req model.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid birth_date format",
		})
	}

	// Ambil user lewat repository
	user, err := s.userRepo.GetByPersonalData(req.Name, birthDate, req.NIK, req.Telp)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data user tidak ditemukan",
		})
	}

	// Reset password default
	defaultPassword := "password123"
	hashedPassword, _ := utils.HashPassword(defaultPassword)
	user.Password = hashedPassword
	user.ForceChangePassword = true

	// Update user
	if err := s.userRepo.Update(fmt.Sprintf("%d", user.ID), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal reset password",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password berhasil di-reset ke default",
	})
}
