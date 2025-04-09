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

func (s *UserService) CreateUser(ctx context.Context, r CreateUserReq) (*LoginUserRes, error) {
	hashPass, err := common.HashPass(r.Password)
	if err != nil {
		return nil, common.NewErrorf(modles.InvalidUserCredentialsError, "Invalid user credentials")
	}

	role, err := modles.NewRole(r.Role)
	if err != nil {
		return nil, common.NewErrorf(modles.InvalidUserRoleError, "Invalid user role")
	}

	b, err := s.userRepo.UserExists(ctx, r.Email)
	if err != nil || b {
		return nil, common.NewErrorf(modles.UserAlreadyExistsError, "user already exists")
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
		return nil, err
	}
	return generateTokens(savedUser, s.jwt)
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

	return generateTokens(gu, s.jwt)
}

func generateTokens(user *modles.User, jwt *common.JWTGenerator) (*LoginUserRes, error) {
	accessToken, accessClaims, err := jwt.CreateToken(int64(user.ID()), user.Email(), user.Role().String(), 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshClaims, err := jwt.CreateToken(int64(user.ID()), user.Email(), user.Role().String(), 15*time.Minute)
	if err != nil {
		return nil, err
	}

	res := &LoginUserRes{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		User: UserInfo{
			ID:    user.ID(),
			Name:  user.Name(),
			Email: user.Email(),
			Role:  user.Role().String(),
		},
	}
	return res, nil
}
