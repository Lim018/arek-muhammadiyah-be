package service

import (
	"errors"
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
)

type ArticleService struct {
	articleRepo  *repository.ArticleRepository
	documentRepo *repository.DocumentRepository // documentRepo bisa disimpan jika dipakai servis lain
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		articleRepo:  repository.NewArticleRepository(),
		documentRepo: repository.NewDocumentRepository(),
	}
}

func (s *ArticleService) GetAllArticles(page, limit int, published *bool) ([]model.Article, model.Pagination, error) {
	offset := (page - 1) * limit
	articles, total, err := s.articleRepo.GetAll(limit, offset, published)
	if err != nil {
		return nil, model.Pagination{}, err
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)
	return articles, pagination, nil
}

func (s *ArticleService) GetArticleByID(id uint) (*model.Article, error) {
	return s.articleRepo.GetByID(id)
}

func (s *ArticleService) GetArticleBySlug(slug string) (*model.Article, error) {
	return s.articleRepo.GetBySlug(slug)
}

func (s *ArticleService) CreateArticle(userID string, req *model.CreateArticleRequest) (*model.Article, error) {
	slug := helper.GenerateSlug(req.Title)

	_, err := s.articleRepo.GetBySlug(slug)
	if err == nil {
		slug = helper.GenerateUniqueSlug(req.Title)
	}

	// 1. Siapkan slice dokumen dari request
	var documents []model.Document
	for _, docReq := range req.Documents {
		documents = append(documents, model.Document{
			Title:       docReq.Title,
			Description: docReq.Description,
			FilePath:    docReq.FilePath,
			FileName:    docReq.FileName,
			FileSize:    docReq.FileSize,
			MimeType:    docReq.MimeType,
			// ArticleID akan diisi oleh repository
		})
	}

	// 2. Siapkan objek article lengkap
	article := &model.Article{
		UserID:       userID,
		CategoryID:   req.CategoryID,
		Title:        req.Title,
		Slug:         slug,
		Content:      req.Content,
		FeatureImage: req.FeatureImage,
		IsPublished:  helper.GetBoolValue(req.IsPublished, false),
		Documents:    documents, // <-- Masukkan dokumen ke sini
	}

	// 3. Panggil repo.Create HANYA SEKALI
	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	// 4. HAPUS loop "for _, docReq := ..." yang lama

	return s.articleRepo.GetByID(article.ID)
}

func (s *ArticleService) UpdateArticle(id uint, req *model.CreateArticleRequest) (*model.Article, error) {
	existing, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	var slug string
	if req.Title != existing.Title {
		slug = helper.GenerateSlug(req.Title)
		existingBySlug, _ := s.articleRepo.GetBySlug(slug)
		if existingBySlug != nil && existingBySlug.ID != id {
			slug = helper.GenerateUniqueSlug(req.Title)
		}
	} else {
		slug = existing.Slug
	}

	// 1. Siapkan data update untuk artikel
	updateData := &model.Article{
		CategoryID:   helper.GetUintPointer(req.CategoryID, existing.CategoryID),
		Title:        helper.GetStringValue(&req.Title, existing.Title),
		Slug:         slug,
		Content:      helper.GetStringValue(&req.Content, existing.Content),
		FeatureImage: helper.GetStringPointer(req.FeatureImage, existing.FeatureImage),
		IsPublished:  helper.GetBoolValue(req.IsPublished, existing.IsPublished),
	}

	// 2. Siapkan dokumen BARU yang akan ditambahkan
	// (Logika ini hanya menambah, tidak menghapus/mengedit dokumen lama)
	var newDocuments []model.Document
	for _, docReq := range req.Documents {
		newDocuments = append(newDocuments, model.Document{
			Title:       docReq.Title,
			Description: docReq.Description,
			FilePath:    docReq.FilePath,
			FileName:    docReq.FileName,
			FileSize:    docReq.FileSize,
			MimeType:    docReq.MimeType,
			// ArticleID akan diisi oleh repository
		})
	}

	// 3. Panggil repo.Update HANYA SEKALI
	// Kita perlu mengubah signature repo.Update untuk menerima dokumen baru
	if err := s.articleRepo.Update(id, updateData, newDocuments); err != nil {
		return nil, err
	}

	// 4. HAPUS loop "for _, docReq := ..." yang lama

	return s.articleRepo.GetByID(id)
}

func (s *ArticleService) DeleteArticle(id uint) error {
	// Panggil repo.Delete, yang sudah transaksional
	_, err := s.articleRepo.GetByID(id)
	if err != nil {
		return errors.New("article not found")
	}
	
	// Tidak perlu loop untuk hapus dokumen, repo.Delete akan menanganinya
	return s.articleRepo.Delete(id)
}

func (s *ArticleService) GetArticlesByCategory(categoryID uint, page, limit int) ([]model.Article, model.Pagination, error) {
	offset := (page - 1) * limit

	articles, err := s.articleRepo.GetByCategory(categoryID, limit, offset)
	if err != nil {
		return nil, model.Pagination{}, err
	}

	total := int64(len(articles))
	if total == 0 {
		return []model.Article{}, helper.CreatePagination(int64(page), int64(limit), 0), nil
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)

	return articles, pagination, nil
}
