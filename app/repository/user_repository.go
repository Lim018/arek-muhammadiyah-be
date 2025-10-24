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

	err = r.db.Preload("Role").Preload("SubVillage").Preload("SubVillage.Village").
		Limit(limit).Offset(offset).Find(&users).Error

	return users, total, err
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").Preload("SubVillage").Preload("SubVillage.Village").
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

func (r *UserRepository) GetBySubVillage(subVillageID uint, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).Where("sub_village_id = ?", subVillageID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").Preload("SubVillage").Preload("SubVillage.Village").
		Where("sub_village_id = ?", subVillageID).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) GetByVillage(villageID uint, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).
		Joins("JOIN sub_villages ON users.sub_village_id = sub_villages.id").
		Where("sub_villages.village_id = ?", villageID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").Preload("SubVillage").Preload("SubVillage.Village").
		Joins("JOIN sub_villages ON users.sub_village_id = sub_villages.id").
		Where("sub_villages.village_id = ?", villageID).
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
	
	// Calculate age for each user
	for i := range users {
		users[i].Age = users[i].GetAge()
	}
	
	return users, total, err
}

func (r *UserRepository) GetMobileUsers(limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).Where("is_mobile = ?", true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").Preload("SubVillage").Preload("SubVillage.Village").
		Where("is_mobile = ?", true).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) GetByGender(gender string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).Where("gender = ?", gender).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").Preload("SubVillage").Preload("SubVillage.Village").
		Where("gender = ?", gender).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}