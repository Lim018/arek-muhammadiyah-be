package repository

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/database"
	"gorm.io/gorm"
)

type VillageRepository struct {
	db *gorm.DB
}

func NewVillageRepository() *VillageRepository {
	return &VillageRepository{
		db: database.DB,
	}
}

func (r *VillageRepository) GetAll(limit, offset int, activeOnly bool) ([]model.Village, int64, error) {
	var villages []model.Village
	var total int64

	query := r.db.Model(&model.Village{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("name ASC").
		Limit(limit).Offset(offset).Find(&villages).Error

	return villages, total, err
}

func (r *VillageRepository) GetByID(id uint) (*model.Village, error) {
	var village model.Village
	err := r.db.First(&village, id).Error
	return &village, err
}

func (r *VillageRepository) Create(village *model.Village) error {
	return r.db.Create(village).Error
}

func (r *VillageRepository) Update(id uint, village *model.Village) error {
	return r.db.Where("id = ?", id).Updates(village).Error
}

func (r *VillageRepository) Delete(id uint) error {
	return r.db.Delete(&model.Village{}, id).Error
}

func (r *VillageRepository) GetWithUserCount() ([]model.VillageWithUserCount, error) {
	var villages []model.VillageWithUserCount

	query := `
		SELECT v.*, COALESCE(user_count, 0) as total_users
		FROM villages v
		LEFT JOIN (
			SELECT village_id, COUNT(*) as user_count 
			FROM users 
			WHERE village_id IS NOT NULL
			GROUP BY village_id
		) u ON v.id = u.village_id
		WHERE v.is_active = true
		ORDER BY v.name ASC
	`

	err := r.db.Raw(query).Scan(&villages).Error
	return villages, err
}