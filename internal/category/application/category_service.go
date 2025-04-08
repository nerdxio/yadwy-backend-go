package application

import (
	"context"
	"go.uber.org/zap"
	"mime/multipart"
	"yadwy-backend/internal/category/domain"
	"yadwy-backend/internal/common"
)

type CategoryService struct {
	repo   domain.CategoryRepo
	files  *common.FileService
	logger *zap.Logger
}

func NewCategoryService(repo domain.CategoryRepo, files *common.FileService, logger *zap.Logger) *CategoryService {
	return &CategoryService{
		repo:   repo,
		files:  files,
		logger: logger,
	}
}

func (s *CategoryService) CreateCategory(name, description string, image *multipart.FileHeader) error {
	url, err := s.files.SaveFile(image)
	if err != nil {
		return common.NewErrorf(domain.FailedToUploadImage, "Failed to upload image: %v", err)
	}

	s.logger.Info("Image uploaded successfully", zap.String("url", url))
	p := domain.Params{
		Name:        name,
		Description: description,
		ImageUrl:    url,
	}
	category := domain.NewCategory(p)

	s.logger.Info("Creating category in repository",
		zap.String("name", category.Name()),
		zap.String("description", category.Description()),
		zap.String("image", category.ImageUrl()))

	_, err = s.repo.CreateCategory(name, description, url)
	if err != nil {
		s.logger.Error("Failed to create category", zap.Error(err))
		return err
	}

	return nil
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := s.repo.GetAllCategories(ctx)
	if err != nil {
		return nil, common.NewErrorf(domain.FailedToReadCategories, "Failed to read categories: %v", err)
	}

	s.logger.Info("Successfully retrieved categories", zap.Int("count", len(categories)))
	return categories, nil
}
