package service

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"

	"github.com/gofiber/fiber/v2"
)

type StatsService struct {
	statsRepo *repository.StatsRepository
}

func NewStatsService() *StatsService {
	return &StatsService{
		statsRepo: repository.NewStatsRepository(),
	}
}

// GetGlobalStats - Get global statistics
func (s *StatsService) GetGlobalStats(c *fiber.Ctx) error {
	// Get total users
	totalUsers, err := s.statsRepo.GetTotalUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to get total users",
		})
	}

	// Get ticket stats
	ticketStats, err := s.statsRepo.GetTicketStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to get ticket stats",
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "Statistics retrieved successfully",
		Data: fiber.Map{
			"total_users":  totalUsers,
			"ticket_stats": ticketStats,
		},
	})
}

// GetCityStats - Get statistics by city
func (s *StatsService) GetCityStats(c *fiber.Ctx) error {
	cityID := c.Params("cityId")

	// Get total users in city
	totalUsers, err := s.statsRepo.GetTotalUsersByCity(cityID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to get total users",
		})
	}

	// Get ticket stats for city
	ticketStats, err := s.statsRepo.GetTicketStatsByCity(cityID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Success: false,
			Message: "Failed to get ticket stats",
		})
	}

	return c.JSON(model.Response{
		Success: true,
		Message: "City statistics retrieved successfully",
		Data: fiber.Map{
			"city_id":      cityID,
			"total_users":  totalUsers,
			"ticket_stats": ticketStats,
		},
	})
}