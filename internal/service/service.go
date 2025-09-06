// Package service provides business logic
package service

import (
	"context"
	"log/slog"
	"time"

	"rest-api-marketplace/internal/entity"
	"rest-api-marketplace/internal/repository"
	"rest-api-marketplace/pkg/auth"
	"rest-api-marketplace/pkg/hash"
)

// UserInput represents user credentials input
type UserInput struct {
	Login    string
	Password string
}

// Tokens contains access and refresh JWT tokens
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// CreateAdInput is used to create a new ad
type CreateAdInput struct {
	Title       string
	Description string
	ImageURL    string
	Price       float64
}

// UpdateAdInput is used to update an existing ad
type UpdateAdInput struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	Price       *float64 `json:"price,omitempty"`
}

// Users defines the interface for user-related operations
type Users interface {
	SignUp(ctx context.Context, input UserInput) (*entity.User, error)
	SignIn(ctx context.Context, input UserInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	createSession(ctx context.Context, id int64) (Tokens, error)
}

// Ads defines the interface for ad-related operations
type Ads interface {
	Create(ctx context.Context, input CreateAdInput, userID int64) (*entity.Ad, error)
	Update(ctx context.Context, adID, userID int64, input UpdateAdInput) (*entity.Ad, error)
	GetByID(ctx context.Context, id int64) (*entity.Ad, error)
	GetByIDWithAuthor(ctx context.Context, id int64, currentUserID *int64) (*entity.AdResponse, error)
	GetAll(ctx context.Context, params entity.GetAdsQuery, currentUserID *int64) ([]entity.AdResponse, error)
	Delete(ctx context.Context, adID, userID int64) error
}

// Services aggregates all service implementations
type Services struct {
	Users Users
	Ads   Ads
}

// Deps contains dependencies required to initialize services
type Deps struct {
	Logger          *slog.Logger
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// NewServices initializes all services with dependencies
func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users, deps.Logger, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)
	adsService := NewAdService(deps.Repos.Ads, deps.Logger)
	return &Services{
		Users: usersService,
		Ads:   adsService,
	}
}
