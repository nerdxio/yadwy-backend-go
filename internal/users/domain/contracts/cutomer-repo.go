package contracts

import (
	"context"
	"yadwy-backend/internal/users/domain/modles"
)

type CustomerRepo interface {
	CreateCustomer(ctx context.Context, customer *modles.Customer) (int, error)
	GetCustomer(ctx context.Context, id int64) (*modles.Customer, error)
}
