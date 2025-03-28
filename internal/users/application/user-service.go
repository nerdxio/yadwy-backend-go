package application

import (
	"yadwy-backend/internal/users/domain/contracts"
)

type UserService struct {
	repo contracts.UserRepo
}

func NewUserService(repo contracts.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(name, email, password, role string) (int, error) {
	id, err := s.repo.CreateUser(name, email, password, role)
	if err != nil {
		return 0, err
	}

	return id, nil
}
