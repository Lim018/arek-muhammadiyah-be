package service

import (
	"encoding/json"
	"os"
	"arek-muhammadiyah-be/app/model"
	"github.com/gofiber/fiber/v2"
)

type WilayahService struct {
	cities []model.City
}

func NewWilayahService() *WilayahService {
	return &WilayahService{}
}

// LoadWilayahData - Load data dari JSON file saat startup
func (s *WilayahService) LoadWilayahData(filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &s.cities)
	if err != nil {
		return err
	}

	return nil
}

// GetAllCities - Get semua kota
func (s *WilayahService) GetAllCities(c *fiber.Ctx) error {
	return c.JSON(model.Response{
		Success: true,
		Message: "Cities retrieved successfully",
		Data:    s.cities,
	})
}

// GetDistricts - Get kecamatan berdasarkan city_id
func (s *WilayahService) GetDistricts(c *fiber.Ctx) error {
	cityID := c.Params("cityId")
	
	for _, city := range s.cities {
		if city.ID == cityID {
			return c.JSON(model.Response{
				Success: true,
				Message: "Districts retrieved successfully",
				Data:    city.Districts,
			})
		}
	}
	
	return c.Status(fiber.StatusNotFound).JSON(model.Response{
		Success: false,
		Message: "City not found",
	})
}

// GetVillages - Get kelurahan berdasarkan district_id
func (s *WilayahService) GetVillages(c *fiber.Ctx) error {
	cityID := c.Params("cityId")
	districtID := c.Params("districtId")
	
	for _, city := range s.cities {
		if city.ID == cityID {
			for _, district := range city.Districts {
				if district.ID == districtID {
					return c.JSON(model.Response{
						Success: true,
						Message: "Villages retrieved successfully",
						Data:    district.Villages,
					})
				}
			}
		}
	}
	
	return c.Status(fiber.StatusNotFound).JSON(model.Response{
		Success: false,
		Message: "District not found",
	})
}

// GetWilayahInfo - Get informasi lengkap wilayah dari village_id
func (s *WilayahService) GetWilayahInfo(villageID string) (cityName, districtName, villageName string) {
	for _, city := range s.cities {
		for _, district := range city.Districts {
			for _, village := range district.Villages {
				if village.ID == villageID {
					return city.Name, district.Name, village.Name
				}
			}
		}
	}
	return "", "", ""
}

// SearchVillages - Search kelurahan by name
func (s *WilayahService) SearchVillages(c *fiber.Ctx) error {
	query := c.Query("q", "")
	
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Success: false,
			Message: "Query parameter 'q' is required",
		})
	}
	
	var results []map[string]interface{}
	
	for _, city := range s.cities {
		for _, district := range city.Districts {
			for _, village := range district.Villages {
				// Simple case-insensitive search
				if contains(village.Name, query) || contains(district.Name, query) {
					results = append(results, map[string]interface{}{
						"village_id":    village.ID,
						"village_name":  village.Name,
						"district_id":   district.ID,
						"district_name": district.Name,
						"city_id":       city.ID,
						"city_name":     city.Name,
						"full_address":  village.Name + ", " + district.Name + ", " + city.Name,
					})
				}
			}
		}
	}
	
	return c.JSON(model.Response{
		Success: true,
		Message: "Villages found",
		Data:    results,
	})
}

// Helper function for case-insensitive search
func contains(str, substr string) bool {
	str = toLower(str)
	substr = toLower(substr)
	return len(str) >= len(substr) && (str == substr || len(str) > len(substr) && hasSubstr(str, substr))
}

func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

func hasSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}