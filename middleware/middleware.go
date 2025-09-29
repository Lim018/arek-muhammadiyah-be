package middleware

import (
	"github.com/Lim018/arek-muhammadiyah-be/config"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Setup logger
	config.SetupLogger()

	// Security middleware
	app.Use(Helmet())
	app.Use(CORS())
	app.Use(Logger())
}