package repository

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/database"
	"gorm.io/gorm"
)

type SubVillageRepository struct {
	db *gorm.DB
}

func NewSubVillageRepository() *SubVillageRepository {
	return &SubVillageRepository{
		db: database.DB,
	}
}

func (r *SubVillageRepository) GetAll(limit, offset int, activeOnly bool) ([]model.SubVillage, int64, error) {
	var subVillages []model.SubVillage
	var total int64

	query := r.db.Model(&model.SubVillage{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Village").Order("name ASC").
		Limit(limit).Offset(offset).Find(&subVillages).Error

	return subVillages, total, err
}

func (r *SubVillageRepository) GetByID(id uint) (*model.SubVillage, error) {
	var subVillage model.SubVillage
	err := r.db.Preload("Village").First(&subVillage, id).Error
	return &subVillage, err
}

func (r *SubVillageRepository) GetByVillageID(villageID uint, limit, offset int) ([]model.SubVillage, int64, error) {
	var subVillages []model.SubVillage
	var total int64

	err := r.db.Model(&model.SubVillage{}).Where("village_id = ?", villageID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Village").
		Where("village_id = ?", villageID).
		Order("name ASC").
		Limit(limit).Offset(offset).Find(&subVillages).Error
	return subVillages, total, err
}

func (r *SubVillageRepository) Create(subVillage *model.SubVillage) error {
	return r.db.Create(subVillage).Error
}

func (r *SubVillageRepository) Update(id uint, subVillage *model.SubVillage) error {
	return r.db.Where("id = ?", id).Updates(subVillage).Error
}

func (r *SubVillageRepository) Delete(id uint) error {
	return r.db.Delete(&model.SubVillage{}, id).Error
}

func (r *SubVillageRepository) GetWithUserCount() ([]model.SubVillageWithUserCount, error) {
	var subVillages []model.SubVillageWithUserCount

	query := `
		SELECT sv.*, COALESCE(user_count, 0) as total_users
		FROM sub_villages sv
		LEFT JOIN (
			SELECT sub_village_id, COUNT(*) as user_count 
			FROM users 
			WHERE sub_village_id IS NOT NULL
			GROUP BY sub_village_id
		) u ON sv.id = u.sub_village_id
		WHERE sv.is_active = true
		ORDER BY sv.name ASC
	`

	err := r.db.Raw(query).Scan(&subVillages).Error
	return subVillages, err
}

func (r *SubVillageRepository) GetWithCompleteStats() ([]model.SubVillageWithStats, error) {
	var subVillages []model.SubVillageWithStats

	query := `
		SELECT 
			sv.id, 
			sv.village_id,
			sv.name, 
			sv.code,
			COALESCE(members, 0) as members,
			COALESCE(tickets, 0) as tickets,
			COALESCE(articles, 0) as articles,
			COALESCE(app_users, 0) as app_users
		FROM sub_villages sv
		LEFT JOIN (
			SELECT sub_village_id, COUNT(*) as members
			FROM users
			WHERE sub_village_id IS NOT NULL
			GROUP BY sub_village_id
		) u ON sv.id = u.sub_village_id
		LEFT JOIN (
			SELECT u.sub_village_id, COUNT(t.id) as tickets
			FROM tickets t
			JOIN users u ON t.user_id = u.id
			WHERE u.sub_village_id IS NOT NULL
			GROUP BY u.sub_village_id
		) t ON sv.id = t.sub_village_id
		LEFT JOIN (
			SELECT u.sub_village_id, COUNT(a.id) as articles
			FROM articles a
			JOIN users u ON a.user_id = u.id
			WHERE u.sub_village_id IS NOT NULL
			GROUP BY u.sub_village_id
		) a ON sv.id = a.sub_village_id
		LEFT JOIN (
			SELECT sub_village_id, COUNT(*) as app_users
			FROM users
			WHERE sub_village_id IS NOT NULL AND is_mobile = true
			GROUP BY sub_village_id
		) au ON sv.id = au.sub_village_id
		WHERE sv.is_active = true
		ORDER BY sv.name ASC
	`

	err := r.db.Raw(query).Scan(&subVillages).Error
	return subVillages, err
}