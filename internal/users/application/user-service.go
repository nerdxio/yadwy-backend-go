package application

import (
	"context"
	"time"
	"yadwy-backend/internal/common"
	"yadwy-backend/internal/users/domain/contracts"
	"yadwy-backend/internal/users/domain/modles"
)

type UserService struct {
	userRepo contracts.UserRepo
	jwt      *common.JWTGenerator
}

func NewUserService(
	repo contracts.UserRepo,
	jwt *common.JWTGenerator) *UserService {
	return &UserService{
		userRepo: repo,
		jwt:      jwt,
	}
}

func (s *UserService) CreateUser(ctx context.Context, r CreateUserReq) (int, error) {

	hashPass, err := common.HashPass(r.Password)
	if err != nil {
		return 0, err
	}

	role, err := modles.NewRole(r.Role)

	if err != nil {
		return 0, err
	}

	_, err = s.userRepo.UserExists(ctx, r.Email)
	if err != nil {
		return 0, err
	}

	user := modles.NewUser(
		0,
		r.Name,
		r.Email,
		hashPass,
		role,
	)

	savedUser, err := s.userRepo.CreateUser(ctx, user)

	if err != nil {
		return 0, err
	}
	return savedUser.ID(), nil
}

func (s *UserService) LoginUser(ctx context.Context, req LoginUserReq) (*LoginUserRes, error) {
	gu, err := s.userRepo.GetUser(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	err = common.CheckPassword(gu.Password(), req.Password)
	if err != nil {
		return nil, err
	}

	accessToken, accessClaims, err := s.jwt.CreateToken(int64(gu.ID()), gu.Email(), "ADMIN", 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshClaims, err := s.jwt.CreateToken(int64(gu.ID()), gu.Email(), "ADMIN", 15*time.Minute)
	if err != nil {
		return nil, err
	}

	res := &LoginUserRes{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		User: UserInfo{
			ID:    gu.ID(),
			Name:  gu.Name(),
			Email: gu.Email(),
			Role:  gu.Role().String(),
		},
	}
	return res, nil
}
