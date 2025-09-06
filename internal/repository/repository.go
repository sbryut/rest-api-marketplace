// Package repository aggregates DB repositories
package repository

import (
	"context"
	"database/sql"

	"rest-api-marketplace/internal/entity"
)

// Users defines user repository interface
type Users interface {
	Create(ctx context.Context, user entity.User) (int64, error)
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error)
	SetSession(ctx context.Context, id int64, session entity.Session) error
}

// Ads defines ad repository interface
type Ads interface {
	Create(ctx context.Context, ad entity.Ad) (int64, error)
	Update(ctx context.Context, id int64, ad entity.Ad) error
	GetByID(ctx context.Context, id int64) (*entity.Ad, error)
	GetByIDWithAuthor(ctx context.Context, id int64) (*entity.AdWithAuthor, error)
	GetAll(ctx context.Context, params entity.GetAdsQuery) ([]entity.AdWithAuthor, error)
	Delete(ctx context.Context, id int64) error
}

// Repositories aggregates all repositories
type Repositories struct {
	Users Users
	Ads   Ads
}

// NewRepositories initializes all repositories
func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
		Ads:   NewAdsRepo(db),
	}
}
