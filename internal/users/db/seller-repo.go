package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"yadwy-backend/internal/users/domain/modles"
)

type SellerDbo struct {
	ID     int64 `db:"id"`
	UserID int64 `db:"user_id"`
}

type SellerRepo struct {
	db *sqlx.DB
}

func NewSellerRepo(db *sqlx.DB) *SellerRepo {
	return &SellerRepo{
		db: db,
	}
}

func (r *SellerRepo) CreateSeller(ctx context.Context, seller *modles.Seller) (int64, error) {
	query := `
		INSERT INTO sellers (user_id)
		VALUES ($1)
		RETURNING id
	`
	var sellerID int64
	err := r.db.QueryRowContext(ctx, query, seller.UserId()).Scan(&sellerID)

	if err != nil {
		return 0, err
	}

	return sellerID, nil
}

func (r *SellerRepo) GetSeller(ctx context.Context, id int64) (*modles.Seller, error) {
	query := `
		SELECT id, user_id
		FROM sellers
		WHERE id = $1
	`
	var dbo SellerDbo
	err := r.db.GetContext(ctx, &dbo, query, id)
	if err != nil {
		return nil, err
	}

	return toDomain(dbo), nil
}

func toDomain(dbo SellerDbo) *modles.Seller {
	return modles.NewSeller(
		dbo.ID,
		dbo.UserID,
	)
}
