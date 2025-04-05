package contracts

import (
	"context"
	"yadwy-backend/internal/users/domain/modles"
)

type AdminRepo interface {
	CreateAdmin(ctx context.Context, admin *modles.Admin) (int, error)
	GetAdmin(ctx context.Context, id int64) (*modles.Admin, error)
}
