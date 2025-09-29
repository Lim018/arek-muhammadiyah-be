package service

import (
	"errors"
	"github.com/Lim018/arek-muhammadiyah-be/app/model"
	"github.com/Lim018/arek-muhammadiyah-be/app/repository"
	"github.com/Lim018/arek-muhammadiyah-be/helper"
	// "strings"
)

type ArticleService struct {
	articleRepo *repository.ArticleRepository
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		articleRepo: repository.NewArticleRepository(),
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
	
	// Check if slug already exists
	_, err := s.articleRepo.GetBySlug(slug)
	if err == nil {
		slug = helper.GenerateUniqueSlug(req.Title)
	}

	article := &model.Article{
		UserID:        userID,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Slug:          slug,
		Content:       req.Content,
		FeaturedImage: req.FeaturedImage,
		IsPublished:   helper.GetBoolValue(req.IsPublished, false),
	}

	err = s.articleRepo.Create(article)
	if err != nil {
		return nil, err
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
		// Check if new slug already exists
		existingBySlug, _ := s.articleRepo.GetBySlug(slug)
		if existingBySlug != nil && existingBySlug.ID != id {
			slug = helper.GenerateUniqueSlug(req.Title)
		}
	} else {
		slug = existing.Slug
	}

	updateData := &model.Article{
		CategoryID:    helper.GetUintPointer(req.CategoryID, existing.CategoryID),
		Title:         helper.GetStringValue(&req.Title, existing.Title),
		Slug:          slug,
		Content:       helper.GetStringValue(&req.Content, existing.Content),
		FeaturedImage: helper.GetStringPointer(req.FeaturedImage, existing.FeaturedImage),
		IsPublished:   helper.GetBoolValue(req.IsPublished, existing.IsPublished),
	}

	err = s.articleRepo.Update(id, updateData)
	if err != nil {
		return nil, err
	}

	return s.articleRepo.GetByID(id)
}

func (s *ArticleService) DeleteArticle(id uint) error {
	_, err := s.articleRepo.GetByID(id)
	if err != nil {
		return errors.New("article not found")
	}

	return s.articleRepo.Delete(id)
}