package service

import (
	"errors"
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
	"time"
)

type TicketService struct {
	ticketRepo *repository.TicketRepository
}

func NewTicketService() *TicketService {
	return &TicketService{
		ticketRepo: repository.NewTicketRepository(),
	}
}

func (s *TicketService) GetAllTickets(page, limit int, status *model.TicketStatus) ([]model.Ticket, model.Pagination, error) {
	offset := (page - 1) * limit
	tickets, total, err := s.ticketRepo.GetAll(limit, offset, status)
	if err != nil {
		return nil, model.Pagination{}, err
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)
	return tickets, pagination, nil
}

func (s *TicketService) GetTicketByID(id uint) (*model.Ticket, error) {
	return s.ticketRepo.GetByID(id)
}

func (s *TicketService) GetUserTickets(userID string, page, limit int) ([]model.Ticket, model.Pagination, error) {
	offset := (page - 1) * limit
	tickets, total, err := s.ticketRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, model.Pagination{}, err
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)
	return tickets, pagination, nil
}

func (s *TicketService) CreateTicket(userID string, req *model.CreateTicketRequest) (*model.Ticket, error) {
	ticket := &model.Ticket{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		Status:      model.TicketStatusUnread,
	}

	err := s.ticketRepo.Create(ticket)
	if err != nil {
		return nil, err
	}

	return s.ticketRepo.GetByID(ticket.ID)
}

func (s *TicketService) UpdateTicket(id uint, req *model.UpdateTicketRequest) (*model.Ticket, error) {
	_, err := s.ticketRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	updateData := &model.Ticket{}
	
	if req.Status != nil {
		updateData.Status = *req.Status
		
		if *req.Status == model.TicketStatusResolved || *req.Status == model.TicketStatusRejected {
			now := time.Now()
			updateData.ResolvedAt = &now
		}
	}
	
	if req.Resolution != nil {
		updateData.Resolution = req.Resolution
	}

	err = s.ticketRepo.Update(id, updateData)
	if err != nil {
		return nil, err
	}

	return s.ticketRepo.GetByID(id)
}

func (s *TicketService) DeleteTicket(id uint) error {
	_, err := s.ticketRepo.GetByID(id)
	if err != nil {
		return errors.New("ticket not found")
	}

	return s.ticketRepo.Delete(id)
}

func (s *TicketService) GetTicketStats() (map[model.TicketStatus]int64, error) {
	return s.ticketRepo.GetCountByStatus()
}