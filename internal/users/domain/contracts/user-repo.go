package contracts

import "yadwy-backend/internal/users/domain/modles"

type UserRepo interface {
	CreateUser(name, email, password string, role string) (int, error)
	ListUsers() ([]modles.User, error)
}
