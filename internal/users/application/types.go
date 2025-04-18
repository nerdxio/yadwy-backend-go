package application

import "time"

// CreateUserReq represents the request payload for creating a new user
// @Description User registration request payload
type CreateUserReq struct {
	Name     string `json:"name" validate:"required" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required" example:"strongpassword123"`
	Role     string `json:"role" validate:"required" role:"CUSTOMER,ADMIN,SELLER" example:"CUSTOMER"`
}

// CreateUserRes represents the response for user creation
// @Description User registration response payload
type CreateUserRes struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
	Role  string `json:"role" example:"CUSTOMER"`
}

// LoginUserReq represents the login request payload
// @Description User login request payload
type LoginUserReq struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required" example:"strongpassword123"`
}

// LoginUserRes represents the login response
// @Description User login response payload
type LoginUserRes struct {
	AccessToken           string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken          string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UserInfo  `json:"user"`
}

// UserInfo represents basic user information
// @Description Basic user information
type UserInfo struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
	Role  string `json:"role" example:"CUSTOMER"`
}
