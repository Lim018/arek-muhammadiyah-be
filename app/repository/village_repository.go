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
	err := r.db.Preload("SubVillages").First(&village, id).Error
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
			SELECT sv.village_id, COUNT(u.id) as user_count 
			FROM users u
			JOIN sub_villages sv ON u.sub_village_id = sv.id
			WHERE u.sub_village_id IS NOT NULL
			GROUP BY sv.village_id
		) u ON v.id = u.village_id
		WHERE v.is_active = true
		ORDER BY v.name ASC
	`

	err := r.db.Raw(query).Scan(&villages).Error
	return villages, err
}

func (r *VillageRepository) GetWithCompleteStats() ([]model.VillageWithStats, error) {
	var villages []model.VillageWithStats

	query := `
		SELECT 
			v.id, 
			v.name, 
			v.code, 
			v.color,
			COALESCE(members, 0) as members,
			COALESCE(tickets, 0) as tickets,
			COALESCE(articles, 0) as articles,
			COALESCE(app_users, 0) as app_users,
			COALESCE(sub_village_count, 0) as sub_villages
		FROM villages v
		LEFT JOIN (
			SELECT sv.village_id, COUNT(u.id) as members
			FROM users u
			JOIN sub_villages sv ON u.sub_village_id = sv.id
			WHERE u.sub_village_id IS NOT NULL
			GROUP BY sv.village_id
		) u ON v.id = u.village_id
		LEFT JOIN (
			SELECT sv.village_id, COUNT(t.id) as tickets
			FROM tickets t
			JOIN users u ON t.user_id = u.id
			JOIN sub_villages sv ON u.sub_village_id = sv.id
			WHERE u.sub_village_id IS NOT NULL
			GROUP BY sv.village_id
		) t ON v.id = t.village_id
		LEFT JOIN (
			SELECT sv.village_id, COUNT(a.id) as articles
			FROM articles a
			JOIN users u ON a.user_id = u.id
			JOIN sub_villages sv ON u.sub_village_id = sv.id
			WHERE u.sub_village_id IS NOT NULL
			GROUP BY sv.village_id
		) a ON v.id = a.village_id
		LEFT JOIN (
			SELECT sv.village_id, COUNT(u.id) as app_users
			FROM users u
			JOIN sub_villages sv ON u.sub_village_id = sv.id
			WHERE u.sub_village_id IS NOT NULL AND u.is_mobile = true
			GROUP BY sv.village_id
		) au ON v.id = au.village_id
		LEFT JOIN (
			SELECT village_id, COUNT(*) as sub_village_count
			FROM sub_villages
			WHERE is_active = true
			GROUP BY village_id
		) sv ON v.id = sv.village_id
		WHERE v.is_active = true
		ORDER BY v.name ASC
	`

	err := r.db.Raw(query).Scan(&villages).Error
	return villages, err
}