package banner

import (
	"context"
	"go.uber.org/zap"
	"mime/multipart"
	"yadwy-backend/internal/common"
)

type Service struct {
	Repo   Repo
	Logger *zap.Logger
	Files  *common.FileService
}

func NewService(repo Repo, logger *zap.Logger, files *common.FileService) *Service {
	return &Service{
		Repo:   repo,
		Logger: logger,
		Files:  files,
	}
}

func (s *Service) GetBanners(ctx context.Context) ([]Banner, error) {
	banners, err := s.Repo.GetBanners(ctx)
	if err != nil {
		s.Logger.Error("Failed to get banners", zap.Error(err))
		return nil, err
	}
	return banners, nil
}

func (s *Service) CreateBanner(ctx context.Context, name string, image *multipart.FileHeader, index int) (Banner, error) {
	url, err := s.Files.SaveFile(image)
	if err != nil {
		s.Logger.Error("Failed to upload image", zap.Error(err))
		return Banner{}, common.NewErrorf(FailedToUploadImage, "Failed to upload image: %v", err)
	}

	s.Logger.Info("Image uploaded successfully", zap.String("url", url))

	banner := Banner{
		Name:     name,
		ImageUrl: url,
		Index:    index,
	}
	banner, err = s.Repo.CreateBanner(ctx, banner)
	if err != nil {
		s.Logger.Error("Failed to create banner", zap.Error(err))
		return Banner{}, err
	}

	return banner, nil
}
