package common

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	secretKey string
}

func NewJWTGenerator(secretKey string) *JWTGenerator {
	return &JWTGenerator{secretKey}
}

// UserClaims represents the custom claims for the JWT token
// @Description JWT token claims containing user information
type UserClaims struct {
	ID                   int64            `json:"id" example:"1"`
	Email                string           `json:"email" example:"john@example.com"`
	Role                 string           `json:"role" example:"CUSTOMER"`
	TokenID              string           `json:"jti,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Subject              string           `json:"sub,omitempty" example:"john@example.com"`
	IssuedAt             *jwt.NumericDate `json:"iat,omitempty" swaggertype:"primitive,integer" example:"1618317375"`
	ExpiresAt            *jwt.NumericDate `json:"exp,omitempty" swaggertype:"primitive,integer" example:"1618317975"`
	jwt.RegisteredClaims `json:"-"`
}

func NewUserClaims(id int64, email, role string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token ID: %w", err)
	}

	return &UserClaims{
		Email: email,
		ID:    id,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

func (maker *JWTGenerator) CreateToken(id int64, email, role string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, role, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error signing token: %w", err)
	}

	return tokenStr, claims, nil
}

func (maker *JWTGenerator) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// verify the signing method
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
