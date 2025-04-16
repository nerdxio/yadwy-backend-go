package domain

import "context"

//go:generate mockgen -destination=mock/cart_repository_mock.go -package=mock yadwy-backend/internal/cart/domain CartRepository

type CartRepository interface {
	CreateCart(ctx context.Context, userID int64) (*Cart, error)
	GetCart(ctx context.Context, userID int64) (*Cart, error)
	AddItem(ctx context.Context, userID int64, productID int64, quantity int) error
	UpdateItem(ctx context.Context, userID int64, productID int64, quantity int) error
	RemoveItem(ctx context.Context, userID int64, productID int64) error
	ClearCart(ctx context.Context, userID int64) error
}
