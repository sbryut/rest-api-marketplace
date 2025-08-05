package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rest-api-marketplace/internal/entity"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) Create(ctx context.Context, user entity.User) (int64, error) {
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Login, user.PasswordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create new user: %w", err)
	}
	return id, nil
}

func (r *UsersRepo) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return &user, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `SELECT id, login, created_at FROM users WHERE id = $1 `
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Login, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error) {
	query := `SELECT id, login, password_hash, created_at
			  FROM users
			  WHERE refresh_token = $1 AND refresh_expires_at > NOW()`

	var user entity.User
	err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by refresh token: %w", err)
	}
	return &user, nil
}

func (r *UsersRepo) SetSession(ctx context.Context, id int64, session entity.Session) error {
	query := `UPDATE users SET refresh_token = $1, refresh_expires_at = $2, last_visit_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, session.RefreshToken, session.ExpiresAt, id)
	return err
}
