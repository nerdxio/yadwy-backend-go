package application

import (
	"context"
	"time"
	"yadwy-backend/internal/common"
	"yadwy-backend/internal/users/domain/contracts"
	"yadwy-backend/internal/users/domain/modles"
)

type UserService struct {
	userRepo     contracts.UserRepo
	customerRepo contracts.CustomerRepo
	sellerRepo   contracts.SellerRepo
	adminRepo    contracts.AdminRepo
	jwt          *common.JWTGenerator
}

func NewUserService(
	repo contracts.UserRepo,
	//customerRepo contracts.CustomerRepo,
	sellerRepo contracts.SellerRepo,
	//adminRepo contracts.AdminRepo,
	jwt *common.JWTGenerator) *UserService {
	return &UserService{
		userRepo: repo,
		//customerRepo: customerRepo,
		sellerRepo: sellerRepo,
		//adminRepo:    adminRepo,
		jwt: jwt,
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

	err = s.assignRole(ctx, savedUser)
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
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User: UserInfo{
			ID:    gu.ID(),
			Name:  gu.Name(),
			Email: gu.Email(),
			Role:  gu.Role().String(),
		},
	}
	return res, nil
}

func (s *UserService) assignRole(ctx context.Context, u *modles.User) error {
	role := u.Role().String()

	switch role {
	case "SELLER":
		seller := modles.NewSeller(0, int64(u.ID()))
		_, err := s.sellerRepo.CreateSeller(ctx, seller)
		if err != nil {
			return err
		}
	case "CUSTOMER":
	// todo: create customer
	case "ADMIN":
		//todo: create admin
	}

	return nil
}
