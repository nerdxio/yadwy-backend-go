package category

import (
	"go.uber.org/zap"
	"mime/multipart"
	"yadwy-backend/internal/common"
)

type CategoryService struct {
	repo   CategoryRepo
	files  *common.FileService
	logger *zap.Logger
}

func NewCategoryService(repo CategoryRepo, files *common.FileService, logger *zap.Logger) *CategoryService {
	return &CategoryService{
		repo:   repo,
		files:  files,
		logger: logger,
	}
}

func (s *CategoryService) Execute(name, description string, image *multipart.FileHeader) error {
	s.logger.Info("Creating category")
	file, err := s.files.SaveFile(image)
	if err != nil {
		return common.NewErrorf(FailedToUploadImage, "Failed to upload image: %v", err)
	}

	s.logger.Info("Image uploaded successfully", zap.String("file", file))
	category := NewCategory(name, description, file)

	s.logger.Info("Creating category in repository",
		zap.String("name", category.name),
		zap.String("description", category.description),
		zap.String("image", category.imageUrl))

	_, err = s.repo.CreateCategory(name, description, file)
	if err != nil {
		s.logger.Error("Failed to create category", zap.Error(err))
		return err
	}

	return nil
}
