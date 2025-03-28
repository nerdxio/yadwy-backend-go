package db

import (
	"context"
	"database/sql"
	"errors"
	"yadwy-backend/internal/users/domain/modles"

	"github.com/jmoiron/sqlx"
)

type UserEntity struct {
	UserID   int    `db:"user_id"`
	Username string `db:"username"`
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

func (r *UserRepo) CreateUser(name, email, password, roleStr string) (int, error) {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING user_id
	`

	var userID int
	err := r.db.QueryRowContext(context.Background(), query, name, email, password).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepo) ListUsers() ([]modles.User, error) {
	query := `
		SELECT user_id, username, email, password
		FROM users
	`

	var entities []UserEntity
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
		SELECT user_id, username, email, password
		FROM users
		WHERE user_id = $1
	`

	var entity UserEntity
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

// Helper functions

// mapEntityToDomain maps a database entity to a domain model
func mapEntityToDomain(entity UserEntity, role modles.Role) (*modles.User, error) {
	user, err := modles.NewUser(entity.Username, entity.Email, entity.Password, role)
	if err != nil {
		return nil, err
	}

	// Since the domain model doesn't expose a setter for ID, we need to create a new user with the ID
	// This is a bit of a hack, but it's necessary to maintain the domain model's encapsulation
	// In a real-world scenario, you might want to add a method to set the ID in the domain model
	return user, nil
}
