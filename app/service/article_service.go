package service

import (
	"errors"
	"arek-muhammadiyah-be/app/model"
	"arek-muhammadiyah-be/app/repository"
	"arek-muhammadiyah-be/helper"
)

type ArticleService struct {
	articleRepo  *repository.ArticleRepository
	documentRepo *repository.DocumentRepository
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

	article := &model.Article{
		UserID:       userID,
		CategoryID:   req.CategoryID,
		Title:        req.Title,
		Slug:         slug,
		Content:      req.Content,
		FeatureImage: req.FeatureImage,
		IsPublished:  helper.GetBoolValue(req.IsPublished, false),
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	for _, docReq := range req.Documents {
		doc := &model.Document{
			ArticleID:   article.ID,
			Title:       docReq.Title,
			Description: docReq.Description,
			FilePath:    docReq.FilePath,
			FileName:    docReq.FileName,
			FileSize:    docReq.FileSize,
			MimeType:    docReq.MimeType,
		}
		_ = s.documentRepo.Create(doc) 
	}

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

	updateData := &model.Article{
		CategoryID:   helper.GetUintPointer(req.CategoryID, existing.CategoryID),
		Title:        helper.GetStringValue(&req.Title, existing.Title),
		Slug:         slug,
		Content:      helper.GetStringValue(&req.Content, existing.Content),
		FeatureImage: helper.GetStringPointer(req.FeatureImage, existing.FeatureImage),
		IsPublished:  helper.GetBoolValue(req.IsPublished, existing.IsPublished),
	}

	if err := s.articleRepo.Update(id, updateData); err != nil {
		return nil, err
	}

	if len(req.Documents) > 0 {
		for _, docReq := range req.Documents {
			doc := &model.Document{
				ArticleID:   id,
				Title:       docReq.Title,
				Description: docReq.Description,
				FilePath:    docReq.FilePath,
				FileName:    docReq.FileName,
				FileSize:    docReq.FileSize,
				MimeType:    docReq.MimeType,
			}
			_ = s.documentRepo.Create(doc)
		}
	}

	return s.articleRepo.GetByID(id)
}

func (s *ArticleService) DeleteArticle(id uint) error {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return errors.New("article not found")
	}
	docs, _, _ := s.documentRepo.GetByArticleID(article.ID, 1000, 0)
	for _, doc := range docs {
		if err := s.documentRepo.Delete(doc.ID); err != nil {
			return err
		}
	}

	return s.articleRepo.Delete(id)
}
