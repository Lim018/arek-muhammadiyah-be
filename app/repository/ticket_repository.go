package repository

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/database"

	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{
		db: database.DB,
	}
}

func (r *TicketRepository) GetAll(limit, offset int, status *model.TicketStatus) ([]model.Ticket, int64, error) {
	var tickets []model.Ticket
	var total int64

	query := r.db.Model(&model.Ticket{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("User").
		Preload("Category").
		Preload("Documents"). 
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tickets).Error

	return tickets, total, err
}

func (r *TicketRepository) GetByID(id uint) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Documents"). 
		First(&ticket, id).Error
	return &ticket, err
}

func (r *TicketRepository) GetByUserID(userID string, limit, offset int) ([]model.Ticket, int64, error) {
	var tickets []model.Ticket
	var total int64

	if err := r.db.Model(&model.Ticket{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Preload("Category").
		Preload("Documents"). 
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tickets).Error

	return tickets, total, err
}

func (r *TicketRepository) Create(ticket *model.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) Update(id uint, ticket *model.Ticket) error {
	return r.db.Where("id = ?", id).Updates(ticket).Error
}

func (r *TicketRepository) Delete(id uint) error {
	return r.db.Delete(&model.Ticket{}, id).Error
}

func (r *TicketRepository) GetCountByStatus() (map[model.TicketStatus]int64, error) {
	var results []struct {
		Status model.TicketStatus
		Count  int64
	}

	err := r.db.Model(&model.Ticket{}).
		Select("status, COUNT(*) AS count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[model.TicketStatus]int64)
	for _, result := range results {
		counts[result.Status] = result.Count
	}

	return counts, nil
}