package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yadwy-backend/internal/users/domain/modles"

	"github.com/jmoiron/sqlx"
)

type UserDbo struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *modles.User) (int, error) {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var userID int
	err := r.db.QueryRowContext(ctx, query, user.Name(), user.Email(), user.Password()).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepo) GetUser(ctx context.Context, email string) (*modles.User, error) {
	var u UserDbo
	err := r.db.GetContext(ctx, &u, "SELECT id,name,email,password FROM users WHERE email = $1", email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return mapEntityToDomain(u, modles.Role(u.Role))
}

func (r *UserRepo) ListUsers() ([]modles.User, error) {
	query := `
		SELECT id, name, email, password
		FROM users
	`

	var entities []UserDbo
	err := r.db.SelectContext(context.Background(), &entities, query)
	if err != nil {
		return nil, err
	}

	// Map database entities to domain models
	users := make([]modles.User, 0, len(entities))
	for _, entity := range entities {
		role := modles.Role(entity.Role)
		if !role.IsValid() {
			return nil, errors.New("invalid role found in database: " + entity.Role)
		}

		user, err := mapEntityToDomain(entity, role)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepo) GetUserByID(id int) (*modles.User, error) {
	query := `
		SELECT id, name, email, password
		FROM users
		WHERE id = $1
	`

	var entity UserDbo
	err := r.db.GetContext(context.Background(), &entity, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	role := modles.Role(entity.Role)
	if !role.IsValid() {
		return nil, errors.New("invalid role found in database: " + entity.Role)
	}

	return mapEntityToDomain(entity, role)
}

func mapEntityToDomain(dbo UserDbo, role modles.Role) (*modles.User, error) {
	user := modles.NewUser(dbo.ID, dbo.Name, dbo.Email, dbo.Password, role)
	return user, nil
}
