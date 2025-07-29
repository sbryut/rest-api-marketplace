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
	Create(ad entity.Ad) entity.Ad
	GetOne(id string) entity.Ad
	GetAll(ctx context.Context) []entity.Ad
	Delete(ad entity.Ad) error
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
