package repository

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/database"
	"gorm.io/gorm"
)

type StatsRepository struct {
	db *gorm.DB
}

func NewStatsRepository() *StatsRepository {
	return &StatsRepository{
		db: database.DB,
	}
}

// GetTotalUsers - Get total users count
func (r *StatsRepository) GetTotalUsers() (int64, error) {
	var total int64
	err := r.db.Model(&model.User{}).Count(&total).Error
	return total, err
}

// GetTicketStats - Get ticket statistics by status
func (r *StatsRepository) GetTicketStats() (*model.TicketStats, error) {
	var stats model.TicketStats
	
	// Count by status
	r.db.Model(&model.Ticket{}).Where("status = ?", model.TicketStatusUnread).Count(&stats.Unread)
	r.db.Model(&model.Ticket{}).Where("status = ?", model.TicketStatusRead).Count(&stats.Read)
	r.db.Model(&model.Ticket{}).Where("status = ?", model.TicketStatusInProgress).Count(&stats.InProgress)
	r.db.Model(&model.Ticket{}).Where("status = ?", model.TicketStatusResolved).Count(&stats.Resolved)
	r.db.Model(&model.Ticket{}).Where("status = ?", model.TicketStatusRejected).Count(&stats.Rejected)
	
	// Total tickets
	err := r.db.Model(&model.Ticket{}).Count(&stats.Total).Error
	
	return &stats, err
}

// GetTotalUsersByCity - Get total users by city
func (r *StatsRepository) GetTotalUsersByCity(cityID string) (int64, error) {
	var total int64
	err := r.db.Model(&model.User{}).
		Where("LEFT(village_id, 4) = ?", cityID).
		Count(&total).Error
	return total, err
}

// GetTicketStatsByCity - Get ticket statistics by city
func (r *StatsRepository) GetTicketStatsByCity(cityID string) (*model.TicketStats, error) {
	var stats model.TicketStats
	
	// Subquery to get user IDs in the city
	userIDs := r.db.Model(&model.User{}).
		Select("id").
		Where("LEFT(village_id, 4) = ?", cityID)
	
	// Count tickets by status for users in this city
	r.db.Model(&model.Ticket{}).
		Where("user_id IN (?)", userIDs).
		Where("status = ?", model.TicketStatusUnread).
		Count(&stats.Unread)
	
	r.db.Model(&model.Ticket{}).
		Where("user_id IN (?)", userIDs).
		Where("status = ?", model.TicketStatusRead).
		Count(&stats.Read)
	
	r.db.Model(&model.Ticket{}).
		Where("user_id IN (?)", userIDs).
		Where("status = ?", model.TicketStatusInProgress).
		Count(&stats.InProgress)
	
	r.db.Model(&model.Ticket{}).
		Where("user_id IN (?)", userIDs).
		Where("status = ?", model.TicketStatusResolved).
		Count(&stats.Resolved)
	
	r.db.Model(&model.Ticket{}).
		Where("user_id IN (?)", userIDs).
		Where("status = ?", model.TicketStatusRejected).
		Count(&stats.Rejected)
	
	// Total tickets
	err := r.db.Model(&model.Ticket{}).
		Where("user_id IN (?)", userIDs).
		Count(&stats.Total).Error
	
	return &stats, err
}