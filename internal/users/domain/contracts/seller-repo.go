package contracts

import (
	"context"
	"yadwy-backend/internal/users/domain/modles"
)

type SellerRepo interface {
	CreateSeller(ctx context.Context, seller *modles.Seller) (int64, error)
	GetSeller(ctx context.Context, id int64) (*modles.Seller, error)
}
