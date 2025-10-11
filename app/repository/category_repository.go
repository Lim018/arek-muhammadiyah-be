package repository

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/database"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		db: database.DB,
	}
}

func (r *CategoryRepository) GetAll(limit, offset int, activeOnly bool) ([]model.Category, int64, error) {
	var categories []model.Category
	var total int64

	query := r.db.Model(&model.Category{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("name ASC").
		Limit(limit).Offset(offset).Find(&categories).Error

	return categories, total, err
}

func (r *CategoryRepository) GetByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) Update(id uint, category *model.Category) error {
	return r.db.Where("id = ?", id).Updates(category).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}