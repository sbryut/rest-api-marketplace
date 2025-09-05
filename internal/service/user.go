package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"rest-api-marketplace/internal/repository"
	"rest-api-marketplace/pkg/auth"
	"rest-api-marketplace/pkg/hash"
	"time"

	"rest-api-marketplace/internal/entity"
)

type UsersService struct {
	repo            repository.Users
	logger          *slog.Logger
	hasher          hash.PasswordHasher
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repo repository.Users, logger *slog.Logger, hasher hash.PasswordHasher, tokenManager auth.TokenManager, tokenTTL, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		logger:          logger,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  tokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input UserInput) (*entity.User, error) {
	const op = "service.UsersService.SignUp"

	if len(input.Login) < 3 || len(input.Login) > 30 {
		return nil, fmt.Errorf("%s: %w: login length", op, entity.ErrInvalidInput)
	}
	if len(input.Password) < 6 {
		return nil, fmt.Errorf("%s: %w: password too short", op, entity.ErrInvalidInput)
	}

	hashedPass, err := s.hasher.Hash(input.Password)
	if err != nil {
		s.logger.Error("failed to hash password", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user := entity.User{
		Login:        input.Login,
		PasswordHash: hashedPass,
		CreatedAt:    time.Now(),
	}

	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		s.logger.Error("failed to create user", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createdUser, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to retrieve created user", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return createdUser, nil
}

func (s *UsersService) SignIn(ctx context.Context, input UserInput) (Tokens, error) {
	const op = "service.UsersService.SignIn"

	user, err := s.repo.GetByLogin(ctx, input.Login)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return Tokens{}, fmt.Errorf("%s: %w", op, entity.ErrInvalidCreds)
		}
		s.logger.Error("failed to get user by login", slog.String("op", op), slog.String("error", err.Error()))
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	if !s.hasher.Check(input.Password, user.PasswordHash) {
		return Tokens{}, entity.ErrInvalidCreds
	}

	return s.createSession(ctx, user.ID)
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	const op = "service.UsersService.RefreshTokens"

	if refreshToken == "" {
		return Tokens{}, fmt.Errorf("%s: %w: empty refresh token", op, entity.ErrInvalidInput)
	}

	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return Tokens{}, fmt.Errorf("%s: %w", op, err)
		}
		s.logger.Error("failed to get user by refresh token", slog.String("op", op), slog.String("error", err.Error()))
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}
	return s.createSession(ctx, user.ID)
}

func (s *UsersService) createSession(ctx context.Context, id int64) (Tokens, error) {
	const op = "service.UsersService.createSession"
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWTToken(id, s.accessTokenTTL)
	if err != nil {
		s.logger.Error("failed to create access token", slog.String("op", op), slog.String("error", err.Error()))
		return res, fmt.Errorf("%s: %w", op, err)
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		s.logger.Error("failed to create refresh token", slog.String("op", op), slog.String("error", err.Error()))
		return res, fmt.Errorf("%s: %w", op, err)
	}

	session := entity.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, id, session)
	if err != nil {
		s.logger.Error("failed to set session", slog.String("op", op), slog.String("error", err.Error()))
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
