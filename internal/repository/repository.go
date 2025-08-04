package repository

import (
	"context"
	"database/sql"
	"rest-api-marketplace/internal/entity"
)

type Users interface {
	Create(ctx context.Context, user entity.User) (int64, error)
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error)
	SetSession(ctx context.Context, id int64, session entity.Session) error
}

type Ads interface {
	Create(ctx context.Context, ad entity.Ad) (int64, error)
	Update(ctx context.Context, id int64, ad entity.Ad) error
	GetById(ctx context.Context, id int64) (*entity.Ad, error)
	GetAll(ctx context.Context, params entity.GetAdsQuery) ([]entity.AdWithAuthor, error)
	Delete(ctx context.Context, id int64) error
}

type Repositories struct {
	Users Users
	Ads   Ads
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
		Ads:   NewAdsRepo(db),
	}
}
