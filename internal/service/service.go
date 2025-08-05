package service

import (
	"context"
	"time"

	"rest-api-marketplace/internal/entity"
	"rest-api-marketplace/internal/repository"
	"rest-api-marketplace/pkg/auth"
	"rest-api-marketplace/pkg/hash"
)

type UserInput struct {
	Login    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type CreateAdInput struct {
	Title       string
	Description string
	ImageURL    string
	Price       float64
}

type UpdateAdInput struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	ImageURl    *string  `json:"image_url,omitempty"`
	Price       *float64 `json:"price,omitempty"`
}

type Users interface {
	SignUp(ctx context.Context, input UserInput) (*entity.User, error)
	SignIn(ctx context.Context, input UserInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	createSession(ctx context.Context, id int64) (Tokens, error)
}

type Ads interface {
	Create(ctx context.Context, input CreateAdInput, userId int64) (*entity.Ad, error)
	Update(ctx context.Context, adId, userId int64, input UpdateAdInput) (*entity.Ad, error)
	GetByID(ctx context.Context, id int64) (*entity.Ad, error)
	GetAll(ctx context.Context, params entity.GetAdsQuery, currentUserId *int64) ([]entity.AdResponse, error)
	Delete(ctx context.Context, adId, userId int64) error
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
