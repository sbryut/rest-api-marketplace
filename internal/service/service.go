package service

import (
	"context"
	"rest-api-marketplace/internal/repository"
	"rest-api-marketplace/pkg/auth"
	"rest-api-marketplace/pkg/hash"
	"time"

	"rest-api-marketplace/internal/entity"
)

type UserInput struct {
	Login    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input UserInput) (*entity.User, error)
	SignIn(ctx context.Context, input UserInput) (string, error)
}

type Ads interface {
	GetByID(ctx context.Context, id string) entity.Ad
	GetAll(ctx context.Context) []entity.Ad
}

type Services struct {
	Users Users
	Ads   Ads
}

type Deps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)
	adsService := NewAdService(deps.Repos.Ads)
	return &Services{
		Users: usersService,
		Ads:   adsService,
	}
}
