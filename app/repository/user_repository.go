package repository

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/database"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

func (r *UserRepository) GetAll(limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").Preload("Village").
		Limit(limit).Offset(offset).Find(&users).Error

	return users, total, err
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").Preload("Village").
		First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(id string, user *model.User) error {
	return r.db.Where("id = ?", id).Updates(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

func (r *UserRepository) GetByVillage(villageID uint, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).Where("village_id = ?", villageID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").Where("village_id = ?", villageID).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) GetWithStats(limit, offset int) ([]model.UserWithStats, int64, error) {
	var users []model.UserWithStats
	var total int64

	err := r.db.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT u.*, 
			   COALESCE(article_count, 0) as total_articles,
			   COALESCE(ticket_count, 0) as total_tickets,
			   COALESCE(document_count, 0) as total_documents
		FROM users u
		LEFT JOIN (
			SELECT user_id, COUNT(*) as article_count 
			FROM articles 
			GROUP BY user_id
		) a ON u.id = a.user_id
		LEFT JOIN (
			SELECT user_id, COUNT(*) as ticket_count 
			FROM tickets 
			GROUP BY user_id
		) t ON u.id = t.user_id
		LEFT JOIN (
			SELECT user_id, COUNT(*) as document_count 
			FROM documents 
			GROUP BY user_id
		) d ON u.id = d.user_id
		LIMIT ? OFFSET ?
	`

	err = r.db.Raw(query, limit, offset).Scan(&users).Error
	return users, total, err
}

func (r *UserRepository) GetByCardStatus(status string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).Where("card_status = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").Preload("Village").
		Where("card_status = ?", status).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) GetCardStatusStats() (map[string]int64, error) {
	var results []struct {
		CardStatus string
		Count      int64
	}

	err := r.db.Model(&model.User{}).
		Select("card_status, count(*) as count").
		Group("card_status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.CardStatus] = result.Count
	}

	return stats, nil
}

func (r *UserRepository) GetMobileUsers(limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).Where("is_mobile = ?", true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").Preload("Village").
		Where("is_mobile = ?", true).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}