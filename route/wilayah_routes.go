package route

import (
	"arek-muhammadiyah-be/app/service"
	"github.com/gofiber/fiber/v2"
)

func SetupWilayahRoutes(app *fiber.App, wilayahService *service.WilayahService) {
	wilayah := app.Group("/api/wilayah")

	wilayah.Get("/cities", wilayahService.GetAllCities)
	wilayah.Get("/cities/:cityId/districts", wilayahService.GetDistricts)
	wilayah.Get("/cities/:cityId/districts/:districtId/villages", wilayahService.GetVillages)
	wilayah.Get("/search", wilayahService.SearchVillages)
}