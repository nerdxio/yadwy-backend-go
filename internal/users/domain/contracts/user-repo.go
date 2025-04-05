package contracts

import (
	"context"
	"yadwy-backend/internal/users/domain/modles"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *modles.User) (*modles.User, error)
	ListUsers() ([]modles.User, error)
	GetUser(ctx context.Context, email string) (*modles.User, error)
	UserExists(ctx context.Context, email string) (bool, error)
}
