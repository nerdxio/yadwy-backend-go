package application

import (
	"context"
	"go.uber.org/zap"
	"mime/multipart"
	"strings"
	"yadwy-backend/internal/common"
	"yadwy-backend/internal/prodcuts/domain"
)

const (
	FailedToUploadImage     = "failed-to-upload-image"
	FailedToCreateProduct   = "failed-to-create-product"
	FailedToRetrieveProduct = "failed-to-retrieve-product"
)

type ProductService struct {
	repo   domain.ProductRepository // Changed from pointer to interface
	files  *common.FileService
	logger *zap.Logger
}

func NewProductService(repo domain.ProductRepository, files *common.FileService, logger *zap.Logger) *ProductService {
	return &ProductService{
		repo:   repo,
		files:  files,
		logger: logger,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, p *domain.Product, images []*multipart.FileHeader) error {
	var savedImages []domain.Image

	for _, img := range images {
		// Extract image type from filename
		parts := strings.Split(img.Filename, ":")
		imageType := "main" // default type
		if len(parts) == 2 {
			imageType = parts[0]
			img.Filename = parts[1]
		}

		url, err := s.files.SaveFile(img)
		if err != nil {
			return common.NewErrorf(FailedToUploadImage, "failed to upload image: %v", err)
		}

		savedImages = append(savedImages, domain.Image{
			URL:  url,
			Type: imageType,
		})

		s.logger.Info("Image uploaded successfully",
			zap.String("url", url),
			zap.String("type", imageType))
	}

	err := s.repo.CreateProduct(ctx, p, savedImages)
	if err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return common.NewErrorf(FailedToCreateProduct, "failed to create product: %v", err)
	}

	s.logger.Info("Product created successfully",
		zap.String("name", p.Name),
		zap.Int("imagesCount", len(savedImages)))

	return nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve product", zap.Error(err))
		return nil, common.NewErrorf(FailedToRetrieveProduct, "failed to retrieve product: %v", err)
	}

	s.logger.Info("Product retrieved successfully",
		zap.Int64("id", id),
		zap.Int("imagesCount", len(product.Images)))

	return product, nil
}
