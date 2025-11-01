package repository

import (
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/database"
	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{
		db: database.DB,
	}
}

func (r *DocumentRepository) GetAll(limit, offset int) ([]model.Document, int64, error) {
	var documents []model.Document
	var total int64

	if err := r.db.Model(&model.Document{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Preload("Ticket").
		Preload("Article").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&documents).Error

	return documents, total, err
}

func (r *DocumentRepository) GetByID(id uint) (*model.Document, error) {
	var document model.Document
	err := r.db.
		Preload("Ticket").
		Preload("Article").
		First(&document, id).Error
	return &document, err
}

func (r *DocumentRepository) GetByTicketID(ticketID uint, limit, offset int) ([]model.Document, int64, error) {
	var documents []model.Document
	var total int64

	if err := r.db.Model(&model.Document{}).
		Where("ticket_id = ?", ticketID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Where("ticket_id = ?", ticketID).
		Preload("Ticket").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&documents).Error

	return documents, total, err
}

func (r *DocumentRepository) GetByArticleID(articleID uint, limit, offset int) ([]model.Document, int64, error) {
	var documents []model.Document
	var total int64

	if err := r.db.Model(&model.Document{}).
		Where("article_id = ?", articleID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Where("article_id = ?", articleID).
		Preload("Article").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&documents).Error

	return documents, total, err
}

func (r *DocumentRepository) Create(document *model.Document) error {
	return r.db.Create(document).Error
}

func (r *DocumentRepository) Update(id uint, document *model.Document) error {
	return r.db.Where("id = ?", id).Updates(document).Error
}

func (r *DocumentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Document{}, id).Error
}

