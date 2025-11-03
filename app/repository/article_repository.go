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

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("User").
		Preload("Category").
		Preload("Documents").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error

	return articles, total, err
}

func (r *ArticleRepository) GetByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Documents").
		First(&article, id).Error
	return &article, err
}

func (r *ArticleRepository) GetBySlug(slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Documents").
		Where("slug = ?", slug).
		First(&article).Error
	return &article, err
}

func (r *ArticleRepository) Create(article *model.Article) error {
	tx := r.db.Begin()

	// 1. Buat artikel (ID akan diisi GORM)
	if err := tx.Create(article).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Jika ada dokumen, set ArticleID-nya
	if len(article.Documents) > 0 {
		for i := range article.Documents {
			article.Documents[i].ArticleID = &article.ID // <-- PERBAIKAN: Gunakan pointer
		}
		// Buat dokumen dalam transaksi yang sama
		if err := tx.Create(&article.Documents).Error; err != nil {
			tx.Rollback() // Batalkan pembuatan artikel jika dokumen gagal
			return err
		}
	}

	return tx.Commit().Error
}

func (r *ArticleRepository) Update(id uint, articleData *model.Article, newDocuments []model.Document) error {
	// MULAI TRANSAKSI
	tx := r.db.Begin()

	// 1. Update data utama artikel
	if err := tx.Model(&model.Article{}).Where("id = ?", id).Updates(articleData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Jika ada dokumen BARU, buat dalam transaksi yang sama
	if len(newDocuments) > 0 {
		for i := range newDocuments {
			newDocuments[i].ArticleID = &id // <-- PERBAIKAN: Set ID & gunakan pointer
		}
		if err := tx.Create(&newDocuments).Error; err != nil {
			tx.Rollback() // Batalkan update artikel jika dokumen baru gagal
			return err
		}
	}

	// SELESAI TRANSAKSI
	return tx.Commit().Error
}

func (r *ArticleRepository) Delete(id uint) error {
	tx := r.db.Begin()
	
	// PERBAIKAN: Hapus dokumen menggunakan GORM dan pointer
	// Ini akan menghapus semua dokumen di mana article_id = id
	if err := tx.Where(&model.Document{ArticleID: &id}).Delete(&model.Document{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Hapus artikel itu sendiri
	if err := tx.Delete(&model.Article{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *ArticleRepository) GetByCategory(categoryID uint, limit, offset int) ([]model.Article, error) {
	var articles []model.Article
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Documents").
		Where("category_id = ? AND is_published = ?", categoryID, true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error
	return articles, err
}