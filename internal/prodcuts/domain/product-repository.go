package domain

import "context"

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product, images []Image) error
	GetProduct(ctx context.Context, id int64) (*Product, error)
}
