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

	err = r.db.Preload("Role").
		Limit(limit).Offset(offset).Find(&users).Error

	return users, total, err
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").
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

func (r *UserRepository) GetMobileUsers(limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).Where("is_mobile = ?", true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").
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

	err = r.db.Preload("Role").
		Where("gender = ?", gender).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

// GetByCity - Get users by city_id (first 4 chars of village_id)
func (r *UserRepository) GetByCity(cityID string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// village_id format: "3576011001"
	// city_id: "3576" (first 4 chars)
	query := r.db.Model(&model.User{}).Where("LEFT(village_id, 4) = ?", cityID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").
		Where("LEFT(village_id, 4) = ?", cityID).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

// GetByDistrict - Get users by district_id (first 6 chars of village_id)
func (r *UserRepository) GetByDistrict(districtID string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// village_id format: "3576011001"
	// district_id: "357601" (first 6 chars)
	query := r.db.Model(&model.User{}).Where("LEFT(village_id, 6) = ?", districtID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").
		Where("LEFT(village_id, 6) = ?", districtID).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

// GetByVillage - Get users by village_id
func (r *UserRepository) GetByVillage(villageID string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := r.db.Model(&model.User{}).Where("village_id = ?", villageID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Role").
		Where("village_id = ?", villageID).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

// GetWilayahStats - Get statistics by wilayah
func (r *UserRepository) GetWilayahStats() (map[string]interface{}, error) {
	var results []struct {
		CityID string
		Total  int64
	}

	// Count by city (first 4 chars of village_id)
	err := r.db.Model(&model.User{}).
		Select("LEFT(village_id, 4) as city_id, COUNT(*) as total").
		Where("village_id IS NOT NULL AND village_id != ''").
		Group("city_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})
	for _, result := range results {
		stats[result.CityID] = result.Total
	}

	return stats, nil
}