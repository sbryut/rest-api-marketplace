package service

import (
	"context"
	"fmt"
	"rest-api-marketplace/internal/repository"
	"rest-api-marketplace/pkg/auth"
	"rest-api-marketplace/pkg/hash"
	"time"

	"rest-api-marketplace/internal/entity"
)

type UsersService struct {
	repo            repository.Users
	hasher          hash.PasswordHasher
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repo repository.Users, tokenManager auth.TokenManager, tokenTTL, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		tokenManager:    tokenManager,
		accessTokenTTL:  tokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input UserInput) (*entity.User, error) {
	// TODO: сделать валидацию лучше
	if len(input.Login) < 3 || len(input.Login) > 30 {
		return nil, fmt.Errorf("login must be between 3 and 50 characters")
	}
	if len(input.Password) < 6 {
		return nil, fmt.Errorf("login must be at least 6 characters")
	}

	existingUser, err := s.repo.GetByLogin(ctx, input.Login)
	if err != nil {
		return nil, fmt.Errorf("UsersService layer error: %w", err)
	}
	if existingUser != nil {
		return nil, entity.ErrUserExists
	}

	hashedPass, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := entity.User{
		Login:        input.Login,
		PasswordHash: hashedPass,
		CreatedAt:    time.Now(),
	}

	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	createdUser, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created user: %w", err)
	}
	return createdUser, nil
}

func (s *UsersService) SignIn(ctx context.Context, input UserInput) (string, error) {
	user, err := s.repo.GetByLogin(ctx, input.Login)
	if err != nil {
		return "", fmt.Errorf("UsersService layer error: %w", err)
	}
	if user == nil {
		return "", entity.ErrUserNotFound
	}

	if !s.hasher.Check(input.Password, user.PasswordHash) {
		return "", entity.ErrInvalidCreds
	}

	token, err := s.tokenManager.NewJWTToken(user.ID, s.accessTokenTTL)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
	// TODO либо тут сделать возврат (Tokens{}, error), хз
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	return s.createSession(ctx, user.ID)
}

func (s *UsersService) createSession(ctx context.Context, id int64) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWTToken(id, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := entity.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, id, session)

	return res, nil
}
