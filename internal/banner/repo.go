package banner

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
	"yadwy-backend/internal/common"
)

type Repo interface {
	GetBanners(ctx context.Context) ([]Banner, error)
	CreateBanner(ctx context.Context, b Banner) (Banner, error)
}

type RepoImpl struct {
	db     *sqlx.DB
	logger *zap.Logger
}

type BannerDbo struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	ImageUrl  string    `db:"image_url"`
	Index     int       `db:"index"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewRepo(db *sqlx.DB, logger *zap.Logger) *RepoImpl {
	return &RepoImpl{
		db:     db,
		logger: logger,
	}
}

func (r *RepoImpl) GetBanners(ctx context.Context) ([]Banner, error) {
	var dbos []BannerDbo
	err := r.db.SelectContext(ctx, &dbos, "SELECT * FROM banners")
	if err != nil {
		r.logger.Error("Failed to get banners", zap.Error(err))
		return nil, common.NewErrorf(FailedToGetAllBanners, "Failed to get all banners: %v", err)
	}
	return mapToBanners(dbos), nil
}

func (r *RepoImpl) CreateBanner(ctx context.Context, b Banner) (Banner, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO banners (name, image_url, index) VALUES ($1, $2, $3) RETURNING id",
		b.Name, b.ImageUrl, b.Index,
	).Scan(&id)
	if err != nil {
		r.logger.Error("Failed to create banner", zap.Error(err))
		return Banner{}, common.NewErrorf(FailedToCreateBanner, "Failed to create banner: %v", err)
	}
	b.Id = id
	return b, nil
}

func mapToBanners(dbos []BannerDbo) []Banner {
	result := make([]Banner, 0, len(dbos))
	for _, dbo := range dbos {
		result = append(result, mapToBanner(dbo))
	}
	return result
}

func mapToBanner(dbo BannerDbo) Banner {
	return Banner{
		Id:        dbo.ID,
		Name:      dbo.Name,
		ImageUrl:  dbo.ImageUrl,
		Index:     dbo.Index,
		CreatedAt: dbo.CreatedAt,
		UpdatedAt: dbo.UpdatedAt,
	}
}
