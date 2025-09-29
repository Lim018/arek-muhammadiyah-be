package route

import (
	"github.com/Lim018/arek-muhammadiyah-be/app/model"
	"github.com/Lim018/arek-muhammadiyah-be/app/service"
	"github.com/Lim018/arek-muhammadiyah-be/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SetupTicketRoutes(app *fiber.App) {
	tickets := app.Group("/api/tickets")
	ticketService := service.NewTicketService()

	// Protected routes
	tickets.Use(middleware.JWTMiddleware())

	// Get all tickets (admin only)
	tickets.Get("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		statusStr := c.Query("status")
		
		var status *model.TicketStatus
		if statusStr != "" {
			s := model.TicketStatus(statusStr)
			status = &s
		}

		tickets, pagination, err := ticketService.GetAllTickets(page, limit, status)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.PaginatedResponse{
			Success:    true,
			Message:    "Tickets retrieved successfully",
			Data:       tickets,
			Pagination: pagination,
		})
	})

	// Get user tickets
	tickets.Get("/my", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		tickets, pagination, err := ticketService.GetUserTickets(userID, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.PaginatedResponse{
			Success:    true,
			Message:    "User tickets retrieved successfully",
			Data:       tickets,
			Pagination: pagination,
		})
	})

	// Get ticket statistics (admin only)
	tickets.Get("/stats", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		stats, err := ticketService.GetTicketStats()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Ticket statistics retrieved successfully",
			Data:    stats,
		})
	})

	tickets.Get("/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		ticket, err := ticketService.GetTicketByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Success: false,
				Message: "Ticket not found",
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Ticket retrieved successfully",
			Data:    ticket,
		})
	})

	tickets.Post("/", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		var req model.CreateTicketRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		ticket, err := ticketService.CreateTicket(userID, &req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(model.Response{
			Success: true,
			Message: "Ticket created successfully",
			Data:    ticket,
		})
	})

	tickets.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		var req model.UpdateTicketRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		ticket, err := ticketService.UpdateTicket(uint(id), &req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Ticket updated successfully",
			Data:    ticket,
		})
	})

	tickets.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		err := ticketService.DeleteTicket(uint(id))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(model.Response{
			Success: true,
			Message: "Ticket deleted successfully",
		})
	})
}