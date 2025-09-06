package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rest-api-marketplace/internal/entity"

	"github.com/lib/pq"
)

// UsersRepo provides DB operations for users
type UsersRepo struct {
	db *sql.DB
}

// NewUsersRepo creates a new UsersRepo instance
func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

// Create inserts a new user and returns its ID
func (r *UsersRepo) Create(ctx context.Context, user entity.User) (int64, error) {
	const op = "repository.UsersRepo.Create"

	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Login, user.PasswordHash).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, entity.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

// GetByLogin retrieves a user by login
func (r *UsersRepo) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	const op = "repository.UsersRepo.GetByLogin"

	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`

	var user entity.User
	err := r.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, entity.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UsersRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	const op = "repository.UsersRepo.GetByID"

	query := `SELECT id, login, created_at FROM users WHERE id = $1 `

	var user entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Login, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, entity.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

// GetByRefreshToken retrieves a user by refresh token
func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error) {
	const op = "repository.UsersRepo.GetByRefreshToken"

	query := `SELECT id, login, password_hash, created_at
			  FROM users
			  WHERE refresh_token = $1 AND refresh_expires_at > NOW()`

	var user entity.User
	err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, entity.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

// SetSession updates a user's refresh token and expiration
func (r *UsersRepo) SetSession(ctx context.Context, id int64, session entity.Session) error {
	const op = "repository.UsersRepo.SetSession"

	query := `UPDATE users SET refresh_token = $1, refresh_expires_at = $2, last_visit_at = NOW() WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, session.RefreshToken, session.ExpiresAt, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
