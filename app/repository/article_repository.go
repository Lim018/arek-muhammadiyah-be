package repository

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/database"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{
		db: database.DB,
	}
}

func (r *ArticleRepository) GetAll(limit, offset int, published *bool) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.Model(&model.Article{})
	if published != nil {
		query = query.Where("is_published = ?", *published)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("User").Preload("Category").
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&articles).Error

	return articles, total, err
}

func (r *ArticleRepository) GetByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("User").Preload("Category").
		First(&article, id).Error
	return &article, err
}

func (r *ArticleRepository) GetBySlug(slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("User").Preload("Category").
		Where("slug = ?", slug).First(&article).Error
	return &article, err
}

func (r *ArticleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

func (r *ArticleRepository) Update(id uint, article *model.Article) error {
	return r.db.Where("id = ?", id).Updates(article).Error
}

func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Article{}, id).Error
}

func (r *ArticleRepository) GetByCategory(categoryID uint, limit, offset int) ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Preload("User").Preload("Category").
		Where("category_id = ? AND is_published = ?", categoryID, true).
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&articles).Error
	return articles, err
}