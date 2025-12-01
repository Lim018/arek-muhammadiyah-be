package route

import (
	"arek-muhammadiyah-be/app/service"
	"arek-muhammadiyah-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupTicketRoutes(app *fiber.App) {
	ticketService := service.NewTicketService()
	tickets := app.Group("/api/tickets", middleware.Authorization())
	
	tickets.Get("/", middleware.AdminOnly(), ticketService.GetAllTickets())
	// Get user tickets  
	tickets.Get("/my", ticketService.GetUserTickets())
	// Get ticket statistics (admin only)
	tickets.Get("/stats", middleware.AdminOnly(), ticketService.GetTicketStats())
	// Get ticket by ID
	tickets.Get("/:id", ticketService.GetTicketByID())
	// Create ticket - dengan atau tanpa files
	tickets.Post("/", ticketService.CreateTicketWithFiles())
	// Upload files to existing ticket
	tickets.Post("/:id/upload-files", ticketService.UploadFilesToTicket())
	// Get ticket files
	tickets.Get("/:id/files", ticketService.GetTicketFiles())
	// Delete file from ticket
	tickets.Delete("/:id/files/:fileId", ticketService.DeleteFile())
	// Update ticket (admin only)
	tickets.Put("/:id", middleware.AdminOnly(), ticketService.UpdateTicket())
	// Delete ticket (admin only) 
	tickets.Delete("/:id", ticketService.DeleteTicket())
	// Serve file untuk download (public route)
	app.Get("/api/files/:fileId", ticketService.ServeFile())
}